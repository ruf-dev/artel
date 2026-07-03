package tract

import (
	"context"
	"time"

	"go.redsock.ru/rerrors"
)

// SweepStaleRuns marks every run/step still 'running' with a start time older than threshold
// as 'failed' — v1 crash honesty: if the process died mid-run, startup marks history
// straight instead of leaving 'running' rows lying forever. Call once at app init.
func (s *Service) SweepStaleRuns(ctx context.Context, threshold time.Time) error {
	err := s.tracts.SweepStaleRuns(ctx, threshold)
	if err != nil {
		return rerrors.Wrap(err, "error sweeping stale tract runs")
	}

	err = s.tracts.SweepStaleRunSteps(ctx, threshold)
	if err != nil {
		return rerrors.Wrap(err, "error sweeping stale tract run steps")
	}
	return nil
}
