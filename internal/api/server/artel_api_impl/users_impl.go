package artel_api_impl

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"

	artel_api "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
)

type UsersImpl struct {
	artel_api.UnimplementedUsersAPIServer
	userSvc service.UserService
}

func NewUsersImpl(userSvc service.UserService) *UsersImpl {
	return &UsersImpl{userSvc: userSvc}
}

func (u *UsersImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterUsersAPIServer(srv, u)
}

func (u *UsersImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := artel_api.RegisterUsersAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering users grpc-gateway handler")
	}

	return "/api/users", gwMux
}

func (u *UsersImpl) CreateUser(ctx context.Context, req *artel_api.CreateUser_Request) (*artel_api.CreateUser_Response, error) {
	err := u.userSvc.CreateUser(ctx, req.Username, req.Password, req.Roles)
	if err != nil {
		return nil, rerrors.Wrap(err, "create user")
	}

	resp := &artel_api.CreateUser_Response{
		Username: req.Username,
		Roles:    req.Roles,
	}
	return resp, nil
}

func (u *UsersImpl) GetUser(ctx context.Context, req *artel_api.GetUser_Request) (*artel_api.GetUser_Response, error) {
	user, err := u.userSvc.GetUser(ctx, req.Username)
	if err != nil {
		return nil, rerrors.Wrap(err, "get user")
	}

	resp := &artel_api.GetUser_Response{
		Username: user.Username,
		Roles:    user.Roles,
	}
	return resp, nil
}

func (u *UsersImpl) UpdateUser(ctx context.Context, req *artel_api.UpdateUser_Request) (*artel_api.UpdateUser_Response, error) {
	err := u.userSvc.UpdateUser(ctx, req.Username, req.Password, req.Roles)
	if err != nil {
		return nil, rerrors.Wrap(err, "update user")
	}

	resp := &artel_api.UpdateUser_Response{
		Username: req.Username,
		Roles:    req.Roles,
	}
	return resp, nil
}

func (u *UsersImpl) DeleteUser(ctx context.Context, req *artel_api.DeleteUser_Request) (*artel_api.DeleteUser_Response, error) {
	err := u.userSvc.DeleteUser(ctx, req.Username)
	if err != nil {
		return nil, rerrors.Wrap(err, "delete user")
	}

	resp := &artel_api.DeleteUser_Response{
		Ok: true,
	}
	return resp, nil
}
