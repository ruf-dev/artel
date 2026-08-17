package skills_api

import (
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
)

// skillToInfo maps a domain.Skill to the wire SkillInfo shape, omitting Body — callers that need
// the body use GetSkill's separate top-level field.
func skillToInfo(skill domain.Skill) *pb.SkillInfo {
	info := &pb.SkillInfo{
		Slug:        skill.Slug,
		Name:        skill.Name,
		Description: skill.Description,
		StorageMode: string(skill.StorageMode),
		IsHotPlug:   skill.IsHotPlug,
		IsSystem:    skill.IsSystem,
	}

	return info
}
