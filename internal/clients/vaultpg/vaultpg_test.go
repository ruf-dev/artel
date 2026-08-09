package vaultpg

import (
	"errors"
	"net/url"
	"testing"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

func TestConfig_DSN(t *testing.T) {
	cases := []struct {
		name string
		cfg  Config
		want string
	}{
		{
			name: "basic",
			cfg: Config{
				Host:     "localhost",
				Port:     5432,
				Database: "vault_abc123",
				User:     "vaultpg_abc123",
				Password: "hunter2",
				SSLMode:  "disable",
			},
			want: "postgres://vaultpg_abc123:hunter2@localhost:5432/vault_abc123?sslmode=disable",
		},
		{
			name: "no sslmode set omits the query param entirely",
			cfg: Config{
				Host:     "db.internal",
				Port:     5433,
				Database: "postgres",
				User:     "admin",
				Password: "pw",
			},
			want: "postgres://admin:pw@db.internal:5433/postgres",
		},
		{
			name: "special characters in user and password are escaped",
			cfg: Config{
				Host:     "localhost",
				Port:     5432,
				Database: "db",
				User:     "user@name",
				Password: "p@ss:word/with?special&chars",
				SSLMode:  "require",
			},
			want: "postgres://user%40name:p%40ss%3Aword%2Fwith%3Fspecial&chars@localhost:5432/db?sslmode=require",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.cfg.DSN()
			if got != tc.want {
				t.Fatalf("DSN() = %q, want %q", got, tc.want)
			}

			// The DSN must always be parseable as a valid URL — lib/pq's URL-form DSN parser
			// relies on this.
			_, err := url.Parse(got)
			if err != nil {
				t.Fatalf("DSN() produced an unparseable URL: %v", err)
			}
		})
	}
}

func TestIsDuplicateObject(t *testing.T) {
	dupErr := &pq.Error{Code: pqerror.DuplicateObject}
	if !isDuplicateObject(dupErr) {
		t.Fatal("expected duplicate_object pq.Error to be classified as duplicate")
	}

	otherPqErr := &pq.Error{Code: pqerror.UniqueViolation}
	if isDuplicateObject(otherPqErr) {
		t.Fatal("expected a different pq.Error code to not be classified as duplicate")
	}

	wrapped := errors.New("pq: role \"x\" already exists")
	if !isDuplicateObject(wrapped) {
		t.Fatal("expected a plain error mentioning 'already exists' to be classified as duplicate")
	}

	unrelated := errors.New("connection refused")
	if isDuplicateObject(unrelated) {
		t.Fatal("expected an unrelated error to not be classified as duplicate")
	}
}

func TestIsDuplicateDatabase(t *testing.T) {
	dupErr := &pq.Error{Code: pqerror.DuplicateDatabase}
	if !isDuplicateDatabase(dupErr) {
		t.Fatal("expected duplicate_database pq.Error to be classified as duplicate")
	}

	otherPqErr := &pq.Error{Code: pqerror.DuplicateObject}
	if isDuplicateDatabase(otherPqErr) {
		t.Fatal("expected a different pq.Error code to not be classified as duplicate database")
	}

	wrapped := errors.New("pq: database \"x\" already exists")
	if !isDuplicateDatabase(wrapped) {
		t.Fatal("expected a plain error mentioning 'already exists' to be classified as duplicate")
	}
}
