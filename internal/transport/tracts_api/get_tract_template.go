package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) GetTractTemplate(
	ctx context.Context,
	req *pb.GetTractTemplate_Request,
) (*pb.GetTractTemplate_Response, error) {
	id, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing template uuid")
	}

	template, err := t.tractSvc.GetTemplate(ctx, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting tract template")
	}

	resp := &pb.GetTractTemplate_Response{
		Template: tractTemplateToProto(template),
	}

	return resp, nil
}
