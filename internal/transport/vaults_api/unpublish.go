package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) UnpublishVault(ctx context.Context, req *pb.UnpublishVault_Request) (*pb.UnpublishVault_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = v.vaultSvc.UnpublishVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "unpublish vault")
	}

	vault, err := v.vaultSvc.GetVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "get vault")
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

	resp := &pb.UnpublishVault_Response{Vault: item}

	return resp, nil
}
