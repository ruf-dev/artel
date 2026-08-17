package subscription

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
)

// unlimitedFeatures grants every known feature — the free implementation reproduces the
// pre-subscriptions-layer behavior of gating nothing.
var unlimitedFeatures = domain.FeatureSet{
	Emails:       true,
	TaskTrackers: true,
	Notes:        true,
	Spreadsheets: true,
	Connectors:   true,
	Tract:        true,
}

// FreeService is used when config.EnvironmentConfig.SubscriptionsEnabled is false: every check
// always passes and quotas are unlimited, with no CouchDB/S3 calls made (GetUsage would be
// meaningless when nothing enforces it).
type FreeService struct{}

func NewFree() *FreeService {
	return &FreeService{}
}

func (s *FreeService) CheckActive(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (s *FreeService) HasFeature(_ context.Context, _ uuid.UUID, _ domain.SubscriptionFeature) (bool, error) {
	return true, nil
}

func (s *FreeService) CheckFeature(_ context.Context, _ uuid.UUID, _ domain.SubscriptionFeature) error {
	return nil
}

func (s *FreeService) GetEffective(_ context.Context, userUuid uuid.UUID) (domain.EffectiveSubscription, error) {
	effective := domain.EffectiveSubscription{
		UserUuid:         userUuid,
		Active:           true,
		Features:         unlimitedFeatures,
		CouchQuotaBytes:  math.MaxInt64,
		S3QuotaBytes:     math.MaxInt64,
		MaxHotPlugSkills: math.MaxInt32,
		MaxTotalSkills:   math.MaxInt32,
	}

	return effective, nil
}

func (s *FreeService) GetUsage(_ context.Context, _ uuid.UUID) (domain.StorageUsage, error) {
	return domain.StorageUsage{}, nil
}

func (s *FreeService) CheckStorageQuota(_ context.Context, _ uuid.UUID) error {
	return nil
}

// CheckSkillLimit always passes — skill caps are unlimited when subscriptions are disabled,
// matching every other check in this no-op implementation.
func (s *FreeService) CheckSkillLimit(_ context.Context, _ uuid.UUID, _ bool) error {
	return nil
}
