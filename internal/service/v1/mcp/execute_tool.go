package mcp

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
	"github.com/ruf-dev/artel/internal/storage"
)

func (s *McpServiceImpl) ExecuteTool(ctx context.Context, keyCtx domain.McpKeyContext, toolName string, params map[string]interface{}) (domain.ToolExecResult, error) {
	if toolName == toolConnections {
		return s.listConnectedMoms(ctx, keyCtx.KeyUuid)
	}

	if toolName == toolConnectionsForTracts {
		return s.listConnectionsForTracts(ctx, keyCtx.UserUuid)
	}

	if executors.IsTractTool(toolName) {
		return s.executeTractTool(ctx, keyCtx.UserUuid, toolName, params)
	}

	client := couchdb.NewLiveSyncClient(keyCtx.CouchURL, keyCtx.CouchDb, keyCtx.CouchUser, keyCtx.CouchPass)

	var bucket storage.BinaryStore
	if keyCtx.S3 != nil {
		cfg := s3client.Config{
			Endpoint:  keyCtx.S3.Endpoint,
			Region:    keyCtx.S3.Region,
			AccessKey: keyCtx.S3.AccessKey,
			SecretKey: keyCtx.S3.SecretKey,
			UseSSL:    keyCtx.S3.UseSSL,
			PathStyle: keyCtx.S3.PathStyle,
		}
		s3cli, err := s3client.New(cfg, keyCtx.S3.Bucket)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "init s3 client")
		}
		bucket = s3cli
	}

	return s.vaultExecutor.Execute(ctx, toolName, client, bucket, params)
}

func (s *McpServiceImpl) executeTractTool(ctx context.Context, userUuid uuid.UUID, toolName string, params map[string]interface{}) (domain.ToolExecResult, error) {
	if s.tractSvc == nil {
		return domain.ToolExecResult{}, user_errors.TractServiceNotConfigured
	}

	uc := user_context.UserContext{UserUuid: userUuid}
	ctx = user_context.WithUserContext(ctx, uc)

	return s.tractExecutor.Execute(ctx, s.tractBaseCtx, s.tractSvc, toolName, params)
}
