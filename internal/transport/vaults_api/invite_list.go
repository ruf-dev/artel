package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) ListInviteLinks(
	ctx context.Context,
	req *pb.ListInviteLinks_Request,
) (*pb.ListInviteLinks_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	invites, err := v.vaultSvc.ListInviteLinks(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list invite links")
	}

	pbInvites := make([]*pb.VaultInviteItem, len(invites))
	for i, inv := range invites {
		pbInvites[i] = &pb.VaultInviteItem{
			Id:        inv.Uuid.String(),
			VaultId:   inv.VaultUuid.String(),
			Role:      string(inv.Role),
			Token:     inv.Token,
			Revoked:   inv.RevokedAt != nil,
			CreatedAt: inv.CreatedAt.String(),
		}
	}

	return &pb.ListInviteLinks_Response{Invites: pbInvites}, nil
}
