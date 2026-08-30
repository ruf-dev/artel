//go:build e2e

package e2e_test

// Covers the five Trello write tools seeded by migration 081_trello_write_tools.sql end to end
// against real Postgres (via tests/docker-compose.yaml + the make test-e2e bootstrap that applies
// the migration) and a mocked Trello REST API (httptest.Server). No external Trello credentials,
// no network egress: one mcp-scoped HttpMiddleware (svcv1.WithMomMcpHttpMiddleware, from PR
// mom-http-middleware) rewrites every outbound trello request onto the mock server.
//
// Run: docker compose -f tests/docker-compose.yaml up -d --wait
//      go test -tags "e2e e2e_bootstrap" -count=1 ./tests/bootstrap/... -run TestEnvSetup
//      go test -tags e2e ./tests/e2e/... -run TestTrelloMom

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/ruf-dev/artel/internal/config"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
	"github.com/ruf-dev/artel/tests/harness"
)

type recordedTrelloRequest struct {
	method string
	path   string
	query  url.Values
}

// mockTrelloServer records every request it receives and answers with a canned body shaped like
// the real Trello API's response for that path, or 400 for any path containing "/bad".
type mockTrelloServer struct {
	mu       sync.Mutex
	requests []recordedTrelloRequest
}

func (m *mockTrelloServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	rec := recordedTrelloRequest{method: r.Method, path: r.URL.Path, query: r.URL.Query()}
	m.requests = append(m.requests, rec)
	m.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(r.URL.Path, "/bad") {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"invalid value"}`))

		return
	}

	switch {
	case r.Method == http.MethodPut && r.URL.Path == "/1/cards/card123":
		_, _ = w.Write([]byte(
			`{"id":"card123","name":"renamed","desc":"new body","idList":"l1","closed":false,"url":"https://trello.com/c/card123"}`,
		))
	case r.Method == http.MethodPost && r.URL.Path == "/1/lists":
		_, _ = w.Write([]byte(`{"id":"list-new","name":"Backlog","closed":false,"idBoard":"b1"}`))
	case r.Method == http.MethodPut && r.URL.Path == "/1/lists/l9":
		_, _ = w.Write([]byte(`{"id":"l9","name":"Archived","closed":true,"idBoard":"b1"}`))
	case r.Method == http.MethodPost && r.URL.Path == "/1/labels":
		_, _ = w.Write([]byte(`{"id":"label-new","name":"bug","color":"red","idBoard":"b1"}`))
	case r.Method == http.MethodPost && r.URL.Path == "/1/cards/card123/attachments":
		_, _ = w.Write([]byte(
			`{"id":"att-1","name":"spec","url":"https://ex.com/x","date":"2024-01-01T00:00:00.000Z","mimeType":"","bytes":0}`,
		))
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"not found"}`))
	}
}

func (m *mockTrelloServer) last() recordedTrelloRequest {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.requests[len(m.requests)-1]
}

// roundTripFunc adapts a plain function to http.RoundTripper.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// redirectTo is the test HttpMiddleware: it rewrites the request's scheme/host (and Host header)
// onto target before handing off to the real transport, leaving path + query untouched.
func redirectTo(target *url.URL) executors.HttpMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripFunc(func(r *http.Request) (*http.Response, error) {
			r.URL.Scheme = target.Scheme
			r.URL.Host = target.Host
			r.Host = target.Host

			return next.RoundTrip(r)
		})
	}
}

type TrelloMomSuite struct {
	suite.Suite

	db   *sql.DB
	svcs *svcv1.Services
	mock *mockTrelloServer
}

func TestTrelloMom(t *testing.T) {
	suite.Run(t, new(TrelloMomSuite))
}

func (s *TrelloMomSuite) SetupSuite() {
	s.db = harness.OpenPostgres(s.T())

	s.mock = &mockTrelloServer{}
	mockServer := httptest.NewServer(s.mock)
	s.T().Cleanup(mockServer.Close)

	target, err := url.Parse(mockServer.URL)
	s.Require().NoError(err)

	cfg := config.EnvironmentConfig{}
	redirect := redirectTo(target)
	middlewareOpt := svcv1.WithMomMcpHttpMiddleware(executors.McpNameTrello, redirect)

	_, s.svcs, _ = harness.BuildServices(s.T(), s.db, cfg, middlewareOpt)
}

func (s *TrelloMomSuite) execTrello(tool string, params map[string]interface{}) (string, error) {
	secrets := map[string]interface{}{"api_key": "k", "api_token": "t"}

	return s.svcs.MomService().ExecuteToolWithSecrets(context.Background(), "trello", tool, secrets, params)
}

func (s *TrelloMomSuite) TestTrelloWriteToolsRoundTrip() {
	s.Run("update_card", func() {
		params := map[string]interface{}{"card_id": "card123", "name": "renamed", "desc": "new body"}

		out, err := s.execTrello("update_card", params)
		s.Require().NoError(err)

		rec := s.mock.last()
		s.Equal(http.MethodPut, rec.method)
		s.Equal("/1/cards/card123", rec.path)
		s.Equal("renamed", rec.query.Get("name"))
		s.Equal("new body", rec.query.Get("desc"))
		s.Equal("k", rec.query.Get("key"))
		s.Equal("t", rec.query.Get("token"))

		var card struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}
		s.Require().NoError(json.Unmarshal([]byte(out), &card))
		s.Equal("card123", card.Id)
		s.Equal("renamed", card.Name)
	})

	s.Run("create_list", func() {
		params := map[string]interface{}{"board_id": "b1", "name": "Backlog"}

		out, err := s.execTrello("create_list", params)
		s.Require().NoError(err)

		rec := s.mock.last()
		s.Equal(http.MethodPost, rec.method)
		s.Equal("/1/lists", rec.path)
		s.Equal("Backlog", rec.query.Get("name"))
		s.Equal("b1", rec.query.Get("idBoard"))
		s.Empty(rec.query.Get("pos"), "unset optional param must be dropped, not sent empty")
		s.Equal("k", rec.query.Get("key"))
		s.Equal("t", rec.query.Get("token"))

		var list struct {
			Id      string `json:"id"`
			Name    string `json:"name"`
			IdBoard string `json:"idBoard"`
		}
		s.Require().NoError(json.Unmarshal([]byte(out), &list))
		s.Equal("list-new", list.Id)
		s.Equal("Backlog", list.Name)
		s.Equal("b1", list.IdBoard)
	})

	s.Run("update_list", func() {
		params := map[string]interface{}{"list_id": "l9", "closed": true}

		out, err := s.execTrello("update_list", params)
		s.Require().NoError(err)

		rec := s.mock.last()
		s.Equal(http.MethodPut, rec.method)
		s.Equal("/1/lists/l9", rec.path)
		s.Equal("true", rec.query.Get("closed"))
		s.Equal("k", rec.query.Get("key"))
		s.Equal("t", rec.query.Get("token"))

		var list struct {
			Id     string `json:"id"`
			Name   string `json:"name"`
			Closed bool   `json:"closed"`
		}
		s.Require().NoError(json.Unmarshal([]byte(out), &list))
		s.Equal("l9", list.Id)
		s.Equal("Archived", list.Name)
		s.True(list.Closed)
	})

	s.Run("create_label", func() {
		params := map[string]interface{}{"board_id": "b1", "name": "bug", "color": "red"}

		out, err := s.execTrello("create_label", params)
		s.Require().NoError(err)

		rec := s.mock.last()
		s.Equal(http.MethodPost, rec.method)
		s.Equal("/1/labels", rec.path)
		s.Equal("bug", rec.query.Get("name"))
		s.Equal("red", rec.query.Get("color"))
		s.Equal("b1", rec.query.Get("idBoard"))
		s.Equal("k", rec.query.Get("key"))
		s.Equal("t", rec.query.Get("token"))

		var label struct {
			Id    string `json:"id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		}
		s.Require().NoError(json.Unmarshal([]byte(out), &label))
		s.Equal("label-new", label.Id)
		s.Equal("bug", label.Name)
		s.Equal("red", label.Color)
	})

	s.Run("add_attachment", func() {
		params := map[string]interface{}{"card_id": "card123", "url": "https://ex.com/x", "name": "spec"}

		out, err := s.execTrello("add_attachment", params)
		s.Require().NoError(err)

		rec := s.mock.last()
		s.Equal(http.MethodPost, rec.method)
		s.Equal("/1/cards/card123/attachments", rec.path)
		s.Equal("https://ex.com/x", rec.query.Get("url"))
		s.Equal("spec", rec.query.Get("name"))
		s.Equal("k", rec.query.Get("key"))
		s.Equal("t", rec.query.Get("token"))

		var att struct {
			Id  string `json:"id"`
			Url string `json:"url"`
		}
		s.Require().NoError(json.Unmarshal([]byte(out), &att))
		s.Equal("att-1", att.Id)
		s.Equal("https://ex.com/x", att.Url)
	})

	s.Run("http 400 surfaces as an error", func() {
		params := map[string]interface{}{"card_id": "bad"}

		_, err := s.execTrello("update_card", params)
		s.Require().Error(err)
		s.Contains(err.Error(), "status 400")
	})
}
