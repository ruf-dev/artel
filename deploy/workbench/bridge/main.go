// Command workbench-bridge is the process a workbench container runs as PID 1 — see
// deploy/workbench/entrypoint.sh and ./README.md for the surrounding architecture.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"workbenchbridge/internal/authlogin"
	"workbenchbridge/internal/chatprotocol"
	"workbenchbridge/internal/claudecli"
	"workbenchbridge/internal/claudesettings"
	"workbenchbridge/internal/envdrop"
	"workbenchbridge/internal/hub"
	"workbenchbridge/internal/permissions"
)

// shutdownTimeout bounds how long the HTTP servers get to drain in-flight requests (there's
// essentially only ever the one open chat WebSocket) before shutdown gives up and exits anyway —
// promptness on `docker stop` matters more than a graceful drain here.
const shutdownTimeout = 5 * time.Second

// listenAddr is where the bridge serves the chat WebSocket. Two other already-finished tracks
// hardcode this exact port — internal/transport/telegram_webhook/session.go's bridgeWSPath
// constant and the frontend's useChatSession.ts — both also assuming path wsPath below. It is
// also the same port the old ttyd process used, so the backend's reverse proxy
// (internal/clients/workbenchdocker/container.go) needs no changes.
const listenAddr = ":7681"

// wsPath is where hub.Hub is mounted. See listenAddr's doc comment.
const wsPath = "/ws"

// historyDir is where the bridge persists each chat session as append-only JSONL, so a consumer
// reconnecting can be replayed accurate history and a restarted bridge process (container
// restart/redeploy) can recover it — see hub.New's doc comment. This must stay in sync with
// workspaceMountPath in internal/clients/workbenchdocker/client.go (the main repo, a different Go
// module this one can't import): it's the fixed path where the workbench's persistent named
// volume is mounted inside the container, which is what makes history survive the container
// being stopped and started again.
const historyDir = "/workspace/.chat-history"

func main() {
	envStore := envdrop.New("")

	settingsPath, err := claudesettings.Write("")
	if err != nil {
		log.Fatalf("workbench-bridge: error writing claude settings file: %v", err)
	}

	chatHub := hub.New(historyDir)
	broker := permissions.NewBroker(chatHub.Broadcast)

	hookServer := startHookServer(broker)
	chatServer := startChatServer(chatHub)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	runner := claudecli.NewRunner("", "", settingsPath, envStore)

	err = ensureAuthenticated(ctx, envStore)
	if err != nil {
		log.Printf("workbench-bridge: subscription login did not complete: %v", err)
	}

	authCompleteEvent := chatprotocol.Event{Type: chatprotocol.EventAuthComplete}
	chatHub.Broadcast(authCompleteEvent)

	runMainLoop(ctx, chatHub, runner, broker)

	shutdown(hookServer, chatServer)
}

// ensureAuthenticated blocks until the bridge has an OAuth token to run `claude` with, if
// WORKBENCH_AUTH_MODE=subscription_login and no token has been injected (or produced by a
// previous run of this same container) yet. api_key mode, or a subscription_login mode that
// already has a token dropped, returns immediately.
//
// The bridge no longer drives its own sign-in flow (see authlogin.WaitForLogin's doc comment) —
// it only waits for someone to have logged in (or to log in) via the container's interactive
// terminal, which writes claude's own credentials file independently of this process.
func ensureAuthenticated(ctx context.Context, envStore *envdrop.Store) error {
	values, err := envStore.Read()
	if err != nil {
		return err
	}

	if values[envdrop.AuthModeVar] != envdrop.AuthModeSubscriptionLogin {
		return nil
	}

	if values[envdrop.OauthTokenVar] != "" {
		return nil
	}

	return authlogin.WaitForLogin(ctx, envStore)
}

// runMainLoop is the bridge's steady-state loop, reached once authentication (if any was needed)
// has completed. It dispatches each turn asynchronously (see the EventUserMessage case) so that it
// keeps servicing Inbound() — in particular EventPermissionDecision — while a turn is in flight.
// A turn's PreToolUse hook call parks in permissions.Broker.resolve until this very loop processes
// the matching EventPermissionDecision and calls broker.Decide; if the loop instead blocked
// synchronously inside RunTurn for the whole turn (as it used to), it could never dequeue that
// decision, deadlocking until the broker's own DecisionWait timeout auto-denied the request out
// from under the loop.
//
// This does not risk two `claude` processes actually running concurrently: Runner.mu (see its doc
// comment) serialises real turn execution regardless of how many goroutines call RunTurn — a
// second concurrent call simply blocks on that mutex until the first turn's process exits.
// RunTurn itself never crashes the process on a bad turn (see its doc comment), so no extra
// recover() is needed around the dispatch here either.
func runMainLoop(ctx context.Context, chatHub *hub.Hub, runner *claudecli.Runner, broker *permissions.Broker) {
	for {
		select {
		case event := <-chatHub.Inbound():
			switch event.Type {
			case chatprotocol.EventUserMessage:
				chatHub.Broadcast(event)
				go runner.RunTurn(ctx, event.Text, chatHub.Broadcast)
			case chatprotocol.EventPermissionDecision:
				if broker.Decide(event.ID, event.Decision) {
					chatHub.Broadcast(event)
				}
			case chatprotocol.EventNewChat:
				runner.SetSessionId("")
				chatHub.Reset(event)
			default:
				log.Printf("workbench-bridge: ignoring unexpected inbound event type %q during normal operation", event.Type)
			}
		case <-ctx.Done():
			return
		}
	}
}

// startHookServer starts the PreToolUse hook endpoint claude's settings file (see
// claudesettings.Write) points at. Loopback-only — see permissions.ListenAddr's doc comment.
func startHookServer(broker *permissions.Broker) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(permissions.HookPath, broker)

	server := &http.Server{
		Addr:    permissions.ListenAddr,
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("workbench-bridge: hook server error: %v", err)
		}
	}()

	return server
}

// startChatServer starts the public-facing chat WebSocket server. chatHub itself answers both the
// WebSocket upgrade at wsPath and a plain-GET status response for anything else it's mounted on
// (see hub.Hub.ServeHTTP's doc comment) — mounted at "/" too, rather than adding a redundant
// /healthz, so a curious operator or a health check hitting any other path still gets that same
// status JSON instead of a 404.
func startChatServer(chatHub *hub.Hub) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(wsPath, chatHub)
	mux.Handle("/", chatHub)

	server := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("workbench-bridge: chat server error: %v", err)
		}
	}()

	return server
}

// shutdown closes both HTTP servers so `docker stop` doesn't have to wait out a SIGKILL timeout —
// this process is PID 1 (see entrypoint.sh), so no shell trap is watching over it the way the old
// ttyd-based entrypoint needed one to.
func shutdown(servers ...*http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	for _, server := range servers {
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("workbench-bridge: error shutting down server: %v", err)
		}
	}

	os.Exit(0)
}
