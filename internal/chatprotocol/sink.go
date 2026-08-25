package chatprotocol

import "context"

// EventSink is the transport seam for a producer of chat events — today Simple Chat's in-process
// agent loop (internal/service/v1/simplechat), which produces Events and consumes permission
// decisions through it without knowing a WebSocket is on the other side.
//
// It lives here, rather than in the simplechat package, so that both service.SimpleChatService
// and the transport layer can name it without either importing an internal/service/v1
// implementation package. It is deliberately in its own file: events.go is mirrored
// byte-identically into deploy/workbench/bridge's standalone module, and this interface is not
// part of that mirrored wire contract.
//
// Implementations are per-connection: one sink per WebSocket, constructed fresh, never shared.
type EventSink interface {
	// Send delivers one outbound event to the client. It is called from the goroutine the turn
	// runs on, which is not the connection's read loop, so implementations must be safe to call
	// concurrently with that read loop.
	Send(event Event) error

	// AwaitPermissionDecision blocks until a permission_decision event carrying requestId
	// arrives from the client, or ctx is done. The connection's read loop is responsible for
	// routing an incoming decision to whichever call is parked on that id.
	AwaitPermissionDecision(ctx context.Context, requestId string) (PermissionDecision, error)
}
