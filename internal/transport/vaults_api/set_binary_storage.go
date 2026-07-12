package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) SetVaultBinaryStorage(
	ctx context.Context,
	req *pb.SetVaultBinaryStorage_Request,
) (*pb.SetVaultBinaryStorage_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = v.vaultSvc.SetUseCouchDBForBinaries(ctx, vaultID, req.UseCouchdb)
	if err != nil {
		return nil, rerrors.Wrap(err, "set vault binary storage")
	}

	return &pb.SetVaultBinaryStorage_Response{}, nil
}
