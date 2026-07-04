package smtp

import (
	"context"
	"fmt"
	"net/smtp"

	"go.redsock.ru/rerrors"
)

type Client struct {
	host     string
	port     int
	email    string
	password string
}

func New(host string, port int, email, password string) *Client {
	return &Client{
		host:     host,
		port:     port,
		email:    email,
		password: password,
	}
}

func (c *Client) Send(_ context.Context, to, subject, body string) error {
	auth := smtp.PlainAuth("", c.email, c.password, c.host)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		c.email, to, subject, body)

	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	err := smtp.SendMail(addr, auth, c.email, []string{to}, []byte(msg))
	if err != nil {
		return rerrors.Wrap(err, "send mail")
	}

	return nil
}
