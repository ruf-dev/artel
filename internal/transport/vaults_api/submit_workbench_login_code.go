package vaults_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// SubmitWorkbenchLoginCode relays the user's pasted OAuth code (or, on a brand-new container,
// their first theme/login-menu keystroke) into a workbench's tmux session — see
// docs/workbench/03_auth_and_login_flow.md, "Mechanism (confirmed)", step 4.
func (v *VaultsImpl) SubmitWorkbenchLoginCode(
	ctx context.Context, req *pb.SubmitWorkbenchLoginCode_Request,
) (*pb.SubmitWorkbenchLoginCode_Response, error) {
	if v.workbenchSvc == nil {
		return nil, user_errors.WorkbenchNotConfigured
	}

	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	err = v.workbenchSvc.SubmitLoginCode(ctx, vaultID, req.Code)
	if err != nil {
		return nil, rerrors.Wrap(err, "submit workbench login code")
	}

	resp := &pb.SubmitWorkbenchLoginCode_Response{}

	return resp, nil
}
