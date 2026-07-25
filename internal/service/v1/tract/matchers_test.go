package tract

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMatchesRequest_EmptyAlwaysMatches(t *testing.T) {
	headers := http.Header{}

	matched := MatchesRequest(domain.TriggerMatchers{}, headers, nil)

	assert.True(t, matched)
}

func TestMatchesRequest_HeaderMatch(t *testing.T) {
	matchers := domain.TriggerMatchers{
		CheckHeaders: []domain.HeaderMatcher{
			{Header: "X-Gitlab-Event", Equals: "Push Hook"},
		},
	}

	headers := http.Header{}
	headers.Set("X-Gitlab-Event", "Push Hook")

	matched := MatchesRequest(matchers, headers, nil)

	assert.True(t, matched)
}

func TestMatchesRequest_HeaderMismatch(t *testing.T) {
	matchers := domain.TriggerMatchers{
		CheckHeaders: []domain.HeaderMatcher{
			{Header: "X-Gitlab-Event", Equals: "Push Hook"},
		},
	}

	headers := http.Header{}
	headers.Set("X-Gitlab-Event", "Merge Request Hook")

	matched := MatchesRequest(matchers, headers, nil)

	assert.False(t, matched)
}

func TestMatchesRequest_ANDSemanticsAcrossHeaders(t *testing.T) {
	matchers := domain.TriggerMatchers{
		CheckHeaders: []domain.HeaderMatcher{
			{Header: "X-Gitlab-Event", Equals: "Push Hook"},
			{Header: "X-Custom", Equals: "expected"},
		},
	}

	headers := http.Header{}
	headers.Set("X-Gitlab-Event", "Push Hook")
	headers.Set("X-Custom", "unexpected")

	matched := MatchesRequest(matchers, headers, nil)

	assert.False(t, matched)
}

func TestMatchesRequest_BodyMatch(t *testing.T) {
	matchers := domain.TriggerMatchers{
		CheckHeaders: []domain.HeaderMatcher{
			{Header: "X-Gitlab-Event", Equals: "Merge Request Hook"},
		},
		CheckBody: []domain.BodyMatcher{
			{Path: "object_attributes.action", Equals: "merge"},
		},
	}

	headers := http.Header{}
	headers.Set("X-Gitlab-Event", "Merge Request Hook")

	body := json.RawMessage(`{"object_attributes": {"action": "merge", "iid": 7}}`)

	matched := MatchesRequest(matchers, headers, body)

	assert.True(t, matched)
}

func TestMatchesRequest_BodyMismatch(t *testing.T) {
	matchers := domain.TriggerMatchers{
		CheckHeaders: []domain.HeaderMatcher{
			{Header: "X-Gitlab-Event", Equals: "Merge Request Hook"},
		},
		CheckBody: []domain.BodyMatcher{
			{Path: "object_attributes.action", Equals: "merge"},
		},
	}

	headers := http.Header{}
	headers.Set("X-Gitlab-Event", "Merge Request Hook")

	body := json.RawMessage(`{"object_attributes": {"action": "open", "iid": 7}}`)

	matched := MatchesRequest(matchers, headers, body)

	assert.False(t, matched)
}

func TestMatchesRequest_BodyPathMissing(t *testing.T) {
	matchers := domain.TriggerMatchers{
		CheckBody: []domain.BodyMatcher{
			{Path: "object_attributes.action", Equals: "merge"},
		},
	}

	headers := http.Header{}

	body := json.RawMessage(`{"object_attributes": {}}`)

	matched := MatchesRequest(matchers, headers, body)

	assert.False(t, matched)
}

func TestMatchesRequest_BodyMalformedJSON(t *testing.T) {
	matchers := domain.TriggerMatchers{
		CheckBody: []domain.BodyMatcher{
			{Path: "object_attributes.action", Equals: "merge"},
		},
	}

	headers := http.Header{}

	body := json.RawMessage(`not json`)

	matched := MatchesRequest(matchers, headers, body)

	assert.False(t, matched)
}

func TestMatchesRequest_EmptyBody_NoCheckBody_StillMatches(t *testing.T) {
	matchers := domain.TriggerMatchers{
		CheckHeaders: []domain.HeaderMatcher{
			{Header: "X-Gitlab-Event", Equals: "Push Hook"},
		},
	}

	headers := http.Header{}
	headers.Set("X-Gitlab-Event", "Push Hook")

	matched := MatchesRequest(matchers, headers, nil)

	assert.True(t, matched)
}
