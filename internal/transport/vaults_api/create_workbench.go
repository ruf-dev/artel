package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) CreateWorkbench(
	ctx context.Context, req *pb.CreateWorkbench_Request,
) (*pb.CreateWorkbench_Response, error) {
	if v.workbenchSvc == nil {
		return nil, user_errors.WorkbenchNotConfigured
	}

	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	wb, err := v.workbenchSvc.CreateWorkbench(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "create workbench")
	}

	resp := &pb.CreateWorkbench_Response{
		Status: string(wb.Status),
	}

	return resp, nil
}
