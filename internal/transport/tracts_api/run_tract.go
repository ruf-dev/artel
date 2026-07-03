package tracts_api

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/tract"
)

// RunTract fires a manual run and returns immediately — the same fire-and-forget shape as the
// tract webhook: the run's outcome is visible via ListRuns/GetRun. The engine is spawned
// against t.baseCtx (the server-lifecycle context), not ctx, which is cancelled the moment
// this RPC returns.
func (t *TractsImpl) RunTract(ctx context.Context, req *pb.RunTract_Request) (*pb.RunTract_Response, error) {
	id, err := uuid.Parse(req.TractUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tract uuid")
	}

	params := []byte(req.Params)
	if len(params) > 0 && !json.Valid(params) {
		return nil, rerrors.Wrap(user_errors.TractRequestFieldInvalidJSON, "error validating run params")
	}

	found, err := t.tractSvc.GetTract(ctx, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting tract")
	}

	// found is captured here (fetched via the authenticated request context) rather than
	// re-fetched inside the goroutine — t.baseCtx carries no user_context, so a GetTract call
	// against it would fail the ownership check.
	go t.runManual(found, params)

	return &pb.RunTract_Response{}, nil
}

func (t *TractsImpl) runManual(tr domain.Tract, params json.RawMessage) {
	_, err := t.tractSvc.StartRun(t.baseCtx, tr, params, tract.StartedByManual, uuid.Nil)
	if err != nil {
		log.Error().Err(err).Str("tract_uuid", tr.Uuid.String()).Msg("manual tract run failed")
	}
}
