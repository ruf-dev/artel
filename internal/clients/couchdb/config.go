package couchdb

import (
	"net/url"

	"go.redsock.ru/rerrors"
)

type Config struct {
	BaseURL  string
	User     string
	Password string
}

func NewConfig(raw string) (Config, error) {
	if raw == "" {
		return Config{}, rerrors.New("couchdb_url is not set")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return Config{}, rerrors.Wrap(err, "invalid COUCHDB_URL")
	}

	var c Config
	if u.User != nil {
		c.User = u.User.Username()
		c.Password, _ = u.User.Password()
		u.User = nil
	}
	c.BaseURL = u.String()

	return c, nil
}
