//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/config"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/transport/auth_api"
	"github.com/ruf-dev/artel/internal/transport/mcp_api"
	"github.com/ruf-dev/artel/internal/transport/mcp_keys_api"
	"github.com/ruf-dev/artel/internal/transport/notes_api"
	"github.com/ruf-dev/artel/internal/transport/skills_api"
	"github.com/ruf-dev/artel/internal/transport/vaults_api"
	"github.com/ruf-dev/artel/tests/harness"
)

// skillCreatorSlug is the always-present, non-vault-backed system skill's slug (see
// internal/service/v1/skills.SkillCreatorSlug — kept as a local literal since that constant lives
// in an internal package this suite doesn't import).
const skillCreatorSlug = "skill-creator"

// createdSlugRe extracts the slug create_skill's MCP tool reports back in its confirmation text
// ("Skill %q created with slug %q." — internal/service/v1/mcp/skills_tools.go's createSkillTool).
var createdSlugRe = regexp.MustCompile(`with slug "([^"]+)"`)

// SkillsSuite exercises the skills feature end to end: skills stored as CouchDB notes under a
// vault's reserved .skills/ folder, the always-present system "skill-creator" skill, both the
// SkillsAPI gRPC surface (dashboard-facing CRUD) and the MCP tool surface (static skill_*
// management tools plus dynamic per-hot-plug-skill skill_<slug> tools) an agent client drives,
// and skill-count quota enforcement mirroring QuotaSuite's storage-quota pattern.
type SkillsSuite struct {
	suite.Suite
	db              *sql.DB
	repos           *repopg.Repos
	svcs            *svcv1.Services
	mcpHdlr         *mcp_api.McpHandler
	couchInstanceID string

	authClient   pb.AuthAPIClient
	vaultClient  pb.VaultsAPIClient
	mcpKeyClient pb.McpKeysAPIClient
	notesClient  pb.NotesAPIClient
	skillsClient pb.SkillsAPIClient
}

func TestSkills(t *testing.T) {
	suite.Run(t, new(SkillsSuite))
}

func (s *SkillsSuite) SetupSuite() {
	ctx := context.Background()

	s.db = harness.OpenPostgres(s.T())

	// SubscriptionsEnabled must be true here: with it false (the E2ESuite default) every quota
	// check runs through FreeService, which always passes and never touches CouchDB — the skill
	// limit scenarios below need the real PaidService.CheckSkillLimit path.
	cfg := config.EnvironmentConfig{}
	cfg.SubscriptionsEnabled = true

	var credsEncrypted bool
	s.repos, s.svcs, credsEncrypted = harness.BuildServices(s.T(), s.db, cfg)

	couchURL := harness.CouchURL(s.T())

	s.couchInstanceID = harness.GetCouchInstance(s.T(), ctx, s.svcs, couchURL)

	s.mcpHdlr = mcp_api.NewMcpHandler(s.svcs.Mcp, s.svcs.Mom)

	s.startGrpcServer(credsEncrypted)
}

// startGrpcServer builds a real *grpc.Server chained with the production auth interceptor,
// registers the AuthAPI, VaultsAPI, McpKeysAPI, NotesAPI and SkillsAPI implementations onto it,
// and serves it over an in-memory bufconn listener so the suite's RPCs travel through the real
// transport + auth stack without binding a TCP port.
func (s *SkillsSuite) startGrpcServer(credsEncrypted bool) {
	authImpl := auth_api.NewAuthImpl(
		s.svcs.Auth, "", s.svcs.S3Instance, s.svcs.CouchInstance,
		false, credsEncrypted, s.svcs.DockerHost, s.svcs.SetupWizard,
	)
	workbenchTerminalShellHandler := vaults_api.NewWorkbenchTerminalShellHandler(
		s.svcs.Auth, s.repos.VaultMembers(), s.svcs.Workbench,
	)
	vaultsImpl := vaults_api.NewVaultsImpl(s.svcs.Vault, s.svcs.Workbench, workbenchTerminalShellHandler)
	mcpKeysImpl := mcp_keys_api.NewMcpKeysImpl(s.svcs.Mcp, s.svcs.Mom)
	notesImpl := notes_api.NewNotesImpl(s.svcs.Notes)
	skillsImpl := skills_api.NewSkillsImpl(s.svcs.SkillsService())

	conn := harness.NewBufconnServer(
		s.T(), s.svcs,
		authImpl.Register, vaultsImpl.Register, mcpKeysImpl.Register, notesImpl.Register, skillsImpl.Register,
	)

	s.authClient = pb.NewAuthAPIClient(conn)
	s.vaultClient = pb.NewVaultsAPIClient(conn)
	s.mcpKeyClient = pb.NewMcpKeysAPIClient(conn)
	s.notesClient = pb.NewNotesAPIClient(conn)
	s.skillsClient = pb.NewSkillsAPIClient(conn)
}

func (s *SkillsSuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err, "close db")
	}
}

type skillsTestUser struct {
	userUuid  uuid.UUID
	authedCtx context.Context
	vaultUuid uuid.UUID
	rawToken  string
}

// setupUser registers a user via the real AuthAPI RPC, creates a vault for them via VaultsAPI, and
// issues an MCP key scoped to it via McpKeysAPI — the minimum a caller needs to drive both the
// SkillsAPI gRPC surface and the MCP tool surface.
func (s *SkillsSuite) setupUser() skillsTestUser {
	ctx := context.Background()

	slug := harness.Slug(s.T())
	email := slug + "@test.local"
	// A fixed password, not slug-derived: bcrypt rejects passwords over 72 bytes, and slug embeds
	// the full (possibly long, subtest-qualified) t.Name() — uniqueness only matters for email.
	password := "test-password-skills"

	registerReq := &pb.Register_Request{Email: email, Password: password}

	registerResp, err := s.authClient.Register(ctx, registerReq)
	s.Require().NoError(err)

	userUuid, err := uuid.Parse(registerResp.Id)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), userUuid)
	})

	// UserName must be non-empty: vault creation derives the CouchDB database name from it, and
	// an empty UserName produces a name starting with "-", which CouchDB rejects. Email/password
	// registration never populates domain.User.Username, so the slug is written directly onto the
	// persisted user row, since this suite authenticates through the real gRPC auth interceptor
	// (which reads UserName off that row).
	_, err = s.db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, slug, userUuid)
	s.Require().NoError(err, "set username stand-in for vault creation")

	// With SubscriptionsEnabled=true, middleware.authWithSession calls subscriptionService.
	// CheckActive on every authenticated RPC — so the subscription has to be activated before
	// CreateVault/CreateMcpKey below, not only once overrideSkillLimits runs later in the test.
	// overrideSkillLimits itself re-activates (Upsert is idempotent) once it pins a skill-count
	// bound.
	_, err = s.repos.Subscriptions().Upsert(ctx, userUuid, true)
	s.Require().NoError(err, "activate subscription so authenticated RPCs pass CheckActive")

	passwordCreds := &pb.PasswordCredentials{Email: email, Password: password}
	loginMethod := &pb.Login_Request_Password{Password: passwordCreds}
	loginReq := &pb.Login_Request{Method: loginMethod}

	loginResp, err := s.authClient.Login(ctx, loginReq)
	s.Require().NoError(err)
	s.Require().NotEmpty(loginResp.Token)

	authedCtx := harness.AuthedContext(ctx, loginResp.Token)

	createVaultReq := &pb.CreateVault_Request{Name: slug + "_vault"}

	createVaultResp, err := s.vaultClient.CreateVault(authedCtx, createVaultReq)
	s.Require().NoError(err)

	vaultUuid, err := uuid.Parse(createVaultResp.Id)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Vault.DeleteVault(context.Background(), vaultUuid)
	})

	createKeyReq := &pb.CreateMcpKey_Request{
		VaultId: createVaultResp.Id,
		Name:    slug + "_key",
	}

	createKeyResp, err := s.mcpKeyClient.CreateMcpKey(authedCtx, createKeyReq)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		keyUuid, parseErr := uuid.Parse(createKeyResp.Key.Id)
		if parseErr == nil {
			_ = s.svcs.Mcp.RevokeKey(context.Background(), keyUuid)
		}
	})

	testUser := skillsTestUser{
		userUuid:  userUuid,
		authedCtx: authedCtx,
		vaultUuid: vaultUuid,
		rawToken:  createKeyResp.RawToken,
	}

	return testUser
}

// overrideSkillLimits activates the user's subscription and pins their effective total/hot-plug
// skill caps — a nil bound leaves that cap on the plan's default (20 total / 5 hot-plug from the
// seeded "basic" plan, see migrations/071_skill_limits.sql), so a test can exhaust one cap without
// also tripping the other. Mirrors QuotaSuite.overrideQuota's exact shape, targeting the skill-cap
// override columns instead of the couch/s3 quota ones.
func (s *SkillsSuite) overrideSkillLimits(userUuid uuid.UUID, totalOverride, hotPlugOverride *int) {
	ctx := context.Background()
	subsRepo := s.repos.Subscriptions()

	_, err := subsRepo.Upsert(ctx, userUuid, true)
	s.Require().NoError(err, "activate subscription")

	sub, err := subsRepo.GetByUser(ctx, userUuid)
	s.Require().NoError(err, "get subscription")

	sub.MaxTotalSkillsOverride = totalOverride
	sub.MaxHotPlugSkillsOverride = hotPlugOverride

	_, err = subsRepo.UpsertOverrides(ctx, sub)
	s.Require().NoError(err, "override subscription skill limits")
}

// mcpToolCall drives a tools/call JSON-RPC request over the real MCP HTTP path (mirrors
// QuotaSuite.callWriteFile), returning the raw "result" payload and the RPC-level error (nil on
// success) separately — unlike callWriteFile, these tests need the result payload itself, not
// just success/failure.
func (s *SkillsSuite) mcpToolCall(rawToken, name string, args map[string]any) (json.RawMessage, *mcpRpcError) {
	callParams := map[string]any{
		"name":      name,
		"arguments": args,
	}
	body := mcpCall("tools/call", callParams)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+rawToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.mcpHdlr.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code)

	var rpcResp struct {
		Result json.RawMessage `json:"result"`
		Error  *mcpRpcError    `json:"error"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &rpcResp)
	s.Require().NoError(err)

	return rpcResp.Result, rpcResp.Error
}

// mcpToolCallText calls a tool via mcpToolCall and extracts the first text content block — the
// shape every skill tool result (both static skill_* management tools and dynamic skill_<slug>
// tools) uses (see mcp_api.textResult/toolResultFromExec).
func (s *SkillsSuite) mcpToolCallText(rawToken, name string, args map[string]any) (string, *mcpRpcError) {
	raw, rpcErr := s.mcpToolCall(rawToken, name, args)
	if rpcErr != nil {
		return "", rpcErr
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	err := json.Unmarshal(raw, &result)
	s.Require().NoError(err)
	s.Require().NotEmpty(result.Content, "tool result had no content blocks")

	return result.Content[0].Text, nil
}

// mcpToolsList drives a tools/list JSON-RPC request and returns just the tool names, for
// assertions on which tools are/aren't currently synthesized.
func (s *SkillsSuite) mcpToolsList(rawToken string) []string {
	emptyParams := map[string]any{}
	body := mcpCall("tools/list", emptyParams)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+rawToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.mcpHdlr.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code)

	var rpcResp struct {
		Result struct {
			Tools []struct {
				Name string `json:"name"`
			} `json:"tools"`
		} `json:"result"`
		Error *mcpRpcError `json:"error"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &rpcResp)
	s.Require().NoError(err)
	s.Require().Nil(rpcResp.Error, "tools/list returned rpc error: %+v", rpcResp.Error)

	names := make([]string, 0, len(rpcResp.Result.Tools))
	for _, t := range rpcResp.Result.Tools {
		names = append(names, t.Name)
	}

	return names
}

// skillListSlugs unmarshals list_skills' JSON-array result text into just the slugs present,
// for membership assertions.
func (s *SkillsSuite) skillListSlugs(rawToken string) []string {
	emptyArgs := map[string]any{}

	text, rpcErr := s.mcpToolCallText(rawToken, "list_skills", emptyArgs)
	s.Require().Nil(rpcErr, "list_skills returned rpc error: %+v", rpcErr)

	var rows []struct {
		Slug string `json:"slug"`
	}

	err := json.Unmarshal([]byte(text), &rows)
	s.Require().NoError(err)

	slugs := make([]string, 0, len(rows))
	for _, r := range rows {
		slugs = append(slugs, r.Slug)
	}

	return slugs
}

// TestCRUDLifecycle drives every SkillsAPI RPC in sequence against a single skill: create with
// hot_plug=false, confirm it (and the always-present system skill-creator skill) show up in
// ListSkills, read it back, update its fields in place, promote it to hot-plug, and finally
// delete it — confirming a GetSkill afterward fails as not-found.
func (s *SkillsSuite) TestCRUDLifecycle() {
	user := s.setupUser()
	vaultId := user.vaultUuid.String()

	createReq := &pb.CreateSkill_Request{
		VaultId:     vaultId,
		Name:        "My First Skill",
		Description: "Does a thing when triggered",
		StorageMode: "none",
		Body:        "body v1",
		HotPlug:     false,
	}

	createResp, err := s.skillsClient.CreateSkill(user.authedCtx, createReq)
	s.Require().NoError(err)
	s.Require().NotNil(createResp.Skill)
	s.False(createResp.Skill.IsHotPlug, "created with hot_plug=false")
	s.False(createResp.Skill.IsSystem)
	s.Equal("none", createResp.Skill.StorageMode)

	slug := createResp.Skill.Slug
	s.Require().NotEmpty(slug)

	listReq := &pb.ListSkills_Request{VaultId: vaultId}

	listResp, err := s.skillsClient.ListSkills(user.authedCtx, listReq)
	s.Require().NoError(err)

	var sawSystemSkill, sawNewSkill bool

	for _, sk := range listResp.Skills {
		switch sk.Slug {
		case skillCreatorSlug:
			sawSystemSkill = true

			s.True(sk.IsSystem, "system skill should report IsSystem")
			s.True(sk.IsHotPlug, "system skill is always hot-plug")
		case slug:
			sawNewSkill = true
		}
	}

	s.True(sawSystemSkill, "ListSkills should always include the system skill-creator skill")
	s.True(sawNewSkill, "ListSkills should include the newly created skill")

	getReq := &pb.GetSkill_Request{VaultId: vaultId, Slug: slug}

	getResp, err := s.skillsClient.GetSkill(user.authedCtx, getReq)
	s.Require().NoError(err)
	s.Equal("body v1", getResp.Body)
	s.Equal("My First Skill", getResp.Skill.Name)

	updateReq := &pb.UpdateSkill_Request{
		VaultId:     vaultId,
		Slug:        slug,
		Name:        "Updated Skill Name",
		Description: "Updated description",
		StorageMode: "freeform_notes",
		Body:        "body v2",
	}

	updateResp, err := s.skillsClient.UpdateSkill(user.authedCtx, updateReq)
	s.Require().NoError(err)
	s.Equal("Updated Skill Name", updateResp.Skill.Name)
	s.Equal(slug, updateResp.Skill.Slug, "update never changes the slug")

	getResp2, err := s.skillsClient.GetSkill(user.authedCtx, getReq)
	s.Require().NoError(err)
	s.Equal("body v2", getResp2.Body, "GetSkill should reflect the update")
	s.Equal("Updated Skill Name", getResp2.Skill.Name)
	s.Equal("freeform_notes", getResp2.Skill.StorageMode)

	hotPlugReq := &pb.SetSkillHotPlug_Request{
		VaultId: vaultId,
		Slug:    slug,
		HotPlug: true,
	}

	hotPlugResp, err := s.skillsClient.SetSkillHotPlug(user.authedCtx, hotPlugReq)
	s.Require().NoError(err)
	s.True(hotPlugResp.Skill.IsHotPlug)

	listResp2, err := s.skillsClient.ListSkills(user.authedCtx, listReq)
	s.Require().NoError(err)

	var sawHotPlugged bool

	for _, sk := range listResp2.Skills {
		if sk.Slug == slug {
			sawHotPlugged = true

			s.True(sk.IsHotPlug, "ListSkills should reflect the hot-plug promotion")
		}
	}

	s.True(sawHotPlugged)

	deleteReq := &pb.DeleteSkill_Request{VaultId: vaultId, Slug: slug}

	_, err = s.skillsClient.DeleteSkill(user.authedCtx, deleteReq)
	s.Require().NoError(err)

	_, err = s.skillsClient.GetSkill(user.authedCtx, getReq)
	s.Require().Error(err, "GetSkill after delete should fail")
	s.Contains(err.Error(), "not found")
}

// TestSlugDedup confirms CreateSkill dedupes slugs derived from the same name: a second skill
// created with an identical name gets a "-2"-suffixed slug rather than colliding with the first.
func (s *SkillsSuite) TestSlugDedup() {
	user := s.setupUser()
	vaultId := user.vaultUuid.String()

	name := "Duplicate Name Skill"

	firstReq := &pb.CreateSkill_Request{
		VaultId: vaultId, Name: name, Description: "first", StorageMode: "none", Body: "b1", HotPlug: false,
	}

	resp1, err := s.skillsClient.CreateSkill(user.authedCtx, firstReq)
	s.Require().NoError(err)

	secondReq := &pb.CreateSkill_Request{
		VaultId: vaultId, Name: name, Description: "second", StorageMode: "none", Body: "b2", HotPlug: false,
	}

	resp2, err := s.skillsClient.CreateSkill(user.authedCtx, secondReq)
	s.Require().NoError(err)

	s.NotEqual(resp1.Skill.Slug, resp2.Skill.Slug, "second create must not collide with the first slug")
	s.Contains(resp2.Skill.Slug, "-2", "dedupeSlug appends -2 on first collision")

	listReq := &pb.ListSkills_Request{VaultId: vaultId}

	listResp, err := s.skillsClient.ListSkills(user.authedCtx, listReq)
	s.Require().NoError(err)

	var slugs []string
	for _, sk := range listResp.Skills {
		slugs = append(slugs, sk.Slug)
	}

	s.Contains(slugs, resp1.Skill.Slug)
	s.Contains(slugs, resp2.Skill.Slug)
}

// TestSystemSkillImmutable confirms every mutating SkillsAPI RPC (Update/Delete/SetSkillHotPlug)
// rejects the "skill-creator" slug with user_errors.SystemSkillNotEditable — it's synthesized
// from a Go constant, never stored in CouchDB, and has nothing to update/move/delete.
func (s *SkillsSuite) TestSystemSkillImmutable() {
	user := s.setupUser()
	vaultId := user.vaultUuid.String()

	const wantMsg = "cannot be modified, hot-plug-toggled, or deleted"

	updateReq := &pb.UpdateSkill_Request{
		VaultId: vaultId, Slug: skillCreatorSlug, Name: "x", Description: "y", StorageMode: "none", Body: "z",
	}

	_, err := s.skillsClient.UpdateSkill(user.authedCtx, updateReq)
	s.Require().Error(err)
	s.Contains(err.Error(), wantMsg)

	deleteReq := &pb.DeleteSkill_Request{VaultId: vaultId, Slug: skillCreatorSlug}

	_, err = s.skillsClient.DeleteSkill(user.authedCtx, deleteReq)
	s.Require().Error(err)
	s.Contains(err.Error(), wantMsg)

	hotPlugReq := &pb.SetSkillHotPlug_Request{
		VaultId: vaultId, Slug: skillCreatorSlug, HotPlug: false,
	}

	_, err = s.skillsClient.SetSkillHotPlug(user.authedCtx, hotPlugReq)
	s.Require().Error(err)
	s.Contains(err.Error(), wantMsg)
}

// TestSkillsFolderHiddenFromNotes confirms the reserved .skills/ folder never leaks into the
// regular NotesAPI.ListNotes surface — skills share the same vault-scoped CouchDB storage as
// notes, but are meant to be invisible there (mirrors the existing "_design/" skip).
func (s *SkillsSuite) TestSkillsFolderHiddenFromNotes() {
	user := s.setupUser()
	vaultId := user.vaultUuid.String()

	createReq := &pb.CreateSkill_Request{
		VaultId: vaultId, Name: "Hidden Test Skill", Description: "d", StorageMode: "none", Body: "b", HotPlug: false,
	}

	_, err := s.skillsClient.CreateSkill(user.authedCtx, createReq)
	s.Require().NoError(err)

	listNotesReq := &pb.ListNotes_Request{VaultId: vaultId}

	notesResp, err := s.notesClient.ListNotes(user.authedCtx, listNotesReq)
	s.Require().NoError(err)

	for _, n := range notesResp.Notes {
		s.False(strings.HasPrefix(n.Path, ".skills/"), "note path %q must not be listed under ListNotes", n.Path)
	}
}

// TestMCPDispatch drives the skill tools exclusively through the real MCP JSON-RPC-over-HTTP
// path (the surface an agent client actually speaks), covering: the six static skill_* builtin
// tools, the always-present dynamic skill_skill-creator tool, and the dynamic skill_<slug> tool
// synthesized for a newly hot-plugged skill (including its disappearance once demoted/deleted).
func (s *SkillsSuite) TestMCPDispatch() {
	user := s.setupUser()

	initialTools := s.mcpToolsList(user.rawToken)
	s.Contains(initialTools, "list_skills")
	s.Contains(initialTools, "get_skill_body")
	s.Contains(initialTools, "create_skill")
	s.Contains(initialTools, "update_skill")
	s.Contains(initialTools, "delete_skill")
	s.Contains(initialTools, "set_skill_hot_plug")
	s.Contains(initialTools, "skill_"+skillCreatorSlug, "the system skill is always hot-plug, so its dynamic tool is always present")

	body := "# MCP Skill\nStep one. Step two."
	createArgs := map[string]any{
		"name":         "MCP Dispatch Skill",
		"description":  "A skill created via the MCP tool surface",
		"storage_mode": "none",
		"body":         body,
		"hot_plug":     true,
	}

	createText, rpcErr := s.mcpToolCallText(user.rawToken, "create_skill", createArgs)
	s.Require().Nil(rpcErr, "create_skill returned rpc error: %+v", rpcErr)

	matches := createdSlugRe.FindStringSubmatch(createText)
	s.Require().Len(matches, 2, "expected to find created slug in create_skill result text: %q", createText)
	slug := matches[1]

	toolsAfterCreate := s.mcpToolsList(user.rawToken)
	s.Contains(toolsAfterCreate, "skill_"+slug, "hot-plugging synthesizes a dynamic skill_<slug> tool")

	// Calling the dynamic skill_<slug> tool with no arguments returns the skill's body verbatim.
	emptyArgs := map[string]any{}

	bodyText, rpcErr := s.mcpToolCallText(user.rawToken, "skill_"+slug, emptyArgs)
	s.Require().Nil(rpcErr, "skill_%s returned rpc error: %+v", slug, rpcErr)
	s.Equal(body, bodyText)

	s.Contains(s.skillListSlugs(user.rawToken), slug, "list_skills should include the new skill")

	getBodyArgs := map[string]any{"slug": slug}

	getBodyText, rpcErr := s.mcpToolCallText(user.rawToken, "get_skill_body", getBodyArgs)
	s.Require().Nil(rpcErr, "get_skill_body returned rpc error: %+v", rpcErr)
	s.Equal(body, getBodyText)

	demoteArgs := map[string]any{"slug": slug, "hot_plug": false}

	_, rpcErr = s.mcpToolCallText(user.rawToken, "set_skill_hot_plug", demoteArgs)
	s.Require().Nil(rpcErr, "set_skill_hot_plug returned rpc error: %+v", rpcErr)

	toolsAfterDemote := s.mcpToolsList(user.rawToken)
	s.NotContains(toolsAfterDemote, "skill_"+slug, "demoting out of hot-plug removes the dynamic tool")

	deleteArgs := map[string]any{"slug": slug}

	_, rpcErr = s.mcpToolCallText(user.rawToken, "delete_skill", deleteArgs)
	s.Require().Nil(rpcErr, "delete_skill returned rpc error: %+v", rpcErr)

	s.NotContains(s.skillListSlugs(user.rawToken), slug, "list_skills should no longer include the deleted skill")
}

// TestTotalSkillLimit overrides a user's total-skill cap down to zero — already-at-quota — and
// confirms CreateSkill (even with hot_plug=false, so only the total cap applies) is rejected
// before anything is written to CouchDB. The system skill-creator skill is synthesized, not
// CouchDB-backed, so it never counts against this live-measured cap — ListSkills still reports
// exactly the one system skill afterward.
func (s *SkillsSuite) TestTotalSkillLimit() {
	user := s.setupUser()
	vaultId := user.vaultUuid.String()

	zero := 0
	s.overrideSkillLimits(user.userUuid, &zero, nil)

	createReq := &pb.CreateSkill_Request{
		VaultId: vaultId, Name: "Over The Cap", Description: "d", StorageMode: "none", Body: "b", HotPlug: false,
	}

	_, err := s.skillsClient.CreateSkill(user.authedCtx, createReq)
	s.Require().Error(err, "create should be rejected once the total skill cap is exhausted")
	s.Contains(err.Error(), "skill limit exceeded")

	listReq := &pb.ListSkills_Request{VaultId: vaultId}

	listResp, err := s.skillsClient.ListSkills(user.authedCtx, listReq)
	s.Require().NoError(err)
	s.Require().Len(listResp.Skills, 1, "only the system skill should be present; nothing was written")
	s.Equal(skillCreatorSlug, listResp.Skills[0].Slug)
}

// TestHotPlugSkillLimit overrides a user's hot-plug-skill cap down to zero (total cap left at the
// plan default) and confirms: a hot_plug=true create is rejected and nothing is written;
// hot_plug=false still succeeds since only the total cap applies there; and promoting that
// library skill afterward via SetSkillHotPlug hits the same hot-plug cap and leaves the skill's
// hot-plug state unchanged.
func (s *SkillsSuite) TestHotPlugSkillLimit() {
	user := s.setupUser()
	vaultId := user.vaultUuid.String()

	zero := 0
	s.overrideSkillLimits(user.userUuid, nil, &zero)

	hotPlugCreateReq := &pb.CreateSkill_Request{
		VaultId: vaultId, Name: "Over The Hot-Plug Cap", Description: "d", StorageMode: "none", Body: "b", HotPlug: true,
	}

	_, err := s.skillsClient.CreateSkill(user.authedCtx, hotPlugCreateReq)
	s.Require().Error(err, "hot-plug create should be rejected once the hot-plug skill cap is exhausted")
	s.Contains(err.Error(), "hot-plug skill limit exceeded")

	listReq := &pb.ListSkills_Request{VaultId: vaultId}

	listResp, err := s.skillsClient.ListSkills(user.authedCtx, listReq)
	s.Require().NoError(err)
	s.Require().Len(listResp.Skills, 1, "only the system skill should be present; nothing was written")
	s.Equal(skillCreatorSlug, listResp.Skills[0].Slug)

	// Only the total cap applies to a library (non-hot-plug) create, and that cap is still at its
	// plan default here, so this should succeed.
	libraryCreateReq := &pb.CreateSkill_Request{
		VaultId: vaultId, Name: "Library Only Skill", Description: "d", StorageMode: "none", Body: "b", HotPlug: false,
	}

	createResp, err := s.skillsClient.CreateSkill(user.authedCtx, libraryCreateReq)
	s.Require().NoError(err, "library-only create should succeed since only the total cap applies")
	slug := createResp.Skill.Slug

	promoteReq := &pb.SetSkillHotPlug_Request{
		VaultId: vaultId, Slug: slug, HotPlug: true,
	}

	_, err = s.skillsClient.SetSkillHotPlug(user.authedCtx, promoteReq)
	s.Require().Error(err, "promoting to hot-plug should hit the same exhausted hot-plug cap")
	s.Contains(err.Error(), "hot-plug skill limit exceeded")

	getReq := &pb.GetSkill_Request{VaultId: vaultId, Slug: slug}

	getResp, err := s.skillsClient.GetSkill(user.authedCtx, getReq)
	s.Require().NoError(err)
	s.False(getResp.Skill.IsHotPlug, "the rejected promotion must not have changed the skill's hot-plug state")
}
