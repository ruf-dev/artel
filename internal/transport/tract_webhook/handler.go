package tract_webhook

import (
	"context"
	"crypto/sha256"
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
	"github.com/ruf-dev/artel/internal/service/v1/tract"
	"github.com/ruf-dev/artel/internal/utils"
)

const (
	pathPrefix   = "/tract/hook/"
	maxBodyBytes = 1 << 20 // 1MB
)

// Handler serves POST /tract/hook/{trigger_uuid}, modeled on
// internal/transport/gitlab_webhook/handler.go: it authenticates the request directly against
// the trigger's own secret (not a user session — there is no user_context on a webhook call),
// so it depends on repository.TriggersRepo rather than going through service.TractService for
// the lookup/token-compare step. StartRun is spawned with baseCtx (the server-lifecycle
// context, not the request context, which is cancelled the moment ServeHTTP returns) so a
// fired run keeps executing after the HTTP response is sent.
type Handler struct {
	triggers repository.TriggersRepo
	tractSvc service.TractService
	baseCtx  context.Context
}

func New(baseCtx context.Context, triggers repository.TriggersRepo, tractSvc service.TractService) *Handler {
	return &Handler{
		triggers: triggers,
		tractSvc: tractSvc,
		baseCtx:  baseCtx,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	rawId := strings.TrimPrefix(r.URL.Path, pathPrefix)

	triggerUuid, err := uuid.Parse(rawId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	result, err := h.triggers.GetByTriggerUuid(ctx, triggerUuid)
	if err != nil {
		log.Error().Err(err).Str("trigger_uuid", triggerUuid.String()).Msg("tract webhook: trigger lookup failed")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if !result.Valid {
		log.Warn().Str("trigger_uuid", triggerUuid.String()).Msg("tract webhook: no trigger found for this webhook id")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	trigger := result.V

	token := r.Header.Get("X-Tract-Token")
	if token == "" {
		token = r.Header.Get("X-Gitlab-Token")
	}

	tokenHash := sha256.Sum256([]byte(token))

	secretMatches := token != "" && subtle.ConstantTimeCompare(tokenHash[:], trigger.SecretHash) == 1
	if !secretMatches {
		log.Warn().
			Str("trigger_uuid", triggerUuid.String()).
			Msg("tract webhook: token does not match configured secret")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	defer utils.CloseWithLog(r.Body, "tract webhook request body")

	limitedBody := io.LimitReader(r.Body, maxBodyBytes)

	body, err := io.ReadAll(limitedBody)
	if err != nil {
		log.Error().Err(err).Str("trigger_uuid", triggerUuid.String()).Msg("tract webhook: failed to read body")
		w.WriteHeader(http.StatusOK)

		return
	}

	// Always 200 from here on — the response never leaks whether the trigger is disabled or
	// which (if any) linked tracts matched their filters.
	if !trigger.Enabled {
		w.WriteHeader(http.StatusOK)

		return
	}

	normalized, err := tract.NormalizePayload(trigger.Source, body)
	if err != nil {
		log.Error().Err(err).Str("trigger_uuid", triggerUuid.String()).Msg("tract webhook: failed to normalize payload")
		w.WriteHeader(http.StatusOK)

		return
	}

	links, err := h.triggers.ListLinksByTrigger(ctx, trigger.Uuid)
	if err != nil {
		log.Error().
			Err(err).
			Str("trigger_uuid", triggerUuid.String()).
			Msg("tract webhook: failed to list linked tracts")
		w.WriteHeader(http.StatusOK)

		return
	}

	for _, link := range links {
		if !link.Tract.Enabled {
			continue
		}

		matched, err := tract.EvaluateTriggerFilters(normalized, link.Filters)
		if err != nil {
			log.Error().
				Err(err).
				Str("trigger_uuid", triggerUuid.String()).
				Str("tract_uuid", link.Tract.Uuid.String()).
				Msg("tract webhook: failed to evaluate link filters")

			continue
		}

		if !matched {
			continue
		}

		linkedTract := link.Tract
		go h.startRun(linkedTract, trigger.Uuid, normalized)
	}

	w.WriteHeader(http.StatusOK)
}

// startRun runs against h.baseCtx (the server-lifecycle context, set at construction) rather
// than the request context, which is cancelled the moment ServeHTTP returns.
func (h *Handler) startRun(linkedTract domain.Tract, triggerUuid uuid.UUID, payload json.RawMessage) {
	_, err := h.tractSvc.StartRun(h.baseCtx, linkedTract, payload, tract.StartedByWebhook, triggerUuid)
	if err != nil {
		log.Error().
			Err(err).
			Str("tract_uuid", linkedTract.Uuid.String()).
			Str("trigger_uuid", triggerUuid.String()).
			Msg("tract webhook: run failed")
	}
}
