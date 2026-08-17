package couchdb

import (
	"reflect"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantFields map[string]string
		wantBody   string
	}{
		{
			name: "simple frontmatter with body",
			content: "---\n" +
				"name: My Skill\n" +
				"description: does a thing\n" +
				"storage_mode: none\n" +
				"---\n" +
				"# Body\n" +
				"instructions here",
			wantFields: map[string]string{
				"name":         "My Skill",
				"description":  "does a thing",
				"storage_mode": "none",
			},
			wantBody: "# Body\ninstructions here",
		},
		{
			name:       "no frontmatter block",
			content:    "just a plain note\nwith no header",
			wantFields: nil,
			wantBody:   "just a plain note\nwith no header",
		},
		{
			name:       "unterminated frontmatter block",
			content:    "---\nname: broken\nno closing delimiter",
			wantFields: nil,
			wantBody:   "---\nname: broken\nno closing delimiter",
		},
		{
			name:       "empty frontmatter block",
			content:    "---\n---\nbody only",
			wantFields: map[string]string{},
			wantBody:   "body only",
		},
		{
			name: "quoted values are unquoted",
			content: "---\n" +
				"name: \"Quoted Name\"\n" +
				"description: 'single quoted'\n" +
				"---\n" +
				"body",
			wantFields: map[string]string{
				"name":        "Quoted Name",
				"description": "single quoted",
			},
			wantBody: "body",
		},
		{
			name: "empty body after closing delimiter",
			content: "---\n" +
				"name: Foo\n" +
				"---",
			wantFields: map[string]string{
				"name": "Foo",
			},
			wantBody: "",
		},
		{
			name: "lines without a colon are skipped",
			content: "---\n" +
				"name: Foo\n" +
				"not a key value line\n" +
				"description: bar\n" +
				"---\n" +
				"body",
			wantFields: map[string]string{
				"name":        "Foo",
				"description": "bar",
			},
			wantBody: "body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, body := ParseFrontmatter(tt.content)

			if !reflect.DeepEqual(fields, tt.wantFields) {
				t.Errorf("fields = %#v, want %#v", fields, tt.wantFields)
			}

			if body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}
