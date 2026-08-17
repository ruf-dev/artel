package skills_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (s *SkillsImpl) ListSkills(ctx context.Context, req *pb.ListSkills_Request) (*pb.ListSkills_Response, error) {
	vaultUuid, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	skillList, err := s.skillsSvc.ListSkills(ctx, vaultUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing skills")
	}

	items := make([]*pb.SkillInfo, len(skillList))
	for i, sk := range skillList {
		items[i] = skillToInfo(sk)
	}

	return &pb.ListSkills_Response{Skills: items}, nil
}
