package notes

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestNormalizeFolderPath(t *testing.T) {
	cases := map[string]string{
		"":            "",
		"/":           "",
		"Work":        "Work",
		"/Work/":      "Work",
		"Work/Notes":  "Work/Notes",
		"/Work/Notes": "Work/Notes",
	}

	for input, want := range cases {
		got := normalizeFolderPath(input)
		if got != want {
			t.Errorf("normalizeFolderPath(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestUnderFolder(t *testing.T) {
	cases := []struct {
		notePath   string
		folderPath string
		want       bool
	}{
		{"Work/Notes/todo.md", "", true},
		{"Work/Notes/todo.md", "Work/Notes", true},
		{"Work/Notes", "Work/Notes", true},
		{"Work/Notes2/todo.md", "Work/Notes", false},
		{"Personal/todo.md", "Work/Notes", false},
		{"Work/Notes/sub/todo.md", "Work/Notes", true},
	}

	for _, c := range cases {
		got := underFolder(c.notePath, c.folderPath)
		if got != c.want {
			t.Errorf("underFolder(%q, %q) = %v, want %v", c.notePath, c.folderPath, got, c.want)
		}
	}
}

func TestRelativeToFolder(t *testing.T) {
	cases := []struct {
		notePath   string
		folderPath string
		want       string
	}{
		{"Work/Notes/todo.md", "", "Work/Notes/todo.md"},
		{"Work/Notes/todo.md", "Work/Notes", "todo.md"},
		{"Work/Notes/sub/todo.md", "Work/Notes", "sub/todo.md"},
	}

	for _, c := range cases {
		got := relativeToFolder(c.notePath, c.folderPath)
		if got != c.want {
			t.Errorf("relativeToFolder(%q, %q) = %q, want %q", c.notePath, c.folderPath, got, c.want)
		}
	}
}

func TestJoinDestPath(t *testing.T) {
	cases := []struct {
		destPath  string
		entryName string
		want      string
	}{
		{"", "todo.md", "todo.md"},
		{"", "/todo.md", "todo.md"},
		{"Imported", "todo.md", "Imported/todo.md"},
		{"Imported", "sub/todo.md", "Imported/sub/todo.md"},
	}

	for _, c := range cases {
		got := joinDestPath(c.destPath, c.entryName)
		if got != c.want {
			t.Errorf("joinDestPath(%q, %q) = %q, want %q", c.destPath, c.entryName, got, c.want)
		}
	}
}

func buildTestZip(t *testing.T, files map[string]string) []byte {
	t.Helper()

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("error creating zip entry %q: %v", name, err)
		}

		_, err = w.Write([]byte(content))
		if err != nil {
			t.Fatalf("error writing zip entry %q: %v", name, err)
		}
	}

	err := zw.Close()
	if err != nil {
		t.Fatalf("error closing zip writer: %v", err)
	}

	return buf.Bytes()
}

func TestReadZipEntries(t *testing.T) {
	files := map[string]string{
		"todo.md":     "# Todo",
		"sub/note.md": "# Sub note",
	}

	zipData := buildTestZip(t, files)

	entries, err := readZipEntries(zipData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != len(files) {
		t.Fatalf("expected %d entries, got %d", len(files), len(entries))
	}

	got := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.Name] = string(e.Content)
	}

	for name, content := range files {
		if got[name] != content {
			t.Errorf("entry %q = %q, want %q", name, got[name], content)
		}
	}
}
