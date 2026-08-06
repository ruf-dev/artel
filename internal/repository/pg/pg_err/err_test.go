package pg_err

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func TestUnwrapPgErr(t *testing.T) {
	t.Run("no rows becomes NotFound", func(t *testing.T) {
		got := UnwrapPgErr(sql.ErrNoRows)
		if !errors.Is(got, user_errors.NotFound) {
			t.Fatalf("expected NotFound, got %v", got)
		}
	})

	t.Run("unique violation becomes AlreadyExists", func(t *testing.T) {
		pqErr := &pq.Error{Code: pqerror.UniqueViolation, Message: "duplicate key value violates unique constraint"}

		got := UnwrapPgErr(pqErr)
		if !errors.Is(got, user_errors.AlreadyExists) {
			t.Fatalf("expected AlreadyExists, got %v", got)
		}
	})

	t.Run("unrelated pg error passes through unchanged", func(t *testing.T) {
		pqErr := &pq.Error{Code: pqerror.NotNullViolation, Message: "null value violates not-null constraint"}

		got := UnwrapPgErr(pqErr)
		if !errors.Is(got, pqErr) {
			t.Fatalf("expected the original pq error to pass through, got %v", got)
		}
	})

	t.Run("non-pg error passes through unchanged", func(t *testing.T) {
		plain := errors.New("some driver error")

		got := UnwrapPgErr(plain)
		if !errors.Is(got, plain) {
			t.Fatalf("expected the original error to pass through, got %v", got)
		}
	})
}
