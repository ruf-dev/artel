package tracts

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
	"github.com/sqlc-dev/pqtype"
	"go.redsock.ru/rerrors"
)

// Repo is the pure-DB layer for tracts + their runs/run-steps. The step tree itself is not a
// separate table — Tract.Definition is persisted as one JSONB column (see migration 033).
type Repo struct {
	q  *artel_q.Queries
	tx tx_manager.TxManager
}

func New(q *artel_q.Queries, tx tx_manager.TxManager) *Repo {
	return &Repo{q: q, tx: tx}
}

// --- tracts ---

func (r *Repo) Create(ctx context.Context, tract domain.Tract) (domain.Tract, error) {
	defJSON, err := json.Marshal(tract.Definition)
	if err != nil {
		return domain.Tract{}, rerrors.Wrap(err, "error marshaling tract definition")
	}

	params := artel_q.InsertTractParams{
		UserID:      tract.UserUuid,
		Name:        tract.Name,
		Description: tract.Description,
		Enabled:     tract.Enabled,
		Definition:  defJSON,
	}

	row, err := r.q.InsertTract(ctx, params)
	if err != nil {
		return domain.Tract{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error inserting tract")
	}

	return tractToDomain(row)
}

func (r *Repo) Get(ctx context.Context, id uuid.UUID) (sql.Null[domain.Tract], error) {
	row, err := r.q.GetTract(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.Tract]{}, nil
		}

		return sql.Null[domain.Tract]{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting tract")
	}

	tract, err := tractToDomain(row)
	if err != nil {
		return sql.Null[domain.Tract]{}, err
	}

	result := sql.Null[domain.Tract]{V: tract, Valid: true}

	return result, nil
}

func (r *Repo) ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.Tract, error) {
	rows, err := r.q.ListTractsByUser(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing tracts")
	}

	tracts := make([]domain.Tract, len(rows))

	for i, row := range rows {
		tract, convErr := tractToDomain(row)
		if convErr != nil {
			return nil, convErr
		}

		tracts[i] = tract
	}

	return tracts, nil
}

// Update overwrites name/description/definition — enabled has its own SetEnabled.
func (r *Repo) Update(ctx context.Context, tract domain.Tract) (domain.Tract, error) {
	defJSON, err := json.Marshal(tract.Definition)
	if err != nil {
		return domain.Tract{}, rerrors.Wrap(err, "error marshaling tract definition")
	}

	params := artel_q.UpdateTractParams{
		ID:          tract.Uuid,
		Name:        tract.Name,
		Description: tract.Description,
		Definition:  defJSON,
	}

	row, err := r.q.UpdateTract(ctx, params)
	if err != nil {
		return domain.Tract{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error updating tract")
	}

	return tractToDomain(row)
}

func (r *Repo) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) error {
	params := artel_q.SetTractEnabledParams{
		ID:      id,
		Enabled: enabled,
	}

	err := r.q.SetTractEnabled(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error setting tract enabled")
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteTract(ctx, id)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error deleting tract")
	}

	return nil
}

func tractToDomain(row artel_q.Tract) (domain.Tract, error) {
	var def domain.TractDefinition

	err := json.Unmarshal(row.Definition, &def)
	if err != nil {
		return domain.Tract{}, rerrors.Wrap(err, "error unmarshaling tract definition")
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

	return tract, nil
}

// --- tract_runs ---

func (r *Repo) InsertRun(ctx context.Context, run domain.TractRun) (domain.TractRun, error) {
	triggerPayload := run.TriggerPayload
	if triggerPayload == nil {
		triggerPayload = json.RawMessage("{}")
	}

	params := artel_q.InsertTractRunParams{
		TractID:       run.TractUuid,
		TriggerID:     uuid.NullUUID{UUID: run.TriggerUuid, Valid: run.TriggerUuid != uuid.Nil},
		StartedBy:     run.StartedBy,
		TriggerOutput: triggerPayload,
	}

	row, err := r.q.InsertTractRun(ctx, params)
	if err != nil {
		return domain.TractRun{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error inserting tract run")
	}

	return runToDomain(row), nil
}

func (r *Repo) GetRun(ctx context.Context, id uuid.UUID) (sql.Null[domain.TractRun], error) {
	row, err := r.q.GetTractRun(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.TractRun]{}, nil
		}

		return sql.Null[domain.TractRun]{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting tract run")
	}

	result := sql.Null[domain.TractRun]{V: runToDomain(row), Valid: true}

	return result, nil
}

func (r *Repo) ListRunsByTract(ctx context.Context, tractUuid uuid.UUID, limit int32) ([]domain.TractRun, error) {
	params := artel_q.ListTractRunsByTractParams{
		TractID: tractUuid,
		Limit:   limit,
	}

	rows, err := r.q.ListTractRunsByTract(ctx, params)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing tract runs")
	}

	runs := make([]domain.TractRun, len(rows))
	for i, row := range rows {
		runs[i] = runToDomain(row)
	}

	return runs, nil
}

func (r *Repo) UpdateRunStatus(ctx context.Context, id uuid.UUID, status domain.TractRunStatus, errMsg string) error {
	params := artel_q.UpdateTractRunStatusParams{
		ID:     id,
		Status: artel_q.TractRunStatus(status),
		Error:  sql.NullString{String: errMsg, Valid: errMsg != ""},
	}

	err := r.q.UpdateTractRunStatus(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error updating tract run status")
	}

	return nil
}

func (r *Repo) SweepStaleRuns(ctx context.Context, threshold time.Time) error {
	err := r.q.SweepStaleTractRuns(ctx, threshold)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error sweeping stale tract runs")
	}

	return nil
}

func runToDomain(row artel_q.TractRun) domain.TractRun {
	run := domain.TractRun{
		Uuid:           row.ID,
		TractUuid:      row.TractID,
		StartedBy:      row.StartedBy,
		Status:         domain.TractRunStatus(row.Status),
		TriggerPayload: row.TriggerOutput,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
	if row.TriggerID.Valid {
		run.TriggerUuid = row.TriggerID.UUID
	}

	if row.Error.Valid {
		run.Error = row.Error.String
	}

	return run
}

// --- tract_run_steps ---

func (r *Repo) InsertRunStep(ctx context.Context, step domain.TractRunStep) (domain.TractRunStep, error) {
	params := artel_q.InsertTractRunStepParams{
		RunID:    step.RunUuid,
		StepID:   step.StepId,
		StepName: step.StepName,
		StepType: step.StepType,
		Input:    pqtype.NullRawMessage{RawMessage: step.Input, Valid: step.Input != nil},
	}

	row, err := r.q.InsertTractRunStep(ctx, params)
	if err != nil {
		return domain.TractRunStep{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error inserting tract run step")
	}

	return runStepToDomain(row), nil
}

func (r *Repo) UpdateRunStepFinish(ctx context.Context, id uuid.UUID, status domain.TractRunStepStatus, output json.RawMessage, errMsg string) error {
	params := artel_q.UpdateTractRunStepFinishParams{
		ID:     id,
		Output: pqtype.NullRawMessage{RawMessage: output, Valid: output != nil},
		Status: artel_q.TractRunStepStatus(status),
		Error:  sql.NullString{String: errMsg, Valid: errMsg != ""},
	}

	err := r.q.UpdateTractRunStepFinish(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error updating tract run step")
	}

	return nil
}

func (r *Repo) ListRunStepsByRun(ctx context.Context, runUuid uuid.UUID) ([]domain.TractRunStep, error) {
	rows, err := r.q.ListTractRunStepsByRun(ctx, runUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing tract run steps")
	}

	steps := make([]domain.TractRunStep, len(rows))
	for i, row := range rows {
		steps[i] = runStepToDomain(row)
	}

	return steps, nil
}

func (r *Repo) SweepStaleRunSteps(ctx context.Context, threshold time.Time) error {
	err := r.q.SweepStaleTractRunSteps(ctx, threshold)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error sweeping stale tract run steps")
	}

	return nil
}

func runStepToDomain(row artel_q.TractRunStep) domain.TractRunStep {
	step := domain.TractRunStep{
		Uuid:      row.ID,
		RunUuid:   row.RunID,
		StepId:    row.StepID,
		StepName:  row.StepName,
		StepType:  row.StepType,
		Status:    domain.TractRunStepStatus(row.Status),
		StartedAt: row.StartedAt,
	}
	if row.Input.Valid {
		step.Input = row.Input.RawMessage
	}

	if row.Output.Valid {
		step.Output = row.Output.RawMessage
	}

	if row.Error.Valid {
		step.Error = row.Error.String
	}

	if row.FinishedAt.Valid {
		step.FinishedAt = row.FinishedAt.Time
	}

	return step
}
