package skills_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (s *SkillsImpl) GetSkill(ctx context.Context, req *pb.GetSkill_Request) (*pb.GetSkill_Response, error) {
	vaultUuid, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	skill, err := s.skillsSvc.GetSkillBody(ctx, vaultUuid, req.Slug)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting skill")
	}

	resp := &pb.GetSkill_Response{
		Skill: skillToInfo(skill),
		Body:  skill.Body,
	}

	return resp, nil
}
