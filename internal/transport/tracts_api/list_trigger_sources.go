package tracts_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (t *TractsImpl) ListTriggerSources(ctx context.Context, _ *pb.ListTriggerSources_Request) (*pb.ListTriggerSources_Response, error) {
	sources, err := t.tractSvc.ListTriggerSources(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing trigger sources")
	}

	resp := &pb.ListTriggerSources_Response{Sources: triggerSourcesToProto(sources)}
	return resp, nil
}
