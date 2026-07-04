package s3client

import (
	"context"

	"go.redsock.ru/rerrors"
)

// TestConnection checks connectivity/credentials against an S3-compatible endpoint without
// needing any particular bucket to exist yet — backs an admin "test connection" API.
func TestConnection(ctx context.Context, cfg Config) error {
	mc, err := newMinioClient(cfg)
	if err != nil {
		return rerrors.Wrap(err, "constructing minio client")
	}

	_, err = mc.ListBuckets(ctx)
	if err != nil {
		return rerrors.Wrap(err, "list buckets")
	}

	return nil
}
