// Package claudesettings writes the settings file the bridge passes to every `claude` invocation
// via --settings, which is what registers the PreToolUse http hook.
//
// The nesting below is not obvious and getting it wrong fails *silently* — Claude Code simply
// never calls the hook, and every gated tool call is then resolved by its own permission flow
// instead of by a human. It was verified end to end against CLI 2.1.234 by pointing a real
// `claude -p` at a throwaway HTTP server and observing the POST arrive:
//
//	{
//	  "hooks": {
//	    "PreToolUse": [
//	      { "matcher": "*", "hooks": [ { "type": "http", "url": "...", "timeout": 300 } ] }
//	    ]
//	  }
//	}
//
// Note the double nesting: the outer array holds *matcher groups*, each of which holds its own
// "hooks" array of handlers. "*" matches every tool.
package claudesettings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"workbenchbridge/internal/permissions"
)

// DefaultPath is where the bridge writes the settings file. Deliberately not $HOME/.claude/
// settings.json: an explicit --settings path is unambiguous about which file is in effect, and
// leaves the user's own in-workspace settings free to layer on top without either clobbering the
// other.
const DefaultPath = "/run/workbench/claude-settings.json"

type settings struct {
	Hooks       hooks             `json:"hooks"`
	Permissions *permissionsBlock `json:"permissions,omitempty"`
}

type hooks struct {
	PreToolUse []matcherGroup `json:"PreToolUse"`
}

// matcherGroup is one entry of the PreToolUse array: which tools it applies to, and the handlers
// to run for them.
type matcherGroup struct {
	Matcher string    `json:"matcher"`
	Hooks   []handler `json:"hooks"`
}

type handler struct {
	Type    string `json:"type"`
	Url     string `json:"url"`
	Timeout int    `json:"timeout"`
}

// permissionsBlock is the "permissions" object in the settings file (named to avoid colliding
// with the sibling workbenchbridge/internal/permissions package this file already imports for the
// hook URL/timeout). Read/Edit calls matching Allow are resolved locally by the `claude` CLI
// itself, without ever reaching the PreToolUse hook — see Write's doc comment on why that's
// paired with narrowing the hook's own matcher.
type permissionsBlock struct {
	Allow []string `json:"allow"`
}

// Write renders the settings file at path (empty means DefaultPath) and returns the path written,
// for the caller to pass to `claude --settings`.
//
// When homeDir is non-empty, the settings file also gets a permissions.allow list auto-approving
// Read and Edit under <homeDir>/vault/** — the workbench's mounted volume (see homeMountPath in
// internal/clients/workbenchdocker/client.go, the main repo, a different Go module this one can't
// import) — so browsing/editing vault files never needs a round trip through the PreToolUse hook.
// This is a silent, no-chat-indicator bypass by design: it only ever widens what two tools that
// are already read-only-or-reversible on vault content can touch, never anything destructive or
// outside the vault. Passing homeDir="" (as tests may) skips the permissions block entirely,
// leaving every tool call gated by the hook exactly as before this field was added.
//
// The PreToolUse hook's own Matcher is narrowed from "*" to excludeReadEditMatcher accordingly:
// Read/Edit no longer need to reach the hook at all once permissions.allow covers them, and
// leaving them in the hook's matcher too would just mean the CLI's own local allow-list decides
// first and the hook is never asked — but *not* narrowing it would be silently equivalent, not
// silently wrong, so this is a belt-and-suspenders exclusion, not a correctness requirement.
func Write(path, homeDir string) (string, error) {
	if path == "" {
		path = DefaultPath
	}

	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return "", err
	}

	hookHandler := handler{
		Type:    "http",
		Url:     permissions.HookUrl,
		Timeout: permissions.HookTimeoutSeconds,
	}

	group := matcherGroup{
		Matcher: matcherFor(homeDir),
		Hooks:   []handler{hookHandler},
	}

	content := settings{
		Hooks: hooks{
			PreToolUse: []matcherGroup{group},
		},
		Permissions: permissionsFor(homeDir),
	}

	encoded, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return "", err
	}

	err = os.WriteFile(path, encoded, 0o644)
	if err != nil {
		return "", err
	}

	return path, nil
}

// allMatcher matches every tool — the original, always-safe PreToolUse matcher used whenever no
// homeDir is configured (nothing is auto-allowed locally, so every tool call must still reach the
// hook).
const allMatcher = "*"

// excludeReadEditMatcher matches every tool name except exactly "Read" and "Edit" — a negative
// lookahead. The matcher field is just a JSON string written here and evaluated entirely inside
// the `claude` CLI's own (Node.js) engine — this Go process never compiles or evaluates it itself
// — so it only has to be valid there, not as a Go regexp: confirmed by hand that Go's regexp
// package (RE2) rejects this exact pattern with "invalid or unsupported Perl syntax: `(?!`",
// since RE2 deliberately excludes lookaround for its linear-time guarantee, while JS's engine (and
// most other mainstream ones) supports it. That difference is exactly why this still needs the
// same manual verification the file's package doc comment describes for the original "*" matcher
// (verified against CLI 2.1.234 with a throwaway HTTP server) — pointing a real `claude -p` at a
// throwaway hook server and confirming Read/Edit calls never arrive while everything else still
// does — before trusting this in production. A silently-wrong matcher here fails open: every tool
// call, including ones meant to stay gated, would simply never reach the hook.
const excludeReadEditMatcher = `^(?!Read$|Edit$).*$`

// matcherFor returns the PreToolUse matcher to use: excludeReadEditMatcher when homeDir is
// configured (Read/Edit are locally auto-allowed via permissions.allow instead — see Write's doc
// comment), else the original allMatcher.
func matcherFor(homeDir string) string {
	if homeDir == "" {
		return allMatcher
	}

	return excludeReadEditMatcher
}

// permissionsFor returns the permissions block Write embeds when homeDir is non-empty, or nil
// (omitted from the encoded JSON, via the "omitempty" tag) when it's empty.
func permissionsFor(homeDir string) *permissionsBlock {
	if homeDir == "" {
		return nil
	}

	vaultGlob := homeDir + "/vault/**"

	allow := []string{
		"Read(" + vaultGlob + ")",
		"Edit(" + vaultGlob + ")",
	}

	return &permissionsBlock{Allow: allow}
}
