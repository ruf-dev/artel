package vaults_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (v *VaultsImpl) ListVaults(ctx context.Context, req *pb.ListVaults_Request) (*pb.ListVaults_Response, error) {
	vaults, err := v.vaultSvc.ListVaults(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vaults")
	}

	items := make([]*pb.VaultItem, 0, len(vaults))
	for _, vault := range vaults {
		item := &pb.VaultItem{
			Id:                 vault.Uuid.String(),
			Name:               vault.Name,
			DbUrl:              vault.CouchDBURL,
			LivesyncPassphrase: vault.LiveSyncPassphrase,
		}
		items = append(items, item)
	}

	resp := &pb.ListVaults_Response{
		Vaults: items,
	}
	return resp, nil
}
