package couchdb

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	kivik "github.com/go-kivik/kivik/v4"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/utils"
	"go.redsock.ru/rerrors"
)

type couchDBUser struct {
	ID       string   `json:"_id"`
	Rev      string   `json:"_rev,omitempty"`
	Name     string   `json:"name"`
	Password string   `json:"password,omitempty"`
	Roles    []string `json:"roles"`
	Type     string   `json:"type"`
}

type UserInfo struct {
	Name  string
	Roles []string
}

type UserListEntry struct {
	Name  string
	Roles []string
}

func (c *Client) getUserDoc(ctx context.Context, username string) (couchDBUser, error) {
	db := c.kivik.DB("_users")
	docID := fmt.Sprintf("org.couchdb.user:%s", username)

	row := db.Get(ctx, docID)
	defer utils.CloseWithLog(row, "get user response")

	var u couchDBUser

	err := row.ScanDoc(&u)
	if err != nil {
		if kivik.HTTPStatus(err) == http.StatusForbidden && strings.Contains(err.Error(), "locked") {
			return couchDBUser{}, rerrors.Wrap(user_errors.CouchDbAccountLocked)
		}

		return couchDBUser{}, rerrors.Wrap(err, "scanning user doc")
	}

	return u, nil
}

func (c *Client) GetUser(ctx context.Context, username string) (UserInfo, error) {
	u, err := c.getUserDoc(ctx, username)
	if err != nil {
		return UserInfo{}, err
	}

	return UserInfo{Name: u.Name, Roles: u.Roles}, nil
}

func (c *Client) ListUsers(ctx context.Context) ([]UserListEntry, error) {
	db := c.kivik.DB("_users")

	rows := db.AllDocs(ctx, kivik.Params(map[string]any{"include_docs": true}))
	defer utils.CloseWithLog(rows, "list users response")

	var users []UserListEntry

	for rows.Next() {
		id, err := rows.ID()
		if err != nil {
			continue
		}

		if !strings.HasPrefix(id, "org.couchdb.user:") {
			continue
		}

		var doc couchDBUser

		err = rows.ScanDoc(&doc)
		if err != nil {
			continue
		}

		entry := UserListEntry{Name: doc.Name, Roles: doc.Roles}
		users = append(users, entry)
	}

	err := rows.Err()
	if err != nil {
		return nil, rerrors.Wrap(err, "iterating user rows")
	}

	return users, nil
}

func (c *Client) UpdateUser(ctx context.Context, username, password string, roles []string) error {
	existing, err := c.getUserDoc(ctx, username)
	if err != nil {
		return rerrors.Wrap(err, "getting user for update")
	}

	db := c.kivik.DB("_users")
	docID := fmt.Sprintf("org.couchdb.user:%s", username)

	u := couchDBUser{
		ID:       docID,
		Rev:      existing.Rev,
		Name:     existing.Name,
		Password: password,
		Roles:    roles,
		Type:     "user",
	}

	_, err = db.Put(ctx, docID, u)
	if err != nil {
		if kivik.HTTPStatus(err) == http.StatusForbidden && strings.Contains(err.Error(), "locked") {
			return rerrors.Wrap(user_errors.CouchDbAccountLocked)
		}

		return rerrors.Wrap(err, "updating user")
	}

	return nil
}

func (c *Client) DeleteUser(ctx context.Context, username string) error {
	existing, err := c.getUserDoc(ctx, username)
	if err != nil {
		return rerrors.Wrap(err, "getting user for delete")
	}

	db := c.kivik.DB("_users")
	docID := fmt.Sprintf("org.couchdb.user:%s", username)

	_, err = db.Delete(ctx, docID, existing.Rev)
	if err != nil {
		if kivik.HTTPStatus(err) == http.StatusForbidden && strings.Contains(err.Error(), "locked") {
			return rerrors.Wrap(user_errors.CouchDbAccountLocked)
		}

		return rerrors.Wrap(err, "deleting user")
	}

	return nil
}
