package postgresinstances

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
)

// fakePostgresInstancesRepo is a hand-rolled fake of repository.PostgresInstances — only the
// methods exercised by a given test set a func field, the rest panic if called unexpectedly.
type fakePostgresInstancesRepo struct {
	registerFunc func(ctx context.Context, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) (uuid.UUID, error)
	getFunc      func(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error)
	listFunc     func(ctx context.Context) ([]domain.PostgresInstance, error)
	updateFunc   func(ctx context.Context, id uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) error
	deleteFunc   func(ctx context.Context, id uuid.UUID) error
	existsFunc   func(ctx context.Context) (bool, error)
}

func (f *fakePostgresInstancesRepo) Register(
	ctx context.Context, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string,
) (uuid.UUID, error) {
	return f.registerFunc(ctx, host, port, adminDatabase, username, passwordPlain, sslMode)
}

func (f *fakePostgresInstancesRepo) Get(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error) {
	return f.getFunc(ctx, id)
}

func (f *fakePostgresInstancesRepo) RandomPick(ctx context.Context) (domain.PostgresInstance, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) PickForUser(ctx context.Context, userID uuid.UUID) (domain.PostgresInstance, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) GetOwned(ctx context.Context, userID uuid.UUID) (sql.Null[domain.PostgresInstance], error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) List(ctx context.Context) ([]domain.PostgresInstance, error) {
	return f.listFunc(ctx)
}

func (f *fakePostgresInstancesRepo) Update(
	ctx context.Context, id uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string,
) error {
	return f.updateFunc(ctx, id, host, port, adminDatabase, username, passwordPlain, sslMode)
}

func (f *fakePostgresInstancesRepo) RegisterOwned(
	ctx context.Context, ownerUserID uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string,
) (uuid.UUID, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return f.deleteFunc(ctx, id)
}

func (f *fakePostgresInstancesRepo) DeleteOwnedIfUnreferenced(ctx context.Context, ownerUserID uuid.UUID) error {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) Exists(ctx context.Context) (bool, error) {
	return f.existsFunc(ctx)
}

func (f *fakePostgresInstancesRepo) WithTx(tx postgres.DB) repository.PostgresInstances {
	panic("not implemented")
}

func TestRegisterPostgresInstance(t *testing.T) {
	wantID := uuid.New()

	repo := &fakePostgresInstancesRepo{
		registerFunc: func(ctx context.Context, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) (uuid.UUID, error) {
			if host != "db.internal" || port != 5432 || adminDatabase != "postgres" || username != "admin" || sslMode != "require" {
				t.Fatalf("unexpected register args: %s %d %s %s %s", host, port, adminDatabase, username, sslMode)
			}

			if string(passwordPlain) != "hunter2" {
				t.Fatalf("unexpected password: %s", passwordPlain)
			}

			return wantID, nil
		},
	}

	svc := &Service{postgresInstancesRepo: repo}

	gotID, err := svc.RegisterPostgresInstance(context.Background(), "db.internal", 5432, "postgres", "admin", "hunter2", "require")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotID != wantID.String() {
		t.Fatalf("got id %q, want %q", gotID, wantID.String())
	}
}

func TestGetPostgresInstance_InvalidUuid(t *testing.T) {
	svc := &Service{postgresInstancesRepo: &fakePostgresInstancesRepo{}}

	_, err := svc.GetPostgresInstance(context.Background(), "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid uuid, got nil")
	}
}

func TestGetPostgresInstance_RepoError(t *testing.T) {
	wantErr := errors.New("boom")

	repo := &fakePostgresInstancesRepo{
		getFunc: func(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error) {
			return domain.PostgresInstance{}, wantErr
		},
	}

	svc := &Service{postgresInstancesRepo: repo}

	_, err := svc.GetPostgresInstance(context.Background(), uuid.New().String())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}

func TestListPostgresInstances(t *testing.T) {
	want := []domain.PostgresInstance{{Uuid: uuid.New(), Host: "a"}, {Uuid: uuid.New(), Host: "b"}}

	repo := &fakePostgresInstancesRepo{
		listFunc: func(ctx context.Context) ([]domain.PostgresInstance, error) {
			return want, nil
		},
	}

	svc := &Service{postgresInstancesRepo: repo}

	got, err := svc.ListPostgresInstances(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("got %d instances, want %d", len(got), len(want))
	}
}

func TestUpdatePostgresInstance(t *testing.T) {
	id := uuid.New()
	called := false

	repo := &fakePostgresInstancesRepo{
		updateFunc: func(ctx context.Context, gotID uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) error {
			called = true

			if gotID != id {
				t.Fatalf("got id %s, want %s", gotID, id)
			}

			return nil
		},
	}

	svc := &Service{postgresInstancesRepo: repo}

	err := svc.UpdatePostgresInstance(context.Background(), id.String(), "host", 5432, "db", "user", "pw", "disable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Fatal("expected repo Update to be called")
	}
}

func TestDeletePostgresInstance(t *testing.T) {
	id := uuid.New()
	called := false

	repo := &fakePostgresInstancesRepo{
		deleteFunc: func(ctx context.Context, gotID uuid.UUID) error {
			called = true

			if gotID != id {
				t.Fatalf("got id %s, want %s", gotID, id)
			}

			return nil
		},
	}

	svc := &Service{postgresInstancesRepo: repo}

	err := svc.DeletePostgresInstance(context.Background(), id.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Fatal("expected repo Delete to be called")
	}
}

func TestHasPostgresInstances(t *testing.T) {
	repo := &fakePostgresInstancesRepo{
		existsFunc: func(ctx context.Context) (bool, error) {
			return true, nil
		},
	}

	svc := &Service{postgresInstancesRepo: repo}

	got, err := svc.HasPostgresInstances(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !got {
		t.Fatal("expected true")
	}
}

func TestTestPostgresInstance_InvalidUuid(t *testing.T) {
	svc := &Service{postgresInstancesRepo: &fakePostgresInstancesRepo{}}

	err := svc.TestPostgresInstance(context.Background(), "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid uuid, got nil")
	}
}

func TestTestPostgresInstance_Unreachable(t *testing.T) {
	repo := &fakePostgresInstancesRepo{
		getFunc: func(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error) {
			return domain.PostgresInstance{
				Host: "127.0.0.1", Port: 1, AdminDatabase: "postgres", Username: "u", Password: "p", SSLMode: "disable",
			}, nil
		},
	}

	svc := &Service{postgresInstancesRepo: repo}

	err := svc.TestPostgresInstance(context.Background(), uuid.New().String())
	if err == nil {
		t.Fatal("expected error pinging an unreachable postgres server, got nil")
	}
}
