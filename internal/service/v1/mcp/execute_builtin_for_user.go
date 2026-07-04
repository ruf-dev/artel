package mcp

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/storage"
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

	var bucket storage.BinaryStore
	if vault.S3InstanceUuid != nil {
		s3Instance, err := s.s3Instances.Get(ctx, *vault.S3InstanceUuid)
		if err != nil {
			return "", rerrors.Wrap(err, "error getting s3 instance")
		}

		cfg := s3client.Config{
			Endpoint:  s3Instance.Endpoint,
			Region:    s3Instance.Region,
			AccessKey: s3Instance.AccessKey,
			SecretKey: s3Instance.SecretKey,
			UseSSL:    s3Instance.UseSSL,
			PathStyle: s3Instance.PathStyle,
		}
		s3cli, err := s3client.New(cfg, vault.S3BucketName)
		if err != nil {
			return "", rerrors.Wrap(err, "error initializing s3 client")
		}
		bucket = s3cli
	}

	result, err := s.vaultExecutor.Execute(ctx, toolName, client, bucket, params)
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
