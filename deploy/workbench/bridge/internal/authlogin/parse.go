package authlogin

import (
	"bytes"
	"regexp"
	"strings"
)

// SignalKind identifies which of setup-token's transcript markers a Signal reports.
type SignalKind int

const (
	// SignalAuthLink: the sign-in URL is now known (URL is populated).
	SignalAuthLink SignalKind = iota
	// SignalCodeNeeded: the "paste code here" prompt is now on screen.
	SignalCodeNeeded
	// SignalError: the CLI rejected the most recently submitted code (Error is populated).
	SignalError
	// SignalToken: the flow completed and produced a long-lived token (Token is populated).
	SignalToken
)

// Signal is one event extracted from claude setup-token's raw pty output.
type Signal struct {
	Kind  SignalKind
	URL   string
	Error string
	Token string
}

// osc8LinkPattern extracts the URI parameter of an OSC 8 hyperlink escape sequence:
// ESC ] 8 ; params ; URI BEL. Confirmed against real CLI 2.1.234 output (captured by running
// `claude setup-token` under a pty and inspecting the raw bytes — see testdata/ and parse_test.go):
// the sign-in link is emitted as
//
//	\x1b]8;id=<id>;https://claude.com/cai/oauth/authorize?...\x07
//
// followed by dimmed visible text re-rendering (part of) the same URL for terminals that don't
// support OSC 8, and terminated later by a bare \x1b]8;;\x07. The CLI redraws this block several
// times as the long URL wraps across terminal columns, so the same link appears many times in a
// row in the raw stream — Parser dedupes on "already saw a link" rather than on URL equality, since
// this is the only shape actually observed.
//
// The escape sequence's own URI field is used rather than the wrapped visible text because the
// latter is split across multiple styled segments mid-word (to fit terminal width) and would have
// to be reassembled; the URI field carries it whole and unstyled every time.
var osc8LinkPattern = regexp.MustCompile(`\x1b\]8;[^;]*;(https://[^\x07\x1b]+)\x07`)

// oauthErrorPattern extracts the message of an "OAuth error: ..." line. Confirmed against real
// output: submitting an invalid code produces (raw, escapes elided)
//
//	OAuth error: Invalid code. Please make sure the full code was copied Press Enter to retry.
//
// on one line, terminated by \r (the CLI runs in raw terminal mode, so lines end in \r rather
// than \n).
var oauthErrorPattern = regexp.MustCompile(`OAuth error:([^\r\x1b]+)`)

// tokenPattern matches claude's long-lived OAuth token prefix (the same "sk-ant-" family as
// regular API keys, just from a different subcommand).
//
// Unlike osc8LinkPattern and oauthErrorPattern, this was NOT confirmed against a real success
// screen — completing the flow requires an actual Claude subscription and a real browser login,
// which is out of scope for this change (see the package doc). If setup-token's actual success
// output doesn't contain a bare sk-ant-... token on its own (e.g. it's boxed, wrapped, or printed
// under a different label), this pattern still matches as long as the token substring itself is
// intact in the raw byte stream, which every other tested string in this package's output has
// been. Re-verify this specific pattern against real output before relying on it in production.
var tokenPattern = regexp.MustCompile(`sk-ant-[A-Za-z0-9_-]{20,}`)

// Parser incrementally scans claude setup-token's raw pty output — control sequences included,
// deliberately not stripped, since the OSC 8 escape is what makes link extraction reliable — and
// reports each new signal at most once (errors are the exception: a user can submit several bad
// codes in a row, and each distinct message is reported).
//
// Not safe for concurrent use; one Parser is used per login attempt.
type Parser struct {
	buf []byte

	sawLink       bool
	sawCodeNeeded bool
	sawToken      bool
	lastError     string
}

// NewParser returns a Parser for one setup-token run.
func NewParser() *Parser {
	return &Parser{}
}

// Feed appends chunk to the parser's accumulated transcript and returns any new signals it
// produces.
//
// The transcript is never trimmed: a full setup-token run is at most a few tens of KB, and
// keeping it whole is what lets a marker whose bytes happen to straddle two separate Read() calls
// still be found once the rest arrives on a later Feed.
func (p *Parser) Feed(chunk []byte) []Signal {
	p.buf = append(p.buf, chunk...)

	var signals []Signal

	if !p.sawLink {
		if m := osc8LinkPattern.FindSubmatch(p.buf); m != nil {
			p.sawLink = true

			signals = append(signals, Signal{Kind: SignalAuthLink, URL: string(m[1])})
		}
	}

	if !p.sawCodeNeeded && bytes.Contains(p.buf, []byte("Paste")) && bytes.Contains(p.buf, []byte("prompted")) {
		p.sawCodeNeeded = true

		signals = append(signals, Signal{Kind: SignalCodeNeeded})
	}

	if loc := oauthErrorPattern.FindSubmatchIndex(p.buf); loc != nil {
		// The capturing group is "everything up to \r or an escape byte". If it currently
		// extends all the way to the end of the buffer, that's only because no terminator has
		// arrived yet — the message may still be growing across a later Feed call, and reporting
		// it now would both cut it short and, worse, report a second (longer) signal once the
		// rest arrives. Only a match that ends *before* the buffer's end proves a real
		// terminator was found, and therefore that the text is final.
		matchEnd := loc[1]
		if matchEnd < len(p.buf) {
			groupStart, groupEnd := loc[2], loc[3]

			text := strings.TrimSpace(string(p.buf[groupStart:groupEnd]))
			if text != "" && text != p.lastError {
				p.lastError = text

				signals = append(signals, Signal{Kind: SignalError, Error: text})
			}
		}
	}

	if !p.sawToken {
		if m := tokenPattern.Find(p.buf); m != nil {
			p.sawToken = true

			signals = append(signals, Signal{Kind: SignalToken, Token: string(m)})
		}
	}

	return signals
}
