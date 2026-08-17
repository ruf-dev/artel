package subscription

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	s3client "github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// PaidService is used when config.EnvironmentConfig.SubscriptionsEnabled is true: it enforces
// each user's effective plan+overrides — feature gating (formerly UserPermissions) and storage
// quotas, measured on demand against CouchDB/S3 rather than a maintained running counter.
type PaidService struct {
	subscriptions  repository.Subscriptions
	vaults         repository.Vaults
	couchInstances repository.CouchInstances
	s3Instances    repository.S3Instances
}

func NewPaid(
	subscriptions repository.Subscriptions,
	vaults repository.Vaults,
	couchInstances repository.CouchInstances,
	s3Instances repository.S3Instances,
) *PaidService {
	service := &PaidService{
		subscriptions:  subscriptions,
		vaults:         vaults,
		couchInstances: couchInstances,
		s3Instances:    s3Instances,
	}

	return service
}

func (s *PaidService) CheckActive(ctx context.Context, userUuid uuid.UUID) error {
	sub, err := s.subscriptions.GetByUser(ctx, userUuid)
	if err != nil {
		return rerrors.Wrap(err, "get subscription")
	}

	if !sub.Active {
		return user_errors.NoActiveSubscription
	}

	return nil
}

func (s *PaidService) GetEffective(ctx context.Context, userUuid uuid.UUID) (domain.EffectiveSubscription, error) {
	effective, err := s.subscriptions.GetWithPlan(ctx, userUuid)
	if err != nil {
		return domain.EffectiveSubscription{}, rerrors.Wrap(err, "error getting effective subscription")
	}

	return effective, nil
}

func (s *PaidService) HasFeature(
	ctx context.Context, userUuid uuid.UUID, feature domain.SubscriptionFeature,
) (bool, error) {
	effective, err := s.GetEffective(ctx, userUuid)
	if err != nil {
		return false, err
	}

	return effective.Features.Has(feature), nil
}

func (s *PaidService) CheckFeature(ctx context.Context, userUuid uuid.UUID, feature domain.SubscriptionFeature) error {
	hasFeature, err := s.HasFeature(ctx, userUuid, feature)
	if err != nil {
		return err
	}

	if !hasFeature {
		return rerrors.Wrap(user_errors.FeatureNotEnabled, string(feature))
	}

	return nil
}

// GetUsage sums CouchDB + S3 storage across every vault userUuid belongs to, measured live —
// no maintained running counter.
func (s *PaidService) GetUsage(ctx context.Context, userUuid uuid.UUID) (domain.StorageUsage, error) {
	memberVaults, err := s.vaults.ListByMembership(ctx, userUuid)
	if err != nil {
		return domain.StorageUsage{}, rerrors.Wrap(err, "error listing vaults for usage")
	}

	var usage domain.StorageUsage

	for _, vault := range memberVaults {
		couchBytes, couchErr := s.couchUsage(ctx, vault)
		if couchErr != nil {
			return domain.StorageUsage{}, couchErr
		}

		usage.CouchBytes += couchBytes

		if vault.S3InstanceUuid == nil {
			continue
		}

		s3Bytes, s3Err := s.s3Usage(ctx, vault)
		if s3Err != nil {
			return domain.StorageUsage{}, s3Err
		}

		usage.S3Bytes += s3Bytes
	}

	return usage, nil
}

func (s *PaidService) couchUsage(ctx context.Context, vault domain.Vault) (int64, error) {
	instance, err := s.couchInstances.Get(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return 0, rerrors.Wrap(err, "error getting couch instance for usage")
	}

	cfg := couchdb.Config{
		BaseURL:  instance.Url,
		User:     instance.Username,
		Password: instance.Password,
	}

	client, err := couchdb.New(cfg)
	if err != nil {
		return 0, rerrors.Wrap(err, "error initializing couch client for usage")
	}

	size, err := client.DatabaseSize(ctx, vault.CouchDBName)
	if err != nil {
		return 0, rerrors.Wrap(err, "error getting couch database size")
	}

	return size, nil
}

func (s *PaidService) s3Usage(ctx context.Context, vault domain.Vault) (int64, error) {
	instance, err := s.s3Instances.Get(ctx, *vault.S3InstanceUuid)
	if err != nil {
		return 0, rerrors.Wrap(err, "error getting s3 instance for usage")
	}

	cfg := s3client.Config{
		Endpoint:  instance.Endpoint,
		Region:    instance.Region,
		AccessKey: instance.AccessKey,
		SecretKey: instance.SecretKey,
		UseSSL:    instance.UseSSL,
		PathStyle: instance.PathStyle,
	}

	client, err := s3client.New(cfg, vault.S3BucketName)
	if err != nil {
		return 0, rerrors.Wrap(err, "error initializing s3 client for usage")
	}

	size, err := client.TotalSize(ctx)
	if err != nil {
		return 0, rerrors.Wrap(err, "error getting s3 bucket size")
	}

	return size, nil
}

// CheckStorageQuota rejects once usage is already at or over quota — a coarse pre-write check
// (on-demand measurement, not byte-exact accounting for the write in flight).
func (s *PaidService) CheckStorageQuota(ctx context.Context, userUuid uuid.UUID) error {
	effective, err := s.GetEffective(ctx, userUuid)
	if err != nil {
		return err
	}

	usage, err := s.GetUsage(ctx, userUuid)
	if err != nil {
		return err
	}

	if usage.CouchBytes >= effective.CouchQuotaBytes {
		return user_errors.CouchStorageQuotaExceeded
	}

	if usage.S3Bytes >= effective.S3QuotaBytes {
		return user_errors.S3StorageQuotaExceeded
	}

	return nil
}

// skillCounts measures a single vault's current skill population live against CouchDB: total
// is every note under either .skills/ subfolder, hotPlug is the subset under
// domain.SkillsActiveFolder. Uses instance-level admin credentials to build the LiveSyncClient
// (this check isn't scoped to a particular caller's couch account, unlike notes.Service), the
// same credential source couchUsage already reads from instance.
func (s *PaidService) skillCounts(ctx context.Context, vault domain.Vault) (total int, hotPlug int, err error) {
	instance, err := s.couchInstances.Get(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return 0, 0, rerrors.Wrap(err, "error getting couch instance for skill limit check")
	}

	client := couchdb.NewLiveSyncClient(instance.Url, vault.CouchDBName, instance.Username, instance.Password)

	notes, err := client.ListSkillNotes(ctx)
	if err != nil {
		return 0, 0, rerrors.Wrap(err, "error listing skill notes")
	}

	for _, n := range notes {
		switch {
		case strings.HasPrefix(n.Path, domain.SkillsActiveFolder):
			total++
			hotPlug++
		case strings.HasPrefix(n.Path, domain.SkillsLibraryFolder):
			total++
		}
	}

	return total, hotPlug, nil
}

// CheckSkillLimit compares vaultUuid's current skill counts against its owner's effective plan
// caps. wantHotPlug scopes which caps are checked: a library-only create/update only needs the
// total cap checked, while creating or hot-plugging a skill needs both.
func (s *PaidService) CheckSkillLimit(ctx context.Context, vaultUuid uuid.UUID, wantHotPlug bool) error {
	vault, err := s.vaults.GetByID(ctx, vaultUuid)
	if err != nil {
		return rerrors.Wrap(err, "error getting vault for skill limit check")
	}

	effective, err := s.GetEffective(ctx, vault.UserUuid)
	if err != nil {
		return err
	}

	total, hotPlug, err := s.skillCounts(ctx, vault)
	if err != nil {
		return err
	}

	if total >= effective.MaxTotalSkills {
		return user_errors.SkillLimitExceeded
	}

	if wantHotPlug && hotPlug >= effective.MaxHotPlugSkills {
		return user_errors.HotPlugSkillLimitExceeded
	}

	return nil
}
