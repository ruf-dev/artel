package tract

import (
	"net/http"
	"testing"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMatchesRequest_EmptyAlwaysMatches(t *testing.T) {
	headers := http.Header{}

	matched := MatchesRequest(domain.TriggerMatchers{}, headers)

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

	matched := MatchesRequest(matchers, headers)

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

	matched := MatchesRequest(matchers, headers)

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

	matched := MatchesRequest(matchers, headers)

	assert.False(t, matched)
}
