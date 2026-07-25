package tract

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizePayload_GitlabPush_AddsBranch(t *testing.T) {
	raw := json.RawMessage(`{"ref": "refs/heads/feature-x"}`)

	normalized, err := NormalizePayload(SourceGitlabPush, raw)
	require.NoError(t, err)

	var payload map[string]interface{}
	err = json.Unmarshal(normalized, &payload)
	require.NoError(t, err)

	assert.Equal(t, "feature-x", payload[fieldBranch])
}

func TestNormalizePayload_GitlabMergeRequest_SurfacesMrIidAndAction(t *testing.T) {
	raw := json.RawMessage(`{
		"object_attributes": {"iid": 7, "action": "merge", "title": "Fix bug"},
		"project": {"id": 42}
	}`)

	normalized, err := NormalizePayload(SourceGitlabMergeRequest, raw)
	require.NoError(t, err)

	var payload map[string]interface{}
	err = json.Unmarshal(normalized, &payload)
	require.NoError(t, err)

	assert.Equal(t, float64(7), payload[fieldMrIid])
	assert.Equal(t, "merge", payload[fieldAction])
	// object_attributes itself is left intact — normalizer only adds top-level convenience fields.
	attrs, ok := payload[fieldObjectAttributes].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Fix bug", attrs["title"])
}

func TestNormalizePayload_GitlabMergeRequest_MissingObjectAttributes(t *testing.T) {
	raw := json.RawMessage(`{"project": {"id": 42}}`)

	normalized, err := NormalizePayload(SourceGitlabMergeRequest, raw)
	require.NoError(t, err)

	var payload map[string]interface{}
	err = json.Unmarshal(normalized, &payload)
	require.NoError(t, err)

	assert.Nil(t, payload[fieldMrIid])
	assert.Nil(t, payload[fieldAction])
}

func TestNormalizePayload_Generic_Passthrough(t *testing.T) {
	raw := json.RawMessage(`{"foo": "bar"}`)

	normalized, err := NormalizePayload(SourceGeneric, raw)
	require.NoError(t, err)

	assert.JSONEq(t, `{"foo": "bar"}`, string(normalized))
}

func TestNormalizePayload_EmptyRaw_DefaultsToEmptyObject(t *testing.T) {
	normalized, err := NormalizePayload(SourceGeneric, nil)
	require.NoError(t, err)

	assert.JSONEq(t, `{}`, string(normalized))
}
