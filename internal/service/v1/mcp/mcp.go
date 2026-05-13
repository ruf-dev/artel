package mcp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

const tokenPrefix = "artel_vtk_"

const bcryptCost = 12

type McpServiceImpl struct {
	mcpKeys        repository.McpKeyRepository
	vaults         repository.Vaults
	couchInstances repository.CouchInstances
}

func New(
	mcpKeys repository.McpKeyRepository,
	vaults repository.Vaults,
	couchInstances repository.CouchInstances,
) *McpServiceImpl {
	return &McpServiceImpl{
		mcpKeys:        mcpKeys,
		vaults:         vaults,
		couchInstances: couchInstances,
	}
}

func (s *McpServiceImpl) CreateKey(ctx context.Context, vaultID uuid.UUID, name string) (rawToken string, key domain.McpKey, err error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return "", domain.McpKey{}, user_errors.Unauthenticated
	}

	keyID := uuid.New()

	secretBytes := make([]byte, 16)
	_, err = rand.Read(secretBytes)
	if err != nil {
		return "", domain.McpKey{}, rerrors.Wrap(err, "generate secret bytes")
	}
	secretHex := hex.EncodeToString(secretBytes)

	uuidHex := strings.ReplaceAll(keyID.String(), "-", "")
	rawToken = tokenPrefix + fmt.Sprintf("%s_%s", uuidHex, secretHex)

	keyPreview := rawToken[:12]

	keyHash, err := bcrypt.GenerateFromPassword([]byte(secretHex), bcryptCost)
	if err != nil {
		return "", domain.McpKey{}, rerrors.Wrap(err, "hash secret")
	}

	key, err = s.mcpKeys.CreateMcpKey(ctx, vaultID, uc.UserUuid, name, keyHash, keyPreview)
	if err != nil {
		return "", domain.McpKey{}, rerrors.Wrap(err, "create mcp key")
	}

	return rawToken, key, nil
}

func (s *McpServiceImpl) ListKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error) {
	keys, err := s.mcpKeys.ListMcpKeysByVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mcp keys")
	}
	return keys, nil
}

func (s *McpServiceImpl) RevokeKey(ctx context.Context, keyID uuid.UUID) error {
	err := s.mcpKeys.RevokeMcpKey(ctx, keyID)
	if err != nil {
		return rerrors.Wrap(err, "revoke mcp key")
	}
	return nil
}

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
		log.Error().Err(err).Str("vault_uuid", mcpKey.VaultUuid.String()).Msg("ResolveKey: get vault failed")
		return domain.McpKeyContext{}, rerrors.Wrap(err, "get vault")
	}

	couchWithAccount, err := s.couchInstances.Pick(ctx, vault.CouchInstanceUuid)
	if err != nil {
		log.Error().Err(err).Str("couch_instance_uuid", vault.CouchInstanceUuid.String()).Msg("ResolveKey: pick couch instance failed")
		return domain.McpKeyContext{}, rerrors.Wrap(err, "pick couch instance")
	}

	result := domain.McpKeyContext{
		KeyUuid:   mcpKey.Uuid,
		VaultUuid: mcpKey.VaultUuid,
		UserUuid:  mcpKey.UserUuid,
		CouchURL:  couchWithAccount.Instance.Url,
		CouchDb:   vault.CouchDBName,
		CouchUser: couchWithAccount.Instance.Username,
		CouchPass: couchWithAccount.Instance.Password,
	}
	return result, nil
}
