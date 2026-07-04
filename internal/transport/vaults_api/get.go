package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
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

	s3InstanceId := ""
	if vault.S3InstanceUuid != nil {
		s3InstanceId = vault.S3InstanceUuid.String()
	}

	resp := &pb.GetVault_Response{
		Id:           vault.Uuid.String(),
		Name:         vault.Name,
		DbUrl:        vault.CouchDBURL,
		S3InstanceId: s3InstanceId,
		S3BucketName: vault.S3BucketName,
	}

	return resp, nil
}
