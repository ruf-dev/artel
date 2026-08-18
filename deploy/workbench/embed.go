// Package workbenchimage embeds the workbench Docker image's build context
// (Dockerfile + entrypoint.sh + tmux.conf) so it can be built directly through the Docker Engine
// API at runtime — see internal/clients/workbenchdocker/image.go — instead of requiring a manual
// `docker build` step on every docker host.
package workbenchimage

import "embed"

//go:embed Dockerfile entrypoint.sh tmux.conf
var Files embed.FS
