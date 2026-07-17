package tracts_api

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// RetryRun fires a new run that replays runUuid's own trigger_uuid/trigger_payload/started_by
// verbatim — no new Trigger row is ever created, so the original and the retried run end up
// pointing at the same trigger (or both uuid.Nil, for a retried manual run). Same fire-and-
// forget shape as RunTract: the run's outcome is visible via ListRuns/GetRun.
func (t *TractsImpl) RetryRun(ctx context.Context, req *pb.RetryRun_Request) (*pb.RetryRun_Response, error) {
	runUuid, err := uuid.Parse(req.RunUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing run uuid")
	}

	run, _, err := t.tractSvc.GetRun(ctx, runUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting run")
	}

	// found is captured here (fetched via the authenticated request context) rather than
	// re-fetched inside the goroutine — t.baseCtx carries no user_context, so a GetTract call
	// against it would fail the ownership check.
	found, err := t.tractSvc.GetTract(ctx, run.TractUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting tract")
	}

	go t.runRetry(found, run)

	return &pb.RetryRun_Response{}, nil
}

func (t *TractsImpl) runRetry(tr domain.Tract, run domain.TractRun) {
	_, err := t.tractSvc.StartRun(t.baseCtx, tr, run.TriggerPayload, run.StartedBy, run.TriggerUuid)
	if err != nil {
		log.Error().
			Err(err).
			Str("tract_uuid", tr.Uuid.String()).
			Str("retry_of_run_uuid", run.Uuid.String()).
			Msg("retry tract run failed")
	}
}
