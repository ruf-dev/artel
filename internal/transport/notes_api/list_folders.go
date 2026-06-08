package notes_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (n *NotesImpl) ListFolders(ctx context.Context, req *pb.ListFolders_Request) (*pb.ListFolders_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	folders, err := n.noteSvc.ListFolders(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list folders")
	}

	return &pb.ListFolders_Response{Folders: folders}, nil
}
