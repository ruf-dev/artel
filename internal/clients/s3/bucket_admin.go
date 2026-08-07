package s3client

import (
	"context"

	"github.com/minio/minio-go/v7"

	"go.redsock.ru/rerrors"
)

// ListBuckets returns the names of every bucket visible to the credentials in cfg. Unlike Client,
// which is scoped to one bucket, this is a package-level admin helper for callers that need to
// enumerate buckets across an entire endpoint (e.g. test cleanup).
func ListBuckets(ctx context.Context, cfg Config) ([]string, error) {
	mc, err := newMinioClient(cfg)
	if err != nil {
		return nil, rerrors.Wrap(err, "constructing minio client")
	}

	buckets, err := mc.ListBuckets(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "listing buckets")
	}

	names := make([]string, 0, len(buckets))
	for _, b := range buckets {
		names = append(names, b.Name)
	}

	return names, nil
}

// DeleteBucket empties bucket (removing every object) and then removes the bucket itself — S3
// requires a bucket to be empty before it can be removed.
func DeleteBucket(ctx context.Context, cfg Config, bucket string) error {
	mc, err := newMinioClient(cfg)
	if err != nil {
		return rerrors.Wrap(err, "constructing minio client")
	}

	listOpts := minio.ListObjectsOptions{Recursive: true}
	objectsCh := mc.ListObjects(ctx, bucket, listOpts)

	removeOpts := minio.RemoveObjectsOptions{}

	for removeErr := range mc.RemoveObjects(ctx, bucket, objectsCh, removeOpts) {
		if removeErr.Err != nil {
			return rerrors.Wrap(removeErr.Err, "removing object during bucket empty")
		}
	}

	err = mc.RemoveBucket(ctx, bucket)
	if err != nil {
		return rerrors.Wrap(err, "removing bucket")
	}

	return nil
}
