package domain

import (
	"time"

	"github.com/google/uuid"
)

// TaskTracker is the task_trackers_api wire-shape view over a trello external_connections row —
// backed entirely by external_connections + MoM (see internal/service/v1/tasktracker), not by a
// dedicated table. Credentials never surface here; they stay inside external_connections'
// encrypted CredentialsJSON and are only touched by the MoM http executor.
type TaskTracker struct {
	Uuid      uuid.UUID
	UserUuid  uuid.UUID
	Type      string
	Name      string
	CreatedAt time.Time
}
