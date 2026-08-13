package vault

import (
	"errors"
	"testing"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// TestValidateSlug_RejectsReservedGithubDocsSlug ensures the reserved sentinel slug used to
// route `/docs/:slug` to the GitHub-backed quick-start resolver (public_docs.Service) can never
// be claimed by a real vault via PublishVault.
func TestValidateSlug_RejectsReservedGithubDocsSlug(t *testing.T) {
	err := validateSlug(domain.ReservedGithubDocsSlug)
	if !errors.Is(err, user_errors.VaultSlugInvalid) {
		t.Fatalf("expected user_errors.VaultSlugInvalid, got %v", err)
	}
}

// TestValidateSlug_AcceptsRegularKebabCaseSlug is a regression guard confirming the reserved-slug
// check added above doesn't reject ordinary valid slugs.
func TestValidateSlug_AcceptsRegularKebabCaseSlug(t *testing.T) {
	err := validateSlug("my-public-docs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
