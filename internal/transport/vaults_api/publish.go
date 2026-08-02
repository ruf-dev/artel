package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) PublishVault(ctx context.Context, req *pb.PublishVault_Request) (*pb.PublishVault_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	vault, err := v.vaultSvc.PublishVault(ctx, vaultID, req.Slug)
	if err != nil {
		return nil, rerrors.Wrap(err, "publish vault")
	}

	s3InstanceId := ""
	if vault.S3InstanceUuid != nil {
		s3InstanceId = vault.S3InstanceUuid.String()
	}

	item := &pb.VaultItem{
		Id:                    vault.Uuid.String(),
		Name:                  vault.Name,
		DbUrl:                 vault.CouchDBURL,
		LivesyncPassphrase:    vault.LiveSyncPassphrase,
		S3InstanceId:          s3InstanceId,
		S3BucketName:          vault.S3BucketName,
		UseCouchdbForBinaries: vault.UseCouchDBForBinaries,
		IsPublic:              vault.IsPublic,
		Slug:                  vault.Slug,
	}

	resp := &pb.PublishVault_Response{Vault: item}

	return resp, nil
}
