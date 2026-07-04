package notes_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (n *NotesImpl) MoveNote(ctx context.Context, req *pb.MoveNote_Request) (*pb.MoveNote_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	err = n.noteSvc.MoveNote(ctx, vaultID, req.OldPath, req.NewPath)
	if err != nil {
		return nil, rerrors.Wrap(err, "error moving note")
	}

	return &pb.MoveNote_Response{}, nil
}
