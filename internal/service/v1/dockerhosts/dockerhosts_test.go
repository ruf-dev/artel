package dockerhosts

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/stretchr/testify/require"
)

// fakeDockerHostsRepo is a minimal in-memory repository.DockerHosts used to test Service without
// a real database — following the fakeCouchInstancesRepo convention in
// internal/service/v1/couchinstances/couchinstances_test.go (a plain struct with in-memory
// state, not an embedded-interface trick).
type fakeDockerHostsRepo struct {
	hosts map[uuid.UUID]domain.DockerHost

	existsResult bool
	existsErr    error
}

func newFakeDockerHostsRepo() *fakeDockerHostsRepo {
	repo := &fakeDockerHostsRepo{
		hosts: map[uuid.UUID]domain.DockerHost{},
	}

	return repo
}

func (f *fakeDockerHostsRepo) Register(_ context.Context, url, caCert, clientCert, clientKey string) (uuid.UUID, error) {
	id := uuid.New()
	f.hosts[id] = domain.DockerHost{
		Uuid:       id,
		Url:        url,
		CaCert:     caCert,
		ClientCert: clientCert,
		ClientKey:  clientKey,
	}

	return id, nil
}

func (f *fakeDockerHostsRepo) Get(_ context.Context, id uuid.UUID) (domain.DockerHost, error) {
	host, ok := f.hosts[id]
	if !ok {
		return domain.DockerHost{}, errors.New("not found")
	}

	// Get is creds-free, mirroring the real repo's Get vs GetWithCreds split.
	host.CaCert = ""
	host.ClientCert = ""
	host.ClientKey = ""

	return host, nil
}

func (f *fakeDockerHostsRepo) GetWithCreds(_ context.Context, id uuid.UUID) (domain.DockerHost, error) {
	host, ok := f.hosts[id]
	if !ok {
		return domain.DockerHost{}, errors.New("not found")
	}

	return host, nil
}

func (f *fakeDockerHostsRepo) List(_ context.Context) ([]domain.DockerHost, error) {
	hosts := make([]domain.DockerHost, 0, len(f.hosts))
	for _, host := range f.hosts {
		host.CaCert = ""
		host.ClientCert = ""
		host.ClientKey = ""
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func (f *fakeDockerHostsRepo) Update(_ context.Context, id uuid.UUID, url string, caCert, clientCert, clientKey *string) error {
	host, ok := f.hosts[id]
	if !ok {
		return errors.New("not found")
	}

	host.Url = url

	if caCert != nil {
		host.CaCert = *caCert
	}

	if clientCert != nil {
		host.ClientCert = *clientCert
	}

	if clientKey != nil {
		host.ClientKey = *clientKey
	}

	f.hosts[id] = host

	return nil
}

func (f *fakeDockerHostsRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(f.hosts, id)

	return nil
}

func (f *fakeDockerHostsRepo) Exists(_ context.Context) (bool, error) {
	if f.existsErr != nil {
		return false, f.existsErr
	}

	return f.existsResult, nil
}

func (f *fakeDockerHostsRepo) PickLeastLoaded(_ context.Context) (domain.DockerHost, error) {
	for _, host := range f.hosts {
		return host, nil
	}

	return domain.DockerHost{}, errors.New("no hosts")
}

func (f *fakeDockerHostsRepo) WithTx(_ sqldb.DB) repository.DockerHosts {
	return f
}

var _ repository.DockerHosts = (*fakeDockerHostsRepo)(nil)

func TestService_RegisterDockerHost(t *testing.T) {
	fakeRepo := newFakeDockerHostsRepo()
	svc := &Service{dockerHostsRepo: fakeRepo}

	id, err := svc.RegisterDockerHost(context.Background(), "tcp://host:2376", "", "", "")
	require.NoError(t, err)
	require.NotEmpty(t, id)

	uid, err := uuid.Parse(id)
	require.NoError(t, err)
	require.Equal(t, "tcp://host:2376", fakeRepo.hosts[uid].Url)
}

func TestService_GetDockerHost(t *testing.T) {
	fakeRepo := newFakeDockerHostsRepo()
	svc := &Service{dockerHostsRepo: fakeRepo}

	id, err := svc.RegisterDockerHost(context.Background(), "tcp://host:2376", "", "", "")
	require.NoError(t, err)

	t.Run("returns the registered host", func(t *testing.T) {
		got, err := svc.GetDockerHost(context.Background(), id)
		require.NoError(t, err)
		require.Equal(t, "tcp://host:2376", got.Url)
	})

	t.Run("propagates a bad uuid", func(t *testing.T) {
		_, err := svc.GetDockerHost(context.Background(), "not-a-uuid")
		require.Error(t, err)
	})

	t.Run("propagates a repo lookup error", func(t *testing.T) {
		_, err := svc.GetDockerHost(context.Background(), uuid.New().String())
		require.Error(t, err)
	})
}

func TestService_ListDockerHosts(t *testing.T) {
	fakeRepo := newFakeDockerHostsRepo()
	svc := &Service{dockerHostsRepo: fakeRepo}

	_, err := svc.RegisterDockerHost(context.Background(), "tcp://host-a:2376", "", "", "")
	require.NoError(t, err)
	_, err = svc.RegisterDockerHost(context.Background(), "tcp://host-b:2376", "", "", "")
	require.NoError(t, err)

	hosts, err := svc.ListDockerHosts(context.Background())
	require.NoError(t, err)
	require.Len(t, hosts, 2)
}

func TestService_UpdateDockerHost(t *testing.T) {
	fakeRepo := newFakeDockerHostsRepo()
	svc := &Service{dockerHostsRepo: fakeRepo}

	id, err := svc.RegisterDockerHost(context.Background(), "tcp://host:2376", "", "", "")
	require.NoError(t, err)

	t.Run("updates the url", func(t *testing.T) {
		err := svc.UpdateDockerHost(context.Background(), id, "tcp://new-host:2376", nil, nil, nil)
		require.NoError(t, err)

		got, err := svc.GetDockerHost(context.Background(), id)
		require.NoError(t, err)
		require.Equal(t, "tcp://new-host:2376", got.Url)
	})

	t.Run("propagates a bad uuid", func(t *testing.T) {
		err := svc.UpdateDockerHost(context.Background(), "not-a-uuid", "tcp://host:2376", nil, nil, nil)
		require.Error(t, err)
	})
}

func TestService_DeleteDockerHost(t *testing.T) {
	fakeRepo := newFakeDockerHostsRepo()
	svc := &Service{dockerHostsRepo: fakeRepo}

	id, err := svc.RegisterDockerHost(context.Background(), "tcp://host:2376", "", "", "")
	require.NoError(t, err)

	t.Run("deletes the host", func(t *testing.T) {
		err := svc.DeleteDockerHost(context.Background(), id)
		require.NoError(t, err)

		_, err = svc.GetDockerHost(context.Background(), id)
		require.Error(t, err)
	})

	t.Run("propagates a bad uuid", func(t *testing.T) {
		err := svc.DeleteDockerHost(context.Background(), "not-a-uuid")
		require.Error(t, err)
	})
}

func TestService_HasDockerHosts(t *testing.T) {
	t.Run("returns true when hosts exist", func(t *testing.T) {
		fakeRepo := newFakeDockerHostsRepo()
		fakeRepo.existsResult = true

		svc := &Service{dockerHostsRepo: fakeRepo}

		exists, err := svc.HasDockerHosts(context.Background())
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("returns false when no hosts exist", func(t *testing.T) {
		fakeRepo := newFakeDockerHostsRepo()
		fakeRepo.existsResult = false

		svc := &Service{dockerHostsRepo: fakeRepo}

		exists, err := svc.HasDockerHosts(context.Background())
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("propagates wrapped repo error", func(t *testing.T) {
		fakeRepo := newFakeDockerHostsRepo()
		fakeRepo.existsErr = errors.New("db unavailable")

		svc := &Service{dockerHostsRepo: fakeRepo}

		exists, err := svc.HasDockerHosts(context.Background())
		require.Error(t, err)
		require.False(t, exists)
		require.ErrorIs(t, err, fakeRepo.existsErr)
	})
}
