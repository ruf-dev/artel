package vaults_api

import (
	"context"

	"github.com/rs/zerolog/log"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) CreateVault(ctx context.Context, req *pb.CreateVault_Request) (*pb.CreateVault_Response, error) {
	vault, err := v.vaultSvc.CreateVault(ctx, req.Name)
	if err != nil {
		return nil, rerrors.Wrap(err, "create vault")
	}

	// CreateWorkbench is a sibling call, not embedded inside VaultService, so the vault package
	// never gains a compile-time dependency on Docker (mirrors LinkS3Bucket) — see
	// docs/workbench/01_data_model_and_lifecycle.md, "Hook points". Failure here must not roll
	// back or fail vault creation: a vault without a workbench is a recoverable, retryable state
	// (CreateWorkbench is idempotent on vault_id), so we only log clearly.
	if v.workbenchSvc != nil {
		_, wbErr := v.workbenchSvc.CreateWorkbench(ctx, vault.Uuid)
		if wbErr != nil {
			log.Error().Err(wbErr).Str("vault_id", vault.Uuid.String()).Msg("error creating workbench for vault")
		}
	}

	resp := &pb.CreateVault_Response{
		Id:    vault.Uuid.String(),
		Name:  vault.Name,
		DbUrl: vault.CouchDBURL,
	}

	return resp, nil
}
