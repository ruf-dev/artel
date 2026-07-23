package workbench

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
)

func TestParseLoginPrompt_UrlPresent(t *testing.T) {
	pane := "Opening browser to sign in...\n" +
		"Browser didn't open? Use the url below to sign in (c to copy)\n" +
		"\n" +
		"https://platform.claude.com/oauth/authorize?code=true&state=abc123\n" +
		"\n" +
		"Paste code here if prompted >"

	got := parseLoginPrompt(pane)
	require.Equal(t, domain.WorkbenchLoginStateURLPresent, got.State)
	require.Equal(t, "https://platform.claude.com/oauth/authorize?code=true&state=abc123", got.URL)
	require.Empty(t, got.ErrorMessage)
}

// TestParseLoginPrompt_UrlWrappedAcrossLines exercises the case tmux's capture-pane hard-wraps a
// long URL at the terminal width — no whitespace at the wrap point, so continuation lines must
// be joined back with no separator.
func TestParseLoginPrompt_UrlWrappedAcrossLines(t *testing.T) {
	pane := "Use the url below to sign in (c to copy)\n" +
		"https://platform.claude.com/oauth/authorize?client_id=abc\n" +
		"&redirect_uri=https%3A%2F%2Fplatform.claude.com%2Foauth%2Fcode%2Fcallback\n" +
		"\n" +
		"Paste code here if prompted >"

	got := parseLoginPrompt(pane)
	require.Equal(t, domain.WorkbenchLoginStateURLPresent, got.State)
	want := "https://platform.claude.com/oauth/authorize?client_id=abc" +
		"&redirect_uri=https%3A%2F%2Fplatform.claude.com%2Foauth%2Fcode%2Fcallback"
	require.Equal(t, want, got.URL)
}

func TestParseLoginPrompt_OAuthErrorPresent(t *testing.T) {
	pane := "https://platform.claude.com/oauth/authorize?code=true\n" +
		"\n" +
		"OAuth error: Invalid code. Please make sure the full code was copied\n" +
		"Paste code here if prompted >"

	got := parseLoginPrompt(pane)
	require.Equal(t, domain.WorkbenchLoginStateError, got.State)
	require.Equal(t, "OAuth error: Invalid code. Please make sure the full code was copied", got.ErrorMessage)
	require.Empty(t, got.URL)
}

// TestParseLoginPrompt_ErrorTakesPriorityOverStaleUrl covers the realistic case: the URL was
// printed earlier and is still visible in the pane's scrollback above a subsequent rejected-code
// error — the error should win, signaling "retry", not "here's a (stale) URL".
func TestParseLoginPrompt_ErrorTakesPriorityOverStaleUrl(t *testing.T) {
	pane := "https://platform.claude.com/oauth/authorize?code=true\n" +
		"OAuth error: Invalid code. Please make sure the full code was copied\n" +
		"Paste code here if prompted >"

	got := parseLoginPrompt(pane)
	require.Equal(t, domain.WorkbenchLoginStateError, got.State)
}

func TestParseLoginPrompt_NeitherPresent_NoMenuMarkers_Authorized(t *testing.T) {
	pane := "│ > hello, what can you help me with today?\n"

	got := parseLoginPrompt(pane)
	require.Equal(t, domain.WorkbenchLoginStateAuthorized, got.State)
}

func TestParseLoginPrompt_NeitherPresent_LoginMenuStillShowing_Pending(t *testing.T) {
	pane := "Select login method:\n" +
		"> Claude account with subscription\n" +
		"  Anthropic Console account\n" +
		"  3rd-party platform"

	got := parseLoginPrompt(pane)
	require.Equal(t, domain.WorkbenchLoginStatePending, got.State)
}

// TestParseLoginPrompt_ThemeScreenStillShowing_KnownGap documents a known, accepted gap (see the
// comment on loginPromptMarkers): the one-time theme-selection screen's exact wording isn't
// quoted in docs/workbench/03_auth_and_login_flow.md, so parseLoginPrompt has no marker for it
// and reports "authorized" rather than "pending" while it's showing. Low practical impact — it
// only affects the first poll or two against a brand-new container.
func TestParseLoginPrompt_ThemeScreenStillShowing_KnownGap(t *testing.T) {
	pane := "Choose the text style that looks best with your terminal:\n" +
		"> Dark mode\n" +
		"  Light mode\n"

	got := parseLoginPrompt(pane)
	require.Equal(t, domain.WorkbenchLoginStateAuthorized, got.State)
}

func TestParseLoginPrompt_EmptyPane_Authorized(t *testing.T) {
	got := parseLoginPrompt("")
	require.Equal(t, domain.WorkbenchLoginStateAuthorized, got.State)
}
