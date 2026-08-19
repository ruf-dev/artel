package authlogin

import (
	"os"
	"strings"
	"testing"
)

// testdata/setup_token_invalid_code.bin is a real, unmodified capture of `claude setup-token`'s
// raw pty output against CLI 2.1.234: the initial screen (spinner, sign-in URL as an OSC 8
// hyperlink, the "Paste code here if prompted" prompt) followed by an invalid code being
// submitted and the CLI's "OAuth error: ..." response. It contains no real credentials — the
// client_id/state/code_challenge values are per-run OAuth PKCE parameters, not secrets, and no
// login was completed.
const fixturePath = "testdata/setup_token_invalid_code.bin"

func readFixture(t *testing.T) []byte {
	t.Helper()

	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("error reading fixture %s: %v", fixturePath, err)
	}

	return data
}

func TestParser_FeedWholeFixture(t *testing.T) {
	data := readFixture(t)

	parser := NewParser()
	signals := parser.Feed(data)

	var (
		gotLink       *Signal
		gotCodeNeeded bool
		gotError      *Signal
	)

	for i := range signals {
		signal := signals[i]

		switch signal.Kind {
		case SignalAuthLink:
			gotLink = &signal
		case SignalCodeNeeded:
			gotCodeNeeded = true
		case SignalError:
			gotError = &signal
		}
	}

	if gotLink == nil {
		t.Fatal("expected a SignalAuthLink, got none")
	}

	wantURL := "https://claude.com/cai/oauth/authorize?code=true&client_id=9d1c250a-e61b-44d9-88ed-5944d1962f5e&response_type=code&redirect_uri=https%3A%2F%2Fplatform.claude.com%2Foauth%2Fcode%2Fcallback&scope=user%3Ainference&code_challenge=xT29T8OoeGfpfiqkh6qSN_GfGBOj_XqOY8qw8fksDHM&code_challenge_method=S256&state=nBH-NbM1Hwu1tZr_Kl0aS7moVCmB8IhZRzmDfpmVkro"
	if gotLink.URL != wantURL {
		t.Errorf("auth link URL mismatch:\n got:  %s\n want: %s", gotLink.URL, wantURL)
	}

	if !gotCodeNeeded {
		t.Error("expected a SignalCodeNeeded, got none")
	}

	if gotError == nil {
		t.Fatal("expected a SignalError, got none")
	}

	wantErrSubstring := "Invalid code"
	if !strings.Contains(gotError.Error, wantErrSubstring) {
		t.Errorf("error message %q does not contain %q", gotError.Error, wantErrSubstring)
	}
}

// TestParser_FeedByteByByte replays the same fixture one byte at a time, to prove Feed's signals
// don't depend on control sequences or marker words arriving intact in a single chunk — real pty
// reads split output at arbitrary byte boundaries.
func TestParser_FeedByteByByte(t *testing.T) {
	data := readFixture(t)

	parser := NewParser()

	var (
		linkCount       int
		codeNeededCount int
		errorCount      int
	)

	for i := range data {
		signals := parser.Feed(data[i : i+1])

		for _, signal := range signals {
			switch signal.Kind {
			case SignalAuthLink:
				linkCount++
			case SignalCodeNeeded:
				codeNeededCount++
			case SignalError:
				errorCount++
			}
		}
	}

	if linkCount != 1 {
		t.Errorf("expected exactly one SignalAuthLink, got %d", linkCount)
	}

	if codeNeededCount != 1 {
		t.Errorf("expected exactly one SignalCodeNeeded, got %d", codeNeededCount)
	}

	if errorCount != 1 {
		t.Errorf("expected exactly one SignalError, got %d", errorCount)
	}
}

func TestParser_NoSignalsOnEmptyInput(t *testing.T) {
	parser := NewParser()

	signals := parser.Feed(nil)
	if len(signals) != 0 {
		t.Errorf("expected no signals from empty input, got %v", signals)
	}
}

// TestParser_AuthLinkOSC8TerminatedByST confirms authLinkPattern matches an OSC 8 hyperlink
// terminated by ST (\x1b\\) rather than the BEL (\x07) terminator seen in the captured fixture —
// the old osc8LinkPattern hardcoded \x07 and could not have matched this.
func TestParser_AuthLinkOSC8TerminatedByST(t *testing.T) {
	parser := NewParser()

	chunk := []byte("\x1b]8;id=1;https://claude.com/cai/oauth/authorize?code=true&client_id=abc&state=xyz\x1b\\")
	signals := parser.Feed(chunk)

	var got *Signal

	for i := range signals {
		if signals[i].Kind == SignalAuthLink {
			got = &signals[i]
		}
	}

	if got == nil {
		t.Fatal("expected a SignalAuthLink, got none")
	}

	want := "https://claude.com/cai/oauth/authorize?code=true&client_id=abc&state=xyz"
	if got.URL != want {
		t.Errorf("auth link URL mismatch:\n got:  %s\n want: %s", got.URL, want)
	}
}

// TestParser_AuthLinkPlainText confirms authLinkPattern matches the sign-in URL even with no OSC
// 8 wrapper at all (e.g. a hypothetical non-tty-detection fallback that prints plain text).
func TestParser_AuthLinkPlainText(t *testing.T) {
	parser := NewParser()

	chunk := []byte("Visit: https://claude.com/cai/oauth/authorize?code=true&client_id=abc&state=xyz\r\n")
	signals := parser.Feed(chunk)

	var got *Signal

	for i := range signals {
		if signals[i].Kind == SignalAuthLink {
			got = &signals[i]
		}
	}

	if got == nil {
		t.Fatal("expected a SignalAuthLink, got none")
	}

	want := "https://claude.com/cai/oauth/authorize?code=true&client_id=abc&state=xyz"
	if got.URL != want {
		t.Errorf("auth link URL mismatch:\n got:  %s\n want: %s", got.URL, want)
	}
}

// TestParser_AuthLinkSplitAcrossReads reproduces a pty Read splitting the OSC 8 sequence mid-URL:
// the first Feed call's chunk ends inside the URL, with no terminator byte yet available. Before
// the fix, authLinkPattern.Find would match everything up to that arbitrary cutoff and latch
// sawLink on the truncated URL, permanently dropping trailing params like redirect_uri.
func TestParser_AuthLinkSplitAcrossReads(t *testing.T) {
	parser := NewParser()

	firstChunk := []byte("\x1b]8;id=1;https://claude.com/cai/oauth/authorize?code=true&client_id=abc&redirect")
	signals := parser.Feed(firstChunk)

	for _, signal := range signals {
		if signal.Kind == SignalAuthLink {
			t.Fatalf("did not expect a SignalAuthLink before the terminator arrives, got URL %q", signal.URL)
		}
	}

	secondChunk := []byte("_uri=https%3A%2F%2Fexample.com%2Fcallback&state=xyz\x07")
	signals = parser.Feed(secondChunk)

	var got *Signal

	for i := range signals {
		if signals[i].Kind == SignalAuthLink {
			got = &signals[i]
		}
	}

	if got == nil {
		t.Fatal("expected a SignalAuthLink once the terminator arrives, got none")
	}

	want := "https://claude.com/cai/oauth/authorize?code=true&client_id=abc&redirect_uri=https%3A%2F%2Fexample.com%2Fcallback&state=xyz"
	if got.URL != want {
		t.Errorf("auth link URL mismatch:\n got:  %s\n want: %s", got.URL, want)
	}
}

func TestParser_TokenSignal(t *testing.T) {
	parser := NewParser()

	chunk := []byte("Your long-lived token: sk-ant-oat01-abcdefghijklmnopqrstuvwxyz0123456789\r\n")
	signals := parser.Feed(chunk)

	var got *Signal

	for i := range signals {
		if signals[i].Kind == SignalToken {
			got = &signals[i]
		}
	}

	if got == nil {
		t.Fatal("expected a SignalToken, got none")
	}

	want := "sk-ant-oat01-abcdefghijklmnopqrstuvwxyz0123456789"
	if got.Token != want {
		t.Errorf("token mismatch:\n got:  %s\n want: %s", got.Token, want)
	}
}
