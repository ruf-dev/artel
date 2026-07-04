package executors_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
)

const (
	testRenamedName    = "renamed"
	testDefinitionKey  = "definition"
	testNewDescription = "new description"
)

// fakeTractService implements executors.TractService with just enough behavior for
// update_tract tests; every other method panics if called.
type fakeTractService struct {
	updateTractFn func(
		ctx context.Context, id uuid.UUID, name string, description string, def domain.TractDefinition,
	) (domain.Tract, []string, error)
}

func (f *fakeTractService) CreateTract(
	_ context.Context,
	_ string,
	_ string,
	_ domain.TractDefinition,
) (domain.Tract, []string, error) {
	panic("not implemented")
}

func (f *fakeTractService) UpdateTract(
	ctx context.Context,
	id uuid.UUID,
	name string,
	description string,
	def domain.TractDefinition,
) (domain.Tract, []string, error) {
	return f.updateTractFn(ctx, id, name, description, def)
}

func (f *fakeTractService) GetTract(_ context.Context, _ uuid.UUID) (domain.Tract, error) {
	panic("not implemented")
}

func (f *fakeTractService) ListTracts(_ context.Context) ([]domain.Tract, error) {
	panic("not implemented")
}

func (f *fakeTractService) CreateTrigger(
	_ context.Context,
	_ string,
	_ string,
	_ string,
	_ json.RawMessage,
	_ domain.ToolSchema,
) (domain.Trigger, string, error) {
	panic("not implemented")
}

func (f *fakeTractService) LinkTrigger(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ []domain.TractCondition) error {
	panic("not implemented")
}

func (f *fakeTractService) ListLinksByTract(_ context.Context, _ uuid.UUID) ([]repository.TractTriggerLink, error) {
	panic("not implemented")
}

func (f *fakeTractService) StartRun(
	_ context.Context,
	_ domain.Tract,
	_ json.RawMessage,
	_ string,
	_ uuid.UUID,
) (domain.TractRun, error) {
	panic("not implemented")
}

func (f *fakeTractService) ListRuns(_ context.Context, _ uuid.UUID, _ int32) ([]domain.TractRun, error) {
	panic("not implemented")
}

func (f *fakeTractService) GetRun(_ context.Context, _ uuid.UUID) (domain.TractRun, []domain.TractRunStep, error) {
	panic("not implemented")
}

func (f *fakeTractService) ListTractTools(_ context.Context) ([]domain.McpToolRef, error) {
	panic("not implemented")
}

func (f *fakeTractService) ListTriggerSources(_ context.Context) ([]domain.TriggerSourcePreset, error) {
	panic("not implemented")
}

func TestUpdateTract_Success(t *testing.T) {
	tractUuid := uuid.New()
	def := domain.TractDefinition{Steps: []domain.TractStep{{Id: "step1", Type: "action"}}}

	updated := domain.Tract{
		Uuid:        tractUuid,
		Name:        testRenamedName,
		Description: testNewDescription,
		Enabled:     true,
		Definition:  def,
	}

	var gotId uuid.UUID

	var gotName string

	var gotDescription string

	var gotDef domain.TractDefinition

	ts := &fakeTractService{
		updateTractFn: func(
			_ context.Context, id uuid.UUID, name string, description string, d domain.TractDefinition,
		) (domain.Tract, []string, error) {
			gotId = id
			gotName = name
			gotDescription = description
			gotDef = d

			return updated, nil, nil
		},
	}

	e := executors.NewTractExecutor()

	defRaw, err := json.Marshal(def)
	if err != nil {
		t.Fatalf("error marshaling definition: %v", err)
	}

	var defParam map[string]interface{}

	err = json.Unmarshal(defRaw, &defParam)
	if err != nil {
		t.Fatalf("error unmarshaling definition param: %v", err)
	}

	params := map[string]interface{}{
		"tract_uuid":      tractUuid.String(),
		"name":            testRenamedName,
		"description":     testNewDescription,
		testDefinitionKey: defParam,
	}

	result, err := e.Execute(context.Background(), context.Background(), ts, executors.ToolUpdateTract, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotId != tractUuid {
		t.Fatalf("expected id %s, got %s", tractUuid, gotId)
	}

	if gotName != testRenamedName {
		t.Fatalf("expected name 'renamed', got %q", gotName)
	}

	if gotDescription != testNewDescription {
		t.Fatalf("expected description 'new description', got %q", gotDescription)
	}

	if len(gotDef.Steps) != 1 || gotDef.Steps[0].Id != "step1" {
		t.Fatalf("expected definition to round-trip with step1, got %+v", gotDef)
	}

	var decoded map[string]interface{}

	err = json.Unmarshal([]byte(result.Text), &decoded)
	if err != nil {
		t.Fatalf("error unmarshaling result: %v", err)
	}

	tractRowRaw, ok := decoded["tract"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected result to contain a tract object, got %+v", decoded)
	}

	if tractRowRaw["uuid"] != tractUuid.String() {
		t.Fatalf("expected result tract uuid %s, got %v", tractUuid, tractRowRaw["uuid"])
	}
}

func TestUpdateTract_MissingTractUuid(t *testing.T) {
	ts := &fakeTractService{}
	e := executors.NewTractExecutor()

	params := map[string]interface{}{
		"name":            testRenamedName,
		testDefinitionKey: map[string]interface{}{},
	}

	_, err := e.Execute(context.Background(), context.Background(), ts, executors.ToolUpdateTract, params)
	if err == nil {
		t.Fatalf("expected error for missing tract_uuid, got nil")
	}
}

func TestUpdateTract_MissingName(t *testing.T) {
	ts := &fakeTractService{}
	e := executors.NewTractExecutor()

	params := map[string]interface{}{
		"tract_uuid":      uuid.New().String(),
		testDefinitionKey: map[string]interface{}{},
	}

	_, err := e.Execute(context.Background(), context.Background(), ts, executors.ToolUpdateTract, params)
	if err == nil {
		t.Fatalf("expected error for missing name, got nil")
	}
}
