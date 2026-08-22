package authlogin

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writeCredentials writes a minimal claude .credentials.json at path with the given access token.
func writeCredentials(t *testing.T, path, accessToken string) {
	t.Helper()

	err := os.MkdirAll(filepath.Dir(path), 0o700)
	if err != nil {
		t.Fatalf("error creating credentials dir: %v", err)
	}

	content := `{"claudeAiOauth":{"accessToken":"` + accessToken + `"}}`

	err = os.WriteFile(path, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("error writing credentials file: %v", err)
	}
}

func TestCredentialsPath_DefaultsToHomeClaudeDir(t *testing.T) {
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	t.Setenv("HOME", "/home/example")

	got := credentialsPath()
	want := filepath.Join("/home/example", ".claude", ".credentials.json")

	if got != want {
		t.Fatalf("credentialsPath() = %q, want %q", got, want)
	}
}

func TestCredentialsPath_HonorsConfigDirOverride(t *testing.T) {
	t.Setenv("CLAUDE_CONFIG_DIR", "/custom/config")

	got := credentialsPath()
	want := filepath.Join("/custom/config", ".credentials.json")

	if got != want {
		t.Fatalf("credentialsPath() = %q, want %q", got, want)
	}
}

func TestReadFreshToken_MissingFileIsNotAnError(t *testing.T) {
	t.Setenv("CLAUDE_CONFIG_DIR", filepath.Join(t.TempDir(), "does-not-exist"))

	token, found, err := readFreshToken(time.Now())
	if err != nil {
		t.Fatalf("readFreshToken() error = %v, want nil", err)
	}

	if found {
		t.Fatalf("readFreshToken() found = true, want false for a missing file")
	}

	if token != "" {
		t.Fatalf("readFreshToken() token = %q, want empty", token)
	}
}

func TestReadFreshToken_StaleFileIsIgnored(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", configDir)

	writeCredentials(t, filepath.Join(configDir, ".credentials.json"), "sk-ant-oat01-stale")

	// since is in the future relative to the file just written, simulating a credentials file
	// left over from an earlier login in this same $HOME rather than one produced by this run.
	since := time.Now().Add(time.Hour)

	_, found, err := readFreshToken(since)
	if err != nil {
		t.Fatalf("readFreshToken() error = %v, want nil", err)
	}

	if found {
		t.Fatalf("readFreshToken() found = true, want false for a file older than since")
	}
}

func TestReadFreshToken_FreshFileReturnsAccessToken(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", configDir)

	since := time.Now()

	writeCredentials(t, filepath.Join(configDir, ".credentials.json"), "sk-ant-oat01-fresh")

	token, found, err := readFreshToken(since)
	if err != nil {
		t.Fatalf("readFreshToken() error = %v, want nil", err)
	}

	if !found {
		t.Fatalf("readFreshToken() found = false, want true for a freshly written file")
	}

	if token != "sk-ant-oat01-fresh" {
		t.Fatalf("readFreshToken() token = %q, want %q", token, "sk-ant-oat01-fresh")
	}
}

func TestReadFreshToken_EmptyAccessTokenIsNotFound(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", configDir)

	since := time.Now()

	writeCredentials(t, filepath.Join(configDir, ".credentials.json"), "")

	_, found, err := readFreshToken(since)
	if err != nil {
		t.Fatalf("readFreshToken() error = %v, want nil", err)
	}

	if found {
		t.Fatalf("readFreshToken() found = true, want false when accessToken is empty")
	}
}
