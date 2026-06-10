package couchdb

import (
	"context"
	"net/http"

	kivik "github.com/go-kivik/kivik/v4"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (c *Client) CreateDatabase(ctx context.Context, name string) error {
	err := c.kivik.CreateDB(ctx, name)
	if err == nil {
		return nil
	}

	switch kivik.HTTPStatus(err) {
	case http.StatusPreconditionFailed:
		return user_errors.CouchDbDatabaseAlreadyExists
	case http.StatusBadRequest:
		return user_errors.InvalidCouchDbDatabaseName
	default:
		return rerrors.Wrap(err, "creating database")
	}
}
