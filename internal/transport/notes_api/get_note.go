package notes_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (n *NotesImpl) GetNote(ctx context.Context, req *pb.GetNote_Request) (*pb.GetNote_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	doc, err := n.noteSvc.GetNote(ctx, vaultID, req.Path)
	if err != nil {
		return nil, rerrors.Wrap(err, "get note")
	}

	return &pb.GetNote_Response{Content: doc.Content}, nil
}
