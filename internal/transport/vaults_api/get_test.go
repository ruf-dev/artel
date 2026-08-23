package vaults_api

import (
	"errors"
	"testing"
)

// TestResolvePendingAuthLink covers resolvePendingAuthLink's decision table: an empty cached link
// short-circuits regardless of loggedIn/err, a confirmed login clears the cache and reports no
// link, a confirmed non-login keeps reporting the cached link without clearing it, and a check
// error fails open — the cached link is still reported and the cache is left alone.
func TestResolvePendingAuthLink(t *testing.T) {
	checkErr := errors.New("boom")

	tests := []struct {
		name            string
		cachedLink      string
		loggedIn        bool
		checkErr        error
		wantPendingLink string
		wantShouldClear bool
	}{
		{
			name:            "empty cached link, logged in",
			cachedLink:      "",
			loggedIn:        true,
			checkErr:        nil,
			wantPendingLink: "",
			wantShouldClear: false,
		},
		{
			name:            "empty cached link, not logged in",
			cachedLink:      "",
			loggedIn:        false,
			checkErr:        nil,
			wantPendingLink: "",
			wantShouldClear: false,
		},
		{
			name:            "empty cached link, check errored",
			cachedLink:      "",
			loggedIn:        false,
			checkErr:        checkErr,
			wantPendingLink: "",
			wantShouldClear: false,
		},
		{
			name:            "cached link, logged in",
			cachedLink:      "https://claude.ai/oauth/authorize",
			loggedIn:        true,
			checkErr:        nil,
			wantPendingLink: "",
			wantShouldClear: true,
		},
		{
			name:            "cached link, not logged in",
			cachedLink:      "https://claude.ai/oauth/authorize",
			loggedIn:        false,
			checkErr:        nil,
			wantPendingLink: "https://claude.ai/oauth/authorize",
			wantShouldClear: false,
		},
		{
			name:            "cached link, check errored (fails open, loggedIn irrelevant)",
			cachedLink:      "https://claude.ai/oauth/authorize",
			loggedIn:        true,
			checkErr:        checkErr,
			wantPendingLink: "https://claude.ai/oauth/authorize",
			wantShouldClear: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pendingLink, shouldClear := resolvePendingAuthLink(tt.cachedLink, tt.loggedIn, tt.checkErr)

			if pendingLink != tt.wantPendingLink {
				t.Fatalf("pendingLink = %q, want %q", pendingLink, tt.wantPendingLink)
			}

			if shouldClear != tt.wantShouldClear {
				t.Fatalf("shouldClear = %v, want %v", shouldClear, tt.wantShouldClear)
			}
		})
	}
}
