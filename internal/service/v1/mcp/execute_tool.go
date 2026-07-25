package mcp

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	s3client "github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
	"github.com/ruf-dev/artel/internal/storage"
	"go.redsock.ru/rerrors"
)

func (s *ServiceImpl) ExecuteTool(
	ctx context.Context,
	keyCtx domain.McpKeyContext,
	toolName string,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	if toolName == toolConnections {
		return s.listConnectedMoms(ctx, keyCtx.KeyUuid)
	}

	if toolName == toolConnectionsForTracts {
		return s.listConnectionsForTracts(ctx, keyCtx.UserUuid)
	}

	if toolName == toolCreateCommunityConnector {
		return s.createCommunityConnector(ctx, keyCtx.UserUuid, params)
	}

	if executors.IsTractTool(toolName) {
		return s.executeTractTool(ctx, keyCtx.UserUuid, toolName, params)
	}

	if toolName == executors.ToolWriteFile {
		err := s.subscriptions.CheckStorageQuota(ctx, keyCtx.UserUuid)
		if err != nil {
			return domain.ToolExecResult{}, err
		}
	}

	client := couchdb.NewLiveSyncClient(keyCtx.CouchURL, keyCtx.CouchDb, keyCtx.CouchUser, keyCtx.CouchPass)

	var bucket storage.BinaryStore

	switch {
	case keyCtx.UseCouchDBForBinaries:
		bucket = couchdb.NewBinaryStoreAdapter(client)
	case keyCtx.S3 != nil:
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

	result, err := s.vaultExecutor.Execute(ctx, toolName, client, bucket, params)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error executing vault tool")
	}

	return result, nil
}

func (s *ServiceImpl) executeTractTool(
	ctx context.Context,
	userUuid uuid.UUID,
	toolName string,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	if s.tractSvc == nil {
		return domain.ToolExecResult{}, user_errors.TractServiceNotConfigured
	}

	uc := user_context.UserContext{UserUuid: userUuid}
	ctx = user_context.WithUserContext(ctx, uc)

	result, err := s.tractExecutor.Execute(ctx, s.tractBaseCtx, s.tractSvc, toolName, params)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error executing tract tool")
	}

	return result, nil
}
