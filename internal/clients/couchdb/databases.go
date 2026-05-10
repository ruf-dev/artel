package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.redsock.ru/rerrors"
)

func (c *Client) DatabaseURL(name string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, name)
}

func (c *Client) ListDatabases(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/_all_dbs", c.baseURL), nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(c.user, c.password)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to execute request")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, rerrors.New(fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(body)))
	}

	var names []string
	err = json.Unmarshal(body, &names)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to parse response")
	}

	return names, nil
}