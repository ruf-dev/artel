// Package vaultbucket resolves the storage.BinaryStore backing a vault's non-markdown
// content — shared by the MCP vault-tool executor and the notes service's folder
// export/import/delete methods so the S3-vs-CouchDB-adapter decision lives in one place.
package vaultbucket

import (
	"context"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	s3client "github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/storage"
	"go.redsock.ru/rerrors"
)

// Resolve returns nil (not an error) when the vault has neither CouchDB-for-binaries nor a
// linked S3 bucket — callers must treat a nil bucket as "no binary storage available".
func Resolve(
	ctx context.Context,
	vault domain.Vault,
	client *couchdb.LiveSyncClient,
	s3Instances repository.S3Instances,
) (storage.BinaryStore, error) {
	switch {
	case vault.UseCouchDBForBinaries:
		return couchdb.NewBinaryStoreAdapter(client), nil
	case vault.S3InstanceUuid != nil:
		s3Instance, err := s3Instances.Get(ctx, *vault.S3InstanceUuid)
		if err != nil {
			return nil, rerrors.Wrap(err, "error getting s3 instance")
		}

		cfg := s3client.Config{
			Endpoint:  s3Instance.Endpoint,
			Region:    s3Instance.Region,
			AccessKey: s3Instance.AccessKey,
			SecretKey: s3Instance.SecretKey,
			UseSSL:    s3Instance.UseSSL,
			PathStyle: s3Instance.PathStyle,
		}

		bucket, err := s3client.New(cfg, vault.S3BucketName)
		if err != nil {
			return nil, rerrors.Wrap(err, "error initializing s3 client")
		}

		return bucket, nil
	default:
		return nil, nil
	}
}
