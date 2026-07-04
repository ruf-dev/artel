package gitlab_webhook

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/utils"
	"go.redsock.ru/rerrors"
)

const (
	pathPrefix         = "/webhooks/gitlab/"
	maxBodyBytes       = 1 << 20 // 1MB
	loggedPayloadBytes = 500
)

type Handler struct {
	externalConns repository.ExternalConnectionRepo
	momSvc        service.MomService
}

func New(externalConns repository.ExternalConnectionRepo, momSvc service.MomService) *Handler {
	return &Handler{
		externalConns: externalConns,
		momSvc:        momSvc,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	rawId := strings.TrimPrefix(r.URL.Path, pathPrefix)

	exConnId, err := uuid.Parse(rawId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	exConn, err := h.externalConns.GetByID(ctx, exConnId)
	if err != nil {
		wrappedErr := rerrors.Wrap(user_errors.GitlabWebhookConnectionNotFound, "error loading external connection")
		log.Error().Err(wrappedErr).Str("external_connection_id", exConnId.String()).Msg("gitlab webhook: connection lookup failed")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if exConn.Provider != domain.ProviderGitlab {
		log.Error().Err(user_errors.GitlabWebhookConnectionNotFound).Str("external_connection_id", exConnId.String()).Msg("gitlab webhook: connection is not a gitlab connection")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	var creds domain.GitlabCredentials

	err = json.Unmarshal(exConn.CredentialsJSON, &creds)
	if err != nil {
		wrappedErr := rerrors.Wrap(user_errors.GitlabWebhookConnectionNotFound, "error unmarshalling gitlab credentials")
		log.Error().Err(wrappedErr).Str("external_connection_id", exConnId.String()).Msg("gitlab webhook: invalid credentials")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	token := r.Header.Get("X-Gitlab-Token")

	secretMatches := creds.WebhookSecret != "" && subtle.ConstantTimeCompare([]byte(token), []byte(creds.WebhookSecret)) == 1
	if !secretMatches {
		log.Warn().Str("external_connection_id", exConnId.String()).Msg(user_errors.GitlabWebhookSecretMismatch.Error())
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	event := r.Header.Get("X-Gitlab-Event")
	log.Info().Str("external_connection_id", exConnId.String()).Str("event", event).Msg("gitlab webhook: received event")

	defer utils.CloseWithLog(r.Body, "gitlab webhook request body")

	limitedBody := io.LimitReader(r.Body, maxBodyBytes)

	body, err := io.ReadAll(limitedBody)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", exConnId.String()).Msg("gitlab webhook: failed to read body")
		w.WriteHeader(http.StatusOK)

		return
	}

	payloadPreview := body
	if len(payloadPreview) > loggedPayloadBytes {
		payloadPreview = payloadPreview[:loggedPayloadBytes]
	}

	log.Debug().Str("external_connection_id", exConnId.String()).Str("event", event).Bytes("payload_preview", payloadPreview).Msg("gitlab webhook: payload received")

	w.WriteHeader(http.StatusOK)
}
