package imap

import (
	"context"
	"fmt"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
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

func (c *Client) ListEmails(_ context.Context, limit int) ([]domain.EmailMeta, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Logout()

	mbox, err := conn.Select("INBOX", true)
	if err != nil {
		return nil, rerrors.Wrap(err, "select inbox")
	}

	if mbox.Messages == 0 {
		return []domain.EmailMeta{}, nil
	}

	from := uint32(1)
	if mbox.Messages > uint32(limit) {
		from = mbox.Messages - uint32(limit) + 1
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, mbox.Messages)

	messages := make(chan *imap.Message, limit)
	err = conn.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}, messages)
	if err != nil {
		return nil, rerrors.Wrap(err, "fetch messages")
	}

	var result []domain.EmailMeta
	for msg := range messages {
		meta := domain.EmailMeta{
			Id:      fmt.Sprintf("%d", msg.Uid),
			Subject: msg.Envelope.Subject,
		}
		if len(msg.Envelope.From) > 0 {
			addr := msg.Envelope.From[0]
			meta.From = addr.Address()
		}
		if msg.Envelope.Date.IsZero() {
			meta.Date = ""
		} else {
			meta.Date = msg.Envelope.Date.UTC().Format("2006-01-02T15:04:05Z")
		}
		result = append(result, meta)
	}
	return result, nil
}

func (c *Client) ReadEmail(_ context.Context, uid string) (domain.EmailMessage, error) {
	conn, err := c.connect()
	if err != nil {
		return domain.EmailMessage{}, err
	}
	defer conn.Logout()

	if _, err := conn.Select("INBOX", true); err != nil {
		return domain.EmailMessage{}, rerrors.Wrap(err, "select inbox")
	}

	var uidNum uint32
	if _, err := fmt.Sscanf(uid, "%d", &uidNum); err != nil {
		return domain.EmailMessage{}, user_errors.InvalidEmailId
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uidNum)

	messages := make(chan *imap.Message, 1)
	section := &imap.BodySectionName{}
	err = conn.UidFetch(seqSet, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages)
	if err != nil {
		return domain.EmailMessage{}, rerrors.Wrap(err, "fetch message")
	}

	msg := <-messages
	if msg == nil {
		return domain.EmailMessage{}, user_errors.EmailMessageNotFound
	}

	result := domain.EmailMessage{
		Id:      uid,
		Subject: msg.Envelope.Subject,
	}
	if len(msg.Envelope.From) > 0 {
		result.From = msg.Envelope.From[0].Address()
	}
	if len(msg.Envelope.To) > 0 {
		addrs := make([]string, len(msg.Envelope.To))
		for i, a := range msg.Envelope.To {
			addrs[i] = a.Address()
		}
		result.To = strings.Join(addrs, ", ")
	}
	if !msg.Envelope.Date.IsZero() {
		result.Date = msg.Envelope.Date.UTC().Format("2006-01-02T15:04:05Z")
	}

	body := msg.GetBody(section)
	if body != nil {
		buf := new(strings.Builder)
		buf.Grow(4096)
		b := make([]byte, 4096)
		for {
			n, readErr := body.Read(b)
			if n > 0 {
				buf.Write(b[:n])
			}
			if readErr != nil {
				break
			}
		}
		result.Body = buf.String()
	}

	return result, nil
}

func (c *Client) ListFolders(_ context.Context) ([]string, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Logout()

	mailboxes := make(chan *imap.MailboxInfo, 16)
	done := make(chan error, 1)
	go func() {
		done <- conn.List("", "*", mailboxes)
	}()

	var folders []string
	for m := range mailboxes {
		folders = append(folders, m.Name)
	}
	if err := <-done; err != nil {
		return nil, rerrors.Wrap(err, "list folders")
	}
	return folders, nil
}

func (c *Client) connect() (*client.Client, error) {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	conn, err := client.DialTLS(addr, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "connect to imap server")
	}

	if err := conn.Login(c.email, c.password); err != nil {
		conn.Logout()
		return nil, rerrors.Wrap(err, "imap login failed")
	}

	return conn, nil
}
