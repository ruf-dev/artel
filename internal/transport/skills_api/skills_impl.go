package skills_api

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"google.golang.org/grpc"
)

type SkillsImpl struct {
	pb.UnimplementedSkillsAPIServer
	skillsSvc service.SkillsService
}

func NewSkillsImpl(skillsSvc service.SkillsService) *SkillsImpl {
	return &SkillsImpl{skillsSvc: skillsSvc}
}

func (s *SkillsImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterSkillsAPIServer(srv, s)
}

func (s *SkillsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := transport.NewGatewayMux()

	err := pb.RegisterSkillsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering skills grpc-gateway handler")
	}

	return "/api/skills/", gwMux
}
