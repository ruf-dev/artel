package tract

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// Webhook trigger source presets — see the Trigger.Source doc comment in domain/tract.go.
const (
	SourceGitlabPush = "gitlab_push"
	SourceGeneric    = "generic"

	gitlabRefBranchPrefix = "refs/heads/"
)

// ListTriggerSources returns the webhook preset catalog backing the ListTriggerSources RPC —
// pure Go data, no repo access. vault_event/schedule/chat_keyword presets are future `kind`s,
// not sources, and are not listed here.
func ListTriggerSources() []domain.TriggerSourcePreset {
	return []domain.TriggerSourcePreset{
		{
			Key:           SourceGitlabPush,
			Description:   "GitLab push event webhook",
			PayloadSchema: gitlabPushPayloadSchema(),
		},
		{
			Key:           SourceGeneric,
			Description:   "Generic webhook — payload shape is declared per-trigger (payload_schema), passed through unchanged",
			PayloadSchema: domain.ToolSchema{},
		},
	}
}

// ListTriggerSources exposes the package-level preset catalog through service.TractService.
func (s *Service) ListTriggerSources(_ context.Context) ([]domain.TriggerSourcePreset, error) {
	return ListTriggerSources(), nil
}

// NormalizePayload applies the source-specific payload normalizer — currently only
// gitlab_push adds a "branch" field. Unknown/generic sources pass the raw payload through
// unchanged. Used by the tract webhook handler before filter evaluation and StartRun.
func NormalizePayload(source string, raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		raw = json.RawMessage(`{}`)
	}

	switch source {
	case SourceGitlabPush:
		return normalizeGitlabPush(raw)
	default:
		return raw, nil
	}
}

// normalizeGitlabPush adds "branch" = ref with the "refs/heads/" prefix trimmed, per the
// design's inspector template requirements ({{ trigger.branch }}).
func normalizeGitlabPush(raw json.RawMessage) (json.RawMessage, error) {
	var payload map[string]interface{}

	err := json.Unmarshal(raw, &payload)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "error unmarshaling gitlab push payload")
	}

	ref, _ := payload[fieldRef].(string)
	payload[fieldBranch] = strings.TrimPrefix(ref, gitlabRefBranchPrefix)

	normalized, err := json.Marshal(payload)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshaling normalized gitlab push payload")
	}

	return normalized, nil
}

// gitlabPushPayloadSchema describes the shape NormalizePayload produces for a gitlab_push
// trigger: ref/branch/user_name/project{id,name,path}/commits[]{id,message,url,author{name,
// email}}/total_commits_count.
func gitlabPushPayloadSchema() domain.ToolSchema {
	author := domain.ToolProperty{
		Type: schemaTypeObject,
		Properties: map[string]domain.ToolProperty{
			fieldName: {Type: "string"},
			"email":   {Type: "string"},
		},
	}

	commitItem := domain.ToolProperty{
		Type: schemaTypeObject,
		Properties: map[string]domain.ToolProperty{
			"id":         {Type: "string", Description: "commit sha"},
			fieldMessage: {Type: "string"},
			"url":        {Type: "string"},
			"author":     author,
		},
	}

	project := domain.ToolProperty{
		Type: schemaTypeObject,
		Properties: map[string]domain.ToolProperty{
			"id":      {Type: "integer"},
			fieldName: {Type: "string"},
			"path":    {Type: "string"},
		},
	}

	schema := domain.ToolSchema{
		Properties: map[string]domain.ToolProperty{
			fieldRef: {Type: "string", Description: "full git ref, e.g. refs/heads/feature-x"},
			fieldBranch: {
				Type:        "string",
				Description: "ref with refs/heads/ trimmed (added by the normalizer)",
			},
			"user_name":           {Type: "string"},
			"project":             project,
			"commits":             {Type: "array", Items: &commitItem},
			"total_commits_count": {Type: "integer"},
		},
		Required: []string{fieldRef, fieldBranch},
	}

	return schema
}
