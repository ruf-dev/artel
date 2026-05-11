package couchdb

import (
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type CreateUserErrorResp struct {
	ErrorCode string `json:"error"`
	Reason    string `json:"reason"`
}

func (e CreateUserErrorResp) Error() error {
	switch e.ErrorCode {
	case ErrorCodeIllegalDatabaseName:
		return user_errors.InvalidCouchDbDatabaseName
	case ErrorCodeAlreadyExists:
		return user_errors.CouchDbDatabaseAlreadyExists
	default:
		return ErrUnknownError
	}
}
