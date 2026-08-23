package vaults_api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/terminalauth"
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

	// authLinksMu guards authLinks: the WS relay goroutines (one pair per open terminal
	// connection, across every vault) and PendingTerminalAuthLink (called from the GetVault RPC
	// goroutine) all touch it concurrently.
	authLinksMu sync.Mutex
	// authLinks caches the most recently detected Claude CLI OAuth sign-in link per vault, so
	// GetVault can surface it without holding a live reference to any particular WS connection.
	// Single-instance-deployment assumption: this is process-local, in-memory state (see
	// docs/architecture.md — no horizontal-scaling/session-affinity setup is documented), an
	// explicit design choice rather than an oversight.
	authLinks map[uuid.UUID]string
}

func NewWorkbenchTerminalShellHandler(
	authSvc terminalAuthService, vaultMembers terminalVaultMembers, workbenchSvc terminalShellTargetResolver,
) *WorkbenchTerminalShellHandler {
	authLinks := make(map[uuid.UUID]string)

	return &WorkbenchTerminalShellHandler{
		authSvc:      authSvc,
		vaultMembers: vaultMembers,
		workbenchSvc: workbenchSvc,
		authLinks:    authLinks,
	}
}

// PendingTerminalAuthLink returns the Claude CLI OAuth sign-in link most recently detected on
// vaultID's terminal WS output, or "" if none has been seen yet on this process (a fresh
// detection overwrites any earlier one — see setPendingTerminalAuthLink).
func (h *WorkbenchTerminalShellHandler) PendingTerminalAuthLink(vaultID uuid.UUID) string {
	h.authLinksMu.Lock()
	defer h.authLinksMu.Unlock()

	return h.authLinks[vaultID]
}

func (h *WorkbenchTerminalShellHandler) setPendingTerminalAuthLink(vaultID uuid.UUID, link string) {
	h.authLinksMu.Lock()
	h.authLinks[vaultID] = link
	h.authLinksMu.Unlock()
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

	rest := r.PathValue(terminalShellRestPathValue)

	// ttyd's own WebSocket endpoint ("/ws") is the only path under this proxy carrying an
	// Upgrade request — everything else (the index page, static assets) is a plain HTTP request
	// httputil.ReverseProxy already handles correctly via its own automatic WebSocket passthrough.
	// Checking the request's own Upgrade/Connection headers (via the canonical
	// websocket.IsWebSocketUpgrade helper) rather than hardcoding rest == "ws" is what buildProxy
	// itself relies on implicitly, and stays correct even if ttyd's own client ever serves its
	// WebSocket off a different path.
	if websocket.IsWebSocketUpgrade(r) {
		h.serveTerminalWebsocket(w, r, target, vaultUuid, rest)

		return
	}

	proxy, err := h.buildProxy(target, rest)
	if err != nil {
		log.Error().Err(err).Str("vault_id", vaultUuid.String()).Msg("error building workbench terminal shell proxy")
		http.Error(w, "error resolving workbench terminal", http.StatusInternalServerError)

		return
	}

	proxy.ServeHTTP(w, r)
}

// terminalShellWsUpgrader upgrades the client side of the terminal-shell WebSocket relay.
// Subprotocols must list "tty" — that's the subprotocol ttyd 1.7.7's own web client opens the
// connection with, and gorilla's Upgrader only completes the handshake carrying a subprotocol the
// client offered if the server explicitly lists it here (see Upgrader.Subprotocols' doc);
// otherwise the negotiated subprotocol is silently dropped.
var terminalShellWsUpgrader = websocket.Upgrader{
	Subprotocols: []string{"tty"},
	CheckOrigin:  allowTerminalShellWsOrigin,
}

// allowTerminalShellWsOrigin accepts every origin: this handler has already authenticated the
// caller and checked vault membership earlier in ServeHTTP, and the browser client always
// same-origins here through the artel UI — nothing extra to check at the WS layer.
func allowTerminalShellWsOrigin(_ *http.Request) bool {
	return true
}

// serveTerminalWebsocket relays the terminal-shell WebSocket upgrade end-to-end: it dials ttyd's
// own "/ws" endpoint first (so a dial failure can still be reported as a normal HTTP error status,
// the same way buildProxy's errorHandler does for the non-WS path), then upgrades the client
// connection, then pumps both directions verbatim while tapping upstream output for the Claude
// CLI's OAuth sign-in link.
func (h *WorkbenchTerminalShellHandler) serveTerminalWebsocket(
	w http.ResponseWriter, r *http.Request, target string, vaultUuid uuid.UUID, rest string,
) {
	upstreamURL, err := terminalShellWebsocketURL(target, rest, r.URL.RawQuery)
	if err != nil {
		log.Error().Err(err).Str("vault_id", vaultUuid.String()).Msg("error building workbench terminal shell websocket target")
		http.Error(w, "error resolving workbench terminal", http.StatusInternalServerError)

		return
	}

	dialHeader := http.Header{}
	dialHeader.Set("Sec-WebSocket-Protocol", "tty")

	upstreamConn, _, err := websocket.DefaultDialer.DialContext(r.Context(), upstreamURL, dialHeader)
	if err != nil {
		log.Error().Err(err).Str("vault_id", vaultUuid.String()).Str("target", target).
			Msg("error dialing workbench terminal shell websocket")

		// Mirrors buildProxy's errorHandler: a failed dial (connection refused, no route, dial
		// timeout — the workbench container itself unreachable) is reported distinctly from a
		// generic failure so the frontend can tell the user which kind of problem this is.
		var opErr *net.OpError
		if errors.As(err, &opErr) && opErr.Op == "dial" {
			http.Error(w, "workbench terminal target is unreachable", http.StatusServiceUnavailable)

			return
		}

		http.Error(w, "error resolving workbench terminal", http.StatusBadGateway)

		return
	}

	clientConn, err := terminalShellWsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Str("vault_id", vaultUuid.String()).Msg("error upgrading workbench terminal shell client websocket")

		closeErr := upstreamConn.Close()
		if closeErr != nil {
			log.Debug().Err(closeErr).Str("vault_id", vaultUuid.String()).Msg("workbench terminal shell upstream websocket close")
		}

		return
	}

	h.relayTerminalWebsocket(vaultUuid, clientConn, upstreamConn)
}

// terminalShellWebsocketURL builds the ws(s)://... URL to dial ttyd's own WebSocket endpoint
// directly, mirroring buildProxy's target/path resolution (route prefix stripped, same
// per-workbench target) but for a WebSocket dial rather than an HTTP reverse-proxy request.
func terminalShellWebsocketURL(target, rest, rawQuery string) (string, error) {
	targetUrl, err := url.Parse(target)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing workbench terminal shell websocket target url")
	}

	scheme := "ws"
	if targetUrl.Scheme == "https" {
		scheme = "wss"
	}

	targetUrl.Scheme = scheme
	targetUrl.Path = "/" + strings.TrimPrefix(rest, "/")
	targetUrl.RawQuery = rawQuery

	return targetUrl.String(), nil
}

// ttydOutputCommand is ttyd's Command.OUTPUT byte: the first byte of every server->client message
// that carries raw pty output (see this file's package-level doc references / the task's ttyd
// protocol notes — SET_WINDOW_TITLE='1' and SET_PREFERENCES='2' are the only other commands ttyd
// 1.7.7 sends, and neither carries pty bytes worth scanning).
const ttydOutputCommand = '0'

// relayTerminalWebsocket pumps ttyd's binary WebSocket protocol verbatim in both directions
// between the authenticated client and the upstream ttyd server, while tapping (never altering)
// every upstream->client OUTPUT frame through a fresh terminalauth.Parser to detect the Claude
// CLI's OAuth sign-in link appearing in the raw pty output. Blocks until either side's connection
// errors or closes, then force-closes both so the other direction's blocked read unblocks too.
func (h *WorkbenchTerminalShellHandler) relayTerminalWebsocket(vaultUuid uuid.UUID, clientConn, upstreamConn *websocket.Conn) {
	parser := terminalauth.NewParser()

	observeUpstreamOutput := func(messageType int, payload []byte) {
		h.observeTerminalOutput(vaultUuid, parser, messageType, payload)
	}

	errc := make(chan error, 2)

	go func() {
		errc <- relayWebsocketMessages(upstreamConn, clientConn, observeUpstreamOutput)
	}()

	go func() {
		errc <- relayWebsocketMessages(clientConn, upstreamConn, nil)
	}()

	<-errc

	closeErr := clientConn.Close()
	if closeErr != nil {
		log.Debug().Err(closeErr).Str("vault_id", vaultUuid.String()).Msg("workbench terminal shell client websocket close")
	}

	closeErr = upstreamConn.Close()
	if closeErr != nil {
		log.Debug().Err(closeErr).Str("vault_id", vaultUuid.String()).Msg("workbench terminal shell upstream websocket close")
	}

	<-errc
}

// relayWebsocketMessages reads messages from src and writes each one verbatim (same message type,
// same payload bytes) to dst until either side errors, returning that error. observe, when
// non-nil, is called with every message's exact type and payload before it is forwarded — used to
// tap upstream output without altering what the client receives.
func relayWebsocketMessages(src, dst *websocket.Conn, observe func(messageType int, payload []byte)) error {
	for {
		messageType, payload, err := src.ReadMessage()
		if err != nil {
			return err
		}

		if observe != nil {
			observe(messageType, payload)
		}

		err = dst.WriteMessage(messageType, payload)
		if err != nil {
			return err
		}
	}
}

// observeTerminalOutput feeds an upstream OUTPUT frame's pty payload into parser and, on a newly
// detected Claude CLI OAuth sign-in link, caches it for vaultUuid. Ignores every other ttyd
// command frame (SET_WINDOW_TITLE, SET_PREFERENCES) — only OUTPUT carries pty bytes worth
// scanning.
func (h *WorkbenchTerminalShellHandler) observeTerminalOutput(
	vaultUuid uuid.UUID, parser *terminalauth.Parser, _ int, payload []byte,
) {
	if len(payload) == 0 || payload[0] != ttydOutputCommand {
		return
	}

	link, found := parser.Feed(payload[1:])
	if !found {
		return
	}

	h.setPendingTerminalAuthLink(vaultUuid, link)
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
