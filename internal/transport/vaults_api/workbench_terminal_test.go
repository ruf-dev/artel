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

// fakeTerminalAuthService is a hand-rolled fake of the handler's terminalAuthService dependency.
type fakeTerminalAuthService struct {
	validateTokenFunc func(ctx context.Context, token string) (domain.User, error)

	seenTokens []string
}

func (f *fakeTerminalAuthService) ValidateToken(ctx context.Context, token string) (domain.User, error) {
	f.seenTokens = append(f.seenTokens, token)

	return f.validateTokenFunc(ctx, token)
}

// fakeTerminalVaultMembers is a hand-rolled fake of the handler's terminalVaultMembers
// dependency.
type fakeTerminalVaultMembers struct {
	getFunc func(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error)
}

func (f *fakeTerminalVaultMembers) Get(
	ctx context.Context, vaultID, userID uuid.UUID,
) (domain.VaultMember, error) {
	return f.getFunc(ctx, vaultID, userID)
}

// fakeTerminalTargetResolver is a hand-rolled fake of the handler's terminalTargetResolver
// dependency.
type fakeTerminalTargetResolver struct {
	resolveFunc func(ctx context.Context, vaultID, userID uuid.UUID) (string, error)
}

func (f *fakeTerminalTargetResolver) ResolveTerminalTarget(
	ctx context.Context, vaultID, userID uuid.UUID,
) (string, error) {
	return f.resolveFunc(ctx, vaultID, userID)
}

// newTerminalMux registers handler under TerminalRoutePattern on a real http.ServeMux, so the
// tests exercise the same {vaultId}/{rest...} wildcard parsing production does rather than
// hand-setting path values.
func newTerminalMux(handler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(TerminalRoutePattern, handler)

	return mux
}

// terminalRequest builds a GET against the terminal route for vaultID with rest appended after
// ".../terminal/", optionally carrying an access_token cookie.
func terminalRequest(vaultID uuid.UUID, rest string, token string) *http.Request {
	target := "/api/vaults/workbench/" + vaultID.String() + "/terminal/" + rest

	req := httptest.NewRequest(http.MethodGet, target, nil)

	if token != "" {
		cookie := &http.Cookie{Name: middleware.AccessTokenCookieName, Value: token}
		req.AddCookie(cookie)
	}

	return req
}

// alwaysMember is a terminalVaultMembers fake that accepts every membership lookup.
func alwaysMember() *fakeTerminalVaultMembers {
	return &fakeTerminalVaultMembers{
		getFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{VaultUuid: vaultID, UserUuid: userID}, nil
		},
	}
}

// authedAs is a terminalAuthService fake that resolves any presented token to userUuid.
func authedAs(userUuid uuid.UUID) *fakeTerminalAuthService {
	return &fakeTerminalAuthService{
		validateTokenFunc: func(ctx context.Context, token string) (domain.User, error) {
			return domain.User{Uuid: userUuid}, nil
		},
	}
}

func TestWorkbenchTerminal_Unauthenticated(t *testing.T) {
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

	workbenchSvc := &fakeTerminalTargetResolver{
		resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
			t.Fatal("terminal target must not be resolved for an unauthenticated request")

			return "", nil
		},
	}

	handler := NewWorkbenchTerminalHandler(authSvc, members, workbenchSvc)

	rec := httptest.NewRecorder()
	newTerminalMux(handler).ServeHTTP(rec, terminalRequest(uuid.New(), "", ""))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	// A missing cookie is still passed through to ValidateToken as an empty token — that's what
	// keeps the no-auth local-dev mode working through this route.
	if len(authSvc.seenTokens) != 1 || authSvc.seenTokens[0] != "" {
		t.Fatalf("ValidateToken tokens = %q, want exactly one empty token", authSvc.seenTokens)
	}
}

func TestWorkbenchTerminal_NonMember(t *testing.T) {
	authSvc := authedAs(uuid.New())

	members := &fakeTerminalVaultMembers{
		getFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{}, user_errors.NotFound
		},
	}

	workbenchSvc := &fakeTerminalTargetResolver{
		resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
			t.Fatal("terminal target must not be resolved for a non-member")

			return "", nil
		},
	}

	handler := NewWorkbenchTerminalHandler(authSvc, members, workbenchSvc)

	rec := httptest.NewRecorder()
	newTerminalMux(handler).ServeHTTP(rec, terminalRequest(uuid.New(), "", "token-1"))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestWorkbenchTerminal_WorkbenchNotRunning(t *testing.T) {
	authSvc := authedAs(uuid.New())

	workbenchSvc := &fakeTerminalTargetResolver{
		resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
			return "", user_errors.WorkbenchNotRunning
		},
	}

	handler := NewWorkbenchTerminalHandler(authSvc, alwaysMember(), workbenchSvc)

	rec := httptest.NewRecorder()
	newTerminalMux(handler).ServeHTTP(rec, terminalRequest(uuid.New(), "", "token-1"))

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestWorkbenchTerminal_InvalidVaultId(t *testing.T) {
	authSvc := authedAs(uuid.New())

	workbenchSvc := &fakeTerminalTargetResolver{
		resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
			t.Fatal("terminal target must not be resolved for an unparseable vault id")

			return "", nil
		},
	}

	handler := NewWorkbenchTerminalHandler(authSvc, alwaysMember(), workbenchSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/vaults/workbench/not-a-uuid/terminal/", nil)

	rec := httptest.NewRecorder()
	newTerminalMux(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// TestWorkbenchTerminal_Proxies stands a httptest.Server in for the in-container ttyd server and
// confirms the route prefix is stripped, the query string survives, and the upstream response
// body/headers/status are forwarded back verbatim.
func TestWorkbenchTerminal_Proxies(t *testing.T) {
	tests := []struct {
		name     string
		rest     string
		wantPath string
	}{
		{name: "terminal root serves ttyd index", rest: "", wantPath: "/"},
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

			workbenchSvc := &fakeTerminalTargetResolver{
				resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
					return ttyd.URL, nil
				},
			}

			handler := NewWorkbenchTerminalHandler(authSvc, alwaysMember(), workbenchSvc)

			req := terminalRequest(uuid.New(), tt.rest+"?arg=xyz", "token-1")

			rec := httptest.NewRecorder()
			newTerminalMux(handler).ServeHTTP(rec, req)

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

// TestWorkbenchTerminal_UnreachableTarget points the resolver at an address nothing is listening
// on and confirms the dial-stage failure is surfaced as 503, not the generic 502 buildProxy's
// errorHandler otherwise falls back to.
func TestWorkbenchTerminal_UnreachableTarget(t *testing.T) {
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

	workbenchSvc := &fakeTerminalTargetResolver{
		resolveFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
			return "http://" + addr, nil
		},
	}

	handler := NewWorkbenchTerminalHandler(authSvc, alwaysMember(), workbenchSvc)

	rec := httptest.NewRecorder()
	newTerminalMux(handler).ServeHTTP(rec, terminalRequest(uuid.New(), "", "token-1"))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}
