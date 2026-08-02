package public_docs_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (p *PublicDocsImpl) GetVaultBySlug(
	ctx context.Context, req *pb.GetVaultBySlug_Request,
) (*pb.GetVaultBySlug_Response, error) {
	vault, err := p.publicDocsSvc.GetVaultBySlug(ctx, req.Slug)
	if err != nil {
		return nil, rerrors.Wrap(err, "get vault by slug")
	}

	resp := &pb.GetVaultBySlug_Response{
		Id:   vault.Uuid.String(),
		Name: vault.Name,
	}

	return resp, nil
}
