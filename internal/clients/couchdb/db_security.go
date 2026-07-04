package couchdb

import (
	"context"

	kivik "github.com/go-kivik/kivik/v4"
	"go.redsock.ru/rerrors"
)

func (c *Client) SetDatabaseSecurity(ctx context.Context, dbName string, memberUsernames ...string) error {
	db := c.kivik.DB(dbName)

	sec := &kivik.Security{
		Admins:  kivik.Members{Names: []string{}, Roles: []string{}},
		Members: kivik.Members{Names: memberUsernames, Roles: []string{}},
	}

	err := db.SetSecurity(ctx, sec)
	if err != nil {
		return rerrors.Wrap(err, "setting database security")
	}

	return nil
}

func (c *Client) GrantDatabaseAccess(ctx context.Context, dbName, username string) error {
	db := c.kivik.DB(dbName)

	sec, err := db.Security(ctx)
	if err != nil {
		return rerrors.Wrap(err, "getting database security")
	}

	for _, name := range sec.Members.Names {
		if name == username {
			return nil
		}
	}

	sec.Members.Names = append(sec.Members.Names, username)

	err = db.SetSecurity(ctx, sec)
	if err != nil {
		return rerrors.Wrap(err, "setting database security")
	}

	return nil
}

func (c *Client) GetUserDatabaseAccess(ctx context.Context, username string) ([]string, error) {
	allDBs, err := c.kivik.AllDBs(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "listing databases")
	}

	var granted []string

	for _, dbName := range allDBs {
		if len(dbName) > 0 && dbName[0] == '_' {
			continue
		}

		sec, err := c.kivik.DB(dbName).Security(ctx)
		if err != nil {
			return nil, rerrors.Wrap(err, "getting security for "+dbName)
		}

		for _, name := range sec.Members.Names {
			if name == username {
				granted = append(granted, dbName)

				break
			}
		}
	}

	return granted, nil
}

func (c *Client) RevokeDatabaseAccess(ctx context.Context, dbName, username string) error {
	db := c.kivik.DB(dbName)

	sec, err := db.Security(ctx)
	if err != nil {
		return rerrors.Wrap(err, "getting database security")
	}

	filtered := make([]string, 0, len(sec.Members.Names))

	for _, name := range sec.Members.Names {
		if name != username {
			filtered = append(filtered, name)
		}
	}

	sec.Members.Names = filtered

	err = db.SetSecurity(ctx, sec)
	if err != nil {
		return rerrors.Wrap(err, "setting database security")
	}

	return nil
}
