package domain

import (
	"time"

	"github.com/google/uuid"
)

type S3Instance struct {
	Uuid      uuid.UUID
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	UseSSL    bool
	PathStyle bool
	CreatedAt time.Time
}
