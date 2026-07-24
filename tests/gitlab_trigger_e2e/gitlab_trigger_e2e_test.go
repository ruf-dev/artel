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
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/service/v1/tract"
	"github.com/ruf-dev/artel/internal/transport/gitlab_webhook"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}

	return def
}

func randomEmail() string {
	return fmt.Sprintf("gitlab_trigger_e2e_%08x@test.local", rand.Uint32())
}

// mockGitlabServer records every create_merge_request / add_comment call it receives and
// answers with just enough shape for the tract engine to parse a usable response: a
// create_merge_request response carries "iid" (referenced by the add_comment step's template).
type mockGitlabServer struct {
	mu               sync.Mutex
	mergeRequestReqs []recordedRequest
	commentReqs      []recordedRequest
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

type GitlabTriggerE2ESuite struct {
	suite.Suite
	db    *sql.DB
	repos *repopg.Repos
	svcs  *svcv1.Services
}

func TestGitlabTriggerE2E(t *testing.T) {
	suite.Run(t, new(GitlabTriggerE2ESuite))
}

func (s *GitlabTriggerE2ESuite) SetupSuite() {
	pgDSN := envOrDefault("PG_DSN", "postgres://artel:artel_db@localhost:15434/artel_db?sslmode=disable")
	db, err := sql.Open("postgres", pgDSN)
	s.Require().NoError(err, "open postgres")

	err = db.Ping()
	s.Require().NoError(err, "ping postgres — is tests/docker-compose.yaml up?")
	s.db = db

	goose.SetLogger(goose.NopLogger())

	err = goose.SetDialect("postgres")
	s.Require().NoError(err)

	err = goose.Up(db, "../../migrations")
	s.Require().NoError(err, "run migrations")

	encKey := make([]byte, 32)
	encryptor, err := cryptoutil.NewAESEncryptor(encKey)
	s.Require().NoError(err, "create AES encryptor")
	s.repos = repopg.New(db, encryptor)

	cfg := config.EnvironmentConfig{}
	svcs, err := svcv1.New(s.repos, cfg)
	s.Require().NoError(err, "init services")
	s.svcs = svcs

	// Tract is normally wired in internal/app/custom.go (not svcv1.New) because its
	// ToolExecutor composes the already-built Mcp/Mom services — mirror that wiring exactly,
	// same as tests/tract_e2e, plus the new TriggerPresets repo.
	tractToolExecutor := tract.NewToolExecutor(s.svcs.McpService(), s.svcs.MomService())
	s.svcs.Tract = tract.New(
		s.repos.Tracts(),
		s.repos.Triggers(),
		s.repos.TriggerPresets(),
		s.repos.ExternalConnections(),
		s.repos.McpDefinitions(),
		tractToolExecutor,
	)
	s.svcs.McpService().SetTractService(context.Background(), s.svcs.TractService())
}

func (s *GitlabTriggerE2ESuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err)
	}
}

// TestGitlabPushTriggersMergeRequestAndComment drives the full preset flow: connect GitLab
// (bypassing the real SSRF-guarded validation call by inserting the connection row directly,
// same "not the subsystem under test" bypass tract_e2e uses for CouchDB provisioning), create a
// tract with two GitLab MoM action steps, create a gitlab_push trigger (provider-linked — no
// per-trigger secret), link it, then POST a simulated GitLab push webhook straight at
// gitlab_webhook.Handler and assert the mocked GitLab API received one create_merge_request call
// followed by one add_comment call with body "Created via Artel".
func (s *GitlabTriggerE2ESuite) TestGitlabPushTriggersMergeRequestAndComment() {
	ctx := context.Background()

	email := randomEmail()
	user, err := s.svcs.Auth.Register(ctx, email, "test-password-e2e")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), user.Uuid)
	})

	uc := user_context.UserContext{UserUuid: user.Uuid}
	userCtx := user_context.WithUserContext(ctx, uc)

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

	conn := domain.ExternalConnection{
		UserUuid:        user.Uuid,
		Provider:        domain.ProviderGitlab,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: credsJSON,
	}

	savedConn, err := s.repos.ExternalConnections().Upsert(ctx, conn)
	s.Require().NoError(err)

	// 2. Create a tract with two GitLab MoM action steps: create a merge request from the
	// pushed branch, then comment on it.
	createMrStep := domain.TractStep{
		Id:             "create_mr",
		Name:           "create_mr",
		Type:           "action",
		Mcp:            "gitlab",
		Tool:           "create_merge_request",
		ConnectionUuid: savedConn.Uuid,
		Params: map[string]string{
			"project_id":    "{{ trigger.project.id }}",
			"source_branch": "{{ trigger.branch }}",
			"target_branch": "main",
			"title":         "Automated MR from {{ trigger.branch }}",
		},
	}

	addCommentStep := domain.TractStep{
		Id:             "add_comment",
		Name:           "add_comment",
		Type:           "action",
		Mcp:            "gitlab",
		Tool:           "add_comment",
		ConnectionUuid: savedConn.Uuid,
		Params: map[string]string{
			"project_id":    "{{ trigger.project.id }}",
			"noteable_type": "merge_requests",
			"noteable_iid":  "{{ create_mr.iid }}",
			"body":          "Created via Artel",
		},
	}

	definition := domain.TractDefinition{Steps: []domain.TractStep{createMrStep, addCommentStep}}
	tr, warnings, err := s.svcs.Tract.CreateTract(userCtx, "e2e gitlab push tract", "", definition)
	s.Require().NoError(err)
	s.Empty(warnings)
	s.T().Cleanup(func() {
		_ = s.svcs.Tract.DeleteTract(context.Background(), tr.Uuid)
	})

	// 3. Create the gitlab_push trigger — provider-linked, so it shares savedConn's webhook
	// URL/secret instead of minting its own.
	trigger, rawToken, err := s.svcs.Tract.CreateTrigger(
		userCtx, "e2e gitlab push trigger", "webhook", tract.SourceGitlabPush, nil, domain.ToolSchema{},
	)
	s.Require().NoError(err)
	s.Empty(rawToken, "provider-linked triggers must not mint their own secret")
	s.T().Cleanup(func() {
		_ = s.svcs.Tract.DeleteTrigger(context.Background(), trigger.Uuid)
	})

	err = s.svcs.Tract.LinkTrigger(userCtx, trigger.Uuid, tr.Uuid, nil)
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
	var run domain.TractRun

	s.Require().Eventually(func() bool {
		runs, listErr := s.svcs.Tract.ListRuns(userCtx, tr.Uuid, 1)
		if listErr != nil || len(runs) == 0 {
			return false
		}

		run = runs[0]

		return run.Status != domain.TractRunRunning
	}, 5*time.Second, 50*time.Millisecond, "tract run did not finish")

	require.Equal(s.T(), domain.TractRunDone, run.Status, "run error: %s", run.Error)

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
