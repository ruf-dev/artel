package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) EnablePostgresDatabase(
	ctx context.Context,
	req *pb.EnablePostgresDatabase_Request,
) (*pb.EnablePostgresDatabase_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	db, err := v.vaultSvc.EnablePostgresDatabase(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "enable postgres database")
	}

	resp := &pb.EnablePostgresDatabase_Response{
		Status:       string(db.Status),
		ErrorMessage: db.ErrorMessage,
	}

	return resp, nil
}
