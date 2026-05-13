package vaults_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (v *VaultsImpl) GetVault(ctx context.Context, req *pb.GetVault_Request) (*pb.GetVault_Response, error) {
	vaultID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	vault, err := v.vaultSvc.GetVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "get vault")
	}

	resp := &pb.GetVault_Response{
		Id:    vault.Uuid.String(),
		Name:  vault.Name,
		DbUrl: vault.CouchDBURL,
	}
	return resp, nil
}
