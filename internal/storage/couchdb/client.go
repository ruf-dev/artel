package couchdb

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/config"
)

type Client struct {
	baseURL string
	http    *http.Client
	user    string
	password string
}

func New(cfg config.CouchDB) *Client {
	baseURL := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	return &Client{
		baseURL:  baseURL,
		http:     &http.Client{},
		user:     cfg.User,
		password: cfg.Password,
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return rerrors.New(fmt.Sprintf("unexpected status code %d: %s", resp.StatusCode, string(body)))
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
		return rerrors.New(fmt.Sprintf("unexpected status code %d: %s", resp.StatusCode, string(body)))
	}

	return nil
}
