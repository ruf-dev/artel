package vaults_api

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// fakeTerminalShellTargetResolver is a hand-rolled fake of the handler's
// terminalShellTargetResolver dependency.
type fakeTerminalShellTargetResolver struct {
	resolveFunc func(ctx context.Context, vaultID, userID uuid.UUID) (string, error)
}

func (f *fakeTerminalShellTargetResolver) ResolveTerminalShellTarget(
	ctx context.Context, vaultID, userID uuid.UUID,
) (string, error) {
	return f.resolveFunc(ctx, vaultID, userID)
}

// newTerminalShellMux registers handler under TerminalShellRoutePattern on a real
// http.ServeMux, so the tests exercise the same {vaultId}/{rest...} wildcard parsing
// production does rather than hand-setting path values.
func newTerminalShellMux(handler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(TerminalShellRoutePattern, handler)

	return mux
}

// terminalShellRequest builds a GET against the terminal-shell route for vaultID with rest
// appended after ".../terminal-shell/", optionally carrying an access_token cookie.
func terminalShellRequest(vaultID uuid.UUID, rest string, token string) *http.Request {
	target := "/api/vaults/workbench/" + vaultID.String() + "/terminal-shell/" + rest

	req := httptest.NewRequest(http.MethodGet, target, nil)

	if token != "" {
		cookie := &http.Cookie{Name: middleware.AccessTokenCookieName, Value: token}
		req.AddCookie(cookie)
	}

	return req
}

func TestWorkbenchTerminalShell_Unauthenticated(t *testing.T) {
	authSvc := &fakeTerminalAuthService{
		validateTokenFunc: func(ctx context.Context, token string) (domain.User, error) {
			return domain.User{}, errors.New("no such session")
		},
	}

	members := &fakeTerminalVaultMembers{
		getFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error) {
			t.Fatal("membership must not be checked for an unauthenticated request")

			return domain.VaultMember{}, nil
		},
	}

	workbenchSvc := &fakeTerminalShellTargetResolver{
		resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
			t.Fatal("terminal shell target must not be resolved for an unauthenticated request")

			return "", nil
		},
	}

	handler := NewWorkbenchTerminalShellHandler(authSvc, members, workbenchSvc)

	rec := httptest.NewRecorder()
	newTerminalShellMux(handler).ServeHTTP(rec, terminalShellRequest(uuid.New(), "", ""))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	if len(authSvc.seenTokens) != 1 || authSvc.seenTokens[0] != "" {
		t.Fatalf("ValidateToken tokens = %q, want exactly one empty token", authSvc.seenTokens)
	}
}

func TestWorkbenchTerminalShell_NonMember(t *testing.T) {
	authSvc := authedAs(uuid.New())

	members := &fakeTerminalVaultMembers{
		getFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{}, user_errors.NotFound
		},
	}

	workbenchSvc := &fakeTerminalShellTargetResolver{
		resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
			t.Fatal("terminal shell target must not be resolved for a non-member")

			return "", nil
		},
	}

	handler := NewWorkbenchTerminalShellHandler(authSvc, members, workbenchSvc)

	rec := httptest.NewRecorder()
	newTerminalShellMux(handler).ServeHTTP(rec, terminalShellRequest(uuid.New(), "", "token-1"))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestWorkbenchTerminalShell_WorkbenchNotRunning(t *testing.T) {
	authSvc := authedAs(uuid.New())

	workbenchSvc := &fakeTerminalShellTargetResolver{
		resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
			return "", user_errors.WorkbenchNotRunning
		},
	}

	handler := NewWorkbenchTerminalShellHandler(authSvc, alwaysMember(), workbenchSvc)

	rec := httptest.NewRecorder()
	newTerminalShellMux(handler).ServeHTTP(rec, terminalShellRequest(uuid.New(), "", "token-1"))

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestWorkbenchTerminalShell_InvalidVaultId(t *testing.T) {
	authSvc := authedAs(uuid.New())

	workbenchSvc := &fakeTerminalShellTargetResolver{
		resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
			t.Fatal("terminal shell target must not be resolved for an unparseable vault id")

			return "", nil
		},
	}

	handler := NewWorkbenchTerminalShellHandler(authSvc, alwaysMember(), workbenchSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/vaults/workbench/not-a-uuid/terminal-shell/", nil)

	rec := httptest.NewRecorder()
	newTerminalShellMux(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// TestWorkbenchTerminalShell_Proxies stands a httptest.Server in for the in-container ttyd
// server and confirms the route prefix is stripped, the query string survives, and the
// upstream response body/headers/status are forwarded back verbatim.
func TestWorkbenchTerminalShell_Proxies(t *testing.T) {
	tests := []struct {
		name     string
		rest     string
		wantPath string
	}{
		{name: "terminal-shell root serves ttyd index", rest: "", wantPath: "/"},
		{name: "websocket endpoint", rest: "ws", wantPath: "/ws"},
		{name: "nested static asset", rest: "static/xterm.css", wantPath: "/static/xterm.css"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath, gotQuery string

			ttydHandler := func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				gotQuery = r.URL.RawQuery

				w.Header().Set("X-Ttyd", "yes")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ttyd payload"))
			}

			ttyd := httptest.NewServer(http.HandlerFunc(ttydHandler))
			defer ttyd.Close()

			authSvc := authedAs(uuid.New())

			workbenchSvc := &fakeTerminalShellTargetResolver{
				resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
					return ttyd.URL, nil
				},
			}

			handler := NewWorkbenchTerminalShellHandler(authSvc, alwaysMember(), workbenchSvc)

			req := terminalShellRequest(uuid.New(), tt.rest+"?arg=xyz", "token-1")

			rec := httptest.NewRecorder()
			newTerminalShellMux(handler).ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}

			if gotPath != tt.wantPath {
				t.Fatalf("upstream path = %q, want %q", gotPath, tt.wantPath)
			}

			if gotQuery != "arg=xyz" {
				t.Fatalf("upstream query = %q, want %q", gotQuery, "arg=xyz")
			}

			if rec.Header().Get("X-Ttyd") != "yes" {
				t.Fatalf("upstream header X-Ttyd not forwarded, got headers %v", rec.Header())
			}

			body, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("unexpected error reading body: %v", err)
			}

			if string(body) != "ttyd payload" {
				t.Fatalf("body = %q, want %q", string(body), "ttyd payload")
			}
		})
	}
}

// TestWorkbenchTerminalShell_UnreachableTarget points the resolver at an address nothing is
// listening on and confirms the dial-stage failure is surfaced as 503, not the generic 502
// buildProxy's errorHandler otherwise falls back to.
func TestWorkbenchTerminalShell_UnreachableTarget(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error listening: %v", err)
	}

	addr := listener.Addr().String()

	err = listener.Close()
	if err != nil {
		t.Fatalf("unexpected error closing listener: %v", err)
	}

	authSvc := authedAs(uuid.New())

	workbenchSvc := &fakeTerminalShellTargetResolver{
		resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
			return "http://" + addr, nil
		},
	}

	handler := NewWorkbenchTerminalShellHandler(authSvc, alwaysMember(), workbenchSvc)

	rec := httptest.NewRecorder()
	newTerminalShellMux(handler).ServeHTTP(rec, terminalShellRequest(uuid.New(), "", "token-1"))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}
