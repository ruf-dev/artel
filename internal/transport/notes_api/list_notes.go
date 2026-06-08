package notes_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (n *NotesImpl) ListNotes(ctx context.Context, req *pb.ListNotes_Request) (*pb.ListNotes_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	noteEntries, err := n.noteSvc.ListNotes(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list notes")
	}

	items := make([]*pb.NoteItem, len(noteEntries))
	for i, entry := range noteEntries {
		items[i] = &pb.NoteItem{Path: entry.Path, Mtime: entry.Mtime}
	}

	return &pb.ListNotes_Response{Notes: items}, nil
}
