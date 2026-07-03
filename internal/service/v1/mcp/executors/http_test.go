package executors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ruf-dev/artel/internal/domain"
)

// TestHttpExecutor_QueryParam_EmptyOmitted verifies that a query param whose
// resolved value is the empty string (e.g. a missing ${{params.x}}) is omitted
// from the outgoing URL entirely, rather than being sent as "x=".
func TestHttpExecutor_QueryParam_EmptyOmitted(t *testing.T) {
	var gotQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	httpAction := &domain.HttpAction{
		Method: http.MethodGet,
		Url:    server.URL,
		Query: map[string]string{
			"assignee_id": "${{params.assignee_id}}",
			"state":       "${{params.state}}",
		},
	}
	action := domain.ToolAction{Http: httpAction}

	params := map[string]interface{}{
		"state": "opened",
	}

	e := NewHttpExecutor()

	_, err := e.Execute(context.Background(), action, nil, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := gotQuery["assignee_id"]; ok {
		t.Fatalf("expected empty query param assignee_id to be omitted, got %v", gotQuery)
	}

	state := gotQuery.Get("state")
	if state != "opened" {
		t.Fatalf("expected state=opened, got %q", state)
	}
}

// TestHttpExecutor_Header_EmptyOmitted verifies that a header whose resolved
// value is the empty string is not sent on the outgoing request.
func TestHttpExecutor_Header_EmptyOmitted(t *testing.T) {
	var gotHeaders http.Header

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	httpAction := &domain.HttpAction{
		Method: http.MethodGet,
		Url:    server.URL,
		Headers: map[string]string{
			"X-Empty":    "${{params.missing}}",
			"X-Provided": "${{params.provided}}",
		},
	}
	action := domain.ToolAction{Http: httpAction}

	params := map[string]interface{}{
		"provided": "value",
	}

	e := NewHttpExecutor()

	_, err := e.Execute(context.Background(), action, nil, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotHeaders.Get("X-Empty") != "" {
		t.Fatalf("expected X-Empty header to be absent, got %q", gotHeaders.Get("X-Empty"))
	}
	if _, ok := gotHeaders["X-Empty"]; ok {
		t.Fatalf("expected X-Empty header to not be set at all, got %v", gotHeaders)
	}

	provided := gotHeaders.Get("X-Provided")
	if provided != "value" {
		t.Fatalf("expected X-Provided=value, got %q", provided)
	}
}

// TestHttpExecutor_BodyKey_EmptyDropped verifies that a JSON body key whose
// rendered value is the empty string is dropped from the marshaled body, while
// a key with a non-empty value is kept.
func TestHttpExecutor_BodyKey_EmptyDropped(t *testing.T) {
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&gotBody)
		if err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	bodyTemplate := []byte(`{"title":"${{params.title}}","assignee_id":"${{params.assignee_id}}"}`)

	httpAction := &domain.HttpAction{
		Method: http.MethodPost,
		Url:    server.URL,
		Body:   bodyTemplate,
	}
	action := domain.ToolAction{Http: httpAction}

	params := map[string]interface{}{
		"title": "my issue",
	}

	e := NewHttpExecutor()

	_, err := e.Execute(context.Background(), action, nil, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := gotBody["assignee_id"]; ok {
		t.Fatalf("expected empty assignee_id key to be dropped from body, got %v", gotBody)
	}

	title, ok := gotBody["title"].(string)
	if !ok || title != "my issue" {
		t.Fatalf("expected title=my issue, got %v", gotBody["title"])
	}
}
