package users

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
	"golang.org/x/net/context"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type UsersRepo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *UsersRepo {
	return &UsersRepo{q: q}
}

func (r *UsersRepo) WithTx(tx *sql.Tx) repository.Users {
	return &UsersRepo{q: r.q.WithTx(tx)}
}

func (r *UsersRepo) Create(ctx context.Context, email, passwordHash string) (domain.User, error) {
	params := artel_q.CreateUserParams{
		Email:        sql.NullString{String: email, Valid: email != ""},
		PasswordHash: passwordHash,
	}
	row, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error creating user")
	}

	u := domain.User{
		Uuid:         row.ID,
		Email:        row.Email.String,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	return u, nil
}

func (r *UsersRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error getting user by id")
	}

	u := domain.User{
		Uuid:         row.ID,
		Email:        row.Email.String,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	return u, nil
}

func (r *UsersRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	nullEmail := sql.NullString{String: email, Valid: email != ""}
	row, err := r.q.GetUserByEmail(ctx, nullEmail)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error getting user by email")
	}

	u := domain.User{
		Uuid:         row.ID,
		Email:        row.Email.String,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	return u, nil
}

func (r *UsersRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteUser(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "error deleting user")
	}

	return nil
}

func (r *UsersRepo) GetByTelegramId(ctx context.Context, telegramId string) (sql.Null[domain.User], error) {
	row, err := r.q.GetUserByTelegramId(ctx, telegramId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.User]{}, nil
		}
		return sql.Null[domain.User]{}, rerrors.Wrap(err, "get user by telegram id")
	}

	u := domain.User{
		Uuid:         row.ID,
		Email:        row.Email.String,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	return sql.Null[domain.User]{V: u, Valid: true}, nil
}

func (r *UsersRepo) CreateByUsername(ctx context.Context, username string) (domain.User, error) {
	row, err := r.q.CreateByUsername(ctx, username)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "create user by username")
	}

	return domain.User{
		Uuid:         row.ID,
		Email:        row.Email.String,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}

func (r *UsersRepo) UpsertTelegramIdentity(ctx context.Context, identity domain.TelegramIdentity) error {
	params := artel_q.UpsertTelegramIdentityParams{
		UserID:     identity.UserUuid,
		TelegramID: identity.TelegramId,
		PhotoUrl:   identity.PhotoUrl,
	}
	err := r.q.UpsertTelegramIdentity(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "upsert telegram identity")
	}
	return nil
}

func (r *UsersRepo) GetTelegramPhotoUrl(ctx context.Context, userUuid uuid.UUID) (string, error) {
	photoUrl, err := r.q.GetTelegramPhotoUrlByUserId(ctx, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", rerrors.Wrap(err, "get telegram photo url by user id")
	}
	return photoUrl, nil
}
