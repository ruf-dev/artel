package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Tract is a user-owned automation workflow. Definition is the entire step tree, persisted as
// a single JSONB document (not normalized rows) — matches the "UI/agents edit whole docs"
// convention: UpdateTract replaces Definition wholesale.
type Tract struct {
	Uuid        uuid.UUID
	UserUuid    uuid.UUID
	Name        string
	Description string
	Enabled     bool
	Definition  TractDefinition
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TractDefinition is the JSON-serialized step tree — json tags live here because the
// engine/validation layer (internal/service/v1/tract) marshals it directly, and
// internal/transport/tracts_api/to_proto.go round-trips it via json.Marshal/Unmarshal as a
// definition_json string field.
type TractDefinition struct {
	Steps []TractStep `json:"steps"`
}

// TractStep is one node of the nested step tree. Type is one of "action" | "condition" |
// "parallel" | "group" (constants defined locally in internal/service/v1/tract, not here).
// Action fields (Mcp/Tool/ConnectionUuid/Params) are only meaningful when Type == "action".
// Conditions only when Type == "condition", with Then/Else as its two branches. Parallel/group
// steps run/contain Steps (parallel: concurrently; group: sequentially, a plain nesting
// container). ConnectionUuid's zero value means no external connection required (builtin
// tools).
type TractStep struct {
	Id             string            `json:"id"`
	Name           string            `json:"name,omitempty"`
	Description    string            `json:"description,omitempty"`
	Type           string            `json:"type"`
	Mcp            string            `json:"mcp,omitempty"`
	Tool           string            `json:"tool,omitempty"`
	ConnectionUuid uuid.UUID         `json:"connection_uuid,omitempty"`
	Params         map[string]string `json:"params,omitempty"`
	Conditions     []TractCondition  `json:"conditions,omitempty"`
	Then           []TractStep       `json:"then,omitempty"`
	Else           []TractStep       `json:"else,omitempty"`
	Steps          []TractStep       `json:"steps,omitempty"`
}

// TractCondition is {left, op, right} — both sides are template strings rendered before
// comparison. json tags: consumed directly by internal/transport/tracts_api/to_proto.go's
// filtersFromJSON via json.Unmarshal.
type TractCondition struct {
	Left  string `json:"left"`
	Op    string `json:"op"`
	Right string `json:"right"`
}

type TractRunStatus string

const (
	TractRunRunning TractRunStatus = "running"
	TractRunDone    TractRunStatus = "done"
	TractRunFailed  TractRunStatus = "failed"
)

// TractRun is one execution of a tract, persisted before the step walk starts
// (persist-before-apply). TriggerUuid zero value = manual/mcp run (no trigger fired it).
type TractRun struct {
	Uuid           uuid.UUID
	TractUuid      uuid.UUID
	TriggerUuid    uuid.UUID
	StartedBy      string // "webhook" | "manual" | "mcp"
	Status         TractRunStatus
	TriggerPayload json.RawMessage // normalized payload / manual form values
	Error          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type TractRunStepStatus string

const (
	TractRunStepRunning TractRunStepStatus = "running"
	TractRunStepDone    TractRunStepStatus = "done"
	TractRunStepFailed  TractRunStepStatus = "failed"
)

// TractRunStep is one executed step within a run (audit trail + state store). StepId is the
// TractStep.Id string (steps live in a JSONB tree now, not their own rows, so there is no
// step row uuid to reference) — StepName/StepType are captured at execution time so history
// reads correctly even if the tract definition changes later.
type TractRunStep struct {
	Uuid       uuid.UUID
	RunUuid    uuid.UUID
	StepId     string
	StepName   string
	StepType   string
	Input      json.RawMessage
	Output     json.RawMessage
	Status     TractRunStepStatus
	Error      string
	StartedAt  time.Time
	FinishedAt time.Time // zero value = not finished yet
}

// Trigger is a standalone, user-owned trigger — reusable across multiple tracts via
// TriggerTractLink. Uuid is the stable primary key (used for owner-facing CRUD/API refs);
// TriggerUuid is a SEPARATE, independently rotatable id embedded in the public webhook URL —
// rotating it invalidates old webhook URLs without disturbing Uuid-keyed references
// (tract links, etc). SecretHash is sha256(raw webhook token); nil for manual triggers.
type Trigger struct {
	Uuid          uuid.UUID
	TriggerUuid   uuid.UUID
	UserUuid      uuid.UUID
	Name          string
	Kind          string // "webhook" | "manual"
	Source        string // "gitlab_push" | "generic"
	Config        json.RawMessage
	PayloadSchema ToolSchema
	SecretHash    []byte
	Enabled       bool
	CreatedAt     time.Time
}

// TriggerTractLink links triggerUuid to tractUuid with filter conditions (AND semantics) that
// gate whether a given webhook delivery starts a run for that specific tract.
type TriggerTractLink struct {
	TriggerUuid uuid.UUID
	TractUuid   uuid.UUID
	Filters     []TractCondition
}

// TriggerSourcePreset describes one webhook trigger source (e.g. "gitlab_push") — read-only
// catalog type, not persisted.
type TriggerSourcePreset struct {
	Key           string
	Description   string
	PayloadSchema ToolSchema
}
