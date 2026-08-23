package claudesettings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// decodedSettings mirrors the encoded settings shape loosely, just enough for tests to assert on
// without depending on this package's own unexported types.
type decodedSettings struct {
	Hooks struct {
		PreToolUse []struct {
			Matcher string `json:"matcher"`
			Hooks   []struct {
				Type    string `json:"type"`
				Url     string `json:"url"`
				Timeout int    `json:"timeout"`
			} `json:"hooks"`
		} `json:"PreToolUse"`
	} `json:"hooks"`
	Permissions *struct {
		Allow []string `json:"allow"`
	} `json:"permissions"`
}

// writeAndDecode calls Write(path, homeDir) against a temp file and decodes the result, failing
// the test on any error.
func writeAndDecode(t *testing.T, path, homeDir string) decodedSettings {
	t.Helper()

	written, err := Write(path, homeDir)
	if err != nil {
		t.Fatalf("Write(%q, %q) returned error: %v", path, homeDir, err)
	}

	raw, err := os.ReadFile(written)
	if err != nil {
		t.Fatalf("error reading written settings file %s: %v", written, err)
	}

	var decoded decodedSettings

	err = json.Unmarshal(raw, &decoded)
	if err != nil {
		t.Fatalf("error decoding written settings file %s: %v", written, err)
	}

	return decoded
}

// TestWrite_EmptyHomeDir exercises Write's original behavior — no homeDir means no permissions
// block at all and the unrestricted "*" PreToolUse matcher, exactly as before homeDir was added.
func TestWrite_EmptyHomeDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")

	decoded := writeAndDecode(t, path, "")

	if decoded.Permissions != nil {
		t.Errorf("Permissions = %+v, want nil when homeDir is empty", decoded.Permissions)
	}

	if len(decoded.Hooks.PreToolUse) != 1 {
		t.Fatalf("expected exactly one PreToolUse matcher group, got %d", len(decoded.Hooks.PreToolUse))
	}

	if decoded.Hooks.PreToolUse[0].Matcher != allMatcher {
		t.Errorf("Matcher = %q, want %q", decoded.Hooks.PreToolUse[0].Matcher, allMatcher)
	}
}

// TestWrite_WithHomeDir exercises Write's new behavior: a homeDir produces a permissions.allow
// list scoped to <homeDir>/vault/** for Read and Edit, and narrows the PreToolUse matcher to
// exclude those two tool names.
func TestWrite_WithHomeDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")

	decoded := writeAndDecode(t, path, "/root")

	if decoded.Permissions == nil {
		t.Fatal("Permissions is nil, want a populated block when homeDir is set")
	}

	wantAllow := []string{"Read(/root/vault/**)", "Edit(/root/vault/**)"}

	if len(decoded.Permissions.Allow) != len(wantAllow) {
		t.Fatalf("Permissions.Allow = %v, want %v", decoded.Permissions.Allow, wantAllow)
	}

	for i, want := range wantAllow {
		if decoded.Permissions.Allow[i] != want {
			t.Errorf("Permissions.Allow[%d] = %q, want %q", i, decoded.Permissions.Allow[i], want)
		}
	}

	if len(decoded.Hooks.PreToolUse) != 1 {
		t.Fatalf("expected exactly one PreToolUse matcher group, got %d", len(decoded.Hooks.PreToolUse))
	}

	gotMatcher := decoded.Hooks.PreToolUse[0].Matcher
	if gotMatcher != excludeReadEditMatcher {
		t.Errorf("Matcher = %q, want %q", gotMatcher, excludeReadEditMatcher)
	}

	if gotMatcher == allMatcher {
		t.Error("Matcher still \"*\" — Read/Edit are no longer meant to reach the hook once homeDir is set")
	}
}

// TestWrite_HookHandlerAlwaysPresent exercises that the PreToolUse hook handler itself (URL,
// type, timeout) is unaffected by homeDir — narrowing the matcher must not also drop the handler
// wired up to whatever tools the matcher does still cover.
func TestWrite_HookHandlerAlwaysPresent(t *testing.T) {
	for _, homeDir := range []string{"", "/root"} {
		path := filepath.Join(t.TempDir(), "settings.json")

		decoded := writeAndDecode(t, path, homeDir)

		if len(decoded.Hooks.PreToolUse) != 1 || len(decoded.Hooks.PreToolUse[0].Hooks) != 1 {
			t.Fatalf("homeDir=%q: expected exactly one hook handler, got PreToolUse=%+v", homeDir, decoded.Hooks.PreToolUse)
		}

		handler := decoded.Hooks.PreToolUse[0].Hooks[0]
		if handler.Type != "http" {
			t.Errorf("homeDir=%q: handler.Type = %q, want %q", homeDir, handler.Type, "http")
		}

		if handler.Url == "" {
			t.Errorf("homeDir=%q: handler.Url is empty", homeDir)
		}

		if handler.Timeout <= 0 {
			t.Errorf("homeDir=%q: handler.Timeout = %d, want > 0", homeDir, handler.Timeout)
		}
	}
}

// TestExcludeReadEditMatcher_NotValidGoRegexp pins down a fact worth keeping documented: Go's
// regexp package (RE2) rejects excludeReadEditMatcher's negative-lookahead syntax outright — RE2
// deliberately excludes lookaround for its linear-time guarantee. That's expected and fine: this
// Go process only ever writes the pattern into the settings JSON, it never compiles or evaluates
// it — the `claude` CLI's own (Node.js) engine does, and that's a different regex flavor that does
// support lookahead. If this ever starts compiling under Go's regexp, the pattern was likely
// simplified in a way that may have changed its meaning, worth a second look.
func TestExcludeReadEditMatcher_NotValidGoRegexp(t *testing.T) {
	_, err := regexp.Compile(excludeReadEditMatcher)
	if err == nil {
		t.Fatal("expected excludeReadEditMatcher to be rejected by Go's RE2-based regexp package (no lookahead support), but it compiled")
	}
}

// TestWrite_CreatesParentDir exercises that Write still creates path's parent directory when it
// doesn't already exist — unrelated to the homeDir signature change, but a regression here would
// break every caller since DefaultPath's directory is never pre-created.
func TestWrite_CreatesParentDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "dir", "settings.json")

	written, err := Write(path, "")
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	if written != path {
		t.Errorf("Write returned %q, want %q", written, path)
	}

	_, err = os.Stat(written)
	if err != nil {
		t.Errorf("written settings file not found at %s: %v", written, err)
	}
}
