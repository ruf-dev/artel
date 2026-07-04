package mcpdefinitions

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

// Upsert writes the mcps row and, for every entry in def.Tools, upserts the matching
// mcp_tools row (never deletes — a tool dropped from def.Tools is left in place, per the
// "seeds must upsert, never delete tools" convention).
func (r *Repo) Upsert(ctx context.Context, def domain.McpDefinition) (domain.McpDefinition, error) {
	params := artel_q.UpsertMcpDefinitionParams{
		Name:        def.Name,
		Author:      def.Author,
		Description: def.Description,
	}

	row, err := r.q.UpsertMcpDefinition(ctx, params)
	if err != nil {
		return domain.McpDefinition{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error upserting mcp definition")
	}

	for _, tool := range def.Tools {
		err = r.upsertTool(ctx, def.Name, tool)
		if err != nil {
			return domain.McpDefinition{}, err
		}
	}

	tools, err := r.listTools(ctx, def.Name)
	if err != nil {
		return domain.McpDefinition{}, err
	}

	result := domain.McpDefinition{
		Name:        row.Name,
		Author:      row.Author,
		Description: row.Description,
		Tools:       tools,
		CreatedAt:   row.CreatedAt,
	}

	return result, nil
}

func (r *Repo) upsertTool(ctx context.Context, mcpName string, tool domain.McpToolDef) error {
	params, err := upsertToolParams(mcpName, tool)
	if err != nil {
		return err
	}

	_, err = r.q.UpsertMcpTool(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error upserting mcp tool")
	}

	return nil
}

func (r *Repo) Get(ctx context.Context, name string) (sql.Null[domain.McpDefinition], error) {
	row, err := r.q.GetMcpDefinition(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.McpDefinition]{}, nil
		}

		return sql.Null[domain.McpDefinition]{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting mcp definition")
	}

	tools, err := r.listTools(ctx, row.Name)
	if err != nil {
		return sql.Null[domain.McpDefinition]{}, err
	}

	def := domain.McpDefinition{
		Name:        row.Name,
		Author:      row.Author,
		Description: row.Description,
		Tools:       tools,
		CreatedAt:   row.CreatedAt,
	}

	result := sql.Null[domain.McpDefinition]{V: def, Valid: true}

	return result, nil
}

func (r *Repo) List(ctx context.Context) ([]domain.McpDefinition, error) {
	rows, err := r.q.ListMcpDefinitions(ctx)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing mcp definitions")
	}

	toolRows, err := r.q.ListAllMcpTools(ctx)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing mcp tools")
	}

	toolsByMcp := make(map[string][]domain.McpToolDef, len(rows))

	for _, tr := range toolRows {
		tool, err := toolRowToDomain(tr)
		if err != nil {
			return nil, err
		}

		toolsByMcp[tr.McpName] = append(toolsByMcp[tr.McpName], tool)
	}

	defs := make([]domain.McpDefinition, len(rows))

	for i, row := range rows {
		tools := toolsByMcp[row.Name]
		if tools == nil {
			tools = []domain.McpToolDef{}
		}

		defs[i] = domain.McpDefinition{
			Name:        row.Name,
			Author:      row.Author,
			Description: row.Description,
			Tools:       tools,
			CreatedAt:   row.CreatedAt,
		}
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

// GetTool looks up a single tool by (mcpName, name) without loading the rest of the MoM's
// tool set — used by the tool-execution and catalog paths that already know which tool
// they want.
func (r *Repo) GetTool(ctx context.Context, mcpName string, toolName string) (sql.Null[domain.McpToolDef], error) {
	params := artel_q.GetMcpToolParams{
		McpName: mcpName,
		Name:    toolName,
	}

	row, err := r.q.GetMcpTool(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.McpToolDef]{}, nil
		}

		return sql.Null[domain.McpToolDef]{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting mcp tool")
	}

	tool, err := toolRowToDomain(row)
	if err != nil {
		return sql.Null[domain.McpToolDef]{}, err
	}

	result := sql.Null[domain.McpToolDef]{V: tool, Valid: true}

	return result, nil
}

// ListAllTools returns every tool across every MoM, paired with its owning MoM's name —
// the catalog backing Tract's "builtins + mcp_tools" action list.
func (r *Repo) ListAllTools(ctx context.Context) ([]domain.McpToolRef, error) {
	rows, err := r.q.ListAllMcpTools(ctx)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing all mcp tools")
	}

	refs := make([]domain.McpToolRef, len(rows))

	for i, row := range rows {
		tool, err := toolRowToDomain(row)
		if err != nil {
			return nil, err
		}

		refs[i] = domain.McpToolRef{McpName: row.McpName, Tool: tool}
	}

	return refs, nil
}

func (r *Repo) listTools(ctx context.Context, mcpName string) ([]domain.McpToolDef, error) {
	rows, err := r.q.ListMcpToolsByMcpName(ctx, mcpName)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing mcp tools")
	}

	tools := make([]domain.McpToolDef, len(rows))

	for i, row := range rows {
		tool, err := toolRowToDomain(row)
		if err != nil {
			return nil, err
		}

		tools[i] = tool
	}

	return tools, nil
}

func upsertToolParams(mcpName string, tool domain.McpToolDef) (artel_q.UpsertMcpToolParams, error) {
	inputSchema := toolSchemaFromDomain(domain.ToolSchema{
		Properties: tool.ApiDescription.Properties,
		Required:   tool.ApiDescription.Required,
	})

	inputSchemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		return artel_q.UpsertMcpToolParams{}, rerrors.Wrap(err, "error marshaling tool input schema")
	}

	outputSchema := toolSchemaFromDomain(tool.OutputSchema)

	outputSchemaJSON, err := json.Marshal(outputSchema)
	if err != nil {
		return artel_q.UpsertMcpToolParams{}, rerrors.Wrap(err, "error marshaling tool output schema")
	}

	actionJSON, err := json.Marshal(actionFromDomain(tool.Action))
	if err != nil {
		return artel_q.UpsertMcpToolParams{}, rerrors.Wrap(err, "error marshaling tool action")
	}

	params := artel_q.UpsertMcpToolParams{
		McpName:      mcpName,
		Name:         tool.ApiDescription.Name,
		Description:  tool.ApiDescription.Description,
		InputSchema:  inputSchemaJSON,
		OutputSchema: outputSchemaJSON,
		Action:       actionJSON,
	}

	return params, nil
}

func toolRowToDomain(row artel_q.McpTool) (domain.McpToolDef, error) {
	var inputSchemaRow toolSchemaRow

	err := json.Unmarshal(row.InputSchema, &inputSchemaRow)
	if err != nil {
		return domain.McpToolDef{}, rerrors.Wrap(err, "error unmarshaling tool input schema")
	}

	var outputSchemaRow toolSchemaRow

	err = json.Unmarshal(row.OutputSchema, &outputSchemaRow)
	if err != nil {
		return domain.McpToolDef{}, rerrors.Wrap(err, "error unmarshaling tool output schema")
	}

	var actionRow toolActionRow

	err = json.Unmarshal(row.Action, &actionRow)
	if err != nil {
		return domain.McpToolDef{}, rerrors.Wrap(err, "error unmarshaling tool action")
	}

	inputSchema := toolSchemaToDomain(inputSchemaRow)

	apiDesc := domain.ToolApiDescription{
		Name:        row.Name,
		Description: row.Description,
		Properties:  inputSchema.Properties,
		Required:    inputSchema.Required,
	}

	tool := domain.McpToolDef{
		ApiDescription: apiDesc,
		OutputSchema:   toolSchemaToDomain(outputSchemaRow),
		Action:         actionToDomain(actionRow),
	}

	return tool, nil
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
			Body:        ar.Http.Body,
			Credentials: ar.Http.Credentials,
		}
	}

	return action
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
			Body:        a.Http.Body,
			Credentials: a.Http.Credentials,
		}
	}

	return ar
}

// toolPropertyToDomain / toolPropertyFromDomain recurse through nested object/array shapes.

func toolPropertyToDomain(p toolPropertyRow) domain.ToolProperty {
	prop := domain.ToolProperty{
		Type:        p.Type,
		Description: p.Description,
		Enum:        p.Enum,
		Required:    p.Required,
	}

	if len(p.Properties) > 0 {
		prop.Properties = make(map[string]domain.ToolProperty, len(p.Properties))
		for k, v := range p.Properties {
			prop.Properties[k] = toolPropertyToDomain(v)
		}
	}

	if p.Items != nil {
		items := toolPropertyToDomain(*p.Items)
		prop.Items = &items
	}

	return prop
}

func toolPropertyFromDomain(p domain.ToolProperty) toolPropertyRow {
	row := toolPropertyRow{
		Type:        p.Type,
		Description: p.Description,
		Enum:        p.Enum,
		Required:    p.Required,
	}

	if len(p.Properties) > 0 {
		row.Properties = make(map[string]toolPropertyRow, len(p.Properties))
		for k, v := range p.Properties {
			row.Properties[k] = toolPropertyFromDomain(v)
		}
	}

	if p.Items != nil {
		items := toolPropertyFromDomain(*p.Items)
		row.Items = &items
	}

	return row
}

func toolSchemaToDomain(s toolSchemaRow) domain.ToolSchema {
	props := make(map[string]domain.ToolProperty, len(s.Properties))
	for k, v := range s.Properties {
		props[k] = toolPropertyToDomain(v)
	}

	return domain.ToolSchema{Properties: props, Required: s.Required}
}

func toolSchemaFromDomain(s domain.ToolSchema) toolSchemaRow {
	props := make(map[string]toolPropertyRow, len(s.Properties))
	for k, v := range s.Properties {
		props[k] = toolPropertyFromDomain(v)
	}

	return toolSchemaRow{Properties: props, Required: s.Required}
}

// private JSON structs — storage representation of the input_schema/output_schema/action
// JSONB columns.

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
	Body        json.RawMessage   `json:"body,omitempty"`
	Credentials string            `json:"credentials"`
}
