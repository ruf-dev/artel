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
	Hooks hooks `json:"hooks"`
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

// Write renders the settings file at path (empty means DefaultPath) and returns the path written,
// for the caller to pass to `claude --settings`.
func Write(path string) (string, error) {
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
		// "*" rather than an enumeration of tool names: any tool Claude Code gains in a future
		// version must land in front of a human too, not slip through unreviewed.
		Matcher: "*",
		Hooks:   []handler{hookHandler},
	}

	content := settings{
		Hooks: hooks{
			PreToolUse: []matcherGroup{group},
		},
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
