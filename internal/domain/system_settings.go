package domain

import "time"

// RegistrationMode controls who may self-register a new account versus requiring an
// administrator to create it. See SystemSettings.
type RegistrationMode string

const (
	RegistrationModeAdminOnly    RegistrationMode = "admin_only"
	RegistrationModeSelfRegister RegistrationMode = "self_register"
)

// SystemSettings is the single-row (id=1) global configuration for the first-run setup wizard
// and instance-wide auth policy — see migrations/064_system_settings.sql.
type SystemSettings struct {
	SetupCompleted      bool
	PasswordAuthEnabled bool
	TelegramAuthEnabled bool
	RegistrationMode    RegistrationMode

	// SetupTokenHash/SetupTokenIssuedAt back the one-time setup token minted on first run to
	// authorize the initial admin creation. Empty/zero once CompleteSetup has cleared them.
	SetupTokenHash     string
	SetupTokenIssuedAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
