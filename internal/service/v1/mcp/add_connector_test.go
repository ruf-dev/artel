package mcp

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// --- hand-rolled fakes (this package has no mock-generation setup) ---

type fakeMcpKeyRepo struct {
	key domain.McpKey
}

func (f *fakeMcpKeyRepo) CreateMcpKey(
	context.Context, uuid.UUID, uuid.UUID, uuid.UUID, string, []byte, string,
) (domain.McpKey, error) {
	panic("not implemented")
}
func (f *fakeMcpKeyRepo) ListMcpKeysByVault(context.Context, uuid.UUID) ([]domain.McpKey, error) {
	panic("not implemented")
}
func (f *fakeMcpKeyRepo) ListMcpKeysByUser(context.Context, uuid.UUID) ([]domain.McpKey, error) {
	panic("not implemented")
}
func (f *fakeMcpKeyRepo) GetMcpKeyByID(context.Context, uuid.UUID) (domain.McpKey, error) {
	return f.key, nil
}
func (f *fakeMcpKeyRepo) ListActiveMcpKeys(context.Context, uuid.UUID) ([]domain.McpKey, error) {
	panic("not implemented")
}
func (f *fakeMcpKeyRepo) RevokeMcpKey(context.Context, uuid.UUID) error { panic("not implemented") }
func (f *fakeMcpKeyRepo) SetMcpKeyAccess(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error {
	panic("not implemented")
}
func (f *fakeMcpKeyRepo) TouchLastAccessed(context.Context, uuid.UUID) error {
	panic("not implemented")
}

type fakeMcpDefinitionsRepo struct{}

func (f *fakeMcpDefinitionsRepo) Upsert(
	context.Context, domain.McpDefinition,
) (domain.McpDefinition, error) {
	panic("not implemented")
}
func (f *fakeMcpDefinitionsRepo) Get(
	_ context.Context, name string,
) (sql.Null[domain.McpDefinition], error) {
	return sql.Null[domain.McpDefinition]{V: domain.McpDefinition{Name: name}, Valid: true}, nil
}
func (f *fakeMcpDefinitionsRepo) List(context.Context) ([]domain.McpDefinition, error) {
	panic("not implemented")
}
func (f *fakeMcpDefinitionsRepo) Delete(context.Context, string) error { panic("not implemented") }
func (f *fakeMcpDefinitionsRepo) GetTool(
	context.Context, string, string,
) (sql.Null[domain.McpToolDef], error) {
	panic("not implemented")
}
func (f *fakeMcpDefinitionsRepo) ListAllTools(context.Context) ([]domain.McpToolRef, error) {
	panic("not implemented")
}

type fakeMcpConnectorsRepo struct {
	inserted bool
}

func (f *fakeMcpConnectorsRepo) Insert(
	_ context.Context, connector domain.McpConnector,
) (domain.McpConnector, error) {
	f.inserted = true
	connector.Uuid = uuid.New()

	return connector, nil
}
func (f *fakeMcpConnectorsRepo) Get(
	context.Context, uuid.UUID, string,
) (sql.Null[domain.McpConnector], error) {
	panic("not implemented")
}
func (f *fakeMcpConnectorsRepo) ListByKey(
	context.Context, uuid.UUID,
) ([]domain.McpConnector, error) {
	panic("not implemented")
}
func (f *fakeMcpConnectorsRepo) Delete(context.Context, uuid.UUID, string) error {
	panic("not implemented")
}

type fakeExternalConnRepo struct {
	conn domain.ExternalConnection
}

func (f *fakeExternalConnRepo) Upsert(
	context.Context, domain.ExternalConnection,
) (domain.ExternalConnection, error) {
	panic("not implemented")
}
func (f *fakeExternalConnRepo) Insert(
	context.Context, domain.ExternalConnection,
) (domain.ExternalConnection, error) {
	panic("not implemented")
}
func (f *fakeExternalConnRepo) GetByID(
	context.Context, uuid.UUID,
) (domain.ExternalConnection, error) {
	return f.conn, nil
}
func (f *fakeExternalConnRepo) GetByUserAndProvider(
	context.Context, uuid.UUID, string,
) (sql.Null[domain.ExternalConnection], error) {
	panic("not implemented")
}
func (f *fakeExternalConnRepo) ListByUser(
	context.Context, uuid.UUID,
) ([]domain.ExternalConnection, error) {
	panic("not implemented")
}
func (f *fakeExternalConnRepo) Delete(context.Context, uuid.UUID, string) error {
	panic("not implemented")
}
func (f *fakeExternalConnRepo) DeleteByID(context.Context, uuid.UUID, uuid.UUID) error {
	panic("not implemented")
}

type fakeSubscriptions struct {
	checkFeatureErr error
}

func (f *fakeSubscriptions) CheckActive(context.Context, uuid.UUID) error { return nil }
func (f *fakeSubscriptions) HasFeature(
	context.Context, uuid.UUID, domain.SubscriptionFeature,
) (bool, error) {
	return f.checkFeatureErr == nil, f.checkFeatureErr
}
func (f *fakeSubscriptions) CheckFeature(
	context.Context, uuid.UUID, domain.SubscriptionFeature,
) error {
	return f.checkFeatureErr
}
func (f *fakeSubscriptions) GetEffective(
	context.Context, uuid.UUID,
) (domain.EffectiveSubscription, error) {
	panic("not implemented")
}
func (f *fakeSubscriptions) GetUsage(context.Context, uuid.UUID) (domain.StorageUsage, error) {
	panic("not implemented")
}
func (f *fakeSubscriptions) CheckStorageQuota(context.Context, uuid.UUID) error {
	panic("not implemented")
}
func (f *fakeSubscriptions) CheckSkillLimit(context.Context, uuid.UUID, bool) error {
	panic("not implemented")
}

func TestAddConnector_NotifyFeatureGate(t *testing.T) {
	userUuid := uuid.New()
	keyUuid := uuid.New()
	connUuid := uuid.New()

	cases := []struct {
		name            string
		mcpName         string
		checkFeatureErr error
		wantErr         error
		wantInserted    bool
	}{
		{
			name:            "telegram + feature absent is refused",
			mcpName:         domain.ProviderTelegram,
			checkFeatureErr: user_errors.FeatureNotEnabled,
			wantErr:         user_errors.FeatureNotEnabled,
			wantInserted:    false,
		},
		{
			name:         "telegram + feature present proceeds",
			mcpName:      domain.ProviderTelegram,
			wantInserted: true,
		},
		{
			name:            "gitlab + feature absent still proceeds (not gated)",
			mcpName:         "gitlab",
			checkFeatureErr: user_errors.FeatureNotEnabled,
			wantInserted:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			connectors := &fakeMcpConnectorsRepo{}

			svc := &ServiceImpl{
				mcpKeys:             &fakeMcpKeyRepo{key: domain.McpKey{Uuid: keyUuid, UserUuid: userUuid}},
				mcpDefinitions:      &fakeMcpDefinitionsRepo{},
				mcpConnectors:       connectors,
				externalConnections: &fakeExternalConnRepo{conn: domain.ExternalConnection{Uuid: connUuid, UserUuid: userUuid}},
				subscriptions:       &fakeSubscriptions{checkFeatureErr: tc.checkFeatureErr},
			}

			ctx := user_context.WithUserContext(
				context.Background(), user_context.UserContext{UserUuid: userUuid, UserName: "tester"},
			)

			_, err := svc.AddConnector(ctx, keyUuid, tc.mcpName, connUuid)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.wantInserted, connectors.inserted)
		})
	}
}
