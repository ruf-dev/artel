// Package simple_chat_api exposes Simple Chat over the network: the CRUD RPCs of
// artel_api.SimpleChatAPI here, and the live turn WebSocket in chat_ws.go. The two halves live
// in one package so the wire contract for a chat thread stays in one place.
package simple_chat_api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"go.redsock.ru/rerrors"
)

type SimpleChatImpl struct {
	artel_api.UnimplementedSimpleChatAPIServer

	simpleChatSvc service.SimpleChatService
}

func New(simpleChatSvc service.SimpleChatService) *SimpleChatImpl {
	return &SimpleChatImpl{simpleChatSvc: simpleChatSvc}
}

func (s *SimpleChatImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterSimpleChatAPIServer(srv, s)
}

func (s *SimpleChatImpl) Gateway(
	ctx context.Context, endpoint string, opts ...grpc.DialOption,
) (string, http.Handler) {
	gwMux := transport.NewGatewayMux()

	err := artel_api.RegisterSimpleChatAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering simple chat grpc-gateway handler")
	}

	return "/api/simple-chats/", gwMux
}

func (s *SimpleChatImpl) CreateSimpleChat(
	ctx context.Context, req *artel_api.CreateSimpleChat_Request,
) (*artel_api.CreateSimpleChat_Response, error) {
	vaultUuid, err := uuid.Parse(req.GetVaultId())
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	chat, err := s.simpleChatSvc.CreateChat(ctx, vaultUuid, req.GetModel(), req.GetVaultAccess())
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating simple chat")
	}

	resp := artel_api.CreateSimpleChat_Response{
		Chat: chatToProto(chat),
	}

	return &resp, nil
}

func (s *SimpleChatImpl) ListSimpleChats(
	ctx context.Context, req *artel_api.ListSimpleChats_Request,
) (*artel_api.ListSimpleChats_Response, error) {
	vaultUuid, err := uuid.Parse(req.GetVaultId())
	if err != nil {
		return nil, rerrors.Wrap(err, "parse vault id")
	}

	chats, err := s.simpleChatSvc.ListChats(ctx, vaultUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing simple chats")
	}

	resp := artel_api.ListSimpleChats_Response{
		Chats: chatsToProto(chats),
	}

	return &resp, nil
}

func (s *SimpleChatImpl) GetSimpleChat(
	ctx context.Context, req *artel_api.GetSimpleChat_Request,
) (*artel_api.GetSimpleChat_Response, error) {
	chatUuid, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, rerrors.Wrap(err, "parse chat id")
	}

	chat, messages, err := s.simpleChatSvc.GetChat(ctx, chatUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting simple chat")
	}

	resp := artel_api.GetSimpleChat_Response{
		Chat:     chatToProto(chat),
		Messages: messagesToProto(messages),
	}

	return &resp, nil
}

func (s *SimpleChatImpl) DeleteSimpleChat(
	ctx context.Context, req *artel_api.DeleteSimpleChat_Request,
) (*artel_api.DeleteSimpleChat_Response, error) {
	chatUuid, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, rerrors.Wrap(err, "parse chat id")
	}

	err = s.simpleChatSvc.DeleteChat(ctx, chatUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error deleting simple chat")
	}

	resp := artel_api.DeleteSimpleChat_Response{}

	return &resp, nil
}
