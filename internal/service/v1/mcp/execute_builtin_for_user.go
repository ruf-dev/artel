package mcp

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// ExecuteBuiltinToolForUser runs a builtin (vault) tool as userUuid rather than through an
// MCP key context — used by the tract engine, which has no key (a tract step only carries a
// user). Builtin vault tools are inherently vault-scoped, so this resolves the user's first
// vault by membership. Ambiguous for users with more than one vault until tract action steps
// carry an explicit vault reference — a known v1 gap, flagged for the tract MCP-authoring
// phase to revisit.
func (s *McpServiceImpl) ExecuteBuiltinToolForUser(ctx context.Context, userUuid uuid.UUID, toolName string, params map[string]interface{}) (string, error) {
	memberVaults, err := s.vaults.ListByMembership(ctx, userUuid)
	if err != nil {
		return "", rerrors.Wrap(err, "error listing vaults for user")
	}
	if len(memberVaults) == 0 {
		return "", rerrors.Wrap(user_errors.NoVaultForBuiltinTool)
	}
	vault := memberVaults[0]

	couchInstance, err := s.couchInstances.Get(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return "", rerrors.Wrap(err, "error getting couch instance")
	}

	client := couchdb.NewLiveSyncClient(couchInstance.Url, vault.CouchDBName, couchInstance.Username, couchInstance.Password)

	result, err := s.vaultExecutor.Execute(ctx, toolName, client, params)
	if err != nil {
		return "", rerrors.Wrap(err, "error executing builtin tool")
	}

	if result.Text != "" {
		return result.Text, nil
	}
	if len(result.Data) > 0 {
		return base64.StdEncoding.EncodeToString(result.Data), nil
	}
	return "", nil
}
