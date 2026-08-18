package workbenchdocker

import (
	"reflect"
	"testing"

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
