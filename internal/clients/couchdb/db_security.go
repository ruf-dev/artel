package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/utils"
)

type securityNamesRoles struct {
	Names []string `json:"names"`
	Roles []string `json:"roles"`
}

type securityDoc struct {
	Admins  securityNamesRoles `json:"admins"`
	Members securityNamesRoles `json:"members"`
}

func (c *Client) SetDatabaseSecurity(ctx context.Context, dbName string, memberUsernames ...string) error {
	doc := securityDoc{
		Admins:  securityNamesRoles{Names: []string{}, Roles: []string{}},
		Members: securityNamesRoles{Names: memberUsernames, Roles: []string{}},
	}

	reqBody, err := json.Marshal(doc)
	if err != nil {
		return rerrors.Wrap(err, "marshal security doc")
	}

	reqReader := bytes.NewReader(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
		fmt.Sprintf("%s/%s/_security", c.baseURL, dbName),
		reqReader,
	)
	if err != nil {
		return rerrors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(c.user, c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute request")
	}
	defer utils.CloseWithLog(resp.Body, "security response body")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		lockedErr := checkAccountLocked(resp.StatusCode, body)
		if lockedErr != nil {
			return lockedErr
		}
		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	return nil
}
