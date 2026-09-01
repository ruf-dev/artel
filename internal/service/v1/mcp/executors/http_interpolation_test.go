package executors

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRenderHttpBody_InterpolationFallback exercises the ${{...}} resolver through the body
// renderer, so it also covers renderBodyValue dropping a key whose rendered value is "".
func TestRenderHttpBody_InterpolationFallback(t *testing.T) {
	cases := []struct {
		name    string
		body    string
		params  map[string]interface{}
		secrets map[string]interface{}
		want    map[string]interface{}
	}{
		{
			name:   "single param term present",
			body:   `{"chat_id":"${{params.chat_id}}","text":"${{params.text}}"}`,
			params: map[string]interface{}{"chat_id": "111", "text": "hi"},
			want:   map[string]interface{}{"chat_id": "111", "text": "hi"},
		},
		{
			name:   "single term absent renders empty and drops the key",
			body:   `{"chat_id":"${{params.chat_id}}","text":"hi"}`,
			params: map[string]interface{}{},
			want:   map[string]interface{}{"text": "hi"},
		},
		{
			name:    "fallback chain keeps the left term when present",
			body:    `{"chat_id":"${{params.chat_id | secrets.chat_id}}"}`,
			params:  map[string]interface{}{"chat_id": "from-params"},
			secrets: map[string]interface{}{"chat_id": "from-secrets"},
			want:    map[string]interface{}{"chat_id": "from-params"},
		},
		{
			name:    "fallback chain falls through to the secret when the param is absent",
			body:    `{"chat_id":"${{params.chat_id | secrets.chat_id}}"}`,
			params:  map[string]interface{}{},
			secrets: map[string]interface{}{"chat_id": "from-secrets"},
			want:    map[string]interface{}{"chat_id": "from-secrets"},
		},
		{
			name:    "fallback chain with both terms absent drops the key",
			body:    `{"chat_id":"${{params.chat_id | secrets.chat_id}}","text":"hi"}`,
			params:  map[string]interface{}{},
			secrets: map[string]interface{}{},
			want:    map[string]interface{}{"text": "hi"},
		},
		{
			name:    "whitespace around the pipe and dot is tolerated",
			body:    `{"chat_id":"${{  params . chat_id  |  secrets . chat_id  }}"}`,
			params:  map[string]interface{}{},
			secrets: map[string]interface{}{"chat_id": "from-secrets"},
			want:    map[string]interface{}{"chat_id": "from-secrets"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rendered, err := renderHttpBody(json.RawMessage(tc.body), tc.params, tc.secrets)
			require.NoError(t, err)

			var got map[string]interface{}

			err = json.Unmarshal(rendered, &got)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

// TestInterpolateParams_LeavesForeignNamespaceTokenIntact guards the URL two-pass flow:
// interpolateParams must not consume a ${{secrets.*}} token, so the later interpolateSecrets
// pass can still resolve it.
func TestInterpolateParams_LeavesForeignNamespaceTokenIntact(t *testing.T) {
	rawUrl := "https://api.telegram.org/bot${{secrets.bot_token}}/sendMessage"

	afterParams := interpolateParams(rawUrl, map[string]interface{}{})
	require.Equal(t, rawUrl, afterParams)

	afterSecrets := interpolateSecrets(afterParams, map[string]interface{}{"bot_token": "T"})
	require.Equal(t, "https://api.telegram.org/botT/sendMessage", afterSecrets)
}
