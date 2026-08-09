package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenPartsCount = 2
	uuidHexLen      = 32
	secretByteLen   = 16
)

func (s *ServiceImpl) ResolveKey(ctx context.Context, rawToken string) (domain.McpKeyContext, error) {
	if !strings.HasPrefix(rawToken, tokenPrefix) {
		return domain.McpKeyContext{}, rerrors.Wrap(user_errors.McpInvalidToken, "prefix missing", tokenPrefix)
	}

	rest := rawToken[len(tokenPrefix):]

	parts := strings.SplitN(rest, "_", tokenPartsCount)
	if len(parts) != tokenPartsCount {
		return domain.McpKeyContext{}, rerrors.Wrap(user_errors.McpInvalidToken, "invalid token format")
	}

	uuidHex := parts[0]
	secretHex := parts[1]

	if len(uuidHex) < uuidHexLen {
		return domain.McpKeyContext{}, rerrors.Wrap(user_errors.McpInvalidToken, "invalid token format")
	}

	uuidFormatted := fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		uuidHex[0:8],
		uuidHex[8:12],
		uuidHex[12:16],
		uuidHex[16:20],
		uuidHex[20:32],
	)

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

	var s3Ctx *domain.McpKeyS3Context

	if vault.S3InstanceUuid != nil {
		s3Instance, s3Err := s.s3Instances.Get(ctx, *vault.S3InstanceUuid)
		if s3Err != nil {
			return domain.McpKeyContext{}, rerrors.Wrap(s3Err, "get s3 instance")
		}

		s3Ctx = &domain.McpKeyS3Context{
			Endpoint:  s3Instance.Endpoint,
			Region:    s3Instance.Region,
			AccessKey: s3Instance.AccessKey,
			SecretKey: s3Instance.SecretKey,
			UseSSL:    s3Instance.UseSSL,
			PathStyle: s3Instance.PathStyle,
			Bucket:    vault.S3BucketName,
		}
	}

	var pgCtx *domain.McpKeyPostgresContext

	vpgDB, err := s.vaultPostgresDatabases.GetByVaultID(ctx, mcpKey.VaultUuid)
	if err != nil {
		return domain.McpKeyContext{}, rerrors.Wrap(err, "get vault postgres database")
	}

	if vpgDB.Valid && vpgDB.V.Status == domain.VaultPostgresStatusReady {
		pgInstance, pgErr := s.postgresInstances.Get(ctx, vpgDB.V.PostgresInstanceUuid)
		if pgErr != nil {
			return domain.McpKeyContext{}, rerrors.Wrap(pgErr, "get postgres instance")
		}

		pgCtx = &domain.McpKeyPostgresContext{
			Host:     pgInstance.Host,
			Port:     pgInstance.Port,
			Database: vpgDB.V.DatabaseName,
			Username: vpgDB.V.RoleUsername,
			Password: vpgDB.V.RolePassword,
			SSLMode:  pgInstance.SSLMode,
		}
	}

	result := domain.McpKeyContext{
		KeyUuid:   mcpKey.Uuid,
		VaultUuid: mcpKey.VaultUuid,
		UserUuid:  mcpKey.UserUuid,
		CouchURL:  couchInstance.Url,
		CouchDb:   vault.CouchDBName,
		CouchUser: couchInstance.Username,
		CouchPass: couchInstance.Password,
		S3:        s3Ctx,
		Postgres:  pgCtx,

		UseCouchDBForBinaries: vault.UseCouchDBForBinaries,
	}

	// context.Background() is intentional: this must outlive the request context, which
	// net/http cancels as soon as the MCP handler's response is written.
	go func() { //nolint:gosec,contextcheck // see comment above
		touchErr := s.mcpKeys.TouchLastAccessed(context.Background(), mcpKey.Uuid)
		if touchErr != nil {
			log.Error().
				Err(touchErr).
				Str("key_uuid", mcpKey.Uuid.String()).
				Msg("ResolveKey: touch last accessed failed")
		}
	}()

	return result, nil
}
