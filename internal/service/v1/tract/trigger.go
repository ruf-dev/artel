package tract

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

const triggerTokenBytes = 32

// CreateTrigger creates a standalone trigger and returns the raw webhook token once — only its
// sha256 (SecretHash) is persisted. If source names a provider-linked preset (e.g. gitlab_push),
// it instead attaches to the caller's existing connection for that provider (see
// createProviderLinkedTrigger) and returns an empty token — there's no per-trigger secret, the
// trigger shares the provider connection's own webhook URL/secret.
func (s *Service) CreateTrigger(
	ctx context.Context,
	name string,
	kind string,
	source string,
	config json.RawMessage,
	payloadSchema domain.ToolSchema,
) (domain.Trigger, string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Trigger{}, "", rerrors.Wrap(user_errors.Unauthenticated)
	}

	preset, err := s.triggerPresets.GetByKey(ctx, source)
	if err != nil {
		return domain.Trigger{}, "", rerrors.Wrap(err, "error getting trigger preset")
	}

	if preset.Valid && preset.V.Provider != "" {
		return s.createProviderLinkedTrigger(ctx, uc.UserUuid, name, kind, preset.V)
	}

	rawToken, err := generateTriggerToken()
	if err != nil {
		return domain.Trigger{}, "", rerrors.Wrap(err, "error generating trigger token")
	}

	hash := sha256.Sum256([]byte(rawToken))

	trigger := domain.Trigger{
		UserUuid:      uc.UserUuid,
		Name:          name,
		Kind:          kind,
		Source:        source,
		Config:        config,
		PayloadSchema: payloadSchema,
		SecretHash:    hash[:],
		Enabled:       true,
	}

	created, err := s.triggers.Create(ctx, trigger)
	if err != nil {
		return domain.Trigger{}, "", rerrors.Wrap(err, "error creating trigger")
	}

	return created, rawToken, nil
}

// createProviderLinkedTrigger resolves the caller's existing connection for preset.Provider and
// links the new trigger to it (trigger_provider_links) instead of minting a private
// trigger_uuid/secret_hash webhook URL — the delivery arrives at the provider's shared inbound
// endpoint (e.g. /webhooks/gitlab/{external_connection_id}) and gitlab_webhook.Handler fans it
// out to every trigger linked there, gated by TriggerMatchers.
func (s *Service) createProviderLinkedTrigger(
	ctx context.Context,
	userUuid uuid.UUID,
	name string,
	kind string,
	preset domain.TriggerPreset,
) (domain.Trigger, string, error) {
	conn, err := s.externalConns.GetByUserAndProvider(ctx, userUuid, preset.Provider)
	if err != nil {
		return domain.Trigger{}, "", rerrors.Wrap(err, "error getting provider connection")
	}

	if !conn.Valid {
		return domain.Trigger{}, "", rerrors.Wrap(user_errors.TriggerProviderConnectionRequired, preset.Provider)
	}

	trigger := domain.Trigger{
		UserUuid:      userUuid,
		Name:          name,
		Kind:          kind,
		Source:        preset.Key,
		Config:        json.RawMessage(`{}`),
		PayloadSchema: preset.PayloadSchema,
		Matchers:      preset.DefaultMatchers,
		Enabled:       true,
	}

	created, err := s.triggers.Create(ctx, trigger)
	if err != nil {
		return domain.Trigger{}, "", rerrors.Wrap(err, "error creating provider-linked trigger")
	}

	err = s.triggers.LinkToProvider(ctx, created.Uuid, conn.V.Uuid)
	if err != nil {
		return domain.Trigger{}, "", rerrors.Wrap(err, "error linking trigger to provider connection")
	}

	return created, "", nil
}

func generateTriggerToken() (string, error) {
	buf := make([]byte, triggerTokenBytes)

	_, err := rand.Read(buf)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading random bytes")
	}

	return hex.EncodeToString(buf), nil
}

func (s *Service) GetTrigger(ctx context.Context, id uuid.UUID) (domain.Trigger, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Trigger{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	return s.getOwnedTrigger(ctx, id, uc.UserUuid)
}

func (s *Service) getOwnedTrigger(ctx context.Context, id uuid.UUID, ownerUuid uuid.UUID) (domain.Trigger, error) {
	result, err := s.triggers.Get(ctx, id)
	if err != nil {
		return domain.Trigger{}, rerrors.Wrap(err, "error getting trigger")
	}

	if !result.Valid {
		return domain.Trigger{}, rerrors.Wrap(user_errors.TriggerNotFound)
	}

	if result.V.UserUuid != ownerUuid {
		return domain.Trigger{}, rerrors.Wrap(user_errors.TriggerNotOwned)
	}

	return result.V, nil
}

func (s *Service) ListTriggers(ctx context.Context) ([]domain.Trigger, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	triggers, err := s.triggers.ListByUser(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing triggers")
	}

	return triggers, nil
}

func (s *Service) SetTriggerEnabled(ctx context.Context, id uuid.UUID, enabled bool) error {
	_, err := s.GetTrigger(ctx, id)
	if err != nil {
		return err
	}

	err = s.triggers.SetEnabled(ctx, id, enabled)
	if err != nil {
		return rerrors.Wrap(err, "error setting trigger enabled")
	}

	return nil
}

func (s *Service) DeleteTrigger(ctx context.Context, id uuid.UUID) error {
	_, err := s.GetTrigger(ctx, id)
	if err != nil {
		return err
	}

	err = s.triggers.Delete(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "error deleting trigger")
	}

	return nil
}

// RotateTriggerToken invalidates id's current webhook URL and mints a new one, returning the raw
// token once — only its sha256 (SecretHash) is persisted. Same one-time-reveal contract as
// CreateTrigger.
func (s *Service) RotateTriggerToken(ctx context.Context, id uuid.UUID) (domain.Trigger, string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Trigger{}, "", rerrors.Wrap(user_errors.Unauthenticated)
	}

	_, err := s.getOwnedTrigger(ctx, id, uc.UserUuid)
	if err != nil {
		return domain.Trigger{}, "", err
	}

	rawToken, err := generateTriggerToken()
	if err != nil {
		return domain.Trigger{}, "", rerrors.Wrap(err, "error generating trigger token")
	}

	hash := sha256.Sum256([]byte(rawToken))
	newTriggerUuid := uuid.New()

	rotated, err := s.triggers.RotateSecret(ctx, id, newTriggerUuid, hash[:])
	if err != nil {
		return domain.Trigger{}, "", rerrors.Wrap(err, "error rotating trigger token")
	}

	return rotated, rawToken, nil
}

// LinkTrigger fans triggerUuid out to tractUuid with the given filters — both must belong to
// the caller.
func (s *Service) LinkTrigger(
	ctx context.Context,
	triggerUuid uuid.UUID,
	tractUuid uuid.UUID,
	filters []domain.TractCondition,
) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	_, err := s.getOwnedTrigger(ctx, triggerUuid, uc.UserUuid)
	if err != nil {
		return err
	}

	tract, err := s.tracts.Get(ctx, tractUuid)
	if err != nil {
		return rerrors.Wrap(err, "error getting tract")
	}

	if !tract.Valid {
		return rerrors.Wrap(user_errors.NotFound)
	}

	if tract.V.UserUuid != uc.UserUuid {
		return rerrors.Wrap(user_errors.TractNotOwned)
	}

	link := domain.TriggerTractLink{
		TriggerUuid: triggerUuid,
		TractUuid:   tractUuid,
		Filters:     filters,
	}

	err = s.triggers.Link(ctx, link)
	if err != nil {
		return rerrors.Wrap(err, "error linking trigger to tract")
	}

	return nil
}

func (s *Service) UnlinkTrigger(ctx context.Context, triggerUuid uuid.UUID, tractUuid uuid.UUID) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	_, err := s.getOwnedTrigger(ctx, triggerUuid, uc.UserUuid)
	if err != nil {
		return err
	}

	err = s.triggers.Unlink(ctx, triggerUuid, tractUuid)
	if err != nil {
		return rerrors.Wrap(err, "error unlinking trigger from tract")
	}

	return nil
}

// ListLinksByTract returns the triggers linked to tractUuid — the tract editor's "wired up
// triggers" view.
func (s *Service) ListLinksByTract(ctx context.Context, tractUuid uuid.UUID) ([]repository.TractTriggerLink, error) {
	_, err := s.GetTract(ctx, tractUuid)
	if err != nil {
		return nil, err
	}

	links, err := s.triggers.ListLinksByTract(ctx, tractUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tract trigger links")
	}

	return links, nil
}
