package tract

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func actionStep(id string, params map[string]string) domain.TractStep {
	step := domain.TractStep{
		Id:     id,
		Name:   id,
		Type:   stepTypeAction,
		Mcp:    builtinMcpName,
		Tool:   "write_file",
		Params: params,
	}

	return step
}

func TestValidateShape_StepIdRules(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"valid lowercase", "list_prs", false},
		{"valid single char", "a", false},
		{"valid with digits", "step_2", false},
		{"reserved trigger", "trigger", true},
		{"uppercase invalid", "ListPrs", true},
		{"starts with digit", "2step", true},
		{"starts with underscore", "_step", true},
		{"contains dash", "list-prs", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := actionStep(tt.id, nil)
			def := domain.TractDefinition{Steps: []domain.TractStep{step}}

			err := validateShape(def)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateShape_DuplicateIds(t *testing.T) {
	stepA := actionStep("dup", nil)
	stepB := actionStep("dup", nil)
	def := domain.TractDefinition{Steps: []domain.TractStep{stepA, stepB}}

	err := validateShape(def)
	assert.Error(t, err)
}

func TestValidateShape_DuplicateAcrossNesting(t *testing.T) {
	inner := actionStep("shared", nil)
	condStep := domain.TractStep{
		Id:         "cond",
		Type:       stepTypeCondition,
		Conditions: []domain.TractCondition{{Left: "1", Op: "==", Right: "1"}},
		Then:       []domain.TractStep{inner},
	}
	outer := actionStep("shared", nil)
	def := domain.TractDefinition{Steps: []domain.TractStep{condStep, outer}}

	err := validateShape(def)
	assert.Error(t, err)
}

func TestValidateShape_FieldUsageByType(t *testing.T) {
	t.Run("action with conditions is rejected", func(t *testing.T) {
		step := actionStep("a", nil)
		step.Conditions = []domain.TractCondition{{Left: "1", Op: "==", Right: "1"}}
		def := domain.TractDefinition{Steps: []domain.TractStep{step}}

		err := validateShape(def)
		assert.Error(t, err)
	})

	t.Run("condition with params is rejected", func(t *testing.T) {
		step := domain.TractStep{
			Id:         "c",
			Type:       stepTypeCondition,
			Conditions: []domain.TractCondition{{Left: "1", Op: "==", Right: "1"}},
			Params:     map[string]string{"x": "y"},
		}
		def := domain.TractDefinition{Steps: []domain.TractStep{step}}

		err := validateShape(def)
		assert.Error(t, err)
	})

	t.Run("condition without conditions is rejected", func(t *testing.T) {
		step := domain.TractStep{Id: "c", Type: stepTypeCondition}
		def := domain.TractDefinition{Steps: []domain.TractStep{step}}

		err := validateShape(def)
		assert.Error(t, err)
	})

	t.Run("parallel with params is rejected", func(t *testing.T) {
		step := domain.TractStep{
			Id:     "p",
			Type:   stepTypeParallel,
			Params: map[string]string{"x": "y"},
			Steps:  []domain.TractStep{actionStep("lane1", nil)},
		}
		def := domain.TractDefinition{Steps: []domain.TractStep{step}}

		err := validateShape(def)
		assert.Error(t, err)
	})

	t.Run("parallel with no children is rejected", func(t *testing.T) {
		step := domain.TractStep{Id: "p", Type: stepTypeParallel}
		def := domain.TractDefinition{Steps: []domain.TractStep{step}}

		err := validateShape(def)
		assert.Error(t, err)
	})

	t.Run("unknown step type is rejected", func(t *testing.T) {
		step := domain.TractStep{Id: "u", Type: "delay"}
		def := domain.TractDefinition{Steps: []domain.TractStep{step}}

		err := validateShape(def)
		assert.Error(t, err)
	})

	t.Run("valid group nests fine", func(t *testing.T) {
		group := domain.TractStep{
			Id:    "g",
			Type:  stepTypeGroup,
			Steps: []domain.TractStep{actionStep("inner", nil)},
		}
		def := domain.TractDefinition{Steps: []domain.TractStep{group}}

		err := validateShape(def)
		assert.NoError(t, err)
	})
}

func TestValidateShape_VisibilityRule(t *testing.T) {
	t.Run("trigger ref always visible", func(t *testing.T) {
		step := actionStep("a", map[string]string{"x": "{{ trigger.branch }}"})
		def := domain.TractDefinition{Steps: []domain.TractStep{step}}

		err := validateShape(def)
		assert.NoError(t, err)
	})

	t.Run("ref to preceding sibling is visible", func(t *testing.T) {
		first := actionStep("first", nil)
		second := actionStep("second", map[string]string{"x": "{{ first.field }}"})
		def := domain.TractDefinition{Steps: []domain.TractStep{first, second}}

		err := validateShape(def)
		assert.NoError(t, err)
	})

	t.Run("ref to a later sibling is rejected", func(t *testing.T) {
		first := actionStep("first", map[string]string{"x": "{{ second.field }}"})
		second := actionStep("second", nil)
		def := domain.TractDefinition{Steps: []domain.TractStep{first, second}}

		err := validateShape(def)
		assert.Error(t, err)
	})

	t.Run("ref to sibling in another parallel lane is rejected", func(t *testing.T) {
		laneA := actionStep("lane_a", nil)
		laneB := actionStep("lane_b", map[string]string{"x": "{{ lane_a.field }}"})
		parallelStep := domain.TractStep{
			Id:    "par",
			Type:  stepTypeParallel,
			Steps: []domain.TractStep{laneA, laneB},
		}
		def := domain.TractDefinition{Steps: []domain.TractStep{parallelStep}}

		err := validateShape(def)
		assert.Error(t, err)
	})

	t.Run("ref to parallel step output after it completes is visible", func(t *testing.T) {
		laneA := actionStep("lane_a", nil)
		parallelStep := domain.TractStep{
			Id:    "par",
			Type:  stepTypeParallel,
			Steps: []domain.TractStep{laneA},
		}
		after := actionStep("after", map[string]string{"x": "{{ lane_a.field }}"})
		def := domain.TractDefinition{Steps: []domain.TractStep{parallelStep, after}}

		err := validateShape(def)
		assert.NoError(t, err)
	})

	t.Run("ref from outside into a condition branch is rejected", func(t *testing.T) {
		inner := actionStep("inner", nil)
		condStep := domain.TractStep{
			Id:         "cond",
			Type:       stepTypeCondition,
			Conditions: []domain.TractCondition{{Left: "1", Op: "==", Right: "1"}},
			Then:       []domain.TractStep{inner},
		}
		after := actionStep("after", map[string]string{"x": "{{ inner.field }}"})
		def := domain.TractDefinition{Steps: []domain.TractStep{condStep, after}}

		err := validateShape(def)
		assert.Error(t, err)
	})

	t.Run("ref to the condition step's own output after it is visible", func(t *testing.T) {
		condStep := domain.TractStep{
			Id:         "cond",
			Type:       stepTypeCondition,
			Conditions: []domain.TractCondition{{Left: "1", Op: "==", Right: "1"}},
		}
		after := actionStep("after", map[string]string{"x": "{{ cond.result }}"})
		def := domain.TractDefinition{Steps: []domain.TractStep{condStep, after}}

		err := validateShape(def)
		assert.NoError(t, err)
	})

	t.Run("then branch cannot see else branch and vice versa", func(t *testing.T) {
		elseStep := actionStep("else_step", nil)
		thenStep := actionStep("then_step", map[string]string{"x": "{{ else_step.field }}"})
		condStep := domain.TractStep{
			Id:         "cond",
			Type:       stepTypeCondition,
			Conditions: []domain.TractCondition{{Left: "1", Op: "==", Right: "1"}},
			Then:       []domain.TractStep{thenStep},
			Else:       []domain.TractStep{elseStep},
		}
		def := domain.TractDefinition{Steps: []domain.TractStep{condStep}}

		err := validateShape(def)
		assert.Error(t, err)
	})

	t.Run("group fans in its own ids to later siblings", func(t *testing.T) {
		inner := actionStep("inner", nil)
		group := domain.TractStep{
			Id:    "g",
			Type:  stepTypeGroup,
			Steps: []domain.TractStep{inner},
		}
		after := actionStep("after", map[string]string{"x": "{{ inner.field }}"})
		def := domain.TractDefinition{Steps: []domain.TractStep{group, after}}

		err := validateShape(def)
		assert.NoError(t, err)
	})
}

func newTestService(mcpDefs *fakeMcpDefsRepo, externalConns *fakeExternalConnsRepo, executor *fakeToolExecutor) *Service {
	svc := New(nil, nil, externalConns, mcpDefs, executor)

	return svc
}

func TestValidateActionTool_BuiltinRules(t *testing.T) {
	mcpDefs := newFakeMcpDefsRepo()
	externalConns := newFakeExternalConnsRepo()
	executor := newFakeToolExecutor("write_file")
	svc := newTestService(mcpDefs, externalConns, executor)

	ownerUuid := uuid.New()
	ctx := context.Background()

	t.Run("known builtin without connection is valid", func(t *testing.T) {
		step := actionStep("a", nil)
		err := svc.validateActionTool(ctx, ownerUuid, step)
		assert.NoError(t, err)
	})

	t.Run("unknown builtin is rejected", func(t *testing.T) {
		step := actionStep("a", nil)
		step.Tool = "not_a_real_tool"
		err := svc.validateActionTool(ctx, ownerUuid, step)
		assert.Error(t, err)
	})

	t.Run("connection_id on builtin is rejected", func(t *testing.T) {
		step := actionStep("a", nil)
		step.ConnectionUuid = uuid.New()
		err := svc.validateActionTool(ctx, ownerUuid, step)
		assert.Error(t, err)
	})
}

func TestValidateActionTool_MomRules(t *testing.T) {
	mcpDefs := newFakeMcpDefsRepo()
	mcpDefs.addTool("gitlab", "create_merge_request")

	externalConns := newFakeExternalConnsRepo()
	executor := newFakeToolExecutor("write_file")
	svc := newTestService(mcpDefs, externalConns, executor)

	ownerUuid := uuid.New()
	otherUserUuid := uuid.New()
	connUuid := uuid.New()
	conn := domain.ExternalConnection{Uuid: connUuid, UserUuid: ownerUuid}
	externalConns.conns[connUuid] = conn

	ctx := context.Background()

	t.Run("unknown mom tool is rejected", func(t *testing.T) {
		step := domain.TractStep{Id: "a", Type: stepTypeAction, Mcp: "gitlab", Tool: "does_not_exist", ConnectionUuid: connUuid}
		err := svc.validateActionTool(ctx, ownerUuid, step)
		assert.Error(t, err)
	})

	t.Run("mom tool without connection is rejected", func(t *testing.T) {
		step := domain.TractStep{Id: "a", Type: stepTypeAction, Mcp: "gitlab", Tool: "create_merge_request"}
		err := svc.validateActionTool(ctx, ownerUuid, step)
		assert.Error(t, err)
	})

	t.Run("mom tool with owned connection is valid", func(t *testing.T) {
		step := domain.TractStep{Id: "a", Type: stepTypeAction, Mcp: "gitlab", Tool: "create_merge_request", ConnectionUuid: connUuid}
		err := svc.validateActionTool(ctx, ownerUuid, step)
		assert.NoError(t, err)
	})

	t.Run("mom tool with connection owned by someone else is rejected", func(t *testing.T) {
		step := domain.TractStep{Id: "a", Type: stepTypeAction, Mcp: "gitlab", Tool: "create_merge_request", ConnectionUuid: connUuid}
		err := svc.validateActionTool(ctx, otherUserUuid, step)
		assert.Error(t, err)
	})
}

func TestCheckTriggerFieldWarnings(t *testing.T) {
	step := actionStep("a", map[string]string{"x": "{{ trigger.branch }}", "y": "{{ trigger.unknown_field }}"})
	def := domain.TractDefinition{Steps: []domain.TractStep{step}}

	t.Run("no schemas means no warnings", func(t *testing.T) {
		warnings := checkTriggerFieldWarnings(def, nil)
		assert.Empty(t, warnings)
	})

	t.Run("field present in schema produces no warning for it", func(t *testing.T) {
		schema := domain.ToolSchema{Properties: map[string]domain.ToolProperty{
			"branch": {Type: "string"},
		}}
		warnings := checkTriggerFieldWarnings(def, []domain.ToolSchema{schema})
		require.Len(t, warnings, 1)
		assert.Contains(t, warnings[0], "unknown_field")
	})

	t.Run("field present in any linked schema is enough", func(t *testing.T) {
		schemaA := domain.ToolSchema{Properties: map[string]domain.ToolProperty{"branch": {Type: "string"}}}
		schemaB := domain.ToolSchema{Properties: map[string]domain.ToolProperty{"unknown_field": {Type: "string"}}}
		warnings := checkTriggerFieldWarnings(def, []domain.ToolSchema{schemaA, schemaB})
		assert.Empty(t, warnings)
	})
}
