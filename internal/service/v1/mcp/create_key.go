package mcp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

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
