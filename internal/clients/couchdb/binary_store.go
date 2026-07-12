package couchdb

import (
	"context"

	"github.com/ruf-dev/artel/internal/storage"
	"go.redsock.ru/rerrors"
)

// BinaryStoreAdapter implements storage.BinaryStore on top of a LiveSyncClient's existing
// binary-file methods, letting a vault route non-markdown MCP file operations to CouchDB
// instead of S3 when its "use CouchDB for binaries" setting is on.
type BinaryStoreAdapter struct {
	client *LiveSyncClient
}

var _ storage.BinaryStore = (*BinaryStoreAdapter)(nil)

func NewBinaryStoreAdapter(client *LiveSyncClient) *BinaryStoreAdapter {
	return &BinaryStoreAdapter{client: client}
}

func (a *BinaryStoreAdapter) Put(ctx context.Context, path string, content []byte, mimeType string) error {
	err := a.client.WriteFile(ctx, path, content)
	if err != nil {
		return rerrors.Wrap(err, "writing file to couchdb")
	}

	return nil
}

func (a *BinaryStoreAdapter) Get(ctx context.Context, path string) (storage.Object, error) {
	file, err := a.client.ReadFile(ctx, path)
	if err != nil {
		return storage.Object{}, rerrors.Wrap(err, "reading file from couchdb")
	}

	obj := storage.Object{
		Path:     path,
		RawBytes: file.RawBytes,
		MimeType: file.MimeType,
		Size:     file.Size,
	}

	return obj, nil
}

func (a *BinaryStoreAdapter) Stat(ctx context.Context, path string) (storage.ObjectEntry, error) {
	entry, err := a.client.StatFile(ctx, path)
	if err != nil {
		return storage.ObjectEntry{}, rerrors.Wrap(err, "stat file in couchdb")
	}

	objEntry := storage.ObjectEntry{
		Path:     entry.Path,
		Size:     entry.Size,
		MimeType: entry.MimeType,
		Mtime:    entry.Mtime,
	}

	return objEntry, nil
}

func (a *BinaryStoreAdapter) Delete(ctx context.Context, path string) error {
	err := a.client.DeleteFile(ctx, path)
	if err != nil {
		return rerrors.Wrap(err, "deleting file from couchdb")
	}

	return nil
}

func (a *BinaryStoreAdapter) Move(ctx context.Context, oldPath, newPath string) error {
	err := a.client.MoveFile(ctx, oldPath, newPath)
	if err != nil {
		return rerrors.Wrap(err, "moving file in couchdb")
	}

	return nil
}

func (a *BinaryStoreAdapter) List(ctx context.Context) ([]storage.ObjectEntry, error) {
	files, err := a.client.ListFiles(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "listing files from couchdb")
	}

	entries := make([]storage.ObjectEntry, len(files))

	for i, f := range files {
		entries[i] = storage.ObjectEntry{Path: f.Path, Size: f.Size, MimeType: f.MimeType, Mtime: f.Mtime}
	}

	return entries, nil
}
