package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/utils"
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

func (c *Client) Send(ctx context.Context, to, subject, body string) (err error) {
	opStart := time.Now()
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	defer func() {
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Str("op", "send").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("smtp client")
		} else {
			log.Ctx(ctx).Debug().Str("op", "send").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("smtp client")
		}
	}()

	auth := smtp.PlainAuth("", c.email, c.password, c.host)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		c.email, to, subject, body)

	err = smtp.SendMail(addr, auth, c.email, []string{to}, []byte(msg))
	if err != nil {
		return rerrors.Wrap(err, "send mail")
	}

	return nil
}

// TestConnection dials the SMTP server, upgrades to TLS when offered, and authenticates, without
// sending any mail, to verify the configured host/port/credentials work. The dial honors ctx
// cancellation, and — when ctx carries a deadline — the whole handshake (STARTTLS/Auth included)
// is bounded by it, so a host that never responds fails once ctx expires instead of hanging
// indefinitely.
func (c *Client) TestConnection(ctx context.Context) (err error) {
	opStart := time.Now()
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	defer func() {
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Str("op", "test_connection").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("smtp client")
		} else {
			log.Ctx(ctx).Debug().Str("op", "test_connection").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("smtp client")
		}
	}()

	var dialer net.Dialer

	rawConn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return rerrors.Wrap(err, "dial smtp server")
	}

	if deadline, ok := ctx.Deadline(); ok {
		err = rawConn.SetDeadline(deadline)
		if err != nil {
			return rerrors.Wrap(err, "set smtp connection deadline")
		}
	}

	conn, err := smtp.NewClient(rawConn, c.host)
	if err != nil {
		return rerrors.Wrap(err, "init smtp client")
	}
	defer utils.CloseWithLog(conn, "smtp connection test")

	startTLSSupported, _ := conn.Extension("STARTTLS")
	if startTLSSupported {
		tlsConfig := &tls.Config{ServerName: c.host}

		err = conn.StartTLS(tlsConfig)
		if err != nil {
			return rerrors.Wrap(err, "starttls smtp server")
		}
	}

	auth := smtp.PlainAuth("", c.email, c.password, c.host)

	err = conn.Auth(auth)
	if err != nil {
		return rerrors.Wrap(err, "smtp auth failed")
	}

	return nil
}
