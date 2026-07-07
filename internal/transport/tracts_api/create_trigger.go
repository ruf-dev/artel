package tracts_api

import (
	"context"
	"encoding/json"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

const webhookPathPrefix = "/tract/hook/"

func (t *TractsImpl) CreateTrigger(
	ctx context.Context,
	req *pb.CreateTrigger_Request,
) (*pb.CreateTrigger_Response, error) {
	config := json.RawMessage(req.Config)
	if len(config) == 0 {
		config = json.RawMessage(`{}`)
	}

	if !json.Valid(config) {
		return nil, rerrors.Wrap(user_errors.TractRequestFieldInvalidJSON, "error validating config")
	}

	schema, err := schemaFromJSON(req.PayloadSchema)
	if err != nil {
		return nil, err
	}

	created, rawToken, err := t.tractSvc.CreateTrigger(ctx, req.Name, req.Kind, req.Source, config, schema)
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating trigger")
	}

	// rawToken is empty for provider-linked triggers (see Service.CreateTrigger) — they share the
	// provider connection's own webhook URL/secret, already shown to the user on the connection
	// management dialog (e.g. ManageGitlabDialog), so there is no fresh URL to reveal here.
	var webhookUrl string
	if rawToken != "" {
		webhookUrl = webhookBaseURL(ctx) + webhookPathPrefix + created.TriggerUuid.String()
	}

	resp := &pb.CreateTrigger_Response{
		Trigger:      triggerToProto(created),
		WebhookUrl:   webhookUrl,
		WebhookToken: rawToken,
	}

	return resp, nil
}
