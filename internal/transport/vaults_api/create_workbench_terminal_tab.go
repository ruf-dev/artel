package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) CreateWorkbenchTerminalTab(
	ctx context.Context, req *pb.CreateWorkbenchTerminalTab_Request,
) (*pb.CreateWorkbenchTerminalTab_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	tab, err := v.workbenchSvc.CreateTerminalTab(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "create terminal tab")
	}

	resp := &pb.CreateWorkbenchTerminalTab_Response{
		Tab: &pb.TerminalTab{
			Id:     tab.ID,
			Name:   tab.Name,
			Active: tab.Active,
		},
	}

	return resp, nil
}
