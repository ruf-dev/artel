//go:build e2e

package gitlab_trigger_e2e_test

// Covers the GitLab push-webhook trigger end-to-end against real Postgres (via
// tests/docker-compose.yaml) and a mocked GitLab REST API (httptest.Server): a simulated
// GitLab "push" webhook delivery hits gitlab_webhook.Handler, which fans out to the linked
// gitlab_push trigger, which starts a tract run that calls the (mocked) GitLab API to create a
// merge request, then comments "Created via Artel" on it. No CouchDB/vault dependency — every
// tract step here is a MoM (GitLab) action, not a builtin vault tool.
//
// Run: docker compose -f tests/docker-compose.yaml up -d
//      go test -tags e2e ./tests/gitlab_trigger_e2e/...

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/domain"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/service/v1/tract"
	"github.com/ruf-dev/artel/internal/service/v1/tract/script"
	"github.com/ruf-dev/artel/internal/transport/auth_api"
	"github.com/ruf-dev/artel/internal/transport/gitlab_webhook"
	"github.com/ruf-dev/artel/internal/transport/tracts_api"
	"github.com/ruf-dev/artel/tests/harness"
)

// mockGitlabServer records every create_merge_request / add_comment call it receives and
// answers with just enough shape for the tract engine to parse a usable response: a
// create_merge_request response carries "iid" (referenced by the add_comment step's template).
type mockGitlabServer struct {
	mu               sync.Mutex
	mergeRequestReqs []recordedRequest
	commentReqs      []recordedRequest
	diffReqs         []recordedRequest
	updateReqs       []recordedRequest
}

type recordedRequest struct {
	path  string
	query url.Values
	body  map[string]interface{}
}

func (m *mockGitlabServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var parsedBody map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&parsedBody)

	rec := recordedRequest{path: r.URL.Path, query: r.URL.Query(), body: parsedBody}

	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/api/v4/user":
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"username": "artel-bot"}`))
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/notes"):
		m.commentReqs = append(m.commentReqs, rec)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id": 1, "body": "Created via Artel"}`))
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/merge_requests"):
		m.mergeRequestReqs = append(m.mergeRequestReqs, rec)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"iid": 7, "id": 123, "title": "Automated MR"}`))
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/diffs"):
		m.diffReqs = append(m.diffReqs, rec)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"old_path": "a.go", "new_path": "a.go", "diff": "@@ -1 +1 @@\n-old\n+new"}]`))
	case r.Method == http.MethodPut && strings.Contains(r.URL.Path, "/merge_requests/"):
		m.updateReqs = append(m.updateReqs, rec)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"iid": 7, "state": "merged"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *mockGitlabServer) MergeRequestCalls() []recordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]recordedRequest(nil), m.mergeRequestReqs...)
}

func (m *mockGitlabServer) CommentCalls() []recordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]recordedRequest(nil), m.commentReqs...)
}

func (m *mockGitlabServer) DiffCalls() []recordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]recordedRequest(nil), m.diffReqs...)
}

func (m *mockGitlabServer) UpdateCalls() []recordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]recordedRequest(nil), m.updateReqs...)
}

type GitlabTriggerE2ESuite struct {
	suite.Suite
	db    *sql.DB
	repos *repopg.Repos
	svcs  *svcv1.Services

	authClient  pb.AuthAPIClient
	tractClient pb.TractsAPIClient
}

func TestGitlabTriggerE2E(t *testing.T) {
	suite.Run(t, new(GitlabTriggerE2ESuite))
}

func (s *GitlabTriggerE2ESuite) SetupSuite() {
	ctx := context.Background()

	s.db = harness.OpenPostgres(s.T())

	cfg := config.EnvironmentConfig{}

	var credsEncrypted bool
	s.repos, s.svcs, credsEncrypted = harness.BuildServices(s.T(), s.db, cfg)

	// Tract is normally wired in internal/app/custom.go (not svcv1.New) because its
	// ToolExecutor composes the already-built Mcp/Mom services — mirror that wiring exactly,
	// same as tests/tract_e2e, plus the new TriggerPresets repo.
	tractToolExecutor := tract.NewToolExecutor(s.svcs.McpService(), s.svcs.MomService())
	tractLlmExecutor := tract.NewLlmExecutor(s.repos.ExternalConnections())
	scriptEngines := script.NewRegistry(script.NewJavaScriptEngine())
	s.svcs.Tract = tract.New(
		s.repos.Tracts(),
		s.repos.TractTemplates(),
		s.repos.Triggers(),
		s.repos.TriggerPresets(),
		s.repos.ExternalConnections(),
		s.repos.McpDefinitions(),
		tractToolExecutor,
		s.svcs.SubscriptionService(),
		scriptEngines,
		tractLlmExecutor,
	)
	s.svcs.McpService().SetTractService(ctx, s.svcs.TractService())

	authImpl := auth_api.NewAuthImpl(
		s.svcs.Auth, "", s.svcs.S3Instance, s.svcs.CouchInstance,
		false, credsEncrypted, s.svcs.DockerHost, s.svcs.SetupWizard,
		false,
	)
	tractsImpl := tracts_api.New(ctx, s.svcs.TractService(), false)

	conn := harness.NewBufconnServer(s.T(), s.svcs, authImpl.Register, tractsImpl.Register)

	s.authClient = pb.NewAuthAPIClient(conn)
	s.tractClient = pb.NewTractsAPIClient(conn)
}

func (s *GitlabTriggerE2ESuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err)
	}
}

// registerAndLogin registers+logs in a fresh user via the real AuthAPI RPCs and returns an
// authenticated context plus the raw string user id (external_connections.user_uuid etc. still
// take the parsed uuid.UUID where needed at call sites).
func (s *GitlabTriggerE2ESuite) registerAndLogin(slug string) (context.Context, string) {
	email := slug + "@test.local"
	// A fixed literal, not slug-derived — some subtest names in this suite are long enough that
	// "test-password-"+slug exceeds bcrypt's 72-byte input cap. Uniqueness comes from email.
	const password = "test-password-e2e"

	registerReq := &pb.Register_Request{
		Email:    email,
		Password: password,
	}
	registerResp, err := s.authClient.Register(context.Background(), registerReq)
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		userUuid, parseErr := uuid.Parse(registerResp.Id)
		if parseErr == nil {
			_ = s.repos.Users().Delete(context.Background(), userUuid)
		}
	})

	passwordCreds := &pb.PasswordCredentials{Email: email, Password: password}
	loginMethod := &pb.Login_Request_Password{Password: passwordCreds}
	loginReq := &pb.Login_Request{Method: loginMethod}
	loginResp, err := s.authClient.Login(context.Background(), loginReq)
	s.Require().NoError(err)
	s.Require().NotEmpty(loginResp.Token)

	authedCtx := harness.AuthedContext(context.Background(), loginResp.Token)

	return authedCtx, registerResp.Id
}

// TestGitlabPushTriggersMergeRequestAndComment drives the full preset flow: connect GitLab
// (bypassing the real SSRF-guarded validation call by inserting the connection row directly,
// same "not the subsystem under test" bypass tract_e2e uses for CouchDB provisioning), create a
// tract with two GitLab MoM action steps, create a gitlab_push trigger (provider-linked — no
// per-trigger secret), link it, then POST a simulated GitLab push webhook straight at
// gitlab_webhook.Handler and assert the mocked GitLab API received one create_merge_request call
// followed by one add_comment call with body "Created via Artel".
func (s *GitlabTriggerE2ESuite) TestGitlabPushTriggersMergeRequestAndComment() {
	slug := harness.Slug(s.T())

	authedCtx, userUuidStr := s.registerAndLogin(slug)

	userUuid, err := uuid.Parse(userUuidStr)
	s.Require().NoError(err)

	mockGitlab := &mockGitlabServer{}
	mockServer := httptest.NewServer(mockGitlab)
	s.T().Cleanup(mockServer.Close)

	// 1. Connect GitLab — insert the external_connections row directly (AddGitlabConnection's
	// SSRF guard rejects loopback hosts like httptest's, and validating that guard isn't what
	// this test is about).
	const webhookSecret = "test-webhook-secret"
	creds := domain.GitlabCredentials{
		PersonalAccessToken: "test-personal-access-token",
		InstanceUrl:         mockServer.URL,
		WebhookSecret:       webhookSecret,
	}

	credsJSON, err := json.Marshal(creds)
	s.Require().NoError(err)

	newConn := domain.ExternalConnection{
		UserUuid:        userUuid,
		Provider:        domain.ProviderGitlab,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: credsJSON,
	}

	savedConn, err := s.repos.ExternalConnections().Upsert(context.Background(), newConn)
	s.Require().NoError(err)

	// 2. Create a tract with two GitLab MoM action steps: create a merge request from the
	// pushed branch, then comment on it.
	createMrAction := &pb.ActionStep{
		Mcp:            "gitlab",
		Tool:           "create_merge_request",
		ConnectionUuid: savedConn.Uuid.String(),
		Params: map[string]string{
			"project_id":    "{{ trigger.project.id }}",
			"source_branch": "{{ trigger.branch }}",
			"target_branch": "main",
			"title":         "Automated MR from {{ trigger.branch }}",
		},
	}
	createMrStep := &pb.TractStep{
		Id:   "create_mr",
		Name: "create_mr",
		Kind: &pb.TractStep_Action{Action: createMrAction},
	}

	addCommentAction := &pb.ActionStep{
		Mcp:            "gitlab",
		Tool:           "add_comment",
		ConnectionUuid: savedConn.Uuid.String(),
		Params: map[string]string{
			"project_id":    "{{ trigger.project.id }}",
			"noteable_type": "merge_requests",
			"noteable_iid":  "{{ create_mr.iid }}",
			"body":          "Created via Artel",
		},
	}
	addCommentStep := &pb.TractStep{
		Id:   "add_comment",
		Name: "add_comment",
		Kind: &pb.TractStep_Action{Action: addCommentAction},
	}

	definition := &pb.TractDefinition{Steps: []*pb.TractStep{createMrStep, addCommentStep}}
	createTractReq := &pb.CreateTract_Request{
		Name:       "e2e gitlab push tract",
		Definition: definition,
	}
	createTractResp, err := s.tractClient.CreateTract(authedCtx, createTractReq)
	s.Require().NoError(err)
	s.Empty(createTractResp.Warnings)

	tractUuid := createTractResp.Tract.Uuid
	s.T().Cleanup(func() {
		deleteTractReq := &pb.DeleteTract_Request{Uuid: tractUuid}
		_, _ = s.tractClient.DeleteTract(authedCtx, deleteTractReq)
	})

	// 3. Create the gitlab_push trigger — provider-linked, so it shares savedConn's webhook
	// URL/secret instead of minting its own.
	createTriggerReq := &pb.CreateTrigger_Request{
		Name:          "e2e gitlab push trigger",
		Kind:          "webhook",
		Source:        tract.SourceGitlabPush,
		PayloadSchema: "{}",
	}
	createTriggerResp, err := s.tractClient.CreateTrigger(authedCtx, createTriggerReq)
	s.Require().NoError(err)
	s.Empty(createTriggerResp.WebhookToken, "provider-linked triggers must not mint their own secret")

	triggerUuid := createTriggerResp.Trigger.Uuid
	s.T().Cleanup(func() {
		deleteTriggerReq := &pb.DeleteTrigger_Request{Uuid: triggerUuid}
		_, _ = s.tractClient.DeleteTrigger(authedCtx, deleteTriggerReq)
	})

	linkTriggerReq := &pb.LinkTrigger_Request{TriggerUuid: triggerUuid, TractUuid: tractUuid}
	_, err = s.tractClient.LinkTrigger(authedCtx, linkTriggerReq)
	s.Require().NoError(err)

	// 4. Fire a simulated GitLab push webhook straight at the shared connection endpoint —
	// the same request shape/headers a real GitLab instance sends.
	pushPayload := map[string]interface{}{
		"ref":       "refs/heads/feature-x",
		"user_name": "tester",
		"project": map[string]interface{}{
			"id":   42,
			"name": "demo",
			"path": "group/demo",
		},
		"commits": []interface{}{},
	}

	body, err := json.Marshal(pushPayload)
	s.Require().NoError(err)

	handler := gitlab_webhook.New(
		context.Background(), s.repos.ExternalConnections(), s.repos.Triggers(), s.svcs.MomService(), s.svcs.TractService(),
	)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/gitlab/"+savedConn.Uuid.String(), bytes.NewReader(body))
	req.Header.Set("X-Gitlab-Event", "Push Hook")
	req.Header.Set("X-Gitlab-Token", webhookSecret)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)

	// 5. DispatchTrigger spawns StartRun asynchronously — poll until the run finishes.
	var runItem *pb.TractRunItem

	s.Require().Eventually(func() bool {
		listRunsReq := &pb.ListRuns_Request{TractUuid: tractUuid, Limit: 1}
		listRunsResp, listErr := s.tractClient.ListRuns(authedCtx, listRunsReq)
		if listErr != nil || len(listRunsResp.Runs) == 0 {
			return false
		}

		runItem = listRunsResp.Runs[0]

		return runItem.Status != string(domain.TractRunRunning)
	}, 5*time.Second, 50*time.Millisecond, "tract run did not finish")

	require.Equal(s.T(), string(domain.TractRunDone), runItem.Status, "run error: %s", runItem.Error)

	// 6. Assert the mocked GitLab API actually saw one create_merge_request call followed by
	// one add_comment call with the expected params. create_merge_request sends its params as a
	// JSON body (per migration 034_gitlab_tract_tools.sql), project_id is a URL path segment.
	mrCalls := mockGitlab.MergeRequestCalls()
	s.Require().Len(mrCalls, 1)
	s.True(strings.HasSuffix(mrCalls[0].path, "/projects/42/merge_requests"))
	s.Equal("feature-x", mrCalls[0].body["source_branch"])
	s.Equal("main", mrCalls[0].body["target_branch"])

	commentCalls := mockGitlab.CommentCalls()
	s.Require().Len(commentCalls, 1)
	s.Equal("Created via Artel", commentCalls[0].query.Get("body"))
	s.True(strings.Contains(commentCalls[0].path, "/merge_requests/7/notes"))
}

// TestGitlabMergeRequestMergedTriggersDiffAndDescriptionUpdate covers the second new preset
// (gitlab_merge_request, added by migration 057_gitlab_mr_merged_and_diff.sql): a simulated
// GitLab "Merge Request Hook" delivery with object_attributes.action == "merge" must pass the
// trigger's CheckBody matcher (narrowing past the shared X-Gitlab-Event: Merge Request Hook
// header, which GitLab also sends for opened/updated/closed) and fire a tract run that calls the
// new get_merge_request_diff tool followed by update_merge_request (writing a description back
// onto the MR) — the two GitLab tool gaps from the same migration.
func (s *GitlabTriggerE2ESuite) TestGitlabMergeRequestMergedTriggersDiffAndDescriptionUpdate() {
	slug := harness.Slug(s.T())

	authedCtx, userUuidStr := s.registerAndLogin(slug)

	userUuid, err := uuid.Parse(userUuidStr)
	s.Require().NoError(err)

	mockGitlab := &mockGitlabServer{}
	mockServer := httptest.NewServer(mockGitlab)
	s.T().Cleanup(mockServer.Close)

	const webhookSecret = "test-webhook-secret-mr"
	creds := domain.GitlabCredentials{
		PersonalAccessToken: "test-personal-access-token",
		InstanceUrl:         mockServer.URL,
		WebhookSecret:       webhookSecret,
	}

	credsJSON, err := json.Marshal(creds)
	s.Require().NoError(err)

	newConn := domain.ExternalConnection{
		UserUuid:        userUuid,
		Provider:        domain.ProviderGitlab,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: credsJSON,
	}

	savedConn, err := s.repos.ExternalConnections().Upsert(context.Background(), newConn)
	s.Require().NoError(err)

	// Tract: read the merged MR's diff, then write a description back onto it — the exact flow
	// the "MR merged" preset + diff/description tools were added for.
	getDiffAction := &pb.ActionStep{
		Mcp:            "gitlab",
		Tool:           "get_merge_request_diff",
		ConnectionUuid: savedConn.Uuid.String(),
		Params: map[string]string{
			"project_id": "{{ trigger.project.id }}",
			"mr_iid":     "{{ trigger.mr_iid }}",
		},
	}
	getDiffStep := &pb.TractStep{
		Id:   "get_diff",
		Name: "get_diff",
		Kind: &pb.TractStep_Action{Action: getDiffAction},
	}

	updateMrAction := &pb.ActionStep{
		Mcp:            "gitlab",
		Tool:           "update_merge_request",
		ConnectionUuid: savedConn.Uuid.String(),
		Params: map[string]string{
			"project_id":  "{{ trigger.project.id }}",
			"mr_iid":      "{{ trigger.mr_iid }}",
			"description": "Work done: reviewed the diff",
		},
	}
	updateMrStep := &pb.TractStep{
		Id:   "update_mr",
		Name: "update_mr",
		Kind: &pb.TractStep_Action{Action: updateMrAction},
	}

	definition := &pb.TractDefinition{Steps: []*pb.TractStep{getDiffStep, updateMrStep}}
	createTractReq := &pb.CreateTract_Request{
		Name:       "e2e gitlab mr merged tract",
		Definition: definition,
	}
	createTractResp, err := s.tractClient.CreateTract(authedCtx, createTractReq)
	s.Require().NoError(err)
	s.Empty(createTractResp.Warnings)

	tractUuid := createTractResp.Tract.Uuid
	s.T().Cleanup(func() {
		deleteTractReq := &pb.DeleteTract_Request{Uuid: tractUuid}
		_, _ = s.tractClient.DeleteTract(authedCtx, deleteTractReq)
	})

	// gitlab_merge_request is provider-linked, same as gitlab_push — shares savedConn's webhook
	// URL/secret, no per-trigger secret.
	createTriggerReq := &pb.CreateTrigger_Request{
		Name:          "e2e gitlab mr merged trigger",
		Kind:          "webhook",
		Source:        tract.SourceGitlabMergeRequest,
		PayloadSchema: "{}",
	}
	createTriggerResp, err := s.tractClient.CreateTrigger(authedCtx, createTriggerReq)
	s.Require().NoError(err)
	s.Empty(createTriggerResp.WebhookToken, "provider-linked triggers must not mint their own secret")

	triggerUuid := createTriggerResp.Trigger.Uuid
	s.T().Cleanup(func() {
		deleteTriggerReq := &pb.DeleteTrigger_Request{Uuid: triggerUuid}
		_, _ = s.tractClient.DeleteTrigger(authedCtx, deleteTriggerReq)
	})

	linkTriggerReq := &pb.LinkTrigger_Request{TriggerUuid: triggerUuid, TractUuid: tractUuid}
	_, err = s.tractClient.LinkTrigger(authedCtx, linkTriggerReq)
	s.Require().NoError(err)

	// Fire a simulated GitLab "Merge Request Hook" delivery with action == "merge" — the shape
	// the trigger's CheckBody matcher requires to fire (opened/updated/closed deliveries share
	// the same X-Gitlab-Event header and must NOT fire this trigger, see the sibling
	// non-matching-action test below).
	mrPayload := map[string]interface{}{
		"object_attributes": map[string]interface{}{
			"iid":           7,
			"title":         "Fix bug",
			"source_branch": "fix-bug",
			"target_branch": "main",
			"action":        "merge",
			"merge_status":  "can_be_merged",
		},
		"project": map[string]interface{}{
			"id":   42,
			"name": "demo",
			"path": "group/demo",
		},
	}

	body, err := json.Marshal(mrPayload)
	s.Require().NoError(err)

	handler := gitlab_webhook.New(
		context.Background(), s.repos.ExternalConnections(), s.repos.Triggers(), s.svcs.MomService(), s.svcs.TractService(),
	)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/gitlab/"+savedConn.Uuid.String(), bytes.NewReader(body))
	req.Header.Set("X-Gitlab-Event", "Merge Request Hook")
	req.Header.Set("X-Gitlab-Token", webhookSecret)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)

	var runItem *pb.TractRunItem

	s.Require().Eventually(func() bool {
		listRunsReq := &pb.ListRuns_Request{TractUuid: tractUuid, Limit: 1}
		listRunsResp, listErr := s.tractClient.ListRuns(authedCtx, listRunsReq)
		if listErr != nil || len(listRunsResp.Runs) == 0 {
			return false
		}

		runItem = listRunsResp.Runs[0]

		return runItem.Status != string(domain.TractRunRunning)
	}, 5*time.Second, 50*time.Millisecond, "tract run did not finish")

	require.Equal(s.T(), string(domain.TractRunDone), runItem.Status, "run error: %s", runItem.Error)

	diffCalls := mockGitlab.DiffCalls()
	s.Require().Len(diffCalls, 1)
	s.True(strings.HasSuffix(diffCalls[0].path, "/projects/42/merge_requests/7/diffs"))

	updateCalls := mockGitlab.UpdateCalls()
	s.Require().Len(updateCalls, 1)
	s.True(strings.HasSuffix(updateCalls[0].path, "/projects/42/merge_requests/7"))
	s.Equal("Work done: reviewed the diff", updateCalls[0].body["description"])
}

// TestGitlabMergeRequestNonMergeActionDoesNotTriggerRun asserts the CheckBody matcher actually
// narrows dispatch: an "update" action delivery shares the exact same X-Gitlab-Event: Merge
// Request Hook header as a "merge" delivery, so without the body-level action check the trigger
// would misfire on every MR event, not just merges.
func (s *GitlabTriggerE2ESuite) TestGitlabMergeRequestNonMergeActionDoesNotTriggerRun() {
	slug := harness.Slug(s.T())

	authedCtx, userUuidStr := s.registerAndLogin(slug)

	userUuid, err := uuid.Parse(userUuidStr)
	s.Require().NoError(err)

	mockGitlab := &mockGitlabServer{}
	mockServer := httptest.NewServer(mockGitlab)
	s.T().Cleanup(mockServer.Close)

	const webhookSecret = "test-webhook-secret-mr-nomatch"
	creds := domain.GitlabCredentials{
		PersonalAccessToken: "test-personal-access-token",
		InstanceUrl:         mockServer.URL,
		WebhookSecret:       webhookSecret,
	}

	credsJSON, err := json.Marshal(creds)
	s.Require().NoError(err)

	newConn := domain.ExternalConnection{
		UserUuid:        userUuid,
		Provider:        domain.ProviderGitlab,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: credsJSON,
	}

	savedConn, err := s.repos.ExternalConnections().Upsert(context.Background(), newConn)
	s.Require().NoError(err)

	getDiffAction := &pb.ActionStep{
		Mcp:            "gitlab",
		Tool:           "get_merge_request_diff",
		ConnectionUuid: savedConn.Uuid.String(),
		Params: map[string]string{
			"project_id": "{{ trigger.project.id }}",
			"mr_iid":     "{{ trigger.mr_iid }}",
		},
	}
	getDiffStep := &pb.TractStep{
		Id:   "get_diff",
		Name: "get_diff",
		Kind: &pb.TractStep_Action{Action: getDiffAction},
	}

	definition := &pb.TractDefinition{Steps: []*pb.TractStep{getDiffStep}}
	createTractReq := &pb.CreateTract_Request{
		Name:       "e2e gitlab mr non-merge tract",
		Definition: definition,
	}
	createTractResp, err := s.tractClient.CreateTract(authedCtx, createTractReq)
	s.Require().NoError(err)
	s.Empty(createTractResp.Warnings)

	tractUuid := createTractResp.Tract.Uuid
	s.T().Cleanup(func() {
		deleteTractReq := &pb.DeleteTract_Request{Uuid: tractUuid}
		_, _ = s.tractClient.DeleteTract(authedCtx, deleteTractReq)
	})

	createTriggerReq := &pb.CreateTrigger_Request{
		Name:          "e2e gitlab mr non-merge trigger",
		Kind:          "webhook",
		Source:        tract.SourceGitlabMergeRequest,
		PayloadSchema: "{}",
	}
	createTriggerResp, err := s.tractClient.CreateTrigger(authedCtx, createTriggerReq)
	s.Require().NoError(err)

	triggerUuid := createTriggerResp.Trigger.Uuid
	s.T().Cleanup(func() {
		deleteTriggerReq := &pb.DeleteTrigger_Request{Uuid: triggerUuid}
		_, _ = s.tractClient.DeleteTrigger(authedCtx, deleteTriggerReq)
	})

	linkTriggerReq := &pb.LinkTrigger_Request{TriggerUuid: triggerUuid, TractUuid: tractUuid}
	_, err = s.tractClient.LinkTrigger(authedCtx, linkTriggerReq)
	s.Require().NoError(err)

	mrPayload := map[string]interface{}{
		"object_attributes": map[string]interface{}{
			"iid":           7,
			"action":        "update",
			"source_branch": "fix-bug",
			"target_branch": "main",
		},
		"project": map[string]interface{}{"id": 42},
	}

	body, err := json.Marshal(mrPayload)
	s.Require().NoError(err)

	handler := gitlab_webhook.New(
		context.Background(), s.repos.ExternalConnections(), s.repos.Triggers(), s.svcs.MomService(), s.svcs.TractService(),
	)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/gitlab/"+savedConn.Uuid.String(), bytes.NewReader(body))
	req.Header.Set("X-Gitlab-Event", "Merge Request Hook")
	req.Header.Set("X-Gitlab-Token", webhookSecret)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)

	// Give any (incorrectly) spawned run a moment to appear, then assert none did.
	time.Sleep(200 * time.Millisecond)

	listRunsReq := &pb.ListRuns_Request{TractUuid: tractUuid, Limit: 10}
	listRunsResp, err := s.tractClient.ListRuns(authedCtx, listRunsReq)
	s.Require().NoError(err)
	s.Empty(listRunsResp.Runs, "an 'update' action delivery must not fire the gitlab_merge_request (merged-only) trigger")

	s.Empty(mockGitlab.DiffCalls(), "get_merge_request_diff must not have been called")
}
