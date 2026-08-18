package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) CloseWorkbenchTerminalTab(
	ctx context.Context, req *pb.CloseWorkbenchTerminalTab_Request,
) (*pb.CloseWorkbenchTerminalTab_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = v.workbenchSvc.CloseTerminalTab(ctx, vaultID, req.TabId)
	if err != nil {
		return nil, rerrors.Wrap(err, "close terminal tab")
	}

	resp := &pb.CloseWorkbenchTerminalTab_Response{}

	return resp, nil
}
