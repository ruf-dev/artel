package mom

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
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

func New(
	mcpDefinitions repository.McpDefinitionsRepo,
	mcpConnectors repository.McpConnectorsRepo,
	externalConns repository.ExternalConnectionRepo,
) *ServiceImpl {
	return &ServiceImpl{
		mcpDefinitions: mcpDefinitions,
		mcpConnectors:  mcpConnectors,
		externalConns:  externalConns,
		emailExecutor:  executors.NewEmailExecutor(),
		httpExecutor:   executors.NewHttpExecutor(),
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

func (s *ServiceImpl) ExecuteToolForKey(ctx context.Context, keyId uuid.UUID, toolName string, params map[string]interface{}) (string, error) {
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

			return s.dispatch(ctx, connector.ExternalConnectionUuid, tool.Action, params)
		}
	}

	return "", user_errors.McpToolNotFound
}

// ExecuteToolForConnection executes a tool against a specific external connection without an
// ownership check — callers (e.g. the tract engine, which has no user_context for
// webhook-triggered runs) must verify connection ownership themselves before calling this.
func (s *ServiceImpl) ExecuteToolForConnection(ctx context.Context, exConnUuid uuid.UUID, mcpName string, toolName string, params map[string]interface{}) (string, error) {
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

		return s.dispatch(ctx, exConnUuid, tool.Action, params)
	}

	return "", user_errors.McpToolNotFound
}

// ExecuteToolForUserConnection executes a tool against a client-supplied external
// connection, but only after verifying that the connection belongs to the authenticated
// caller. Used by the web UI Toolbox, where the external_connection_id is supplied by the
// client and must not be trusted blindly.
func (s *ServiceImpl) ExecuteToolForUserConnection(ctx context.Context, exConnUuid uuid.UUID, mcpName string, toolName string, params map[string]interface{}) (string, error) {
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

func (s *ServiceImpl) dispatch(ctx context.Context, exConnUuid uuid.UUID, action domain.ToolAction, params map[string]interface{}) (string, error) {
	exConn, err := s.externalConns.GetByID(ctx, exConnUuid)
	if err != nil {
		return "", rerrors.Wrap(err, "error getting external connection")
	}

	switch {
	case action.Imap != nil || action.Smtp != nil:
		return s.dispatchEmail(ctx, exConn, action, params)
	case action.Http != nil:
		return s.dispatchHttp(ctx, exConn, action, params)
	default:
		return "", user_errors.McpActionMissing
	}
}

func (s *ServiceImpl) dispatchEmail(ctx context.Context, exConn domain.ExternalConnection, action domain.ToolAction, params map[string]interface{}) (string, error) {
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

func (s *ServiceImpl) dispatchHttp(ctx context.Context, exConn domain.ExternalConnection, action domain.ToolAction, params map[string]interface{}) (string, error) {
	if exConn.Provider != action.Http.Credentials {
		return "", user_errors.McpCredentialsProviderMismatch
	}

	var secrets map[string]interface{}

	err := json.Unmarshal(exConn.CredentialsJSON, &secrets)
	if err != nil {
		return "", rerrors.Wrap(err, "error unmarshaling http credentials")
	}

	result, err := s.httpExecutor.Execute(ctx, action, secrets, params)
	if err != nil {
		return "", rerrors.Wrap(err, "error executing http tool")
	}

	return result, nil
}
