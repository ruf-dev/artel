package users

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type UsersRepo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *UsersRepo {
	return &UsersRepo{q: q}
}

func (r *UsersRepo) Create(ctx context.Context, email string) (domain.User, error) {
	row, err := r.q.CreateUser(ctx, email)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error creating user")
	}

	u := domain.User{
		Uuid:      row.ID,
		Email:     row.Email,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	return u, nil
}

func (r *UsersRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error getting user by id")
	}

	u := domain.User{
		Uuid:      row.ID,
		Email:     row.Email,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	return u, nil
}

func (r *UsersRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error getting user by email")
	}

	u := domain.User{
		Uuid:      row.ID,
		Email:     row.Email,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
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
