package couchdb

import (
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

var ErrUnknownError = rerrors.New("unknown error")

const (
	ErrorCodeIllegalDatabaseName = "illegal_database_name"
	ErrorCodeAlreadyExists       = "file_exists"
)

type CreateDbErrorResp struct {
	ErrorCode string `json:"error"`
	Reason    string `json:"reason"`
}

func (e CreateDbErrorResp) Error() error {
	switch e.ErrorCode {
	case ErrorCodeIllegalDatabaseName:
		return user_errors.InvalidDatabaseName
	case ErrorCodeAlreadyExists:
		return user_errors.DatabaseAlreadyExists
	default:
		return ErrUnknownError
	}
}
