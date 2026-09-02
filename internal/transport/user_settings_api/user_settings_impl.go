// Package user_settings_api exposes UserSettingsService over the network as
// artel_api.UserSettingsAPI: per-user preference storage synced across devices.
package user_settings_api

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"go.redsock.ru/rerrors"
)

type UserSettingsImpl struct {
	artel_api.UnimplementedUserSettingsAPIServer

	userSettingsSvc service.UserSettingsService
}

func New(userSettingsSvc service.UserSettingsService) *UserSettingsImpl {
	return &UserSettingsImpl{userSettingsSvc: userSettingsSvc}
}

func (s *UserSettingsImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterUserSettingsAPIServer(srv, s)
}

func (s *UserSettingsImpl) Gateway(
	ctx context.Context, endpoint string, opts ...grpc.DialOption,
) (string, http.Handler) {
	gwMux := transport.NewGatewayMux()

	err := artel_api.RegisterUserSettingsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering user settings grpc-gateway handler")
	}

	return "/api/user-settings/", gwMux
}

func (s *UserSettingsImpl) GetUserSettings(
	ctx context.Context, _ *artel_api.GetUserSettings_Request,
) (*artel_api.GetUserSettings_Response, error) {
	settings, err := s.userSettingsSvc.GetUserSettings(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting user settings")
	}

	resp := artel_api.GetUserSettings_Response{
		UserPrompt:            settings.UserPrompt,
		LikedOpenrouterModels: settings.LikedOpenrouterModels,
		LastUsedModel:         settings.LastUsedModel,
	}

	return &resp, nil
}

func (s *UserSettingsImpl) SetLikedModels(
	ctx context.Context, req *artel_api.SetLikedModels_Request,
) (*artel_api.SetLikedModels_Response, error) {
	err := s.userSettingsSvc.SetLikedModels(ctx, req.GetLikedOpenrouterModels())
	if err != nil {
		return nil, rerrors.Wrap(err, "error setting liked models")
	}

	return &artel_api.SetLikedModels_Response{}, nil
}

func (s *UserSettingsImpl) SetLastUsedModel(
	ctx context.Context, req *artel_api.SetLastUsedModel_Request,
) (*artel_api.SetLastUsedModel_Response, error) {
	err := s.userSettingsSvc.SetLastUsedModel(ctx, req.GetModel())
	if err != nil {
		return nil, rerrors.Wrap(err, "error setting last used model")
	}

	return &artel_api.SetLastUsedModel_Response{}, nil
}
