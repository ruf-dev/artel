package email_accounts_api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
)

type EmailAccountsImpl struct {
	pb.UnimplementedEmailAccountsAPIServer
	emailSvc service.EmailService
}

func NewEmailAccountsImpl(emailSvc service.EmailService) *EmailAccountsImpl {
	return &EmailAccountsImpl{emailSvc: emailSvc}
}

func (e *EmailAccountsImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterEmailAccountsAPIServer(srv, e)
}

func (e *EmailAccountsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := pb.RegisterEmailAccountsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering email accounts grpc-gateway handler")
	}

	return "/api/email-accounts/", gwMux
}
