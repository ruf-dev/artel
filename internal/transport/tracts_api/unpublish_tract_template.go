package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) UnpublishTractTemplate(
	ctx context.Context,
	req *pb.UnpublishTractTemplate_Request,
) (*pb.UnpublishTractTemplate_Response, error) {
	templateUuid, err := uuid.Parse(req.TemplateUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing template uuid")
	}

	err = t.tractSvc.UnpublishTemplate(ctx, templateUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error unpublishing tract template")
	}

	return &pb.UnpublishTractTemplate_Response{}, nil
}
