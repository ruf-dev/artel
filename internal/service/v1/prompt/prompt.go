package prompt

import (
	"context"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"go.redsock.ru/rerrors"
)

type Service struct {
	repo repository.Prompts
}

func New(repo repository.Prompts) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListPrompts(ctx context.Context, params service.ListPromptsParams) ([]domain.Prompt, int64, error) {
	var offset uint32
	if params.Page > 0 && params.PageSize > 0 {
		offset = (params.Page - 1) * params.PageSize
	}

	result, total, err := s.repo.List(ctx, repository.ListPromptsParams{
		Ids:    params.Ids,
		Limit:  params.PageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, rerrors.Wrap(err, "list prompts")
	}

	return result, total, nil
}
