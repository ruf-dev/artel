package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) LinkS3Bucket(
	ctx context.Context,
	req *pb.LinkS3Bucket_Request,
) (*pb.LinkS3Bucket_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	var s3InstanceIDPtr *uuid.UUID
	if req.S3InstanceId != "" {
		var s3InstanceID uuid.UUID
		s3InstanceID, err = uuid.Parse(req.S3InstanceId)
		if err != nil {
			return nil, rerrors.Wrap(err, "parse s3 instance id")
		}
		s3InstanceIDPtr = &s3InstanceID
	}

	err = v.vaultSvc.LinkS3Bucket(ctx, vaultID, s3InstanceIDPtr, req.BucketName)
	if err != nil {
		return nil, rerrors.Wrap(err, "link s3 bucket")
	}

	return &pb.LinkS3Bucket_Response{}, nil
}
