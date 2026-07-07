package tract

import (
	"net/http"

	"github.com/ruf-dev/artel/internal/domain"
)

// MatchesRequest decides whether an inbound delivery should fire a trigger configured with
// matchers — AND semantics across all CheckHeaders entries (same philosophy as
// EvaluateTriggerFilters), so empty matchers always match. Used by the gitlab_webhook handler to
// pick which trigger(s), among possibly several sharing one provider connection, a given
// delivery is for (e.g. gitlab_push only fires on X-Gitlab-Event: Push Hook).
func MatchesRequest(matchers domain.TriggerMatchers, headers http.Header) bool {
	for _, m := range matchers.CheckHeaders {
		if headers.Get(m.Header) != m.Equals {
			return false
		}
	}

	return true
}
