package domain

import (
	"time"

	"github.com/google/uuid"
)

// DockerHost is one admin-registered Docker daemon endpoint in the workbench pool. Unlike
// CouchInstance/S3Instance it has no single credential blob; instead it carries three optional
// TLS/mTLS fields (see
// migrations/062_docker_hosts_tls.sql) for the remote-daemon case, empty for a local
// unix-socket/unauthenticated-tcp daemon.
type DockerHost struct {
	Uuid      uuid.UUID
	Url       string
	CreatedAt time.Time

	// CaCert/ClientCert/ClientKey are the decrypted, in-memory PEM blobs used to build a
	// TLS-configured Docker client (internal/clients/workbenchdocker). They are never populated
	// by Get/List (creds-free, admin list/view queries) — only by DockerHosts.GetWithCreds,
	// mirroring the S3Instances Get vs GetWithCreds split. Empty means "no TLS configured".
	CaCert     string
	ClientCert string
	ClientKey  string
}
