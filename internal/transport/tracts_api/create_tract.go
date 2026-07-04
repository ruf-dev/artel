package tracts_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) CreateTract(ctx context.Context, req *pb.CreateTract_Request) (*pb.CreateTract_Response, error) {
	def, err := definitionFromJSON(req.Definition)
	if err != nil {
		return nil, err
	}

	created, warnings, err := t.tractSvc.CreateTract(ctx, req.Name, req.Description, def)
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating tract")
	}

	item, err := tractToProto(created)
	if err != nil {
		return nil, err
	}

	resp := &pb.CreateTract_Response{
		Tract:    item,
		Warnings: warnings,
	}

	return resp, nil
}
