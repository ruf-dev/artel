package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (v *VaultsImpl) UpdateVaultPrompts(ctx context.Context, req *pb.UpdateVaultPrompts_Request) (*pb.UpdateVaultPrompts_Response, error) {
	id, err := uuid.Parse(req.GetVaultId())
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}
	err = v.vaultSvc.UpdatePrompts(ctx, id, req.GetPrompt(), req.GetUseSystemPrompt())
	if err != nil {
		return nil, rerrors.Wrap(err, "update vault prompts")
	}
	return &pb.UpdateVaultPrompts_Response{}, nil
}
