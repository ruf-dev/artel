package user_errors

import (
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc/codes"
)

var (
	Unauthenticated = rerrors.New("unauthenticated", codes.Unauthenticated)

	InvalidDatabaseName   = rerrors.New("invalid database name", codes.InvalidArgument)
	DatabaseAlreadyExists = rerrors.New("database already exists", codes.FailedPrecondition)
)
