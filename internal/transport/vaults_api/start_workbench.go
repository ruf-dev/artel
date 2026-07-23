package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) StartWorkbench(
	ctx context.Context, req *pb.StartWorkbench_Request,
) (*pb.StartWorkbench_Response, error) {
	if v.workbenchSvc == nil {
		return nil, user_errors.WorkbenchNotConfigured
	}

	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	authMode := domain.WorkbenchAuthMode(req.AuthMode)

	wb, err := v.workbenchSvc.StartWorkbench(ctx, vaultID, authMode)
	if err != nil {
		return nil, rerrors.Wrap(err, "start workbench")
	}

	resp := &pb.StartWorkbench_Response{
		Status:   string(wb.Status),
		AuthMode: string(wb.AuthMode),
	}

	return resp, nil
}
