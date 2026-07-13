package notes_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (n *NotesImpl) DeleteFolder(ctx context.Context, req *pb.DeleteFolder_Request) (*pb.DeleteFolder_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	deletedCount, failedPaths, err := n.noteSvc.DeleteFolder(ctx, vaultID, req.Path)
	if err != nil {
		return nil, rerrors.Wrap(err, "error deleting folder")
	}

	resp := &pb.DeleteFolder_Response{
		DeletedCount: int32(deletedCount),
		FailedPaths:  failedPaths,
	}

	return resp, nil
}
