package vaults_api

import (
	"time"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

const watchWorkbenchLoginPollInterval = 700 * time.Millisecond

// WatchWorkbenchLogin streams periodic snapshots of a workbench's subscription_login TUI flow
// (docs/workbench/03_auth_and_login_flow.md) until it reaches domain.WorkbenchLoginStateAuthorized
// or the client disconnects. A thin polling wrapper around WorkbenchService.GetLoginPrompt,
// modeled directly on TractsImpl.WatchRun (internal/transport/tracts_api/watch_run.go) — tmux
// capture-pane is cheap, so the same poll interval is reused rather than inventing a new one.
func (v *VaultsImpl) WatchWorkbenchLogin(
	req *pb.WatchWorkbenchLogin_Request, stream pb.VaultsAPI_WatchWorkbenchLoginServer,
) error {
	if v.workbenchSvc == nil {
		return user_errors.WorkbenchNotConfigured
	}

	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return rerrors.Wrap(err, "parse vault id")
	}

	for {
		prompt, err := v.workbenchSvc.GetLoginPrompt(stream.Context(), vaultID)
		if err != nil {
			return rerrors.Wrap(err, "get workbench login prompt")
		}

		resp := &pb.WatchWorkbenchLogin_Response{
			State:        string(prompt.State),
			Url:          prompt.URL,
			ErrorMessage: prompt.ErrorMessage,
		}

		err = stream.Send(resp)
		if err != nil {
			return rerrors.Wrap(err, "error sending workbench login prompt snapshot")
		}

		if prompt.State == domain.WorkbenchLoginStateAuthorized {
			return nil
		}

		select {
		case <-stream.Context().Done():
			return nil
		case <-time.After(watchWorkbenchLoginPollInterval):
		}
	}
}
