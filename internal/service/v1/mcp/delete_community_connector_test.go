package mcp_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/mcp"
	"github.com/ruf-dev/artel/internal/service/v1/subscription"
	"go.redsock.ru/rerrors"
)

func TestDeleteCommunityConnector_NonAdminRejected(t *testing.T) {
	ownerUuid := uuid.New()

	owned := domain.McpDefinition{Name: "acme", OwnerUserUuid: &ownerUuid, IsCommunity: true}
	defsRepo := newFakeMcpDefinitionsRepo(owned)

	authSvc := newFakeAuthService() // not registered as admin

	svc := mcp.New(nil, nil, nil, nil, nil, nil, defsRepo, nil, subscription.NewFree(), authSvc)

	uc := user_context.UserContext{UserUuid: ownerUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	err := svc.DeleteCommunityConnector(ctx, "acme")
	if err == nil {
		t.Fatal("expected error for non-admin caller")
	}

	if _, stillThere := defsRepo.defs["acme"]; !stillThere {
		t.Fatal("expected connector to remain — repo should not be touched before the admin check")
	}
}

func TestDeleteCommunityConnector_NonOwnerRejectedWithNotFound(t *testing.T) {
	ownerUuid := uuid.New()
	otherAdminUuid := uuid.New()

	owned := domain.McpDefinition{Name: "acme", OwnerUserUuid: &ownerUuid, IsCommunity: true}
	defsRepo := newFakeMcpDefinitionsRepo(owned)

	authSvc := newFakeAuthService()
	authSvc.admins[otherAdminUuid] = true

	svc := mcp.New(nil, nil, nil, nil, nil, nil, defsRepo, nil, subscription.NewFree(), authSvc)

	uc := user_context.UserContext{UserUuid: otherAdminUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	err := svc.DeleteCommunityConnector(ctx, "acme")
	if err == nil {
		t.Fatal("expected error for non-owner caller")
	}

	if !rerrors.Is(err, user_errors.NotFound) {
		t.Fatalf("expected user_errors.NotFound (not a distinct forbidden error), got %v", err)
	}

	if _, stillThere := defsRepo.defs["acme"]; !stillThere {
		t.Fatal("expected connector to remain")
	}
}

func TestDeleteCommunityConnector_OwnerDeletesSuccessfully(t *testing.T) {
	ownerUuid := uuid.New()

	owned := domain.McpDefinition{Name: "acme", OwnerUserUuid: &ownerUuid, IsCommunity: true}
	defsRepo := newFakeMcpDefinitionsRepo(owned)

	authSvc := newFakeAuthService()
	authSvc.admins[ownerUuid] = true

	svc := mcp.New(nil, nil, nil, nil, nil, nil, defsRepo, nil, subscription.NewFree(), authSvc)

	uc := user_context.UserContext{UserUuid: ownerUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	err := svc.DeleteCommunityConnector(ctx, "acme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, stillThere := defsRepo.defs["acme"]; stillThere {
		t.Fatal("expected connector to be deleted")
	}
}
