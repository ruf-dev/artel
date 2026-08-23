package vaults_api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

// resolvePendingAuthLink decides what GetVault should report for pending_terminal_auth_link and
// whether the cached entry should now be cleared, given the terminal-scrollback-detected
// cachedLink and a live docker-exec login check. checkErr fails open — a transient exec failure
// shouldn't flip a real pending banner off, so the cached link is still reported and the cache is
// left alone to be re-checked on the next GetVault call.
func resolvePendingAuthLink(cachedLink string, loggedIn bool, checkErr error) (pendingLink string, shouldClear bool) {
	if cachedLink == "" {
		return "", false
	}

	if checkErr != nil {
		return cachedLink, false
	}

	if loggedIn {
		return "", true
	}

	return cachedLink, false
}

func (v *VaultsImpl) GetVault(ctx context.Context, req *pb.GetVault_Request) (*pb.GetVault_Response, error) {
	vaultID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	vault, err := v.vaultSvc.GetVault(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "get vault")
	}

	s3InstanceId := ""
	if vault.S3InstanceUuid != nil {
		s3InstanceId = vault.S3InstanceUuid.String()
	}

	workbenchExists := false
	workbenchStatus := ""
	wb, wbErr := v.workbenchSvc.GetWorkbench(ctx, vaultID)
	if wbErr == nil {
		workbenchExists = true
		workbenchStatus = string(wb.Status)
	} else if !errors.Is(wbErr, sql.ErrNoRows) {
		// A vault without a workbench row is an expected, non-fatal state (predates this
		// feature, or the create-vault hook's CreateWorkbench call failed separately) — only
		// log genuine lookup failures, never fail GetVault over supplementary data.
		log.Error().Err(wbErr).Str("vault_id", req.Id).Msg("error getting workbench status for vault")
	}

	postgresEnabled := false
	postgresStatus := "not_enabled"
	pgDB, pgErr := v.vaultSvc.GetPostgresDatabase(ctx, vaultID)
	if pgErr != nil {
		// Mirrors the workbench lookup above: a lookup failure over supplementary data must not
		// fail GetVault as a whole, only skip enrichment for this field.
		log.Error().Err(pgErr).Str("vault_id", req.Id).Msg("error getting postgres database status for vault")
	} else if pgDB.Valid {
		postgresEnabled = true
		postgresStatus = string(pgDB.V.Status)
	}

	pendingAuthLink := ""
	if workbenchExists && wb.Status == domain.WorkbenchStatusRunning {
		cachedLink := v.terminalAuthLinks.PendingTerminalAuthLink(vaultID)

		if cachedLink != "" {
			loggedIn, checkErr := v.workbenchSvc.IsClaudeLoggedIn(ctx, vaultID)
			if checkErr != nil {
				log.Error().Err(checkErr).Str("vault_id", req.Id).Msg("error checking claude login status for vault")
			}

			var shouldClear bool

			pendingAuthLink, shouldClear = resolvePendingAuthLink(cachedLink, loggedIn, checkErr)

			if shouldClear {
				v.terminalAuthLinks.ClearPendingTerminalAuthLink(vaultID)
			}
		}
	}

	resp := &pb.GetVault_Response{
		Id:                      vault.Uuid.String(),
		Name:                    vault.Name,
		DbUrl:                   vault.CouchDBURL,
		S3InstanceId:            s3InstanceId,
		S3BucketName:            vault.S3BucketName,
		WorkbenchExists:         workbenchExists,
		WorkbenchStatus:         workbenchStatus,
		PostgresEnabled:         postgresEnabled,
		PostgresStatus:          postgresStatus,
		PendingTerminalAuthLink: pendingAuthLink,
	}

	return resp, nil
}
