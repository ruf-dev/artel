package public_docs_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (p *PublicDocsImpl) GetDefaultVault(
	ctx context.Context, _ *pb.GetDefaultVault_Request,
) (*pb.GetDefaultVault_Response, error) {
	vault, err := p.publicDocsSvc.GetDefaultVault(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "get default vault")
	}

	resp := &pb.GetDefaultVault_Response{
		Id:   vault.Uuid.String(),
		Name: vault.Name,
		Slug: vault.Slug,
	}

	return resp, nil
}
