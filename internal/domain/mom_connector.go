package domain

import (
	"time"

	"github.com/google/uuid"
)

// McpConnector links an McpKey to an McpDefinition via an ExternalConnection.
// Stored in the mcp_connectors table.
type McpConnector struct {
	Uuid                   uuid.UUID
	McpKeyUuid             uuid.UUID
	McpName                string
	ExternalConnectionUuid uuid.UUID
	CreatedAt              time.Time
}
