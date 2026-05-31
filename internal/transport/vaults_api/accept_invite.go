package vaults_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (v *VaultsImpl) AcceptInvite(ctx context.Context, req *pb.AcceptInvite_Request) (*pb.AcceptInvite_Response, error) {
	err := v.vaultSvc.AcceptInvite(ctx, req.Token)
	if err != nil {
		return nil, rerrors.Wrap(err, "accept invite")
	}

	return &pb.AcceptInvite_Response{}, nil
}
