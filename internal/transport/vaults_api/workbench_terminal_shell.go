package vaults_api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

const (
	// TerminalShellRoutePattern is the http.ServeMux pattern this handler is registered under in
	// internal/app/custom.go. Everything after ".../terminal-shell/" is forwarded verbatim to the
	// workbench's in-container ttyd server, which serves its own index page, static assets and
	// the "/ws" WebSocket endpoint off its own root — hence the trailing multi-wildcard.
	//
	// Distinct from TerminalRoutePattern (which stays pointed at the in-container chat bridge):
	// this is the raw interactive terminal — tmux windows, each running its own `claude` TUI, with
	// zero shared history with the chat bridge. It sits under the "/api/vaults/" prefix for the
	// same reason TerminalRoutePattern does — see that constant's doc comment.
	TerminalShellRoutePattern = "/api/vaults/workbench/{vaultId}/terminal-shell/{rest...}"

	// terminalShellVaultIdPathValue / terminalShellRestPathValue are the wildcard names in
	// TerminalShellRoutePattern, read back via http.Request.PathValue.
	terminalShellVaultIdPathValue = "vaultId"
	terminalShellRestPathValue    = "rest"
)

// terminalShellTargetResolver is the narrow subset of service.WorkbenchService this handler
// depends on — just the ttyd base-URL lookup, not the whole workbench lifecycle surface.
type terminalShellTargetResolver interface {
	ResolveTerminalShellTarget(ctx context.Context, vaultID, userID uuid.UUID) (string, error)
}

// WorkbenchTerminalShellHandler reverse-proxies a vault member's browser straight to the ttyd
// server running inside that vault's workbench container, giving the user a real interactive
// terminal on the workbench's tmux session — one tab per tmux window, each its own live `claude`
// TUI, entirely independent of the chat bridge WorkbenchTerminalHandler proxies to.
//
// It is registered with transport.AddHttpHandler, i.e. outside the gRPC auth interceptor chain,
// so it authenticates and authorizes each request itself — same as WorkbenchTerminalHandler and
// the other raw handlers in internal/app/custom.go (mcp_api's OAuth endpoints, gitlab_webhook).
type WorkbenchTerminalShellHandler struct {
	authSvc      terminalAuthService
	vaultMembers terminalVaultMembers
	workbenchSvc terminalShellTargetResolver
}

func NewWorkbenchTerminalShellHandler(
	authSvc terminalAuthService, vaultMembers terminalVaultMembers, workbenchSvc terminalShellTargetResolver,
) *WorkbenchTerminalShellHandler {
	return &WorkbenchTerminalShellHandler{
		authSvc:      authSvc,
		vaultMembers: vaultMembers,
		workbenchSvc: workbenchSvc,
	}
}

func (h *WorkbenchTerminalShellHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vaultUuid, err := uuid.Parse(r.PathValue(terminalShellVaultIdPathValue))
	if err != nil {
		http.Error(w, "invalid vault id", http.StatusBadRequest)

		return
	}

	user, err := h.authenticate(ctx, r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)

		return
	}

	_, err = h.vaultMembers.Get(ctx, vaultUuid, user.Uuid)
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)

		return
	}

	target, err := h.workbenchSvc.ResolveTerminalShellTarget(ctx, vaultUuid, user.Uuid)
	if err != nil {
		if rerrors.Is(err, user_errors.WorkbenchNotRunning) {
			http.Error(w, user_errors.WorkbenchNotRunning.Error(), http.StatusConflict)

			return
		}

		log.Error().Err(err).Str("vault_id", vaultUuid.String()).Msg("error resolving workbench terminal shell target")
		http.Error(w, "error resolving workbench terminal", http.StatusInternalServerError)

		return
	}

	proxy, err := h.buildProxy(target, r.PathValue(terminalShellRestPathValue))
	if err != nil {
		log.Error().Err(err).Str("vault_id", vaultUuid.String()).Msg("error building workbench terminal shell proxy")
		http.Error(w, "error resolving workbench terminal", http.StatusInternalServerError)

		return
	}

	proxy.ServeHTTP(w, r)
}

// authenticate resolves the caller from the access_token cookie, the same credential the
// grpc-gateway's CookieToMetadataAnnotator forwards as "authorization" metadata for regular
// RPCs. An absent cookie is deliberately passed through to ValidateToken as an empty token
// rather than short-circuited: it fails the session lookup exactly like an invalid one, while
// still honouring the no-auth local-dev mode, where ValidateToken returns the fixed dev user
// regardless of what was presented (see auth.Service.ValidateToken).
func (h *WorkbenchTerminalShellHandler) authenticate(ctx context.Context, r *http.Request) (domain.User, error) {
	token := ""

	cookie, err := r.Cookie(middleware.AccessTokenCookieName)
	if err == nil {
		token = cookie.Value
	}

	user, err := h.authSvc.ValidateToken(ctx, token)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error validating workbench terminal shell access token")
	}

	return user, nil
}

// buildProxy builds a ReverseProxy forwarding to target (the ttyd base URL) with the route's
// own prefix stripped, so ttyd sees the paths it serves off its own root ("/", "/ws",
// "/token", ...). A proxy is built per request rather than cached because target is
// per-workbench and changes whenever a container is restarted; this is cheap, since a nil
// Transport means every proxy shares http.DefaultTransport's connection pool.
//
// WebSocket upgrades need no special handling here: httputil.ReverseProxy detects an Upgrade
// request and switches to raw bidirectional copying itself, and it reads the upgrade type off
// the outbound request *after* Rewrite runs, so the rewritten request keeps it.
func (h *WorkbenchTerminalShellHandler) buildProxy(target string, rest string) (*httputil.ReverseProxy, error) {
	targetUrl, err := url.Parse(target)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing workbench terminal shell target url")
	}

	upstreamPath := "/" + strings.TrimPrefix(rest, "/")

	rewrite := func(pr *httputil.ProxyRequest) {
		pr.Out.URL.Scheme = targetUrl.Scheme
		pr.Out.URL.Host = targetUrl.Host
		pr.Out.URL.Path = upstreamPath
		pr.Out.URL.RawPath = ""
		pr.Out.Host = targetUrl.Host
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error().Err(err).Str("target", target).Msg("error proxying workbench terminal shell request")

		// http.Transport wraps a failed dial in *url.Error; errors.As unwraps through that to the
		// underlying *net.OpError regardless of the exact cause (connection refused, no route,
		// dial timeout). Surfacing this distinctly from a generic 502 lets the frontend tell the
		// user the workbench container itself is unreachable, rather than a vague proxy failure.
		var opErr *net.OpError
		if errors.As(err, &opErr) && opErr.Op == "dial" {
			http.Error(w, "workbench terminal target is unreachable", http.StatusServiceUnavailable)

			return
		}

		w.WriteHeader(http.StatusBadGateway)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite:      rewrite,
		ErrorHandler: errorHandler,
	}

	return proxy, nil
}
