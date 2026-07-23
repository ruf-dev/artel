package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) StopWorkbench(
	ctx context.Context, req *pb.StopWorkbench_Request,
) (*pb.StopWorkbench_Response, error) {
	if v.workbenchSvc == nil {
		return nil, user_errors.WorkbenchNotConfigured
	}

	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = v.workbenchSvc.StopWorkbench(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "stop workbench")
	}

	resp := &pb.StopWorkbench_Response{
		Status: string(domain.WorkbenchStatusStopped),
	}

	return resp, nil
}
