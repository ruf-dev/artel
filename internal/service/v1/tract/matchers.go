package tract

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ruf-dev/artel/internal/domain"
)

// MatchesRequest decides whether an inbound delivery should fire a trigger configured with
// matchers — AND semantics across all CheckHeaders/CheckBody entries (same philosophy as
// EvaluateTriggerFilters), so empty matchers always match. Used by the gitlab_webhook handler to
// pick which trigger(s), among possibly several sharing one provider connection, a given
// delivery is for (e.g. gitlab_push only fires on X-Gitlab-Event: Push Hook; gitlab_merge_request
// additionally narrows on CheckBody since GitLab sends the same X-Gitlab-Event: Merge Request
// Hook header for opened/updated/merged/closed alike). body may be empty/nil when matchers has no
// CheckBody entries.
func MatchesRequest(matchers domain.TriggerMatchers, headers http.Header, body json.RawMessage) bool {
	for _, m := range matchers.CheckHeaders {
		if headers.Get(m.Header) != m.Equals {
			return false
		}
	}

	if len(matchers.CheckBody) == 0 {
		return true
	}

	var payload map[string]interface{}

	err := json.Unmarshal(body, &payload)
	if err != nil {
		return false
	}

	for _, m := range matchers.CheckBody {
		if bodyFieldString(payload, m.Path) != m.Equals {
			return false
		}
	}

	return true
}

// bodyFieldString walks a dot-separated path (e.g. "object_attributes.action") into a decoded
// JSON body and returns the string value found there, or "" if any segment is missing or the
// value at that path isn't a string/object as expected.
func bodyFieldString(payload map[string]interface{}, path string) string {
	var current interface{} = payload

	for _, seg := range strings.Split(path, ".") {
		m, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}

		current, ok = m[seg]
		if !ok {
			return ""
		}
	}

	s, _ := current.(string)

	return s
}
