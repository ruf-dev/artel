package tracts_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (t *TractsImpl) ListTractTools(ctx context.Context, _ *pb.ListTractTools_Request) (*pb.ListTractTools_Response, error) {
	tools, err := t.tractSvc.ListTractTools(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tract tools")
	}

	resp := &pb.ListTractTools_Response{Tools: toolRefsToProto(tools)}
	return resp, nil
}
