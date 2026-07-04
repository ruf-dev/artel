package tasktracker

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/trello"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

type Service struct {
	trackers repository.TaskTrackerRepo
}

func New(trackers repository.TaskTrackerRepo) *Service {
	return &Service{trackers: trackers}
}

func (s *Service) AddTracker(
	ctx context.Context,
	tracker domain.TaskTracker,
	creds trello.TaskTrackerCredentials,
) (domain.TaskTracker, []domain.TrelloBoard, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.TaskTracker{}, nil, user_errors.Unauthenticated
	}

	tc := creds.GetClient()

	member, err := tc.GetMember(ctx)
	if err != nil {
		if errors.Is(err, trello.ErrUnauthorized) {
			return domain.TaskTracker{}, nil, user_errors.TrelloInvalidCredentials
		}

		return domain.TaskTracker{}, nil, rerrors.Wrap(err, "error verifying trello credentials")
	}

	boards, err := tc.ListBoards(ctx)
	if err != nil {
		return domain.TaskTracker{}, nil, rerrors.Wrap(err, "error listing trello boards")
	}

	if tCreds, ok := creds.(trello.TrelloCredentials); ok {
		tracker.ApiKey = tCreds.ApiKey
		tracker.ApiToken = tCreds.ApiToken
	}

	tracker.UserUuid = uc.UserUuid
	tracker.Name = member.FullName

	created, err := s.trackers.Insert(ctx, tracker)
	if err != nil {
		return domain.TaskTracker{}, nil, rerrors.Wrap(err, "error inserting task tracker")
	}

	return created, boards, nil
}

func (s *Service) ListTrackers(ctx context.Context) ([]domain.TaskTracker, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, user_errors.Unauthenticated
	}

	trackers, err := s.trackers.ListByUser(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing task trackers")
	}

	return trackers, nil
}

func (s *Service) DeleteTracker(ctx context.Context, trackerUuid uuid.UUID) error {
	err := s.trackers.Delete(ctx, trackerUuid)
	if err != nil {
		return rerrors.Wrap(err, "error deleting task tracker")
	}

	return nil
}

func (s *Service) ListTrelloBoards(ctx context.Context, trackerUuid uuid.UUID) ([]domain.TrelloBoard, error) {
	tracker, err := s.trackers.GetByUuid(ctx, trackerUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting task tracker")
	}

	tc := trello.New(tracker.ApiKey, tracker.ApiToken)

	boards, err := tc.ListBoards(ctx)
	if err != nil {
		if errors.Is(err, trello.ErrUnauthorized) {
			return nil, user_errors.TrelloInvalidCredentials
		}

		return nil, rerrors.Wrap(err, "error listing trello boards")
	}

	return boards, nil
}
