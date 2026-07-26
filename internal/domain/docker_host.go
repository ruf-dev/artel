package domain

import (
	"time"

	"github.com/google/uuid"
)

// DockerHost is one admin-registered Docker daemon endpoint in the workbench pool — see
// docs/workbench/02_docker_topology.md. It carries no credentials (unlike CouchInstance/
// S3Instance), just a bare connection string.
type DockerHost struct {
	Uuid      uuid.UUID
	Url       string
	CreatedAt time.Time
}
