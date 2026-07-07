package triggers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"go.redsock.ru/rerrors"
)

// Repo is the pure-DB layer for standalone triggers and their tract links (see migration 033:
// triggers, tract_trigger_links).
type Repo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) Create(ctx context.Context, trigger domain.Trigger) (domain.Trigger, error) {
	config := trigger.Config
	if config == nil {
		config = json.RawMessage("{}")
	}

	schemaJSON, err := marshalToolSchema(trigger.PayloadSchema)
	if err != nil {
		return domain.Trigger{}, err
	}

	matchersJSON, err := marshalMatchers(trigger.Matchers)
	if err != nil {
		return domain.Trigger{}, err
	}

	params := artel_q.InsertTriggerParams{
		UserID:        trigger.UserUuid,
		Name:          trigger.Name,
		Kind:          trigger.Kind,
		Source:        trigger.Source,
		Config:        config,
		PayloadSchema: schemaJSON,
		SecretHash:    trigger.SecretHash,
		TokenSuffix:   sql.NullString{String: trigger.TokenSuffix, Valid: trigger.TokenSuffix != ""},
		Matchers:      matchersJSON,
		Enabled:       trigger.Enabled,
	}

	row, err := r.q.InsertTrigger(ctx, params)
	if err != nil {
		return domain.Trigger{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error inserting trigger")
	}

	return triggerToDomain(row)
}

func (r *Repo) Get(ctx context.Context, id uuid.UUID) (sql.Null[domain.Trigger], error) {
	row, err := r.q.GetTrigger(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.Trigger]{}, nil
		}

		return sql.Null[domain.Trigger]{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting trigger")
	}

	trigger, err := triggerToDomain(row)
	if err != nil {
		return sql.Null[domain.Trigger]{}, err
	}

	result := sql.Null[domain.Trigger]{V: trigger, Valid: true}

	return result, nil
}

// GetByTriggerUuid looks up a trigger by its rotatable webhook routing id — the inbound
// webhook handler only knows the routing id embedded in the fired URL.
func (r *Repo) GetByTriggerUuid(ctx context.Context, triggerUuid uuid.UUID) (sql.Null[domain.Trigger], error) {
	row, err := r.q.GetTriggerByTriggerUuid(ctx, triggerUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.Trigger]{}, nil
		}

		return sql.Null[domain.Trigger]{}, rerrors.Wrap(
			pg_err.UnwrapPgErr(err),
			"error getting trigger by trigger uuid",
		)
	}

	trigger, err := triggerToDomain(row)
	if err != nil {
		return sql.Null[domain.Trigger]{}, err
	}

	result := sql.Null[domain.Trigger]{V: trigger, Valid: true}

	return result, nil
}

func (r *Repo) ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.Trigger, error) {
	rows, err := r.q.ListTriggersByUser(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing triggers")
	}

	triggers := make([]domain.Trigger, len(rows))

	for i, row := range rows {
		trigger, convErr := triggerToDomain(row)
		if convErr != nil {
			return nil, convErr
		}

		triggers[i] = trigger
	}

	return triggers, nil
}

func (r *Repo) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) error {
	params := artel_q.SetTriggerEnabledParams{
		ID:      id,
		Enabled: enabled,
	}

	err := r.q.SetTriggerEnabled(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error setting trigger enabled")
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteTrigger(ctx, id)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error deleting trigger")
	}

	return nil
}

// RotateSecret overwrites trigger_uuid and secret_hash in place, keyed by the trigger's stable
// primary key id — invalidates the trigger's current webhook URL/token without touching
// anything else about the trigger or its tract links.
func (r *Repo) RotateSecret(
	ctx context.Context,
	id uuid.UUID,
	newTriggerUuid uuid.UUID,
	secretHash []byte,
	tokenSuffix string,
) (domain.Trigger, error) {
	params := artel_q.RotateTriggerSecretParams{
		ID:          id,
		TriggerUuid: newTriggerUuid,
		SecretHash:  secretHash,
		TokenSuffix: sql.NullString{String: tokenSuffix, Valid: tokenSuffix != ""},
	}

	row, err := r.q.RotateTriggerSecret(ctx, params)
	if err != nil {
		return domain.Trigger{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error rotating trigger secret")
	}

	return triggerToDomain(row)
}

// Link upserts the (tractUuid, triggerUuid) link's filters — ON CONFLICT DO UPDATE, so
// re-linking an already-linked pair just replaces its filters.
func (r *Repo) Link(ctx context.Context, link domain.TriggerTractLink) error {
	filtersJSON, err := marshalFilters(link.Filters)
	if err != nil {
		return err
	}

	params := artel_q.LinkTriggerToTractParams{
		TractID:   link.TractUuid,
		TriggerID: link.TriggerUuid,
		Filters:   filtersJSON,
	}

	err = r.q.LinkTriggerToTract(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error linking trigger to tract")
	}

	return nil
}

func (r *Repo) Unlink(ctx context.Context, triggerUuid uuid.UUID, tractUuid uuid.UUID) error {
	params := artel_q.UnlinkTriggerFromTractParams{
		TriggerID: triggerUuid,
		TractID:   tractUuid,
	}

	err := r.q.UnlinkTriggerFromTract(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error unlinking trigger from tract")
	}

	return nil
}

// ListLinksByTract returns tractUuid's linked triggers with Trigger populated (Tract left
// zero) — the tract editor's "wired up triggers" view.
func (r *Repo) ListLinksByTract(ctx context.Context, tractUuid uuid.UUID) ([]repository.TractTriggerLink, error) {
	rows, err := r.q.ListTriggerLinksByTract(ctx, tractUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing tract trigger links")
	}

	links := make([]repository.TractTriggerLink, len(rows))

	for i, row := range rows {
		schema, convErr := unmarshalToolSchema(row.PayloadSchema)
		if convErr != nil {
			return nil, convErr
		}

		matchers, convErr := unmarshalMatchers(row.Matchers)
		if convErr != nil {
			return nil, convErr
		}

		trigger := domain.Trigger{
			Uuid:          row.ID,
			TriggerUuid:   row.TriggerUuid,
			UserUuid:      row.UserID,
			Name:          row.Name,
			Kind:          row.Kind,
			Source:        row.Source,
			Config:        row.Config,
			PayloadSchema: schema,
			SecretHash:    row.SecretHash,
			TokenSuffix:   row.TokenSuffix.String,
			Matchers:      matchers,
			Enabled:       row.Enabled,
			CreatedAt:     row.CreatedAt,
		}

		filters, convErr := unmarshalFilters(row.Filters)
		if convErr != nil {
			return nil, convErr
		}

		links[i] = repository.TractTriggerLink{
			TractUuid:   tractUuid,
			TriggerUuid: trigger.Uuid,
			Trigger:     trigger,
			Filters:     filters,
		}
	}

	return links, nil
}

// LinkToProvider attaches triggerId to a shared provider connection (see trigger_provider_links)
// instead of it minting its own trigger_uuid/secret_hash webhook URL.
func (r *Repo) LinkToProvider(ctx context.Context, triggerId uuid.UUID, externalConnectionId uuid.UUID) error {
	params := artel_q.InsertTriggerProviderLinkParams{
		TriggerID:            triggerId,
		ExternalConnectionID: externalConnectionId,
	}

	err := r.q.InsertTriggerProviderLink(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error linking trigger to provider connection")
	}

	return nil
}

// ListByExternalConnection returns every trigger (Matchers populated) linked to one shared
// provider connection — the gitlab_webhook handler's fan-out lookup.
func (r *Repo) ListByExternalConnection(ctx context.Context, externalConnectionId uuid.UUID) ([]domain.Trigger, error) {
	rows, err := r.q.ListTriggersByExternalConnection(ctx, externalConnectionId)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing triggers by external connection")
	}

	triggers := make([]domain.Trigger, len(rows))

	for i, row := range rows {
		trigger, convErr := triggerToDomain(row)
		if convErr != nil {
			return nil, convErr
		}

		triggers[i] = trigger
	}

	return triggers, nil
}

// ListLinksByTrigger returns triggerUuid's linked tracts with Tract populated (Trigger left
// zero) — the webhook handler's fan-out: one delivery may start runs on several tracts.
func (r *Repo) ListLinksByTrigger(ctx context.Context, triggerUuid uuid.UUID) ([]repository.TractTriggerLink, error) {
	rows, err := r.q.ListTriggerLinksByTrigger(ctx, triggerUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing trigger tract links")
	}

	links := make([]repository.TractTriggerLink, len(rows))

	for i, row := range rows {
		var def domain.TractDefinition

		unmarshalErr := json.Unmarshal(row.Definition, &def)
		if unmarshalErr != nil {
			return nil, rerrors.Wrap(unmarshalErr, "error unmarshaling tract definition")
		}

		tract := domain.Tract{
			Uuid:        row.ID,
			UserUuid:    row.UserID,
			Name:        row.Name,
			Description: row.Description,
			Enabled:     row.Enabled,
			Definition:  def,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		}

		filters, convErr := unmarshalFilters(row.Filters)
		if convErr != nil {
			return nil, convErr
		}

		links[i] = repository.TractTriggerLink{
			TractUuid:   tract.Uuid,
			TriggerUuid: triggerUuid,
			Tract:       tract,
			Filters:     filters,
		}
	}

	return links, nil
}

func triggerToDomain(row artel_q.Trigger) (domain.Trigger, error) {
	schema, err := unmarshalToolSchema(row.PayloadSchema)
	if err != nil {
		return domain.Trigger{}, err
	}

	matchers, err := unmarshalMatchers(row.Matchers)
	if err != nil {
		return domain.Trigger{}, err
	}

	trigger := domain.Trigger{
		Uuid:          row.ID,
		TriggerUuid:   row.TriggerUuid,
		UserUuid:      row.UserID,
		Name:          row.Name,
		Kind:          row.Kind,
		Source:        row.Source,
		Config:        row.Config,
		PayloadSchema: schema,
		SecretHash:    row.SecretHash,
		TokenSuffix:   row.TokenSuffix.String,
		Matchers:      matchers,
		Enabled:       row.Enabled,
		CreatedAt:     row.CreatedAt,
	}

	return trigger, nil
}

// tool schema (de)serialization — payload_schema is stored in the same {properties, required}
// shape as mcp_tools' input/output schemas, mirrored here as private row types (see
// mcpdefinitions repo) so domain.ToolSchema/ToolProperty stay free of persistence JSON tags.

func marshalToolSchema(schema domain.ToolSchema) (json.RawMessage, error) {
	row := toolSchemaRow{
		Properties: make(map[string]toolPropertyRow, len(schema.Properties)),
		Required:   schema.Required,
	}
	for k, v := range schema.Properties {
		row.Properties[k] = toolPropertyRowFromDomain(v)
	}

	data, err := json.Marshal(row)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshaling trigger payload schema")
	}

	return data, nil
}

func unmarshalToolSchema(raw json.RawMessage) (domain.ToolSchema, error) {
	if len(raw) == 0 {
		return domain.ToolSchema{}, nil
	}

	var row toolSchemaRow

	err := json.Unmarshal(raw, &row)
	if err != nil {
		return domain.ToolSchema{}, rerrors.Wrap(err, "error unmarshaling trigger payload schema")
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

// matchers (de)serialization — mirrors the toolSchemaRow pattern above so domain.TriggerMatchers
// stays free of persistence JSON tags.

func marshalMatchers(matchers domain.TriggerMatchers) (json.RawMessage, error) {
	row := triggerMatchersRow{
		CheckHeaders: make([]headerMatcherRow, len(matchers.CheckHeaders)),
	}
	for i, m := range matchers.CheckHeaders {
		row.CheckHeaders[i] = headerMatcherRow{Header: m.Header, Equals: m.Equals}
	}

	data, err := json.Marshal(row)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshaling trigger matchers")
	}

	return data, nil
}

func unmarshalMatchers(raw json.RawMessage) (domain.TriggerMatchers, error) {
	if len(raw) == 0 {
		return domain.TriggerMatchers{}, nil
	}

	var row triggerMatchersRow

	err := json.Unmarshal(raw, &row)
	if err != nil {
		return domain.TriggerMatchers{}, rerrors.Wrap(err, "error unmarshaling trigger matchers")
	}

	matchers := domain.TriggerMatchers{
		CheckHeaders: make([]domain.HeaderMatcher, len(row.CheckHeaders)),
	}
	for i, m := range row.CheckHeaders {
		matchers.CheckHeaders[i] = domain.HeaderMatcher{Header: m.Header, Equals: m.Equals}
	}

	return matchers, nil
}

type triggerMatchersRow struct {
	CheckHeaders []headerMatcherRow `json:"check_headers,omitempty"`
}

type headerMatcherRow struct {
	Header string `json:"header"`
	Equals string `json:"equals"`
}

// filters (de)serialization — domain.TractCondition already carries json tags (see
// internal/domain/tract.go), so this just wraps json.(Un)marshal.

func marshalFilters(filters []domain.TractCondition) (json.RawMessage, error) {
	if filters == nil {
		filters = []domain.TractCondition{}
	}

	data, err := json.Marshal(filters)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshaling trigger link filters")
	}

	return data, nil
}

func unmarshalFilters(raw json.RawMessage) ([]domain.TractCondition, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var filters []domain.TractCondition

	err := json.Unmarshal(raw, &filters)
	if err != nil {
		return nil, rerrors.Wrap(err, "error unmarshaling trigger link filters")
	}

	return filters, nil
}
