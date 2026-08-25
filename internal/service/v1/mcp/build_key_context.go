package mcp

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

// BuildKeyContext builds a domain.McpKeyContext for an already-authenticated in-process caller
// (e.g. Simple Chat's agent loop) acting as userUuid within vaultUuid — unlike ResolveKey, there
// is no mcp_keys token involved, so KeyUuid is left as uuid.Nil on the returned context (nothing
// downstream depends on it besides ResolveKey's own TouchLastAccessed bookkeeping, which this
// path intentionally skips). Requires userUuid to actually be a member of vaultUuid — returns
// an error otherwise (mirrors internal/service/v1/workbench/workbench.go's CreateWorkbench
// membership check).
func (s *ServiceImpl) BuildKeyContext(ctx context.Context, vaultUuid, userUuid uuid.UUID) (domain.McpKeyContext, error) {
	_, err := s.vaultMembers.Get(ctx, vaultUuid, userUuid)
	if err != nil {
		return domain.McpKeyContext{}, rerrors.Wrap(err, "error checking vault membership")
	}

	result, err := s.buildKeyContext(ctx, uuid.Nil, vaultUuid, userUuid)
	if err != nil {
		return domain.McpKeyContext{}, rerrors.Wrap(err, "error building key context")
	}

	return result, nil
}
