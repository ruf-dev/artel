package notes_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (n *NotesImpl) MoveFolder(ctx context.Context, req *pb.MoveFolder_Request) (*pb.MoveFolder_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	movedCount, err := n.noteSvc.MoveFolder(ctx, vaultID, req.OldPath, req.NewPath)
	if err != nil {
		return nil, rerrors.Wrap(err, "error moving folder")
	}

	resp := &pb.MoveFolder_Response{
		MovedCount: int32(movedCount),
	}

	return resp, nil
}
