package pg_err

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	errors "go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func UnwrapPgErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return user_errors.NotFound
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case pqerror.UniqueViolation:
			return errors.Wrap(user_errors.AlreadyExists)
		}
	}

	return err
}
