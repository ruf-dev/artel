package skills_api

import (
	"testing"

	"github.com/ruf-dev/artel/internal/domain"
)

func TestSkillToInfo(t *testing.T) {
	skill := domain.Skill{
		Slug:        "my-skill",
		Name:        "My Skill",
		Description: "does things",
		StorageMode: domain.SkillStorageFreeform,
		Body:        "instructional body text — must not leak into SkillInfo",
		IsHotPlug:   true,
		IsSystem:    false,
	}

	info := skillToInfo(skill)

	if info.Slug != skill.Slug {
		t.Errorf("Slug = %q, want %q", info.Slug, skill.Slug)
	}

	if info.Name != skill.Name {
		t.Errorf("Name = %q, want %q", info.Name, skill.Name)
	}

	if info.Description != skill.Description {
		t.Errorf("Description = %q, want %q", info.Description, skill.Description)
	}

	if info.StorageMode != string(skill.StorageMode) {
		t.Errorf("StorageMode = %q, want %q", info.StorageMode, skill.StorageMode)
	}

	if info.IsHotPlug != skill.IsHotPlug {
		t.Errorf("IsHotPlug = %v, want %v", info.IsHotPlug, skill.IsHotPlug)
	}

	if info.IsSystem != skill.IsSystem {
		t.Errorf("IsSystem = %v, want %v", info.IsSystem, skill.IsSystem)
	}
}

func TestSkillToInfo_SystemSkill(t *testing.T) {
	skill := domain.Skill{
		Slug:        "skill-creator",
		Name:        "Skill Creator",
		StorageMode: domain.SkillStorageNone,
		IsHotPlug:   true,
		IsSystem:    true,
	}

	info := skillToInfo(skill)

	if !info.IsSystem {
		t.Errorf("IsSystem = false, want true for the system skill")
	}

	if !info.IsHotPlug {
		t.Errorf("IsHotPlug = false, want true for the always-hot-plug system skill")
	}
}
