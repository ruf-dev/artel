package tracts_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (t *TractsImpl) ListTriggers(ctx context.Context, _ *pb.ListTriggers_Request) (*pb.ListTriggers_Response, error) {
	triggers, err := t.tractSvc.ListTriggers(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing triggers")
	}

	resp := &pb.ListTriggers_Response{Triggers: triggersToProto(triggers)}
	return resp, nil
}
