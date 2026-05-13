package domain

import "time"

type PendingAuthCode struct {
	Code          string
	RawToken      string
	CodeChallenge string
	RedirectUri   string
	ClientId      string
	ExpiresAt     time.Time
}
