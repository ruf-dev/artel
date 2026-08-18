package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) ListWorkbenchTerminalTabs(
	ctx context.Context, req *pb.ListWorkbenchTerminalTabs_Request,
) (*pb.ListWorkbenchTerminalTabs_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	tabs, err := v.workbenchSvc.ListTerminalTabs(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list terminal tabs")
	}

	pbTabs := make([]*pb.TerminalTab, len(tabs))
	for i, t := range tabs {
		pbTabs[i] = &pb.TerminalTab{
			Id:     t.ID,
			Name:   t.Name,
			Active: t.Active,
		}
	}

	return &pb.ListWorkbenchTerminalTabs_Response{Tabs: pbTabs}, nil
}
