package vaults_api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/ruf-dev/artel/internal/middleware"
	"github.com/ruf-dev/artel/internal/utils"
)

// authLinkForWsTest is split across two scripted upstream frames below (see
// TestWorkbenchTerminalShell_Websocket_RelaysAndDetectsAuthLink) to exercise the same
// cross-message-boundary case internal/terminalauth's own unit tests cover, but end to end through
// the real handler.
const authLinkForWsTest = "https://claude.com/cai/oauth/authorize?client_id=abc&redirect_uri=https://example.com"

// newFakeTtydWebsocketServer stands in for the in-container ttyd server: it upgrades exactly one
// connection and writes outputFrames to it in order (each already framed with its ttyd command
// byte, e.g. '0' for OUTPUT, '1' for SET_WINDOW_TITLE), then blocks on a single ReadMessage —
// handing whatever the client sends back onto clientMsg — before closing.
//
// advance, when non-nil, gates each frame write on a receive from that channel — without this, a
// gorilla WriteMessage call returns as soon as the payload is queued in the OS socket buffer, not
// once the client has actually read it, so a test asserting "state X hasn't happened yet" right
// after reading frame N off the client could otherwise race against the relay goroutine having
// already sped ahead and processed frame N+1 (or later) server-side. Pacing one frame per receive
// from advance, with the test only sending on it right before reading the corresponding frame off
// the client, keeps at most one frame in flight and makes those "not yet" assertions deterministic.
func newFakeTtydWebsocketServer(t *testing.T, outputFrames [][]byte, advance <-chan struct{}, clientMsg chan<- []byte) *httptest.Server {
	t.Helper()

	upgrader := websocket.Upgrader{
		Subprotocols: []string{"tty"},
		CheckOrigin:  allowTerminalShellWsOrigin,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("fake ttyd: error upgrading: %v", err)

			return
		}
		defer utils.CloseWithLog(conn, "fake ttyd upstream websocket")

		for _, frame := range outputFrames {
			if advance != nil {
				<-advance
			}

			writeErr := conn.WriteMessage(websocket.BinaryMessage, frame)
			if writeErr != nil {
				t.Errorf("fake ttyd: error writing frame: %v", writeErr)

				return
			}
		}

		_, msg, readErr := conn.ReadMessage()
		if readErr == nil {
			clientMsg <- msg
		}
	}

	return httptest.NewServer(http.HandlerFunc(handler))
}

// TestWorkbenchTerminalShell_Websocket_RelaysAndDetectsAuthLink dials the real handler's
// terminal-shell WS upgrade path end to end against a fake upstream ttyd server, and confirms:
//   - every frame the fake upstream sends reaches the client unmodified, in order, with its
//     original message type preserved (including a non-OUTPUT SET_WINDOW_TITLE frame, which must
//     pass through untouched even though it's never fed to the parser);
//   - PendingTerminalAuthLink stays empty until the OUTPUT frame completing the split auth link has
//     actually reached the client (proving detection happened, not just that it eventually would);
//   - PendingTerminalAuthLink returns the correct URL once that frame arrives;
//   - the client -> upstream direction also relays verbatim.
func TestWorkbenchTerminalShell_Websocket_RelaysAndDetectsAuthLink(t *testing.T) {
	titleFrame := append([]byte{'1'}, []byte("my terminal title")...)
	plainOutputFrame := append([]byte{'0'}, []byte("hello terminal\r\n")...)
	// The auth link's OSC-8 URI field, split across two OUTPUT frames mid-URL — mirrors a ttyd
	// OUTPUT frame boundary landing inside the link, same case internal/terminalauth/parser_test.go
	// covers directly.
	linkPartOneFrame := append([]byte{'0'}, []byte("prefix \x1b]8;;https://claude.com/cai/oauth/authorize?client_id=abc&redirect_uri")...)
	linkPartTwoFrame := append([]byte{'0'}, []byte("=https://example.com\x07 trailer")...)

	outputFrames := [][]byte{titleFrame, plainOutputFrame, linkPartOneFrame, linkPartTwoFrame}

	advance := make(chan struct{})
	clientMsg := make(chan []byte, 1)
	ttyd := newFakeTtydWebsocketServer(t, outputFrames, advance, clientMsg)
	defer ttyd.Close()

	vaultID := uuid.New()
	userID := uuid.New()

	workbenchSvc := &fakeTerminalShellTargetResolver{
		resolveFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (string, error) {
			return ttyd.URL, nil
		},
	}

	handler := NewWorkbenchTerminalShellHandler(authedAs(userID), alwaysMember(), workbenchSvc)

	proxy := httptest.NewServer(newTerminalShellMux(handler))
	defer proxy.Close()

	wsURL := "ws" + strings.TrimPrefix(proxy.URL, "http") +
		"/api/vaults/workbench/" + vaultID.String() + "/terminal-shell/ws"

	dialHeader := http.Header{}
	dialHeader.Set("Cookie", middleware.AccessTokenCookieName+"=token-1")

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, dialHeader)
	if err != nil {
		t.Fatalf("unexpected error dialing terminal-shell websocket: %v", err)
	}
	defer utils.CloseWithLog(clientConn, "terminal-shell websocket test client")

	// advanceAndReadFrame releases the fake upstream's next scripted write and reads the resulting
	// relayed frame off the client — see newFakeTtydWebsocketServer's doc comment for why the
	// pacing matters here.
	advanceAndReadFrame := func(want []byte) {
		t.Helper()

		advance <- struct{}{}

		msgType, payload, readErr := clientConn.ReadMessage()
		if readErr != nil {
			t.Fatalf("unexpected error reading relayed frame: %v", readErr)
		}

		if msgType != websocket.BinaryMessage {
			t.Fatalf("relayed message type = %d, want %d (BinaryMessage)", msgType, websocket.BinaryMessage)
		}

		if !bytes.Equal(payload, want) {
			t.Fatalf("relayed payload = %q, want %q", payload, want)
		}
	}

	advanceAndReadFrame(titleFrame)

	if got := handler.PendingTerminalAuthLink(vaultID); got != "" {
		t.Fatalf("PendingTerminalAuthLink after title frame = %q, want empty", got)
	}

	advanceAndReadFrame(plainOutputFrame)

	if got := handler.PendingTerminalAuthLink(vaultID); got != "" {
		t.Fatalf("PendingTerminalAuthLink after plain output frame = %q, want empty", got)
	}

	advanceAndReadFrame(linkPartOneFrame)

	// The link is still incomplete at this point (no terminator byte has arrived yet) — must not
	// be detected from the truncated half alone.
	if got := handler.PendingTerminalAuthLink(vaultID); got != "" {
		t.Fatalf("PendingTerminalAuthLink after truncated link half = %q, want empty", got)
	}

	advanceAndReadFrame(linkPartTwoFrame)

	// The frame completing the link has now reached the client over the same WS connection the
	// server tapped it on its way through: the relay's observe() call happens strictly before its
	// WriteMessage call in program order, and the payload cannot arrive at the client before that
	// WriteMessage call was made, so detection is guaranteed to have already happened by this
	// point — no polling/sleep needed.
	if got := handler.PendingTerminalAuthLink(vaultID); got != authLinkForWsTest {
		t.Fatalf("PendingTerminalAuthLink after completed link = %q, want %q", got, authLinkForWsTest)
	}

	// Confirm the client -> upstream direction also relays verbatim.
	outbound := []byte("ping-from-client")

	err = clientConn.WriteMessage(websocket.TextMessage, outbound)
	if err != nil {
		t.Fatalf("unexpected error writing client message: %v", err)
	}

	select {
	case got := <-clientMsg:
		if !bytes.Equal(got, outbound) {
			t.Fatalf("upstream received %q, want %q", got, outbound)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for upstream to receive the client's message")
	}
}

// TestWorkbenchTerminalShell_Websocket_NoLinkStaysEmpty confirms PendingTerminalAuthLink stays
// empty for a vault whose terminal output never contains the auth link.
func TestWorkbenchTerminalShell_Websocket_NoLinkStaysEmpty(t *testing.T) {
	outputFrames := [][]byte{append([]byte{'0'}, []byte("nothing interesting here\r\n")...)}

	clientMsg := make(chan []byte, 1)
	ttyd := newFakeTtydWebsocketServer(t, outputFrames, nil, clientMsg)
	defer ttyd.Close()

	vaultID := uuid.New()
	userID := uuid.New()

	workbenchSvc := &fakeTerminalShellTargetResolver{
		resolveFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (string, error) {
			return ttyd.URL, nil
		},
	}

	handler := NewWorkbenchTerminalShellHandler(authedAs(userID), alwaysMember(), workbenchSvc)

	proxy := httptest.NewServer(newTerminalShellMux(handler))
	defer proxy.Close()

	wsURL := "ws" + strings.TrimPrefix(proxy.URL, "http") +
		"/api/vaults/workbench/" + vaultID.String() + "/terminal-shell/ws"

	dialHeader := http.Header{}
	dialHeader.Set("Cookie", middleware.AccessTokenCookieName+"=token-1")

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, dialHeader)
	if err != nil {
		t.Fatalf("unexpected error dialing terminal-shell websocket: %v", err)
	}
	defer utils.CloseWithLog(clientConn, "terminal-shell websocket test client")

	_, _, err = clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("unexpected error reading relayed frame: %v", err)
	}

	if got := handler.PendingTerminalAuthLink(vaultID); got != "" {
		t.Fatalf("PendingTerminalAuthLink = %q, want empty", got)
	}
}

// lowerAuthLinkQuietWindow overrides the package-level authLinkQuietWindow var for the duration of
// one test, restoring the previous value via t.Cleanup so the override never leaks into other
// tests in this package. authLinkQuietWindow is a var rather than a const specifically so tests
// like this one can shrink it and exercise the "elapsed"/"not elapsed" clearing behavior without
// real multi-second sleeps.
func lowerAuthLinkQuietWindow(t *testing.T, window time.Duration) {
	t.Helper()

	previous := authLinkQuietWindow
	authLinkQuietWindow = window

	t.Cleanup(func() {
		authLinkQuietWindow = previous
	})
}

// dialTerminalShellWebsocket dials the terminal-shell WS route for vaultID against proxy and
// returns the client connection, failing the test on any error.
func dialTerminalShellWebsocket(t *testing.T, proxyURL string, vaultID uuid.UUID) *websocket.Conn {
	t.Helper()

	wsURL := "ws" + strings.TrimPrefix(proxyURL, "http") +
		"/api/vaults/workbench/" + vaultID.String() + "/terminal-shell/ws"

	dialHeader := http.Header{}
	dialHeader.Set("Cookie", middleware.AccessTokenCookieName+"=token-1")

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, dialHeader)
	if err != nil {
		t.Fatalf("unexpected error dialing terminal-shell websocket: %v", err)
	}

	return clientConn
}

// TestWorkbenchTerminalShell_Websocket_AuthLinkPersistsAcrossRedraws confirms that several frames
// redrawing the *same* auth link within the quiet window keep PendingTerminalAuthLink pointing at
// it, rather than treating terminalauth.Parser.Feed's found=false on a repeat redraw as a signal to
// clear the cache.
func TestWorkbenchTerminalShell_Websocket_AuthLinkPersistsAcrossRedraws(t *testing.T) {
	lowerAuthLinkQuietWindow(t, time.Hour)

	linkFrame := append([]byte{'0'}, []byte("sign in: "+authLinkForWsTest+"\r\n")...)
	redrawFrame := append([]byte{'0'}, []byte("sign in: "+authLinkForWsTest+"\r\n")...)

	outputFrames := [][]byte{linkFrame, redrawFrame, redrawFrame}

	advance := make(chan struct{})
	clientMsg := make(chan []byte, 1)
	ttyd := newFakeTtydWebsocketServer(t, outputFrames, advance, clientMsg)
	defer ttyd.Close()

	vaultID := uuid.New()
	userID := uuid.New()

	workbenchSvc := &fakeTerminalShellTargetResolver{
		resolveFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (string, error) {
			return ttyd.URL, nil
		},
	}

	handler := NewWorkbenchTerminalShellHandler(authedAs(userID), alwaysMember(), workbenchSvc)

	proxy := httptest.NewServer(newTerminalShellMux(handler))
	defer proxy.Close()

	clientConn := dialTerminalShellWebsocket(t, proxy.URL, vaultID)
	defer utils.CloseWithLog(clientConn, "terminal-shell websocket test client")

	for range outputFrames {
		advance <- struct{}{}

		_, _, err := clientConn.ReadMessage()
		if err != nil {
			t.Fatalf("unexpected error reading relayed frame: %v", err)
		}

		if got := handler.PendingTerminalAuthLink(vaultID); got != authLinkForWsTest {
			t.Fatalf("PendingTerminalAuthLink = %q, want %q", got, authLinkForWsTest)
		}
	}
}

// TestWorkbenchTerminalShell_Websocket_AuthLinkClearedAfterQuietWindow confirms a cached auth link
// is cleared once a frame with no link content arrives after authLinkQuietWindow has elapsed since
// the link was last affirmed, but is NOT cleared by a no-link frame arriving before the window
// elapses.
func TestWorkbenchTerminalShell_Websocket_AuthLinkClearedAfterQuietWindow(t *testing.T) {
	quietWindow := 50 * time.Millisecond
	lowerAuthLinkQuietWindow(t, quietWindow)

	linkFrame := append([]byte{'0'}, []byte("sign in: "+authLinkForWsTest+"\r\n")...)
	noLinkFrame := append([]byte{'0'}, []byte("still working...\r\n")...)

	// Two no-link frames: the first is read back well before the quiet window elapses (asserting
	// "not cleared yet"), then the test sleeps past the window before releasing the second (asserting
	// "cleared once the window has genuinely elapsed").
	outputFrames := [][]byte{linkFrame, noLinkFrame, noLinkFrame}

	advance := make(chan struct{})
	clientMsg := make(chan []byte, 1)
	ttyd := newFakeTtydWebsocketServer(t, outputFrames, advance, clientMsg)
	defer ttyd.Close()

	vaultID := uuid.New()
	userID := uuid.New()

	workbenchSvc := &fakeTerminalShellTargetResolver{
		resolveFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (string, error) {
			return ttyd.URL, nil
		},
	}

	handler := NewWorkbenchTerminalShellHandler(authedAs(userID), alwaysMember(), workbenchSvc)

	proxy := httptest.NewServer(newTerminalShellMux(handler))
	defer proxy.Close()

	clientConn := dialTerminalShellWebsocket(t, proxy.URL, vaultID)
	defer utils.CloseWithLog(clientConn, "terminal-shell websocket test client")

	// linkFrame: the link is detected and cached.
	advance <- struct{}{}

	_, _, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("unexpected error reading relayed frame: %v", err)
	}

	if got := handler.PendingTerminalAuthLink(vaultID); got != authLinkForWsTest {
		t.Fatalf("PendingTerminalAuthLink after link frame = %q, want %q", got, authLinkForWsTest)
	}

	// First noLinkFrame: arrives well before the quiet window elapses, so the entry must survive.
	advance <- struct{}{}

	_, _, err = clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("unexpected error reading relayed frame: %v", err)
	}

	if got := handler.PendingTerminalAuthLink(vaultID); got != authLinkForWsTest {
		t.Fatalf("PendingTerminalAuthLink after no-link frame before quiet window elapsed = %q, want %q (not cleared)",
			got, authLinkForWsTest)
	}

	// Sleep past the quiet window (measured from linkFrame's lastAffirmedAt), then send a second
	// no-link frame — this one must clear the cache.
	time.Sleep(quietWindow * 3)

	advance <- struct{}{}

	_, _, err = clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("unexpected error reading relayed frame: %v", err)
	}

	if got := handler.PendingTerminalAuthLink(vaultID); got != "" {
		t.Fatalf("PendingTerminalAuthLink after no-link frame past quiet window = %q, want empty (cleared)", got)
	}
}
