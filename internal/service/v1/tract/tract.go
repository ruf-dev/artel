package tract

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// builtinMcpName is the Mcp sentinel value for builtin (vault) tool actions — "artel" per the
// plan's step grammar, as opposed to a real MoM name.
const builtinMcpName = "artel"

const (
	stepTypeAction    = "action"
	stepTypeCondition = "condition"
	stepTypeParallel  = "parallel"
	stepTypeGroup     = "group"
)

const reservedStepId = "trigger"

var stepIdPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Service implements service.TractService and (via engine.go) the tract run engine.
type Service struct {
	tracts         repository.TractsRepo
	triggers       repository.TriggersRepo
	triggerPresets repository.TriggerPresetsRepo
	externalConns  repository.ExternalConnectionRepo
	mcpDefs        repository.McpDefinitionsRepo
	toolExecutor   ToolExecutor
}

func New(
	tracts repository.TractsRepo,
	triggers repository.TriggersRepo,
	triggerPresets repository.TriggerPresetsRepo,
	externalConns repository.ExternalConnectionRepo,
	mcpDefs repository.McpDefinitionsRepo,
	toolExecutor ToolExecutor,
) *Service {
	return &Service{
		tracts:         tracts,
		triggers:       triggers,
		triggerPresets: triggerPresets,
		externalConns:  externalConns,
		mcpDefs:        mcpDefs,
		toolExecutor:   toolExecutor,
	}
}

// CreateTract validates def and persists a new tract owned by the caller. The returned
// warnings are non-fatal hints (e.g. a trigger.* ref not found in any linked trigger's
// payload schema) — always empty on create, since a brand-new tract has no linked triggers
// yet to check against.
func (s *Service) CreateTract(
	ctx context.Context,
	name string,
	description string,
	def domain.TractDefinition,
) (domain.Tract, []string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Tract{}, nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	err := validateShape(def)
	if err != nil {
		return domain.Tract{}, nil, err
	}

	err = s.validateTools(ctx, uc.UserUuid, def)
	if err != nil {
		return domain.Tract{}, nil, err
	}

	tract := domain.Tract{
		UserUuid:    uc.UserUuid,
		Name:        name,
		Description: description,
		Enabled:     true,
		Definition:  def,
	}

	created, err := s.tracts.Create(ctx, tract)
	if err != nil {
		return domain.Tract{}, nil, rerrors.Wrap(err, "error creating tract")
	}

	return created, checkTriggerFieldWarnings(def, nil), nil
}

func (s *Service) GetTract(ctx context.Context, id uuid.UUID) (domain.Tract, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Tract{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	result, err := s.tracts.Get(ctx, id)
	if err != nil {
		return domain.Tract{}, rerrors.Wrap(err, "error getting tract")
	}

	if !result.Valid {
		return domain.Tract{}, rerrors.Wrap(user_errors.NotFound)
	}

	if result.V.UserUuid != uc.UserUuid {
		return domain.Tract{}, rerrors.Wrap(user_errors.TractNotOwned)
	}

	return result.V, nil
}

func (s *Service) ListTracts(ctx context.Context) ([]domain.Tract, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	tracts, err := s.tracts.ListByUser(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tracts")
	}

	return tracts, nil
}

// UpdateTract overwrites name, description and the whole definition document — matching the
// "UI/agents edit whole docs" JSONB convention. Trigger.* refs are checked against the
// tract's currently linked triggers' payload schemas.
func (s *Service) UpdateTract(
	ctx context.Context,
	id uuid.UUID,
	name string,
	description string,
	def domain.TractDefinition,
) (domain.Tract, []string, error) {
	existing, err := s.GetTract(ctx, id)
	if err != nil {
		return domain.Tract{}, nil, err
	}

	err = validateShape(def)
	if err != nil {
		return domain.Tract{}, nil, err
	}

	err = s.validateTools(ctx, existing.UserUuid, def)
	if err != nil {
		return domain.Tract{}, nil, err
	}

	links, err := s.triggers.ListLinksByTract(ctx, id)
	if err != nil {
		return domain.Tract{}, nil, rerrors.Wrap(err, "error listing tract trigger links")
	}

	schemas := make([]domain.ToolSchema, len(links))
	for i, link := range links {
		schemas[i] = link.Trigger.PayloadSchema
	}

	existing.Name = name
	existing.Description = description
	existing.Definition = def

	updated, err := s.tracts.Update(ctx, existing)
	if err != nil {
		return domain.Tract{}, nil, rerrors.Wrap(err, "error updating tract")
	}

	return updated, checkTriggerFieldWarnings(def, schemas), nil
}

func (s *Service) SetTractEnabled(ctx context.Context, id uuid.UUID, enabled bool) error {
	_, err := s.GetTract(ctx, id)
	if err != nil {
		return err
	}

	err = s.tracts.SetEnabled(ctx, id, enabled)
	if err != nil {
		return rerrors.Wrap(err, "error setting tract enabled")
	}

	return nil
}

func (s *Service) DeleteTract(ctx context.Context, id uuid.UUID) error {
	_, err := s.GetTract(ctx, id)
	if err != nil {
		return err
	}

	err = s.tracts.Delete(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "error deleting tract")
	}

	return nil
}

// validateShape checks step id format/uniqueness, per-type field usage, and the template-ref
// visibility rule. Pure — no repo access.
func validateShape(def domain.TractDefinition) error {
	ids := map[string]struct{}{}

	err := walkIdsAndShape(def.Steps, ids)
	if err != nil {
		return err
	}

	_, err = checkVisibility(def.Steps, map[string]struct{}{})

	return err
}

func walkIdsAndShape(steps []domain.TractStep, ids map[string]struct{}) error {
	for _, step := range steps {
		err := validateStepId(step.Id, ids)
		if err != nil {
			return err
		}

		err = validateStepShape(step)
		if err != nil {
			return err
		}

		switch step.Type {
		case stepTypeCondition:
			err = walkIdsAndShape(step.Then, ids)
			if err != nil {
				return err
			}

			err = walkIdsAndShape(step.Else, ids)
			if err != nil {
				return err
			}
		case stepTypeParallel, stepTypeGroup:
			err = walkIdsAndShape(step.Steps, ids)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateStepId(id string, ids map[string]struct{}) error {
	if id == reservedStepId || !stepIdPattern.MatchString(id) {
		return rerrors.Wrap(user_errors.TractStepIdInvalid, id)
	}

	_, exists := ids[id]
	if exists {
		return rerrors.Wrap(user_errors.TractStepIdDuplicate, id)
	}

	ids[id] = struct{}{}

	return nil
}

func validateStepShape(step domain.TractStep) error {
	switch step.Type {
	case stepTypeAction:
		if len(step.Conditions) > 0 || len(step.Then) > 0 || len(step.Else) > 0 || len(step.Steps) > 0 {
			return rerrors.Wrap(user_errors.TractConditionsOnNonBranch, step.Id)
		}

		if step.Mcp == "" || step.Tool == "" {
			return rerrors.Wrap(user_errors.TractToolNotFound, step.Id)
		}

		return nil

	case stepTypeCondition:
		if len(step.Params) > 0 || step.Mcp != "" || step.Tool != "" || step.ConnectionUuid != uuid.Nil ||
			len(step.Steps) > 0 {
			return rerrors.Wrap(user_errors.TractParamsOnNonAction, step.Id)
		}

		if len(step.Conditions) == 0 {
			return rerrors.Wrap(
				user_errors.TractConditionsOnNonBranch,
				fmt.Sprintf("%s: condition step requires at least one condition", step.Id),
			)
		}

		return nil

	case stepTypeParallel, stepTypeGroup:
		if len(step.Params) > 0 || step.Mcp != "" || step.Tool != "" || step.ConnectionUuid != uuid.Nil ||
			len(step.Conditions) > 0 || len(step.Then) > 0 || len(step.Else) > 0 {
			return rerrors.Wrap(user_errors.TractParamsOnNonAction, step.Id)
		}

		if len(step.Steps) == 0 {
			return rerrors.Wrap(
				user_errors.TractUnknownStepType,
				fmt.Sprintf("%s: %s step requires at least one child step", step.Id, step.Type),
			)
		}

		return nil

	default:
		return rerrors.Wrap(user_errors.TractUnknownStepType, step.Type)
	}
}

// checkVisibility walks steps in execution order, verifying every template ref (from Params
// and Conditions) resolves to "trigger" or a step that is definitely executed before it:
// a preceding sibling at any ancestor level. Parallel lanes never see each other while
// running, but all of a parallel/group step's descendant ids fan in for steps that follow it.
// A condition step's own id becomes visible after it, but ids declared inside its Then/Else
// branches never escape the branch. Returns the ids this step list makes visible to whatever
// follows it (used by the parallel/group cases to fan those ids upward).
func checkVisibility(steps []domain.TractStep, visible map[string]struct{}) (map[string]struct{}, error) {
	local := make(map[string]struct{}, len(visible))
	for id := range visible {
		local[id] = struct{}{}
	}

	produced := map[string]struct{}{}

	for _, step := range steps {
		err := checkStepRefs(step, local)
		if err != nil {
			return nil, err
		}

		switch step.Type {
		case stepTypeAction:
			local[step.Id] = struct{}{}
			produced[step.Id] = struct{}{}

		case stepTypeCondition:
			_, err = checkVisibility(step.Then, local)
			if err != nil {
				return nil, err
			}

			_, err = checkVisibility(step.Else, local)
			if err != nil {
				return nil, err
			}

			local[step.Id] = struct{}{}
			produced[step.Id] = struct{}{}

		case stepTypeParallel:
			fannedIn := map[string]struct{}{}

			for _, lane := range step.Steps {
				laneSteps := []domain.TractStep{lane}

				laneProduced, err := checkVisibility(laneSteps, local)
				if err != nil {
					return nil, err
				}

				for id := range laneProduced {
					fannedIn[id] = struct{}{}
				}
			}

			for id := range fannedIn {
				local[id] = struct{}{}
				produced[id] = struct{}{}
			}

		case stepTypeGroup:
			groupProduced, err := checkVisibility(step.Steps, local)
			if err != nil {
				return nil, err
			}

			for id := range groupProduced {
				local[id] = struct{}{}
				produced[id] = struct{}{}
			}
		}
	}

	return produced, nil
}

func checkStepRefs(step domain.TractStep, visible map[string]struct{}) error {
	var sources []string
	for _, v := range step.Params {
		sources = append(sources, v)
	}

	for _, c := range step.Conditions {
		sources = append(sources, c.Left, c.Right)
	}

	for _, source := range sources {
		refs, err := extractRefs(source)
		if err != nil {
			return err
		}

		for _, ref := range refs {
			base, _, err := splitRef(ref)
			if err != nil {
				return err
			}

			if base == reservedStepId {
				continue
			}

			_, ok := visible[base]
			if !ok {
				return rerrors.Wrap(
					user_errors.TractInvalidTemplateRef,
					fmt.Sprintf("step %q references %q", step.Id, ref),
				)
			}
		}
	}

	return nil
}

// validateTools checks every action step's tool reference: a builtin name known to
// toolExecutor, or an (mcp, tool) pair present in mcp_tools; connection_id is required for
// MoM tools (and must belong to ownerUuid), forbidden for builtins.
func (s *Service) validateTools(ctx context.Context, ownerUuid uuid.UUID, def domain.TractDefinition) error {
	return walkActions(def.Steps, func(step domain.TractStep) error {
		return s.validateActionTool(ctx, ownerUuid, step)
	})
}

func walkActions(steps []domain.TractStep, fn func(domain.TractStep) error) error {
	for _, step := range steps {
		switch step.Type {
		case stepTypeAction:
			err := fn(step)
			if err != nil {
				return err
			}
		case stepTypeCondition:
			err := walkActions(step.Then, fn)
			if err != nil {
				return err
			}

			err = walkActions(step.Else, fn)
			if err != nil {
				return err
			}
		case stepTypeParallel, stepTypeGroup:
			err := walkActions(step.Steps, fn)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) validateActionTool(ctx context.Context, ownerUuid uuid.UUID, step domain.TractStep) error {
	if step.Mcp == builtinMcpName {
		if !s.toolExecutor.IsBuiltinTool(step.Tool) {
			return rerrors.Wrap(user_errors.TractToolNotFound, step.Id)
		}

		if step.ConnectionUuid != uuid.Nil {
			return rerrors.Wrap(user_errors.TractConnectionForbidden, step.Id)
		}

		return nil
	}

	tool, err := s.mcpDefs.GetTool(ctx, step.Mcp, step.Tool)
	if err != nil {
		return rerrors.Wrap(err, "error checking tract action tool")
	}

	if !tool.Valid {
		return rerrors.Wrap(user_errors.TractToolNotFound, step.Id)
	}

	if step.ConnectionUuid == uuid.Nil {
		return rerrors.Wrap(user_errors.TractConnectionRequired, step.Id)
	}

	conn, err := s.externalConns.GetByID(ctx, step.ConnectionUuid)
	if err != nil {
		return rerrors.Wrap(err, "error getting external connection")
	}

	if conn.UserUuid != ownerUuid {
		return rerrors.Wrap(user_errors.TractConnectionNotOwned, step.Id)
	}

	return nil
}

// checkTriggerFieldWarnings returns a warning for every distinct top-level `trigger.<field>`
// ref in def that isn't declared by any of schemas — a best-effort hint (matches the "hint
// UI" language in the plan), not a full type check. No schemas known (e.g. a tract with no
// linked triggers yet) means nothing can be verified, so no warnings are produced.
func checkTriggerFieldWarnings(def domain.TractDefinition, schemas []domain.ToolSchema) []string {
	if len(schemas) == 0 {
		return nil
	}

	fields := collectTriggerFields(def.Steps)
	if len(fields) == 0 {
		return nil
	}

	seen := map[string]struct{}{}

	var warnings []string

	for _, field := range fields {
		_, already := seen[field]
		if already {
			continue
		}

		seen[field] = struct{}{}

		if !anySchemaDeclaresField(schemas, field) {
			warnings = append(
				warnings,
				fmt.Sprintf(
					"template references trigger.%s, which is not declared by any linked trigger's payload schema",
					field,
				),
			)
		}
	}

	return warnings
}

func collectTriggerFields(steps []domain.TractStep) []string {
	var fields []string

	for _, step := range steps {
		var sources []string
		for _, v := range step.Params {
			sources = append(sources, v)
		}

		for _, c := range step.Conditions {
			sources = append(sources, c.Left, c.Right)
		}

		for _, source := range sources {
			refs, err := extractRefs(source)
			if err != nil {
				continue
			}

			for _, ref := range refs {
				base, segments, err := splitRef(ref)
				if err != nil || base != reservedStepId || len(segments) == 0 || segments[0].isIndex {
					continue
				}

				fields = append(fields, segments[0].field)
			}
		}

		fields = append(fields, collectTriggerFields(step.Then)...)
		fields = append(fields, collectTriggerFields(step.Else)...)
		fields = append(fields, collectTriggerFields(step.Steps)...)
	}

	return fields
}

func anySchemaDeclaresField(schemas []domain.ToolSchema, field string) bool {
	for _, schema := range schemas {
		_, ok := schema.Properties[field]
		if ok {
			return true
		}
	}

	return false
}
