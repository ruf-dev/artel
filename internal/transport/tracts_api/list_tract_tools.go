package tracts_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) ListTractTools(ctx context.Context, _ *pb.ListTractTools_Request) (*pb.ListTractTools_Response, error) {
	tools, err := t.tractSvc.ListTractTools(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tract tools")
	}

	resp := &pb.ListTractTools_Response{Tools: toolRefsToProto(tools)}

	return resp, nil
}
