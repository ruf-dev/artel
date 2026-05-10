package users

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
)

type UsersRepo struct {
	db sqldb.DB
}

func New(db sqldb.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) Create(ctx context.Context, email string) (domain.User, error) {
	query := `INSERT INTO users (email) VALUES ($1) RETURNING id, email, created_at, updated_at`

	var u domain.User
	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&u.Uuid, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error creating user")
	}

	return u, nil
}

func (r *UsersRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	query := `SELECT id, email, created_at, updated_at FROM users WHERE id = $1`

	var u domain.User
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&u.Uuid, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error getting user by id")
	}

	return u, nil
}

func (r *UsersRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT id, email, created_at, updated_at FROM users WHERE email = $1`

	var u domain.User
	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&u.Uuid, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error getting user by email")
	}

	return u, nil
}
