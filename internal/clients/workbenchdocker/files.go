package workbenchdocker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"

	"github.com/ruf-dev/artel/internal/utils"

	"go.redsock.ru/rerrors"
)

// WriteFilesToVolume writes files (keys are paths relative to the workspace root, e.g.
// "notes/foo.md") into containerID's workspace volume by building an in-memory tar archive and
// uploading it via the Docker Engine API's archive endpoint. Works whether the container is
// running or stopped — the copy API operates on the container's filesystem directly, not on a
// live process inside it — so this is used both to materialize a vault into a freshly-created,
// not-yet-started container (StartWorkbench) and, in principle, against an already-running one.
func (c *Client) WriteFilesToVolume(ctx context.Context, containerID string, files map[string][]byte) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	for path, content := range files {
		header := &tar.Header{
			Name: path,
			Mode: 0o644,
			Size: int64(len(content)),
		}

		err := tw.WriteHeader(header)
		if err != nil {
			return rerrors.Wrap(err, "error writing tar header for "+path)
		}

		_, err = tw.Write(content)
		if err != nil {
			return rerrors.Wrap(err, "error writing tar content for "+path)
		}
	}

	err := tw.Close()
	if err != nil {
		return rerrors.Wrap(err, "error closing tar writer")
	}

	copyOptions := container.CopyToContainerOptions{}

	err = c.cli.CopyToContainer(ctx, containerID, homeMountPath, buf, copyOptions)
	if err != nil {
		return rerrors.Wrap(err, "error copying files to workbench volume")
	}

	return nil
}

// ReadFilesFromVolume reads every regular file under containerID's workspace volume and returns
// it as a map keyed by path relative to the workspace root — the mirror image of
// WriteFilesToVolume, used to pull the workbench's edits back out on stop. Like
// WriteFilesToVolume, this works whether the container is running or stopped.
//
// The Docker Engine API's archive endpoint always wraps the requested path in one leading
// directory component named after the last path segment of homeMountPath (i.e. "vault/
// notes/foo.md" rather than "notes/foo.md") — that leading component is stripped off here so
// callers see paths relative to the workspace root, matching WriteFilesToVolume's own path
// shape.
func (c *Client) ReadFilesFromVolume(ctx context.Context, containerID string) (map[string][]byte, error) {
	reader, _, err := c.cli.CopyFromContainer(ctx, containerID, homeMountPath)
	if err != nil {
		return nil, rerrors.Wrap(err, "error copying files from workbench volume")
	}
	defer utils.CloseWithLog(reader, "workbench volume archive reader")

	files := make(map[string][]byte)

	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, rerrors.Wrap(err, "error reading workbench volume archive")
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		relPath := stripLeadingDir(header.Name)
		if relPath == "" {
			continue
		}

		content, err := io.ReadAll(tr)
		if err != nil {
			return nil, rerrors.Wrap(err, "error reading tar content for "+header.Name)
		}

		files[relPath] = content
	}

	return files, nil
}

// stripLeadingDir drops the first "/"-separated path component of name, returning "" if name
// has no separator at all (i.e. it IS the leading component — the workspace root directory entry
// itself).
func stripLeadingDir(name string) string {
	name = strings.TrimPrefix(name, "/")

	idx := strings.Index(name, "/")
	if idx < 0 {
		return ""
	}

	return name[idx+1:]
}
