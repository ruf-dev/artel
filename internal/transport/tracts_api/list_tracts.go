package tracts_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

const lastRunLimit = 1

// ListTracts assembles each card's linked-trigger summaries and last-run badge from the
// per-tract link/run listings service.TractService already exposes — no new aggregation
// method needed on the service.
func (t *TractsImpl) ListTracts(ctx context.Context, _ *pb.ListTracts_Request) (*pb.ListTracts_Response, error) {
	tracts, err := t.tractSvc.ListTracts(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tracts")
	}

	items := make([]*pb.TractItem, len(tracts))

	for i, tr := range tracts {
		item, err := t.tractWithSummary(ctx, tr)
		if err != nil {
			return nil, err
		}

		items[i] = item
	}

	resp := &pb.ListTracts_Response{Tracts: items}

	return resp, nil
}

func (t *TractsImpl) tractWithSummary(ctx context.Context, tr domain.Tract) (*pb.TractItem, error) {
	links, err := t.tractSvc.ListLinksByTract(ctx, tr.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tract trigger links")
	}

	runs, err := t.tractSvc.ListRuns(ctx, tr.Uuid, lastRunLimit)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tract runs")
	}

	var lastRun *domain.TractRun
	if len(runs) > 0 {
		lastRun = &runs[0]
	}

	return tractToProtoWithSummary(tr, links, lastRun)
}
