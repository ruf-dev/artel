package vaults_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (v *VaultsImpl) ListMembers(ctx context.Context, req *pb.ListMembers_Request) (*pb.ListMembers_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	members, err := v.vaultSvc.ListMembers(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list members")
	}

	pbMembers := make([]*pb.VaultMemberInfo, len(members))
	for i, m := range members {
		pbMembers[i] = &pb.VaultMemberInfo{
			Id:       m.Uuid.String(),
			UserId:   m.UserUuid.String(),
			Email:    m.Email,
			Username: m.Username,
			Role:     string(m.Role),
		}
	}

	return &pb.ListMembers_Response{Members: pbMembers}, nil
}
