// Package workbenchimage embeds the workbench Docker image's build context (Dockerfile,
// entrypoint.sh, tmux.conf, and bridge.tar — a packaged copy of the bridge/ Go module tree, see
// ./gen_bridge_tar.go) so it can be built directly through the Docker Engine API at runtime — see
// internal/clients/workbenchdocker/image.go — instead of requiring a manual `docker build` step
// on every docker host.
//
// bridge.tar exists, rather than naming the bridge directory directly in the //go:embed pattern
// below, because Go's //go:embed refuses to embed anything inside a directory that contains its
// own go.mod — confirmed empirically (go1.26.3): "pattern bridge: cannot embed directory bridge:
// in different module", identically whether the pattern names the directory, a subdirectory, or
// go.mod itself, and unaffected by the "all:" prefix. bridge/ has to keep its own go.mod (it's a
// genuinely separate module — see bridge/README.md), so gen_bridge_tar.go packages it into one
// ordinary file instead, which //go:embed has no objection to. Regenerate it after any change
// under bridge/ with `go generate ./deploy/workbench/...`; internal/clients/workbenchdocker/
// image.go unpacks it back into the Docker build context and hashes its bytes into the
// content-addressed image tag exactly like the other embedded files, so a stale bridge.tar is
// the only way a bridge source change could go unnoticed.
package workbenchimage

import "embed"

//go:generate go run gen_bridge_tar.go

//go:embed Dockerfile entrypoint.sh tmux.conf bridge.tar
var Files embed.FS
