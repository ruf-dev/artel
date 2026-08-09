package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) DisablePostgresDatabase(
	ctx context.Context,
	req *pb.DisablePostgresDatabase_Request,
) (*pb.DisablePostgresDatabase_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = v.vaultSvc.DisablePostgresDatabase(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "disable postgres database")
	}

	resp := &pb.DisablePostgresDatabase_Response{
		Status: "not_enabled",
	}

	return resp, nil
}
