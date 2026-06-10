package admin_couch_api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"

	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
)

type AdminCouchImpl struct {
	artel_api.UnimplementedAdminCouchAPIServer
	svc service.AdminCouchService
}

func New(svc service.AdminCouchService) *AdminCouchImpl {
	return &AdminCouchImpl{svc: svc}
}

func (a *AdminCouchImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterAdminCouchAPIServer(srv, a)
}

func (a *AdminCouchImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := artel_api.RegisterAdminCouchAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering admin couch grpc-gateway handler")
	}

	return "/api/admin_couch/", gwMux
}

func (a *AdminCouchImpl) ListCouchUsers(ctx context.Context, req *artel_api.ListCouchUsers_Request) (*artel_api.ListCouchUsers_Response, error) {
	users, err := a.svc.ListUsers(ctx, req.InstanceId)
	if err != nil {
		return nil, rerrors.Wrap(err, "listing couch users")
	}

	entries := make([]*artel_api.CouchUserEntry, len(users))
	for i, u := range users {
		entry := &artel_api.CouchUserEntry{Name: u.Name, Roles: u.Roles}
		entries[i] = entry
	}

	resp := &artel_api.ListCouchUsers_Response{Users: entries}
	return resp, nil
}

func (a *AdminCouchImpl) DeleteCouchUser(ctx context.Context, req *artel_api.DeleteCouchUser_Request) (*artel_api.DeleteCouchUser_Response, error) {
	err := a.svc.DeleteUser(ctx, req.InstanceId, req.Username)
	if err != nil {
		return nil, rerrors.Wrap(err, "deleting couch user")
	}

	resp := &artel_api.DeleteCouchUser_Response{}
	return resp, nil
}

func (a *AdminCouchImpl) ChangeCouchUserPassword(ctx context.Context, req *artel_api.ChangeCouchUserPassword_Request) (*artel_api.ChangeCouchUserPassword_Response, error) {
	err := a.svc.ChangeUserPassword(ctx, req.InstanceId, req.Username, req.NewPassword)
	if err != nil {
		return nil, rerrors.Wrap(err, "changing couch user password")
	}

	resp := &artel_api.ChangeCouchUserPassword_Response{}
	return resp, nil
}

func (a *AdminCouchImpl) ListCouchDatabases(ctx context.Context, req *artel_api.ListCouchDatabases_Request) (*artel_api.ListCouchDatabases_Response, error) {
	databases, err := a.svc.ListDatabases(ctx, req.InstanceId)
	if err != nil {
		return nil, rerrors.Wrap(err, "listing couch databases")
	}

	resp := &artel_api.ListCouchDatabases_Response{Databases: databases}
	return resp, nil
}

func (a *AdminCouchImpl) GrantDatabaseAccess(ctx context.Context, req *artel_api.GrantDatabaseAccess_Request) (*artel_api.GrantDatabaseAccess_Response, error) {
	err := a.svc.GrantDatabaseAccess(ctx, req.InstanceId, req.DbName, req.Username)
	if err != nil {
		return nil, rerrors.Wrap(err, "granting database access")
	}

	resp := &artel_api.GrantDatabaseAccess_Response{}
	return resp, nil
}

func (a *AdminCouchImpl) RevokeDatabaseAccess(ctx context.Context, req *artel_api.RevokeDatabaseAccess_Request) (*artel_api.RevokeDatabaseAccess_Response, error) {
	err := a.svc.RevokeDatabaseAccess(ctx, req.InstanceId, req.DbName, req.Username)
	if err != nil {
		return nil, rerrors.Wrap(err, "revoking database access")
	}

	resp := &artel_api.RevokeDatabaseAccess_Response{}
	return resp, nil
}

func (a *AdminCouchImpl) GetUserDatabaseAccess(ctx context.Context, req *artel_api.GetUserDatabaseAccess_Request) (*artel_api.GetUserDatabaseAccess_Response, error) {
	databases, err := a.svc.GetUserDatabaseAccess(ctx, req.InstanceId, req.Username)
	if err != nil {
		return nil, rerrors.Wrap(err, "getting user database access")
	}

	resp := &artel_api.GetUserDatabaseAccess_Response{Databases: databases}
	return resp, nil
}
