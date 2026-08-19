package domain

import (
	"time"

	"github.com/google/uuid"
)

type WorkbenchStatus string

const (
	WorkbenchStatusCreated     WorkbenchStatus = "created"
	WorkbenchStatusConfiguring WorkbenchStatus = "configuring"
	WorkbenchStatusRunning     WorkbenchStatus = "running"
	WorkbenchStatusStopped     WorkbenchStatus = "stopped"
	WorkbenchStatusRemoved     WorkbenchStatus = "removed"
)

type WorkbenchAuthMode string

const (
	WorkbenchAuthModeAPIKey            WorkbenchAuthMode = "api_key"
	WorkbenchAuthModeSubscriptionLogin WorkbenchAuthMode = "subscription_login"
)

// Workbench is a thin reflection of the Docker container backing a vault's cloud workbench;
// Status moves through the WorkbenchStatus values above as CreateWorkbench/StartWorkbench/
// StopWorkbench/DeleteWorkbench (internal/service/v1/workbench) proceed.
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
	// DockerHostUuid is the docker_hosts row this workbench's container/volume live on, assigned
	// at CreateWorkbench time by picking the least-loaded host — see
	// internal/service/v1/workbench/workbench.go's resolveClient. nil only for rows created
	// before the docker_hosts pool existed.
	DockerHostUuid *uuid.UUID
}

// TerminalTab is one tmux window inside a workbench's tmux session, surfaced to the browser
// as a terminal tab. ID is tmux's own #{window_id} (e.g. "@1") — stable across window
// renumbering, unlike #{window_index}.
type TerminalTab struct {
	ID     string
	Name   string
	Active bool
}
