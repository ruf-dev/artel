package domain

import (
	"time"

	"github.com/google/uuid"
)

type McpSpreadsheet struct {
	Uuid                 uuid.UUID
	UserUuid             uuid.UUID
	ExternalConnectionId uuid.UUID
	SpreadsheetId        string
	Name                 string
	CreatedAt            time.Time
}
