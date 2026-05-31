package mcpkeys

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
)

type McpKeyRepo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *McpKeyRepo {
	return &McpKeyRepo{q: q}
}

func (r *McpKeyRepo) CreateMcpKey(ctx context.Context, vaultID, userID, keyID uuid.UUID, name string, keyHash []byte, keyPreview string) (domain.McpKey, error) {
	params := artel_q.CreateMcpKeyParams{
		ID:         keyID,
		VaultID:    vaultID,
		UserID:     userID,
		Name:       name,
		KeyHash:    keyHash,
		KeyPreview: keyPreview,
	}
	row, err := r.q.CreateMcpKey(ctx, params)
	if err != nil {
		return domain.McpKey{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "store mcp key")
	}

	return toMcpKey(row), nil
}

func (r *McpKeyRepo) ListMcpKeysByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error) {
	rows, err := r.q.ListMcpKeysByVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "list mcp keys by vault")
	}

	keys := make([]domain.McpKey, 0, len(rows))
	for _, row := range rows {
		keys = append(keys, toMcpKey(row))
	}
	return keys, nil
}

func (r *McpKeyRepo) ListMcpKeysByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.McpKey, error) {
	rows, err := r.q.ListMcpKeysByUser(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "list mcp keys by user")
	}

	keys := make([]domain.McpKey, 0, len(rows))
	for _, row := range rows {
		keys = append(keys, toMcpKey(row))
	}
	return keys, nil
}

func (r *McpKeyRepo) GetMcpKeyByID(ctx context.Context, id uuid.UUID) (domain.McpKey, error) {
	row, err := r.q.GetMcpKeyByID(ctx, id)
	if err != nil {
		return domain.McpKey{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "get mcp key by id")
	}

	return toMcpKey(row), nil
}

func (r *McpKeyRepo) ListActiveMcpKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error) {
	rows, err := r.q.ListActiveMcpKeys(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "list active mcp keys")
	}

	keys := make([]domain.McpKey, 0, len(rows))
	for _, row := range rows {
		keys = append(keys, toMcpKey(row))
	}
	return keys, nil
}

func (r *McpKeyRepo) RevokeMcpKey(ctx context.Context, id uuid.UUID) error {
	err := r.q.RevokeMcpKey(ctx, id)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "revoke mcp key")
	}
	return nil
}

func (r *McpKeyRepo) SetMcpKeyAccess(ctx context.Context, keyUuid, userUuid, vaultUuid uuid.UUID, emailAccountUuid *uuid.UUID) error {
	var nullEmailID uuid.NullUUID
	if emailAccountUuid != nil {
		nullEmailID = uuid.NullUUID{UUID: *emailAccountUuid, Valid: true}
	}

	err := r.q.SetMcpKeyAccess(ctx, artel_q.SetMcpKeyAccessParams{
		ID:             keyUuid,
		UserID:         userUuid,
		VaultID:        vaultUuid,
		EmailAccountID: nullEmailID,
	})
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "set mcp key access")
	}
	return nil
}

func (r *McpKeyRepo) TouchLastAccessed(ctx context.Context, keyUuid uuid.UUID) error {
	err := r.q.TouchMcpKeyLastAccessed(ctx, keyUuid)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "touch mcp key last accessed")
	}
	return nil
}

func toMcpKey(row artel_q.McpKey) domain.McpKey {
	key := domain.McpKey{
		Uuid:       row.ID,
		VaultUuid:  row.VaultID,
		UserUuid:   row.UserID,
		Name:       row.Name,
		KeyHash:    row.KeyHash,
		KeyPreview: row.KeyPreview,
		CreatedAt:  row.CreatedAt,
	}
	if row.RevokedAt.Valid {
		key.RevokedAt = &row.RevokedAt.Time
	}
	if row.EmailAccountID.Valid {
		id := row.EmailAccountID.UUID
		key.EmailAccountUuid = &id
	}
	if row.LastAccessedAt.Valid {
		key.LastAccessedAt = &row.LastAccessedAt.Time
	}
	return key
}
