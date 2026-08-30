package mom

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
	"go.redsock.ru/rerrors"
)

type ServiceImpl struct {
	mcpDefinitions repository.McpDefinitionsRepo
	mcpConnectors  repository.McpConnectorsRepo
	externalConns  repository.ExternalConnectionRepo
	emailExecutor  *executors.EmailExecutor
	httpExecutor   *executors.HttpExecutor
}

// Option configures a mom.ServiceImpl at construction time.
type Option func(*momConfig)

type momConfig struct {
	execOpts []executors.HttpExecutorOption
}

// WithToolHttpMiddleware registers an http.RoundTripper middleware scoped to one MoM tool
// ("<mcp>.<tool>") on the MoM http executor. Forwarded verbatim to
// executors.WithToolHttpMiddleware.
func WithToolHttpMiddleware(tool artel_q.McpToolName, mw executors.HttpMiddleware) Option {
	return func(cfg *momConfig) {
		execOpt := executors.WithToolHttpMiddleware(tool, mw)
		cfg.execOpts = append(cfg.execOpts, execOpt)
	}
}

// WithMcpHttpMiddleware registers an http.RoundTripper middleware scoped to every tool of one
// MoM on the MoM http executor. Forwarded verbatim to executors.WithMcpHttpMiddleware.
func WithMcpHttpMiddleware(mcp artel_q.McpName, mw executors.HttpMiddleware) Option {
	return func(cfg *momConfig) {
		execOpt := executors.WithMcpHttpMiddleware(mcp, mw)
		cfg.execOpts = append(cfg.execOpts, execOpt)
	}
}

func New(
	mcpDefinitions repository.McpDefinitionsRepo,
	mcpConnectors repository.McpConnectorsRepo,
	externalConns repository.ExternalConnectionRepo,
	opts ...Option,
) *ServiceImpl {
	cfg := &momConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	return &ServiceImpl{
		mcpDefinitions: mcpDefinitions,
		mcpConnectors:  mcpConnectors,
		externalConns:  externalConns,
		emailExecutor:  executors.NewEmailExecutor(),
		httpExecutor:   executors.NewHttpExecutor(cfg.execOpts...),
	}
}

func (s *ServiceImpl) ListToolsForKey(ctx context.Context, keyId uuid.UUID) ([]domain.McpToolDef, error) {
	connectors, err := s.mcpConnectors.ListByKey(ctx, keyId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing mcp connectors")
	}

	var tools []domain.McpToolDef

	for _, connector := range connectors {
		def, err := s.mcpDefinitions.Get(ctx, connector.McpName)
		if err != nil {
			return nil, rerrors.Wrap(err, "error getting mcp definition")
		}

		if !def.Valid {
			continue
		}

		tools = append(tools, def.V.Tools...)
	}

	return tools, nil
}

func (s *ServiceImpl) ExecuteToolForKey(
	ctx context.Context,
	keyId uuid.UUID,
	toolName string,
	params map[string]interface{},
) (string, error) {
	connectors, err := s.mcpConnectors.ListByKey(ctx, keyId)
	if err != nil {
		return "", rerrors.Wrap(err, "error listing mcp connectors")
	}

	for _, connector := range connectors {
		def, err := s.mcpDefinitions.Get(ctx, connector.McpName)
		if err != nil {
			return "", rerrors.Wrap(err, "error getting mcp definition")
		}

		if !def.Valid {
			continue
		}

		for _, tool := range def.V.Tools {
			if tool.ApiDescription.Name != toolName {
				continue
			}

			ident := executors.ToolIdent{McpName: connector.McpName, ToolName: tool.ApiDescription.Name}

			return s.dispatch(ctx, ident, connector.ExternalConnectionUuid, tool.Action, params)
		}
	}

	return "", user_errors.McpToolNotFound
}

// ExecuteToolForConnection executes a tool against a specific external connection without an
// ownership check — callers (e.g. the tract engine, which has no user_context for
// webhook-triggered runs) must verify connection ownership themselves before calling this.
func (s *ServiceImpl) ExecuteToolForConnection(
	ctx context.Context,
	exConnUuid uuid.UUID,
	mcpName string,
	toolName string,
	params map[string]interface{},
) (string, error) {
	def, err := s.mcpDefinitions.Get(ctx, mcpName)
	if err != nil {
		return "", rerrors.Wrap(err, "error getting mcp definition")
	}

	if !def.Valid {
		return "", user_errors.McpToolNotFound
	}

	for _, tool := range def.V.Tools {
		if tool.ApiDescription.Name != toolName {
			continue
		}

		ident := executors.ToolIdent{McpName: mcpName, ToolName: toolName}

		return s.dispatch(ctx, ident, exConnUuid, tool.Action, params)
	}

	return "", user_errors.McpToolNotFound
}

// ExecuteToolForUserConnection executes a tool against a client-supplied external
// connection, but only after verifying that the connection belongs to the authenticated
// caller. Used by the web UI Toolbox, where the external_connection_id is supplied by the
// client and must not be trusted blindly.
func (s *ServiceImpl) ExecuteToolForUserConnection(
	ctx context.Context,
	exConnUuid uuid.UUID,
	mcpName string,
	toolName string,
	params map[string]interface{},
) (string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return "", user_errors.Unauthenticated
	}

	exConn, err := s.externalConns.GetByID(ctx, exConnUuid)
	if err != nil {
		return "", rerrors.Wrap(err, "error getting external connection")
	}

	if exConn.UserUuid != uc.UserUuid {
		return "", user_errors.McpConnectionNotOwned
	}

	return s.ExecuteToolForConnection(ctx, exConnUuid, mcpName, toolName, params)
}

// ExecuteToolWithSecrets executes a MoM http tool's action directly against a caller-supplied
// secrets map, without requiring a persisted external_connections row — used to validate
// externally-supplied credentials against the real provider before anything is saved.
func (s *ServiceImpl) ExecuteToolWithSecrets(
	ctx context.Context,
	mcpName string,
	toolName string,
	secrets map[string]interface{},
	params map[string]interface{},
) (string, error) {
	def, err := s.mcpDefinitions.Get(ctx, mcpName)
	if err != nil {
		return "", rerrors.Wrap(err, "error getting mcp definition")
	}

	if !def.Valid {
		return "", user_errors.McpToolNotFound
	}

	for _, tool := range def.V.Tools {
		if tool.ApiDescription.Name != toolName {
			continue
		}

		if tool.Action.Http == nil {
			return "", user_errors.McpActionMissing
		}

		ident := executors.ToolIdent{McpName: mcpName, ToolName: toolName}

		result, err := s.httpExecutor.Execute(ctx, ident, tool.Action, secrets, params)
		if err != nil {
			return "", rerrors.Wrap(err, "error executing http tool")
		}

		return result, nil
	}

	return "", user_errors.McpToolNotFound
}

func (s *ServiceImpl) dispatch(
	ctx context.Context,
	ident executors.ToolIdent,
	exConnUuid uuid.UUID,
	action domain.ToolAction,
	params map[string]interface{},
) (string, error) {
	exConn, err := s.externalConns.GetByID(ctx, exConnUuid)
	if err != nil {
		return "", rerrors.Wrap(err, "error getting external connection")
	}

	switch {
	case action.Imap != nil || action.Smtp != nil:
		return s.dispatchEmail(ctx, exConn, action, params)
	case action.Http != nil:
		return s.dispatchHttp(ctx, ident, exConn, action, params)
	default:
		return "", user_errors.McpActionMissing
	}
}

func (s *ServiceImpl) dispatchEmail(
	ctx context.Context,
	exConn domain.ExternalConnection,
	action domain.ToolAction,
	params map[string]interface{},
) (string, error) {
	var creds domain.EmailCredentials

	err := json.Unmarshal(exConn.CredentialsJSON, &creds)
	if err != nil {
		return "", rerrors.Wrap(err, "error unmarshaling email credentials")
	}

	result, err := s.emailExecutor.Execute(ctx, action, creds, params)
	if err != nil {
		return "", rerrors.Wrap(err, "error executing email tool")
	}

	return result, nil
}

func (s *ServiceImpl) dispatchHttp(
	ctx context.Context,
	ident executors.ToolIdent,
	exConn domain.ExternalConnection,
	action domain.ToolAction,
	params map[string]interface{},
) (string, error) {
	if exConn.Provider != action.Http.Credentials {
		return "", user_errors.McpCredentialsProviderMismatch
	}

	var secrets map[string]interface{}

	// UseNumber preserves numeric credential fields (e.g. account/workspace IDs) as their
	// original digit string via json.Number, instead of decoding them into float64 and
	// losing precision/formatting when they're later interpolated into a URL.
	decoder := json.NewDecoder(bytes.NewReader(exConn.CredentialsJSON))
	decoder.UseNumber()

	err := decoder.Decode(&secrets)
	if err != nil {
		return "", rerrors.Wrap(err, "error unmarshaling http credentials")
	}

	result, err := s.httpExecutor.Execute(ctx, ident, action, secrets, params)
	if err != nil {
		return "", rerrors.Wrap(err, "error executing http tool")
	}

	return result, nil
}
