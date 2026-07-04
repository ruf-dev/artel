package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) RemoveMember(
	ctx context.Context,
	req *pb.RemoveMember_Request,
) (*pb.RemoveMember_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse user id")
	}

	err = v.vaultSvc.RemoveMember(ctx, vaultID, userID)
	if err != nil {
		return nil, rerrors.Wrap(err, "remove member")
	}

	resp := &pb.RemoveMember_Response{}

	return resp, nil
}
