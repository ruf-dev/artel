package tracts_api

import (
	"encoding/json"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

const timeFormat = "2006-01-02T15:04:05Z"

// tractToProto renders a bare TractItem — no linked-trigger summaries or last-run info.
// ListTracts overlays those via tractToProtoWithSummary.
func tractToProto(t domain.Tract) (*pb.TractItem, error) {
	defJSON, err := json.Marshal(t.Definition)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshaling tract definition")
	}

	item := &pb.TractItem{
		Uuid:        t.Uuid.String(),
		Name:        t.Name,
		Description: t.Description,
		Enabled:     t.Enabled,
		Definition:  string(defJSON),
		CreatedAt:   t.CreatedAt.UTC().Format(timeFormat),
		UpdatedAt:   t.UpdatedAt.UTC().Format(timeFormat),
	}
	return item, nil
}

// tractToProtoWithSummary adds the linked-trigger summaries and last-run badge ListTracts
// shows on each card. lastRun is nil when the tract has never run.
func tractToProtoWithSummary(t domain.Tract, links []repository.TractTriggerLink, lastRun *domain.TractRun) (*pb.TractItem, error) {
	item, err := tractToProto(t)
	if err != nil {
		return nil, err
	}

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

	return item, nil
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
	inputSchema := domain.ToolSchema{Properties: ref.Tool.ApiDescription.Properties, Required: ref.Tool.ApiDescription.Required}

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

func triggerSourceToProto(p domain.TriggerSourcePreset) *pb.TriggerSourceItem {
	item := &pb.TriggerSourceItem{
		Key:           p.Key,
		Description:   p.Description,
		PayloadSchema: schemaToJSON(p.PayloadSchema),
	}
	return item
}

func triggerSourcesToProto(presets []domain.TriggerSourcePreset) []*pb.TriggerSourceItem {
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
	row := toolSchemaRow{Properties: make(map[string]toolPropertyRow, len(schema.Properties)), Required: schema.Required}
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
		return domain.ToolSchema{}, rerrors.Wrap(user_errors.TractRequestFieldInvalidJSON, "error unmarshaling payload_schema")
	}

	schema := domain.ToolSchema{Properties: make(map[string]domain.ToolProperty, len(row.Properties)), Required: row.Required}
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

// -- JSON <-> domain.TractDefinition / []domain.TractCondition --
//
// Unlike ToolSchema/ToolProperty, domain.TractDefinition and domain.TractCondition already
// carry json tags (see internal/domain/tract.go) — the engine/validation layer marshals them
// directly — so these just wrap json.(Un)marshal with the transport's error convention.

func definitionFromJSON(raw string) (domain.TractDefinition, error) {
	var def domain.TractDefinition
	err := json.Unmarshal([]byte(raw), &def)
	if err != nil {
		return domain.TractDefinition{}, rerrors.Wrap(user_errors.TractRequestFieldInvalidJSON, "error unmarshaling definition")
	}
	return def, nil
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
