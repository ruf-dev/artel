package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) DeleteWorkbench(
	ctx context.Context, req *pb.DeleteWorkbench_Request,
) (*pb.DeleteWorkbench_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = v.workbenchSvc.DeleteWorkbench(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "delete workbench")
	}

	resp := &pb.DeleteWorkbench_Response{
		Status: string(domain.WorkbenchStatusRemoved),
	}

	return resp, nil
}
