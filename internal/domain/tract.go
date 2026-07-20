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

// TractTemplate is an immutable published snapshot of a Tract, browsable/copyable by any user
// (not just its owner). Definition is copied at publish time — it is NOT a live view of the
// source tract; later edits to the source tract never propagate here. Every ConnectionUuid
// inside Definition's step tree is always zero: those uuids referenced per-user
// external_connections rows belonging to the publisher's account, meaningless (and unsafe to
// expose) to anyone who installs the template into their own account, so they are stripped to
// nil at publish time. SourceTractUuid is zero if the source tract has since been deleted
// (source_tract_id is ON DELETE SET NULL — the template snapshot is designed to outlive it).
type TractTemplate struct {
	Uuid            uuid.UUID
	SourceTractUuid uuid.UUID
	OwnerUuid       uuid.UUID
	Name            string
	Description     string
	Definition      TractDefinition
	Category        string
	InstallCount    int
	PublishedAt     time.Time
	UpdatedAt       time.Time
}

// TractDefinition is the JSON-serialized step tree — json tags live here because the
// engine/validation layer (internal/service/v1/tract) marshals it directly for Postgres JSONB
// storage. The wire (proto) shape is a typed oneof-per-step-kind message, mapped to/from this
// flat struct by internal/transport/tracts_api/to_proto.go's stepToProto/stepFromProto.
type TractDefinition struct {
	Steps []TractStep `json:"steps"`
}

// TractStep is one node of the nested step tree. Type is one of "action" | "condition" |
// "parallel" | "group" | "script" (constants defined locally in internal/service/v1/tract,
// not here). Action fields (Mcp/Tool/ConnectionUuid/Params) are only meaningful when Type ==
// "action". Conditions only when Type == "condition", with Then/Else as its two branches.
// Parallel/group steps run/contain Steps (parallel: concurrently; group: sequentially, a
// plain nesting container). ConnectionUuid's zero value means no external connection
// required (builtin tools). Script fields (Language/Code/InputParams/OutputParams) are only
// meaningful when Type == "script" — Params is reused there too, as the template-expression
// binding for each declared InputParams entry (same role it plays for action steps).
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
	Language       ScriptLanguage    `json:"language,omitempty"`
	Code           string            `json:"code,omitempty"`
	InputParams    []ScriptParam     `json:"input_params,omitempty"`
	OutputParams   []ScriptParam     `json:"output_params,omitempty"`
}

// ScriptLanguage is the engine id a script step runs under — a closed set, mapped to/from a
// real numbered proto enum at the transport boundary (see internal/transport/tracts_api).
type ScriptLanguage string

const (
	ScriptLanguageUnspecified ScriptLanguage = ""
	ScriptLanguageJavaScript  ScriptLanguage = "javascript"
)

// ScriptParam is one named, typed, ordered entry in a script step's declared input or
// output list — ordered (unlike ToolSchema.Properties, a map) because order drives the
// function signature generated around the step's user-authored Code.
type ScriptParam struct {
	Name     string
	Property ToolProperty
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
	// TokenSuffix is the last 4 hex chars of the raw webhook token, kept in plaintext (unlike
	// SecretHash) purely so the UI can render a persistent "••••1234" reminder; never enough to
	// reconstruct the token. Empty for manual/provider-linked triggers.
	TokenSuffix string
	// Matchers gates which triggers (among possibly several sharing one provider link, see
	// trigger_provider_links) actually fire for a given inbound delivery. Seeded from the chosen
	// TriggerPreset's DefaultMatchers at creation time; empty always matches (standalone
	// triggers, which are the only recipient of their own webhook URL anyway).
	Matchers  TriggerMatchers
	Enabled   bool
	CreatedAt time.Time
}

// TriggerMatchers decides whether an inbound delivery should fire a given trigger. Only
// CheckHeaders exists today (AND semantics across entries, same philosophy as
// EvaluateTriggerFilters); future check kinds (body field match, mail subject/from match) extend
// this struct the same way TractCondition covers filter ops.
type TriggerMatchers struct {
	CheckHeaders []HeaderMatcher
}

type HeaderMatcher struct {
	Header string
	Equals string
}

// TriggerTractLink links triggerUuid to tractUuid with filter conditions (AND semantics) that
// gate whether a given webhook delivery starts a run for that specific tract.
type TriggerTractLink struct {
	TriggerUuid uuid.UUID
	TractUuid   uuid.UUID
	Filters     []TractCondition
}

// TriggerPreset describes one webhook trigger source (e.g. "gitlab_push") — backed by the
// trigger_presets table (seeded via migration, see 038_trigger_presets_and_links.sql), not
// hardcoded Go. Provider is the external_connections.provider a trigger created from this preset
// attaches to (see trigger_provider_links); empty means standalone (own trigger_uuid/secret_hash
// webhook URL, user-declared PayloadSchema).
type TriggerPreset struct {
	Key             string
	Category        string
	Label           string
	Description     string
	Provider        string
	PayloadSchema   ToolSchema
	DefaultMatchers TriggerMatchers
}
