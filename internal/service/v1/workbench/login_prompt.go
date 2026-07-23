package workbench

import (
	"regexp"
	"strings"

	"github.com/ruf-dev/artel/internal/domain"
)

// urlStartRegexp matches the start of an OAuth authorize URL as printed by claude's
// subscription-login fallback (docs/workbench/03_auth_and_login_flow.md, "Task 0 findings").
// Deliberately loose (just "https://" + non-whitespace) rather than matching the exact printed
// wording, since this is TUI output that could change with future `claude` CLI versions.
var urlStartRegexp = regexp.MustCompile(`https://\S+`)

// oauthErrorMarker is the (lowercased) prefix of the line claude prints when a submitted code is
// rejected — see docs/workbench/03_auth_and_login_flow.md, "Task 0 findings".
const oauthErrorMarker = "oauth error"

// loginPromptMarkers are (lowercased) substrings that indicate the login-method menu or the
// "paste code" input is still on screen, even before any URL has been printed. Used to
// distinguish "no prompt yet" (domain.WorkbenchLoginStatePending) from "login completed"
// (domain.WorkbenchLoginStateAuthorized) when neither a URL nor an error line is present.
//
// Known gap: docs/workbench/03_auth_and_login_flow.md confirms a one-time theme-selection
// screen precedes the login-method menu on a brand-new container, but doesn't quote its exact
// wording, so it has no marker here — a poll landing exactly on that screen reads as
// "authorized" rather than "pending" until the caller advances past it (e.g. by relaying an
// Enter keystroke via SubmitLoginCode, per the doc's step 6). Low practical impact: it only
// affects the first poll or two against a brand-new container.
var loginPromptMarkers = []string{
	"claude account with subscription",
	"anthropic console account",
	"3rd-party platform",
	"opening browser to sign in",
	"use the url below to sign in",
	"paste code here",
}

// parseLoginPrompt derives the current subscription_login TUI state from a tmux capture-pane
// snapshot of the workbench's session. Pure and side-effect-free so it's unit-testable directly
// against literal TUI text, without a fake Docker client. See
// docs/workbench/03_auth_and_login_flow.md, "Mechanism (confirmed)", steps 2-3 and 5-6.
func parseLoginPrompt(pane string) domain.WorkbenchLoginPrompt {
	errMsg, hasError := findOAuthError(pane)
	if hasError {
		prompt := domain.WorkbenchLoginPrompt{
			State:        domain.WorkbenchLoginStateError,
			ErrorMessage: errMsg,
		}

		return prompt
	}

	url, hasURL := findLoginURL(pane)
	if hasURL {
		prompt := domain.WorkbenchLoginPrompt{
			State: domain.WorkbenchLoginStateURLPresent,
			URL:   url,
		}

		return prompt
	}

	paneLower := strings.ToLower(pane)
	for _, marker := range loginPromptMarkers {
		if strings.Contains(paneLower, marker) {
			prompt := domain.WorkbenchLoginPrompt{State: domain.WorkbenchLoginStatePending}

			return prompt
		}
	}

	prompt := domain.WorkbenchLoginPrompt{State: domain.WorkbenchLoginStateAuthorized}

	return prompt
}

// findOAuthError returns the trimmed text of the first line starting with oauthErrorMarker
// (case-insensitively), if any.
func findOAuthError(pane string) (string, bool) {
	lines := strings.Split(pane, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(trimmed), oauthErrorMarker) {
			return trimmed, true
		}
	}

	return "", false
}

// findLoginURL locates the OAuth authorize URL in pane, joining wrapped continuation lines back
// together. tmux's capture-pane hard-wraps long lines at the terminal width by inserting a plain
// newline with no separator, so a long URL can be split across two or more lines with no
// whitespace at the join point. A line is only treated as a continuation if it contains no
// whitespace itself and doesn't look like the start of the next prompt (the "Paste code here"
// line) — this is a best-effort heuristic, not an exact terminal-wrap reconstruction.
func findLoginURL(pane string) (string, bool) {
	lines := strings.Split(pane, "\n")

	for i, line := range lines {
		match := urlStartRegexp.FindString(line)
		if match == "" {
			continue
		}

		url := match

		for _, next := range lines[i+1:] {
			trimmed := strings.TrimSpace(next)
			if trimmed == "" {
				break
			}

			if strings.ContainsAny(trimmed, " \t") {
				break
			}

			if strings.HasPrefix(strings.ToLower(trimmed), "paste code here") {
				break
			}

			url += trimmed
		}

		return url, true
	}

	return "", false
}
