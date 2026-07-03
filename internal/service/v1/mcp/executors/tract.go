package executors

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/requesthost"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/tract"
)

const (
	ToolListTractActions = "list_tract_actions"
	ToolCreateTract      = "create_tract"
	ToolUpdateTract      = "update_tract"
	ToolCreateTrigger    = "create_trigger"
	ToolLinkTrigger      = "link_trigger"
	ToolListTracts       = "list_tracts"
	ToolGetTract         = "get_tract"
	ToolGetTractRuns     = "get_tract_runs"
	ToolRunTract         = "run_tract"
)

const (
	triggerKindWebhook = "webhook"
	webhookPathPrefix  = "/tract/hook/"
	timeFormat         = "2006-01-02T15:04:05Z"
	defaultRunsLimit   = 20
)

// TractService is the narrow slice of service.TractService the tract-authoring builtin tools
// need. Defined here — not imported from internal/service — so this package never depends on
// the top-level service interfaces package; internal/service/v1/mcp passes its concrete
// service.TractService value in per call, which satisfies this interface structurally.
type TractService interface {
	CreateTract(ctx context.Context, name string, description string, def domain.TractDefinition) (domain.Tract, []string, error)
	UpdateTract(ctx context.Context, id uuid.UUID, name string, description string, def domain.TractDefinition) (domain.Tract, []string, error)
	GetTract(ctx context.Context, id uuid.UUID) (domain.Tract, error)
	ListTracts(ctx context.Context) ([]domain.Tract, error)
	CreateTrigger(ctx context.Context, name string, kind string, source string, config json.RawMessage, payloadSchema domain.ToolSchema) (domain.Trigger, string, error)
	LinkTrigger(ctx context.Context, triggerUuid uuid.UUID, tractUuid uuid.UUID, filters []domain.TractCondition) error
	ListLinksByTract(ctx context.Context, tractUuid uuid.UUID) ([]repository.TractTriggerLink, error)
	StartRun(ctx context.Context, tract domain.Tract, payload json.RawMessage, startedBy string, triggerUuid uuid.UUID) (domain.TractRun, error)
	ListRuns(ctx context.Context, tractUuid uuid.UUID, limit int32) ([]domain.TractRun, error)
	GetRun(ctx context.Context, id uuid.UUID) (domain.TractRun, []domain.TractRunStep, error)
	ListTractTools(ctx context.Context) ([]domain.McpToolRef, error)
	ListTriggerSources(ctx context.Context) ([]domain.TriggerSourcePreset, error)
}

// TractExecutor implements the tract-authoring builtin tools (list_tract_actions, create_tract,
// ...) that let an agent build and run tracts over MCP. It holds no state of its own — the
// tract service, the server-lifecycle context for run_tract's async StartRun call, and the
// caller's user_context are all supplied per call by internal/service/v1/mcp, which owns the
// service.TractService dependency (wired via McpService.SetTractService).
type TractExecutor struct{}

func NewTractExecutor() *TractExecutor {
	return &TractExecutor{}
}

// IsTractTool reports whether name is one of the tract-authoring builtin tools.
func IsTractTool(name string) bool {
	switch name {
	case ToolListTractActions, ToolCreateTract, ToolUpdateTract, ToolCreateTrigger, ToolLinkTrigger,
		ToolListTracts, ToolGetTract, ToolGetTractRuns, ToolRunTract:
		return true
	}
	return false
}

// Execute dispatches one tract-authoring builtin tool call. ctx must already carry the caller's
// user_context (the tract service's CreateTract/GetTract/... all require it for ownership
// checks). baseCtx is the server-lifecycle context run_tract spawns StartRun against, since the
// per-request ctx is cancelled once the MCP handler's response is written.
func (e *TractExecutor) Execute(ctx context.Context, baseCtx context.Context, ts TractService, toolName string, params map[string]interface{}) (domain.ToolExecResult, error) {
	switch toolName {
	case ToolListTractActions:
		return e.listTractActions(ctx, ts)
	case ToolCreateTract:
		return e.createTract(ctx, ts, params)
	case ToolUpdateTract:
		return e.updateTract(ctx, ts, params)
	case ToolCreateTrigger:
		return e.createTrigger(ctx, ts, params)
	case ToolLinkTrigger:
		return e.linkTrigger(ctx, ts, params)
	case ToolListTracts:
		return e.listTracts(ctx, ts)
	case ToolGetTract:
		return e.getTract(ctx, ts, params)
	case ToolGetTractRuns:
		return e.getTractRuns(ctx, ts, params)
	case ToolRunTract:
		return e.runTract(ctx, baseCtx, ts, params)
	default:
		return domain.ToolExecResult{}, rerrors.Wrap(user_errors.McpTractUnknownToolName, toolName)
	}
}

// -- list_tract_actions --

type tractActionEntry struct {
	Mcp          string        `json:"mcp"`
	Tool         string        `json:"tool"`
	Description  string        `json:"description"`
	InputSchema  toolSchemaRow `json:"input_schema"`
	OutputSchema toolSchemaRow `json:"output_schema"`
}

type triggerSourceEntry struct {
	Key           string        `json:"key"`
	Description   string        `json:"description"`
	PayloadSchema toolSchemaRow `json:"payload_schema"`
}

type tractActionsResult struct {
	Actions        []tractActionEntry   `json:"actions"`
	TriggerSources []triggerSourceEntry `json:"trigger_sources"`
}

func (e *TractExecutor) listTractActions(ctx context.Context, ts TractService) (domain.ToolExecResult, error) {
	refs, err := ts.ListTractTools(ctx)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error listing tract tools")
	}

	sources, err := ts.ListTriggerSources(ctx)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error listing trigger sources")
	}

	actions := make([]tractActionEntry, len(refs))
	for i, ref := range refs {
		inputSchema := domain.ToolSchema{Properties: ref.Tool.ApiDescription.Properties, Required: ref.Tool.ApiDescription.Required}
		actions[i] = tractActionEntry{
			Mcp:          ref.McpName,
			Tool:         ref.Tool.ApiDescription.Name,
			Description:  ref.Tool.ApiDescription.Description,
			InputSchema:  toolSchemaToRow(inputSchema),
			OutputSchema: toolSchemaToRow(ref.Tool.OutputSchema),
		}
	}

	sourceEntries := make([]triggerSourceEntry, len(sources))
	for i, s := range sources {
		sourceEntries[i] = triggerSourceEntry{
			Key:           s.Key,
			Description:   s.Description,
			PayloadSchema: toolSchemaToRow(s.PayloadSchema),
		}
	}

	result := tractActionsResult{Actions: actions, TriggerSources: sourceEntries}

	text, err := marshalResult(result)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

// -- create_tract --

type tractRow struct {
	Uuid        string                  `json:"uuid"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Enabled     bool                    `json:"enabled"`
	Definition  *domain.TractDefinition `json:"definition,omitempty"`
	CreatedAt   string                  `json:"created_at,omitempty"`
	UpdatedAt   string                  `json:"updated_at,omitempty"`
}

func tractToRow(t domain.Tract, includeDefinition bool) tractRow {
	row := tractRow{
		Uuid:        t.Uuid.String(),
		Name:        t.Name,
		Description: t.Description,
		Enabled:     t.Enabled,
		CreatedAt:   t.CreatedAt.UTC().Format(timeFormat),
		UpdatedAt:   t.UpdatedAt.UTC().Format(timeFormat),
	}
	if includeDefinition {
		def := t.Definition
		row.Definition = &def
	}
	return row
}

type webhookTriggerRow struct {
	TriggerUuid  string `json:"trigger_uuid"`
	WebhookUrl   string `json:"webhook_url"`
	WebhookToken string `json:"webhook_token"`
}

func (e *TractExecutor) createTract(ctx context.Context, ts TractService, params map[string]interface{}) (domain.ToolExecResult, error) {
	name, ok := paramString(params, "name")
	if !ok || name == "" {
		return domain.ToolExecResult{}, user_errors.McpTractNameRequired
	}

	description, _ := paramString(params, "description")

	defParam, ok := paramObject(params, "definition")
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpTractDefinitionRequired
	}

	var def domain.TractDefinition
	err := decodeInto(defParam, &def)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error decoding definition")
	}

	created, warnings, err := ts.CreateTract(ctx, name, description, def)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error creating tract")
	}

	result := map[string]interface{}{
		"tract":    tractToRow(created, true),
		"warnings": warnings,
	}

	webhookParam, ok := paramObject(params, "webhook_trigger")
	if ok {
		webhookResult, err := e.createLinkedWebhook(ctx, ts, created.Uuid, webhookParam)
		if err != nil {
			return domain.ToolExecResult{}, err
		}
		result["webhook_trigger"] = webhookResult
	}

	text, err := marshalResult(result)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

func (e *TractExecutor) createLinkedWebhook(ctx context.Context, ts TractService, tractUuid uuid.UUID, wt map[string]interface{}) (webhookTriggerRow, error) {
	name, ok := paramString(wt, "name")
	if !ok || name == "" {
		return webhookTriggerRow{}, user_errors.McpTractNameRequired
	}

	source, ok := paramString(wt, "source")
	if !ok || source == "" {
		return webhookTriggerRow{}, user_errors.McpTriggerSourceRequired
	}

	var filters []domain.TractCondition
	rawFilters, ok := wt["filters"]
	if ok {
		err := decodeInto(rawFilters, &filters)
		if err != nil {
			return webhookTriggerRow{}, rerrors.Wrap(err, "error decoding webhook_trigger filters")
		}
	}

	emptyConfig := json.RawMessage(`{}`)
	emptySchema := domain.ToolSchema{}
	trig, rawToken, err := ts.CreateTrigger(ctx, name, triggerKindWebhook, source, emptyConfig, emptySchema)
	if err != nil {
		return webhookTriggerRow{}, rerrors.Wrap(err, "error creating trigger")
	}

	err = ts.LinkTrigger(ctx, trig.Uuid, tractUuid, filters)
	if err != nil {
		return webhookTriggerRow{}, rerrors.Wrap(err, "error linking trigger to tract")
	}

	row := webhookTriggerRow{
		TriggerUuid:  trig.Uuid.String(),
		WebhookUrl:   buildWebhookUrl(ctx, trig.TriggerUuid),
		WebhookToken: rawToken,
	}
	return row, nil
}

func buildWebhookUrl(ctx context.Context, triggerUuid uuid.UUID) string {
	host, ok := requesthost.FromContext(ctx)
	if !ok || host == "" {
		return webhookPathPrefix + triggerUuid.String()
	}
	return "https://" + host + webhookPathPrefix + triggerUuid.String()
}

// -- update_tract --

func (e *TractExecutor) updateTract(ctx context.Context, ts TractService, params map[string]interface{}) (domain.ToolExecResult, error) {
	id, err := paramUuid(params, "tract_uuid", user_errors.McpTractUuidRequired, user_errors.McpTractUuidInvalid)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	name, ok := paramString(params, "name")
	if !ok || name == "" {
		return domain.ToolExecResult{}, user_errors.McpTractNameRequired
	}

	description, _ := paramString(params, "description")

	defParam, ok := paramObject(params, "definition")
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpTractDefinitionRequired
	}

	var def domain.TractDefinition
	err = decodeInto(defParam, &def)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error decoding definition")
	}

	updated, warnings, err := ts.UpdateTract(ctx, id, name, description, def)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error updating tract")
	}

	result := map[string]interface{}{
		"tract":    tractToRow(updated, true),
		"warnings": warnings,
	}

	text, err := marshalResult(result)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

// -- create_trigger --

func (e *TractExecutor) createTrigger(ctx context.Context, ts TractService, params map[string]interface{}) (domain.ToolExecResult, error) {
	name, ok := paramString(params, "name")
	if !ok || name == "" {
		return domain.ToolExecResult{}, user_errors.McpTractNameRequired
	}

	source, ok := paramString(params, "source")
	if !ok || source == "" {
		return domain.ToolExecResult{}, user_errors.McpTriggerSourceRequired
	}

	schema, err := toolSchemaFromParam(params["payload_schema"])
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error decoding payload_schema")
	}

	emptyConfig := json.RawMessage(`{}`)
	trig, rawToken, err := ts.CreateTrigger(ctx, name, triggerKindWebhook, source, emptyConfig, schema)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error creating trigger")
	}

	result := webhookTriggerRow{
		TriggerUuid:  trig.Uuid.String(),
		WebhookUrl:   buildWebhookUrl(ctx, trig.TriggerUuid),
		WebhookToken: rawToken,
	}

	text, err := marshalResult(result)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

// -- link_trigger --

func (e *TractExecutor) linkTrigger(ctx context.Context, ts TractService, params map[string]interface{}) (domain.ToolExecResult, error) {
	triggerUuid, err := paramUuid(params, "trigger_uuid", user_errors.McpTriggerUuidRequired, user_errors.McpTriggerUuidInvalid)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	tractUuid, err := paramUuid(params, "tract_uuid", user_errors.McpTractUuidRequired, user_errors.McpTractUuidInvalid)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	var filters []domain.TractCondition
	rawFilters, ok := params["filters"]
	if ok {
		err = decodeInto(rawFilters, &filters)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "error decoding filters")
		}
	}

	err = ts.LinkTrigger(ctx, triggerUuid, tractUuid, filters)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error linking trigger to tract")
	}

	result := map[string]interface{}{"message": "trigger linked"}
	text, err := marshalResult(result)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

// -- list_tracts --

func (e *TractExecutor) listTracts(ctx context.Context, ts TractService) (domain.ToolExecResult, error) {
	tracts, err := ts.ListTracts(ctx)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error listing tracts")
	}

	rows := make([]tractRow, len(tracts))
	for i, t := range tracts {
		rows[i] = tractToRow(t, false)
	}

	text, err := marshalResult(rows)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

// -- get_tract --

type triggerLinkRow struct {
	TriggerUuid string                  `json:"trigger_uuid"`
	Name        string                  `json:"name"`
	Kind        string                  `json:"kind"`
	Source      string                  `json:"source"`
	Enabled     bool                    `json:"enabled"`
	Filters     []domain.TractCondition `json:"filters,omitempty"`
}

type tractDetailRow struct {
	Tract    tractRow         `json:"tract"`
	Triggers []triggerLinkRow `json:"triggers"`
}

func (e *TractExecutor) getTract(ctx context.Context, ts TractService, params map[string]interface{}) (domain.ToolExecResult, error) {
	id, err := paramUuid(params, "tract_uuid", user_errors.McpTractUuidRequired, user_errors.McpTractUuidInvalid)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	tr, err := ts.GetTract(ctx, id)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error getting tract")
	}

	links, err := ts.ListLinksByTract(ctx, id)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error listing tract trigger links")
	}

	triggerRows := make([]triggerLinkRow, len(links))
	for i, link := range links {
		triggerRows[i] = triggerLinkRow{
			TriggerUuid: link.Trigger.Uuid.String(),
			Name:        link.Trigger.Name,
			Kind:        link.Trigger.Kind,
			Source:      link.Trigger.Source,
			Enabled:     link.Trigger.Enabled,
			Filters:     link.Filters,
		}
	}

	detail := tractDetailRow{
		Tract:    tractToRow(tr, true),
		Triggers: triggerRows,
	}

	text, err := marshalResult(detail)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

// -- get_tract_runs --

type runStepRow struct {
	StepId     string          `json:"step_id"`
	StepName   string          `json:"step_name"`
	StepType   string          `json:"step_type"`
	Status     string          `json:"status"`
	Input      json.RawMessage `json:"input,omitempty"`
	Output     json.RawMessage `json:"output,omitempty"`
	Error      string          `json:"error,omitempty"`
	StartedAt  string          `json:"started_at,omitempty"`
	FinishedAt string          `json:"finished_at,omitempty"`
}

type runRow struct {
	Uuid           string          `json:"uuid"`
	Status         string          `json:"status"`
	StartedBy      string          `json:"started_by"`
	TriggerPayload json.RawMessage `json:"trigger_payload,omitempty"`
	Error          string          `json:"error,omitempty"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
	Steps          []runStepRow    `json:"steps"`
}

func runToRow(r domain.TractRun, steps []domain.TractRunStep) runRow {
	stepRows := make([]runStepRow, len(steps))
	for i, s := range steps {
		stepRows[i] = stepToRow(s)
	}

	row := runRow{
		Uuid:           r.Uuid.String(),
		Status:         string(r.Status),
		StartedBy:      r.StartedBy,
		TriggerPayload: r.TriggerPayload,
		Error:          r.Error,
		CreatedAt:      r.CreatedAt.UTC().Format(timeFormat),
		UpdatedAt:      r.UpdatedAt.UTC().Format(timeFormat),
		Steps:          stepRows,
	}
	return row
}

func stepToRow(s domain.TractRunStep) runStepRow {
	finishedAt := ""
	if !s.FinishedAt.IsZero() {
		finishedAt = s.FinishedAt.UTC().Format(timeFormat)
	}

	row := runStepRow{
		StepId:     s.StepId,
		StepName:   s.StepName,
		StepType:   s.StepType,
		Status:     string(s.Status),
		Input:      s.Input,
		Output:     s.Output,
		Error:      s.Error,
		StartedAt:  s.StartedAt.UTC().Format(timeFormat),
		FinishedAt: finishedAt,
	}
	return row
}

func (e *TractExecutor) getTractRuns(ctx context.Context, ts TractService, params map[string]interface{}) (domain.ToolExecResult, error) {
	id, err := paramUuid(params, "tract_uuid", user_errors.McpTractUuidRequired, user_errors.McpTractUuidInvalid)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	limit := paramInt32(params, "limit", defaultRunsLimit)

	runs, err := ts.ListRuns(ctx, id, limit)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error listing tract runs")
	}

	rows := make([]runRow, len(runs))
	for i, run := range runs {
		_, steps, err := ts.GetRun(ctx, run.Uuid)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "error getting tract run")
		}
		rows[i] = runToRow(run, steps)
	}

	text, err := marshalResult(rows)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

// -- run_tract --

func (e *TractExecutor) runTract(ctx context.Context, baseCtx context.Context, ts TractService, params map[string]interface{}) (domain.ToolExecResult, error) {
	id, err := paramUuid(params, "tract_uuid", user_errors.McpTractUuidRequired, user_errors.McpTractUuidInvalid)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	tr, err := ts.GetTract(ctx, id)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error getting tract")
	}

	payload := json.RawMessage(`{}`)
	runParams, ok := params["params"]
	if ok {
		encoded, err := json.Marshal(runParams)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "error marshaling run params")
		}
		payload = encoded
	}

	go e.runAsync(baseCtx, ts, tr, payload)

	result := map[string]interface{}{"message": "tract run started", "tract_uuid": tr.Uuid.String()}
	text, err := marshalResult(result)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

// runAsync spawns StartRun against baseCtx (the server-lifecycle context) rather than the
// per-request ctx, which net/http cancels once the MCP handler's response is written — mirrors
// tracts_api.TractsImpl.runManual.
func (e *TractExecutor) runAsync(baseCtx context.Context, ts TractService, tr domain.Tract, payload json.RawMessage) {
	triggerUuid := uuid.Nil
	_, err := ts.StartRun(baseCtx, tr, payload, tract.StartedByMcp, triggerUuid)
	if err != nil {
		log.Error().Err(err).Str("tract_uuid", tr.Uuid.String()).Msg("mcp tract run failed")
	}
}

// -- param helpers --

func paramString(params map[string]interface{}, key string) (string, bool) {
	v, ok := params[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func paramObject(params map[string]interface{}, key string) (map[string]interface{}, bool) {
	v, ok := params[key]
	if !ok {
		return nil, false
	}
	m, ok := v.(map[string]interface{})
	return m, ok
}

func paramInt32(params map[string]interface{}, key string, fallback int32) int32 {
	v, ok := params[key]
	if !ok {
		return fallback
	}

	switch n := v.(type) {
	case float64:
		return int32(n)
	case int:
		return int32(n)
	case int32:
		return n
	case int64:
		return int32(n)
	}
	return fallback
}

func paramUuid(params map[string]interface{}, key string, requiredErr error, invalidErr error) (uuid.UUID, error) {
	s, ok := paramString(params, key)
	if !ok || s == "" {
		return uuid.Nil, requiredErr
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, invalidErr
	}
	return id, nil
}

// decodeInto round-trips v (already-decoded JSON, e.g. a map[string]interface{} from tool call
// arguments) through json.Marshal/Unmarshal into out — the same technique tracts_api uses to
// turn transport-layer JSON into domain.TractDefinition/[]domain.TractCondition, both of which
// already carry json tags.
func decodeInto(v interface{}, out interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return rerrors.Wrap(err, "error marshaling param")
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return rerrors.Wrap(err, "error unmarshaling param")
	}
	return nil
}

// -- JSON <-> domain.ToolSchema/ToolProperty --
//
// domain.ToolSchema/ToolProperty carry no json tags (kept persistence-agnostic), so every layer
// that needs a wire representation mirrors them locally with tagged row types — see
// internal/transport/tracts_api/to_proto.go for the original. Reused here both to decode the
// optional payload_schema tool-call param and to render list_tract_actions' catalog schemas.

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

func toolSchemaFromParam(v interface{}) (domain.ToolSchema, error) {
	if v == nil {
		return domain.ToolSchema{}, nil
	}

	var row toolSchemaRow
	err := decodeInto(v, &row)
	if err != nil {
		return domain.ToolSchema{}, err
	}

	schema := domain.ToolSchema{Properties: make(map[string]domain.ToolProperty, len(row.Properties)), Required: row.Required}
	for k, pr := range row.Properties {
		schema.Properties[k] = toolPropertyRowToDomain(pr)
	}
	return schema, nil
}

func toolSchemaToRow(schema domain.ToolSchema) toolSchemaRow {
	row := toolSchemaRow{Properties: make(map[string]toolPropertyRow, len(schema.Properties)), Required: schema.Required}
	for k, v := range schema.Properties {
		row.Properties[k] = toolPropertyRowFromDomain(v)
	}
	return row
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

// -- tool catalog --
//
// Output schemas below are hints for tool-picker/port-chip UIs, not enforced at runtime (see
// executors/vault.go). list_tracts/get_tract_runs marshal a JSON array as their text result;
// since ToolSchema (unlike ToolProperty) has no array/root type of its own, Properties there
// describes the shape of each array element rather than the response envelope itself.

func TractToolDefinitions() []domain.McpToolDef {
	filterItemSchema := domain.ToolProperty{
		Type: "object",
		Properties: map[string]domain.ToolProperty{
			"left":  {Type: "string", Description: "Template string, e.g. \"${trigger.branch}\""},
			"op":    {Type: "string", Description: "Comparison operator, e.g. \"==\", \"!=\", \"contains\""},
			"right": {Type: "string", Description: "Template string"},
		},
		Required: []string{"left", "op", "right"},
	}

	return []domain.McpToolDef{
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolListTractActions,
				Description: "List the action catalog available to tract steps (builtin tools + every connected MoM tool), plus available webhook trigger source presets",
				Properties:  map[string]domain.ToolProperty{},
				Required:    []string{},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"actions": {
						Type:        "array",
						Description: "Available tract action tools",
						Items: &domain.ToolProperty{
							Type: "object",
							Properties: map[string]domain.ToolProperty{
								"mcp":           {Type: "string", Description: "MoM name, or \"artel\" for builtins"},
								"tool":          {Type: "string"},
								"description":   {Type: "string"},
								"input_schema":  {Type: "object", Description: "{properties, required} input schema"},
								"output_schema": {Type: "object", Description: "{properties, required} output schema"},
							},
							Required: []string{"mcp", "tool", "description"},
						},
					},
					"trigger_sources": {
						Type:        "array",
						Description: "Available webhook trigger source presets",
						Items: &domain.ToolProperty{
							Type: "object",
							Properties: map[string]domain.ToolProperty{
								"key":            {Type: "string"},
								"description":    {Type: "string"},
								"payload_schema": {Type: "object"},
							},
							Required: []string{"key", "description"},
						},
					},
				},
				Required: []string{"actions", "trigger_sources"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolCreateTract,
				Description: "Create a new tract (workflow automation). Optionally also create and link a webhook trigger in the same call. Use list_tract_actions first to see available action tools.",
				Properties: map[string]domain.ToolProperty{
					"name":        {Type: "string", Description: "Tract display name"},
					"description": {Type: "string", Description: "Tract description (optional)"},
					"definition":  {Type: "object", Description: "Tract definition: {steps: [...]}"},
					"webhook_trigger": {
						Type:        "object",
						Description: "Optional: also create a webhook trigger and link it to this tract",
						Properties: map[string]domain.ToolProperty{
							"name":   {Type: "string"},
							"source": {Type: "string", Description: "Trigger source preset key, e.g. gitlab_push or generic"},
							"filters": {
								Type:  "array",
								Items: &filterItemSchema,
							},
						},
						Required: []string{"name", "source"},
					},
				},
				Required: []string{"name", "definition"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"tract": {
						Type: "object",
						Properties: map[string]domain.ToolProperty{
							"uuid":        {Type: "string"},
							"name":        {Type: "string"},
							"description": {Type: "string"},
							"enabled":     {Type: "boolean"},
							"definition":  {Type: "object"},
							"created_at":  {Type: "string"},
							"updated_at":  {Type: "string"},
						},
					},
					"warnings": {Type: "array", Description: "Non-fatal validation warnings", Items: &domain.ToolProperty{Type: "string"}},
					"webhook_trigger": {
						Type:        "object",
						Description: "Present only when webhook_trigger was requested",
						Properties: map[string]domain.ToolProperty{
							"trigger_uuid":  {Type: "string"},
							"webhook_url":   {Type: "string"},
							"webhook_token": {Type: "string", Description: "One-time secret — save it now, it is never shown again"},
						},
					},
				},
				Required: []string{"tract", "warnings"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolUpdateTract,
				Description: "Replace an existing tract's name, description and whole definition. This is a wholesale overwrite, not a partial patch — pass the full desired definition, not just the changed steps. Use get_tract first to fetch the current definition to edit.",
				Properties: map[string]domain.ToolProperty{
					"tract_uuid":  {Type: "string"},
					"name":        {Type: "string", Description: "Tract display name"},
					"description": {Type: "string", Description: "Tract description (optional)"},
					"definition":  {Type: "object", Description: "Tract definition: {steps: [...]}"},
				},
				Required: []string{"tract_uuid", "name", "definition"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"tract": {
						Type: "object",
						Properties: map[string]domain.ToolProperty{
							"uuid":        {Type: "string"},
							"name":        {Type: "string"},
							"description": {Type: "string"},
							"enabled":     {Type: "boolean"},
							"definition":  {Type: "object"},
							"created_at":  {Type: "string"},
							"updated_at":  {Type: "string"},
						},
					},
					"warnings": {Type: "array", Description: "Non-fatal validation warnings", Items: &domain.ToolProperty{Type: "string"}},
				},
				Required: []string{"tract", "warnings"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolCreateTrigger,
				Description: "Create a standalone webhook trigger, not yet linked to any tract. Use link_trigger to wire it up afterwards.",
				Properties: map[string]domain.ToolProperty{
					"name":           {Type: "string"},
					"source":         {Type: "string", Description: "Trigger source preset key, e.g. gitlab_push or generic"},
					"payload_schema": {Type: "object", Description: "Optional {properties, required} schema describing the webhook payload shape"},
				},
				Required: []string{"name", "source"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"trigger_uuid":  {Type: "string"},
					"webhook_url":   {Type: "string"},
					"webhook_token": {Type: "string", Description: "One-time secret — save it now, it is never shown again"},
				},
				Required: []string{"trigger_uuid", "webhook_url", "webhook_token"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolLinkTrigger,
				Description: "Link an existing trigger to an existing tract, with optional filters",
				Properties: map[string]domain.ToolProperty{
					"trigger_uuid": {Type: "string"},
					"tract_uuid":   {Type: "string"},
					"filters": {
						Type:  "array",
						Items: &filterItemSchema,
					},
				},
				Required: []string{"trigger_uuid", "tract_uuid"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"message": {Type: "string"},
				},
				Required: []string{"message"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolListTracts,
				Description: "List all tracts owned by the caller",
				Properties:  map[string]domain.ToolProperty{},
				Required:    []string{},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"uuid":        {Type: "string", Description: "Tract uuid (each array entry)"},
					"name":        {Type: "string"},
					"description": {Type: "string"},
					"enabled":     {Type: "boolean"},
					"created_at":  {Type: "string"},
					"updated_at":  {Type: "string"},
				},
				Required: []string{"uuid", "name", "enabled", "created_at", "updated_at"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolGetTract,
				Description: "Get one tract's full definition plus its linked triggers",
				Properties: map[string]domain.ToolProperty{
					"tract_uuid": {Type: "string"},
				},
				Required: []string{"tract_uuid"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"tract": {
						Type: "object",
						Properties: map[string]domain.ToolProperty{
							"uuid":        {Type: "string"},
							"name":        {Type: "string"},
							"description": {Type: "string"},
							"enabled":     {Type: "boolean"},
							"definition":  {Type: "object"},
							"created_at":  {Type: "string"},
							"updated_at":  {Type: "string"},
						},
					},
					"triggers": {
						Type: "array",
						Items: &domain.ToolProperty{
							Type: "object",
							Properties: map[string]domain.ToolProperty{
								"trigger_uuid": {Type: "string"},
								"name":         {Type: "string"},
								"kind":         {Type: "string"},
								"source":       {Type: "string"},
								"enabled":      {Type: "boolean"},
								"filters":      {Type: "array", Items: &filterItemSchema},
							},
						},
					},
				},
				Required: []string{"tract", "triggers"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolGetTractRuns,
				Description: "Get a tract's most recent runs (most recent first), each with its ordered step rows, for debugging a run",
				Properties: map[string]domain.ToolProperty{
					"tract_uuid": {Type: "string"},
					"limit":      {Type: "integer", Description: "Max runs to return (default 20)"},
				},
				Required: []string{"tract_uuid"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"uuid":            {Type: "string", Description: "Run uuid (each array entry)"},
					"status":          {Type: "string", Enum: []string{"running", "done", "failed"}},
					"started_by":      {Type: "string", Enum: []string{"webhook", "manual", "mcp"}},
					"trigger_payload": {Type: "object"},
					"error":           {Type: "string"},
					"created_at":      {Type: "string"},
					"updated_at":      {Type: "string"},
					"steps": {
						Type: "array",
						Items: &domain.ToolProperty{
							Type: "object",
							Properties: map[string]domain.ToolProperty{
								"step_id":     {Type: "string"},
								"step_name":   {Type: "string"},
								"step_type":   {Type: "string"},
								"status":      {Type: "string", Enum: []string{"running", "done", "failed"}},
								"input":       {Type: "object"},
								"output":      {Type: "object"},
								"error":       {Type: "string"},
								"started_at":  {Type: "string"},
								"finished_at": {Type: "string"},
							},
						},
					},
				},
				Required: []string{"uuid", "status", "started_by", "created_at", "updated_at", "steps"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolRunTract,
				Description: "Manually start a run of a tract, optionally seeding it with trigger-payload-shaped params",
				Properties: map[string]domain.ToolProperty{
					"tract_uuid": {Type: "string"},
					"params":     {Type: "object", Description: "Optional params, used as this run's trigger payload"},
				},
				Required: []string{"tract_uuid"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"message":    {Type: "string"},
					"tract_uuid": {Type: "string"},
				},
				Required: []string{"message", "tract_uuid"},
			},
		},
	}
}
