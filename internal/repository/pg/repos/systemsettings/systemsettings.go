package systemsettings

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) WithTx(tx *sql.Tx) repository.SystemSettingsRepo {
	return &Repo{q: r.q.WithTx(tx)}
}

func (r *Repo) Get(ctx context.Context) (domain.SystemSettings, error) {
	row, err := r.q.GetSystemSettings(ctx)
	if err != nil {
		return domain.SystemSettings{}, rerrors.Wrap(err, "error getting system settings")
	}

	return settingsFromRow(row), nil
}

func (r *Repo) GetForUpdate(ctx context.Context) (domain.SystemSettings, error) {
	row, err := r.q.GetSystemSettingsForUpdate(ctx)
	if err != nil {
		return domain.SystemSettings{}, rerrors.Wrap(err, "error getting system settings for update")
	}

	return settingsFromRow(row), nil
}

func (r *Repo) UpdateAuthMethods(ctx context.Context, passwordEnabled, telegramEnabled bool) error {
	params := artel_q.UpdateSystemSettingsAuthMethodsParams{
		PasswordAuthEnabled: passwordEnabled,
		TelegramAuthEnabled: telegramEnabled,
	}

	err := r.q.UpdateSystemSettingsAuthMethods(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error updating system settings auth methods")
	}

	return nil
}

func (r *Repo) UpdateRegistrationMode(ctx context.Context, mode domain.RegistrationMode) error {
	err := r.q.UpdateSystemSettingsRegistrationMode(ctx, string(mode))
	if err != nil {
		return rerrors.Wrap(err, "error updating system settings registration mode")
	}

	return nil
}

func (r *Repo) SetSetupToken(ctx context.Context, tokenHash string, issuedAt time.Time) error {
	params := artel_q.SetSetupTokenParams{
		SetupTokenHash:     sql.NullString{String: tokenHash, Valid: tokenHash != ""},
		SetupTokenIssuedAt: sql.NullTime{Time: issuedAt, Valid: !issuedAt.IsZero()},
	}

	err := r.q.SetSetupToken(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error setting setup token")
	}

	return nil
}

func (r *Repo) CompleteSetup(ctx context.Context) error {
	err := r.q.CompleteSetup(ctx)
	if err != nil {
		return rerrors.Wrap(err, "error completing setup")
	}

	return nil
}

func (r *Repo) UpdateDefaultDocsVault(ctx context.Context, vaultUuid *uuid.UUID) error {
	param := uuid.NullUUID{Valid: vaultUuid != nil}
	if vaultUuid != nil {
		param.UUID = *vaultUuid
	}

	err := r.q.UpdateSystemSettingsDefaultDocsVault(ctx, param)
	if err != nil {
		return rerrors.Wrap(err, "error updating system settings default docs vault")
	}

	return nil
}

func (r *Repo) UpdateDefaultDocsSource(ctx context.Context, source domain.DocsSource) error {
	err := r.q.UpdateSystemSettingsDefaultDocsSource(ctx, string(source))
	if err != nil {
		return rerrors.Wrap(err, "error updating system settings default docs source")
	}

	return nil
}

func (r *Repo) UpdateSystemPrompt(ctx context.Context, prompt string) error {
	err := r.q.UpdateSystemPrompt(ctx, prompt)
	if err != nil {
		return rerrors.Wrap(err, "error updating system prompt")
	}
	return nil
}

func settingsFromRow(row artel_q.SystemSetting) domain.SystemSettings {
	settings := domain.SystemSettings{
		SystemPrompt:        row.SystemPrompt,
		SetupCompleted:      row.SetupCompleted,
		PasswordAuthEnabled: row.PasswordAuthEnabled,
		TelegramAuthEnabled: row.TelegramAuthEnabled,
		RegistrationMode:    domain.RegistrationMode(row.RegistrationMode),
		SetupTokenHash:      row.SetupTokenHash.String,
		SetupTokenIssuedAt:  row.SetupTokenIssuedAt.Time,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
		DefaultDocsSource:   domain.DocsSource(row.DefaultDocsSource),
	}

	if row.DefaultDocsVaultID.Valid {
		vaultUuid := row.DefaultDocsVaultID.UUID
		settings.DefaultDocsVaultID = &vaultUuid
	}

	return settings
}
