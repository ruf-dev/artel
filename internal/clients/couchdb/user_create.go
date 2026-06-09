package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (c *Client) CreateUser(ctx context.Context, username, password string, roles []string) error {
	u := couchDBUser{
		ID:       fmt.Sprintf("org.couchdb.user:%s", username),
		Name:     username,
		Password: password,
		Roles:    roles,
		Type:     "user",
	}

	reqBody, err := json.Marshal(u)
	if err != nil {
		return rerrors.Wrap(err, "failed to marshal user")
	}

	url := fmt.Sprintf("%s/_users/org.couchdb.user:%s", c.baseURL, username)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(reqBody))
	if err != nil {
		return rerrors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(c.user, c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		return nil
	}
	if resp.StatusCode == http.StatusConflict {
		return rerrors.Wrap(user_errors.UserAlreadyExistInCouchDb)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read response body")
	}

	lockedErr := checkAccountLocked(resp.StatusCode, body)
	if lockedErr != nil {
		return lockedErr
	}

	return rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
}
