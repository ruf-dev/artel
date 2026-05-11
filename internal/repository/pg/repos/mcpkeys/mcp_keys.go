package mcpkeys

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/domain"
)

type McpKeyRepo struct {
}

func New() *McpKeyRepo {
	return &McpKeyRepo{}
}

func (r *McpKeyRepo) CreateMcpKey(ctx context.Context, vaultID, userID uuid.UUID, name string, keyHash []byte, keyPreview string) (domain.McpKey, error) {
	return domain.McpKey{}, nil
}

func (r *McpKeyRepo) ListMcpKeysByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error) {
	return nil, nil
}

func (r *McpKeyRepo) GetMcpKeyByID(ctx context.Context, id uuid.UUID) (domain.McpKey, error) {
	return domain.McpKey{}, nil
}

func (r *McpKeyRepo) ListActiveMcpKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error) {
	return nil, nil
}

func (r *McpKeyRepo) RevokeMcpKey(ctx context.Context, id uuid.UUID) error {
	return nil
}
