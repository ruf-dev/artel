package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) CreateInviteLink(
	ctx context.Context,
	req *pb.CreateInviteLink_Request,
) (*pb.CreateInviteLink_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	role := artel_q.VaultRole(req.Role)
	if role != artel_q.VaultRoleReader && role != artel_q.VaultRoleMaintainer {
		return nil, rerrors.New("role must be 'reader' or 'maintainer'")
	}

	invite, err := v.vaultSvc.CreateInviteLink(ctx, vaultID, role)
	if err != nil {
		return nil, rerrors.Wrap(err, "create invite link")
	}

	pbInvite := &pb.VaultInviteItem{
		Id:        invite.Uuid.String(),
		VaultId:   invite.VaultUuid.String(),
		Role:      string(invite.Role),
		Token:     invite.Token,
		Revoked:   invite.RevokedAt != nil,
		CreatedAt: invite.CreatedAt.String(),
	}

	return &pb.CreateInviteLink_Response{Invite: pbInvite}, nil
}
