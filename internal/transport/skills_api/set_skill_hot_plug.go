package skills_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (s *SkillsImpl) SetSkillHotPlug(
	ctx context.Context, req *pb.SetSkillHotPlug_Request,
) (*pb.SetSkillHotPlug_Response, error) {
	vaultUuid, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	skill, err := s.skillsSvc.SetSkillHotPlug(ctx, vaultUuid, req.Slug, req.HotPlug)
	if err != nil {
		return nil, rerrors.Wrap(err, "error setting skill hot-plug state")
	}

	return &pb.SetSkillHotPlug_Response{Skill: skillToInfo(skill)}, nil
}
