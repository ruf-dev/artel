package executors

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

// stubRoundTripper is the base of every wound chain under test: it records the request it was
// handed and answers with a fixed 200 / "{}" response.
type stubRoundTripper struct {
	gotReq *http.Request
}

func (s *stubRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	s.gotReq = req

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("{}")),
		Header:     make(http.Header),
	}

	return resp, nil
}

// roundTripFunc adapts a plain function to http.RoundTripper.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// tagMiddleware appends its name to trace and stamps an X-Mw header when its wrapper actually
// runs, so a test can assert both scoping and the LIFO ordering.
func tagMiddleware(name string, trace *[]string) HttpMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripFunc(func(req *http.Request) (*http.Response, error) {
			*trace = append(*trace, name)
			req.Header.Add("X-Mw", name)

			return next.RoundTrip(req)
		})
	}
}

func getAction() domain.ToolAction {
	httpAction := &domain.HttpAction{
		Method: http.MethodGet,
		Url:    "http://mw.test/x",
	}

	return domain.ToolAction{Http: httpAction}
}

func execTool(t *testing.T, e *HttpExecutor, mcp, tool string) {
	t.Helper()

	ident := ToolIdent{McpName: mcp, ToolName: tool}

	out, err := e.Execute(context.Background(), ident, getAction(), nil, nil)
	require.NoError(t, err)
	require.Equal(t, "{}", out)
}

func TestHttpExecutor_NoRegistrationsUseBaseClient(t *testing.T) {
	stub := &stubRoundTripper{}
	e := newHttpExecutor(stub)

	execTool(t, e, "trello", "create_card")

	require.Same(t, e.base, e.clientFor(ToolIdent{McpName: "trello", ToolName: "create_card"}))
	require.Empty(t, stub.gotReq.Header.Values("X-Mw"))
}

func TestHttpExecutor_ToolScopedRegistrationRoutesOnlyThatTool(t *testing.T) {
	var trace []string

	stub := &stubRoundTripper{}
	opt := WithToolHttpMiddleware(artel_q.McpToolNameTrellocreateCard, tagMiddleware("create-card-mw", &trace))
	e := newHttpExecutor(stub, opt)

	execTool(t, e, "trello", "create_card")
	require.Equal(t, []string{"create-card-mw"}, trace)
	require.Equal(t, "create-card-mw", stub.gotReq.Header.Get("X-Mw"))

	trace = nil

	execTool(t, e, "trello", "list_cards")
	require.Empty(t, trace, "a sibling tool of the same mcp must not route through a tool-scoped registration")
	require.Empty(t, stub.gotReq.Header.Values("X-Mw"))
}

func TestHttpExecutor_McpScopedRegistrationRoutesEveryToolOfThatMcp(t *testing.T) {
	var trace []string

	stub := &stubRoundTripper{}
	opt := WithMcpHttpMiddleware(artel_q.McpNameTrello, tagMiddleware("trello-mw", &trace))
	e := newHttpExecutor(stub, opt)

	execTool(t, e, "trello", "create_card")
	execTool(t, e, "trello", "move_card")
	require.Equal(t, []string{"trello-mw", "trello-mw"}, trace)

	trace = nil

	execTool(t, e, "gitlab", "add_comment")
	require.Empty(t, trace, "a tool of another mcp must not route through an mcp-scoped registration")
	require.Empty(t, stub.gotReq.Header.Values("X-Mw"))
}

func TestHttpExecutor_ToolAndMcpRegistrationsCompose(t *testing.T) {
	var trace []string

	stub := &stubRoundTripper{}
	mcpOpt := WithMcpHttpMiddleware(artel_q.McpNameTrello, tagMiddleware("mcp", &trace))
	toolOpt := WithToolHttpMiddleware(artel_q.McpToolNameTrellocreateCard, tagMiddleware("tool", &trace))
	e := newHttpExecutor(stub, mcpOpt, toolOpt)

	execTool(t, e, "trello", "create_card")

	require.Equal(t, []string{"tool", "mcp"}, trace, "tool-scoped wrapper is outermost, mcp-scoped runs inside it")
	require.Equal(t, []string{"tool", "mcp"}, stub.gotReq.Header.Values("X-Mw"))
}

func TestHttpExecutor_SameScopeRegistrationsAreLifo(t *testing.T) {
	var trace []string

	stub := &stubRoundTripper{}
	first := WithMcpHttpMiddleware(artel_q.McpNameTrello, tagMiddleware("A", &trace))
	second := WithMcpHttpMiddleware(artel_q.McpNameTrello, tagMiddleware("B", &trace))
	e := newHttpExecutor(stub, first, second)

	execTool(t, e, "trello", "create_card")

	require.Equal(t, []string{"B", "A"}, trace, "last-registered middleware for a scope runs first")
	require.Equal(t, []string{"B", "A"}, stub.gotReq.Header.Values("X-Mw"))
}

func TestHttpExecutor_ToolClientStillIncludesItsMcpMiddleware(t *testing.T) {
	var trace []string

	stub := &stubRoundTripper{}
	mcpOpt := WithMcpHttpMiddleware(artel_q.McpNameTrello, tagMiddleware("trello-mcp", &trace))
	toolOpt := WithToolHttpMiddleware(artel_q.McpToolNameTrellocreateCard, tagMiddleware("create-card", &trace))
	e := newHttpExecutor(stub, mcpOpt, toolOpt)

	execTool(t, e, "trello", "create_card")

	require.Contains(t, trace, "trello-mcp", "a tool with its own client still gets its mcp's mcp-scoped middleware")
	require.Equal(t, "trello-mcp", trace[len(trace)-1], "the mcp-scoped wrapper is innermost, so it runs last")
}
