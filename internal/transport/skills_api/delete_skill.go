package skills_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (s *SkillsImpl) DeleteSkill(ctx context.Context, req *pb.DeleteSkill_Request) (*pb.DeleteSkill_Response, error) {
	vaultUuid, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	err = s.skillsSvc.DeleteSkill(ctx, vaultUuid, req.Slug)
	if err != nil {
		return nil, rerrors.Wrap(err, "error deleting skill")
	}

	return &pb.DeleteSkill_Response{}, nil
}
