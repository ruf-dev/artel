package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (s *McpServiceImpl) ResolveKey(ctx context.Context, rawToken string) (domain.McpKeyContext, error) {
	if !strings.HasPrefix(rawToken, tokenPrefix) {
		return domain.McpKeyContext{}, rerrors.Wrap(user_errors.McpInvalidToken, "prefix missing", tokenPrefix)
	}

	rest := rawToken[len(tokenPrefix):]
	parts := strings.SplitN(rest, "_", 2)
	if len(parts) != 2 {
		return domain.McpKeyContext{}, rerrors.Wrap(user_errors.McpInvalidToken, "invalid token format")
	}
	uuidHex := parts[0]
	secretHex := parts[1]
	if len(uuidHex) < 32 {
		return domain.McpKeyContext{}, rerrors.Wrap(user_errors.McpInvalidToken, "invalid token format")
	}

	uuidFormatted := fmt.Sprintf("%s-%s-%s-%s-%s", uuidHex[0:8], uuidHex[8:12], uuidHex[12:16], uuidHex[16:20], uuidHex[20:32])
	keyUUID, err := uuid.Parse(uuidFormatted)
	if err != nil {
		return domain.McpKeyContext{}, rerrors.Wrap(err, "parse key uuid")
	}

	mcpKey, err := s.mcpKeys.GetMcpKeyByID(ctx, keyUUID)
	if err != nil {
		return domain.McpKeyContext{}, rerrors.Wrap(err, "get mcp key")
	}

	if mcpKey.RevokedAt != nil {
		return domain.McpKeyContext{}, rerrors.Wrap(user_errors.McpKeyRevoked)
	}

	err = bcrypt.CompareHashAndPassword(mcpKey.KeyHash, []byte(secretHex))
	if err != nil {
		return domain.McpKeyContext{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	vault, err := s.vaults.GetByID(ctx, mcpKey.VaultUuid)
	if err != nil {
		return domain.McpKeyContext{}, rerrors.Wrap(err, "get vault")
	}

	couchInstance, err := s.couchInstances.Get(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return domain.McpKeyContext{}, rerrors.Wrap(err, "get couch instance")
	}

	result := domain.McpKeyContext{
		KeyUuid:   mcpKey.Uuid,
		VaultUuid: mcpKey.VaultUuid,
		UserUuid:  mcpKey.UserUuid,
		CouchURL:  couchInstance.Url,
		CouchDb:   vault.CouchDBName,
		CouchUser: couchInstance.Username,
		CouchPass: couchInstance.Password,
	}

	go func() {
		if touchErr := s.mcpKeys.TouchLastAccessed(context.Background(), mcpKey.Uuid); touchErr != nil {
			log.Error().Err(touchErr).Str("key_uuid", mcpKey.Uuid.String()).Msg("ResolveKey: touch last accessed failed")
		}
	}()

	return result, nil
}
