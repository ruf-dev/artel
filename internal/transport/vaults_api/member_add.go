package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"go.redsock.ru/rerrors"
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

	role := artel_q.VaultRole(req.Role)
	if role != artel_q.VaultRoleReader && role != artel_q.VaultRoleMaintainer {
		role = artel_q.VaultRoleReader
	}

	err = v.vaultSvc.AddMember(ctx, vaultID, userID, role)
	if err != nil {
		return nil, rerrors.Wrap(err, "add member")
	}

	resp := &pb.AddMember_Response{}

	return resp, nil
}
