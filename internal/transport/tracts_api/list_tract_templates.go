package tracts_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) ListTractTemplates(
	ctx context.Context,
	req *pb.ListTractTemplates_Request,
) (*pb.ListTractTemplates_Response, error) {
	templates, err := t.tractSvc.ListTemplates(ctx, req.Category, req.MineOnly)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tract templates")
	}

	resp := &pb.ListTractTemplates_Response{
		Templates: tractTemplateSummariesToProto(templates),
	}

	return resp, nil
}
