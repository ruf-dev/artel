package tract

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/tract/script"
	"go.redsock.ru/rerrors"
	"golang.org/x/sync/errgroup"
)

const (
	StartedByWebhook = "webhook"
	StartedByManual  = "manual"
	StartedByMcp     = "mcp"

	branchThen = "then"
	branchElse = "else"
)

// runState carries the parts of one run shared across the (possibly concurrent) step walk:
// the tract/run identity and the resolver context — trigger payload plus every step output
// recorded so far, keyed by step id. Guarded by mu since parallel steps write concurrently.
type runState struct {
	tract          domain.Tract
	run            domain.TractRun
	triggerPayload interface{}

	mu      sync.Mutex
	outputs map[string]interface{}
}

// snapshotResolver copies the current outputs into a resolver — a snapshot rather than a
// live view, so a step's template rendering never observes a sibling's output appearing
// mid-render.
func (rs *runState) snapshotResolver() *resolver {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	outputsCopy := make(map[string]interface{}, len(rs.outputs))
	for k, v := range rs.outputs {
		outputsCopy[k] = v
	}

	return &resolver{trigger: rs.triggerPayload, outputs: outputsCopy}
}

func (rs *runState) setOutput(stepId string, value interface{}) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.outputs[stepId] = value
}

// StartRun persists a tract_runs row (persist-before-apply) then walks tract.Definition.Steps.
// triggerUuid is uuid.Nil for manual/mcp runs. The caller decides whether to run this
// synchronously or as `go engine.StartRun(...)` against a server-lifecycle context — StartRun
// itself does not spawn a goroutine.
func (s *Service) StartRun(
	ctx context.Context,
	tract domain.Tract,
	payload json.RawMessage,
	startedBy string,
	triggerUuid uuid.UUID,
) (domain.TractRun, error) {
	err := s.subscriptions.CheckFeature(ctx, tract.UserUuid, domain.FeatureTract)
	if err != nil {
		return domain.TractRun{}, err
	}

	if len(payload) == 0 {
		payload = json.RawMessage(`{}`)
	}

	runInput := domain.TractRun{
		TractUuid:      tract.Uuid,
		TriggerUuid:    triggerUuid,
		StartedBy:      startedBy,
		TriggerPayload: payload,
	}

	run, err := s.tracts.InsertRun(ctx, runInput)
	if err != nil {
		return domain.TractRun{}, rerrors.Wrap(err, "error inserting tract run")
	}

	var triggerPayload interface{}

	err = json.Unmarshal(payload, &triggerPayload)
	if err != nil {
		triggerPayload = map[string]interface{}{}
	}

	state := &runState{
		tract:          tract,
		run:            run,
		triggerPayload: triggerPayload,
		outputs:        map[string]interface{}{},
	}

	runErr := s.executeSteps(ctx, state, tract.Definition.Steps)

	finalStatus := domain.TractRunDone
	errMsg := ""

	if runErr != nil {
		finalStatus = domain.TractRunFailed
		errMsg = runErr.Error()
		log.Error().
			Err(runErr).
			Str("run_uuid", run.Uuid.String()).
			Str("tract_uuid", tract.Uuid.String()).
			Msg("tract run failed")
	}

	err = s.tracts.UpdateRunStatus(ctx, run.Uuid, finalStatus, errMsg)
	if err != nil {
		log.Error().Err(err).Str("run_uuid", run.Uuid.String()).Msg("error updating tract run status")
	}

	run.Status = finalStatus
	run.Error = errMsg

	return run, runErr
}

func (s *Service) executeSteps(ctx context.Context, state *runState, steps []domain.TractStep) error {
	for _, step := range steps {
		err := s.executeStep(ctx, state, step)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) executeStep(ctx context.Context, state *runState, step domain.TractStep) error {
	select {
	case <-ctx.Done():
		return rerrors.Wrap(ctx.Err())
	default:
	}

	switch step.Type {
	case stepTypeAction:
		return s.executeAction(ctx, state, step)
	case stepTypeCondition:
		return s.executeCondition(ctx, state, step)
	case stepTypeParallel:
		return s.executeParallel(ctx, state, step)
	case stepTypeGroup:
		return s.executeSteps(ctx, state, step.Steps)
	case stepTypeScript:
		return s.executeScript(ctx, state, step)
	default:
		// Validation already rejects unknown step types; this is a defensive fallback.
		return rerrors.Wrap(user_errors.TractUnknownStepType, step.Type)
	}
}

func (s *Service) executeAction(ctx context.Context, state *runState, step domain.TractStep) error {
	rslv := state.snapshotResolver()

	rendered, err := renderParams(rslv, step.Params)
	if err != nil {
		return rerrors.Wrap(err, "error rendering action params: "+step.Id)
	}

	inputJSON, err := json.Marshal(rendered)
	if err != nil {
		return rerrors.Wrap(err, "error marshaling rendered params: "+step.Id)
	}

	runStep := domain.TractRunStep{
		RunUuid:  state.run.Uuid,
		StepId:   step.Id,
		StepName: step.Name,
		StepType: step.Type,
		Input:    inputJSON,
	}

	insertedStep, err := s.tracts.InsertRunStep(ctx, runStep)
	if err != nil {
		return rerrors.Wrap(err, "error inserting run step: "+step.Id)
	}

	resultText, err := s.executeTool(ctx, state.tract, step, rendered)
	if err != nil {
		finishErr := s.tracts.UpdateRunStepFinish(ctx, insertedStep.Uuid, domain.TractRunStepFailed, nil, err.Error())
		if finishErr != nil {
			log.Error().Err(finishErr).Str("step_id", step.Id).Msg("error updating failed run step")
		}

		return rerrors.Wrap(err, "error executing tool for step: "+step.Id)
	}

	output := parseToolOutput(resultText)

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return rerrors.Wrap(err, "error marshaling step output: "+step.Id)
	}

	err = s.tracts.UpdateRunStepFinish(ctx, insertedStep.Uuid, domain.TractRunStepDone, outputJSON, "")
	if err != nil {
		return rerrors.Wrap(err, "error finishing run step: "+step.Id)
	}

	state.setOutput(step.Id, output)

	return nil
}

// executeScript mirrors executeAction's plumbing (snapshot → render → record → run →
// finish → publish), swapping the MoM/builtin tool dispatch for the step's registered
// script.Engine. Engine presence for step.Language is already guaranteed by
// validateScriptEngines at save time.
func (s *Service) executeScript(ctx context.Context, state *runState, step domain.TractStep) error {
	rslv := state.snapshotResolver()

	rendered, err := renderParams(rslv, step.Params)
	if err != nil {
		return rerrors.Wrap(err, "error rendering script params: "+step.Id)
	}

	inputJSON, err := json.Marshal(rendered)
	if err != nil {
		return rerrors.Wrap(err, "error marshaling rendered params: "+step.Id)
	}

	runStep := domain.TractRunStep{
		RunUuid:  state.run.Uuid,
		StepId:   step.Id,
		StepName: step.Name,
		StepType: step.Type,
		Input:    inputJSON,
	}

	insertedStep, err := s.tracts.InsertRunStep(ctx, runStep)
	if err != nil {
		return rerrors.Wrap(err, "error inserting run step: "+step.Id)
	}

	runInput := script.RunInput{
		Body:         step.Code,
		InputParams:  step.InputParams,
		OutputParams: step.OutputParams,
		Args:         rendered,
	}

	output, err := s.scriptEngines[step.Language].Run(ctx, runInput)
	if err != nil {
		finishErr := s.tracts.UpdateRunStepFinish(ctx, insertedStep.Uuid, domain.TractRunStepFailed, nil, err.Error())
		if finishErr != nil {
			log.Error().Err(finishErr).Str("step_id", step.Id).Msg("error updating failed run step")
		}

		return rerrors.Wrap(err, "error executing script for step: "+step.Id)
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return rerrors.Wrap(err, "error marshaling step output: "+step.Id)
	}

	err = s.tracts.UpdateRunStepFinish(ctx, insertedStep.Uuid, domain.TractRunStepDone, outputJSON, "")
	if err != nil {
		return rerrors.Wrap(err, "error finishing run step: "+step.Id)
	}

	state.setOutput(step.Id, output)

	return nil
}

func renderParams(rslv *resolver, params map[string]string) (map[string]interface{}, error) {
	rendered := make(map[string]interface{}, len(params))

	for key, value := range params {
		result, err := rslv.render(value)
		if err != nil {
			return nil, rerrors.Wrap(err, "error rendering param: "+key)
		}

		rendered[key] = result
	}

	return rendered, nil
}

// executeTool dispatches to the builtin or MoM path. For MoM tools it re-verifies that the
// connection belongs to the tract owner — the tract owner is the principal for
// webhook-triggered runs, which have no user_context.
func (s *Service) executeTool(
	ctx context.Context,
	tract domain.Tract,
	step domain.TractStep,
	params map[string]interface{},
) (string, error) {
	if step.Mcp == builtinMcpName {
		res, err := s.toolExecutor.ExecuteBuiltinTool(ctx, tract.UserUuid, step.Tool, params)
		if err != nil {
			return "", rerrors.Wrap(err, "error executing builtin tool: "+step.Tool)
		}

		return res, nil
	}

	conn, err := s.externalConns.GetByID(ctx, step.ConnectionUuid)
	if err != nil {
		return "", rerrors.Wrap(err, "error getting external connection")
	}

	if conn.UserUuid != tract.UserUuid {
		return "", rerrors.Wrap(user_errors.TractConnectionNotOwned, step.Id)
	}

	res, err := s.toolExecutor.ExecuteMomTool(ctx, step.ConnectionUuid, step.Mcp, step.Tool, params)
	if err != nil {
		return "", rerrors.Wrap(err, "error executing mom tool: "+step.Id)
	}

	return res, nil
}

// parseToolOutput stores a tool's result as parsed JSON when it is valid JSON, or wraps it in
// {"raw": "<text>"} otherwise — either way the recorded output is always a JSON value,
// resolvable by later steps' templates.
func parseToolOutput(text string) interface{} {
	var parsed interface{}

	err := json.Unmarshal([]byte(text), &parsed)
	if err != nil {
		return map[string]interface{}{"raw": text}
	}

	return parsed
}

type renderedCondition struct {
	Left  interface{} `json:"left"`
	Op    string      `json:"op"`
	Right interface{} `json:"right"`
}

func (s *Service) executeCondition(ctx context.Context, state *runState, step domain.TractStep) error {
	rslv := state.snapshotResolver()

	rendered := make([]renderedCondition, len(step.Conditions))
	overallResult := true

	for i, cond := range step.Conditions {
		result, left, right, err := evaluate(cond, rslv)
		if err != nil {
			return rerrors.Wrap(err, "error evaluating condition: "+step.Id)
		}

		rendered[i] = renderedCondition{Left: left, Op: cond.Op, Right: right}

		if !result {
			overallResult = false
		}
	}

	inputJSON, err := json.Marshal(rendered)
	if err != nil {
		return rerrors.Wrap(err, "error marshaling condition input: "+step.Id)
	}

	branch := branchElse
	branchSteps := step.Else

	if overallResult {
		branch = branchThen
		branchSteps = step.Then
	}

	output := map[string]interface{}{"result": overallResult, fieldBranch: branch}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return rerrors.Wrap(err, "error marshaling condition output: "+step.Id)
	}

	runStep := domain.TractRunStep{
		RunUuid:  state.run.Uuid,
		StepId:   step.Id,
		StepName: step.Name,
		StepType: step.Type,
		Input:    inputJSON,
	}

	insertedStep, err := s.tracts.InsertRunStep(ctx, runStep)
	if err != nil {
		return rerrors.Wrap(err, "error inserting condition run step: "+step.Id)
	}

	err = s.tracts.UpdateRunStepFinish(ctx, insertedStep.Uuid, domain.TractRunStepDone, outputJSON, "")
	if err != nil {
		return rerrors.Wrap(err, "error finishing condition run step: "+step.Id)
	}

	state.setOutput(step.Id, output)

	return s.executeSteps(ctx, state, branchSteps)
}

// executeParallel runs each lane concurrently over an errgroup-derived context: the first
// lane to fail cancels the group's context, which executeStep observes at the top of every
// recursive call (and which propagates into any in-flight tool call built with that context).
func (s *Service) executeParallel(ctx context.Context, state *runState, step domain.TractStep) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for _, lane := range step.Steps {
		lane := lane

		group.Go(func() error {
			return s.executeStep(groupCtx, state, lane)
		})
	}

	err := group.Wait()
	if err != nil {
		return rerrors.Wrap(err, "error in parallel step: "+step.Id)
	}

	return nil
}
