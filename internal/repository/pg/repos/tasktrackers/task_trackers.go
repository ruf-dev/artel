package tasktrackers

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q             *artel_q.Queries
	encryptionKey []byte
}

func New(q *artel_q.Queries, encryptionKey []byte) *Repo {
	return &Repo{
		q:             q,
		encryptionKey: encryptionKey,
	}
}

func (r *Repo) Insert(ctx context.Context, tracker domain.TaskTracker) (domain.TaskTracker, error) {
	apiKeyEnc, err := cryptoutil.Encrypt(r.encryptionKey, []byte(tracker.ApiKey))
	if err != nil {
		return domain.TaskTracker{}, rerrors.Wrap(err, "error encrypting api key")
	}

	apiTokenEnc, err := cryptoutil.Encrypt(r.encryptionKey, []byte(tracker.ApiToken))
	if err != nil {
		return domain.TaskTracker{}, rerrors.Wrap(err, "error encrypting api token")
	}

	params := artel_q.InsertTaskTrackerParams{
		UserID:      tracker.UserUuid,
		Type:        tracker.Type,
		Name:        tracker.Name,
		ApiKeyEnc:   apiKeyEnc,
		ApiTokenEnc: apiTokenEnc,
	}

	row, err := r.q.InsertTaskTracker(ctx, params)
	if err != nil {
		return domain.TaskTracker{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error inserting task tracker")
	}

	return toDomain(row, tracker.ApiKey, tracker.ApiToken), nil
}

func (r *Repo) GetByUuid(ctx context.Context, id uuid.UUID) (domain.TaskTracker, error) {
	row, err := r.q.GetTaskTrackerByUuid(ctx, id)
	if err != nil {
		return domain.TaskTracker{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting task tracker")
	}

	apiKey, err := cryptoutil.Decrypt(r.encryptionKey, row.ApiKeyEnc)
	if err != nil {
		return domain.TaskTracker{}, rerrors.Wrap(err, "error decrypting api key")
	}

	apiToken, err := cryptoutil.Decrypt(r.encryptionKey, row.ApiTokenEnc)
	if err != nil {
		return domain.TaskTracker{}, rerrors.Wrap(err, "error decrypting api token")
	}

	return toDomain(row, string(apiKey), string(apiToken)), nil
}

func (r *Repo) ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.TaskTracker, error) {
	rows, err := r.q.ListTaskTrackersByUser(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing task trackers")
	}

	trackers := make([]domain.TaskTracker, len(rows))

	for i, row := range rows {
		apiKey, err := cryptoutil.Decrypt(r.encryptionKey, row.ApiKeyEnc)
		if err != nil {
			return nil, rerrors.Wrap(err, "error decrypting api key")
		}

		apiToken, err := cryptoutil.Decrypt(r.encryptionKey, row.ApiTokenEnc)
		if err != nil {
			return nil, rerrors.Wrap(err, "error decrypting api token")
		}

		trackers[i] = toDomain(row, string(apiKey), string(apiToken))
	}

	return trackers, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteTaskTracker(ctx, id)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error deleting task tracker")
	}

	return nil
}

func toDomain(row artel_q.TaskTracker, apiKey, apiToken string) domain.TaskTracker {
	return domain.TaskTracker{
		Uuid:      row.ID,
		UserUuid:  row.UserID,
		Type:      row.Type,
		Name:      row.Name,
		ApiKey:    apiKey,
		ApiToken:  apiToken,
		CreatedAt: row.CreatedAt,
	}
}
