package workbenchdocker

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"sort"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/pkg/jsonmessage"

	"go.redsock.ru/rerrors"

	workbenchimage "github.com/ruf-dev/artel/deploy/workbench"
	"github.com/ruf-dev/artel/internal/utils"
)

// workbenchImageTagPrefix is the repository name every workbench image is tagged under — see
// workbenchImageTag for how the tag suffix is derived.
const workbenchImageTagPrefix = "artel-workbench:"

// workbenchImageTagHashLen is the number of hex characters of the content hash kept in the
// image tag — enough to make an accidental collision practically impossible while keeping tags
// short and readable.
const workbenchImageTagHashLen = 12

// EnsureImage returns the tag of the workbench Docker image on the daemon this Client talks to,
// building it from the embedded deploy/workbench build context (workbenchimage.Files) if an
// image with that tag isn't already present. The tag is a deterministic hash of the embedded
// files' content (see workbenchImageTag), so a changed Dockerfile/entrypoint.sh automatically
// produces a new tag and triggers a rebuild on next use, while an unchanged build context reuses
// whatever was already built — on this daemon or a previous run against it.
func (c *Client) EnsureImage(ctx context.Context) (string, error) {
	tag, err := workbenchImageTag()
	if err != nil {
		return "", rerrors.Wrap(err, "error computing workbench image tag")
	}

	_, err = c.cli.ImageInspect(ctx, tag)
	if err == nil {
		return tag, nil
	}

	if !errdefs.IsNotFound(err) {
		return "", rerrors.Wrap(err, "error inspecting workbench image")
	}

	err = c.buildImage(ctx, tag)
	if err != nil {
		return "", rerrors.Wrap(err, "error building workbench image")
	}

	return tag, nil
}

// buildImage sends the embedded deploy/workbench build context to the daemon and tags the
// resulting image as tag.
func (c *Client) buildImage(ctx context.Context, tag string) error {
	buildContext, err := workbenchBuildContext()
	if err != nil {
		return rerrors.Wrap(err, "error assembling workbench image build context")
	}

	buildOptions := build.ImageBuildOptions{
		Tags:   []string{tag},
		Remove: true,
	}

	resp, err := c.cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return rerrors.Wrap(err, "error sending workbench image build request")
	}
	defer utils.CloseWithLog(resp.Body, "workbench image build response body")

	// A non-nil error from ImageBuild above only means the HTTP request itself failed. The
	// daemon streams build progress/errors back as a sequence of JSON messages in resp.Body, so
	// a mid-build failure (e.g. a RUN step exiting non-zero) only ever surfaces here.
	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, io.Discard, 0, false, nil)
	if err != nil {
		return rerrors.Wrap(err, "error in workbench image build stream")
	}

	return nil
}

// workbenchImageTag computes a deterministic tag for the workbench image from the content of
// the embedded deploy/workbench build context: each embedded file's name and content are hashed
// together (sorted by filename, so the result doesn't depend on directory-listing order), and
// the first workbenchImageTagHashLen hex characters of the resulting sha256 digest become the
// tag suffix.
func workbenchImageTag() (string, error) {
	entries, err := fs.ReadDir(workbenchimage.Files, ".")
	if err != nil {
		return "", rerrors.Wrap(err, "error listing embedded workbench build context")
	}

	names := make([]string, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		names = append(names, entry.Name())
	}

	sort.Strings(names)

	hasher := sha256.New()

	for _, name := range names {
		content, err := fs.ReadFile(workbenchimage.Files, name)
		if err != nil {
			return "", rerrors.Wrap(err, "error reading embedded workbench build context file "+name)
		}

		_, err = hasher.Write([]byte(name + "\x00"))
		if err != nil {
			return "", rerrors.Wrap(err, "error hashing embedded workbench build context file name "+name)
		}

		_, err = hasher.Write(content)
		if err != nil {
			return "", rerrors.Wrap(err, "error hashing embedded workbench build context file "+name)
		}
	}

	sum := hasher.Sum(nil)
	hexSum := hex.EncodeToString(sum)

	return workbenchImageTagPrefix + hexSum[:workbenchImageTagHashLen], nil
}

// workbenchBuildContext packs the embedded deploy/workbench build context (workbenchimage.Files)
// into an in-memory tar archive, as required by the Docker Engine API's ImageBuild endpoint.
func workbenchBuildContext() (*bytes.Buffer, error) {
	var tarBuf bytes.Buffer

	tarWriter := tar.NewWriter(&tarBuf)

	walkErr := fs.WalkDir(workbenchimage.Files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(workbenchimage.Files, path)
		if err != nil {
			return err
		}

		header := tar.Header{
			Name: path,
			Mode: 0644,
			Size: int64(len(content)),
		}

		err = tarWriter.WriteHeader(&header)
		if err != nil {
			return err
		}

		_, err = tarWriter.Write(content)

		return err
	})
	if walkErr != nil {
		return nil, rerrors.Wrap(walkErr, "error walking embedded workbench build context")
	}

	err := tarWriter.Close()
	if err != nil {
		return nil, rerrors.Wrap(err, "error closing workbench build context tar writer")
	}

	return &tarBuf, nil
}
