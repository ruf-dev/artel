package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.redsock.ru/rerrors"
)

type couchDBUser struct {
	ID       string   `json:"_id"`
	Name     string   `json:"name"`
	Password string   `json:"password,omitempty"`
	Roles    []string `json:"roles"`
	Type     string   `json:"type"`
}

type UserInfo struct {
	Name  string
	Roles []string
}

func (c *Client) GetUser(ctx context.Context, username string) (UserInfo, error) {
	url := fmt.Sprintf("%s/_users/org.couchdb.user:%s", c.baseURL, username)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return UserInfo{}, rerrors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(c.user, c.password)

	resp, err := c.http.Do(req)
	if err != nil {
		return UserInfo{}, rerrors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInfo{}, rerrors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return UserInfo{}, rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	var u couchDBUser
	err = json.Unmarshal(body, &u)
	if err != nil {
		return UserInfo{}, rerrors.Wrap(err, "failed to parse response")
	}

	info := UserInfo{
		Name:  u.Name,
		Roles: u.Roles,
	}
	return info, nil
}

func (c *Client) UpdateUser(ctx context.Context, username, password string, roles []string) error {
	existing, err := c.GetUser(ctx, username)
	if err != nil {
		return rerrors.Wrap(err, "get user for update")
	}

	u := couchDBUser{
		ID:       fmt.Sprintf("org.couchdb.user:%s", username),
		Name:     existing.Name,
		Password: password,
		Roles:    roles,
		Type:     "user",
	}

	body, err := json.Marshal(u)
	if err != nil {
		return rerrors.Wrap(err, "failed to marshal user")
	}

	url := fmt.Sprintf("%s/_users/org.couchdb.user:%s", c.baseURL, username)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(respBody)))
	}

	return nil
}

func (c *Client) DeleteUser(ctx context.Context, username string) error {
	url := fmt.Sprintf("%s/_users/org.couchdb.user:%s", c.baseURL, username)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return rerrors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(c.user, c.password)

	resp, err := c.http.Do(req)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(respBody)))
	}

	return nil
}
