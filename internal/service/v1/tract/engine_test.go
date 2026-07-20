package tract

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/v1/subscription"
	"github.com/ruf-dev/artel/internal/service/v1/tract/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newEngineTestService(executor *fakeToolExecutor) (*Service, *fakeTractsRepo, *fakeExternalConnsRepo) {
	tracts := newFakeTractsRepo()
	templates := newFakeTractTemplatesRepo()
	externalConns := newFakeExternalConnsRepo()
	mcpDefs := newFakeMcpDefsRepo()
	scriptEngines := script.NewRegistry(script.NewJavaScriptEngine())
	svc := New(tracts, templates, nil, nil, externalConns, mcpDefs, executor, subscription.NewFree(), scriptEngines)

	return svc, tracts, externalConns
}

func testTract(userUuid uuid.UUID, steps []domain.TractStep) domain.Tract {
	def := domain.TractDefinition{Steps: steps}
	tract := domain.Tract{
		Uuid:       uuid.New(),
		UserUuid:   userUuid,
		Name:       "test tract",
		Enabled:    true,
		Definition: def,
	}

	return tract
}

func TestEngine_SimpleActionRun(t *testing.T) {
	executor := newFakeToolExecutor("write_file")
	svc, tracts, _ := newEngineTestService(executor)

	userUuid := uuid.New()
	step := actionStep("write", map[string]string{"path": "note.md"})
	tract := testTract(userUuid, []domain.TractStep{step})

	run, err := svc.StartRun(
		context.Background(),
		tract,
		json.RawMessage(`{"branch":"main"}`),
		StartedByManual,
		uuid.Nil,
	)
	require.NoError(t, err)
	assert.Equal(t, domain.TractRunDone, run.Status)

	steps := tracts.stepsByStepId(run.Uuid)
	require.Contains(t, steps, "write")
	assert.Equal(t, domain.TractRunStepDone, steps["write"].Status)
}

func TestEngine_ConditionThenElse(t *testing.T) {
	executor := newFakeToolExecutor("write_file")

	t.Run("then branch taken", func(t *testing.T) {
		svc, tracts, _ := newEngineTestService(executor)
		userUuid := uuid.New()

		thenStep := actionStep("then_action", nil)
		elseStep := actionStep("else_action", nil)
		condStep := domain.TractStep{
			Id:         "cond",
			Type:       stepTypeCondition,
			Conditions: []domain.TractCondition{{Left: "{{ trigger.flag }}", Op: "==", Right: "true"}},
			Then:       []domain.TractStep{thenStep},
			Else:       []domain.TractStep{elseStep},
		}
		tract := testTract(userUuid, []domain.TractStep{condStep})

		run, err := svc.StartRun(
			context.Background(),
			tract,
			json.RawMessage(`{"flag":true}`),
			StartedByManual,
			uuid.Nil,
		)
		require.NoError(t, err)
		assert.Equal(t, domain.TractRunDone, run.Status)

		steps := tracts.stepsByStepId(run.Uuid)
		assert.Contains(t, steps, "then_action")
		assert.NotContains(t, steps, "else_action")

		var output map[string]interface{}

		err = json.Unmarshal(steps["cond"].Output, &output)
		require.NoError(t, err)
		assert.Equal(t, "then", output["branch"])
		assert.Equal(t, true, output["result"])
	})

	t.Run("else branch taken", func(t *testing.T) {
		svc, tracts, _ := newEngineTestService(executor)
		userUuid := uuid.New()

		thenStep := actionStep("then_action", nil)
		elseStep := actionStep("else_action", nil)
		condStep := domain.TractStep{
			Id:         "cond",
			Type:       stepTypeCondition,
			Conditions: []domain.TractCondition{{Left: "{{ trigger.flag }}", Op: "==", Right: "true"}},
			Then:       []domain.TractStep{thenStep},
			Else:       []domain.TractStep{elseStep},
		}
		tract := testTract(userUuid, []domain.TractStep{condStep})

		run, err := svc.StartRun(
			context.Background(),
			tract,
			json.RawMessage(`{"flag":false}`),
			StartedByManual,
			uuid.Nil,
		)
		require.NoError(t, err)
		assert.Equal(t, domain.TractRunDone, run.Status)

		steps := tracts.stepsByStepId(run.Uuid)
		assert.Contains(t, steps, "else_action")
		assert.NotContains(t, steps, "then_action")
	})
}

func TestEngine_GroupNesting(t *testing.T) {
	executor := newFakeToolExecutor("write_file")
	svc, tracts, _ := newEngineTestService(executor)
	userUuid := uuid.New()

	first := actionStep("first", nil)
	second := actionStep("second", map[string]string{"ref": "{{ first.ok }}"})
	group := domain.TractStep{
		Id:    "g",
		Type:  stepTypeGroup,
		Steps: []domain.TractStep{first, second},
	}
	tract := testTract(userUuid, []domain.TractStep{group})

	run, err := svc.StartRun(context.Background(), tract, json.RawMessage(`{}`), StartedByManual, uuid.Nil)
	require.NoError(t, err)
	assert.Equal(t, domain.TractRunDone, run.Status)

	steps := tracts.stepsByStepId(run.Uuid)
	require.Contains(t, steps, "first")
	require.Contains(t, steps, "second")

	var secondInput map[string]interface{}

	err = json.Unmarshal(steps["second"].Input, &secondInput)
	require.NoError(t, err)
	assert.Equal(t, true, secondInput["ref"])
}

func TestEngine_ParallelFailingLaneCancelsSiblings(t *testing.T) {
	executor := newFakeToolExecutor("write_file")
	svc, tracts, _ := newEngineTestService(executor)
	userUuid := uuid.New()

	failStep := domain.TractStep{Id: "fail_lane", Type: stepTypeAction, Mcp: builtinMcpName, Tool: "fail"}
	slowStep := domain.TractStep{Id: "slow_lane", Type: stepTypeAction, Mcp: builtinMcpName, Tool: "slow"}
	parallelStep := domain.TractStep{
		Id:    "par",
		Type:  stepTypeParallel,
		Steps: []domain.TractStep{failStep, slowStep},
	}
	tract := testTract(userUuid, []domain.TractStep{parallelStep})

	start := time.Now()
	run, err := svc.StartRun(context.Background(), tract, json.RawMessage(`{}`), StartedByManual, uuid.Nil)
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Equal(t, domain.TractRunFailed, run.Status)
	// The slow lane's fake tool would take 5s to time out on its own; cancellation from the
	// failing sibling must cut that short.
	assert.Less(t, elapsed, 2*time.Second)

	steps := tracts.stepsByStepId(run.Uuid)
	require.Contains(t, steps, "fail_lane")
	assert.Equal(t, domain.TractRunStepFailed, steps["fail_lane"].Status)
}

func TestEngine_OutputParsing(t *testing.T) {
	executor := newFakeToolExecutor("write_file")
	svc, tracts, _ := newEngineTestService(executor)
	userUuid := uuid.New()

	t.Run("valid json output stored as-is", func(t *testing.T) {
		step := actionStep("json_step", nil)
		tract := testTract(userUuid, []domain.TractStep{step})

		run, err := svc.StartRun(context.Background(), tract, json.RawMessage(`{}`), StartedByManual, uuid.Nil)
		require.NoError(t, err)

		steps := tracts.stepsByStepId(run.Uuid)

		var output map[string]interface{}

		err = json.Unmarshal(steps["json_step"].Output, &output)
		require.NoError(t, err)
		assert.Equal(t, true, output["ok"])
	})

	t.Run("non-json output wrapped as raw", func(t *testing.T) {
		step := domain.TractStep{Id: "raw_step", Type: stepTypeAction, Mcp: builtinMcpName, Tool: "raw_text"}
		tract := testTract(userUuid, []domain.TractStep{step})

		run, err := svc.StartRun(context.Background(), tract, json.RawMessage(`{}`), StartedByManual, uuid.Nil)
		require.NoError(t, err)

		steps := tracts.stepsByStepId(run.Uuid)

		var output map[string]interface{}

		err = json.Unmarshal(steps["raw_step"].Output, &output)
		require.NoError(t, err)
		assert.Equal(t, "plain text result, not json", output["raw"])
	})
}

func TestEngine_ActionFailureFailsRun(t *testing.T) {
	executor := newFakeToolExecutor("write_file")
	svc, tracts, _ := newEngineTestService(executor)
	userUuid := uuid.New()

	failStep := domain.TractStep{Id: "fail_step", Type: stepTypeAction, Mcp: builtinMcpName, Tool: "fail"}
	afterStep := actionStep("after", nil)
	tract := testTract(userUuid, []domain.TractStep{failStep, afterStep})

	run, err := svc.StartRun(context.Background(), tract, json.RawMessage(`{}`), StartedByManual, uuid.Nil)
	require.Error(t, err)
	assert.Equal(t, domain.TractRunFailed, run.Status)
	assert.NotEmpty(t, run.Error)

	steps := tracts.stepsByStepId(run.Uuid)
	assert.Equal(t, domain.TractRunStepFailed, steps["fail_step"].Status)
	assert.NotContains(t, steps, "after")
}

func TestEngine_ScriptStep(t *testing.T) {
	executor := newFakeToolExecutor("write_file")
	svc, tracts, _ := newEngineTestService(executor)
	userUuid := uuid.New()

	step := domain.TractStep{
		Id:       "sum_nums",
		Name:     "sum_nums",
		Type:     stepTypeScript,
		Language: domain.ScriptLanguageJavaScript,
		Code:     "for (let i = 0; i < nums.length; i++) { total += nums[i]; }",
		InputParams: []domain.ScriptParam{
			{Name: "nums", Property: domain.ToolProperty{Type: "array"}},
		},
		OutputParams: []domain.ScriptParam{
			{Name: "total", Property: domain.ToolProperty{Type: "number"}},
		},
		Params: map[string]string{"nums": "{{ trigger.nums }}"},
	}
	tract := testTract(userUuid, []domain.TractStep{step})

	run, err := svc.StartRun(context.Background(), tract, json.RawMessage(`{"nums":[1,2,3]}`), StartedByManual, uuid.Nil)
	require.NoError(t, err)
	assert.Equal(t, domain.TractRunDone, run.Status)

	steps := tracts.stepsByStepId(run.Uuid)
	require.Contains(t, steps, "sum_nums")
	assert.Equal(t, domain.TractRunStepDone, steps["sum_nums"].Status)

	var output map[string]interface{}

	err = json.Unmarshal(steps["sum_nums"].Output, &output)
	require.NoError(t, err)
	assert.Equal(t, float64(6), output["total"])
}

func TestEngine_ScriptStepOutputFeedsLaterStep(t *testing.T) {
	executor := newFakeToolExecutor("write_file")
	svc, tracts, _ := newEngineTestService(executor)
	userUuid := uuid.New()

	doubleStep := domain.TractStep{
		Id:       "double",
		Name:     "double",
		Type:     stepTypeScript,
		Language: domain.ScriptLanguageJavaScript,
		Code:     "result = n * 2;",
		InputParams: []domain.ScriptParam{
			{Name: "n", Property: domain.ToolProperty{Type: "number"}},
		},
		OutputParams: []domain.ScriptParam{
			{Name: "result", Property: domain.ToolProperty{Type: "number"}},
		},
		Params: map[string]string{"n": "{{ trigger.n }}"},
	}
	after := actionStep("after", map[string]string{"ref": "{{ double.result }}"})
	tract := testTract(userUuid, []domain.TractStep{doubleStep, after})

	run, err := svc.StartRun(context.Background(), tract, json.RawMessage(`{"n":21}`), StartedByManual, uuid.Nil)
	require.NoError(t, err)
	assert.Equal(t, domain.TractRunDone, run.Status)

	steps := tracts.stepsByStepId(run.Uuid)
	require.Contains(t, steps, "after")

	var afterInput map[string]interface{}

	err = json.Unmarshal(steps["after"].Input, &afterInput)
	require.NoError(t, err)
	assert.InDelta(t, 42, afterInput["ref"], 0.0001)
}

func TestEngine_MomToolOwnershipCheck(t *testing.T) {
	executor := newFakeToolExecutor()
	svc, tracts, externalConns := newEngineTestService(executor)

	ownerUuid := uuid.New()
	otherUserUuid := uuid.New()
	connUuid := uuid.New()
	conn := domain.ExternalConnection{Uuid: connUuid, UserUuid: otherUserUuid}
	externalConns.conns[connUuid] = conn

	step := domain.TractStep{
		Id:             "mom_step",
		Type:           stepTypeAction,
		Mcp:            "gitlab",
		Tool:           "create_merge_request",
		ConnectionUuid: connUuid,
	}
	tract := testTract(ownerUuid, []domain.TractStep{step})

	run, err := svc.StartRun(context.Background(), tract, json.RawMessage(`{}`), StartedByManual, uuid.Nil)
	require.Error(t, err)
	assert.Equal(t, domain.TractRunFailed, run.Status)

	steps := tracts.stepsByStepId(run.Uuid)
	assert.Equal(t, domain.TractRunStepFailed, steps["mom_step"].Status)
}

func TestEngine_MomToolExecutesWhenOwned(t *testing.T) {
	executor := newFakeToolExecutor()
	svc, tracts, externalConns := newEngineTestService(executor)

	ownerUuid := uuid.New()
	connUuid := uuid.New()
	conn := domain.ExternalConnection{Uuid: connUuid, UserUuid: ownerUuid}
	externalConns.conns[connUuid] = conn

	step := domain.TractStep{
		Id:             "mom_step",
		Type:           stepTypeAction,
		Mcp:            "gitlab",
		Tool:           "create_merge_request",
		ConnectionUuid: connUuid,
	}
	tract := testTract(ownerUuid, []domain.TractStep{step})

	run, err := svc.StartRun(context.Background(), tract, json.RawMessage(`{}`), StartedByManual, uuid.Nil)
	require.NoError(t, err)
	assert.Equal(t, domain.TractRunDone, run.Status)

	steps := tracts.stepsByStepId(run.Uuid)
	assert.Equal(t, domain.TractRunStepDone, steps["mom_step"].Status)
}
