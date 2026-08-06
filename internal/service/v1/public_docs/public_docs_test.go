package public_docs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// fakeVaults is a hand-written repository.Vaults exposing only GetBySlug's behavior as
// configurable — the rest of the interface is unused by Service.
type fakeVaults struct {
	getBySlug func(ctx context.Context, slug string) (domain.Vault, error)
}

func (f *fakeVaults) Upsert(context.Context, uuid.UUID, uuid.UUID, string, string, string, string) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) GetByID(context.Context, uuid.UUID) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) GetByNameAndUser(context.Context, uuid.UUID, string) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) UpdateStatus(context.Context, uuid.UUID, string) error { return nil }

func (f *fakeVaults) SetLiveSyncPassphrase(context.Context, uuid.UUID, string) error { return nil }

func (f *fakeVaults) ListByMembership(context.Context, uuid.UUID) ([]domain.Vault, error) {
	return nil, nil
}

func (f *fakeVaults) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeVaults) LinkS3Bucket(context.Context, uuid.UUID, uuid.UUID, string) error { return nil }

func (f *fakeVaults) UnlinkS3Bucket(context.Context, uuid.UUID) error { return nil }

func (f *fakeVaults) SetUseCouchDBForBinaries(context.Context, uuid.UUID, bool) error { return nil }

func (f *fakeVaults) PublishVault(context.Context, uuid.UUID, string) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) UnpublishVault(context.Context, uuid.UUID) error { return nil }

func (f *fakeVaults) GetBySlug(ctx context.Context, slug string) (domain.Vault, error) {
	return f.getBySlug(ctx, slug)
}

func (f *fakeVaults) WithTx(postgres.DB) repository.Vaults { return f }

// fakeCouchInstances is a hand-written repository.CouchInstances exposing only Get's behavior
// as configurable — the rest of the interface is unused by Service.
type fakeCouchInstances struct {
	get func(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error)
}

func (f *fakeCouchInstances) Register(context.Context, string, string, []byte) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (f *fakeCouchInstances) Get(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error) {
	if f.get == nil {
		return domain.CouchInstance{Uuid: id}, nil
	}

	return f.get(ctx, id)
}

func (f *fakeCouchInstances) RandomPick(context.Context) (domain.CouchInstanceWithAccount, error) {
	return domain.CouchInstanceWithAccount{}, nil
}

func (f *fakeCouchInstances) List(context.Context) ([]domain.CouchInstance, error) { return nil, nil }

func (f *fakeCouchInstances) Update(context.Context, uuid.UUID, string, string, []byte) error {
	return nil
}

func (f *fakeCouchInstances) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeCouchInstances) Exists(context.Context) (bool, error) { return false, nil }

func (f *fakeCouchInstances) WithTx(postgres.DB) repository.CouchInstances { return f }

func newTestService(vaults *fakeVaults, instances *fakeCouchInstances) *Service {
	return &Service{
		vaults:         vaults,
		couchInstances: instances,
	}
}

func TestGetVaultBySlug_PublicVaultSucceeds(t *testing.T) {
	vault := domain.Vault{Uuid: uuid.New(), Name: "Public Vault", IsPublic: true, Slug: "public-vault"}

	vaults := &fakeVaults{
		getBySlug: func(context.Context, string) (domain.Vault, error) {
			return vault, nil
		},
	}

	svc := newTestService(vaults, &fakeCouchInstances{})

	// Deliberately plain context.Background() — this path must work with no user_context at
	// all, since PublicDocsAPI callers are genuinely unauthenticated.
	got, err := svc.GetVaultBySlug(context.Background(), "public-vault")
	require.NoError(t, err)
	require.Equal(t, vault.Uuid, got.Uuid)
	require.Equal(t, vault.Name, got.Name)
}

func TestGetVaultBySlug_NonPublicVaultNotFound(t *testing.T) {
	vault := domain.Vault{Uuid: uuid.New(), Name: "Private Vault", IsPublic: false, Slug: "private-vault"}

	vaults := &fakeVaults{
		getBySlug: func(context.Context, string) (domain.Vault, error) {
			return vault, nil
		},
	}

	svc := newTestService(vaults, &fakeCouchInstances{})

	_, err := svc.GetVaultBySlug(context.Background(), "private-vault")
	require.ErrorIs(t, err, user_errors.NotFound)
}

func TestGetVaultBySlug_UnknownSlugNotFound(t *testing.T) {
	vaults := &fakeVaults{
		getBySlug: func(context.Context, string) (domain.Vault, error) {
			return domain.Vault{}, rerrors.Wrap(user_errors.NotFound)
		},
	}

	svc := newTestService(vaults, &fakeCouchInstances{})

	_, err := svc.GetVaultBySlug(context.Background(), "does-not-exist")
	require.ErrorIs(t, err, user_errors.NotFound)
}

// newFakeCouchServer starts an httptest.Server that answers just enough of the CouchDB HTTP API
// for LiveSyncClient.ListNotes (via kivik's AllDocs) to succeed against a single fixture note
// document — exercising the real anonymous-read path (admin instance credentials, no
// user_context) end-to-end instead of stopping at vault resolution.
func newFakeCouchServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]any{"couchdb": "Welcome", "version": "3.0.0"}

		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	})

	mux.HandleFunc("/testdb/_all_docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		row := map[string]any{
			"id":  "note.md",
			"key": "note.md",
			"value": map[string]any{
				"rev": "1-abc",
			},
			"doc": map[string]any{
				"_id":   "note.md",
				"_rev":  "1-abc",
				"type":  "plain",
				"mtime": 123,
				"size":  5,
			},
		}
		response := map[string]any{
			"total_rows": 1,
			"offset":     0,
			"rows":       []any{row},
		}

		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	})

	mux.HandleFunc("/testdb/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]any{"db_name": "testdb"}

		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	})

	return httptest.NewServer(mux)
}

func TestListNotes_PublicVaultSucceeds(t *testing.T) {
	server := newFakeCouchServer(t)
	t.Cleanup(server.Close)

	instanceID := uuid.New()
	vault := domain.Vault{
		Uuid:              uuid.New(),
		CouchInstanceUuid: instanceID,
		CouchDBName:       "testdb",
		IsPublic:          true,
		Slug:              "public-vault",
	}

	vaults := &fakeVaults{
		getBySlug: func(context.Context, string) (domain.Vault, error) {
			return vault, nil
		},
	}
	instances := &fakeCouchInstances{
		get: func(context.Context, uuid.UUID) (domain.CouchInstance, error) {
			return domain.CouchInstance{
				Uuid:     instanceID,
				Url:      strings.TrimSuffix(server.URL, "/"),
				Username: "admin",
				Password: "admin",
			}, nil
		},
	}

	svc := newTestService(vaults, instances)

	notes, err := svc.ListNotes(context.Background(), "public-vault")
	require.NoError(t, err)
	require.Len(t, notes, 1)
	require.Equal(t, "note.md", notes[0].Path)
}

func TestListFolders_NonPublicVaultNotFound(t *testing.T) {
	vault := domain.Vault{Uuid: uuid.New(), IsPublic: false, Slug: "private-vault"}

	vaults := &fakeVaults{
		getBySlug: func(context.Context, string) (domain.Vault, error) {
			return vault, nil
		},
	}

	svc := newTestService(vaults, &fakeCouchInstances{})

	// No user_context anywhere in this test — confirms the read path never depends on one.
	_, err := svc.ListFolders(context.Background(), "private-vault")
	require.ErrorIs(t, err, user_errors.NotFound)
}
