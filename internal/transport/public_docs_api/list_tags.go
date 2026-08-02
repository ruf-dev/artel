package public_docs_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (p *PublicDocsImpl) ListTags(
	ctx context.Context, req *pb.PublicDocsListTags_Request,
) (*pb.PublicDocsListTags_Response, error) {
	tags, err := p.publicDocsSvc.ListTags(ctx, req.Slug)
	if err != nil {
		return nil, rerrors.Wrap(err, "list tags")
	}

	resp := &pb.PublicDocsListTags_Response{Tags: tags}

	return resp, nil
}
