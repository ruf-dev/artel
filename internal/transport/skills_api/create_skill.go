package skills_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

func (s *SkillsImpl) CreateSkill(ctx context.Context, req *pb.CreateSkill_Request) (*pb.CreateSkill_Response, error) {
	vaultUuid, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	storageMode := domain.SkillStorageMode(req.StorageMode)

	skill, err := s.skillsSvc.CreateSkill(ctx, vaultUuid, req.Name, req.Description, storageMode, req.Body, req.HotPlug)
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating skill")
	}

	return &pb.CreateSkill_Response{Skill: skillToInfo(skill)}, nil
}
