package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) PublishTractTemplate(
	ctx context.Context,
	req *pb.PublishTractTemplate_Request,
) (*pb.PublishTractTemplate_Response, error) {
	tractUuid, err := uuid.Parse(req.TractUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tract uuid")
	}

	template, err := t.tractSvc.PublishTemplate(ctx, tractUuid, req.Category)
	if err != nil {
		return nil, rerrors.Wrap(err, "error publishing tract template")
	}

	resp := &pb.PublishTractTemplate_Response{
		Template: tractTemplateToProto(template),
	}

	return resp, nil
}
