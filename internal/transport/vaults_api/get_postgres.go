package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) GetPostgresDatabase(
	ctx context.Context,
	req *pb.GetPostgresDatabase_Request,
) (*pb.GetPostgresDatabase_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	db, err := v.vaultSvc.GetPostgresDatabase(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "get postgres database")
	}

	if !db.Valid {
		resp := &pb.GetPostgresDatabase_Response{
			Status: "not_enabled",
		}

		return resp, nil
	}

	resp := &pb.GetPostgresDatabase_Response{
		Status:       string(db.V.Status),
		ErrorMessage: db.V.ErrorMessage,
	}

	return resp, nil
}
