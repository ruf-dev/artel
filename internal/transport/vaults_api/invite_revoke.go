package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) RevokeInviteLink(ctx context.Context, req *pb.RevokeInviteLink_Request) (*pb.RevokeInviteLink_Response, error) {
	inviteID, err := uuid.Parse(req.InviteId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse invite id")
	}

	err = v.vaultSvc.RevokeInviteLink(ctx, inviteID)
	if err != nil {
		return nil, rerrors.Wrap(err, "revoke invite link")
	}

	return &pb.RevokeInviteLink_Response{}, nil
}
