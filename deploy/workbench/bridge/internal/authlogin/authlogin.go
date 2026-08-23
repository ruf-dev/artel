// Package authlogin passively detects an OAuth token for `claude` becoming available in this
// container, without driving any sign-in flow itself.
//
// The bridge used to spawn `claude setup-token` under a pty and drive it to completion. That was
// removed: it competed with a user logging in through the container's own interactive terminal
// (a separate `claude` session reachable independently of this process — see
// internal/transport/vaults_api/workbench_terminal_shell.go in the main module), and nothing in
// the current frontend can complete the device-flow prompts it needed anyway. Now this package
// only watches for the credentials file `claude login` (and `claude setup-token`, if ever run by
// hand) writes as a side effect of a successful sign-in, however that sign-in happened.
//
// # The credentials file
//
// `claude` persists a successful login to $CLAUDE_CONFIG_DIR/.credentials.json (default
// $HOME/.claude/.credentials.json) — a {"claudeAiOauth":{"accessToken":"sk-ant-...",...}}
// document, read by every later `claude` invocation (headless or interactive) regardless of env
// vars. $CLAUDE_CONFIG_DIR's resolution order matters here because it's the same override the
// bridge's own claudecli.Runner and claudesettings package have to agree with — a mismatch would
// have this package watching a file `claude` itself never reads or writes.
package authlogin

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"workbenchbridge/internal/envdrop"
)

// credentialsPollInterval bounds how quickly WaitForLogin notices a credentials file appearing
// after a user completes a terminal login. Short enough not to be perceptible against a human
// finishing a login flow.
const credentialsPollInterval = 500 * time.Millisecond

// WaitForLogin blocks until claude's credentials file holds a valid OAuth token, then persists it
// into envStore and returns nil.
//
// An already-existing, already-valid credentials file counts immediately — there is no freshness
// or mtime gating here, deliberately: the whole point of this package's passive design is that a
// terminal login which happened before this call started (even a container restart or two ago)
// is just as good as one that happens while this call is waiting. If no valid file is present
// yet, WaitForLogin polls every credentialsPollInterval until one appears or ctx is cancelled, in
// which case it returns ctx.Err().
func WaitForLogin(ctx context.Context, envStore *envdrop.Store) error {
	token, found, err := readFreshToken(time.Time{})
	if err == nil && found {
		return envStore.Write(envdrop.OauthTokenVar, token)
	}

	ticker := time.NewTicker(credentialsPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			token, found, err := readFreshToken(time.Time{})
			if err == nil && found {
				return envStore.Write(envdrop.OauthTokenVar, token)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// credentialsPath resolves where `claude` stores its own persistent credentials, matching the
// resolution order confirmed in the CLI's own bundled source: $CLAUDE_CONFIG_DIR/.credentials.json
// when that env var is set, else $HOME/.claude/.credentials.json.
func credentialsPath() string {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, ".credentials.json")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "/root"
	}

	return filepath.Join(home, ".claude", ".credentials.json")
}

// credentialsFile is the subset of claude's own .credentials.json this package reads — see the
// package doc's "The credentials file" section.
type credentialsFile struct {
	ClaudeAiOauth struct {
		AccessToken string `json:"accessToken"`
	} `json:"claudeAiOauth"`
}

// readFreshToken reports the access token in claude's credentials file, but only if that file's
// mtime is after since — a missing file, or one that predates since, is not an error, just a
// "not yet" (found=false). Passing the zero time.Time disables this gating entirely, since every
// real file's ModTime is after it; WaitForLogin always does this, but the parameter is kept
// (rather than dropping it and always reading unconditionally) so a future caller that does need
// freshness gating doesn't have to reintroduce it here.
func readFreshToken(since time.Time) (string, bool, error) {
	path := credentialsPath()

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}

		return "", false, err
	}

	if !info.ModTime().After(since) {
		return "", false, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", false, err
	}

	var creds credentialsFile

	err = json.Unmarshal(data, &creds)
	if err != nil {
		return "", false, err
	}

	if creds.ClaudeAiOauth.AccessToken == "" {
		return "", false, nil
	}

	return creds.ClaudeAiOauth.AccessToken, true, nil
}
