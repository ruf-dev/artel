package externalconnections

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// fakeUsersRepo is a hand-rolled fake of repository.Users, scoped to GetTelegramIdByUserId.
type fakeUsersRepo struct {
	telegramIdFunc func(ctx context.Context, userUuid uuid.UUID) (sql.Null[string], error)
}

func (f *fakeUsersRepo) GetTelegramIdByUserId(
	ctx context.Context, userUuid uuid.UUID,
) (sql.Null[string], error) {
	return f.telegramIdFunc(ctx, userUuid)
}

func (f *fakeUsersRepo) Create(context.Context, string, string) (domain.User, error) {
	panic("not implemented")
}
func (f *fakeUsersRepo) GetByID(context.Context, uuid.UUID) (domain.User, error) {
	panic("not implemented")
}
func (f *fakeUsersRepo) GetByEmail(context.Context, string) (domain.User, error) {
	panic("not implemented")
}
func (f *fakeUsersRepo) FindByEmail(context.Context, string) (sql.Null[domain.User], error) {
	panic("not implemented")
}
func (f *fakeUsersRepo) GetByTelegramId(context.Context, string) (sql.Null[domain.User], error) {
	panic("not implemented")
}
func (f *fakeUsersRepo) CreateByUsername(context.Context, string, string) (domain.User, error) {
	panic("not implemented")
}
func (f *fakeUsersRepo) UpsertTelegramIdentity(context.Context, domain.TelegramIdentity) error {
	panic("not implemented")
}
func (f *fakeUsersRepo) GetTelegramPhotoUrl(context.Context, uuid.UUID) (string, error) {
	panic("not implemented")
}
func (f *fakeUsersRepo) UpdatePhotoUrl(context.Context, uuid.UUID, string) error {
	panic("not implemented")
}
func (f *fakeUsersRepo) UpdatePasswordHash(context.Context, uuid.UUID, string) error {
	panic("not implemented")
}
func (f *fakeUsersRepo) ListAll(context.Context, domain.ListUsersReq) ([]domain.User, int64, error) {
	panic("not implemented")
}
func (f *fakeUsersRepo) GetDetailsById(context.Context, uuid.UUID) (domain.UserDetails, error) {
	panic("not implemented")
}
func (f *fakeUsersRepo) Delete(context.Context, uuid.UUID) error {
	panic("not implemented")
}
func (f *fakeUsersRepo) WithTx(*sql.Tx) repository.Users {
	panic("not implemented")
}

// fakeSubscriptionService is a hand-rolled fake of service.SubscriptionService, scoped to
// CheckFeature.
type fakeSubscriptionService struct {
	checkFeatureErr error
}

func (f *fakeSubscriptionService) CheckActive(context.Context, uuid.UUID) error { return nil }
func (f *fakeSubscriptionService) HasFeature(
	context.Context, uuid.UUID, domain.SubscriptionFeature,
) (bool, error) {
	return f.checkFeatureErr == nil, f.checkFeatureErr
}
func (f *fakeSubscriptionService) CheckFeature(
	context.Context, uuid.UUID, domain.SubscriptionFeature,
) error {
	return f.checkFeatureErr
}
func (f *fakeSubscriptionService) GetEffective(
	context.Context, uuid.UUID,
) (domain.EffectiveSubscription, error) {
	panic("not implemented")
}
func (f *fakeSubscriptionService) GetUsage(context.Context, uuid.UUID) (domain.StorageUsage, error) {
	panic("not implemented")
}
func (f *fakeSubscriptionService) CheckStorageQuota(context.Context, uuid.UUID) error {
	panic("not implemented")
}
func (f *fakeSubscriptionService) CheckSkillLimit(context.Context, uuid.UUID, bool) error {
	panic("not implemented")
}

const testArtelBotToken = "artel-bot-token-xyz"

func TestEnsureArtelTelegramConnection(t *testing.T) {
	existingRowUuid := uuid.New()

	cases := []struct {
		name            string
		telegramId      sql.Null[string]
		checkFeatureErr error
		upsertUuid      uuid.UUID
		wantErr         error
		wantChatId      int64
		wantUuid        uuid.UUID
	}{
		{
			name:       "telegram identity absent",
			telegramId: sql.Null[string]{Valid: false},
			upsertUuid: uuid.New(),
			wantErr:    user_errors.NotifyNoTelegramIdentity,
		},
		{
			name:            "feature not enabled",
			telegramId:      sql.Null[string]{V: "555", Valid: true},
			checkFeatureErr: user_errors.FeatureNotEnabled,
			upsertUuid:      uuid.New(),
			wantErr:         user_errors.FeatureNotEnabled,
		},
		{
			name:       "identity present, no existing managed row - inserts",
			telegramId: sql.Null[string]{V: "123456", Valid: true},
			upsertUuid: uuid.New(),
			wantChatId: 123456,
		},
		{
			name:       "identity present, existing managed row - updates in place, same uuid, refreshed chat id",
			telegramId: sql.Null[string]{V: "987654", Valid: true},
			upsertUuid: existingRowUuid,
			wantChatId: 987654,
			wantUuid:   existingRowUuid,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var captured domain.ExternalConnection

			repo := &fakeExternalConnectionRepo{
				upsertFunc: func(
					_ context.Context, conn domain.ExternalConnection,
				) (domain.ExternalConnection, error) {
					captured = conn
					conn.Uuid = tc.upsertUuid

					return conn, nil
				},
			}

			svc := &Service{
				connections:   repo,
				users:         &fakeUsersRepo{telegramIdFunc: func(context.Context, uuid.UUID) (sql.Null[string], error) { return tc.telegramId, nil }},
				subscriptions: &fakeSubscriptionService{checkFeatureErr: tc.checkFeatureErr},
				artelBotToken: testArtelBotToken,
			}

			gotUuid, err := svc.EnsureArtelTelegramConnection(withUser(uuid.New()))

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)

			if tc.wantUuid != uuid.Nil {
				require.Equal(t, tc.wantUuid, gotUuid)
			} else {
				require.Equal(t, tc.upsertUuid, gotUuid)
			}

			require.Equal(t, domain.ProviderTelegram, captured.Provider)

			var creds domain.TelegramCredentials

			require.NoError(t, json.Unmarshal(captured.CredentialsJSON, &creds))
			require.Equal(t, testArtelBotToken, creds.BotToken)
			require.Equal(t, tc.wantChatId, creds.ChatID)

			var meta telegramConnectionMeta

			require.NoError(t, json.Unmarshal(captured.Metadata, &meta))
			require.Equal(t, managedArtelBot, meta.Managed)
		})
	}
}

func TestEnsureArtelTelegramConnection_Unauthenticated(t *testing.T) {
	svc := &Service{}

	_, err := svc.EnsureArtelTelegramConnection(context.Background())
	require.ErrorIs(t, err, user_errors.Unauthenticated)
}
