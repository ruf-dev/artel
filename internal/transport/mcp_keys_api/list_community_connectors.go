package mcp_keys_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

func (m *McpKeysImpl) ListCommunityConnectors(
	ctx context.Context,
	_ *pb.ListCommunityConnectors_Request,
) (*pb.ListCommunityConnectors_Response, error) {
	candidates, err := m.mcpSvc.ListCommunityConnectors(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list community connectors")
	}

	out := make([]*pb.CommunityConnectorInfo, 0, len(candidates))

	for _, c := range candidates {
		connector := communityConnectorToProto(c)
		out = append(out, connector)
	}

	response := &pb.ListCommunityConnectors_Response{Connectors: out}

	return response, nil
}

func communityConnectorToProto(c domain.MomCandidate) *pb.CommunityConnectorInfo {
	tools := make([]*pb.McpToolInfo, 0, len(c.Tools))

	for _, t := range c.Tools {
		tools = append(tools, momToolToProto(t))
	}

	connector := &pb.CommunityConnectorInfo{
		Name:              c.Name,
		Author:            c.Author,
		Description:       c.Description,
		Tools:             tools,
		ViewerIsOwner:     c.ViewerIsOwner,
		ViewerIsConnected: len(c.Connections) > 0,
	}

	return connector
}

// momToolToProto mirrors the tool-mapping block in ListMomCandidates (list_mom_candidates.go) —
// factored out here so both handlers share it instead of duplicating the smtp/imap oneof +
// params mapping.
func momToolToProto(t domain.McpToolDef) *pb.McpToolInfo {
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

	params := make(map[string]*pb.ToolParamDef, len(t.ApiDescription.Properties))

	for name, prop := range t.ApiDescription.Properties {
		paramDef := toolPropertyToProto(prop)
		params[name] = paramDef
	}

	tool.Params = params

	return tool
}
