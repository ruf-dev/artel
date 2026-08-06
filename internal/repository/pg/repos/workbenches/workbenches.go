package workbenches

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q *artel_q.Queries
}

func New(db postgres.DB) *Repo {
	return &Repo{
		q: artel_q.New(db),
	}
}

func (r *Repo) Create(
	ctx context.Context,
	vaultID, userID uuid.UUID,
	volumeName string,
	dockerHostID uuid.UUID,
) (domain.Workbench, error) {
	params := artel_q.CreateWorkbenchParams{
		VaultID:      vaultID,
		UserID:       userID,
		VolumeName:   volumeName,
		DockerHostID: uuid.NullUUID{UUID: dockerHostID, Valid: true},
	}

	row, err := r.q.CreateWorkbench(ctx, params)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error creating workbench")
	}

	w := rowToWorkbench(
		row.ID,
		row.VaultID,
		row.UserID,
		row.Status,
		row.AuthMode,
		row.ContainerID,
		row.VolumeName,
		row.CreatedAt,
		row.StartedAt,
		row.StoppedAt,
		row.DockerHostID,
	)

	return w, nil
}

func (r *Repo) GetByVaultID(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error) {
	row, err := r.q.GetWorkbenchByVaultID(ctx, vaultID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault id")
	}

	w := rowToWorkbench(
		row.ID,
		row.VaultID,
		row.UserID,
		row.Status,
		row.AuthMode,
		row.ContainerID,
		row.VolumeName,
		row.CreatedAt,
		row.StartedAt,
		row.StoppedAt,
		row.DockerHostID,
	)

	return w, nil
}

func (r *Repo) MarkRunning(ctx context.Context, vaultID uuid.UUID, authMode domain.WorkbenchAuthMode) error {
	params := artel_q.MarkWorkbenchRunningParams{
		VaultID: vaultID,
		AuthMode: artel_q.NullWorkbenchAuthMode{
			WorkbenchAuthMode: artel_q.WorkbenchAuthMode(authMode),
			Valid:             authMode != "",
		},
	}

	err := r.q.MarkWorkbenchRunning(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error marking workbench running")
	}

	return nil
}

func (r *Repo) MarkContainerCreated(ctx context.Context, vaultID uuid.UUID, containerID string) error {
	params := artel_q.MarkWorkbenchContainerCreatedParams{
		VaultID:     vaultID,
		ContainerID: sql.NullString{String: containerID, Valid: containerID != ""},
	}

	err := r.q.MarkWorkbenchContainerCreated(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error marking workbench container created")
	}

	return nil
}

func (r *Repo) MarkConfiguring(ctx context.Context, vaultID uuid.UUID) error {
	err := r.q.MarkWorkbenchConfiguring(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error marking workbench configuring")
	}

	return nil
}

func (r *Repo) MarkStopped(ctx context.Context, vaultID uuid.UUID) error {
	err := r.q.MarkWorkbenchStopped(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error marking workbench stopped")
	}

	return nil
}

func (r *Repo) MarkRemoved(ctx context.Context, vaultID uuid.UUID) error {
	err := r.q.MarkWorkbenchRemoved(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error marking workbench removed")
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, vaultID uuid.UUID) error {
	err := r.q.DeleteWorkbench(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error deleting workbench")
	}

	return nil
}

func (r *Repo) WithTx(tx postgres.DB) repository.Workbenches {
	return New(tx)
}

func rowToWorkbench(
	id, vaultID, userID uuid.UUID,
	status artel_q.WorkbenchStatus,
	authMode artel_q.NullWorkbenchAuthMode,
	containerID sql.NullString,
	volumeName string,
	createdAt time.Time,
	startedAt, stoppedAt sql.NullTime,
	dockerHostID uuid.NullUUID,
) domain.Workbench {
	w := domain.Workbench{
		Uuid:       id,
		VaultUuid:  vaultID,
		UserUuid:   userID,
		Status:     domain.WorkbenchStatus(status),
		VolumeName: volumeName,
		CreatedAt:  createdAt,
	}

	if authMode.Valid {
		w.AuthMode = domain.WorkbenchAuthMode(authMode.WorkbenchAuthMode)
	}

	if containerID.Valid {
		w.ContainerId = containerID.String
	}

	if startedAt.Valid {
		startedAtVal := startedAt.Time
		w.StartedAt = &startedAtVal
	}

	if stoppedAt.Valid {
		stoppedAtVal := stoppedAt.Time
		w.StoppedAt = &stoppedAtVal
	}

	if dockerHostID.Valid {
		dockerHostIDVal := dockerHostID.UUID
		w.DockerHostUuid = &dockerHostIDVal
	}

	return w
}
