package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailAccount struct {
	Uuid      uuid.UUID
	UserUuid  uuid.UUID
	Email     string
	ImapHost  string
	ImapPort  int
	SmtpHost  string
	SmtpPort  int
	Password  string
	CreatedAt time.Time
}

type EmailMeta struct {
	Id      string
	From    string
	Subject string
	Date    string
}

type EmailMessage struct {
	Id      string
	From    string
	To      string
	Subject string
	Date    string
	Body    string
}
