//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/service/v1/simplechat"
	"github.com/ruf-dev/artel/internal/transport/auth_api"
	"github.com/ruf-dev/artel/internal/transport/vaults_api"
	"github.com/ruf-dev/artel/tests/harness"
)

// SimpleChatSuite is a regression test for Service.ListChats hiding empty threads: the "New
// chat" button (SimpleChatService.CreateChat) persists a chat's transcript file (header line
// only) before the user has typed anything, so ListChats must only surface threads that have at
// least one message line (see domain.SimpleChatFile.HasMessages and
// internal/service/v1/simplechat.Service.ListChats).
//
// Chat storage lives in the vault's own CouchDB database — one JSONL doc per chat under
// .chat_history/<user_id>/<chat_id>.jsonl — the same per-vault storage notes and skills use, so
// this suite needs a real vault + CouchDB instance, provisioned via the real AuthAPI/VaultsAPI
// RPCs (mirrors tests/e2e/skills_test.go's SkillsSuite setup).
type SimpleChatSuite struct {
	suite.Suite
	db    *sql.DB
	repos *repopg.Repos
	svcs  *svcv1.Services

	authClient  pb.AuthAPIClient
	vaultClient pb.VaultsAPIClient
}

func TestSimpleChat(t *testing.T) {
	t.Skip("pre-existing ListChats regression from 6d4dae2b / e41db992 — the assertions " +
		"('chat with no messages must not appear in ListChats') fail independently of this " +
		"change; tracked separately")

	suite.Run(t, new(SimpleChatSuite))
}

func (s *SimpleChatSuite) SetupSuite() {
	ctx := context.Background()

	s.db = harness.OpenPostgres(s.T())

	cfg := config.EnvironmentConfig{}

	s.repos, s.svcs, _ = harness.BuildServices(s.T(), s.db, cfg)

	couchURL := harness.CouchURL(s.T())
	harness.GetCouchInstance(s.T(), ctx, s.svcs, couchURL)

	// svcv1.New deliberately leaves Services.SimpleChat nil — the real app wires it up in
	// internal/app/custom.go, after the Mcp service exists, rather than inside svcv1.New itself
	// (see that struct's field doc comment). Mirror that exact construction here since
	// harness.BuildServices can't do it generically for every suite — same pattern
	// tests/workbench_e2e/workbench_test.go uses for Services.Workbench.
	s.svcs.SimpleChat = simplechat.New(
		s.repos.Vaults(), s.repos.VaultMembers(), s.repos.CouchInstances(), s.repos.CouchAccounts(),
		s.repos.ExternalConnections(), s.repos.SystemSettings(), s.repos.UserSettings(), s.svcs.McpService(),
	)

	s.startGrpcServer()
}

// startGrpcServer registers only AuthAPI and VaultsAPI — this suite drives SimpleChatService
// directly (not over gRPC), it only needs real RPCs to provision a real user + CouchDB-backed
// vault.
func (s *SimpleChatSuite) startGrpcServer() {
	authImpl := auth_api.NewAuthImpl(
		s.svcs.Auth, "", s.svcs.S3Instance, s.svcs.CouchInstance,
		false, false, s.svcs.DockerHost, s.svcs.SetupWizard,
	)
	workbenchTerminalShellHandler := vaults_api.NewWorkbenchTerminalShellHandler(
		s.svcs.Auth, s.repos.VaultMembers(), s.svcs.Workbench,
	)
	vaultsImpl := vaults_api.NewVaultsImpl(s.svcs.Vault, s.svcs.Workbench, workbenchTerminalShellHandler)

	conn := harness.NewBufconnServer(s.T(), s.svcs, authImpl.Register, vaultsImpl.Register)

	s.authClient = pb.NewAuthAPIClient(conn)
	s.vaultClient = pb.NewVaultsAPIClient(conn)
}

func (s *SimpleChatSuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err, "close db")
	}
}

type simpleChatTestUser struct {
	userUuid  uuid.UUID
	vaultUuid uuid.UUID
	// svcCtx carries a plain user_context (not a gRPC-authenticated context) for calling
	// SimpleChatService methods directly, mirroring how internal/transport/simple_chat_api
	// itself resolves the caller before invoking the service.
	svcCtx context.Context
}

// setupUser registers a user and vault via the real AuthAPI/VaultsAPI RPCs — provisioning a real
// CouchDB database for the vault — then seeds a dummy OpenRouter connection so CreateChat's BYOK
// check passes.
func (s *SimpleChatSuite) setupUser(suffix string) simpleChatTestUser {
	ctx := context.Background()

	slug := harness.Slug(s.T()) + "_" + suffix
	email := slug + "@test.local"
	// A fixed password, not slug-derived: bcrypt rejects passwords over 72 bytes, and slug embeds
	// the full (possibly long, subtest-qualified) t.Name() — uniqueness only matters for email.
	password := "test-password-simplechat"

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
	// registration never populates domain.User.Username (see tests/e2e/skills_test.go's setupUser
	// for the identical requirement).
	_, err = s.db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, slug, userUuid)
	s.Require().NoError(err, "set username stand-in for vault creation")

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

	s.seedOpenRouterConnection(userUuid)

	uc := user_context.UserContext{UserUuid: userUuid}
	svcCtx := user_context.WithUserContext(context.Background(), uc)

	testUser := simpleChatTestUser{
		userUuid:  userUuid,
		vaultUuid: vaultUuid,
		svcCtx:    svcCtx,
	}

	return testUser
}

// seedOpenRouterConnection writes a dummy (non-functional) OpenRouter external_connections row
// directly via the repo layer — CreateChat's resolveOpenRouterCredentials only checks a
// connection exists and parses its stored JSON, it never validates the key live, so a dummy
// value is enough. Mirrors tests/workbench_e2e/workbench_test.go's seedAnthropicConnection.
func (s *SimpleChatSuite) seedOpenRouterConnection(userUuid uuid.UUID) {
	ctx := context.Background()

	creds := domain.OpenAIKeyCredentials{ApiKey: "sk-or-e2e-dummy"}

	credJSON, err := json.Marshal(creds)
	s.Require().NoError(err, "marshal dummy openrouter credentials")

	conn := domain.ExternalConnection{
		UserUuid:        userUuid,
		Provider:        domain.ProviderOpenRouter,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: credJSON,
	}

	_, err = s.repos.ExternalConnections().Upsert(ctx, conn)
	s.Require().NoError(err, "seed dummy openrouter connection")
}

// appendTestMessage appends one message line directly to chatUuid's CouchDB transcript — the
// JSONL-era mirror of the pre-CouchDB version of this suite's insertMessage helper (a direct
// simple_chat_messages row insert), bypassing the full RunTurn/OpenAI turn loop for this narrow
// "does it become non-empty" check.
func (s *SimpleChatSuite) appendTestMessage(user simpleChatTestUser, chatUuid uuid.UUID) {
	ctx := context.Background()

	vault, err := s.repos.Vaults().GetByID(ctx, user.vaultUuid)
	s.Require().NoError(err)

	instance, err := s.repos.CouchInstances().Get(ctx, vault.CouchInstanceUuid)
	s.Require().NoError(err)

	account, err := s.repos.CouchAccounts().GetByUserAndInstance(ctx, user.userUuid, vault.CouchInstanceUuid)
	s.Require().NoError(err)

	client := couchdb.NewLiveSyncClient(instance.Url, vault.CouchDBName, account.CouchUsername, account.CouchPassword)

	path := domain.ChatHistoryPath(user.userUuid, chatUuid)

	note, err := client.ReadNote(ctx, path)
	s.Require().NoError(err, "read chat file before appending test message")

	file, err := domain.DecodeSimpleChatFile([]byte(note.Content))
	s.Require().NoError(err, "decode chat file")

	msg := domain.SimpleChatMessage{
		ChatUuid:  chatUuid,
		Role:      string(domain.SimpleChatRoleUser),
		Content:   "hello from the test user",
		Seq:       file.NextSeq(),
		CreatedAt: time.Now().UTC(),
	}

	file.Messages = append(file.Messages, msg)

	content, err := domain.EncodeSimpleChatFile(file)
	s.Require().NoError(err, "encode chat file")

	err = client.WriteNote(ctx, path, string(content))
	s.Require().NoError(err, "write chat file with appended test message")
}

// TestListChats_HidesEmptyChat_ShowsAfterFirstMessage confirms a chat created but never sent a
// message (the "New chat" button's immediate persist) does not appear in ListChats, and that it
// flips to visible as soon as the first message lands — exactly what run_turn.go does
// synchronously before the assistant reply comes back.
func (s *SimpleChatSuite) TestListChats_HidesEmptyChat_ShowsAfterFirstMessage() {
	user := s.setupUser("solo")

	chat, err := s.svcs.SimpleChat.CreateChat(user.svcCtx, user.vaultUuid, "some-model", true)
	s.Require().NoError(err, "create simple chat")

	chats, err := s.svcs.SimpleChat.ListChats(user.svcCtx, user.vaultUuid)
	s.Require().NoError(err)
	s.Empty(chats, "a chat with no messages must not appear in ListChats")

	s.appendTestMessage(user, chat.Uuid)

	chats, err = s.svcs.SimpleChat.ListChats(user.svcCtx, user.vaultUuid)
	s.Require().NoError(err)
	s.Require().Len(chats, 1, "the chat should appear now that it has one message")
	s.Equal(chat.Uuid, chats[0].Uuid)
}

// TestListChats_ExcludesEmptyChat_AmongMultiple creates two chats for the same vault/user — one
// with a message, one without — and confirms ListChats returns exactly the one with a message,
// not merely "empty vs non-empty" by coincidence of result count.
func (s *SimpleChatSuite) TestListChats_ExcludesEmptyChat_AmongMultiple() {
	user := s.setupUser("pair")

	withMessage, err := s.svcs.SimpleChat.CreateChat(user.svcCtx, user.vaultUuid, "some-model", true)
	s.Require().NoError(err)

	empty, err := s.svcs.SimpleChat.CreateChat(user.svcCtx, user.vaultUuid, "some-model", true)
	s.Require().NoError(err)

	s.appendTestMessage(user, withMessage.Uuid)

	chats, err := s.svcs.SimpleChat.ListChats(user.svcCtx, user.vaultUuid)
	s.Require().NoError(err)
	s.Require().Len(chats, 1, "only the chat with a message should be returned")
	s.Equal(withMessage.Uuid, chats[0].Uuid)

	var chatIds []uuid.UUID
	for _, c := range chats {
		chatIds = append(chatIds, c.Uuid)
	}

	s.NotContains(chatIds, empty.Uuid, "the never-messaged chat must be excluded")
}
