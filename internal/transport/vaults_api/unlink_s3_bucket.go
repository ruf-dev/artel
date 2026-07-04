package vaults_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (v *VaultsImpl) UnlinkS3Bucket(ctx context.Context, req *pb.UnlinkS3Bucket_Request) (*pb.UnlinkS3Bucket_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = v.vaultSvc.UnlinkS3Bucket(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "unlink s3 bucket")
	}

	return &pb.UnlinkS3Bucket_Response{}, nil
}
