package externalconnections

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// fakePostgresInstancesRepo is a hand-rolled fake of repository.PostgresInstances, scoped to the
// GetOwned/Update/RegisterOwned methods syncOwnedPostgresInstance uses.
type fakePostgresInstancesRepo struct {
	getOwnedFunc      func(ctx context.Context, userID uuid.UUID) (sql.Null[domain.PostgresInstance], error)
	updateFunc        func(ctx context.Context, id uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) error
	registerOwnedFunc func(ctx context.Context, ownerUserID uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) (uuid.UUID, error)
}

func (f *fakePostgresInstancesRepo) Register(ctx context.Context, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) (uuid.UUID, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) Get(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) RandomPick(ctx context.Context) (domain.PostgresInstance, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) PickForUser(ctx context.Context, userID uuid.UUID) (domain.PostgresInstance, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) GetOwned(ctx context.Context, userID uuid.UUID) (sql.Null[domain.PostgresInstance], error) {
	return f.getOwnedFunc(ctx, userID)
}

func (f *fakePostgresInstancesRepo) List(ctx context.Context) ([]domain.PostgresInstance, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) Update(ctx context.Context, id uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) error {
	return f.updateFunc(ctx, id, host, port, adminDatabase, username, passwordPlain, sslMode)
}

func (f *fakePostgresInstancesRepo) RegisterOwned(ctx context.Context, ownerUserID uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) (uuid.UUID, error) {
	return f.registerOwnedFunc(ctx, ownerUserID, host, port, adminDatabase, username, passwordPlain, sslMode)
}

func (f *fakePostgresInstancesRepo) Delete(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) DeleteOwnedIfUnreferenced(ctx context.Context, ownerUserID uuid.UUID) error {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) Exists(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) WithTx(tx postgres.DB) repository.PostgresInstances {
	panic("not implemented")
}

func withUser(userUuid uuid.UUID) context.Context {
	uc := user_context.UserContext{UserUuid: userUuid, UserName: "tester"}

	return user_context.WithUserContext(context.Background(), uc)
}

func TestCheckPostgresConnection_Unauthenticated(t *testing.T) {
	svc := &Service{}

	err := svc.CheckPostgresConnection(context.Background(), "localhost", 5432, "postgres", "u", "p", "disable")
	if !errors.Is(err, user_errors.Unauthenticated) {
		t.Fatalf("expected Unauthenticated, got %v", err)
	}
}

func TestCheckPostgresConnection_Unreachable(t *testing.T) {
	svc := &Service{}

	// Port 1 on loopback: nothing listens there, so the connection attempt fails fast without
	// needing a real postgres server in the test environment.
	err := svc.CheckPostgresConnection(withUser(uuid.New()), "127.0.0.1", 1, "postgres", "u", "p", "disable")
	if !errors.Is(err, user_errors.PostgresConnectionValidationFailed) {
		t.Fatalf("expected PostgresConnectionValidationFailed, got %v", err)
	}
}

func TestAddPostgresConnection_Unauthenticated(t *testing.T) {
	svc := &Service{}

	_, err := svc.AddPostgresConnection(context.Background(), "localhost", 5432, "postgres", "u", "p", "disable")
	if !errors.Is(err, user_errors.Unauthenticated) {
		t.Fatalf("expected Unauthenticated, got %v", err)
	}
}

func TestAddPostgresConnection_ValidationFailure(t *testing.T) {
	svc := &Service{}

	_, err := svc.AddPostgresConnection(withUser(uuid.New()), "127.0.0.1", 1, "postgres", "u", "p", "disable")
	if !errors.Is(err, user_errors.PostgresConnectionValidationFailed) {
		t.Fatalf("expected PostgresConnectionValidationFailed, got %v", err)
	}
}

func TestSyncOwnedPostgresInstance_UpdatesExisting(t *testing.T) {
	userUuid := uuid.New()
	ownedID := uuid.New()
	updateCalled := false

	repo := &fakePostgresInstancesRepo{
		getOwnedFunc: func(ctx context.Context, gotUserID uuid.UUID) (sql.Null[domain.PostgresInstance], error) {
			return sql.Null[domain.PostgresInstance]{Valid: true, V: domain.PostgresInstance{Uuid: ownedID}}, nil
		},
		updateFunc: func(ctx context.Context, id uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) error {
			updateCalled = true

			if id != ownedID {
				t.Fatalf("expected update on owned id %s, got %s", ownedID, id)
			}

			return nil
		},
	}

	svc := &Service{postgresInstances: repo}

	err := svc.syncOwnedPostgresInstance(context.Background(), userUuid, "host", 5432, "db", "user", "pw", "disable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !updateCalled {
		t.Fatal("expected repo.Update to be called for an already-owned instance")
	}
}

func TestSyncOwnedPostgresInstance_RegistersNew(t *testing.T) {
	userUuid := uuid.New()
	registerCalled := false

	repo := &fakePostgresInstancesRepo{
		getOwnedFunc: func(ctx context.Context, gotUserID uuid.UUID) (sql.Null[domain.PostgresInstance], error) {
			return sql.Null[domain.PostgresInstance]{Valid: false}, nil
		},
		registerOwnedFunc: func(ctx context.Context, ownerUserID uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) (uuid.UUID, error) {
			registerCalled = true

			if ownerUserID != userUuid {
				t.Fatalf("expected owner %s, got %s", userUuid, ownerUserID)
			}

			return uuid.New(), nil
		},
	}

	svc := &Service{postgresInstances: repo}

	err := svc.syncOwnedPostgresInstance(context.Background(), userUuid, "host", 5432, "db", "user", "pw", "disable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !registerCalled {
		t.Fatal("expected repo.RegisterOwned to be called when no owned instance exists yet")
	}
}
