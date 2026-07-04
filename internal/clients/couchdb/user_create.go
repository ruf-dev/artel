package couchdb

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	kivik "github.com/go-kivik/kivik/v4"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (c *Client) CreateUser(ctx context.Context, username, password string, roles []string) error {
	docID := fmt.Sprintf("org.couchdb.user:%s", username)

	u := couchDBUser{
		ID:       docID,
		Name:     username,
		Password: password,
		Roles:    roles,
		Type:     "user",
	}

	db := c.kivik.DB("_users")

	_, err := db.Put(ctx, docID, u)
	if err == nil {
		return nil
	}

	switch kivik.HTTPStatus(err) {
	case http.StatusConflict:
		return rerrors.Wrap(user_errors.UserAlreadyExistInCouchDb)
	case http.StatusForbidden:
		if strings.Contains(err.Error(), "locked") {
			return rerrors.Wrap(user_errors.CouchDbAccountLocked)
		}

		return rerrors.Wrap(err, "creating user")
	default:
		return rerrors.Wrap(err, "creating user")
	}
}
