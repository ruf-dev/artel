package couchdb

import (
	"go.redsock.ru/rerrors"
)

var ErrUnknownError = rerrors.New("unknown error")

const (
	ErrorCodeIllegalDatabaseName = "illegal_database_name"
	ErrorCodeAlreadyExists       = "file_exists"
)
