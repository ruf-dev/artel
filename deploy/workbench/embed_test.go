package workbenchimage

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"testing"
)

// TestFiles_EmbedsFlatBuildContext guards the //go:embed directive's shape: exactly Dockerfile,
// entrypoint.sh, tmux.conf, and bridge.tar must be embedded, each non-empty.
func TestFiles_EmbedsFlatBuildContext(t *testing.T) {
	wantFiles := []string{"Dockerfile", "entrypoint.sh", "tmux.conf", "bridge.tar"}

	for _, path := range wantFiles {
		content, err := Files.ReadFile(path)
		if err != nil {
			t.Errorf("expected embedded file %q, got error: %v", path, err)

			continue
		}

		if len(content) == 0 {
			t.Errorf("embedded file %q is empty", path)
		}
	}
}

// TestFiles_BridgeTarContainsModuleTree guards gen_bridge_tar.go's output shape: bridge.tar must
// unpack to the bridge/ module tree — including a nested file (go.sum) and a file one level
// deeper still (a package under internal/) — with paths already prefixed "bridge/...", since
// internal/clients/workbenchdocker/image.go re-emits its entries into the Docker build context
// unchanged.
func TestFiles_BridgeTarContainsModuleTree(t *testing.T) {
	content, err := Files.ReadFile("bridge.tar")
	if err != nil {
		t.Fatalf("error reading embedded bridge.tar: %v", err)
	}

	wantEntries := map[string]bool{
		"bridge/go.mod":                          false,
		"bridge/go.sum":                          false,
		"bridge/main.go":                         false,
		"bridge/internal/chatprotocol/events.go": false,
	}

	tarReader := tar.NewReader(bytes.NewReader(content))

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			t.Fatalf("error reading bridge.tar entry: %v", err)
		}

		if _, ok := wantEntries[header.Name]; ok {
			wantEntries[header.Name] = true
		}
	}

	for name, found := range wantEntries {
		if !found {
			t.Errorf("expected bridge.tar to contain %q, not found", name)
		}
	}
}
