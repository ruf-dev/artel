package users

import (
	context0 "context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"github.com/ruf-dev/artel/internal/utils"
)

type UsersRepo struct {
	q  *artel_q.Queries
	db postgres.DB
}

func New(q *artel_q.Queries, db postgres.DB) *UsersRepo {
	return &UsersRepo{q: q, db: db}
}

func (r *UsersRepo) WithTx(tx *sql.Tx) repository.Users {
	return &UsersRepo{q: r.q.WithTx(tx), db: r.db}
}

func (r *UsersRepo) Create(ctx context0.Context, email, passwordHash string) (domain.User, error) {
	params := artel_q.CreateUserParams{
		Email:        sql.NullString{String: email, Valid: email != ""},
		PasswordHash: passwordHash,
	}

	row, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return domain.User{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error creating user")
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

func (r *UsersRepo) GetByID(ctx context0.Context, id uuid.UUID) (domain.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error getting user by id")
	}

	u := domain.User{
		Uuid:         row.ID,
		Email:        row.Email.String,
		Username:     row.Username,
		PhotoUrl:     row.PhotoUrl,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}

	return u, nil
}

func (r *UsersRepo) GetByEmail(ctx context0.Context, email string) (domain.User, error) {
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

func (r *UsersRepo) FindByEmail(ctx context0.Context, email string) (sql.Null[domain.User], error) {
	nullEmail := sql.NullString{String: email, Valid: email != ""}

	row, err := r.q.GetUserByEmail(ctx, nullEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.User]{}, nil
		}

		return sql.Null[domain.User]{}, rerrors.Wrap(err, "error finding user by email")
	}

	u := domain.User{
		Uuid:         row.ID,
		Email:        row.Email.String,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}

	return sql.Null[domain.User]{V: u, Valid: true}, nil
}

func (r *UsersRepo) Delete(ctx context0.Context, id uuid.UUID) error {
	err := r.q.DeleteUser(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "error deleting user")
	}

	return nil
}

func (r *UsersRepo) GetByTelegramId(ctx context0.Context, telegramId string) (sql.Null[domain.User], error) {
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

func (r *UsersRepo) CreateByUsername(ctx context0.Context, username, photoUrl string) (domain.User, error) {
	params := artel_q.CreateByUsernameParams{
		Username: username,
		PhotoUrl: photoUrl,
	}

	row, err := r.q.CreateByUsername(ctx, params)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "create user by username")
	}

	return domain.User{
		Uuid:         row.ID,
		Email:        row.Email.String,
		Username:     row.Username,
		PhotoUrl:     row.PhotoUrl,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}

func (r *UsersRepo) UpsertTelegramIdentity(ctx context0.Context, identity domain.TelegramIdentity) error {
	params := artel_q.UpsertTelegramIdentityParams{
		UserID:     identity.UserUuid,
		TelegramID: identity.TelegramId,
	}

	err := r.q.UpsertTelegramIdentity(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "upsert telegram identity")
	}

	return nil
}

func (r *UsersRepo) UpdatePhotoUrl(ctx context0.Context, userUuid uuid.UUID, photoUrl string) error {
	params := artel_q.UpdateUserPhotoUrlParams{
		ID:       userUuid,
		PhotoUrl: photoUrl,
	}

	err := r.q.UpdateUserPhotoUrl(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "update user photo url")
	}

	return nil
}

func (r *UsersRepo) UpdatePasswordHash(ctx context0.Context, userUuid uuid.UUID, passwordHash string) error {
	params := artel_q.UpdateUserPasswordHashParams{
		ID:           userUuid,
		PasswordHash: passwordHash,
	}

	err := r.q.UpdateUserPasswordHash(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error updating user password hash")
	}

	return nil
}

func (r *UsersRepo) GetTelegramPhotoUrl(ctx context0.Context, userUuid uuid.UUID) (string, error) {
	photoUrl, err := r.q.GetTelegramPhotoUrlByUserId(ctx, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", rerrors.Wrap(err, "get telegram photo url by user id")
	}

	return photoUrl, nil
}

func (r *UsersRepo) GetTelegramIdByUserId(ctx context0.Context, userUuid uuid.UUID) (sql.Null[string], error) {
	telegramId, err := r.q.GetTelegramIdByUserId(ctx, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[string]{}, nil
		}

		return sql.Null[string]{}, rerrors.Wrap(err, "get telegram id by user id")
	}

	return sql.Null[string]{V: telegramId, Valid: true}, nil
}

func (r *UsersRepo) ListAll(ctx context0.Context, req domain.ListUsersReq) ([]domain.User, int64, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	base := psql.Select("id::text", "username", "email", "created_at").From("users")
	countBase := psql.Select("COUNT(*)").From("users")

	if req.Search != "" {
		like := "%" + req.Search + "%"
		cond := sq.Or{sq.ILike{"username": like}, sq.ILike{"email": like}}
		base = base.Where(cond)
		countBase = countBase.Where(cond)
	}

	if req.Paging.Limit > 0 {
		base = base.Limit(uint64(req.Paging.Limit))
	}

	base = base.Offset(uint64(req.Paging.Offset)).OrderBy("created_at DESC")

	countSQL, countArgs, err := countBase.ToSql()
	if err != nil {
		return nil, 0, rerrors.Wrap(err, "error building count query")
	}

	var total int64

	err = r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, rerrors.Wrap(err, "error counting users")
	}

	listSQL, listArgs, err := base.ToSql()
	if err != nil {
		return nil, 0, rerrors.Wrap(err, "error building list query")
	}

	rows, err := r.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, pg_err.UnwrapPgErr(err)
	}

	defer utils.CloseWithLog(rows, "list users rows")

	var result []domain.User

	for rows.Next() {
		var idStr string

		var u domain.User

		var email sql.NullString

		err = rows.Scan(&idStr, &u.Username, &email, &u.CreatedAt)
		if err != nil {
			return nil, 0, rerrors.Wrap(err, "error scanning user row")
		}

		u.Uuid, err = uuid.Parse(idStr)
		if err != nil {
			return nil, 0, rerrors.Wrap(err, "error parsing user uuid")
		}

		u.Email = email.String
		result = append(result, u)
	}

	return result, total, nil
}

func (r *UsersRepo) GetDetailsById(ctx context0.Context, id uuid.UUID) (domain.UserDetails, error) {
	row, err := r.q.GetUserDetails(ctx, id)
	if err != nil {
		return domain.UserDetails{}, rerrors.Wrap(err, "error getting user details")
	}

	u := domain.User{
		Uuid:         row.ID,
		Username:     row.Username,
		Email:        row.Email.String,
		PhotoUrl:     row.PhotoUrl,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	permissions := domain.UserPermissions{
		UserUuid:        row.ID,
		IsAdministrator: row.IsAdministrator.Bool,
	}
	// EffectiveSubscription is left zero-valued here — callers (auth.GetMe, adminusers.GetUser)
	// populate it via SubscriptionService.GetEffective, since it's a plan+override merge rather
	// than a column this repo-layer query can join.
	details := domain.UserDetails{
		User:        u,
		Permissions: permissions,
	}

	return details, nil
}
