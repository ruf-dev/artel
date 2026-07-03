package mcp

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/domain"
)

type fakeExternalConnectionRepo struct {
	conns []domain.ExternalConnection
}

func (f *fakeExternalConnectionRepo) Upsert(_ context.Context, conn domain.ExternalConnection) (domain.ExternalConnection, error) {
	return conn, nil
}

func (f *fakeExternalConnectionRepo) GetByID(_ context.Context, _ uuid.UUID) (domain.ExternalConnection, error) {
	return domain.ExternalConnection{}, nil
}

func (f *fakeExternalConnectionRepo) GetByUserAndProvider(_ context.Context, _ uuid.UUID, _ string) (sql.Null[domain.ExternalConnection], error) {
	return sql.Null[domain.ExternalConnection]{}, nil
}

func (f *fakeExternalConnectionRepo) ListByUser(_ context.Context, _ uuid.UUID) ([]domain.ExternalConnection, error) {
	return f.conns, nil
}

func (f *fakeExternalConnectionRepo) Delete(_ context.Context, _ uuid.UUID, _ string) error {
	return nil
}

func TestListConnectionsForTracts(t *testing.T) {
	connUuid := uuid.New()
	userUuid := uuid.New()

	repo := &fakeExternalConnectionRepo{
		conns: []domain.ExternalConnection{
			{Uuid: connUuid, UserUuid: userUuid, Provider: domain.ProviderGitlab},
		},
	}

	svc := New(nil, nil, nil, nil, nil, nil, repo)

	result, err := svc.listConnectionsForTracts(context.Background(), userUuid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rows []map[string]string
	err = json.Unmarshal([]byte(result.Text), &rows)
	if err != nil {
		t.Fatalf("error unmarshaling result: %v", err)
	}

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	if rows[0]["uuid"] != connUuid.String() {
		t.Fatalf("expected uuid %s, got %s", connUuid.String(), rows[0]["uuid"])
	}

	if rows[0]["provider"] != domain.ProviderGitlab {
		t.Fatalf("expected provider %s, got %s", domain.ProviderGitlab, rows[0]["provider"])
	}
}
