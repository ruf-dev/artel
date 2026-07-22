package imap

import (
	"context"
	"fmt"
	"math"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/utils"

	"github.com/rs/zerolog/log"
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

// ListEmailsOptions bounds a list_emails page. AfterUid/BeforeUid are raw UID strings (same
// format as ReadEmail's uid param); at most one may be set by the caller.
type ListEmailsOptions struct {
	Limit     int
	AfterUid  *string
	BeforeUid *string
}

func (c *Client) ListEmails(ctx context.Context, opts ListEmailsOptions) (emails []domain.EmailMeta, err error) {
	opStart := time.Now()
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	defer func() {
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Str("op", "list_emails").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("imap client")
		} else {
			log.Ctx(ctx).Debug().Str("op", "list_emails").Str("host", addr).Int("count", len(emails)).Dur("dur", time.Since(opStart)).Msg("imap client")
		}
	}()

	conn, err := c.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer utils.CallWithLog(conn.Logout, "logout error during listing emails")

	mbox, err := conn.Select("INBOX", true)
	if err != nil {
		return nil, rerrors.Wrap(err, "select inbox")
	}

	if mbox.Messages == 0 {
		return []domain.EmailMeta{}, nil
	}

	var afterUid, beforeUid *uint32

	if opts.AfterUid != nil {
		v, parseErr := parseUid(*opts.AfterUid)
		if parseErr != nil {
			return nil, parseErr
		}

		afterUid = &v
	}

	if opts.BeforeUid != nil {
		v, parseErr := parseUid(*opts.BeforeUid)
		if parseErr != nil {
			return nil, parseErr
		}

		beforeUid = &v
	}

	start, stop, hasCursor, ok := emailCursorUIDBounds(afterUid, beforeUid)
	if hasCursor && !ok {
		return []domain.EmailMeta{}, nil
	}

	seqSet := new(imap.SeqSet)

	if hasCursor {
		seqSet.AddRange(start, stop)
	} else {
		from := uint32(1)
		if mbox.Messages > uint32(opts.Limit) {
			from = mbox.Messages - uint32(opts.Limit) + 1
		}

		seqSet.AddRange(from, mbox.Messages)
	}

	result, err := fetchEmailMeta(conn, seqSet, hasCursor)
	if err != nil {
		return nil, err
	}

	if hasCursor {
		result = trimEmailPage(result, opts.Limit, beforeUid != nil)
	}

	return result, nil
}

// emailCursorUIDBounds computes the inclusive UID range to fetch for a cursor-based list_emails
// page. stop == 0 means unbounded ("*", up to the newest UID) — imap.SeqSet.AddRange already
// treats 0 as that sentinel. hasCursor is false when neither bound is set (caller falls back to
// the existing sequence-number path). ok is false when hasCursor is true but the window is
// provably empty (no UID can be < 1, and MaxUint32 can't be exceeded), so the caller should
// return an empty page without an IMAP round trip.
func emailCursorUIDBounds(afterUid, beforeUid *uint32) (start, stop uint32, hasCursor, ok bool) {
	switch {
	case afterUid != nil:
		if *afterUid == math.MaxUint32 {
			return 0, 0, true, false
		}

		return *afterUid + 1, 0, true, true

	case beforeUid != nil:
		if *beforeUid <= 1 {
			return 0, 0, true, false
		}

		return 1, *beforeUid - 1, true, true

	default:
		return 0, 0, false, false
	}
}

// trimEmailPage bounds an ascending-UID-sorted page to at most limit entries. after_uid pages
// keep the front (smallest UIDs — closest to the cursor, continuing the forward walk);
// before_uid pages keep the back (largest UIDs still below the cursor — also closest to it).
// Both keep ascending order, matching list_emails' single ordering contract.
func trimEmailPage(emails []domain.EmailMeta, limit int, keepNewestEnd bool) []domain.EmailMeta {
	if limit <= 0 || len(emails) <= limit {
		return emails
	}

	if keepNewestEnd {
		return emails[len(emails)-limit:]
	}

	return emails[:limit]
}

// parseUid parses a caller-supplied UID string (read_email's id, list_emails' after_uid/
// before_uid) into the numeric IMAP UID it identifies.
func parseUid(s string) (uint32, error) {
	var uid uint32

	_, err := fmt.Sscanf(s, "%d", &uid)
	if err != nil {
		return 0, user_errors.InvalidEmailId
	}

	return uid, nil
}

// fetchEmailMeta runs seqSet through FETCH (or UID FETCH when uid) and drains responses
// concurrently, so it's safe regardless of how many messages seqSet matches. A fixed-size
// buffered channel drained only after the command returns (as ReadEmail does, and as ListEmails
// used to) deadlocks once matches exceed the buffer — a risk cursor ranges (open-ended or wide)
// introduce, since their match count isn't known up front. ListFolders below already uses this
// concurrent-drain pattern for the same reason. Results are sorted ascending by UID before
// mapping, since IMAP doesn't guarantee response order.
func fetchEmailMeta(conn *client.Client, seqSet *imap.SeqSet, uid bool) ([]domain.EmailMeta, error) {
	fetch := conn.Fetch
	if uid {
		fetch = conn.UidFetch
	}

	messages := make(chan *imap.Message, 16)
	done := make(chan error, 1)

	go func() {
		done <- fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}, messages)
	}()

	var msgs []*imap.Message
	for msg := range messages {
		msgs = append(msgs, msg)
	}

	err := <-done
	if err != nil {
		return nil, rerrors.Wrap(err, "fetch messages")
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Uid < msgs[j].Uid
	})

	result := make([]domain.EmailMeta, 0, len(msgs))
	for _, msg := range msgs {
		result = append(result, emailMetaFromMessage(msg))
	}

	return result, nil
}

func emailMetaFromMessage(msg *imap.Message) domain.EmailMeta {
	meta := domain.EmailMeta{
		Id:      fmt.Sprintf("%d", msg.Uid),
		Subject: msg.Envelope.Subject,
	}

	if len(msg.Envelope.From) > 0 {
		meta.From = msg.Envelope.From[0].Address()
	}

	if !msg.Envelope.Date.IsZero() {
		meta.Date = msg.Envelope.Date.UTC().Format("2006-01-02T15:04:05Z")
	}

	return meta
}

func (c *Client) ReadEmail(ctx context.Context, uid string) (message domain.EmailMessage, err error) {
	opStart := time.Now()
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	defer func() {
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Str("op", "read_email").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("imap client")
		} else {
			log.Ctx(ctx).Debug().Str("op", "read_email").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("imap client")
		}
	}()

	conn, err := c.connect(ctx)
	if err != nil {
		return domain.EmailMessage{}, err
	}
	defer utils.CallWithLog(conn.Logout, "logout error during reading emails")

	_, err = conn.Select("INBOX", true)
	if err != nil {
		return domain.EmailMessage{}, rerrors.Wrap(err, "select inbox")
	}

	var uidNum uint32
	_, err = fmt.Sscanf(uid, "%d", &uidNum)
	if err != nil {
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

// TestConnection dials the IMAP server and logs in, then immediately logs out, to verify the
// configured host/port/credentials work without fetching any mailbox data. Unlike connect(), the
// dial (and the deadline covering the subsequent login) is bounded by ctx's deadline, so a host
// that never responds fails once ctx expires instead of hanging indefinitely.
func (c *Client) TestConnection(ctx context.Context) (err error) {
	opStart := time.Now()

	dialer := &net.Dialer{}

	if deadline, ok := ctx.Deadline(); ok {
		dialer.Timeout = time.Until(deadline)
	}

	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	defer func() {
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Str("op", "test_connection").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("imap client")
		} else {
			log.Ctx(ctx).Debug().Str("op", "test_connection").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("imap client")
		}
	}()

	conn, err := client.DialWithDialerTLS(dialer, addr, nil)
	if err != nil {
		return rerrors.Wrap(err, "connect to imap server")
	}
	defer utils.CallWithLog(conn.Logout, "logout error during connection test")

	err = conn.Login(c.email, c.password)
	if err != nil {
		return rerrors.Wrap(err, "imap login failed")
	}

	return nil
}

func (c *Client) ListFolders(ctx context.Context) (names []string, err error) {
	opStart := time.Now()
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	defer func() {
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Str("op", "list_folders").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("imap client")
		} else {
			log.Ctx(ctx).Debug().Str("op", "list_folders").Str("host", addr).Int("count", len(names)).Dur("dur", time.Since(opStart)).Msg("imap client")
		}
	}()

	conn, err := c.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer utils.CallWithLog(conn.Logout, "logout error during listing folders")

	mailboxes := make(chan *imap.MailboxInfo, 16)
	done := make(chan error, 1)

	go func() {
		done <- conn.List("", "*", mailboxes)
	}()

	var folders []string
	for m := range mailboxes {
		folders = append(folders, m.Name)
	}

	err = <-done
	if err != nil {
		return nil, rerrors.Wrap(err, "list folders")
	}

	return folders, nil
}

func (c *Client) connect(ctx context.Context) (conn *client.Client, err error) {
	opStart := time.Now()
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	defer func() {
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Str("op", "connect").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("imap client")
		} else {
			log.Ctx(ctx).Debug().Str("op", "connect").Str("host", addr).Dur("dur", time.Since(opStart)).Msg("imap client")
		}
	}()

	conn, err = client.DialTLS(addr, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "connect to imap server")
	}

	err = conn.Login(c.email, c.password)
	if err != nil {
		utils.CallWithLog(conn.Logout, "logout error during connect")

		return nil, rerrors.Wrap(err, "imap login failed")
	}

	return conn, nil
}
