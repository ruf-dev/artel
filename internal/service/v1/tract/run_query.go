package tract

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// ListRuns returns the most recent runs of tractUuid (most recent first), capped at limit —
// the tract card's run history / "last run" badge. tractUuid must belong to the caller.
func (s *Service) ListRuns(ctx context.Context, tractUuid uuid.UUID, limit int32) ([]domain.TractRun, error) {
	_, err := s.GetTract(ctx, tractUuid)
	if err != nil {
		return nil, err
	}

	runs, err := s.tracts.ListRunsByTract(ctx, tractUuid, limit)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tract runs")
	}
	return runs, nil
}

// GetRun returns one run plus its ordered step rows — the run detail view (logs, "last
// output" per node keyed by step_id). Ownership is checked via the run's tract.
func (s *Service) GetRun(ctx context.Context, id uuid.UUID) (domain.TractRun, []domain.TractRunStep, error) {
	result, err := s.tracts.GetRun(ctx, id)
	if err != nil {
		return domain.TractRun{}, nil, rerrors.Wrap(err, "error getting tract run")
	}
	if !result.Valid {
		return domain.TractRun{}, nil, rerrors.Wrap(user_errors.NotFound)
	}

	_, err = s.GetTract(ctx, result.V.TractUuid)
	if err != nil {
		return domain.TractRun{}, nil, err
	}

	steps, err := s.tracts.ListRunStepsByRun(ctx, id)
	if err != nil {
		return domain.TractRun{}, nil, rerrors.Wrap(err, "error listing tract run steps")
	}

	return result.V, steps, nil
}
