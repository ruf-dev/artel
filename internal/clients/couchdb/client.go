package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.redsock.ru/rerrors"
)

type Config struct {
	BaseURL  string
	User     string
	Password string
}

type Client struct {
	baseURL  string
	http     *http.Client
	user     string
	password string
}

func New(cfg Config) *Client {
	return &Client{
		baseURL:  cfg.BaseURL,
		http:     &http.Client{},
		user:     cfg.User,
		password: cfg.Password,
	}
}

func (c *Client) Config() Config {
	return Config{
		BaseURL:  c.baseURL,
		User:     c.user,
		Password: c.password,
	}
}

func (c *Client) CreateDatabase(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/%s", c.baseURL, name), nil)
	if err != nil {
		return rerrors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(c.user, c.password)

	resp, err := c.http.Do(req)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	switch resp.StatusCode {
	case http.StatusBadRequest,
		http.StatusPreconditionFailed:
		var er CreateDbErrorResp
		err = json.Unmarshal(body, &er)
		if err != nil {
			return rerrors.Wrap(err, "failed to unmarshal response body")
		}

		return rerrors.Wrap(er.Error())
	default:
		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}
}

func (c *Client) DeleteDatabase(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", c.baseURL, name), nil)
	if err != nil {
		return rerrors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(c.user, c.password)

	resp, err := c.http.Do(req)
	if err != nil {
		return rerrors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rerrors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	return nil
}
