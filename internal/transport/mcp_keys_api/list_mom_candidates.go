package mcp_keys_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/transport/external_connections_api"
	"go.redsock.ru/rerrors"
)

func (m *McpKeysImpl) ListMomCandidates(
	ctx context.Context,
	_ *pb.ListMomCandidates_Request,
) (*pb.ListMomCandidates_Response, error) {
	candidates, err := m.mcpSvc.ListMomCandidates(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mom candidates")
	}

	out := make([]*pb.MomCandidate, 0, len(candidates))

	for _, c := range candidates {
		connections := make([]*pb.ExternalConnectionInfo, 0, len(c.Connections))
		for _, conn := range c.Connections {
			connections = append(connections, external_connections_api.ConnectionToProto(conn))
		}

		tools := make([]*pb.McpToolInfo, 0, len(c.Tools))

		for _, t := range c.Tools {
			tools = append(tools, momToolToProto(t))
		}

		candidate := &pb.MomCandidate{
			Name:        c.Name,
			Author:      c.Author,
			Description: c.Description,
			Connections: connections,
			Tools:       tools,
		}
		out = append(out, candidate)
	}

	return &pb.ListMomCandidates_Response{Candidates: out}, nil
}

func toolPropertyToProto(prop domain.ToolProperty) *pb.ToolParamDef {
	paramDef := &pb.ToolParamDef{Description: prop.Description}

	if len(prop.Enum) > 0 {
		enumParam := &pb.EnumParam{Values: prop.Enum}
		enumKind := &pb.ToolParamDef_EnumParam{EnumParam: enumParam}
		paramDef.Kind = enumKind

		return paramDef
	}

	if prop.Type == "integer" || prop.Type == "number" {
		intParam := &pb.IntegerParam{}
		intKind := &pb.ToolParamDef_IntegerParam{IntegerParam: intParam}
		paramDef.Kind = intKind

		return paramDef
	}

	strParam := &pb.StringParam{}
	strKind := &pb.ToolParamDef_StringParam{StringParam: strParam}
	paramDef.Kind = strKind

	return paramDef
}

func imapOperationToProto(op domain.ImapOperation) pb.ImapOperation {
	switch op {
	case domain.IMAP_OP_LIST_FOLDERS:
		return pb.ImapOperation_IMAP_OP_LIST_FOLDERS
	case domain.IMAP_OP_LIST_MESSAGES:
		return pb.ImapOperation_IMAP_OP_LIST_MESSAGES
	case domain.IMAP_OP_FETCH_MESSAGE:
		return pb.ImapOperation_IMAP_OP_FETCH_MESSAGE
	default:
		return pb.ImapOperation_IMAP_OP_UNSPECIFIED
	}
}
