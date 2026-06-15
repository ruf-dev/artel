package mcpdefinitions

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
)

type Repo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) Upsert(ctx context.Context, def domain.McpDefinition) (domain.McpDefinition, error) {
	toolsJSON, err := toolsToJSON(def.Tools)
	if err != nil {
		return domain.McpDefinition{}, rerrors.Wrap(err, "error marshaling tools")
	}

	params := artel_q.UpsertMcpDefinitionParams{
		Name:        def.Name,
		Author:      def.Author,
		Description: def.Description,
		Tools:       toolsJSON,
	}

	row, err := r.q.UpsertMcpDefinition(ctx, params)
	if err != nil {
		return domain.McpDefinition{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error upserting mcp definition")
	}

	return toDomain(row)
}

func (r *Repo) Get(ctx context.Context, name string) (sql.Null[domain.McpDefinition], error) {
	row, err := r.q.GetMcpDefinition(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.McpDefinition]{}, nil
		}
		return sql.Null[domain.McpDefinition]{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting mcp definition")
	}

	def, err := toDomain(row)
	if err != nil {
		return sql.Null[domain.McpDefinition]{}, err
	}

	result := sql.Null[domain.McpDefinition]{V: def, Valid: true}
	return result, nil
}

func (r *Repo) List(ctx context.Context) ([]domain.McpDefinition, error) {
	rows, err := r.q.ListMcpDefinitions(ctx)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing mcp definitions")
	}

	defs := make([]domain.McpDefinition, len(rows))
	for i, row := range rows {
		def, err := toDomain(row)
		if err != nil {
			return nil, err
		}
		defs[i] = def
	}
	return defs, nil
}

func (r *Repo) Delete(ctx context.Context, name string) error {
	err := r.q.DeleteMcpDefinition(ctx, name)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error deleting mcp definition")
	}
	return nil
}

func toDomain(row artel_q.Mcp) (domain.McpDefinition, error) {
	var toolRows []toolDefRow
	err := json.Unmarshal(row.Tools, &toolRows)
	if err != nil {
		return domain.McpDefinition{}, rerrors.Wrap(err, "error unmarshaling mcp tools")
	}

	tools := make([]domain.McpToolDef, len(toolRows))
	for i, tr := range toolRows {
		tools[i] = toolDefToDomain(tr)
	}

	return domain.McpDefinition{
		Name:        row.Name,
		Author:      row.Author,
		Description: row.Description,
		Tools:       tools,
		CreatedAt:   row.CreatedAt,
	}, nil
}

func toolDefToDomain(tr toolDefRow) domain.McpToolDef {
	props := make(map[string]domain.ToolProperty, len(tr.ApiDescription.Properties))
	for k, v := range tr.ApiDescription.Properties {
		props[k] = domain.ToolProperty{
			Type:        v.Type,
			Description: v.Description,
		}
	}

	apiDesc := domain.ToolApiDescription{
		Name:        tr.ApiDescription.Name,
		Description: tr.ApiDescription.Description,
		Properties:  props,
		Required:    tr.ApiDescription.Required,
	}

	return domain.McpToolDef{
		ApiDescription: apiDesc,
		Action:         actionToDomain(tr.Action),
	}
}

func actionToDomain(ar toolActionRow) domain.ToolAction {
	var action domain.ToolAction
	if ar.Imap != nil {
		action.Imap = &domain.ImapAction{Operation: domain.ImapOperation(ar.Imap.Operation)}
	}
	if ar.Smtp != nil {
		action.Smtp = &domain.SmtpAction{Operation: domain.SmtpOperation(ar.Smtp.Operation)}
	}
	if ar.Http != nil {
		action.Http = &domain.HttpAction{
			Method:      ar.Http.Method,
			Url:         ar.Http.Url,
			Headers:     ar.Http.Headers,
			Query:       ar.Http.Query,
			Credentials: ar.Http.Credentials,
		}
	}
	return action
}

func toolsToJSON(tools []domain.McpToolDef) (json.RawMessage, error) {
	rows := make([]toolDefRow, len(tools))
	for i, t := range tools {
		rows[i] = toolDefFromDomain(t)
	}

	data, err := json.Marshal(rows)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(data), nil
}

func toolDefFromDomain(t domain.McpToolDef) toolDefRow {
	props := make(map[string]toolPropertyRow, len(t.ApiDescription.Properties))
	for k, v := range t.ApiDescription.Properties {
		props[k] = toolPropertyRow{Type: v.Type, Description: v.Description}
	}

	return toolDefRow{
		ApiDescription: apiDescriptionRow{
			Name:        t.ApiDescription.Name,
			Description: t.ApiDescription.Description,
			Properties:  props,
			Required:    t.ApiDescription.Required,
		},
		Action: actionFromDomain(t.Action),
	}
}

func actionFromDomain(a domain.ToolAction) toolActionRow {
	var ar toolActionRow
	if a.Imap != nil {
		ar.Imap = &imapActionRow{Operation: string(a.Imap.Operation)}
	}
	if a.Smtp != nil {
		ar.Smtp = &smtpActionRow{Operation: string(a.Smtp.Operation)}
	}
	if a.Http != nil {
		ar.Http = &httpActionRow{
			Method:      a.Http.Method,
			Url:         a.Http.Url,
			Headers:     a.Http.Headers,
			Query:       a.Http.Query,
			Credentials: a.Http.Credentials,
		}
	}
	return ar
}

// private JSON structs — storage representation of the tools JSONB column.

type toolDefRow struct {
	ApiDescription apiDescriptionRow `json:"api_description"`
	Action         toolActionRow     `json:"action"`
}

type apiDescriptionRow struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Properties  map[string]toolPropertyRow `json:"properties"`
	Required    []string                   `json:"required"`
}

type toolPropertyRow struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type toolActionRow struct {
	Imap *imapActionRow `json:"imap"`
	Smtp *smtpActionRow `json:"smtp"`
	Http *httpActionRow `json:"http"`
}

type imapActionRow struct {
	Operation string `json:"operation"`
}

type smtpActionRow struct {
	Operation string `json:"operation"`
}

type httpActionRow struct {
	Method      string            `json:"method"`
	Url         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Query       map[string]string `json:"query"`
	Credentials string            `json:"credentials"`
}
