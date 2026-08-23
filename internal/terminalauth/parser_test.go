package terminalauth

import "testing"

const sampleLink = "https://claude.com/cai/oauth/authorize?client_id=abc&redirect_uri=https://example.com"

// TestParser_Feed covers single-call detection, a no-link chunk, and a link split across two
// Feed calls that must not be detected until the second call supplies the terminator.
func TestParser_Feed(t *testing.T) {
	tests := []struct {
		name      string
		feeds     []string
		wantURL   string
		wantFound bool
	}{
		{
			name:      "no link present",
			feeds:     []string{"just some regular pty output, nothing to see here\r\n"},
			wantURL:   "",
			wantFound: false,
		},
		{
			name:      "link complete in a single chunk, terminated by a control byte",
			feeds:     []string{"prefix text \x1b]8;;" + sampleLink + "\x07 trailer"},
			wantURL:   sampleLink,
			wantFound: true,
		},
		{
			name: "link split across two Feed calls is not detected until the terminator arrives",
			feeds: []string{
				"prefix text \x1b]8;;https://claude.com/cai/oauth/authorize?client_id=abc&redirect_uri",
				"=https://example.com\x07 trailer",
			},
			wantURL:   sampleLink,
			wantFound: true,
		},
		{
			name:      "link with no terminator yet is never reported",
			feeds:     []string{"prefix text \x1b]8;;" + sampleLink},
			wantURL:   "",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()

			var gotURL string
			var gotFound bool

			for _, chunk := range tt.feeds {
				url, found := p.Feed([]byte(chunk))
				if found {
					gotURL = url
					gotFound = true
				}
			}

			if gotFound != tt.wantFound {
				t.Fatalf("found = %v, want %v", gotFound, tt.wantFound)
			}

			if gotURL != tt.wantURL {
				t.Fatalf("url = %q, want %q", gotURL, tt.wantURL)
			}
		})
	}
}

// TestParser_Feed_SplitAcrossCallsNotDetectedEarly is an explicit two-step check (rather than
// folding the intermediate assertion into the table above) that the first Feed call of a split
// link reports nothing at all, only the second one does.
func TestParser_Feed_SplitAcrossCallsNotDetectedEarly(t *testing.T) {
	p := NewParser()

	url, found := p.Feed([]byte("prefix \x1b]8;;https://claude.com/cai/oauth/authorize?client_id=abc&redirect_uri"))
	if found {
		t.Fatalf("first (truncated) chunk reported found=true, url=%q; want not found yet", url)
	}

	url, found = p.Feed([]byte("=https://example.com\x07 trailer"))
	if !found {
		t.Fatal("second chunk (completing the link) did not report found=true")
	}

	if url != sampleLink {
		t.Fatalf("url = %q, want %q", url, sampleLink)
	}
}

// TestParser_Feed_DedupesRepeatedLink confirms the CLI redrawing the same link several times (as
// it repaints the terminal) only produces one signal.
func TestParser_Feed_DedupesRepeatedLink(t *testing.T) {
	p := NewParser()

	chunk := []byte("\x1b]8;;" + sampleLink + "\x07")

	_, found := p.Feed(chunk)
	if !found {
		t.Fatal("first occurrence: want found=true")
	}

	_, found = p.Feed(chunk)
	if found {
		t.Fatal("repeated identical link: want found=false")
	}
}

// TestParser_Feed_ReportsANewLinkAfterAPriorOne confirms a Parser is not a one-shot latch: a
// terminal WS connection is long-lived and a user may run `claude setup-token` again later in the
// same session, producing a distinct link that must still be detected.
func TestParser_Feed_ReportsANewLinkAfterAPriorOne(t *testing.T) {
	p := NewParser()

	_, found := p.Feed([]byte("\x1b]8;;" + sampleLink + "\x07"))
	if !found {
		t.Fatal("first link: want found=true")
	}

	secondLink := "https://claude.com/cai/oauth/authorize?client_id=xyz&redirect_uri=https://example.com/other"

	url, found := p.Feed([]byte("\x1b]8;;" + secondLink + "\x07"))
	if !found {
		t.Fatal("second, distinct link: want found=true")
	}

	if url != secondLink {
		t.Fatalf("url = %q, want %q", url, secondLink)
	}
}

// TestParser_Feed_BufferTrimDoesNotBreakASlowlyArrivingLink confirms a large amount of unrelated
// pty output preceding the link does not prevent detection once the buffer is capped/trimmed —
// the link itself is always at the tail of the buffer relative to the unrelated bytes ahead of it.
func TestParser_Feed_BufferTrimDoesNotBreakASlowlyArrivingLink(t *testing.T) {
	p := NewParser()

	noise := make([]byte, maxBufferSize*3)
	for i := range noise {
		noise[i] = 'x'
	}

	_, found := p.Feed(noise)
	if found {
		t.Fatal("pure noise: want found=false")
	}

	url, found := p.Feed([]byte("\x1b]8;;" + sampleLink + "\x07"))
	if !found {
		t.Fatal("link fed after a large noise buffer: want found=true")
	}

	if url != sampleLink {
		t.Fatalf("url = %q, want %q", url, sampleLink)
	}
}
