package mcp_keys_api

import (
	"context"

	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
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

		out = append(out, &pb.MomCandidate{
			Name:        c.Name,
			Author:      c.Author,
			Description: c.Description,
			Connections: connections,
		})
	}

	return &pb.ListMomCandidates_Response{Candidates: out}, nil
}
