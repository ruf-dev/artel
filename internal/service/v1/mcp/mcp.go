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
	mcpKeys         repository.McpKeyRepository
	vaults          repository.Vaults
	vaultMembers    repository.VaultMembers
	emailAccounts   repository.EmailAccounts
	couchInstances  repository.CouchInstances
	userPermissions repository.UserPermissionsRepo
}

func New(
	mcpKeys repository.McpKeyRepository,
	vaults repository.Vaults,
	vaultMembers repository.VaultMembers,
	emailAccounts repository.EmailAccounts,
	couchInstances repository.CouchInstances,
	userPermissions repository.UserPermissionsRepo,
) *McpServiceImpl {
	return &McpServiceImpl{
		mcpKeys:         mcpKeys,
		vaults:          vaults,
		vaultMembers:    vaultMembers,
		emailAccounts:   emailAccounts,
		couchInstances:  couchInstances,
		userPermissions: userPermissions,
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

	key, err = s.mcpKeys.CreateMcpKey(ctx, vaultID, uc.UserUuid, keyID, name, keyHash, keyPreview)
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

func (s *McpServiceImpl) ListUserKeys(ctx context.Context) ([]domain.McpKey, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	keys, err := s.mcpKeys.ListMcpKeysByUser(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "list user mcp keys")
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

func (s *McpServiceImpl) SetKeyAccess(ctx context.Context, keyID uuid.UUID, vaultID uuid.UUID, emailAccountID *uuid.UUID) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return user_errors.Unauthenticated
	}

	_, err := s.vaultMembers.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return rerrors.Wrap(user_errors.Unauthenticated, "user is not a member of target vault")
	}

	if emailAccountID != nil {
		var account domain.EmailAccount
		account, err = s.emailAccounts.GetByUuid(ctx, *emailAccountID)
		if err != nil {
			return rerrors.Wrap(err, "get email account")
		}
		if account.UserUuid != uc.UserUuid {
			return user_errors.Unauthenticated
		}
	}

	err = s.mcpKeys.SetMcpKeyAccess(ctx, keyID, uc.UserUuid, vaultID, emailAccountID)
	if err != nil {
		return rerrors.Wrap(err, "set mcp key access")
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

	perms, permErr := s.userPermissions.Get(ctx, mcpKey.UserUuid)
	if permErr != nil {
		log.Error().Err(permErr).Msg("ResolveKey: get user permissions failed")
	} else {
		result.HasEmails = perms.HasEmails
	}

	go func() {
		if touchErr := s.mcpKeys.TouchLastAccessed(context.Background(), mcpKey.Uuid); touchErr != nil {
			log.Error().Err(touchErr).Str("key_uuid", mcpKey.Uuid.String()).Msg("ResolveKey: touch last accessed failed")
		}
	}()

	return result, nil
}
