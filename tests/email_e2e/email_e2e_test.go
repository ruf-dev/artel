//go:build e2e

package email_e2e_test

// Covers Part A's agentic IMAP write operations (MoveMessage, CreateFolder, SetSeen) plus the new
// folder-scoped ListEmails/ReadEmail, end to end against a real IMAP/SMTP server: greenmail
// (tests/docker-compose.yaml), a throwaway in-memory mail server with no persistence.
//
// Unlike the other tests/*_e2e suites, this one does NOT bring up Postgres/CouchDB/S3 or run the
// tests/bootstrap provisioning — there's no DB row involved in exercising an IMAP client method,
// so the suite drives internal/clients/imap.Client directly rather than going through the mcp_tools
// row + mom.ExecuteToolForKey dispatch path.
//
// TLS: greenmail presents a self-signed certificate on its IMAPS test port. The production dial
// path (executors.EmailExecutor.executeImap → imap.New with no options) intentionally keeps a
// strict certificate check — see imap.WithTLSConfig's doc comment — so it can't reach greenmail
// as-is. This suite instead constructs its own imap.Client via
// imap.New(..., imap.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})), which is the sanctioned
// test-only seam, and asserts against the imap.Client methods directly (the real surface Part A
// added). The executor's params-decoding switch (internal/service/v1/mcp/executors/email.go) is
// covered by the package's existing unit tests, which don't need real IMAP.
//
// Run: docker compose -f tests/docker-compose.yaml up -d --wait
//      go test -tags e2e ./tests/email_e2e/...

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/ruf-dev/artel/internal/clients/imap"
	"github.com/ruf-dev/artel/internal/domain"
)

func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}

	return def
}

func envIntOrDefault(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	parsed, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	return parsed
}

const (
	greenmailUser = "tester@localhost"
	greenmailPass = "pw"

	pollInterval = 200 * time.Millisecond
	pollTimeout  = 15 * time.Second
)

type EmailSuite struct {
	suite.Suite

	client   *imap.Client
	smtpAddr string
}

func TestEmail(t *testing.T) {
	suite.Run(t, new(EmailSuite))
}

func (s *EmailSuite) SetupSuite() {
	host := envOrDefault("GREENMAIL_HOST", "localhost")
	imapPort := envIntOrDefault("GREENMAIL_IMAPS_PORT", 13993)
	smtpPort := envIntOrDefault("GREENMAIL_SMTP_PORT", 13025)

	s.client = imap.New(host, imapPort, greenmailUser, greenmailPass,
		imap.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	s.smtpAddr = fmt.Sprintf("%s:%d", host, smtpPort)

	err := s.waitForImapReady()
	s.Require().NoError(err, "greenmail IMAPS not reachable — is tests/docker-compose.yaml up (--wait)?")
}

// waitForImapReady retries a cheap IMAP call until it succeeds or pollTimeout elapses. The
// compose healthcheck only proves the TCP port is accepting connections, not that greenmail's
// IMAP protocol handler is ready to authenticate.
func (s *EmailSuite) waitForImapReady() error {
	deadline := time.Now().Add(pollTimeout)

	var lastErr error

	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_, lastErr = s.client.ListFolders(ctx)
		cancel()

		if lastErr == nil {
			return nil
		}

		time.Sleep(pollInterval)
	}

	return lastErr
}

func (s *EmailSuite) sendTestEmail(subject, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		greenmailUser, greenmailUser, subject, body)

	return smtp.SendMail(s.smtpAddr, nil, greenmailUser, []string{greenmailUser}, []byte(msg))
}

// waitForSubjectInFolder polls ListEmails(folder) until a message with the exact subject shows
// up (SMTP delivery into the mailbox greenmail exposes over IMAP isn't synchronous with the SMTP
// accept) or pollTimeout elapses, returning the matched meta.
func (s *EmailSuite) waitForSubjectInFolder(folder, subject string) domain.EmailMeta {
	deadline := time.Now().Add(pollTimeout)

	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		emails, err := s.client.ListEmails(ctx, imap.ListEmailsOptions{Folder: folder, Limit: 50})
		cancel()

		s.Require().NoError(err)

		for _, e := range emails {
			if e.Subject == subject {
				return e
			}
		}

		time.Sleep(pollInterval)
	}

	s.FailNowf("message not found", "subject %q never appeared in folder %q", subject, folder)

	return domain.EmailMeta{}
}

func (s *EmailSuite) subjectPresentInFolder(folder, subject string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	emails, err := s.client.ListEmails(ctx, imap.ListEmailsOptions{Folder: folder, Limit: 100})
	s.Require().NoError(err)

	for _, e := range emails {
		if e.Subject == subject {
			return true
		}
	}

	return false
}

// Test1_SendAppearsInInbox: a message sent via SMTP shows up in INBOX via the folder-scoped
// ListEmails added in Part A.
func (s *EmailSuite) Test1_SendAppearsInInbox() {
	subject := fmt.Sprintf("part-a send %d", time.Now().UnixNano())

	err := s.sendTestEmail(subject, "hello from the send test")
	s.Require().NoError(err)

	meta := s.waitForSubjectInFolder("INBOX", subject)
	s.Require().Equal(subject, meta.Subject)
	s.Require().NotEmpty(meta.Id)
}

// Test2_CreateFolderAndMoveMessage: CreateFolder makes a new mailbox, MoveMessage relocates a
// message into it — it disappears from the source folder's ListEmails and becomes visible (both
// via ListEmails and ReadEmail) in the destination folder.
func (s *EmailSuite) Test2_CreateFolderAndMoveMessage() {
	stamp := time.Now().UnixNano()
	subject := fmt.Sprintf("part-a move %d", stamp)
	archiveFolder := fmt.Sprintf("Archive-%d", stamp)

	err := s.sendTestEmail(subject, "hello from the move test")
	s.Require().NoError(err)

	meta := s.waitForSubjectInFolder("INBOX", subject)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = s.client.CreateFolder(ctx, archiveFolder)
	cancel()
	s.Require().NoError(err)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err = s.client.MoveMessage(ctx, meta.Id, "INBOX", archiveFolder)
	cancel()
	s.Require().NoError(err)

	s.Require().False(s.subjectPresentInFolder("INBOX", subject), "message still in INBOX after MoveMessage")

	movedMeta := s.waitForSubjectInFolder(archiveFolder, subject)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	full, err := s.client.ReadEmail(ctx, movedMeta.Id, archiveFolder)
	cancel()
	s.Require().NoError(err)
	s.Require().Equal(subject, full.Subject)
}

// Test3_SetSeen: SetSeen(true) then SetSeen(false) both succeed against a real message, and the
// message stays readable afterward. The imap.Client surface added in Part A doesn't expose flag
// state through ListEmails/ReadEmail (domain.EmailMeta/EmailMessage carry no Flags field), so this
// doesn't assert on \Seen directly — it asserts the round trip is a no-op error-wise and the
// message isn't disturbed, which is what the task scoped as acceptable coverage for this
// operation.
func (s *EmailSuite) Test3_SetSeen() {
	subject := fmt.Sprintf("part-a seen %d", time.Now().UnixNano())

	err := s.sendTestEmail(subject, "hello from the seen test")
	s.Require().NoError(err)

	meta := s.waitForSubjectInFolder("INBOX", subject)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = s.client.SetSeen(ctx, meta.Id, "INBOX", true)
	cancel()
	s.Require().NoError(err)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err = s.client.SetSeen(ctx, meta.Id, "INBOX", false)
	cancel()
	s.Require().NoError(err)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	full, err := s.client.ReadEmail(ctx, meta.Id, "INBOX")
	cancel()
	s.Require().NoError(err)
	s.Require().Equal(subject, full.Subject)
}
