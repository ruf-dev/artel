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

func main() {
	envStore := envdrop.New("")

	settingsPath, err := claudesettings.Write("")
	if err != nil {
		log.Fatalf("workbench-bridge: error writing claude settings file: %v", err)
	}

	chatHub := hub.New()
	broker := permissions.NewBroker(chatHub.Broadcast)

	hookServer := startHookServer(broker)
	chatServer := startChatServer(chatHub)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	runner := claudecli.NewRunner("", "", settingsPath, envStore)

	err = ensureAuthenticated(ctx, envStore, chatHub)
	if err != nil {
		log.Printf("workbench-bridge: subscription login did not complete: %v", err)
	}

	authCompleteEvent := chatprotocol.Event{Type: chatprotocol.EventAuthComplete}
	chatHub.Broadcast(authCompleteEvent)

	runMainLoop(ctx, chatHub, runner, broker)

	shutdown(hookServer, chatServer)
}

// ensureAuthenticated runs the subscription-login device flow to completion before the bridge
// accepts any user_message, if WORKBENCH_AUTH_MODE=subscription_login and no OAuth token has been
// injected (or produced by a previous run of this same container) yet. api_key mode, or a
// subscription_login mode that already has a token dropped, returns immediately.
func ensureAuthenticated(ctx context.Context, envStore *envdrop.Store, chatHub *hub.Hub) error {
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

	codes := make(chan string, 1)
	done := make(chan struct{})

	go relayAuthCodes(ctx, chatHub, codes, done)
	defer close(done)

	loginRunner := authlogin.NewRunner("")

	return loginRunner.Run(ctx, envStore, chatHub.Broadcast, codes)
}

// relayAuthCodes forwards every auth_code_submit event the hub receives to codes, for as long as
// the login flow (ensureAuthenticated) is running. Any other inbound event type is dropped: a
// consumer sending a user_message or permission_decision before login has completed has nothing
// to act on yet.
//
// The forwarding send is itself inside a select on done/ctx.Done(), not a bare `codes <- ...`:
// codes is a size-1 buffered channel only drained by authlogin.Runner.Run's own select loop, and
// once Run has returned (success, failure, or ctx cancellation) nothing reads it again. Without
// this, a code arriving in that narrow window would block the send forever — done and ctx.Done()
// are both already true by the time this goroutine could next be scheduled, but only if the send
// itself can observe them.
func relayAuthCodes(ctx context.Context, chatHub *hub.Hub, codes chan<- string, done <-chan struct{}) {
	for {
		select {
		case event := <-chatHub.Inbound():
			if event.Type == chatprotocol.EventAuthCodeSubmit {
				select {
				case codes <- event.Code:
				case <-done:
					return
				case <-ctx.Done():
					return
				}
			}
		case <-done:
			return
		case <-ctx.Done():
			return
		}
	}
}

// runMainLoop is the bridge's steady-state loop, reached once authentication (if any was needed)
// has completed. It runs exactly one claude turn at a time — RunTurn itself never crashes the
// process on a bad turn (see its doc comment) — so no extra recover() is needed here as long as
// this loop stays synchronous, which it deliberately does: a second concurrent turn would
// interleave two processes resuming the same session.
func runMainLoop(ctx context.Context, chatHub *hub.Hub, runner *claudecli.Runner, broker *permissions.Broker) {
	for {
		select {
		case event := <-chatHub.Inbound():
			switch event.Type {
			case chatprotocol.EventUserMessage:
				runner.RunTurn(ctx, event.Text, chatHub.Broadcast)
			case chatprotocol.EventPermissionDecision:
				broker.Decide(event.ID, event.Decision)
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
