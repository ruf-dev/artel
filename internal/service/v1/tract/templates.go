package tract

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// PublishTemplate snapshots tractUuid's current definition as a new public template — a frozen
// copy (see domain.TractTemplate doc), not a live view. Every action step's ConnectionUuid, and
// every llm_call step's LlmConnectionUuid, is stripped to nil before persisting: both reference
// the publisher's external_connections rows, meaningless (and not ownership-checkable) for
// anyone else.
func (s *Service) PublishTemplate(ctx context.Context, tractUuid uuid.UUID, category string) (domain.TractTemplate, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.TractTemplate{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	source, err := s.GetTract(ctx, tractUuid)
	if err != nil {
		return domain.TractTemplate{}, err
	}

	def, err := deepCopyDefinition(source.Definition)
	if err != nil {
		return domain.TractTemplate{}, err
	}

	err = walkActions(def.Steps, func(step *domain.TractStep) error {
		step.ConnectionUuid = uuid.Nil

		return nil
	})
	if err != nil {
		return domain.TractTemplate{}, err
	}

	err = walkLlmCallStepsMut(def.Steps, func(step *domain.TractStep) error {
		step.LlmConnectionUuid = uuid.Nil

		return nil
	})
	if err != nil {
		return domain.TractTemplate{}, err
	}

	template := domain.TractTemplate{
		SourceTractUuid: source.Uuid,
		OwnerUuid:       uc.UserUuid,
		Name:            source.Name,
		Description:     source.Description,
		Definition:      def,
		Category:        category,
	}

	created, err := s.templates.Create(ctx, template)
	if err != nil {
		return domain.TractTemplate{}, rerrors.Wrap(err, "error creating tract template")
	}

	return created, nil
}

// UnpublishTemplate removes a template — publisher-only. Already-instantiated copies (real
// tracts rows created from it) are unaffected; they're independent rows, not linked back live.
// Built-in templates (OwnerUuid == uuid.Nil, see domain.TractTemplate doc) can never be
// unpublished by any real user: an authenticated caller's UserUuid is never uuid.Nil, so the
// ownership check below already rejects every caller for them — no separate guard needed.
func (s *Service) UnpublishTemplate(ctx context.Context, templateUuid uuid.UUID) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	result, err := s.templates.Get(ctx, templateUuid)
	if err != nil {
		return rerrors.Wrap(err, "error getting tract template")
	}

	if !result.Valid {
		return rerrors.Wrap(user_errors.NotFound)
	}

	if result.V.OwnerUuid != uc.UserUuid {
		return rerrors.Wrap(user_errors.TractTemplateNotOwned)
	}

	err = s.templates.Delete(ctx, templateUuid)
	if err != nil {
		return rerrors.Wrap(err, "error deleting tract template")
	}

	return nil
}

func (s *Service) ListTemplates(ctx context.Context, category string, mineOnly bool) ([]domain.TractTemplate, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	if mineOnly {
		templates, err := s.templates.ListByOwner(ctx, uc.UserUuid)
		if err != nil {
			return nil, rerrors.Wrap(err, "error listing tract templates by owner")
		}

		return templates, nil
	}

	templates, err := s.templates.ListAll(ctx, category)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tract templates")
	}

	return templates, nil
}

func (s *Service) GetTemplate(ctx context.Context, templateUuid uuid.UUID) (domain.TractTemplate, error) {
	_, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.TractTemplate{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	result, err := s.templates.Get(ctx, templateUuid)
	if err != nil {
		return domain.TractTemplate{}, rerrors.Wrap(err, "error getting tract template")
	}

	if !result.Valid {
		return domain.TractTemplate{}, rerrors.Wrap(user_errors.NotFound)
	}

	return result.V, nil
}

// InstantiateTemplate copies templateUuid into a brand-new tract owned by the caller. connections
// maps MoM name -> connection uuid (the caller's own external_connections rows) for every
// non-builtin action step the template uses, plus — reusing the same map — LLM provider name
// (e.g. domain.ProviderAnthropic) -> connection uuid for every llm_call step the template uses;
// neither key space collides since no MoM is ever named after an LLM provider. An action step
// whose Mcp has no entry, or an llm_call step whose provider has no entry, keeps its connection
// uuid nil and createTractInternal's validateTools/validateLlmConnections will reject the whole
// call — instantiate is all-or-nothing, never a partially-wired draft. Trigger links are never
// copied — a fresh copy always starts with zero triggers.
func (s *Service) InstantiateTemplate(
	ctx context.Context, templateUuid uuid.UUID, name, description string, connections map[string]uuid.UUID,
) (domain.Tract, []string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Tract{}, nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	result, err := s.templates.Get(ctx, templateUuid)
	if err != nil {
		return domain.Tract{}, nil, rerrors.Wrap(err, "error getting tract template")
	}

	if !result.Valid {
		return domain.Tract{}, nil, rerrors.Wrap(user_errors.NotFound)
	}

	def, err := deepCopyDefinition(result.V.Definition)
	if err != nil {
		return domain.Tract{}, nil, err
	}

	err = walkActions(def.Steps, func(step *domain.TractStep) error {
		if step.Mcp != builtinMcpName {
			step.ConnectionUuid = connections[step.Mcp]
		}

		return nil
	})
	if err != nil {
		return domain.Tract{}, nil, err
	}

	err = walkLlmCallStepsMut(def.Steps, func(step *domain.TractStep) error {
		step.LlmConnectionUuid = connections[domain.ProviderAnthropic]

		return nil
	})
	if err != nil {
		return domain.Tract{}, nil, err
	}

	created, warnings, err := s.createTractInternal(ctx, uc.UserUuid, name, description, def)
	if err != nil {
		return domain.Tract{}, nil, err
	}

	incErr := s.templates.IncrementInstallCount(ctx, templateUuid)
	if incErr != nil {
		log.Error().
			Err(incErr).
			Str("template_uuid", templateUuid.String()).
			Msg("error incrementing tract template install count")
	}

	return created, warnings, nil
}

// deepCopyDefinition returns an independent copy of def — a JSON marshal/unmarshal round trip is
// the simplest correct deep copy given TractDefinition's nested step slices.
func deepCopyDefinition(def domain.TractDefinition) (domain.TractDefinition, error) {
	data, err := json.Marshal(def)
	if err != nil {
		return domain.TractDefinition{}, rerrors.Wrap(err, "error marshaling tract definition")
	}

	var copyDef domain.TractDefinition

	err = json.Unmarshal(data, &copyDef)
	if err != nil {
		return domain.TractDefinition{}, rerrors.Wrap(err, "error unmarshaling tract definition")
	}

	return copyDef, nil
}
