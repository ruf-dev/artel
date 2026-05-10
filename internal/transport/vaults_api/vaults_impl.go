package vaults_api

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

type VaultsImpl struct {
	artel_api.UnimplementedVaultsAPIServer
	vaultSvc service.VaultService
}

func NewVaultsImpl(vaultSvc service.VaultService) *VaultsImpl {
	return &VaultsImpl{vaultSvc: vaultSvc}
}

func (v *VaultsImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterVaultsAPIServer(srv, v)
}

func (v *VaultsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := artel_api.RegisterVaultsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering vaults grpc-gateway handler")
	}

	return "/api/vaults", gwMux
}

func (v *VaultsImpl) CreateVault(ctx context.Context, req *artel_api.CreateVault_Request) (*artel_api.CreateVault_Response, error) {
	err := v.vaultSvc.CreateVault(ctx, req.Name)
	if err != nil {
		return nil, rerrors.Wrap(err, "create vault")
	}

	vault, err := v.vaultSvc.GetVault(ctx, req.Name)
	if err != nil {
		return nil, rerrors.Wrap(err, "get vault after create")
	}

	resp := &artel_api.CreateVault_Response{
		Name:  vault.Name,
		DbUrl: vault.CouchDBURL,
	}
	return resp, nil
}

func (v *VaultsImpl) GetVault(ctx context.Context, req *artel_api.GetVault_Request) (*artel_api.GetVault_Response, error) {
	vault, err := v.vaultSvc.GetVault(ctx, req.Name)
	if err != nil {
		return nil, rerrors.Wrap(err, "get vault")
	}

	resp := &artel_api.GetVault_Response{
		Name:  vault.Name,
		DbUrl: vault.CouchDBURL,
	}
	return resp, nil
}

func (v *VaultsImpl) ListVaults(ctx context.Context, req *artel_api.ListVaults_Request) (*artel_api.ListVaults_Response, error) {
	vaults, err := v.vaultSvc.ListVaults(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vaults")
	}

	items := make([]*artel_api.VaultItem, 0, len(vaults))
	for _, vault := range vaults {
		item := &artel_api.VaultItem{
			Name:  vault.Name,
			DbUrl: vault.CouchDBURL,
		}
		items = append(items, item)
	}

	resp := &artel_api.ListVaults_Response{
		Vaults: items,
	}
	return resp, nil
}

func (v *VaultsImpl) DeleteVault(ctx context.Context, req *artel_api.DeleteVault_Request) (*artel_api.DeleteVault_Response, error) {
	err := v.vaultSvc.DeleteVault(ctx, req.Name)
	if err != nil {
		return nil, rerrors.Wrap(err, "delete vault")
	}

	resp := &artel_api.DeleteVault_Response{
		Ok: true,
	}
	return resp, nil
}
