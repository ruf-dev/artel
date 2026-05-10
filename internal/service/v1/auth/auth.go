package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
)

type Service struct {
	signingKey []byte
}

func New(signingKey []byte) *Service {
	return &Service{
		signingKey: signingKey,
	}
}

type tokenPayload struct {
	UserUuid  uuid.UUID `json:"user_uuid"`
	ExpiresAt int64     `json:"exp"`
}

func (s *Service) GenerateToken(userUuid uuid.UUID) (string, error) {
	payload := tokenPayload{
		UserUuid:  userUuid,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", rerrors.Wrap(err, "marshal token payload")
	}

	encoded := base64.RawURLEncoding.EncodeToString(payloadBytes)

	mac := hmac.New(sha256.New, s.signingKey)
	_, err = mac.Write([]byte(encoded))
	if err != nil {
		return "", rerrors.Wrap(err, "compute token signature")
	}

	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return encoded + "." + sig, nil
}

func (s *Service) AuthWithToken(ctx context.Context, token string) (uuid.UUID, error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return uuid.UUID{}, rerrors.New("invalid token format")
	}

	mac := hmac.New(sha256.New, s.signingKey)
	_, err := mac.Write([]byte(parts[0]))
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "compute token signature")
	}

	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[1]), []byte(expectedSig)) {
		return uuid.UUID{}, rerrors.New("invalid token signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "decode token payload")
	}

	var payload tokenPayload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "unmarshal token payload")
	}

	if time.Now().Unix() > payload.ExpiresAt {
		return uuid.UUID{}, rerrors.New("token expired")
	}

	return payload.UserUuid, nil
}
