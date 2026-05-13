package vaults_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (v *VaultsImpl) CreateVault(ctx context.Context, req *pb.CreateVault_Request) (*pb.CreateVault_Response, error) {
	vault, err := v.vaultSvc.CreateVault(ctx, req.Name)
	if err != nil {
		return nil, rerrors.Wrap(err, "create vault")
	}

	resp := &pb.CreateVault_Response{
		Id:    vault.Uuid.String(),
		Name:  vault.Name,
		DbUrl: vault.CouchDBURL,
	}
	return resp, nil
}
