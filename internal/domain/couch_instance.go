package domain

import (
	"time"

	"github.com/google/uuid"
)

type CouchInstance struct {
	Uuid      uuid.UUID
	Url       string
	Username  string
	Password  string
	CreatedAt time.Time
}

type CouchInstanceWithAccount struct {
	Instance CouchInstance
	Account  *CouchAccount
}
