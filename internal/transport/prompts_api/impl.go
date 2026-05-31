package prompts_api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
)

var promptIdToString = map[pb.PromptId]string{
	pb.PromptId_PROMPT_ID_PERSONAL_VAULT_SETUP: "personal_vault_setup",
}

var stringToPromptId = map[string]pb.PromptId{
	"personal_vault_setup": pb.PromptId_PROMPT_ID_PERSONAL_VAULT_SETUP,
}

type PromptsImpl struct {
	pb.UnimplementedPromptsAPIServer
	promptSvc service.PromptService
}

func NewPromptsImpl(promptSvc service.PromptService) *PromptsImpl {
	return &PromptsImpl{promptSvc: promptSvc}
}

func (p *PromptsImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterPromptsAPIServer(srv, p)
}

func (p *PromptsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := pb.RegisterPromptsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering prompts grpc-gateway handler")
	}

	return "/api/prompts/", gwMux
}

func (p *PromptsImpl) ListPrompts(ctx context.Context, req *pb.ListPrompts_Request) (*pb.ListPrompts_Response, error) {
	var ids []string
	for _, pid := range req.GetIds() {
		if s, ok := promptIdToString[pid]; ok {
			ids = append(ids, s)
		}
	}

	prompts, total, err := p.promptSvc.ListPrompts(ctx, service.ListPromptsParams{
		Ids:      ids,
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
	})
	if err != nil {
		return nil, err
	}

	items := make([]*pb.PromptItem, 0, len(prompts))
	for _, pr := range prompts {
		pid := stringToPromptId[pr.Id]
		items = append(items, &pb.PromptItem{
			Id:     pid,
			Prompt: pr.Text,
		})
	}

	return &pb.ListPrompts_Response{
		Prompts: items,
		Total:   uint32(total),
	}, nil
}
