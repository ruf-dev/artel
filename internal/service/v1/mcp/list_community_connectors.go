package mcp

import (
	"context"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// ListCommunityConnectors returns admin-authored community MoMs (IsCommunity == true) only —
// system MoMs seeded via migration are excluded. Shares its connection-grouping and candidate-
// building helpers with ListMomCandidates; the two only differ in which mcpDefinitions.List()
// rows they turn into candidates.
func (s *ServiceImpl) ListCommunityConnectors(ctx context.Context) ([]domain.MomCandidate, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	defs, err := s.mcpDefinitions.List(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mcp definitions")
	}

	conns, err := s.externalConnections.ListByUser(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "list external connections")
	}

	connsByProvider := groupConnectionsByProvider(conns)

	candidates := make([]domain.MomCandidate, 0, len(defs))

	for _, def := range defs {
		if !def.IsCommunity {
			continue
		}

		candidates = append(candidates, buildMomCandidate(def, connsByProvider, uc.UserUuid))
	}

	sortMomCandidates(candidates)

	return candidates, nil
}
