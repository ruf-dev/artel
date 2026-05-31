package couchdb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service/user_errors"
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
		http:     &http.Client{Transport: newLoggingTransport()},
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

func (c *Client) Setup(ctx context.Context) error {
	for _, db := range []string{"_users", "_replicator"} {
		err := c.CreateDatabase(ctx, db)
		if err != nil && !errors.Is(err, user_errors.CouchDbDatabaseAlreadyExists) {
			return rerrors.Wrap(err, "create system db "+db)
		}
	}
	return nil
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
