// Package triggerpresets is the pure-DB read layer for the trigger_presets catalog table (see
// migration 038_trigger_presets_and_links.sql) — the webhook trigger presets shown in the "Add
// trigger" picker (gitlab_push, generic, ...), replacing the hardcoded Go slice that used to live
// in internal/service/v1/tract.
package triggerpresets

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) List(ctx context.Context) ([]domain.TriggerPreset, error) {
	rows, err := r.q.ListTriggerPresets(ctx)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing trigger presets")
	}

	presets := make([]domain.TriggerPreset, len(rows))

	for i, row := range rows {
		preset, convErr := presetToDomain(row)
		if convErr != nil {
			return nil, convErr
		}

		presets[i] = preset
	}

	return presets, nil
}

func (r *Repo) GetByKey(ctx context.Context, key string) (sql.Null[domain.TriggerPreset], error) {
	row, err := r.q.GetTriggerPresetByKey(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.TriggerPreset]{}, nil
		}

		return sql.Null[domain.TriggerPreset]{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting trigger preset")
	}

	preset, err := presetToDomain(row)
	if err != nil {
		return sql.Null[domain.TriggerPreset]{}, err
	}

	result := sql.Null[domain.TriggerPreset]{V: preset, Valid: true}

	return result, nil
}

func presetToDomain(row artel_q.TriggerPreset) (domain.TriggerPreset, error) {
	schema, err := unmarshalToolSchema(row.PayloadSchema)
	if err != nil {
		return domain.TriggerPreset{}, err
	}

	matchers, err := unmarshalMatchers(row.DefaultMatchers)
	if err != nil {
		return domain.TriggerPreset{}, err
	}

	preset := domain.TriggerPreset{
		Key:             row.Key,
		Category:        row.Category,
		Label:           row.Label,
		Description:     row.Description,
		Provider:        row.Provider.String,
		PayloadSchema:   schema,
		DefaultMatchers: matchers,
	}

	return preset, nil
}

// tool schema (de)serialization — mirrors the private row types in
// internal/repository/pg/repos/triggers/triggers.go so domain.ToolSchema/ToolProperty stay free
// of persistence JSON tags. Duplicated rather than shared per that file's own precedent (see its
// comment referencing the mcpdefinitions repo).

func unmarshalToolSchema(raw json.RawMessage) (domain.ToolSchema, error) {
	if len(raw) == 0 {
		return domain.ToolSchema{}, nil
	}

	var row toolSchemaRow

	err := json.Unmarshal(raw, &row)
	if err != nil {
		return domain.ToolSchema{}, rerrors.Wrap(err, "error unmarshaling trigger preset payload schema")
	}

	schema := domain.ToolSchema{
		Properties: make(map[string]domain.ToolProperty, len(row.Properties)),
		Required:   row.Required,
		IsArray:    row.IsArray,
	}
	for k, v := range row.Properties {
		schema.Properties[k] = toolPropertyRowToDomain(v)
	}

	if row.Items != nil {
		items := toolPropertyRowToDomain(*row.Items)
		schema.Items = &items
	}

	return schema, nil
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

type toolSchemaRow struct {
	Properties map[string]toolPropertyRow `json:"properties"`
	Required   []string                   `json:"required"`
	IsArray    bool                       `json:"isArray,omitempty"`
	Items      *toolPropertyRow           `json:"items,omitempty"`
}

type toolPropertyRow struct {
	Type        string                     `json:"type"`
	Description string                     `json:"description,omitempty"`
	Enum        []string                   `json:"enum,omitempty"`
	Properties  map[string]toolPropertyRow `json:"properties,omitempty"`
	Items       *toolPropertyRow           `json:"items,omitempty"`
	Required    []string                   `json:"required,omitempty"`
}

// matchers (de)serialization — mirrors internal/repository/pg/repos/triggers/triggers.go's
// triggerMatchersRow, duplicated for the same reason as the schema row types above.

func unmarshalMatchers(raw json.RawMessage) (domain.TriggerMatchers, error) {
	if len(raw) == 0 {
		return domain.TriggerMatchers{}, nil
	}

	var row triggerMatchersRow

	err := json.Unmarshal(raw, &row)
	if err != nil {
		return domain.TriggerMatchers{}, rerrors.Wrap(err, "error unmarshaling trigger preset default matchers")
	}

	matchers := domain.TriggerMatchers{
		CheckHeaders: make([]domain.HeaderMatcher, len(row.CheckHeaders)),
		CheckBody:    make([]domain.BodyMatcher, len(row.CheckBody)),
	}
	for i, m := range row.CheckHeaders {
		matchers.CheckHeaders[i] = domain.HeaderMatcher{Header: m.Header, Equals: m.Equals}
	}
	for i, m := range row.CheckBody {
		matchers.CheckBody[i] = domain.BodyMatcher{Path: m.Path, Equals: m.Equals}
	}

	return matchers, nil
}

type triggerMatchersRow struct {
	CheckHeaders []headerMatcherRow `json:"check_headers,omitempty"`
	CheckBody    []bodyMatcherRow   `json:"check_body,omitempty"`
}

type headerMatcherRow struct {
	Header string `json:"header"`
	Equals string `json:"equals"`
}

type bodyMatcherRow struct {
	Path   string `json:"path"`
	Equals string `json:"equals"`
}
