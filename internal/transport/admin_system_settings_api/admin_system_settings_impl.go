package admin_system_settings_api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
)

// AdminSystemSettingsImpl implements the AdminSystemSettingsAPI proto interface — every RPC
// here is admin-only (see GrpcAdminInterceptor's admin-path list in internal/app/custom.go).
type AdminSystemSettingsImpl struct {
	artel_api.UnimplementedAdminSystemSettingsAPIServer
	svc service.AdminSystemSettingsService
}

func New(svc service.AdminSystemSettingsService) *AdminSystemSettingsImpl {
	return &AdminSystemSettingsImpl{svc: svc}
}

func (a *AdminSystemSettingsImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterAdminSystemSettingsAPIServer(srv, a)
}

func (a *AdminSystemSettingsImpl) Gateway(
	ctx context.Context, endpoint string, opts ...grpc.DialOption,
) (string, http.Handler) {
	gwMux := transport.NewGatewayMux()

	err := artel_api.RegisterAdminSystemSettingsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering admin system settings grpc-gateway handler")
	}

	return "/api/admin_system_settings/", gwMux
}

func (a *AdminSystemSettingsImpl) GetSettings(
	ctx context.Context,
	_ *artel_api.GetSettings_Request,
) (*artel_api.GetSettings_Response, error) {
	settings, defaultVault, err := a.svc.GetSettings(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting system settings")
	}

	resp := &artel_api.GetSettings_Response{
		PasswordAuthEnabled: settings.PasswordAuthEnabled,
		TelegramAuthEnabled: settings.TelegramAuthEnabled,
		RegistrationMode:    registrationModeToProto(settings.RegistrationMode),
		DefaultDocsSource:   docsSourceToProto(settings.DefaultDocsSource),
		SystemPrompt:        settings.SystemPrompt,
	}

	if defaultVault.Uuid != uuid.Nil {
		resp.DefaultDocsVaultId = defaultVault.Uuid.String()
		resp.DefaultDocsVaultName = defaultVault.Name
		resp.DefaultDocsVaultSlug = defaultVault.Slug
	}

	return resp, nil
}

func (a *AdminSystemSettingsImpl) UpdateSystemPrompt(ctx context.Context, req *artel_api.UpdateSystemPrompt_Request) (*artel_api.UpdateSystemPrompt_Response, error) {
	err := a.svc.UpdateSystemPrompt(ctx, req.GetPrompt())
	if err != nil {
		return nil, rerrors.Wrap(err, "error updating system prompt")
	}
	return &artel_api.UpdateSystemPrompt_Response{}, nil
}

func (a *AdminSystemSettingsImpl) UpdateDefaultDocsVault(
	ctx context.Context,
	req *artel_api.UpdateDefaultDocsVault_Request,
) (*artel_api.UpdateDefaultDocsVault_Response, error) {
	var vaultUuid *uuid.UUID

	if req.GetVaultId() != "" {
		parsed, err := uuid.Parse(req.GetVaultId())
		if err != nil {
			return nil, rerrors.Wrap(err, "error parsing vault id")
		}

		vaultUuid = &parsed
	}

	source := docsSourceFromProto(req.GetSource())

	err := a.svc.UpdateDefaultDocsVault(ctx, vaultUuid, source)
	if err != nil {
		return nil, rerrors.Wrap(err, "error updating default docs vault")
	}

	return &artel_api.UpdateDefaultDocsVault_Response{}, nil
}

func (a *AdminSystemSettingsImpl) UpdateAuthMethods(
	ctx context.Context,
	req *artel_api.UpdateAuthMethods_Request,
) (*artel_api.UpdateAuthMethods_Response, error) {
	err := a.svc.UpdateAuthMethods(ctx, req.GetPasswordEnabled(), req.GetTelegramEnabled())
	if err != nil {
		return nil, rerrors.Wrap(err, "error updating auth methods")
	}

	return &artel_api.UpdateAuthMethods_Response{}, nil
}

func (a *AdminSystemSettingsImpl) UpdateRegistrationMode(
	ctx context.Context,
	req *artel_api.UpdateRegistrationMode_Request,
) (*artel_api.UpdateRegistrationMode_Response, error) {
	mode := registrationModeFromProto(req.GetMode())

	err := a.svc.UpdateRegistrationMode(ctx, mode)
	if err != nil {
		return nil, rerrors.Wrap(err, "error updating registration mode")
	}

	return &artel_api.UpdateRegistrationMode_Response{}, nil
}

func registrationModeFromProto(mode artel_api.RegistrationMode) domain.RegistrationMode {
	if mode == artel_api.RegistrationMode_SELF_REGISTER {
		return domain.RegistrationModeSelfRegister
	}

	return domain.RegistrationModeAdminOnly
}

func registrationModeToProto(mode domain.RegistrationMode) artel_api.RegistrationMode {
	if mode == domain.RegistrationModeSelfRegister {
		return artel_api.RegistrationMode_SELF_REGISTER
	}

	return artel_api.RegistrationMode_ADMIN_ONLY
}

func docsSourceFromProto(source artel_api.DocsSource) domain.DocsSource {
	if source == artel_api.DocsSource_GITHUB {
		return domain.DocsSourceGithub
	}

	return domain.DocsSourceVault
}

func docsSourceToProto(source domain.DocsSource) artel_api.DocsSource {
	if source == domain.DocsSourceGithub {
		return artel_api.DocsSource_GITHUB
	}

	return artel_api.DocsSource_VAULT
}
