package vaults_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (v *VaultsImpl) DeleteVault(ctx context.Context, req *pb.DeleteVault_Request) (*pb.DeleteVault_Response, error) {
	vaultID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = v.vaultSvc.DeleteVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "delete vault")
	}

	resp := &pb.DeleteVault_Response{}
	return resp, nil
}
