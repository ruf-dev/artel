package workbenchdocker

import (
	"reflect"
	"testing"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
)

// TestParseTmuxWindowList exercises parseTmuxWindowList against tmux list-windows output shaped
// by tmuxWindowListFormat ("#{window_id}\t#{window_active}\t#{window_name}", one line per
// window) — see tmux_tabs.go's doc comment on why tab is the field separator.
func TestParseTmuxWindowList(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   []domain.TerminalTab
	}{
		{
			name:   "well-formed multi-line output",
			output: "@1\t1\tclaude\n@2\t0\tbash\n@3\t0\tvim\n",
			want: []domain.TerminalTab{
				{ID: "@1", Name: "claude", Active: true},
				{ID: "@2", Name: "bash", Active: false},
				{ID: "@3", Name: "vim", Active: false},
			},
		},
		{
			name:   "single window",
			output: "@1\t1\tclaude\n",
			want: []domain.TerminalTab{
				{ID: "@1", Name: "claude", Active: true},
			},
		},
		{
			name:   "trailing newline does not produce a phantom entry",
			output: "@1\t1\tclaude\n@2\t0\tbash\n\n",
			want: []domain.TerminalTab{
				{ID: "@1", Name: "claude", Active: true},
				{ID: "@2", Name: "bash", Active: false},
			},
		},
		{
			name:   "window name containing spaces parses correctly since only tab delimits",
			output: "@1\t1\tclaude working on tmux_tabs.go\n",
			want: []domain.TerminalTab{
				{ID: "@1", Name: "claude working on tmux_tabs.go", Active: true},
			},
		},
		{
			name:   "window_active 1 maps to true",
			output: "@1\t1\tclaude\n",
			want: []domain.TerminalTab{
				{ID: "@1", Name: "claude", Active: true},
			},
		},
		{
			name:   "window_active 0 maps to false",
			output: "@1\t0\tclaude\n",
			want: []domain.TerminalTab{
				{ID: "@1", Name: "claude", Active: false},
			},
		},
		{
			name:   "empty output produces no tabs",
			output: "",
			want:   []domain.TerminalTab{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTmuxWindowList(tt.output)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseTmuxWindowList(%q) = %#v, want %#v", tt.output, got, tt.want)
			}
		})
	}
}

// TestIsSessionMissingErr exercises isSessionMissingErr against nil, generic, and tmux's
// session/server-gone errors — including KillTmuxWindow's distinct "can't find window" wording,
// to prove there's no false-positive collision between the two substrings (see tmux_tabs.go's
// doc comment on isSessionMissingErr).
func TestIsSessionMissingErr(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "generic error not matching either substring",
			err:  rerrors.New("tmux exited with status 1: can't find window: @9"),
			want: false,
		},
		{
			name: "can't find session",
			err:  rerrors.New("tmux exited with status 1: can't find session: workbench"),
			want: true,
		},
		{
			name: "no server running",
			err:  rerrors.New("tmux exited with status 1: no server running on /tmp/tmux-0/default"),
			want: true,
		},
		{
			name: "can't find session wrapped a second time",
			err:  rerrors.Wrap(rerrors.New("tmux exited with status 1: can't find session: workbench"), "listing tmux windows"),
			want: true,
		},
		{
			name: "empty message string",
			err:  rerrors.New(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSessionMissingErr(tt.err)

			if got != tt.want {
				t.Fatalf("isSessionMissingErr(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
