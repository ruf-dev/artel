package tract

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// Webhook trigger source presets — see the Trigger.Source doc comment in domain/tract.go. Preset
// metadata (description/schema/category/provider/default matchers) lives in the trigger_presets
// table (migration 038); these keys are just the stable identifiers NormalizePayload's Go switch
// dispatches on for the presets that need a normalizer.
const (
	SourceGitlabPush         = "gitlab_push"
	SourceGitlabMergeRequest = "gitlab_merge_request"
	SourceGeneric            = "generic"

	gitlabRefBranchPrefix = "refs/heads/"
)

// ListTriggerSources returns the webhook preset catalog backing the ListTriggerSources RPC, read
// from the trigger_presets table (seeded via migration) rather than hardcoded Go data — see
// domain.TriggerPreset.
func (s *Service) ListTriggerSources(ctx context.Context) ([]domain.TriggerPreset, error) {
	presets, err := s.triggerPresets.List(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing trigger presets")
	}

	return presets, nil
}

// NormalizePayload applies the source-specific payload normalizer — gitlab_push adds a "branch"
// field, gitlab_merge_request adds "mr_iid"/"action". Unknown/generic sources pass the raw
// payload through unchanged. Used by the tract webhook handler before filter evaluation and
// StartRun.
func NormalizePayload(source string, raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		raw = json.RawMessage(`{}`)
	}

	switch source {
	case SourceGitlabPush:
		return normalizeGitlabPush(raw)
	case SourceGitlabMergeRequest:
		return normalizeGitlabMergeRequest(raw)
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

// normalizeGitlabMergeRequest surfaces "mr_iid" and "action" (object_attributes.iid/action) at
// the top level, mirroring normalizeGitlabPush's "branch" — the existing GitLab MR tools
// (get_merge_request_diff, update_merge_request) already take a mr_iid param, so this lets a
// tract template read {{ trigger.mr_iid }} directly instead of
// {{ trigger.object_attributes.iid }}.
func normalizeGitlabMergeRequest(raw json.RawMessage) (json.RawMessage, error) {
	var payload map[string]interface{}

	err := json.Unmarshal(raw, &payload)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "error unmarshaling gitlab merge request payload")
	}

	attrs, _ := payload[fieldObjectAttributes].(map[string]interface{})
	payload[fieldMrIid] = attrs[fieldIid]
	payload[fieldAction] = attrs[fieldAction]

	normalized, err := json.Marshal(payload)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshaling normalized gitlab merge request payload")
	}

	return normalized, nil
}
