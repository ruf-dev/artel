package admin_users_api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AdminUsersImpl struct {
	artel_api.UnimplementedAdminUsersAPIServer
	svc service.AdminUsersService
}

func New(svc service.AdminUsersService) *AdminUsersImpl {
	return &AdminUsersImpl{svc: svc}
}

func (a *AdminUsersImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterAdminUsersAPIServer(srv, a)
}

func (a *AdminUsersImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := artel_api.RegisterAdminUsersAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering admin users grpc-gateway handler")
	}

	return "/api/admin_users/", gwMux
}

func (a *AdminUsersImpl) ListArtelUsers(
	ctx context.Context,
	req *artel_api.ListArtelUsers_Request,
) (*artel_api.ListArtelUsers_Response, error) {
	paging := domain.Paging{Limit: req.Limit, Offset: req.Offset}
	listReq := domain.ListUsersReq{Paging: paging, Search: req.Search}

	users, total, err := a.svc.ListUsers(ctx, listReq)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing artel users")
	}

	entries := make([]*artel_api.ArtelUserEntry, len(users))

	for i, u := range users {
		entry := &artel_api.ArtelUserEntry{
			UserId:    u.Uuid.String(),
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: timestamppb.New(u.CreatedAt),
		}
		entries[i] = entry
	}

	resp := &artel_api.ListArtelUsers_Response{Users: entries, Total: total}

	return resp, nil
}

func (a *AdminUsersImpl) GetArtelUser(
	ctx context.Context,
	req *artel_api.GetArtelUser_Request,
) (*artel_api.GetArtelUser_Response, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing user id")
	}

	details, err := a.svc.GetUser(ctx, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting artel user")
	}

	userDetails := &artel_api.ArtelUserDetails{
		UserId:             details.Uuid.String(),
		Username:           details.Username,
		Email:              details.Email,
		IsAdministrator:    details.Permissions.IsAdministrator,
		HasEmails:          details.EffectiveSubscription.Features.Emails,
		HasTaskTrackers:    details.EffectiveSubscription.Features.TaskTrackers,
		HasNotes:           details.EffectiveSubscription.Features.Notes,
		HasSpreadsheets:    details.EffectiveSubscription.Features.Spreadsheets,
		SubscriptionActive: details.EffectiveSubscription.Active,
	}
	resp := &artel_api.GetArtelUser_Response{User: userDetails}

	return resp, nil
}

func (a *AdminUsersImpl) GetUserSessions(
	ctx context.Context,
	req *artel_api.GetUserSessions_Request,
) (*artel_api.GetUserSessions_Response, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing user id")
	}

	sessions, err := a.svc.GetUserSessions(ctx, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting user sessions")
	}

	entries := make([]*artel_api.UserSession, len(sessions))

	for i, s := range sessions {
		entry := &artel_api.UserSession{
			SessionId: s.Uuid.String(),
			CreatedAt: timestamppb.New(s.CreatedAt),
			ExpiresAt: timestamppb.New(s.ExpiresAt),
		}
		entries[i] = entry
	}

	resp := &artel_api.GetUserSessions_Response{Sessions: entries}

	return resp, nil
}
