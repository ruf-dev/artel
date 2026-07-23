package domain

import (
	"time"

	"github.com/google/uuid"
)

type WorkbenchStatus string

const (
	WorkbenchStatusCreated WorkbenchStatus = "created"
	WorkbenchStatusRunning WorkbenchStatus = "running"
	WorkbenchStatusStopped WorkbenchStatus = "stopped"
	WorkbenchStatusRemoved WorkbenchStatus = "removed"
)

type WorkbenchAuthMode string

const (
	WorkbenchAuthModeAPIKey            WorkbenchAuthMode = "api_key"
	WorkbenchAuthModeSubscriptionLogin WorkbenchAuthMode = "subscription_login"
)

// WorkbenchLoginState is the state of the in-container subscription_login TUI flow, derived by
// parsing a tmux capture-pane snapshot of the workbench's session — see
// docs/workbench/03_auth_and_login_flow.md, "Mechanism (confirmed)".
type WorkbenchLoginState string

const (
	// WorkbenchLoginStatePending: claude hasn't printed the OAuth authorize URL yet (e.g. the
	// one-time theme-selection screen or the login-method menu is still showing). Nothing
	// actionable for the caller besides "keep polling".
	WorkbenchLoginStatePending WorkbenchLoginState = "pending"
	// WorkbenchLoginStateURLPresent: the OAuth authorize URL has been printed and claude is
	// waiting for a pasted code. URL is populated.
	WorkbenchLoginStateURLPresent WorkbenchLoginState = "url_present"
	// WorkbenchLoginStateError: a previously submitted code was rejected ("OAuth error: ...").
	// ErrorMessage is populated; the caller should let the user retry.
	WorkbenchLoginStateError WorkbenchLoginState = "error"
	// WorkbenchLoginStateAuthorized: neither the login prompt/menu nor an error is visible
	// anymore — login completed.
	WorkbenchLoginStateAuthorized WorkbenchLoginState = "authorized"
)

// WorkbenchLoginPrompt is a snapshot of the subscription_login TUI flow's current state. URL is
// only meaningful when State == WorkbenchLoginStateURLPresent; ErrorMessage only when
// State == WorkbenchLoginStateError.
type WorkbenchLoginPrompt struct {
	State        WorkbenchLoginState
	URL          string
	ErrorMessage string
}

// Workbench is a thin reflection of the Docker container backing a vault's cloud workbench —
// see docs/workbench/01_data_model_and_lifecycle.md for the full state machine.
type Workbench struct {
	Uuid        uuid.UUID
	VaultUuid   uuid.UUID
	UserUuid    uuid.UUID
	Status      WorkbenchStatus
	AuthMode    WorkbenchAuthMode // "" when not yet started
	ContainerId string            // "" only in the creation-failed/retry edge case
	VolumeName  string
	CreatedAt   time.Time
	StartedAt   *time.Time // nil until started
	StoppedAt   *time.Time // nil until stopped
}
