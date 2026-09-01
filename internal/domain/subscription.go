package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionFeature string

const (
	FeatureEmails       SubscriptionFeature = "emails"
	FeatureTaskTrackers SubscriptionFeature = "task_trackers"
	FeatureNotes        SubscriptionFeature = "notes"
	FeatureSpreadsheets SubscriptionFeature = "spreadsheets"
	FeatureConnectors   SubscriptionFeature = "connectors"
	FeatureTract        SubscriptionFeature = "tract"
	FeaturePostgres     SubscriptionFeature = "postgres"
	FeatureNotify       SubscriptionFeature = "notify"
)

type FeatureSet struct {
	Emails       bool `json:"emails"`
	TaskTrackers bool `json:"task_trackers"`
	Notes        bool `json:"notes"`
	Spreadsheets bool `json:"spreadsheets"`
	Connectors   bool `json:"connectors"`
	Tract        bool `json:"tract"`
	Postgres     bool `json:"postgres"`
	Notify       bool `json:"notify"`
}

func (f FeatureSet) Has(feature SubscriptionFeature) bool {
	switch feature {
	case FeatureEmails:
		return f.Emails
	case FeatureTaskTrackers:
		return f.TaskTrackers
	case FeatureNotes:
		return f.Notes
	case FeatureSpreadsheets:
		return f.Spreadsheets
	case FeatureConnectors:
		return f.Connectors
	case FeatureTract:
		return f.Tract
	case FeaturePostgres:
		return f.Postgres
	case FeatureNotify:
		return f.Notify
	default:
		return false
	}
}

// With returns a copy of f with feature set to value — used to apply a sparse override map on
// top of a plan's FeatureSet.
func (f FeatureSet) With(feature SubscriptionFeature, value bool) FeatureSet {
	switch feature {
	case FeatureEmails:
		f.Emails = value
	case FeatureTaskTrackers:
		f.TaskTrackers = value
	case FeatureNotes:
		f.Notes = value
	case FeatureSpreadsheets:
		f.Spreadsheets = value
	case FeatureConnectors:
		f.Connectors = value
	case FeatureTract:
		f.Tract = value
	case FeaturePostgres:
		f.Postgres = value
	case FeatureNotify:
		f.Notify = value
	}

	return f
}

type SubscriptionPlan struct {
	PlanKey          string
	CouchQuotaBytes  int64
	S3QuotaBytes     int64
	MaxHotPlugSkills int
	MaxTotalSkills   int
	Features         FeatureSet
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Subscription is the per-user row: which plan the user is on, plus admin-granted overrides on
// top of that plan — sparse (only keys present in FeatureOverrides override the plan's
// FeatureSet) so most users have an empty override map and simply inherit their plan.
type Subscription struct {
	UserUuid                 uuid.UUID
	Active                   bool
	PlanKey                  string
	FeatureOverrides         map[SubscriptionFeature]bool
	CouchQuotaOverrideBytes  *int64
	S3QuotaOverrideBytes     *int64
	MaxHotPlugSkillsOverride *int
	MaxTotalSkillsOverride   *int
}

// EffectiveSubscription is the merged plan+override view services actually consult — the
// plan's defaults with the user's overrides applied on top.
type EffectiveSubscription struct {
	UserUuid         uuid.UUID
	Active           bool
	PlanKey          string
	Features         FeatureSet
	CouchQuotaBytes  int64
	S3QuotaBytes     int64
	MaxHotPlugSkills int
	MaxTotalSkills   int
}

// StorageUsage is a point-in-time measurement of a user's storage footprint across all of
// their vaults: CouchDB (markdown notes, and binaries when UseCouchDBForBinaries) plus S3
// (binaries, for vaults with a bucket linked).
type StorageUsage struct {
	CouchBytes int64
	S3Bytes    int64
}

// MaxSingleFileBytes caps the size of any single write (note content, imported entry, or MCP
// file write), independent of plan quota. CheckStorageQuota only rejects once measured usage is
// already at/over quota, so a user's very first write starts at zero usage and would otherwise
// pass unchecked no matter how large — this bounds that gap.
const MaxSingleFileBytes = 10 * 1024 * 1024
