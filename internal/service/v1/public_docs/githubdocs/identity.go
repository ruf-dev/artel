package githubdocs

import (
	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
)

// githubIdentityUUID is a fixed, obviously-synthetic UUID that can never collide with a real
// vault's generated UUID — it identifies the GitHub-backed docs "vault" wherever a domain.Vault
// is expected (e.g. as the resolved vault for domain.ReservedGithubDocsSlug).
var githubIdentityUUID = uuid.MustParse("00000000-0000-0000-0000-0000000decaf")

// IdentityVault returns a synthetic domain.Vault representing the GitHub-backed docs "vault".
// Only Uuid/Name/Slug/IsPublic are ever read by callers, so every other field is left
// zero-valued.
func IdentityVault() domain.Vault {
	vault := domain.Vault{
		Uuid:     githubIdentityUUID,
		Name:     "Artel Quick Start",
		Slug:     domain.ReservedGithubDocsSlug,
		IsPublic: true,
	}

	return vault
}
