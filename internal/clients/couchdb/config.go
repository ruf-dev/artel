package couchdb

import (
	"net/url"
	"os"

	"go.redsock.ru/rerrors"
)

type Config struct {
	BaseURL  string
	User     string
	Password string
}

func NewConfig() (Config, error) {
	raw := os.Getenv("COUCHDB_URL")
	if raw == "" {
		return Config{}, rerrors.New("COUCHDB_URL is not set")
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
