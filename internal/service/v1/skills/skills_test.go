package skills

import (
	"testing"

	"github.com/ruf-dev/artel/internal/domain"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "simple lowercase", in: "my skill", want: "my-skill"},
		{name: "mixed case", in: "Weekly Report Drafting", want: "weekly-report-drafting"},
		{name: "punctuation collapses to single dash", in: "foo!!bar??baz", want: "foo-bar-baz"},
		{name: "leading and trailing punctuation trimmed", in: "  -Hello World-  ", want: "hello-world"},
		{name: "digits preserved", in: "Q3 2024 Report", want: "q3-2024-report"},
		{name: "empty input falls back", in: "", want: "skill"},
		{name: "only punctuation falls back", in: "!!!", want: "skill"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slugify(tt.in)
			if got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestDedupeSlug(t *testing.T) {
	tests := []struct {
		name  string
		base  string
		taken map[string]bool
		want  string
	}{
		{
			name:  "no collision",
			base:  "my-skill",
			taken: map[string]bool{},
			want:  "my-skill",
		},
		{
			name:  "single collision suffixes -2",
			base:  "my-skill",
			taken: map[string]bool{"my-skill": true},
			want:  "my-skill-2",
		},
		{
			name:  "multiple collisions find next free suffix",
			base:  "my-skill",
			taken: map[string]bool{"my-skill": true, "my-skill-2": true, "my-skill-3": true},
			want:  "my-skill-4",
		},
		{
			name:  "collision against reserved system slug",
			base:  SkillCreatorSlug,
			taken: map[string]bool{SkillCreatorSlug: true},
			want:  SkillCreatorSlug + "-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dedupeSlug(tt.base, tt.taken)
			if got != tt.want {
				t.Errorf("dedupeSlug(%q, %v) = %q, want %q", tt.base, tt.taken, got, tt.want)
			}
		})
	}
}

func TestSkillPath(t *testing.T) {
	tests := []struct {
		name    string
		slug    string
		hotPlug bool
		want    string
	}{
		{name: "hot-plug goes to active folder", slug: "foo", hotPlug: true, want: ".skills/active/foo.md"},
		{name: "non-hot-plug goes to library folder", slug: "foo", hotPlug: false, want: ".skills/library/foo.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := skillPath(tt.slug, tt.hotPlug)
			if got != tt.want {
				t.Errorf("skillPath(%q, %v) = %q, want %q", tt.slug, tt.hotPlug, got, tt.want)
			}
		})
	}
}

func TestSkillFromNote(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		content string
		want    domain.Skill
	}{
		{
			name: "active skill parses hot-plug true and slug from path",
			path: ".skills/active/my-skill.md",
			content: "---\n" +
				"name: My Skill\n" +
				"description: does a thing\n" +
				"storage_mode: freeform_notes\n" +
				"---\n" +
				"body text",
			want: domain.Skill{
				Slug:        "my-skill",
				Name:        "My Skill",
				Description: "does a thing",
				StorageMode: domain.SkillStorageFreeform,
				Body:        "body text",
				IsHotPlug:   true,
				IsSystem:    false,
			},
		},
		{
			name: "library skill parses hot-plug false and slug from path",
			path: ".skills/library/other-skill.md",
			content: "---\n" +
				"name: Other Skill\n" +
				"description: does another thing\n" +
				"storage_mode: structured\n" +
				"---\n" +
				"other body",
			want: domain.Skill{
				Slug:        "other-skill",
				Name:        "Other Skill",
				Description: "does another thing",
				StorageMode: domain.SkillStorageStructured,
				Body:        "other body",
				IsHotPlug:   false,
				IsSystem:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := skillFromNote(tt.path, tt.content)
			if got != tt.want {
				t.Errorf("skillFromNote() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

// TestRenderNoteAndSkillFromNoteRoundTrip guards the CreateSkill/UpdateSkill write path against
// the ListSkills/GetSkillBody read path drifting apart: whatever renderNote writes, skillFromNote
// must parse back out unchanged (aside from slug/hot-plug, which come from the path, not the
// frontmatter).
func TestRenderNoteAndSkillFromNoteRoundTrip(t *testing.T) {
	name := "My Skill"
	description := "triggers on X and Y"
	storageMode := domain.SkillStorageStructured
	body := "# Steps\n1. do this\n2. do that"

	content := renderNote(name, description, storageMode, body)

	got := skillFromNote(".skills/active/my-skill.md", content)

	if got.Name != name {
		t.Errorf("Name = %q, want %q", got.Name, name)
	}

	if got.Description != description {
		t.Errorf("Description = %q, want %q", got.Description, description)
	}

	if got.StorageMode != storageMode {
		t.Errorf("StorageMode = %q, want %q", got.StorageMode, storageMode)
	}

	if got.Body != body {
		t.Errorf("Body = %q, want %q", got.Body, body)
	}

	if !got.IsHotPlug {
		t.Errorf("IsHotPlug = false, want true (path was under active/)")
	}
}
