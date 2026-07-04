package notes_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (n *NotesImpl) SaveNote(ctx context.Context, req *pb.SaveNote_Request) (*pb.SaveNote_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	err = n.noteSvc.SaveNote(ctx, vaultID, req.Path, req.Content)
	if err != nil {
		return nil, rerrors.Wrap(err, "error saving note")
	}

	return &pb.SaveNote_Response{}, nil
}
