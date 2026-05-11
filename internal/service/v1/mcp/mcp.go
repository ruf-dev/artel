package mcp

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
)

type McpServiceImpl struct {
	repo repository.McpKeyRepository
}

func New(repo repository.McpKeyRepository) *McpServiceImpl {
	return &McpServiceImpl{repo: repo}
}

func (s *McpServiceImpl) CreateKey(ctx context.Context, vaultID uuid.UUID, name string) (rawToken string, key domain.McpKey, err error) {
	return "", domain.McpKey{}, nil
}

func (s *McpServiceImpl) ListKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error) {
	return nil, nil
}

func (s *McpServiceImpl) RevokeKey(ctx context.Context, keyID uuid.UUID) error {
	return nil
}

func (s *McpServiceImpl) ResolveKey(ctx context.Context, rawToken string) (domain.McpKeyContext, error) {
	return domain.McpKeyContext{}, nil
}
