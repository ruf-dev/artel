package public_docs_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (p *PublicDocsImpl) ListFolders(
	ctx context.Context, req *pb.PublicDocsListFolders_Request,
) (*pb.PublicDocsListFolders_Response, error) {
	folders, err := p.publicDocsSvc.ListFolders(ctx, req.Slug)
	if err != nil {
		return nil, rerrors.Wrap(err, "list folders")
	}

	resp := &pb.PublicDocsListFolders_Response{Folders: folders}

	return resp, nil
}
