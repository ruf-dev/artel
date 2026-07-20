package tract

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"go.redsock.ru/rerrors"
)

var (
	errFakeToolFailed   = errors.New("fake tool failed")
	errFakeToolTimedOut = errors.New("fake tool timed out waiting for cancellation")
)

// fakeTractsRepo is an in-memory repository.TractsRepo used by engine tests to observe
// run/run-step persistence without a real database.
type fakeTractsRepo struct {
	mu    sync.Mutex
	runs  map[uuid.UUID]domain.TractRun
	steps map[uuid.UUID]domain.TractRunStep
}

func newFakeTractsRepo() *fakeTractsRepo {
	repo := &fakeTractsRepo{
		runs:  map[uuid.UUID]domain.TractRun{},
		steps: map[uuid.UUID]domain.TractRunStep{},
	}

	return repo
}

func (f *fakeTractsRepo) Create(_ context.Context, tract domain.Tract) (domain.Tract, error) {
	return tract, nil
}

func (f *fakeTractsRepo) Get(_ context.Context, _ uuid.UUID) (sql.Null[domain.Tract], error) {
	return sql.Null[domain.Tract]{}, nil
}

func (f *fakeTractsRepo) ListByUser(_ context.Context, _ uuid.UUID) ([]domain.Tract, error) {
	return nil, nil
}

func (f *fakeTractsRepo) Update(_ context.Context, tract domain.Tract) (domain.Tract, error) {
	return tract, nil
}

func (f *fakeTractsRepo) SetEnabled(_ context.Context, _ uuid.UUID, _ bool) error {
	return nil
}

func (f *fakeTractsRepo) Delete(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (f *fakeTractsRepo) InsertRun(_ context.Context, run domain.TractRun) (domain.TractRun, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	run.Uuid = uuid.New()
	f.runs[run.Uuid] = run

	return run, nil
}

func (f *fakeTractsRepo) UpdateRunStatus(
	_ context.Context,
	id uuid.UUID,
	status domain.TractRunStatus,
	errMsg string,
) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	run := f.runs[id]
	run.Status = status
	run.Error = errMsg
	f.runs[id] = run

	return nil
}

func (f *fakeTractsRepo) GetRun(_ context.Context, id uuid.UUID) (sql.Null[domain.TractRun], error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	run, ok := f.runs[id]
	if !ok {
		return sql.Null[domain.TractRun]{}, nil
	}

	result := sql.Null[domain.TractRun]{V: run, Valid: true}

	return result, nil
}

func (f *fakeTractsRepo) ListRunsByTract(_ context.Context, _ uuid.UUID, _ int32) ([]domain.TractRun, error) {
	return nil, nil
}

func (f *fakeTractsRepo) SweepStaleRuns(_ context.Context, _ time.Time) error {
	return nil
}

func (f *fakeTractsRepo) InsertRunStep(_ context.Context, step domain.TractRunStep) (domain.TractRunStep, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	step.Uuid = uuid.New()
	f.steps[step.Uuid] = step

	return step, nil
}

func (f *fakeTractsRepo) UpdateRunStepFinish(
	_ context.Context,
	id uuid.UUID,
	status domain.TractRunStepStatus,
	output json.RawMessage,
	errMsg string,
) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	step := f.steps[id]
	step.Status = status
	step.Output = output
	step.Error = errMsg
	f.steps[id] = step

	return nil
}

func (f *fakeTractsRepo) ListRunStepsByRun(_ context.Context, runUuid uuid.UUID) ([]domain.TractRunStep, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var steps []domain.TractRunStep
	for _, step := range f.steps {
		if step.RunUuid == runUuid {
			steps = append(steps, step)
		}
	}

	return steps, nil
}

func (f *fakeTractsRepo) SweepStaleRunSteps(_ context.Context, _ time.Time) error {
	return nil
}

func (f *fakeTractsRepo) stepsByStepId(runUuid uuid.UUID) map[string]domain.TractRunStep {
	f.mu.Lock()
	defer f.mu.Unlock()

	result := map[string]domain.TractRunStep{}

	for _, step := range f.steps {
		if step.RunUuid == runUuid {
			result[step.StepId] = step
		}
	}

	return result
}

// fakeExternalConnsRepo is a minimal in-memory repository.ExternalConnectionRepo.
type fakeExternalConnsRepo struct {
	conns map[uuid.UUID]domain.ExternalConnection
}

func newFakeExternalConnsRepo() *fakeExternalConnsRepo {
	repo := &fakeExternalConnsRepo{conns: map[uuid.UUID]domain.ExternalConnection{}}

	return repo
}

func (f *fakeExternalConnsRepo) Upsert(
	_ context.Context,
	conn domain.ExternalConnection,
) (domain.ExternalConnection, error) {
	f.conns[conn.Uuid] = conn

	return conn, nil
}

func (f *fakeExternalConnsRepo) Insert(
	_ context.Context,
	conn domain.ExternalConnection,
) (domain.ExternalConnection, error) {
	f.conns[conn.Uuid] = conn

	return conn, nil
}

func (f *fakeExternalConnsRepo) GetByID(_ context.Context, id uuid.UUID) (domain.ExternalConnection, error) {
	conn, ok := f.conns[id]
	if !ok {
		return domain.ExternalConnection{}, sql.ErrNoRows
	}

	return conn, nil
}

func (f *fakeExternalConnsRepo) GetByUserAndProvider(
	_ context.Context,
	userUuid uuid.UUID,
	provider string,
) (sql.Null[domain.ExternalConnection], error) {
	for _, conn := range f.conns {
		if conn.UserUuid == userUuid && conn.Provider == provider {
			result := sql.Null[domain.ExternalConnection]{V: conn, Valid: true}

			return result, nil
		}
	}

	return sql.Null[domain.ExternalConnection]{}, nil
}

func (f *fakeExternalConnsRepo) ListByUser(_ context.Context, _ uuid.UUID) ([]domain.ExternalConnection, error) {
	return nil, nil
}

func (f *fakeExternalConnsRepo) Delete(_ context.Context, _ uuid.UUID, _ string) error {
	return nil
}

func (f *fakeExternalConnsRepo) DeleteByID(_ context.Context, _ uuid.UUID, id uuid.UUID) error {
	delete(f.conns, id)

	return nil
}

// fakeMcpDefsRepo is a minimal in-memory repository.McpDefinitionsRepo, controllable via
// tools (mcpName -> toolName -> exists).
type fakeMcpDefsRepo struct {
	tools map[string]map[string]bool
}

func newFakeMcpDefsRepo() *fakeMcpDefsRepo {
	repo := &fakeMcpDefsRepo{tools: map[string]map[string]bool{}}

	return repo
}

func (f *fakeMcpDefsRepo) addTool(mcpName string, toolName string) {
	if f.tools[mcpName] == nil {
		f.tools[mcpName] = map[string]bool{}
	}

	f.tools[mcpName][toolName] = true
}

func (f *fakeMcpDefsRepo) Upsert(_ context.Context, def domain.McpDefinition) (domain.McpDefinition, error) {
	return def, nil
}

func (f *fakeMcpDefsRepo) Get(_ context.Context, _ string) (sql.Null[domain.McpDefinition], error) {
	return sql.Null[domain.McpDefinition]{}, nil
}

func (f *fakeMcpDefsRepo) List(_ context.Context) ([]domain.McpDefinition, error) {
	return nil, nil
}

func (f *fakeMcpDefsRepo) Delete(_ context.Context, _ string) error {
	return nil
}

func (f *fakeMcpDefsRepo) GetTool(
	_ context.Context,
	mcpName string,
	toolName string,
) (sql.Null[domain.McpToolDef], error) {
	if !f.tools[mcpName][toolName] {
		return sql.Null[domain.McpToolDef]{}, nil
	}

	tool := domain.McpToolDef{ApiDescription: domain.ToolApiDescription{Name: toolName}}
	result := sql.Null[domain.McpToolDef]{V: tool, Valid: true}

	return result, nil
}

func (f *fakeMcpDefsRepo) ListAllTools(_ context.Context) ([]domain.McpToolRef, error) {
	return nil, nil
}

// fakeTriggersRepo is a minimal in-memory repository.TriggersRepo used by trigger tests.
type fakeTriggersRepo struct {
	mu             sync.Mutex
	triggers       map[uuid.UUID]domain.Trigger
	providerLinks  map[uuid.UUID]uuid.UUID // triggerId -> externalConnectionId
	linksByTrigger map[uuid.UUID][]repository.TractTriggerLink
}

func newFakeTriggersRepo() *fakeTriggersRepo {
	repo := &fakeTriggersRepo{
		triggers:       map[uuid.UUID]domain.Trigger{},
		providerLinks:  map[uuid.UUID]uuid.UUID{},
		linksByTrigger: map[uuid.UUID][]repository.TractTriggerLink{},
	}

	return repo
}

func (f *fakeTriggersRepo) Create(_ context.Context, trigger domain.Trigger) (domain.Trigger, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	trigger.Uuid = uuid.New()
	trigger.TriggerUuid = uuid.New()
	f.triggers[trigger.Uuid] = trigger

	return trigger, nil
}

func (f *fakeTriggersRepo) Get(_ context.Context, id uuid.UUID) (sql.Null[domain.Trigger], error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	trigger, ok := f.triggers[id]
	if !ok {
		return sql.Null[domain.Trigger]{}, nil
	}

	result := sql.Null[domain.Trigger]{V: trigger, Valid: true}

	return result, nil
}

func (f *fakeTriggersRepo) GetByTriggerUuid(_ context.Context, _ uuid.UUID) (sql.Null[domain.Trigger], error) {
	return sql.Null[domain.Trigger]{}, nil
}

func (f *fakeTriggersRepo) ListByUser(_ context.Context, _ uuid.UUID) ([]domain.Trigger, error) {
	return nil, nil
}

func (f *fakeTriggersRepo) SetEnabled(_ context.Context, _ uuid.UUID, _ bool) error {
	return nil
}

func (f *fakeTriggersRepo) Delete(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (f *fakeTriggersRepo) RotateSecret(
	_ context.Context, _ uuid.UUID, _ uuid.UUID, _ []byte, _ string,
) (domain.Trigger, error) {
	return domain.Trigger{}, nil
}

func (f *fakeTriggersRepo) Link(_ context.Context, _ domain.TriggerTractLink) error {
	return nil
}

func (f *fakeTriggersRepo) Unlink(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	return nil
}

func (f *fakeTriggersRepo) ListLinksByTract(_ context.Context, _ uuid.UUID) ([]repository.TractTriggerLink, error) {
	return nil, nil
}

func (f *fakeTriggersRepo) ListLinksByTrigger(_ context.Context, triggerUuid uuid.UUID) ([]repository.TractTriggerLink, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.linksByTrigger[triggerUuid], nil
}

func (f *fakeTriggersRepo) LinkToProvider(_ context.Context, triggerId uuid.UUID, externalConnectionId uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.providerLinks[triggerId] = externalConnectionId

	return nil
}

func (f *fakeTriggersRepo) ListByExternalConnection(_ context.Context, externalConnectionId uuid.UUID) ([]domain.Trigger, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var triggers []domain.Trigger

	for triggerId, connId := range f.providerLinks {
		if connId == externalConnectionId {
			triggers = append(triggers, f.triggers[triggerId])
		}
	}

	return triggers, nil
}

// fakeTriggerPresetsRepo is a minimal in-memory repository.TriggerPresetsRepo.
type fakeTriggerPresetsRepo struct {
	presets map[string]domain.TriggerPreset
}

func newFakeTriggerPresetsRepo() *fakeTriggerPresetsRepo {
	repo := &fakeTriggerPresetsRepo{presets: map[string]domain.TriggerPreset{}}

	return repo
}

func (f *fakeTriggerPresetsRepo) addPreset(preset domain.TriggerPreset) {
	f.presets[preset.Key] = preset
}

func (f *fakeTriggerPresetsRepo) List(_ context.Context) ([]domain.TriggerPreset, error) {
	presets := make([]domain.TriggerPreset, 0, len(f.presets))
	for _, preset := range f.presets {
		presets = append(presets, preset)
	}

	return presets, nil
}

func (f *fakeTriggerPresetsRepo) GetByKey(_ context.Context, key string) (sql.Null[domain.TriggerPreset], error) {
	preset, ok := f.presets[key]
	if !ok {
		return sql.Null[domain.TriggerPreset]{}, nil
	}

	result := sql.Null[domain.TriggerPreset]{V: preset, Valid: true}

	return result, nil
}

// fakeTractTemplatesRepo is a minimal in-memory repository.TractTemplatesRepo used by
// publish/unpublish/list/get/instantiate template tests.
type fakeTractTemplatesRepo struct {
	mu            sync.Mutex
	templates     map[uuid.UUID]domain.TractTemplate
	installCounts map[uuid.UUID]int
}

func newFakeTractTemplatesRepo() *fakeTractTemplatesRepo {
	repo := &fakeTractTemplatesRepo{
		templates:     map[uuid.UUID]domain.TractTemplate{},
		installCounts: map[uuid.UUID]int{},
	}

	return repo
}

func (f *fakeTractTemplatesRepo) Create(_ context.Context, template domain.TractTemplate) (domain.TractTemplate, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	template.Uuid = uuid.New()
	f.templates[template.Uuid] = template

	return template, nil
}

func (f *fakeTractTemplatesRepo) Get(_ context.Context, id uuid.UUID) (sql.Null[domain.TractTemplate], error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	template, ok := f.templates[id]
	if !ok {
		return sql.Null[domain.TractTemplate]{}, nil
	}

	template.InstallCount = f.installCounts[id]
	result := sql.Null[domain.TractTemplate]{V: template, Valid: true}

	return result, nil
}

func (f *fakeTractTemplatesRepo) ListAll(_ context.Context, category string) ([]domain.TractTemplate, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var templates []domain.TractTemplate

	for _, template := range f.templates {
		if category == "" || template.Category == category {
			templates = append(templates, template)
		}
	}

	return templates, nil
}

func (f *fakeTractTemplatesRepo) ListByOwner(_ context.Context, ownerUuid uuid.UUID) ([]domain.TractTemplate, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var templates []domain.TractTemplate

	for _, template := range f.templates {
		if template.OwnerUuid == ownerUuid {
			templates = append(templates, template)
		}
	}

	return templates, nil
}

func (f *fakeTractTemplatesRepo) Delete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.templates, id)

	return nil
}

func (f *fakeTractTemplatesRepo) IncrementInstallCount(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.installCounts[id]++

	return nil
}

var (
	_ repository.TractsRepo             = (*fakeTractsRepo)(nil)
	_ repository.TractTemplatesRepo     = (*fakeTractTemplatesRepo)(nil)
	_ repository.ExternalConnectionRepo = (*fakeExternalConnsRepo)(nil)
	_ repository.McpDefinitionsRepo     = (*fakeMcpDefsRepo)(nil)
	_ repository.TriggersRepo           = (*fakeTriggersRepo)(nil)
	_ repository.TriggerPresetsRepo     = (*fakeTriggerPresetsRepo)(nil)
)

// fakeToolExecutor is a controllable ToolExecutor for engine/validation tests.
type fakeToolExecutor struct {
	builtins map[string]bool

	mu    sync.Mutex
	calls []string
}

func newFakeToolExecutor(builtinNames ...string) *fakeToolExecutor {
	names := map[string]bool{}
	for _, name := range builtinNames {
		names[name] = true
	}

	executor := &fakeToolExecutor{builtins: names}

	return executor
}

func (f *fakeToolExecutor) IsBuiltinTool(toolName string) bool {
	return f.builtins[toolName]
}

func (f *fakeToolExecutor) recordCall(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.calls = append(f.calls, name)
}

func (f *fakeToolExecutor) ExecuteBuiltinTool(
	ctx context.Context,
	_ uuid.UUID,
	toolName string,
	params map[string]interface{},
) (string, error) {
	f.recordCall(toolName)

	return dispatchFakeTool(ctx, toolName, params)
}

func (f *fakeToolExecutor) ExecuteMomTool(
	ctx context.Context,
	_ uuid.UUID,
	_ string,
	toolName string,
	params map[string]interface{},
) (string, error) {
	f.recordCall(toolName)

	return dispatchFakeTool(ctx, toolName, params)
}

func (f *fakeToolExecutor) ListBuiltinTools(_ context.Context) ([]domain.McpToolDef, error) {
	return nil, nil
}

var _ ToolExecutor = (*fakeToolExecutor)(nil)

// dispatchFakeTool implements a few named test-tool behaviors shared by engine tests:
// "fail" errors immediately, "slow" blocks until ctx is cancelled (used to verify parallel
// cancellation), "raw_text" returns a non-JSON string (raw-wrap path), anything else returns
// a small JSON object echoing back its params.
func dispatchFakeTool(ctx context.Context, toolName string, params map[string]interface{}) (string, error) {
	switch toolName {
	case "fail":
		return "", errFakeToolFailed
	case "slow":
		select {
		case <-ctx.Done():
			return "", rerrors.Wrap(ctx.Err())
		case <-time.After(5 * time.Second):
			return "", errFakeToolTimedOut
		}
	case "raw_text":
		return "plain text result, not json", nil
	default:
		result := map[string]interface{}{"ok": true, "tool": toolName, "params": params}

		data, err := json.Marshal(result)
		if err != nil {
			return "", rerrors.Wrap(err, "failed to marshal result")
		}

		return string(data), nil
	}
}
