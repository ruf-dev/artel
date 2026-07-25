package mcp_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/v1/mcp"
	"github.com/ruf-dev/artel/internal/service/v1/subscription"
)

func TestListCommunityConnectors_ExcludesSystemMoms(t *testing.T) {
	adminUuid := uuid.New()

	systemMom := domain.McpDefinition{
		Name:          "email",
		Author:        "Artel",
		Description:   "system mom",
		OwnerUserUuid: nil,
		IsCommunity:   false,
	}
	community := domain.McpDefinition{
		Name:          "acme",
		Author:        "admin@example.com",
		Description:   "community mom",
		OwnerUserUuid: &adminUuid,
		IsCommunity:   true,
	}

	defsRepo := newFakeMcpDefinitionsRepo(systemMom, community)
	connsRepo := &fakeExternalConnectionRepo{}

	authSvc := newFakeAuthService()

	svc := mcp.New(nil, nil, nil, nil, nil, nil, defsRepo, connsRepo, subscription.NewFree(), authSvc)

	uc := user_context.UserContext{UserUuid: adminUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	results, err := svc.ListCommunityConnectors(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 community connector (system mom excluded), got %d", len(results))
	}

	if results[0].Name != "acme" {
		t.Fatalf("expected acme, got %s", results[0].Name)
	}

	if !results[0].ViewerIsOwner {
		t.Fatal("expected caller to be marked as owner of their own connector")
	}
}
