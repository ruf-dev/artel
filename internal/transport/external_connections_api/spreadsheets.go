package external_connections_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

func (e *ExternalConnectionsImpl) AddSpreadsheet(
	ctx context.Context,
	req *pb.AddSpreadsheet_Request,
) (*pb.AddSpreadsheet_Response, error) {
	spreadsheet, err := e.svc.AddSpreadsheet(ctx, req.SpreadsheetId, req.Name)
	if err != nil {
		return nil, rerrors.Wrap(err, "error adding spreadsheet")
	}

	resp := &pb.AddSpreadsheet_Response{
		Spreadsheet: spreadsheetToProto(spreadsheet),
	}

	return resp, nil
}

func (e *ExternalConnectionsImpl) ListSpreadsheets(
	ctx context.Context,
	_ *pb.ListSpreadsheets_Request,
) (*pb.ListSpreadsheets_Response, error) {
	spreadsheets, err := e.svc.ListSpreadsheets(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing spreadsheets")
	}

	items := make([]*pb.Spreadsheet, len(spreadsheets))
	for i, s := range spreadsheets {
		items[i] = spreadsheetToProto(s)
	}

	resp := &pb.ListSpreadsheets_Response{
		Spreadsheets: items,
	}

	return resp, nil
}

func (e *ExternalConnectionsImpl) RemoveSpreadsheet(
	ctx context.Context,
	req *pb.RemoveSpreadsheet_Request,
) (*pb.RemoveSpreadsheet_Response, error) {
	err := e.svc.RemoveSpreadsheet(ctx, req.SpreadsheetId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error removing spreadsheet")
	}

	return &pb.RemoveSpreadsheet_Response{}, nil
}

func spreadsheetToProto(s domain.McpSpreadsheet) *pb.Spreadsheet {
	return &pb.Spreadsheet{
		Id:        s.SpreadsheetId,
		Name:      s.Name,
		CreatedAt: s.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}
