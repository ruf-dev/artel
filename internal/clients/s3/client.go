package s3client

import (
	"bytes"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/storage"
)

type Config struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	UseSSL    bool
	PathStyle bool // required true for Garage (https://garagehq.deuxfleurs.fr/) and most self-hosted S3-compatible stores
}

type Client struct {
	mc     *minio.Client
	bucket string
}

var _ storage.BinaryStore = (*Client)(nil)

// New constructs a client scoped to one bucket. PathStyle:true maps to
// minio.Options{BucketLookup: minio.BucketLookupPath} — required for Garage.
func New(cfg Config, bucket string) (*Client, error) {
	mc, err := newMinioClient(cfg)
	if err != nil {
		return nil, rerrors.Wrap(err, "creating minio client")
	}

	c := &Client{
		mc:     mc,
		bucket: bucket,
	}
	return c, nil
}

func newMinioClient(cfg Config) (*minio.Client, error) {
	creds := credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, "")

	lookup := minio.BucketLookupAuto
	if cfg.PathStyle {
		lookup = minio.BucketLookupPath
	}

	opts := &minio.Options{
		Creds:        creds,
		Secure:       cfg.UseSSL,
		Region:       cfg.Region,
		BucketLookup: lookup,
	}

	mc, err := minio.New(cfg.Endpoint, opts)
	if err != nil {
		return nil, rerrors.Wrap(err, "constructing minio client")
	}
	return mc, nil
}

// EnsureBucket idempotently creates the bucket if it doesn't exist (BucketExists -> MakeBucket).
// Mirrors internal/service/v1/vault/vault.go's tolerance of CouchDbDatabaseAlreadyExists for
// adminClient.CreateDatabase — if MakeBucket fails because the bucket now exists (race), that's
// not an error either.
func (c *Client) EnsureBucket(ctx context.Context) error {
	exists, err := c.BucketExists(ctx)
	if err != nil {
		return rerrors.Wrap(err, "checking bucket existence")
	}
	if exists {
		return nil
	}

	makeBucketOpts := minio.MakeBucketOptions{}
	err = c.mc.MakeBucket(ctx, c.bucket, makeBucketOpts)
	if err != nil {
		exists, existsErr := c.mc.BucketExists(ctx, c.bucket)
		if existsErr == nil && exists {
			return nil
		}
		return rerrors.Wrap(err, "making bucket")
	}
	return nil
}

// BucketExists reports whether the bucket already exists, without creating it.
func (c *Client) BucketExists(ctx context.Context) (bool, error) {
	exists, err := c.mc.BucketExists(ctx, c.bucket)
	if err != nil {
		return false, rerrors.Wrap(err, "checking bucket exists")
	}
	return exists, nil
}

func (c *Client) Put(ctx context.Context, path string, content []byte, mimeType string) error {
	reader := bytes.NewReader(content)

	putOpts := minio.PutObjectOptions{ContentType: mimeType}
	_, err := c.mc.PutObject(ctx, c.bucket, path, reader, int64(len(content)), putOpts)
	if err != nil {
		return rerrors.Wrap(err, "putting object")
	}
	return nil
}

func (c *Client) Get(ctx context.Context, path string) (storage.Object, error) {
	getOpts := minio.GetObjectOptions{}
	obj, err := c.mc.GetObject(ctx, c.bucket, path, getOpts)
	if err != nil {
		return storage.Object{}, rerrors.Wrap(err, "getting object")
	}
	defer obj.Close()

	stat, err := obj.Stat()
	if err != nil {
		return storage.Object{}, rerrors.Wrap(err, "stating object")
	}

	rawBytes, err := io.ReadAll(obj)
	if err != nil {
		return storage.Object{}, rerrors.Wrap(err, "reading object body")
	}

	object := storage.Object{
		Path:     path,
		RawBytes: rawBytes,
		MimeType: stat.ContentType,
		Size:     stat.Size,
	}
	return object, nil
}

func (c *Client) Stat(ctx context.Context, path string) (storage.ObjectEntry, error) {
	statOpts := minio.StatObjectOptions{}
	stat, err := c.mc.StatObject(ctx, c.bucket, path, statOpts)
	if err != nil {
		return storage.ObjectEntry{}, rerrors.Wrap(err, "stating object")
	}

	entry := storage.ObjectEntry{
		Path:     path,
		Size:     stat.Size,
		MimeType: stat.ContentType,
		Mtime:    stat.LastModified.UnixMilli(),
	}
	return entry, nil
}

func (c *Client) Delete(ctx context.Context, path string) error {
	removeOpts := minio.RemoveObjectOptions{}
	err := c.mc.RemoveObject(ctx, c.bucket, path, removeOpts)
	if err != nil {
		return rerrors.Wrap(err, "removing object")
	}
	return nil
}

// Move copies the object to newPath then removes oldPath — minio-go has no native rename/move.
func (c *Client) Move(ctx context.Context, oldPath, newPath string) error {
	src := minio.CopySrcOptions{
		Bucket: c.bucket,
		Object: oldPath,
	}
	dst := minio.CopyDestOptions{
		Bucket: c.bucket,
		Object: newPath,
	}

	_, err := c.mc.CopyObject(ctx, dst, src)
	if err != nil {
		return rerrors.Wrap(err, "copying object to new path")
	}

	err = c.Delete(ctx, oldPath)
	if err != nil {
		return rerrors.Wrap(err, "removing object at old path")
	}
	return nil
}

func (c *Client) List(ctx context.Context) ([]storage.ObjectEntry, error) {
	listOpts := minio.ListObjectsOptions{Recursive: true}
	objectCh := c.mc.ListObjects(ctx, c.bucket, listOpts)

	var entries []storage.ObjectEntry
	for obj := range objectCh {
		if obj.Err != nil {
			return nil, rerrors.Wrap(obj.Err, "listing objects")
		}

		entry := storage.ObjectEntry{
			Path:     obj.Key,
			Size:     obj.Size,
			MimeType: obj.ContentType,
			Mtime:    obj.LastModified.UnixMilli(),
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
