package tracts_api

import (
	"encoding/json"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

const timeFormat = "2006-01-02T15:04:05Z"

// tractToProto renders a bare TractItem — no linked-trigger summaries or last-run info.
// ListTracts overlays those via tractToProtoWithSummary.
func tractToProto(t domain.Tract) *pb.TractItem {
	item := &pb.TractItem{
		Uuid:        t.Uuid.String(),
		Name:        t.Name,
		Description: t.Description,
		Enabled:     t.Enabled,
		Definition:  definitionToProto(t.Definition),
		CreatedAt:   t.CreatedAt.UTC().Format(timeFormat),
		UpdatedAt:   t.UpdatedAt.UTC().Format(timeFormat),
	}

	return item
}

// tractToProtoWithSummary adds the linked-trigger summaries and last-run badge ListTracts
// shows on each card. lastRun is nil when the tract has never run.
func tractToProtoWithSummary(
	t domain.Tract,
	links []repository.TractTriggerLink,
	lastRun *domain.TractRun,
) *pb.TractItem {
	item := tractToProto(t)

	item.Triggers = make([]*pb.TractTriggerSummary, len(links))
	for i, link := range links {
		item.Triggers[i] = triggerSummaryToProto(link.Trigger)
	}

	if lastRun != nil {
		item.LastRun = &pb.TractLastRun{
			Status: string(lastRun.Status),
			At:     lastRun.CreatedAt.UTC().Format(timeFormat),
		}
	}

	return item
}

func triggerSummaryToProto(t domain.Trigger) *pb.TractTriggerSummary {
	summary := &pb.TractTriggerSummary{
		Uuid:   t.Uuid.String(),
		Name:   t.Name,
		Kind:   t.Kind,
		Source: t.Source,
	}

	return summary
}

func runToProto(r domain.TractRun) *pb.TractRunItem {
	triggerUuid := ""
	if r.TriggerUuid != uuid.Nil {
		triggerUuid = r.TriggerUuid.String()
	}

	item := &pb.TractRunItem{
		Uuid:           r.Uuid.String(),
		TractUuid:      r.TractUuid.String(),
		TriggerUuid:    triggerUuid,
		Status:         string(r.Status),
		StartedBy:      r.StartedBy,
		TriggerPayload: string(r.TriggerPayload),
		Error:          r.Error,
		CreatedAt:      r.CreatedAt.UTC().Format(timeFormat),
		UpdatedAt:      r.UpdatedAt.UTC().Format(timeFormat),
	}

	return item
}

func runsToProto(runs []domain.TractRun) []*pb.TractRunItem {
	items := make([]*pb.TractRunItem, len(runs))
	for i, r := range runs {
		items[i] = runToProto(r)
	}

	return items
}

func runStepToProto(s domain.TractRunStep) *pb.TractRunStepItem {
	finishedAt := ""
	if !s.FinishedAt.IsZero() {
		finishedAt = s.FinishedAt.UTC().Format(timeFormat)
	}

	item := &pb.TractRunStepItem{
		StepId:     s.StepId,
		StepName:   s.StepName,
		StepType:   s.StepType,
		Status:     string(s.Status),
		Input:      string(s.Input),
		Output:     string(s.Output),
		Error:      s.Error,
		StartedAt:  s.StartedAt.UTC().Format(timeFormat),
		FinishedAt: finishedAt,
	}

	return item
}

func runStepsToProto(steps []domain.TractRunStep) []*pb.TractRunStepItem {
	items := make([]*pb.TractRunStepItem, len(steps))
	for i, s := range steps {
		items[i] = runStepToProto(s)
	}

	return items
}

func toolRefToProto(ref domain.McpToolRef) *pb.TractToolItem {
	inputSchema := domain.ToolSchema{
		Properties: ref.Tool.ApiDescription.Properties,
		Required:   ref.Tool.ApiDescription.Required,
	}

	item := &pb.TractToolItem{
		Mcp:          ref.McpName,
		Tool:         ref.Tool.ApiDescription.Name,
		Description:  ref.Tool.ApiDescription.Description,
		InputSchema:  schemaToJSON(inputSchema),
		OutputSchema: schemaToJSON(ref.Tool.OutputSchema),
	}

	return item
}

func toolRefsToProto(refs []domain.McpToolRef) []*pb.TractToolItem {
	items := make([]*pb.TractToolItem, len(refs))
	for i, ref := range refs {
		items[i] = toolRefToProto(ref)
	}

	return items
}

func triggerSourceToProto(p domain.TriggerPreset) *pb.TriggerSourceItem {
	item := &pb.TriggerSourceItem{
		Key:           p.Key,
		Description:   p.Description,
		PayloadSchema: schemaToJSON(p.PayloadSchema),
		Category:      p.Category,
		Label:         p.Label,
		Provider:      p.Provider,
	}

	return item
}

func triggerSourcesToProto(presets []domain.TriggerPreset) []*pb.TriggerSourceItem {
	items := make([]*pb.TriggerSourceItem, len(presets))
	for i, p := range presets {
		items[i] = triggerSourceToProto(p)
	}

	return items
}

func triggerToProto(t domain.Trigger) *pb.TriggerItem {
	item := &pb.TriggerItem{
		Uuid:          t.Uuid.String(),
		Name:          t.Name,
		Kind:          t.Kind,
		Source:        t.Source,
		Config:        string(t.Config),
		PayloadSchema: schemaToJSON(t.PayloadSchema),
		TriggerUuid:   t.TriggerUuid.String(),
		Enabled:       t.Enabled,
		CreatedAt:     t.CreatedAt.UTC().Format(timeFormat),
	}

	return item
}

func triggersToProto(triggers []domain.Trigger) []*pb.TriggerItem {
	items := make([]*pb.TriggerItem, len(triggers))
	for i, t := range triggers {
		items[i] = triggerToProto(t)
	}

	return items
}

// -- JSON <-> domain.ToolSchema/ToolProperty --
//
// domain.ToolSchema/ToolProperty carry no json tags (kept persistence-agnostic — see
// internal/repository/pg/repos/triggers/triggers.go) so every layer that needs a wire
// representation mirrors them locally with tagged row types. tracts_api's wire shape is
// {"properties": {...}, "required": [...]}, matching how the DB itself stores payload_schema
// and mcp_tools' input/output schemas.

type toolSchemaRow struct {
	Properties map[string]toolPropertyRow `json:"properties"`
	Required   []string                   `json:"required"`
}

type toolPropertyRow struct {
	Type        string                     `json:"type"`
	Description string                     `json:"description,omitempty"`
	Enum        []string                   `json:"enum,omitempty"`
	Properties  map[string]toolPropertyRow `json:"properties,omitempty"`
	Items       *toolPropertyRow           `json:"items,omitempty"`
	Required    []string                   `json:"required,omitempty"`
}

func schemaToJSON(schema domain.ToolSchema) string {
	row := toolSchemaRow{
		Properties: make(map[string]toolPropertyRow, len(schema.Properties)),
		Required:   schema.Required,
	}
	for k, v := range schema.Properties {
		row.Properties[k] = toolPropertyRowFromDomain(v)
	}

	data, err := json.Marshal(row)
	if err != nil {
		return "{}"
	}

	return string(data)
}

func schemaFromJSON(raw string) (domain.ToolSchema, error) {
	if raw == "" {
		return domain.ToolSchema{}, nil
	}

	var row toolSchemaRow

	err := json.Unmarshal([]byte(raw), &row)
	if err != nil {
		return domain.ToolSchema{}, rerrors.Wrap(
			user_errors.TractRequestFieldInvalidJSON,
			"error unmarshaling payload_schema",
		)
	}

	schema := domain.ToolSchema{
		Properties: make(map[string]domain.ToolProperty, len(row.Properties)),
		Required:   row.Required,
	}
	for k, v := range row.Properties {
		schema.Properties[k] = toolPropertyRowToDomain(v)
	}

	return schema, nil
}

func toolPropertyRowFromDomain(p domain.ToolProperty) toolPropertyRow {
	row := toolPropertyRow{
		Type:        p.Type,
		Description: p.Description,
		Enum:        p.Enum,
		Required:    p.Required,
	}

	if len(p.Properties) > 0 {
		row.Properties = make(map[string]toolPropertyRow, len(p.Properties))
		for k, v := range p.Properties {
			row.Properties[k] = toolPropertyRowFromDomain(v)
		}
	}

	if p.Items != nil {
		items := toolPropertyRowFromDomain(*p.Items)
		row.Items = &items
	}

	return row
}

func toolPropertyRowToDomain(p toolPropertyRow) domain.ToolProperty {
	prop := domain.ToolProperty{
		Type:        p.Type,
		Description: p.Description,
		Enum:        p.Enum,
		Required:    p.Required,
	}

	if len(p.Properties) > 0 {
		prop.Properties = make(map[string]domain.ToolProperty, len(p.Properties))
		for k, v := range p.Properties {
			prop.Properties[k] = toolPropertyRowToDomain(v)
		}
	}

	if p.Items != nil {
		items := toolPropertyRowToDomain(*p.Items)
		prop.Items = &items
	}

	return prop
}

// -- pb.TractDefinition/TractStep (oneof wire shape) <-> domain.TractDefinition/TractStep (flat) --
//
// The wire shape uses a oneof per step kind so invalid field combinations (e.g. action fields
// set alongside conditions) are unrepresentable on the wire. domain.TractStep stays flat with a
// Type string discriminant, matching what internal/service/v1/tract's validation is keyed off
// of — these functions are the only place that bridges the two shapes.

func conditionToProto(c domain.TractCondition) *pb.TractCondition {
	return &pb.TractCondition{Left: c.Left, Op: c.Op, Right: c.Right}
}

func conditionFromProto(c *pb.TractCondition) domain.TractCondition {
	return domain.TractCondition{Left: c.Left, Op: c.Op, Right: c.Right}
}

func conditionsToProto(cs []domain.TractCondition) []*pb.TractCondition {
	items := make([]*pb.TractCondition, len(cs))
	for i, c := range cs {
		items[i] = conditionToProto(c)
	}

	return items
}

func conditionsFromProto(cs []*pb.TractCondition) []domain.TractCondition {
	items := make([]domain.TractCondition, len(cs))
	for i, c := range cs {
		items[i] = conditionFromProto(c)
	}

	return items
}

func stepToProto(s domain.TractStep) *pb.TractStep {
	step := &pb.TractStep{
		Id:          s.Id,
		Name:        s.Name,
		Description: s.Description,
	}

	switch s.Type {
	case "action":
		connectionUuid := ""
		if s.ConnectionUuid != uuid.Nil {
			connectionUuid = s.ConnectionUuid.String()
		}

		step.Kind = &pb.TractStep_Action{Action: &pb.ActionStep{
			Mcp:            s.Mcp,
			Tool:           s.Tool,
			ConnectionUuid: connectionUuid,
			Params:         s.Params,
		}}
	case "condition":
		step.Kind = &pb.TractStep_Condition{Condition: &pb.ConditionStep{
			Conditions: conditionsToProto(s.Conditions),
			Then:       stepsToProto(s.Then),
			Else:       stepsToProto(s.Else),
		}}
	case "parallel":
		step.Kind = &pb.TractStep_Parallel{Parallel: &pb.ParallelStep{Steps: stepsToProto(s.Steps)}}
	case "group":
		step.Kind = &pb.TractStep_Group{Group: &pb.GroupStep{Steps: stepsToProto(s.Steps)}}
	}

	return step
}

func stepsToProto(steps []domain.TractStep) []*pb.TractStep {
	items := make([]*pb.TractStep, len(steps))
	for i, s := range steps {
		items[i] = stepToProto(s)
	}

	return items
}

func stepFromProto(s *pb.TractStep) (domain.TractStep, error) {
	step := domain.TractStep{
		Id:          s.Id,
		Name:        s.Name,
		Description: s.Description,
	}

	switch kind := s.Kind.(type) {
	case *pb.TractStep_Action:
		step.Type = "action"
		step.Mcp = kind.Action.Mcp
		step.Tool = kind.Action.Tool
		step.Params = kind.Action.Params

		if kind.Action.ConnectionUuid != "" {
			id, err := uuid.Parse(kind.Action.ConnectionUuid)
			if err != nil {
				return domain.TractStep{}, rerrors.Wrap(
					user_errors.TractRequestFieldInvalidJSON,
					"error parsing step connection_uuid",
				)
			}

			step.ConnectionUuid = id
		}
	case *pb.TractStep_Condition:
		step.Type = "condition"
		step.Conditions = conditionsFromProto(kind.Condition.Conditions)

		then, err := stepsFromProto(kind.Condition.Then)
		if err != nil {
			return domain.TractStep{}, err
		}

		els, err := stepsFromProto(kind.Condition.Else)
		if err != nil {
			return domain.TractStep{}, err
		}

		step.Then = then
		step.Else = els
	case *pb.TractStep_Parallel:
		step.Type = "parallel"

		steps, err := stepsFromProto(kind.Parallel.Steps)
		if err != nil {
			return domain.TractStep{}, err
		}

		step.Steps = steps
	case *pb.TractStep_Group:
		step.Type = "group"

		steps, err := stepsFromProto(kind.Group.Steps)
		if err != nil {
			return domain.TractStep{}, err
		}

		step.Steps = steps
	default:
		return domain.TractStep{}, rerrors.Wrap(
			user_errors.TractRequestFieldInvalidJSON,
			"tract step is missing a kind",
		)
	}

	return step, nil
}

func stepsFromProto(steps []*pb.TractStep) ([]domain.TractStep, error) {
	items := make([]domain.TractStep, len(steps))
	for i, s := range steps {
		step, err := stepFromProto(s)
		if err != nil {
			return nil, err
		}

		items[i] = step
	}

	return items, nil
}

func definitionToProto(def domain.TractDefinition) *pb.TractDefinition {
	return &pb.TractDefinition{Steps: stepsToProto(def.Steps)}
}

func definitionFromProto(def *pb.TractDefinition) (domain.TractDefinition, error) {
	if def == nil {
		return domain.TractDefinition{}, nil
	}

	steps, err := stepsFromProto(def.Steps)
	if err != nil {
		return domain.TractDefinition{}, err
	}

	return domain.TractDefinition{Steps: steps}, nil
}

func filtersFromJSON(raw string) ([]domain.TractCondition, error) {
	if raw == "" {
		return nil, nil
	}

	var filters []domain.TractCondition

	err := json.Unmarshal([]byte(raw), &filters)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.TractRequestFieldInvalidJSON, "error unmarshaling filters")
	}

	return filters, nil
}
