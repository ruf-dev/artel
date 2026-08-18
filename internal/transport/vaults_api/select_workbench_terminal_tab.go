package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) SelectWorkbenchTerminalTab(
	ctx context.Context, req *pb.SelectWorkbenchTerminalTab_Request,
) (*pb.SelectWorkbenchTerminalTab_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = v.workbenchSvc.SelectTerminalTab(ctx, vaultID, req.TabId)
	if err != nil {
		return nil, rerrors.Wrap(err, "select terminal tab")
	}

	resp := &pb.SelectWorkbenchTerminalTab_Response{}

	return resp, nil
}
