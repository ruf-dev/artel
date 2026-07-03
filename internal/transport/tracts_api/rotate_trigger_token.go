package tracts_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (t *TractsImpl) RotateTriggerToken(ctx context.Context, req *pb.RotateTriggerToken_Request) (*pb.RotateTriggerToken_Response, error) {
	id, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing trigger uuid")
	}

	rotated, rawToken, err := t.tractSvc.RotateTriggerToken(ctx, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error rotating trigger token")
	}

	webhookUrl := webhookBaseURL(ctx) + webhookPathPrefix + rotated.TriggerUuid.String()

	resp := &pb.RotateTriggerToken_Response{
		Trigger:      triggerToProto(rotated),
		WebhookUrl:   webhookUrl,
		WebhookToken: rawToken,
	}
	return resp, nil
}
