package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) DeleteVault(ctx context.Context, req *pb.DeleteVault_Request) (*pb.DeleteVault_Response, error) {
	vaultID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	// DeleteWorkbenchesForVault must run to completion before the vault delete, not alongside it:
	// the workbenches rows cascade on vault delete via ON DELETE CASCADE, but the Docker
	// containers/volumes do not — deleting the vault first would orphan every member's live
	// container (and live API key) with nothing in Postgres pointing at it anymore. Unlike
	// CreateVault's best-effort hook, a failure here must fail the whole call rather than
	// proceeding. See docs/workbench/01_data_model_and_lifecycle.md, "Vault deletion".
	err = v.workbenchSvc.DeleteWorkbenchesForVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "error deleting workbenches for vault")
	}

	err = v.vaultSvc.DeleteVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "delete vault")
	}

	resp := &pb.DeleteVault_Response{}

	return resp, nil
}
