package pendingauthcodes

import (
	"context"
	"time"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type PendingAuthCodesRepo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *PendingAuthCodesRepo {
	return &PendingAuthCodesRepo{q: q}
}

func (r *PendingAuthCodesRepo) Create(ctx context.Context, code, rawToken, codeChallenge, redirectUri, clientId string, expiresAt time.Time) error {
	err := r.q.CreatePendingAuthCode(ctx, artel_q.CreatePendingAuthCodeParams{
		Code:          code,
		RawToken:      rawToken,
		CodeChallenge: codeChallenge,
		RedirectUri:   redirectUri,
		ClientID:      clientId,
		ExpiresAt:     expiresAt,
	})
	if err != nil {
		return rerrors.Wrap(err, "create pending auth code")
	}
	return nil
}

func (r *PendingAuthCodesRepo) Get(ctx context.Context, code string) (domain.PendingAuthCode, error) {
	row, err := r.q.GetPendingAuthCode(ctx, code)
	if err != nil {
		return domain.PendingAuthCode{}, rerrors.Wrap(err, "get pending auth code")
	}
	return domain.PendingAuthCode{
		Code:          row.Code,
		RawToken:      row.RawToken,
		CodeChallenge: row.CodeChallenge,
		RedirectUri:   row.RedirectUri,
		ClientId:      row.ClientID,
		ExpiresAt:     row.ExpiresAt,
	}, nil
}

func (r *PendingAuthCodesRepo) Delete(ctx context.Context, code string) error {
	err := r.q.DeletePendingAuthCode(ctx, code)
	if err != nil {
		return rerrors.Wrap(err, "delete pending auth code")
	}
	return nil
}

func (r *PendingAuthCodesRepo) DeleteExpired(ctx context.Context) error {
	err := r.q.DeleteExpiredPendingAuthCodes(ctx)
	if err != nil {
		return rerrors.Wrap(err, "delete expired pending auth codes")
	}
	return nil
}
