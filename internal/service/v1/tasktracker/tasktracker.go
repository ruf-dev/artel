// Package tasktracker backs the TaskTrackersAPI wire contract (unchanged) on top of the
// generic external_connections + MoM path instead of a dedicated task_trackers table/client.
// Every "tracker" is a trello external_connections row; "boards" are fetched live via MoM's
// trello list_boards tool rather than a bespoke HTTP client.
package tasktracker

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

const trelloMcpName = "trello"

// trelloConnectionMeta mirrors the unexported type of the same name in
// internal/service/v1/externalconnections/external_connections.go — the JSON shape
// AddTrelloConnection writes into external_connections.metadata ({"full_name": "..."}).
// Duplicated here (rather than exported) since ListTrackers needs to read raw
// external_connections rows directly to avoid ExternalConnectionService.ListConnections'
// generic (Google-only) display-name extraction, which doesn't know about trello's metadata
// shape.
type trelloConnectionMeta struct {
	FullName string `json:"full_name"`
}

type Service struct {
	externalConns    repository.ExternalConnectionRepo
	externalConnsSvc service.ExternalConnectionService
	momSvc           service.MomService
}

func New(
	externalConns repository.ExternalConnectionRepo,
	externalConnsSvc service.ExternalConnectionService,
	momSvc service.MomService,
) *Service {
	return &Service{
		externalConns:    externalConns,
		externalConnsSvc: externalConnsSvc,
		momSvc:           momSvc,
	}
}

func (s *Service) AddTracker(
	ctx context.Context,
	apiKey, apiToken string,
) (domain.TaskTracker, []domain.TrelloBoard, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.TaskTracker{}, nil, user_errors.Unauthenticated
	}

	meta, err := s.externalConnsSvc.AddTrelloConnection(ctx, apiKey, apiToken)
	if err != nil {
		return domain.TaskTracker{}, nil, rerrors.Wrap(err, "error adding trello connection")
	}

	boards, err := s.listBoards(ctx, meta.Uuid)
	if err != nil {
		return domain.TaskTracker{}, nil, err
	}

	tracker := domain.TaskTracker{
		Uuid:      meta.Uuid,
		UserUuid:  uc.UserUuid,
		Type:      domain.ProviderTrello,
		Name:      meta.DisplayName,
		CreatedAt: meta.CreatedAt,
	}

	return tracker, boards, nil
}

func (s *Service) ListTrackers(ctx context.Context) ([]domain.TaskTracker, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	conns, err := s.externalConns.ListByUser(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing external connections")
	}

	var trackers []domain.TaskTracker

	for _, conn := range conns {
		if conn.Provider != domain.ProviderTrello {
			continue
		}

		trackers = append(trackers, toTaskTracker(conn))
	}

	return trackers, nil
}

func (s *Service) DeleteTracker(ctx context.Context, trackerUuid uuid.UUID) error {
	err := s.externalConnsSvc.DisconnectConnection(ctx, trackerUuid.String())
	if err != nil {
		return rerrors.Wrap(err, "error deleting task tracker")
	}

	return nil
}

func (s *Service) ListTrelloBoards(ctx context.Context, trackerUuid uuid.UUID) ([]domain.TrelloBoard, error) {
	boards, err := s.listBoards(ctx, trackerUuid)
	if err != nil {
		return nil, err
	}

	return boards, nil
}

// listBoards runs the trello MoM's list_boards tool against exConnUuid, verifying ownership
// (ExecuteToolForUserConnection checks exConn.UserUuid against the caller's user_context).
func (s *Service) listBoards(ctx context.Context, exConnUuid uuid.UUID) ([]domain.TrelloBoard, error) {
	params := map[string]interface{}{}

	result, err := s.momSvc.ExecuteToolForUserConnection(ctx, exConnUuid, trelloMcpName, "list_boards", params)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing trello boards")
	}

	var boards []domain.TrelloBoard

	err = json.Unmarshal([]byte(result), &boards)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing trello boards")
	}

	return boards, nil
}

func toTaskTracker(conn domain.ExternalConnection) domain.TaskTracker {
	var meta trelloConnectionMeta

	err := json.Unmarshal(conn.Metadata, &meta)
	if err != nil {
		log.Error().Err(err).Str("connection", conn.Uuid.String()).Msg("error parsing trello connection metadata")
	}

	return domain.TaskTracker{
		Uuid:      conn.Uuid,
		UserUuid:  conn.UserUuid,
		Type:      domain.ProviderTrello,
		Name:      meta.FullName,
		CreatedAt: conn.CreatedAt,
	}
}
