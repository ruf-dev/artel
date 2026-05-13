package vaults_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (v *VaultsImpl) AddMember(ctx context.Context, req *pb.AddMember_Request) (*pb.AddMember_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse user id")
	}

	err = v.vaultSvc.AddMember(ctx, vaultID, userID)
	if err != nil {
		return nil, rerrors.Wrap(err, "add member")
	}

	resp := &pb.AddMember_Response{}
	return resp, nil
}
