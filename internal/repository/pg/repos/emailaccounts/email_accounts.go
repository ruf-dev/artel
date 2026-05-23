package emailaccounts

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
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

func (r *Repo) Insert(ctx context.Context, account domain.EmailAccount) (domain.EmailAccount, error) {
	passwordEnc, err := cryptoutil.Encrypt(r.encryptionKey, []byte(account.Password))
	if err != nil {
		return domain.EmailAccount{}, rerrors.Wrap(err, "encrypt password")
	}

	params := artel_q.InsertEmailAccountParams{
		UserID:      account.UserUuid,
		Email:       account.Email,
		ImapHost:    account.ImapHost,
		ImapPort:    int32(account.ImapPort),
		SmtpHost:    account.SmtpHost,
		SmtpPort:    int32(account.SmtpPort),
		PasswordEnc: passwordEnc,
	}
	row, err := r.q.InsertEmailAccount(ctx, params)
	if err != nil {
		return domain.EmailAccount{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "insert email account")
	}

	return toDomain(row, account.Password), nil
}

func (r *Repo) GetByUuid(ctx context.Context, id uuid.UUID) (domain.EmailAccount, error) {
	row, err := r.q.GetEmailAccountByUuid(ctx, id)
	if err != nil {
		return domain.EmailAccount{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "get email account")
	}

	passwordPlain, err := cryptoutil.Decrypt(r.encryptionKey, row.PasswordEnc)
	if err != nil {
		return domain.EmailAccount{}, rerrors.Wrap(err, "decrypt password")
	}

	return toDomain(row, string(passwordPlain)), nil
}

func (r *Repo) ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.EmailAccount, error) {
	rows, err := r.q.ListEmailAccountsByUser(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "list email accounts")
	}

	accounts := make([]domain.EmailAccount, len(rows))
	for i, row := range rows {
		passwordPlain, err := cryptoutil.Decrypt(r.encryptionKey, row.PasswordEnc)
		if err != nil {
			return nil, rerrors.Wrap(err, "decrypt password")
		}
		accounts[i] = toDomain(row, string(passwordPlain))
	}
	return accounts, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteEmailAccount(ctx, id)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "delete email account")
	}
	return nil
}

func toDomain(row artel_q.EmailAccount, password string) domain.EmailAccount {
	return domain.EmailAccount{
		Uuid:      row.ID,
		UserUuid:  row.UserID,
		Email:     row.Email,
		ImapHost:  row.ImapHost,
		ImapPort:  int(row.ImapPort),
		SmtpHost:  row.SmtpHost,
		SmtpPort:  int(row.SmtpPort),
		Password:  password,
		CreatedAt: row.CreatedAt,
	}
}
