package notes_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (n *NotesImpl) ListTags(ctx context.Context, req *pb.ListTags_Request) (*pb.ListTags_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	tags, err := n.noteSvc.ListTags(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list tags")
	}

	return &pb.ListTags_Response{Tags: tags}, nil
}
