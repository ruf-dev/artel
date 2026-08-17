package mcp

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// fakeSkillsService is a hand-rolled test double for service.SkillsService — this package has no
// existing mock-generation setup, so it implements the interface directly like the sqlmock-based
// executors package tests do for *sql.DB.
type fakeSkillsService struct {
	createSkillFn   func(ctx context.Context, vaultUuid uuid.UUID, name, description string, storageMode domain.SkillStorageMode, body string, hotPlug bool) (domain.Skill, error)
	setHotPlugFn    func(ctx context.Context, vaultUuid uuid.UUID, slug string, hotPlug bool) (domain.Skill, error)
	deleteSkillFn   func(ctx context.Context, vaultUuid uuid.UUID, slug string) error
	deleteCalls     []string
	setHotPlugCalls []bool
}

func (f *fakeSkillsService) ListSkills(context.Context, uuid.UUID) ([]domain.Skill, error) {
	return nil, nil
}

func (f *fakeSkillsService) GetSkillBody(context.Context, uuid.UUID, string) (domain.Skill, error) {
	return domain.Skill{}, nil
}

func (f *fakeSkillsService) CreateSkill(
	ctx context.Context, vaultUuid uuid.UUID, name, description string,
	storageMode domain.SkillStorageMode, body string, hotPlug bool,
) (domain.Skill, error) {
	return f.createSkillFn(ctx, vaultUuid, name, description, storageMode, body, hotPlug)
}

func (f *fakeSkillsService) UpdateSkill(
	context.Context, uuid.UUID, string, string, string, domain.SkillStorageMode, string,
) (domain.Skill, error) {
	return domain.Skill{}, nil
}

func (f *fakeSkillsService) SetSkillHotPlug(
	ctx context.Context, vaultUuid uuid.UUID, slug string, hotPlug bool,
) (domain.Skill, error) {
	f.setHotPlugCalls = append(f.setHotPlugCalls, hotPlug)

	return f.setHotPlugFn(ctx, vaultUuid, slug, hotPlug)
}

func (f *fakeSkillsService) DeleteSkill(ctx context.Context, vaultUuid uuid.UUID, slug string) error {
	f.deleteCalls = append(f.deleteCalls, slug)

	return f.deleteSkillFn(ctx, vaultUuid, slug)
}

func TestSkillToolCollides(t *testing.T) {
	tests := []struct {
		name      string
		isBuiltin func(string) bool
		slug      string
		want      bool
	}{
		{
			name:      "no collision",
			isBuiltin: func(string) bool { return false },
			slug:      "weekly-report",
			want:      false,
		},
		{
			name: "collides when skill_<slug> matches an existing builtin tool name",
			isBuiltin: func(name string) bool {
				return name == "skill_list_files"
			},
			slug: "list_files",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := skillToolCollides(tt.isBuiltin, tt.slug)
			if got != tt.want {
				t.Errorf("skillToolCollides(...) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateSkillTool_MissingParams(t *testing.T) {
	svc := &ServiceImpl{skillsSvc: &fakeSkillsService{}}
	keyCtx := domain.McpKeyContext{VaultUuid: uuid.New(), UserUuid: uuid.New()}

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr error
	}{
		{
			name:    "missing name",
			params:  map[string]interface{}{},
			wantErr: user_errors.SkillNameRequired,
		},
		{
			name:    "missing description",
			params:  map[string]interface{}{fieldName: "My Skill"},
			wantErr: user_errors.SkillDescriptionRequired,
		},
		{
			name: "missing storage_mode",
			params: map[string]interface{}{
				fieldName:        "My Skill",
				fieldDescription: "does things",
			},
			wantErr: user_errors.SkillStorageModeRequired,
		},
		{
			name: "invalid storage_mode",
			params: map[string]interface{}{
				fieldName:        "My Skill",
				fieldDescription: "does things",
				fieldStorageMode: "not-a-real-mode",
			},
			wantErr: user_errors.SkillStorageModeInvalid,
		},
		{
			name: "missing body",
			params: map[string]interface{}{
				fieldName:        "My Skill",
				fieldDescription: "does things",
				fieldStorageMode: string(domain.SkillStorageNone),
			},
			wantErr: user_errors.SkillBodyRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.createSkillTool(context.Background(), keyCtx, tt.params)
			if !rerrors.Is(err, tt.wantErr) {
				t.Fatalf("createSkillTool() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateSkillTool_Success_NoCollision(t *testing.T) {
	created := domain.Skill{Slug: "my-skill", Name: "My Skill", IsHotPlug: true}
	fake := &fakeSkillsService{
		createSkillFn: func(context.Context, uuid.UUID, string, string, domain.SkillStorageMode, string, bool) (domain.Skill, error) {
			return created, nil
		},
	}
	svc := &ServiceImpl{skillsSvc: fake}
	keyCtx := domain.McpKeyContext{VaultUuid: uuid.New(), UserUuid: uuid.New()}
	params := map[string]interface{}{
		fieldName:        "My Skill",
		fieldDescription: "does things",
		fieldStorageMode: string(domain.SkillStorageNone),
		fieldBody:        "do the thing",
		fieldHotPlug:     true,
	}

	result, err := svc.createSkillTool(context.Background(), keyCtx, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fake.deleteCalls) != 0 {
		t.Fatalf("expected no rollback delete call, got %v", fake.deleteCalls)
	}

	if result.Text == "" {
		t.Fatalf("expected a non-empty confirmation text")
	}
}

// TestCreateSkillTool_NoRealCollisionToday documents that under the actual builtin tool name
// set (none of which starts with "skill_"), createSkillTool's rollback branch is unreachable —
// skillToolCollides is exactly what would gate it (see TestSkillToolCollides, which exercises
// that gate directly against a fake predicate). This pins the current behavior: hot-plugging a
// skill slugified from an existing builtin tool's own name never rolls back, because the
// resulting skill_<name> tool name is still distinct from every builtin tool name.
func TestCreateSkillTool_NoRealCollisionToday(t *testing.T) {
	created := domain.Skill{Slug: "list_files", Name: "List Files", IsHotPlug: true}
	fake := &fakeSkillsService{
		createSkillFn: func(context.Context, uuid.UUID, string, string, domain.SkillStorageMode, string, bool) (domain.Skill, error) {
			return created, nil
		},
	}
	svc := &ServiceImpl{skillsSvc: fake}
	keyCtx := domain.McpKeyContext{VaultUuid: uuid.New(), UserUuid: uuid.New()}
	params := map[string]interface{}{
		fieldName:        "List Files",
		fieldDescription: "does things",
		fieldStorageMode: string(domain.SkillStorageNone),
		fieldBody:        "do the thing",
		fieldHotPlug:     true,
	}

	_, err := svc.createSkillTool(context.Background(), keyCtx, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fake.deleteCalls) != 0 {
		t.Fatalf("expected no rollback delete call under the real (non-colliding) builtin tool set, got %v", fake.deleteCalls)
	}
}

func TestSetSkillHotPlugTool_MissingParams(t *testing.T) {
	svc := &ServiceImpl{skillsSvc: &fakeSkillsService{}}
	keyCtx := domain.McpKeyContext{VaultUuid: uuid.New(), UserUuid: uuid.New()}

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr error
	}{
		{
			name:    "missing slug",
			params:  map[string]interface{}{},
			wantErr: user_errors.SkillSlugRequired,
		},
		{
			name:    "missing hot_plug",
			params:  map[string]interface{}{fieldSlug: "my-skill"},
			wantErr: user_errors.SkillHotPlugRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.setSkillHotPlugTool(context.Background(), keyCtx, tt.params)
			if !rerrors.Is(err, tt.wantErr) {
				t.Fatalf("setSkillHotPlugTool() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetSkillHotPlugTool_Success(t *testing.T) {
	fake := &fakeSkillsService{
		setHotPlugFn: func(context.Context, uuid.UUID, string, bool) (domain.Skill, error) {
			return domain.Skill{Slug: "my-skill", IsHotPlug: true}, nil
		},
	}
	svc := &ServiceImpl{skillsSvc: fake}
	keyCtx := domain.McpKeyContext{VaultUuid: uuid.New(), UserUuid: uuid.New()}
	params := map[string]interface{}{fieldSlug: "my-skill", fieldHotPlug: true}

	_, err := svc.setSkillHotPlugTool(context.Background(), keyCtx, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fake.deleteCalls) != 0 {
		t.Fatalf("expected no rollback for set_skill_hot_plug, got delete calls %v", fake.deleteCalls)
	}
}
