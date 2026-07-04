package couchdb

import (
	"context"
	"fmt"

	"go.redsock.ru/rerrors"
)

func (c *Client) DatabaseURL(name string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, name)
}

func (c *Client) ListDatabases(ctx context.Context) ([]string, error) {
	dbs, err := c.kivik.AllDBs(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "listing databases")
	}

	return dbs, nil
}

func (c *Client) DeleteDatabase(ctx context.Context, name string) error {
	err := c.kivik.DestroyDB(ctx, name)
	if err != nil {
		return rerrors.Wrap(err, "deleting database")
	}

	return nil
}
