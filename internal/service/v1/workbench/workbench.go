// Package workbench implements service.WorkbenchService: it manages the per-(vault, user)
// Docker container backing a vault member's cloud workbench — every vault member gets their own
// workbench. See docs/workbench/01_data_model_and_lifecycle.md for the full state machine and
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
	"github.com/ruf-dev/artel/internal/middleware/user_context"
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
	ContainerAddress(ctx context.Context, containerID string) (string, error)
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
	vaultMembers    repository.VaultMembers
	dockerHostsRepo repository.DockerHosts

	// newDockerClient builds a fresh dockerClient pointed at a given docker host's Url (plus its
	// decrypted TLS material, empty for a local/unauthenticated daemon). A fresh client is built
	// per-call rather than cached, since a workbench pool now has more than one host — mirrors
	// how internal/service/v1/couchinstances/couchinstances.go builds a fresh couchdb.New(cfg)
	// per call instead of caching a client.
	newDockerClient func(host string, tlsCfg workbenchdocker.TLSConfig) (dockerClient, error)

	externalConnections externalConnectionService
}

func New(
	workbenches repository.Workbenches,
	vaults repository.Vaults,
	vaultMembers repository.VaultMembers,
	dockerHosts repository.DockerHosts,
	externalConnections externalConnectionService,
	newDockerClient func(host string, tlsCfg workbenchdocker.TLSConfig) (dockerClient, error),
) *Service {
	if newDockerClient == nil {
		newDockerClient = newRealDockerClient
	}

	return &Service{
		workbenchesRepo: workbenches,
		vaultsRepo:      vaults,
		vaultMembers:    vaultMembers,
		dockerHostsRepo: dockerHosts,

		newDockerClient:     newDockerClient,
		externalConnections: externalConnections,
	}
}

// newRealDockerClient is the production default for Service.newDockerClient — wraps
// workbenchdocker.New so tests can inject a fake instead without a live Docker daemon.
func newRealDockerClient(host string, tlsCfg workbenchdocker.TLSConfig) (dockerClient, error) {
	client, err := workbenchdocker.New(host, tlsCfg)
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating workbench docker client")
	}

	return client, nil
}

// resolveClient resolves the dockerClient for wb's assigned docker host. Returns
// user_errors.WorkbenchMissingDockerHost for a workbenches row that predates the docker_hosts
// pool (nullable docker_host_id, no backfill — see migrations/061_docker_hosts.sql).
//
// It fetches the host via GetWithCreds (rather than Get) so a remote host registered with TLS
// material gets it threaded into newDockerClient — a local/unauthenticated host's CaCert/
// ClientCert/ClientKey just come back empty, which newDockerClient/workbenchdocker.New treat as
// "no TLS", so this is a no-op for the existing unix-socket/dind path.
func (s *Service) resolveClient(ctx context.Context, wb domain.Workbench) (dockerClient, error) {
	if wb.DockerHostUuid == nil {
		return nil, user_errors.WorkbenchMissingDockerHost
	}

	host, err := s.dockerHostsRepo.GetWithCreds(ctx, *wb.DockerHostUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting docker host")
	}

	tlsCfg := workbenchdocker.TLSConfig{
		CaCert:     host.CaCert,
		ClientCert: host.ClientCert,
		ClientKey:  host.ClientKey,
	}

	client, err := s.newDockerClient(host.Url, tlsCfg)
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating docker client")
	}

	return client, nil
}

// CreateWorkbench records a 'configuring' workbenches row for the caller's (vaultID, userID)
// pair, then provisions the Docker volume and (unstarted) container backing it, and finally
// records the container id with status='created'. Safely re-invocable: if a workbenches row for
// (vaultID, userID) already exists (e.g. a retry after a partial failure), it's returned as-is
// rather than duplicated — see docs/workbench/01_data_model_and_lifecycle.md, "Hook points".
// Every vault member gets their own workbench row, so two different members of the same vault
// never share one.
//
// The caller must already be a member of vaultID — returns
// user_errors.WorkbenchRequiresVaultMembership otherwise.
//
// The row is inserted before the Docker calls (rather than after, as previously) so that a
// (vaultID, userID) pair that never gets a workbench created for it is still driven by a caller
// explicitly invoking this method (see CreateWorkbench RPC), and so a caller/UI can observe a
// 'configuring' state while Docker provisioning is in flight. If CreateVolume or
// CreateContainer fails after the row is inserted, the row is intentionally left in
// 'configuring' rather than rolled back or cleaned up here — DeleteWorkbench already tolerates
// partial Docker state (missing container/volume) when tearing down, so a stuck 'configuring'
// row is recoverable by deleting and retrying rather than needing bespoke rollback logic here.
//
// The docker host backing the new workbench is picked once, up front, by
// DockerHosts.PickLeastLoaded — spreading new workbenches across the registered pool. Returns
// user_errors.NoDockerHostsAvailable if the pool is empty.
func (s *Service) CreateWorkbench(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Workbench{}, rerrors.Wrap(user_errors.Unauthenticated)
	}
	userID := uc.UserUuid

	_, err := s.vaultMembers.Get(ctx, vaultID, userID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(user_errors.WorkbenchRequiresVaultMembership)
	}

	existing, err := s.workbenchesRepo.GetByVaultAndUser(ctx, vaultID, userID)
	if err == nil {
		return existing, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault and user")
	}

	host, err := s.dockerHostsRepo.PickLeastLoaded(ctx)
	if err != nil {
		// PickLeastLoaded's underlying sqlc query is a `:one` query against zero rows when the
		// docker_hosts pool is empty — the repo layer runs that through pg_err.UnwrapPgErr (same
		// as couchinstances.go's RandomPick), which turns sql.ErrNoRows into user_errors.NotFound
		// rather than surfacing sql.ErrNoRows itself, so that's what's checked for here.
		if errors.Is(err, user_errors.NotFound) {
			return domain.Workbench{}, user_errors.NoDockerHostsAvailable
		}

		return domain.Workbench{}, rerrors.Wrap(err, "error picking least loaded docker host")
	}

	// PickLeastLoaded's query doesn't select the TLS columns (it's a plain load-balancing pick,
	// same shape as Get) — fetch the decrypted creds separately so a remote TLS-secured host
	// picked here gets the same treatment as one resolved via resolveClient.
	hostWithCreds, err := s.dockerHostsRepo.GetWithCreds(ctx, host.Uuid)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting docker host creds")
	}

	tlsCfg := workbenchdocker.TLSConfig{
		CaCert:     hostWithCreds.CaCert,
		ClientCert: hostWithCreds.ClientCert,
		ClientKey:  hostWithCreds.ClientKey,
	}

	docker, err := s.newDockerClient(host.Url, tlsCfg)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error creating docker client")
	}

	// Two members of the same vault must never collide on Docker resource names, hence userID is
	// part of the name alongside vaultID.
	volumeName := fmt.Sprintf("workbench-%s-%s", vaultID, userID)

	_, err = s.workbenchesRepo.Create(ctx, vaultID, userID, volumeName, host.Uuid)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error creating workbench")
	}

	err = docker.CreateVolume(ctx, volumeName)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error creating workbench volume")
	}

	createOpts := workbenchdocker.CreateOpts{
		Name:       volumeName,
		VolumeName: volumeName,
	}

	containerID, err := docker.CreateContainer(ctx, createOpts)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error creating workbench container")
	}

	err = s.workbenchesRepo.MarkContainerCreated(ctx, vaultID, userID, containerID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error marking workbench container created")
	}

	workbench, err := s.workbenchesRepo.GetByVaultAndUser(ctx, vaultID, userID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault and user")
	}

	return workbench, nil
}

// ResolveTerminalTarget returns the base URL ("http://<container-ip>:<ttyd-port>") of userID's
// running workbench for vaultID's in-container ttyd server, as reachable from this process over
// the workbench docker network. Returns user_errors.WorkbenchNotRunning when the workbench isn't
// in the 'running' state — a stopped/never-started container has no ttyd listening, and no IP
// assigned on the network either.
//
// userID is an explicit parameter rather than resolved from ctx here: this method's only caller
// (internal/transport/vaults_api/workbench_terminal.go) sits outside the gRPC auth interceptor
// chain (it authenticates the request itself, over a cookie, before this is ever called) so ctx
// never carries an injected user_context.UserContext the way an intercepted RPC's does.
//
// The address is resolved per call rather than cached on the workbenches row: Docker reassigns
// a container's network IP on every start, so a stored value would silently go stale across a
// stop/start cycle.
func (s *Service) ResolveTerminalTarget(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
	wb, err := s.workbenchesRepo.GetByVaultAndUser(ctx, vaultID, userID)
	if err != nil {
		return "", rerrors.Wrap(err, "error getting workbench by vault and user")
	}

	if wb.Status != domain.WorkbenchStatusRunning {
		return "", user_errors.WorkbenchNotRunning
	}

	docker, err := s.resolveClient(ctx, wb)
	if err != nil {
		return "", rerrors.Wrap(err, "error resolving docker client")
	}

	address, err := docker.ContainerAddress(ctx, wb.ContainerId)
	if err != nil {
		return "", rerrors.Wrap(err, "error resolving workbench container address")
	}

	return "http://" + address, nil
}

// GetWorkbench returns the calling user's own workbench for vaultID.
func (s *Service) GetWorkbench(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Workbench{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	workbench, err := s.workbenchesRepo.GetByVaultAndUser(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench")
	}

	return workbench, nil
}

// StartWorkbench starts the calling user's workbench container for vaultID under authMode — see
// docs/workbench/03_auth_and_login_flow.md for both modes' designs. Any mode besides the two
// domain.WorkbenchAuthMode constants fails with user_errors.WorkbenchAuthModeNotImplemented
// rather than silently falling back.
func (s *Service) StartWorkbench(
	ctx context.Context, vaultID uuid.UUID, authMode domain.WorkbenchAuthMode,
) (domain.Workbench, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.Workbench{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	switch authMode {
	case domain.WorkbenchAuthModeAPIKey:
		return s.startWithApiKey(ctx, vaultID, uc.UserUuid)
	case domain.WorkbenchAuthModeSubscriptionLogin:
		return s.startWithSubscriptionLogin(ctx, vaultID, uc.UserUuid)
	default:
		return domain.Workbench{}, user_errors.WorkbenchAuthModeNotImplemented
	}
}

// startWithApiKey resolves the workbench owner's BYOK Anthropic api key (failing fast with
// user_errors.WorkbenchMissingAnthropicConnection if none is connected), starts the container
// with it injected as ANTHROPIC_API_KEY, then records status='running'/started_at.
func (s *Service) startWithApiKey(ctx context.Context, vaultID, userID uuid.UUID) (domain.Workbench, error) {
	wb, err := s.workbenchesRepo.GetByVaultAndUser(ctx, vaultID, userID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault and user")
	}

	docker, err := s.resolveClient(ctx, wb)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error resolving docker client")
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

	err = s.workbenchesRepo.MarkConfiguring(ctx, vaultID, userID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error marking workbench configuring")
	}

	err = docker.StartContainer(ctx, wb.ContainerId, env)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error starting workbench container")
	}

	return s.markRunningAndReload(ctx, vaultID, userID, domain.WorkbenchAuthModeAPIKey)
}

// startWithSubscriptionLogin starts the container with no auth env vars injected at all — see
// docs/workbench/03_auth_and_login_flow.md, "Mechanism (confirmed)", step 1. The workbench is
// recorded as status='running' as soon as the container starts, independent of whether the
// in-container login TUI flow has completed — see
// docs/workbench/01_data_model_and_lifecycle.md's state table: "running" means "container
// started", not "authenticated". The user drives the actual in-container login flow themselves,
// interactively, through the workbench terminal (ttyd, proxied by
// internal/transport/vaults_api/workbench_terminal.go) — see ResolveTerminalTarget.
func (s *Service) startWithSubscriptionLogin(ctx context.Context, vaultID, userID uuid.UUID) (domain.Workbench, error) {
	wb, err := s.workbenchesRepo.GetByVaultAndUser(ctx, vaultID, userID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault and user")
	}

	docker, err := s.resolveClient(ctx, wb)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error resolving docker client")
	}

	err = s.workbenchesRepo.MarkConfiguring(ctx, vaultID, userID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error marking workbench configuring")
	}

	err = docker.StartContainer(ctx, wb.ContainerId, nil)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error starting workbench container")
	}

	return s.markRunningAndReload(ctx, vaultID, userID, domain.WorkbenchAuthModeSubscriptionLogin)
}

// markRunningAndReload records status='running'/started_at under authMode and returns the
// refreshed row, shared by both StartWorkbench auth-mode branches.
func (s *Service) markRunningAndReload(
	ctx context.Context, vaultID, userID uuid.UUID, authMode domain.WorkbenchAuthMode,
) (domain.Workbench, error) {
	err := s.workbenchesRepo.MarkRunning(ctx, vaultID, userID, authMode)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error marking workbench running")
	}

	wb, err := s.workbenchesRepo.GetByVaultAndUser(ctx, vaultID, userID)
	if err != nil {
		return domain.Workbench{}, rerrors.Wrap(err, "error getting workbench by vault and user")
	}

	return wb, nil
}

// StopWorkbench stops the calling user's running workbench container for vaultID, retaining it
// (and its volume) for a later restart, and records status='stopped'/stopped_at.
func (s *Service) StopWorkbench(ctx context.Context, vaultID uuid.UUID) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	wb, err := s.workbenchesRepo.GetByVaultAndUser(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return rerrors.Wrap(err, "error getting workbench by vault and user")
	}

	docker, err := s.resolveClient(ctx, wb)
	if err != nil {
		return rerrors.Wrap(err, "error resolving docker client")
	}

	err = docker.StopContainer(ctx, wb.ContainerId)
	if err != nil {
		return rerrors.Wrap(err, "error stopping workbench container")
	}

	err = s.workbenchesRepo.MarkStopped(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return rerrors.Wrap(err, "error marking workbench stopped")
	}

	return nil
}

// DeleteWorkbench tears down the calling user's own workbench for vaultID: stops (if a
// container was created) and removes the container, removes the volume, then deletes the DB row
// — in that order, so a failure partway through never leaves the DB pointing at nothing while
// Docker still holds a live container (see docs/workbench/01_data_model_and_lifecycle.md, "Vault
// deletion"). A missing workbench row is treated as already-deleted, not an error. "Not found"
// from Docker (container/volume already gone, e.g. a retry picking up after a prior partial
// failure) is likewise tolerated at each step rather than failing the whole call.
func (s *Service) DeleteWorkbench(ctx context.Context, vaultID uuid.UUID) error {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return rerrors.Wrap(user_errors.Unauthenticated)
	}

	wb, err := s.workbenchesRepo.GetByVaultAndUser(ctx, vaultID, uc.UserUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return rerrors.Wrap(err, "error getting workbench by vault and user")
	}

	return s.deleteWorkbenchRow(ctx, wb)
}

// DeleteWorkbenchesForVault tears down every member's workbench for vaultID — used by full vault
// deletion, which must clean up every user's Docker resources, not just the caller's. Unlike
// DeleteWorkbench, a failure tearing down one row doesn't stop the rest: every row's teardown is
// attempted, and any failures are joined into a single returned error, so a full vault delete
// never leaves some members' containers orphaned just because another member's teardown hit an
// error.
func (s *Service) DeleteWorkbenchesForVault(ctx context.Context, vaultID uuid.UUID) error {
	rows, err := s.workbenchesRepo.ListByVaultID(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error listing workbenches for vault")
	}

	var teardownErrs []error

	for _, wb := range rows {
		err = s.deleteWorkbenchRow(ctx, wb)
		if err != nil {
			teardownErrs = append(teardownErrs, err)
		}
	}

	if len(teardownErrs) > 0 {
		joined := errors.Join(teardownErrs...)

		return rerrors.Wrap(joined, "error deleting one or more workbenches for vault")
	}

	return nil
}

// deleteWorkbenchRow tears down a single workbench row's Docker container (if created) and
// volume, then deletes the DB row — the shared teardown steps of DeleteWorkbench, factored out
// so DeleteWorkbenchesForVault can run the same logic per-row without duplicating it.
func (s *Service) deleteWorkbenchRow(ctx context.Context, wb domain.Workbench) error {
	docker, err := s.resolveClient(ctx, wb)
	if err != nil {
		return rerrors.Wrap(err, "error resolving docker client")
	}

	if wb.ContainerId != "" {
		err = docker.StopContainer(ctx, wb.ContainerId)
		if err != nil && !errdefs.IsNotFound(err) {
			return rerrors.Wrap(err, "error stopping workbench container")
		}

		err = docker.RemoveContainer(ctx, wb.ContainerId)
		if err != nil && !errdefs.IsNotFound(err) {
			return rerrors.Wrap(err, "error removing workbench container")
		}
	}

	err = docker.RemoveVolume(ctx, wb.VolumeName)
	if err != nil && !errdefs.IsNotFound(err) {
		return rerrors.Wrap(err, "error removing workbench volume")
	}

	err = s.workbenchesRepo.Delete(ctx, wb.VaultUuid, wb.UserUuid)
	if err != nil {
		return rerrors.Wrap(err, "error deleting workbench")
	}

	return nil
}
