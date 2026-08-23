// Package terminalauth scans a workbench's raw ttyd terminal output for the Claude CLI's OAuth
// sign-in link, the same one deploy/workbench/bridge/internal/authlogin.Parser looks for on the
// chat bridge's transport — but the interactive terminal (internal/transport/vaults_api's
// WorkbenchTerminalShellHandler, proxying ttyd's own binary WebSocket protocol) is a separate
// transport with no visibility into the bridge's parser, so it needs its own tap.
//
// Unlike authlogin.Parser, which is used once per short-lived `claude setup-token` run, a Parser
// here lives for the duration of an interactive terminal WebSocket connection — potentially hours
// — so its buffer is capped (see maxBufferSize) rather than kept whole, and it can report more
// than one link over its lifetime (a user may run `claude setup-token` again later in the same
// terminal session).
package terminalauth

import "regexp"

// maxBufferSize caps the rolling buffer Feed accumulates across calls. A terminal WS connection
// is long-lived, so the buffer is trimmed from the front (keeping the tail) once it exceeds this
// size, so memory doesn't grow unbounded over a multi-hour session. The sign-in link itself is
// only ~150-250 bytes, well under this cap, so trimming never truncates an in-progress match — see
// Feed's doc comment.
const maxBufferSize = 16 * 1024

// authLinkPattern mirrors deploy/workbench/bridge/internal/authlogin.authLinkPattern exactly: it
// extracts the sign-in URL by keying off its known stable prefix and stopping at the first control
// character, rather than requiring the exact OSC 8 escape framing the CLI happens to use — this
// tolerates the link appearing either OSC-8-hyperlink-wrapped or as plain unstyled text. See that
// file's doc comment for the full rationale, confirmed against real CLI output.
var authLinkPattern = regexp.MustCompile(`https://claude\.com/cai/oauth/authorize\?[^\x00-\x1f\x7f]+`)

// Parser incrementally scans a terminal WS connection's raw ttyd OUTPUT payload bytes (ANSI/OSC
// escapes included, not stripped — the pattern only needs its own character run, not the
// surrounding escape structure) and reports the Claude CLI OAuth sign-in link each time a new one
// appears complete in the stream.
//
// One Parser is used per active terminal WS connection; not safe for concurrent use.
type Parser struct {
	buf     []byte
	lastURL string
}

// NewParser returns a Parser for one terminal WS connection.
func NewParser() *Parser {
	return &Parser{}
}

// Feed appends chunk (an OUTPUT frame's payload, with the leading command byte already stripped by
// the caller) to the parser's rolling buffer and reports (url, true) the first time a complete,
// new sign-in link is found — "new" meaning different from the last link this Parser reported, so
// the CLI redrawing the same link several times as it repaints the terminal doesn't produce
// repeated signals. Returns ("", false) when no new complete link is present yet.
func (p *Parser) Feed(chunk []byte) (string, bool) {
	p.buf = append(p.buf, chunk...)

	url, found := p.tryMatch()

	if len(p.buf) > maxBufferSize {
		p.buf = p.buf[len(p.buf)-maxBufferSize:]
	}

	return url, found
}

// tryMatch looks for a complete link in the current buffer. A match whose end coincides with the
// buffer's current end is only proof the URL run hasn't been cut off by a real terminator yet — it
// may just be where this Feed call's chunk happened to stop, e.g. a ttyd OUTPUT frame splitting the
// OSC 8 sequence mid-URL (the link is ~150-250 bytes, well within a single frame but not
// guaranteed to arrive in one). Only a match ending strictly before the buffer's end proves a real
// terminator was found downstream, and therefore that the URL is complete.
func (p *Parser) tryMatch() (string, bool) {
	loc := authLinkPattern.FindIndex(p.buf)
	if loc == nil || loc[1] >= len(p.buf) {
		return "", false
	}

	url := string(p.buf[loc[0]:loc[1]])

	// Drop everything through the end of this match, whether or not it turns out to be newly
	// reported below. Without this, a stale already-reported match sitting earlier in the buffer
	// (regexp.FindIndex always returns the leftmost match) would keep hiding a genuinely new,
	// later link from being found on a subsequent Feed call, since the leftmost match is checked
	// first every time.
	p.buf = p.buf[loc[1]:]

	if url == p.lastURL {
		return "", false
	}

	p.lastURL = url

	return url, true
}
