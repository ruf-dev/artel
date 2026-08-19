// Package envdrop reads (and writes) the env-file drop directory the backend injects a workbench
// container's secrets into.
//
// The Docker Engine API cannot attach environment variables to an already-created container, so
// internal/clients/workbenchdocker's injectEnv instead docker-execs one `printf %s "$VALUE" >
// $DIR/$NAME` per variable into a tmpfs-mounted directory (mode 0700). The bridge reads that
// directory back before every `claude` invocation rather than caching it once at startup, because
// the backend injects into a container that is already running — the drop can (and for
// subscription login, does) change after the bridge has started.
package envdrop

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	// DefaultDir is the in-container path the backend drops env files into. Must match envDropDir
	// in internal/clients/workbenchdocker/container.go.
	DefaultDir = "/run/workbench/env"

	// AnthropicApiKeyVar / OauthTokenVar / AuthModeVar are the three variables the bridge itself
	// cares about; anything else found in the drop directory is still passed through to `claude`
	// untouched.
	AnthropicApiKeyVar = "ANTHROPIC_API_KEY"
	OauthTokenVar      = "CLAUDE_CODE_OAUTH_TOKEN"
	AuthModeVar        = "WORKBENCH_AUTH_MODE"

	// AuthModeApiKey / AuthModeSubscriptionLogin are AuthModeVar's two values, matching
	// domain.WorkbenchAuthMode's wire values in the main module.
	AuthModeApiKey            = "api_key"
	AuthModeSubscriptionLogin = "subscription_login"
)

// Store reads and writes one env-file drop directory.
type Store struct {
	dir string
}

// New returns a Store over dir. An empty dir means DefaultDir.
func New(dir string) *Store {
	if dir == "" {
		dir = DefaultDir
	}

	return &Store{dir: dir}
}

// Read returns every variable currently present in the drop directory, keyed by file name. A
// missing directory is not an error — it just means the backend hasn't injected anything yet
// (or, in a locally-run bridge, that there is no drop at all), which every caller has to tolerate
// anyway.
//
// Values are taken byte-exactly except for a single trailing newline, which is stripped: the
// backend writes with `printf %s` and so never adds one, but a human debugging a container with
// `echo` would, and a stray "\n" inside an API key is a uniquely annoying failure to diagnose.
func (s *Store) Read() (map[string]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}

		return nil, err
	}

	values := make(map[string]string, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		content, err := os.ReadFile(filepath.Join(s.dir, entry.Name()))
		if err != nil {
			return nil, err
		}

		values[entry.Name()] = strings.TrimSuffix(string(content), "\n")
	}

	return values, nil
}

// Write persists one variable into the drop directory, creating the directory if the tmpfs mount
// left it empty. Used by the subscription-login flow to store the OAuth token `claude
// setup-token` produces, so subsequent turns (and any later `claude` invocation in this
// container's lifetime) pick it up through the same path the backend's own injection uses.
//
// The 0700/0600 modes mirror the `umask 077` the backend's injectEnv exec applies, keeping a
// bridge-written file indistinguishable from a backend-written one.
func (s *Store) Write(name, value string) error {
	err := os.MkdirAll(s.dir, 0o700)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(s.dir, name), []byte(value), 0o600)
}

// Environ renders values as a `KEY=VALUE` slice suitable for exec.Cmd.Env, merged on top of base
// (typically os.Environ()): a key present in both wins from values, so a freshly injected secret
// always beats a stale one inherited from the bridge's own process environment.
//
// The result is sorted, purely so it is deterministic and therefore testable.
func Environ(base []string, values map[string]string) []string {
	merged := make(map[string]string, len(base)+len(values))

	for _, entry := range base {
		key, value, found := strings.Cut(entry, "=")
		if !found {
			continue
		}

		merged[key] = value
	}

	for key, value := range values {
		merged[key] = value
	}

	environ := make([]string, 0, len(merged))

	for key, value := range merged {
		environ = append(environ, key+"="+value)
	}

	sort.Strings(environ)

	return environ
}
