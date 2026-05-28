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

func (r *UsersRepo) GetByTelegramId(ctx context.Context, telegramId string) (domain.User, error) {
	row, err := r.q.GetUserByTelegramId(ctx, telegramId)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error getting user by telegram id")
	}

	u := domain.User{
		Uuid:         row.ID,
		Email:        row.Email.String,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	return u, nil
}

func (r *UsersRepo) UpsertByTelegramId(ctx context.Context, telegramId string, username string, photoUrl string) (domain.User, error) {
	userRow, err := r.q.GetUserByTelegramId(ctx, telegramId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, rerrors.Wrap(err, "get user by telegram id")
	}

	if err == nil {
		params := artel_q.TouchTelegramIdentityParams{
			TelegramID: telegramId,
			PhotoUrl:   photoUrl,
		}
		touchErr := r.q.TouchTelegramIdentity(ctx, params)
		if touchErr != nil {
			return domain.User{}, rerrors.Wrap(touchErr, "touch telegram identity")
		}
		u := domain.User{
			Uuid:         userRow.ID,
			Email:        userRow.Email.String,
			Username:     userRow.Username,
			PasswordHash: userRow.PasswordHash,
			CreatedAt:    userRow.CreatedAt,
			UpdatedAt:    userRow.UpdatedAt,
		}
		return u, nil
	}

	newUser, err := r.q.CreateTelegramUser(ctx, username)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "create telegram user")
	}

	params := artel_q.InsertTelegramIdentityParams{
		UserID:     newUser.ID,
		TelegramID: telegramId,
		PhotoUrl:   photoUrl,
	}
	err = r.q.InsertTelegramIdentity(ctx, params)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "insert telegram identity")
	}

	u := domain.User{
		Uuid:         newUser.ID,
		Email:        newUser.Email.String,
		Username:     newUser.Username,
		PasswordHash: newUser.PasswordHash,
		CreatedAt:    newUser.CreatedAt,
		UpdatedAt:    newUser.UpdatedAt,
	}
	return u, nil
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
