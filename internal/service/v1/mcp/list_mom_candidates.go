package mcp

import (
	"context"
	"sort"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (s *ServiceImpl) ListMomCandidates(ctx context.Context) ([]domain.MomCandidate, error) {
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
		candidates = append(candidates, buildMomCandidate(def, connsByProvider, uc.UserUuid))
	}

	sortMomCandidates(candidates)

	return candidates, nil
}

// groupConnectionsByProvider indexes conns (the caller's own external connections) by
// provider — shared by ListMomCandidates and ListCommunityConnectors, which only differ in
// which mcpDefinitions.List() rows they turn into candidates.
func groupConnectionsByProvider(conns []domain.ExternalConnection) map[string][]domain.ExternalConnectionMeta {
	connsByProvider := make(map[string][]domain.ExternalConnectionMeta)

	for _, c := range conns {
		connsByProvider[c.Provider] = append(connsByProvider[c.Provider], domain.ExternalConnectionMeta{
			Uuid:         c.Uuid,
			Provider:     c.Provider,
			ProviderType: c.ProviderType,
			Metadata:     c.Metadata,
			CreatedAt:    c.CreatedAt,
			UpdatedAt:    c.UpdatedAt,
		})
	}

	return connsByProvider
}

// buildMomCandidate pairs def with whichever of the caller's connections satisfy its tools'
// required providers, and marks ViewerIsOwner against callerUuid.
func buildMomCandidate(
	def domain.McpDefinition,
	connsByProvider map[string][]domain.ExternalConnectionMeta,
	callerUuid uuid.UUID,
) domain.MomCandidate {
	var matched []domain.ExternalConnectionMeta
	for _, provider := range requiredProviders(def) {
		matched = append(matched, connsByProvider[provider]...)
	}

	candidate := domain.MomCandidate{
		Name:          def.Name,
		Author:        def.Author,
		Description:   def.Description,
		Connections:   matched,
		Tools:         def.Tools,
		OwnerUserUuid: def.OwnerUserUuid,
		ViewerIsOwner: def.OwnerUserUuid != nil && *def.OwnerUserUuid == callerUuid,
	}

	return candidate
}

func sortMomCandidates(candidates []domain.MomCandidate) {
	sort.SliceStable(candidates, func(i, j int) bool {
		iHas, jHas := len(candidates[i].Connections) > 0, len(candidates[j].Connections) > 0
		if iHas != jHas {
			return iHas
		}

		return candidates[i].Name < candidates[j].Name
	})
}

// requiredProviders returns the distinct external_connections.provider values
// that def's tools need credentials for.
func requiredProviders(def domain.McpDefinition) []string {
	seen := make(map[string]bool)

	var providers []string

	add := func(provider string) {
		if provider == "" || seen[provider] {
			return
		}

		seen[provider] = true

		providers = append(providers, provider)
	}

	for _, tool := range def.Tools {
		switch {
		case tool.Action.Imap != nil, tool.Action.Smtp != nil:
			add(domain.ProviderEmail)
		case tool.Action.Http != nil:
			add(tool.Action.Http.Credentials)
		}
	}

	return providers
}
