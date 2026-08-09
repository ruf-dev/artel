package vaults_api

import (
	"context"

	"github.com/rs/zerolog/log"
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

		postgresEnabled := false
		postgresStatus := "not_enabled"
		pgDB, pgErr := v.vaultSvc.GetPostgresDatabase(ctx, vault.Uuid)
		if pgErr != nil {
			// Supplementary data only — a lookup failure here must not fail ListVaults as a
			// whole, only skip enrichment for this vault's postgres fields.
			log.Error().Err(pgErr).Str("vault_id", vault.Uuid.String()).
				Msg("error getting postgres database status for vault")
		} else if pgDB.Valid {
			postgresEnabled = true
			postgresStatus = string(pgDB.V.Status)
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
			Role:                  vault.MyRole,
			PostgresEnabled:       postgresEnabled,
			PostgresStatus:        postgresStatus,
		}
		items = append(items, item)
	}

	resp := &pb.ListVaults_Response{
		Vaults: items,
	}

	return resp, nil
}
