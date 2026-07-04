package vaults_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) ListVaults(ctx context.Context, req *pb.ListVaults_Request) (*pb.ListVaults_Response, error) {
	vaults, err := v.vaultSvc.ListVaults(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vaults")
	}

	items := make([]*pb.VaultItem, 0, len(vaults))

	for _, vault := range vaults {
		s3InstanceId := ""
		if vault.S3InstanceUuid != nil {
			s3InstanceId = vault.S3InstanceUuid.String()
		}

		item := &pb.VaultItem{
			Id:                 vault.Uuid.String(),
			Name:               vault.Name,
			DbUrl:              vault.CouchDBURL,
			LivesyncPassphrase: vault.LiveSyncPassphrase,
			S3InstanceId:       s3InstanceId,
			S3BucketName:       vault.S3BucketName,
		}
		items = append(items, item)
	}

	resp := &pb.ListVaults_Response{
		Vaults: items,
	}

	return resp, nil
}
