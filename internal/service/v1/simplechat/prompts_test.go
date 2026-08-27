package simplechat

import "testing"

func TestMergePrompts(t *testing.T) {
	got := mergePrompts("admin", "personal", "vault", true)
	want := "System instructions:\nadmin\n\nUser instructions:\npersonal\n\nVault instructions:\nvault"
	if got != want {
		t.Fatalf("mergePrompts() = %q, want %q", got, want)
	}
}

func TestMergePromptsOmitsDisabledOrEmptySections(t *testing.T) {
	got := mergePrompts("admin", "  ", "vault", false)
	want := "Vault instructions:\nvault"
	if got != want {
		t.Fatalf("mergePrompts() = %q, want %q", got, want)
	}
}
