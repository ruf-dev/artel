package domain

import (
	"time"

	"github.com/google/uuid"
)

// RegistrationMode controls who may self-register a new account versus requiring an
// administrator to create it. See SystemSettings.
type RegistrationMode string

const (
	RegistrationModeAdminOnly    RegistrationMode = "admin_only"
	RegistrationModeSelfRegister RegistrationMode = "self_register"
)

// DocsSource selects what an unauthenticated visitor hitting `/docs` (no slug) is redirected to
// — either the admin-picked vault (DefaultDocsVaultID) or the built-in GitHub-backed quick-start
// guide. See public_docs.Service.GetDefaultVault.
type DocsSource string

const (
	DocsSourceVault  DocsSource = "vault"
	DocsSourceGithub DocsSource = "github"
)

// SystemSettings is the single-row (id=1) global configuration for the first-run setup wizard
// and instance-wide auth policy — see migrations/064_system_settings.sql and, for
// DefaultDocsVaultID, migrations/069_default_docs_vault.sql.
type SystemSettings struct {
	SystemPrompt        string
	SetupCompleted      bool
	PasswordAuthEnabled bool
	TelegramAuthEnabled bool
	RegistrationMode    RegistrationMode

	// SetupTokenHash/SetupTokenIssuedAt back the one-time setup token minted on first run to
	// authorize the initial admin creation. Empty/zero once CompleteSetup has cleared them.
	SetupTokenHash     string
	SetupTokenIssuedAt time.Time

	// DefaultDocsVaultID is the admin-configured vault an unauthenticated visitor hitting
	// `/docs` (no slug) is redirected to. Nil when unset — see public_docs.Service.GetDefaultVault.
	DefaultDocsVaultID *uuid.UUID

	// DefaultDocsSource selects between DefaultDocsVaultID and the built-in GitHub quick-start
	// guide for that same `/docs` redirect — see public_docs.Service.GetDefaultVault.
	DefaultDocsSource DocsSource

	CreatedAt time.Time
	UpdatedAt time.Time
}
