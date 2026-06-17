package mcp_keys_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/transport/external_connections_api"
)

func (m *McpKeysImpl) ListMomCandidates(ctx context.Context, _ *pb.ListMomCandidates_Request) (*pb.ListMomCandidates_Response, error) {
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
			tool := &pb.McpToolInfo{
				Name:        t.ApiDescription.Name,
				Description: t.ApiDescription.Description,
			}
			if t.Action.Smtp != nil {
				op := pb.SmtpOperation_SMTP_OP_UNSPECIFIED
				if t.Action.Smtp.Operation == domain.SMTP_OP_SEND {
					op = pb.SmtpOperation_SMTP_OP_SEND
				}
				smtpAction := &pb.SmtpToolAction{Operation: op}
				smtpWrapper := &pb.McpToolInfo_Smtp{Smtp: smtpAction}
				tool.Action = smtpWrapper
			}
			if t.Action.Imap != nil {
				op := imapOperationToProto(t.Action.Imap.Operation)
				imapAction := &pb.ImapToolAction{Operation: op}
				imapWrapper := &pb.McpToolInfo_Imap{Imap: imapAction}
				tool.Action = imapWrapper
			}
			tools = append(tools, tool)
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
