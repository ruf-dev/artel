package public_docs_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (p *PublicDocsImpl) GetNote(
	ctx context.Context, req *pb.PublicDocsGetNote_Request,
) (*pb.PublicDocsGetNote_Response, error) {
	doc, err := p.publicDocsSvc.GetNote(ctx, req.Slug, req.Path)
	if err != nil {
		return nil, rerrors.Wrap(err, "get note")
	}

	resp := &pb.PublicDocsGetNote_Response{Content: doc.Content}

	return resp, nil
}
