package couchinstances

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/stretchr/testify/require"
)

// fakeCouchInstancesRepo is a minimal in-memory repository.CouchInstances used to test
// Service.HasCouchInstances without a real database — following the fakeTractsRepo convention
// in internal/service/v1/tract/mocks_test.go (a plain struct with in-memory state, not an
// embedded-interface trick).
type fakeCouchInstancesRepo struct {
	instances map[uuid.UUID]domain.CouchInstance

	existsResult bool
	existsErr    error
}

func newFakeCouchInstancesRepo() *fakeCouchInstancesRepo {
	repo := &fakeCouchInstancesRepo{
		instances: map[uuid.UUID]domain.CouchInstance{},
	}

	return repo
}

func (f *fakeCouchInstancesRepo) Register(
	_ context.Context, url, username string, _ []byte,
) (uuid.UUID, error) {
	id := uuid.New()
	f.instances[id] = domain.CouchInstance{
		Uuid:     id,
		Url:      url,
		Username: username,
	}

	return id, nil
}

func (f *fakeCouchInstancesRepo) Get(_ context.Context, id uuid.UUID) (domain.CouchInstance, error) {
	instance, ok := f.instances[id]
	if !ok {
		return domain.CouchInstance{}, errors.New("not found")
	}

	return instance, nil
}

func (f *fakeCouchInstancesRepo) RandomPick(_ context.Context) (domain.CouchInstanceWithAccount, error) {
	for _, instance := range f.instances {
		result := domain.CouchInstanceWithAccount{Instance: instance}

		return result, nil
	}

	return domain.CouchInstanceWithAccount{}, errors.New("no instances")
}

// PickForUser is unexercised by TestService_HasCouchInstances — it just falls back to RandomPick,
// since this fake has no notion of ownership.
func (f *fakeCouchInstancesRepo) PickForUser(ctx context.Context, _ uuid.UUID) (domain.CouchInstanceWithAccount, error) {
	return f.RandomPick(ctx)
}

// GetOwned is unexercised by TestService_HasCouchInstances — always reports no owned instance.
func (f *fakeCouchInstancesRepo) GetOwned(_ context.Context, _ uuid.UUID) (sql.Null[domain.CouchInstance], error) {
	return sql.Null[domain.CouchInstance]{}, nil
}

// RegisterOwned is unexercised by TestService_HasCouchInstances — behaves like Register.
func (f *fakeCouchInstancesRepo) RegisterOwned(
	_ context.Context, _ uuid.UUID, url, username string, _ []byte,
) (uuid.UUID, error) {
	id := uuid.New()
	f.instances[id] = domain.CouchInstance{
		Uuid:     id,
		Url:      url,
		Username: username,
	}

	return id, nil
}

// DeleteOwnedIfUnreferenced is unexercised by TestService_HasCouchInstances — always a no-op.
func (f *fakeCouchInstancesRepo) DeleteOwnedIfUnreferenced(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (f *fakeCouchInstancesRepo) List(_ context.Context) ([]domain.CouchInstance, error) {
	instances := make([]domain.CouchInstance, 0, len(f.instances))
	for _, instance := range f.instances {
		instances = append(instances, instance)
	}

	return instances, nil
}

func (f *fakeCouchInstancesRepo) Update(_ context.Context, id uuid.UUID, url, username string, _ []byte) error {
	instance, ok := f.instances[id]
	if !ok {
		return errors.New("not found")
	}

	instance.Url = url
	instance.Username = username
	f.instances[id] = instance

	return nil
}

func (f *fakeCouchInstancesRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(f.instances, id)

	return nil
}

func (f *fakeCouchInstancesRepo) Exists(_ context.Context) (bool, error) {
	if f.existsErr != nil {
		return false, f.existsErr
	}

	return f.existsResult, nil
}

func (f *fakeCouchInstancesRepo) WithTx(_ postgres.DB) repository.CouchInstances {
	return f
}

var _ repository.CouchInstances = (*fakeCouchInstancesRepo)(nil)

func TestService_HasCouchInstances(t *testing.T) {
	t.Run("returns true when instances exist", func(t *testing.T) {
		fakeRepo := newFakeCouchInstancesRepo()
		fakeRepo.existsResult = true

		svc := &Service{couchInstancesRepo: fakeRepo}

		exists, err := svc.HasCouchInstances(context.Background())
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("returns false when no instances exist", func(t *testing.T) {
		fakeRepo := newFakeCouchInstancesRepo()
		fakeRepo.existsResult = false

		svc := &Service{couchInstancesRepo: fakeRepo}

		exists, err := svc.HasCouchInstances(context.Background())
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("propagates wrapped repo error", func(t *testing.T) {
		fakeRepo := newFakeCouchInstancesRepo()
		fakeRepo.existsErr = errors.New("db unavailable")

		svc := &Service{couchInstancesRepo: fakeRepo}

		exists, err := svc.HasCouchInstances(context.Background())
		require.Error(t, err)
		require.False(t, exists)
		require.ErrorIs(t, err, fakeRepo.existsErr)
	})
}
