package user_errors

import (
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc/codes"
)

var (
	Unauthenticated = rerrors.New("unauthenticated", codes.Unauthenticated)
	NotFound        = rerrors.New("not found", codes.NotFound)
	AlreadyExists   = rerrors.New("already exists", codes.AlreadyExists)

	InvalidCouchDbDatabaseName   = rerrors.New("invalid database name", codes.InvalidArgument)
	CouchDbDatabaseAlreadyExists = rerrors.New("database already exists", codes.FailedPrecondition)
	UserAlreadyExistInCouchDb    = rerrors.New("user already exists in couch db", codes.AlreadyExists)
)
