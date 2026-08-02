package public_docs_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (p *PublicDocsImpl) ListNotes(
	ctx context.Context, req *pb.PublicDocsListNotes_Request,
) (*pb.PublicDocsListNotes_Response, error) {
	noteEntries, err := p.publicDocsSvc.ListNotes(ctx, req.Slug)
	if err != nil {
		return nil, rerrors.Wrap(err, "list notes")
	}

	items := make([]*pb.PublicDocsNoteItem, len(noteEntries))
	for i, entry := range noteEntries {
		items[i] = &pb.PublicDocsNoteItem{Path: entry.Path, Mtime: entry.Mtime}
	}

	resp := &pb.PublicDocsListNotes_Response{Notes: items}

	return resp, nil
}
