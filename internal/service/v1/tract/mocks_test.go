package tract

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
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

func (f *fakeTractsRepo) UpdateRunStatus(_ context.Context, id uuid.UUID, status domain.TractRunStatus, errMsg string) error {
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

func (f *fakeTractsRepo) UpdateRunStepFinish(_ context.Context, id uuid.UUID, status domain.TractRunStepStatus, output json.RawMessage, errMsg string) error {
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

func (f *fakeExternalConnsRepo) Upsert(_ context.Context, conn domain.ExternalConnection) (domain.ExternalConnection, error) {
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

func (f *fakeExternalConnsRepo) GetByUserAndProvider(_ context.Context, _ uuid.UUID, _ string) (sql.Null[domain.ExternalConnection], error) {
	return sql.Null[domain.ExternalConnection]{}, nil
}

func (f *fakeExternalConnsRepo) ListByUser(_ context.Context, _ uuid.UUID) ([]domain.ExternalConnection, error) {
	return nil, nil
}

func (f *fakeExternalConnsRepo) Delete(_ context.Context, _ uuid.UUID, _ string) error {
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

func (f *fakeMcpDefsRepo) GetTool(_ context.Context, mcpName string, toolName string) (sql.Null[domain.McpToolDef], error) {
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

var (
	_ repository.TractsRepo             = (*fakeTractsRepo)(nil)
	_ repository.ExternalConnectionRepo = (*fakeExternalConnsRepo)(nil)
	_ repository.McpDefinitionsRepo     = (*fakeMcpDefsRepo)(nil)
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

func (f *fakeToolExecutor) ExecuteBuiltinTool(ctx context.Context, _ uuid.UUID, toolName string, params map[string]interface{}) (string, error) {
	f.recordCall(toolName)

	return dispatchFakeTool(ctx, toolName, params)
}

func (f *fakeToolExecutor) ExecuteMomTool(ctx context.Context, _ uuid.UUID, _ string, toolName string, params map[string]interface{}) (string, error) {
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
