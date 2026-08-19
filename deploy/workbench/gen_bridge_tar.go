//go:build ignore

// Command gen_bridge_tar packages deploy/workbench/bridge — a separate Go module, see its own
// go.mod and README.md — into deploy/workbench/bridge.tar, a single ordinary file that
// embed.go's //go:embed directive can name directly.
//
// This indirection exists because Go's //go:embed refuses to embed anything inside a directory
// that contains its own go.mod. Confirmed empirically (go1.26.3) against a minimal repro: the
// exact error is
//
//	pattern bridge: cannot embed directory bridge: in different module
//
// and it fires identically whether the pattern names the directory itself, a subdirectory that
// doesn't contain go.mod, the "all:" prefix, or even go.mod as a standalone file pattern — there
// is no path-selection trick around it. bridge/ has to keep its own go.mod (see its README: it
// must not import github.com/ruf-dev/artel, and a nested go.mod is what makes that a compile
// error instead of a convention someone can accidentally violate), so this repackages the module
// tree as one opaque file instead.
//
// Regenerate after any change under bridge/ — including a go.sum update from `go mod tidy` run
// there — with:
//
//	go generate ./deploy/workbench/...
//
// internal/clients/workbenchdocker/image.go unpacks bridge.tar back into the Docker build
// context it sends to the daemon and folds its bytes into the content-addressed image tag exactly
// like Dockerfile/entrypoint.sh; a stale bridge.tar is the only way a bridge source change could
// fail to produce a new image tag, since nothing else notices the two have diverged.
package main

import (
	"archive/tar"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// root is relative to this file's directory (deploy/workbench), matching how `go generate` and
// `go run` both invoke it with that as the working directory.
const root = "bridge"

const outputName = "bridge.tar"

func main() {
	var paths []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		paths = append(paths, path)

		return nil
	})
	if err != nil {
		log.Fatalf("gen_bridge_tar: error walking %s: %v", root, err)
	}

	// Sorted so the resulting tar (and therefore its hash, which feeds workbenchImageTag) doesn't
	// depend on the filesystem's directory-listing order.
	sort.Strings(paths)

	out, err := os.Create(outputName)
	if err != nil {
		log.Fatalf("gen_bridge_tar: error creating %s: %v", outputName, err)
	}
	defer out.Close()

	err = writeTar(out, paths)
	if err != nil {
		log.Fatalf("gen_bridge_tar: %v", err)
	}
}

func writeTar(out *os.File, paths []string) error {
	tarWriter := tar.NewWriter(out)

	for _, path := range paths {
		err := addFile(tarWriter, path)
		if err != nil {
			return err
		}
	}

	return tarWriter.Close()
}

// addFile writes one file into tarWriter under its slash-form path (already "bridge/..." since
// root is walked from deploy/workbench), with a fixed mode and no mtime — nothing under bridge/
// needs to be executable, and a fixed mode keeps the archive's bytes a pure function of file
// content, matching workbenchImageTag's determinism goal on the consuming side.
func addFile(tarWriter *tar.Writer, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	header := tar.Header{
		Name: filepath.ToSlash(path),
		Mode: 0o644,
		Size: int64(len(content)),
	}

	err = tarWriter.WriteHeader(&header)
	if err != nil {
		return err
	}

	_, err = tarWriter.Write(content)

	return err
}
