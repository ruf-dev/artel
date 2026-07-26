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

	// DeleteWorkbench must run to completion before the vault delete, not alongside it: the
	// workbenches row cascades on vault delete via ON DELETE CASCADE, but the Docker
	// container/volume do not — deleting the vault first would orphan a live container (and a
	// live API key) with nothing in Postgres pointing at it anymore. Unlike CreateVault's
	// best-effort hook, a failure here must fail the whole call rather than proceeding.
	// See docs/workbench/01_data_model_and_lifecycle.md, "Vault deletion".
	err = v.workbenchSvc.DeleteWorkbench(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "error deleting workbench for vault")
	}

	err = v.vaultSvc.DeleteVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "delete vault")
	}

	resp := &pb.DeleteVault_Response{}

	return resp, nil
}
