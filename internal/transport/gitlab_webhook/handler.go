package gitlab_webhook

import (
	"context"
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
	"github.com/ruf-dev/artel/internal/service/v1/tract"
	"github.com/ruf-dev/artel/internal/utils"
	"go.redsock.ru/rerrors"
)

const (
	pathPrefix         = "/webhooks/gitlab/"
	maxBodyBytes       = 1 << 20 // 1MB
	loggedPayloadBytes = 500
)

// Handler serves POST /webhooks/gitlab/{external_connection_id} — the shared inbound endpoint
// for every trigger linked to one GitLab connection (see trigger_provider_links). Unlike
// tract_webhook.Handler (one trigger, its own secret/URL), auth here is against the connection's
// own webhook secret, and a single delivery may fan out to N triggers, each gated by its own
// TriggerMatchers (tract.MatchesRequest) — e.g. a gitlab_push trigger only fires on
// X-Gitlab-Event: Push Hook, so a Merge Request delivery to the same connection URL doesn't
// misfire it.
type Handler struct {
	externalConns repository.ExternalConnectionRepo
	triggers      repository.TriggersRepo
	momSvc        service.MomService
	tractSvc      service.TractService
	baseCtx       context.Context
}

func New(
	baseCtx context.Context,
	externalConns repository.ExternalConnectionRepo,
	triggers repository.TriggersRepo,
	momSvc service.MomService,
	tractSvc service.TractService,
) *Handler {
	return &Handler{
		externalConns: externalConns,
		triggers:      triggers,
		momSvc:        momSvc,
		tractSvc:      tractSvc,
		baseCtx:       baseCtx,
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
		log.Error().
			Err(wrappedErr).
			Str("external_connection_id", exConnId.String()).
			Msg("gitlab webhook: connection lookup failed")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if exConn.Provider != domain.ProviderGitlab {
		log.Error().
			Err(user_errors.GitlabWebhookConnectionNotFound).
			Str("external_connection_id", exConnId.String()).
			Msg("gitlab webhook: connection is not a gitlab connection")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	var creds domain.GitlabCredentials

	err = json.Unmarshal(exConn.CredentialsJSON, &creds)
	if err != nil {
		wrappedErr := rerrors.Wrap(
			user_errors.GitlabWebhookConnectionNotFound,
			"error unmarshalling gitlab credentials",
		)
		log.Error().
			Err(wrappedErr).
			Str("external_connection_id", exConnId.String()).
			Msg("gitlab webhook: invalid credentials")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	token := r.Header.Get("X-Gitlab-Token")

	secretMatches := creds.WebhookSecret != "" &&
		subtle.ConstantTimeCompare([]byte(token), []byte(creds.WebhookSecret)) == 1
	if !secretMatches {
		log.Warn().Str("external_connection_id", exConnId.String()).Msg(user_errors.GitlabWebhookSecretMismatch.Error())
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	event := r.Header.Get("X-Gitlab-Event")
	log.Info().
		Str("external_connection_id", exConnId.String()).
		Str("event", event).
		Msg("gitlab webhook: received event")

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

	log.Debug().
		Str("external_connection_id", exConnId.String()).
		Str("event", event).
		Bytes("payload_preview", payloadPreview).
		Msg("gitlab webhook: payload received")

	// Always 200 from here on — same "never leak internal state" posture as tract_webhook.
	linkedTriggers, err := h.triggers.ListByExternalConnection(ctx, exConnId)
	if err != nil {
		log.Error().
			Err(err).
			Str("external_connection_id", exConnId.String()).
			Msg("gitlab webhook: failed to list linked triggers")
		w.WriteHeader(http.StatusOK)

		return
	}

	for _, trigger := range linkedTriggers {
		if !tract.MatchesRequest(trigger.Matchers, r.Header) {
			continue
		}

		dispatchErr := tract.DispatchTrigger(ctx, h.baseCtx, h.triggers, h.tractSvc, trigger, body)
		if dispatchErr != nil {
			log.Error().
				Err(dispatchErr).
				Str("external_connection_id", exConnId.String()).
				Str("trigger_uuid", trigger.Uuid.String()).
				Msg("gitlab webhook: dispatch failed")
		}
	}

	w.WriteHeader(http.StatusOK)
}
