package skills_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

func (s *SkillsImpl) UpdateSkill(ctx context.Context, req *pb.UpdateSkill_Request) (*pb.UpdateSkill_Response, error) {
	vaultUuid, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	storageMode := domain.SkillStorageMode(req.StorageMode)

	skill, err := s.skillsSvc.UpdateSkill(ctx, vaultUuid, req.Slug, req.Name, req.Description, storageMode, req.Body)
	if err != nil {
		return nil, rerrors.Wrap(err, "error updating skill")
	}

	return &pb.UpdateSkill_Response{Skill: skillToInfo(skill)}, nil
}
