package tract

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
)

// DispatchTrigger runs trigger's full delivery pipeline against rawBody: normalize the payload,
// list linked tracts, evaluate each link's filters, and spawn StartRun (on baseCtx, so runs
// outlive the inbound HTTP request) for every match. Shared by tract_webhook (one trigger, its
// own webhook URL) and gitlab_webhook (N triggers fanned out from one shared provider
// connection) so neither transport handler duplicates this loop.
func DispatchTrigger(
	ctx context.Context,
	baseCtx context.Context,
	triggers repository.TriggersRepo,
	tractSvc service.TractService,
	trigger domain.Trigger,
	rawBody json.RawMessage,
) error {
	if !trigger.Enabled {
		return nil
	}

	normalized, err := NormalizePayload(trigger.Source, rawBody)
	if err != nil {
		return err
	}

	links, err := triggers.ListLinksByTrigger(ctx, trigger.Uuid)
	if err != nil {
		return err
	}

	for _, link := range links {
		if !link.Tract.Enabled {
			continue
		}

		matched, matchErr := EvaluateTriggerFilters(normalized, link.Filters)
		if matchErr != nil {
			log.Error().
				Err(matchErr).
				Str("trigger_uuid", trigger.Uuid.String()).
				Str("tract_uuid", link.Tract.Uuid.String()).
				Msg("dispatch trigger: failed to evaluate link filters")

			continue
		}

		if !matched {
			continue
		}

		linkedTract := link.Tract
		go startRun(baseCtx, tractSvc, linkedTract, trigger.Uuid, normalized)
	}

	return nil
}

// startRun runs against baseCtx (the server-lifecycle context) rather than the request context,
// which is cancelled the moment the transport handler's ServeHTTP returns.
func startRun(baseCtx context.Context, tractSvc service.TractService, linkedTract domain.Tract, triggerUuid uuid.UUID, payload json.RawMessage) {
	_, err := tractSvc.StartRun(baseCtx, linkedTract, payload, StartedByWebhook, triggerUuid)
	if err != nil {
		log.Error().
			Err(err).
			Str("tract_uuid", linkedTract.Uuid.String()).
			Str("trigger_uuid", triggerUuid.String()).
			Msg("dispatch trigger: run failed")
	}
}
