package mcp_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/mcp"
	"github.com/ruf-dev/artel/internal/service/v1/subscription"
)

// fakeMcpDefinitionsRepo is an in-memory repository.McpDefinitionsRepo test double.
type fakeMcpDefinitionsRepo struct {
	defs        map[string]domain.McpDefinition
	upsertCalls []domain.McpDefinition
}

func newFakeMcpDefinitionsRepo(defs ...domain.McpDefinition) *fakeMcpDefinitionsRepo {
	byName := make(map[string]domain.McpDefinition, len(defs))
	for _, d := range defs {
		byName[d.Name] = d
	}

	return &fakeMcpDefinitionsRepo{defs: byName}
}

func (f *fakeMcpDefinitionsRepo) Upsert(_ context.Context, def domain.McpDefinition) (domain.McpDefinition, error) {
	f.upsertCalls = append(f.upsertCalls, def)
	f.defs[def.Name] = def

	return def, nil
}

func (f *fakeMcpDefinitionsRepo) Get(_ context.Context, name string) (sql.Null[domain.McpDefinition], error) {
	def, ok := f.defs[name]
	if !ok {
		return sql.Null[domain.McpDefinition]{}, nil
	}

	result := sql.Null[domain.McpDefinition]{V: def, Valid: true}

	return result, nil
}

func (f *fakeMcpDefinitionsRepo) List(_ context.Context) ([]domain.McpDefinition, error) {
	list := make([]domain.McpDefinition, 0, len(f.defs))
	for _, d := range f.defs {
		list = append(list, d)
	}

	return list, nil
}

func (f *fakeMcpDefinitionsRepo) Delete(_ context.Context, name string) error {
	delete(f.defs, name)

	return nil
}

func (f *fakeMcpDefinitionsRepo) GetTool(
	_ context.Context, _ string, _ string,
) (sql.Null[domain.McpToolDef], error) {
	return sql.Null[domain.McpToolDef]{}, nil
}

func (f *fakeMcpDefinitionsRepo) ListAllTools(_ context.Context) ([]domain.McpToolRef, error) {
	return nil, nil
}

// fakeAuthService is a minimal service.AuthService test double — only CheckIsAdmin and GetMe
// are exercised by create_community_connector/delete_community_connector.
type fakeAuthService struct {
	admins map[uuid.UUID]bool
	emails map[uuid.UUID]string
}

func newFakeAuthService() *fakeAuthService {
	return &fakeAuthService{admins: map[uuid.UUID]bool{}, emails: map[uuid.UUID]string{}}
}

func (f *fakeAuthService) Register(_ context.Context, _, _ string) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeAuthService) Login(_ context.Context, _, _ string) (domain.Session, error) {
	return domain.Session{}, nil
}

func (f *fakeAuthService) Logout(_ context.Context, _ string) error {
	return nil
}

func (f *fakeAuthService) ValidateToken(_ context.Context, _ string) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeAuthService) LoginViaTelegram(_ context.Context, _ string) (domain.Session, error) {
	return domain.Session{}, nil
}

func (f *fakeAuthService) GetMe(_ context.Context, userUuid uuid.UUID) (domain.UserDetails, error) {
	details := domain.UserDetails{}
	details.Uuid = userUuid
	details.Email = f.emails[userUuid]

	return details, nil
}

func (f *fakeAuthService) CheckIsAdmin(_ context.Context, userUuid uuid.UUID) error {
	if !f.admins[userUuid] {
		return user_errors.NotAdmin
	}

	return nil
}

func (f *fakeAuthService) Refresh(_ context.Context, _ string) (domain.Session, error) {
	return domain.Session{}, nil
}

func (f *fakeAuthService) EnsureNoAuthUser(_ context.Context) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeAuthService) RegisterAdmin(_ context.Context, _, _ string) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeAuthService) CreateUserUnchecked(_ context.Context, _, _ string) (domain.User, error) {
	return domain.User{}, nil
}

func (f *fakeAuthService) LoginOrRegisterAdminViaTelegram(_ context.Context, _ string) (domain.Session, error) {
	return domain.Session{}, nil
}

func (f *fakeAuthService) ChangePassword(_ context.Context, _ uuid.UUID, _ string) error {
	return nil
}

func communityConnectorParams(name string) map[string]interface{} {
	return map[string]interface{}{
		"name":        name,
		"description": "A test connector",
		"tools": []interface{}{
			map[string]interface{}{
				"api_description": map[string]interface{}{
					"name":        "do_thing",
					"description": "Does a thing",
					"properties":  map[string]interface{}{},
					"required":    []interface{}{},
				},
				"action": map[string]interface{}{
					"http": map[string]interface{}{
						"method":      "GET",
						"url":         "https://example.com/thing",
						"credentials": name,
					},
				},
			},
		},
	}
}

func TestCreateCommunityConnector_NonAdminRejected(t *testing.T) {
	userUuid := uuid.New()

	authSvc := newFakeAuthService() // no admins registered

	defsRepo := newFakeMcpDefinitionsRepo()

	svc := mcp.New(nil, nil, nil, nil, nil, nil, defsRepo, nil, subscription.NewFree(), authSvc)

	keyCtx := domain.McpKeyContext{UserUuid: userUuid}
	params := communityConnectorParams("acme")

	_, err := svc.ExecuteTool(context.Background(), keyCtx, "create_community_connector", params)
	if err == nil {
		t.Fatal("expected error for non-admin caller, got nil")
	}

	if len(defsRepo.upsertCalls) != 0 {
		t.Fatalf("expected repo never touched, got %d upsert calls", len(defsRepo.upsertCalls))
	}
}

func TestCreateCommunityConnector_AdminCreatesNew(t *testing.T) {
	adminUuid := uuid.New()

	authSvc := newFakeAuthService()
	authSvc.admins[adminUuid] = true
	authSvc.emails[adminUuid] = "admin@example.com"

	defsRepo := newFakeMcpDefinitionsRepo()

	svc := mcp.New(nil, nil, nil, nil, nil, nil, defsRepo, nil, subscription.NewFree(), authSvc)

	keyCtx := domain.McpKeyContext{UserUuid: adminUuid}
	params := communityConnectorParams("acme")

	result, err := svc.ExecuteTool(context.Background(), keyCtx, "create_community_connector", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Text == "" {
		t.Fatal("expected non-empty result text")
	}

	if len(defsRepo.upsertCalls) != 1 {
		t.Fatalf("expected 1 upsert call, got %d", len(defsRepo.upsertCalls))
	}

	saved := defsRepo.upsertCalls[0]
	if !saved.IsCommunity {
		t.Fatal("expected IsCommunity to be true")
	}

	if saved.OwnerUserUuid == nil || *saved.OwnerUserUuid != adminUuid {
		t.Fatalf("expected OwnerUserUuid %s, got %v", adminUuid, saved.OwnerUserUuid)
	}

	if saved.Author != "admin@example.com" {
		t.Fatalf("expected author admin@example.com, got %s", saved.Author)
	}
}

func TestCreateCommunityConnector_AdminRecreatesOwnConnector(t *testing.T) {
	adminUuid := uuid.New()

	authSvc := newFakeAuthService()
	authSvc.admins[adminUuid] = true
	authSvc.emails[adminUuid] = "admin@example.com"

	existing := domain.McpDefinition{
		Name:          "acme",
		Author:        "admin@example.com",
		Description:   "old description",
		OwnerUserUuid: &adminUuid,
		IsCommunity:   true,
	}
	defsRepo := newFakeMcpDefinitionsRepo(existing)

	svc := mcp.New(nil, nil, nil, nil, nil, nil, defsRepo, nil, subscription.NewFree(), authSvc)

	keyCtx := domain.McpKeyContext{UserUuid: adminUuid}
	params := communityConnectorParams("acme")

	_, err := svc.ExecuteTool(context.Background(), keyCtx, "create_community_connector", params)
	if err != nil {
		t.Fatalf("unexpected error re-creating own connector: %v", err)
	}

	if len(defsRepo.upsertCalls) != 1 {
		t.Fatalf("expected 1 upsert call, got %d", len(defsRepo.upsertCalls))
	}
}

func TestCreateCommunityConnector_RejectsSystemMomName(t *testing.T) {
	adminUuid := uuid.New()

	authSvc := newFakeAuthService()
	authSvc.admins[adminUuid] = true

	systemMom := domain.McpDefinition{
		Name:          "email",
		Author:        "Artel",
		Description:   "system mom",
		OwnerUserUuid: nil, // system MoM
		IsCommunity:   false,
	}
	defsRepo := newFakeMcpDefinitionsRepo(systemMom)

	svc := mcp.New(nil, nil, nil, nil, nil, nil, defsRepo, nil, subscription.NewFree(), authSvc)

	keyCtx := domain.McpKeyContext{UserUuid: adminUuid}
	params := communityConnectorParams("email")

	_, err := svc.ExecuteTool(context.Background(), keyCtx, "create_community_connector", params)
	if err == nil {
		t.Fatal("expected error overwriting a system mom name, got nil")
	}

	if len(defsRepo.upsertCalls) != 0 {
		t.Fatalf("expected repo never touched, got %d upsert calls", len(defsRepo.upsertCalls))
	}
}

func TestCreateCommunityConnector_RejectsOtherAdminsConnector(t *testing.T) {
	callingAdmin := uuid.New()
	otherAdmin := uuid.New()

	authSvc := newFakeAuthService()
	authSvc.admins[callingAdmin] = true

	othersConnector := domain.McpDefinition{
		Name:          "acme",
		Author:        "other-admin@example.com",
		Description:   "someone else's",
		OwnerUserUuid: &otherAdmin,
		IsCommunity:   true,
	}
	defsRepo := newFakeMcpDefinitionsRepo(othersConnector)

	svc := mcp.New(nil, nil, nil, nil, nil, nil, defsRepo, nil, subscription.NewFree(), authSvc)

	keyCtx := domain.McpKeyContext{UserUuid: callingAdmin}
	params := communityConnectorParams("acme")

	_, err := svc.ExecuteTool(context.Background(), keyCtx, "create_community_connector", params)
	if err == nil {
		t.Fatal("expected error overwriting another admin's connector, got nil")
	}

	if len(defsRepo.upsertCalls) != 0 {
		t.Fatalf("expected repo never touched, got %d upsert calls", len(defsRepo.upsertCalls))
	}
}
