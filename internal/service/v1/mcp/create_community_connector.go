package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// communityConnectorInput is the wire shape create_community_connector accepts — the same
// {name, description, tools[{api_description, action}]} MoM record documented in
// pkg/mom_examples/README.md. It reuses domain.ToolApiDescription/domain.ToolAction directly:
// every one of their fields is a single word, so encoding/json's case-insensitive field
// matching (e.g. "properties" -> Properties, "operation" -> Operation) maps the wire keys onto
// them without any extra tags or a second recursive parser — only the multi-word
// api_description key needs an explicit tag.
type communityConnectorInput struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Tools       []communityConnectorTool `json:"tools"`
}

type communityConnectorTool struct {
	ApiDescription domain.ToolApiDescription `json:"api_description"`
	Action         domain.ToolAction         `json:"action"`
}

// createCommunityConnector implements the create_community_connector builtin tool: an
// administrator's external AI agent has already turned a user's Swagger/OpenAPI spec into an
// MoM-shaped JSON tool definition itself, and calls this to persist it. Only administrators may
// call it successfully; re-calling with the same name upserts the caller's own existing
// connector but is rejected against a system MoM or another admin's connector.
func (s *ServiceImpl) createCommunityConnector(
	ctx context.Context,
	userUuid uuid.UUID,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	err := s.authSvc.CheckIsAdmin(ctx, userUuid)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "check caller is admin")
	}

	input, err := parseCommunityConnectorInput(params)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	err = validateCommunityConnectorInput(input)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	existing, err := s.mcpDefinitions.Get(ctx, input.Name)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "get existing mcp definition")
	}

	if existing.Valid {
		ownedBySomeoneElse := existing.V.OwnerUserUuid == nil || *existing.V.OwnerUserUuid != userUuid
		if ownedBySomeoneElse {
			return domain.ToolExecResult{}, user_errors.McpCommunityConnectorNameTaken
		}
	}

	author, err := s.authSvc.GetMe(ctx, userUuid)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "get calling admin's identity")
	}

	tools := make([]domain.McpToolDef, len(input.Tools))
	for i, t := range input.Tools {
		tools[i] = domain.McpToolDef{ApiDescription: t.ApiDescription, Action: t.Action}
	}

	def := domain.McpDefinition{
		Name:          input.Name,
		Author:        author.Email,
		Description:   input.Description,
		Tools:         tools,
		OwnerUserUuid: &userUuid,
		IsCommunity:   true,
	}

	saved, err := s.mcpDefinitions.Upsert(ctx, def)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "upsert mcp definition")
	}

	text := fmt.Sprintf("Community connector %q created with %d tool(s).", saved.Name, len(saved.Tools))
	result := domain.ToolExecResult{Text: text}

	return result, nil
}

func parseCommunityConnectorInput(params map[string]interface{}) (communityConnectorInput, error) {
	raw, err := json.Marshal(params)
	if err != nil {
		return communityConnectorInput{}, rerrors.Wrap(err, "error marshaling create_community_connector params")
	}

	var input communityConnectorInput

	err = json.Unmarshal(raw, &input)
	if err != nil {
		return communityConnectorInput{}, rerrors.Wrap(err, "error unmarshaling create_community_connector params")
	}

	return input, nil
}

func validateCommunityConnectorInput(input communityConnectorInput) error {
	if input.Name == "" {
		return user_errors.McpCommunityConnectorNameRequired
	}

	if input.Description == "" {
		return user_errors.McpCommunityConnectorDescriptionRequired
	}

	if len(input.Tools) == 0 {
		return user_errors.McpCommunityConnectorToolsRequired
	}

	return nil
}
