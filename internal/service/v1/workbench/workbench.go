// Package workbench implements service.WorkbenchService: it manages the per-vault Docker
// container backing a vault's cloud workbench. See
// docs/workbench/01_data_model_and_lifecycle.md for the full state machine and
// docs/workbench/02_docker_topology.md for the Docker client this composes.
package workbench

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/containerd/errdefs"
	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/workbenchdocker"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// anthropicApiKeyEnvVar is the env var name StartContainer injects the BYOK key under —
// docs/workbench/03_auth_and_login_flow.md, "api_key mode".
const anthropicApiKeyEnvVar = "ANTHROPIC_API_KEY"

// dockerClient is the narrow subset of workbenchdocker.Client's methods this service depends
// on. workbenchdocker.Client is a concrete struct — this interface exists so tests can exercise
// Service's branching logic against a fake, without a live Docker daemon.
type dockerClient interface {
	CreateVolume(ctx context.Context, name string) error
	RemoveVolume(ctx context.Context, name string) error
	CreateContainer(ctx context.Context, opts workbenchdocker.CreateOpts) (string, error)
	StartContainer(ctx context.Context, containerID string, env map[string]string) error
	StopContainer(ctx context.Context, containerID string) error
	RemoveContainer(ctx context.Context, containerID string) error
	CapturePane(ctx context.Context, containerID string) (string, error)
	SendKeys(ctx context.Context, containerID string, keys string) error
}

// externalConnectionService is the narrow subset of service.ExternalConnectionService this
// service depends on, to resolve the BYOK Anthropic key for api_key auth mode without a
// compile-time dependency on the whole external-connections service.
type externalConnectionService interface {
	GetAnthropicApiKey(ctx context.Context, userUuid uuid.UUID) (string, error)
}

type Service struct {
	workbenchesRepo repository.Workbenches
	vaultsRepo      repository.Vaults

	docker              dockerClient
	externalConnections externalConnectionService
}

func New(
	workbenches repository.Workbenches,
	vaults repository.Vaults,
	docker dockerClient,
	externalConnections externalConnectionService,
) *Service {
	return &Service{
		workbenchesRepo: workbenches,
		vaultsRepo:      vaults,

		docker:              docker,
		externalConnections: externalConnections,
	}
}

// CreateWorkbench records a 'configuring' workbenches row for vaultID, then provisions the
// Docker volume and (unstarted) container backing it, and finally records the container id
// with status='created'. Safely re-invocable: if a workbenches row for vaultID already exists
// (e.g. a retry after a partial failure), it's returned as-is rather than duplicated — see
// docs/workbench/01_data_model_and_lifecycle.md, "Hook points".
//
// The row is inserted before the Docker calls (rather than after, as previously) so that a
// vaultID that never gets a workbench created for it is still driven by a caller explicitly
// invoking this method (see CreateWorkbench RPC), and so a caller/UI can observe a
// 'configuring' state while Docker provisioning is in flight. If CreateVolume or
// CreateContainer fails after the row is inserted, the row is intentionally left in
// 'configuring' rather than rolled back or cleaned up here — DeleteWorkbench already tolerates
// partial Docker state (missing container/volume) when tearing down, so a stuck 'configuring'
// row is recoverable by deleting and retrying rather than needing bespoke rollback logic here.
func (s *Service) CreateWorkbench(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error) {
	existing, err := s.workbenchesRepo.GetByVaultID(ctx, vaultID)
	if err == nil {
		return existing, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault id")
	}

	vault, err := s.vaultsRepo.GetByID(ctx, vaultID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting vault by id")
	}

	volumeName := fmt.Sprintf("workbench-%s", vaultID.String())

	_, err = s.workbenchesRepo.Create(ctx, vaultID, vault.UserUuid, volumeName)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error creating workbench")
	}

	err = s.docker.CreateVolume(ctx, volumeName)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error creating workbench volume")
	}

	createOpts := workbenchdocker.CreateOpts{
		Name:       volumeName,
		VolumeName: volumeName,
	}

	containerID, err := s.docker.CreateContainer(ctx, createOpts)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error creating workbench container")
	}

	err = s.workbenchesRepo.MarkContainerCreated(ctx, vaultID, containerID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error marking workbench container created")
	}

	workbench, err := s.workbenchesRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault id")
	}

	return workbench, nil
}

// GetWorkbench returns the workbench for vaultID.
func (s *Service) GetWorkbench(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error) {
	workbench, err := s.workbenchesRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench")
	}

	return workbench, nil
}

// StartWorkbench starts vaultID's workbench container under authMode — see
// docs/workbench/03_auth_and_login_flow.md for both modes' designs. Any mode besides the two
// domain.WorkbenchAuthMode constants fails with user_errors.WorkbenchAuthModeNotImplemented
// rather than silently falling back.
func (s *Service) StartWorkbench(
	ctx context.Context, vaultID uuid.UUID, authMode domain.WorkbenchAuthMode,
) (domain.Workbench, error) {
	switch authMode {
	case domain.WorkbenchAuthModeAPIKey:
		return s.startWithApiKey(ctx, vaultID)
	case domain.WorkbenchAuthModeSubscriptionLogin:
		return s.startWithSubscriptionLogin(ctx, vaultID)
	default:
		return domain.Workbench{}, user_errors.WorkbenchAuthModeNotImplemented
	}
}

// startWithApiKey resolves the workbench owner's BYOK Anthropic api key (failing fast with
// user_errors.WorkbenchMissingAnthropicConnection if none is connected), starts the container
// with it injected as ANTHROPIC_API_KEY, then records status='running'/started_at.
func (s *Service) startWithApiKey(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error) {
	wb, err := s.workbenchesRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault id")
	}

	apiKey, err := s.externalConnections.GetAnthropicApiKey(ctx, wb.UserUuid)
	if err != nil {
		if errors.Is(err, user_errors.LlmKeyRequired) {
			return domain.Workbench{}, user_errors.WorkbenchMissingAnthropicConnection
		}

		return domain.Workbench{}, rerrors.Wrap(err, "error resolving anthropic api key for workbench owner")
	}

	env := map[string]string{
		anthropicApiKeyEnvVar: apiKey,
	}

	err = s.workbenchesRepo.MarkConfiguring(ctx, vaultID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error marking workbench configuring")
	}

	err = s.docker.StartContainer(ctx, wb.ContainerId, env)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error starting workbench container")
	}

	return s.markRunningAndReload(ctx, vaultID, domain.WorkbenchAuthModeAPIKey)
}

// startWithSubscriptionLogin starts the container with no auth env vars injected at all — see
// docs/workbench/03_auth_and_login_flow.md, "Mechanism (confirmed)", step 1. The workbench is
// recorded as status='running' as soon as the container starts, independent of whether the
// in-container login TUI flow has completed — see
// docs/workbench/01_data_model_and_lifecycle.md's state table: "running" means "container
// started", not "authenticated". Callers observe/drive the actual login flow separately via
// GetLoginPrompt/SubmitLoginCode.
func (s *Service) startWithSubscriptionLogin(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error) {
	wb, err := s.workbenchesRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault id")
	}

	err = s.workbenchesRepo.MarkConfiguring(ctx, vaultID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error marking workbench configuring")
	}

	err = s.docker.StartContainer(ctx, wb.ContainerId, nil)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error starting workbench container")
	}

	return s.markRunningAndReload(ctx, vaultID, domain.WorkbenchAuthModeSubscriptionLogin)
}

// markRunningAndReload records status='running'/started_at under authMode and returns the
// refreshed row, shared by both StartWorkbench auth-mode branches.
func (s *Service) markRunningAndReload(
	ctx context.Context, vaultID uuid.UUID, authMode domain.WorkbenchAuthMode,
) (domain.Workbench, error) {
	err := s.workbenchesRepo.MarkRunning(ctx, vaultID, authMode)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error marking workbench running")
	}

	wb, err := s.workbenchesRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault id")
	}

	return wb, nil
}

// GetLoginPrompt returns a snapshot of vaultID's subscription_login TUI flow, derived by
// capturing the workbench's tmux pane and parsing it — see
// docs/workbench/03_auth_and_login_flow.md, "Mechanism (confirmed)", steps 2-3 and 5-6.
func (s *Service) GetLoginPrompt(ctx context.Context, vaultID uuid.UUID) (domain.WorkbenchLoginPrompt, error) {
	wb, err := s.workbenchesRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return domain.WorkbenchLoginPrompt{}, rerrors.Wrap(err, "error getting workbench by vault id")
	}

	pane, err := s.docker.CapturePane(ctx, wb.ContainerId)
	if err != nil {
		return domain.WorkbenchLoginPrompt{}, rerrors.Wrap(err, "error capturing workbench tmux pane")
	}

	prompt := parseLoginPrompt(pane)

	return prompt, nil
}

// SubmitLoginCode relays the user's pasted OAuth code (or, on a brand-new container, their
// first theme/menu keystroke) back into vaultID's workbench tmux session — see
// docs/workbench/03_auth_and_login_flow.md, "Mechanism (confirmed)", step 4.
func (s *Service) SubmitLoginCode(ctx context.Context, vaultID uuid.UUID, code string) error {
	wb, err := s.workbenchesRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error getting workbench by vault id")
	}

	err = s.docker.SendKeys(ctx, wb.ContainerId, code)
	if err != nil {
		return rerrors.Wrap(err, "error relaying login code into workbench session")
	}

	return nil
}

// StopWorkbench stops vaultID's running workbench container, retaining it (and its volume) for
// a later restart, and records status='stopped'/stopped_at.
func (s *Service) StopWorkbench(ctx context.Context, vaultID uuid.UUID) error {
	wb, err := s.workbenchesRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error getting workbench by vault id")
	}

	err = s.docker.StopContainer(ctx, wb.ContainerId)
	if err != nil {
		return rerrors.Wrap(err, "error stopping workbench container")
	}

	err = s.workbenchesRepo.MarkStopped(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error marking workbench stopped")
	}

	return nil
}

// DeleteWorkbench tears down vaultID's workbench: stops (if a container was created) and
// removes the container, removes the volume, then deletes the DB row — in that order, so a
// failure partway through never leaves the DB pointing at nothing while Docker still holds a
// live container (see docs/workbench/01_data_model_and_lifecycle.md, "Vault deletion"). A
// missing workbench row is treated as already-deleted, not an error. "Not found" from Docker
// (container/volume already gone, e.g. a retry picking up after a prior partial failure) is
// likewise tolerated at each step rather than failing the whole call.
func (s *Service) DeleteWorkbench(ctx context.Context, vaultID uuid.UUID) error {
	wb, err := s.workbenchesRepo.GetByVaultID(ctx, vaultID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return rerrors.Wrap(err, "error getting workbench by vault id")
	}

	if wb.ContainerId != "" {
		err = s.docker.StopContainer(ctx, wb.ContainerId)
		if err != nil && !errdefs.IsNotFound(err) {
			return rerrors.Wrap(err, "error stopping workbench container")
		}

		err = s.docker.RemoveContainer(ctx, wb.ContainerId)
		if err != nil && !errdefs.IsNotFound(err) {
			return rerrors.Wrap(err, "error removing workbench container")
		}
	}

	err = s.docker.RemoveVolume(ctx, wb.VolumeName)
	if err != nil && !errdefs.IsNotFound(err) {
		return rerrors.Wrap(err, "error removing workbench volume")
	}

	err = s.workbenchesRepo.Delete(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error deleting workbench")
	}

	return nil
}
