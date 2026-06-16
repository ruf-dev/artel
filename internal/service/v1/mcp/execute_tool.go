package mcp

import (
	"context"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
)

func (s *McpServiceImpl) ExecuteTool(ctx context.Context, keyCtx domain.McpKeyContext, toolName string, params map[string]interface{}) (domain.ToolExecResult, error) {
	if toolName == toolConnections {
		return s.listConnectedMoms(ctx, keyCtx.KeyUuid)
	}

	client := couchdb.NewLiveSyncClient(keyCtx.CouchURL, keyCtx.CouchDb, keyCtx.CouchUser, keyCtx.CouchPass)
	return s.vaultExecutor.Execute(ctx, toolName, client, params)
}
