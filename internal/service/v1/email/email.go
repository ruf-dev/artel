package email

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/imap"
	"github.com/ruf-dev/artel/internal/clients/smtp"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
)

const defaultListLimit = 20

type EmailServiceImpl struct {
	accounts repository.EmailAccounts
}

func New(accounts repository.EmailAccounts) *EmailServiceImpl {
	return &EmailServiceImpl{accounts: accounts}
}

func (s *EmailServiceImpl) AddAccount(ctx context.Context, account domain.EmailAccount) (domain.EmailAccount, error) {
	result, err := s.accounts.Insert(ctx, account)
	if err != nil {
		return domain.EmailAccount{}, rerrors.Wrap(err, "add email account")
	}
	return result, nil
}

func (s *EmailServiceImpl) ListAccounts(ctx context.Context, userUuid uuid.UUID) ([]domain.EmailAccount, error) {
	accounts, err := s.accounts.ListByUser(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "list email accounts")
	}
	return accounts, nil
}

func (s *EmailServiceImpl) DeleteAccount(ctx context.Context, accountUuid uuid.UUID) error {
	if err := s.accounts.Delete(ctx, accountUuid); err != nil {
		return rerrors.Wrap(err, "delete email account")
	}
	return nil
}

func (s *EmailServiceImpl) ListFolders(ctx context.Context, accountUuid uuid.UUID) ([]string, error) {
	account, err := s.accounts.GetByUuid(ctx, accountUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "get email account")
	}
	c := imap.New(account.ImapHost, account.ImapPort, account.Email, account.Password)
	return c.ListFolders(ctx)
}

func (s *EmailServiceImpl) ListEmails(ctx context.Context, accountUuid uuid.UUID, limit int) ([]domain.EmailMeta, error) {
	account, err := s.accounts.GetByUuid(ctx, accountUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "get email account")
	}

	if limit <= 0 {
		limit = defaultListLimit
	}

	c := imap.New(account.ImapHost, account.ImapPort, account.Email, account.Password)
	emails, err := c.ListEmails(ctx, limit)
	if err != nil {
		return nil, rerrors.Wrap(err, "list emails")
	}
	return emails, nil
}

func (s *EmailServiceImpl) ReadEmail(ctx context.Context, accountUuid uuid.UUID, id string) (domain.EmailMessage, error) {
	account, err := s.accounts.GetByUuid(ctx, accountUuid)
	if err != nil {
		return domain.EmailMessage{}, rerrors.Wrap(err, "get email account")
	}

	c := imap.New(account.ImapHost, account.ImapPort, account.Email, account.Password)
	msg, err := c.ReadEmail(ctx, id)
	if err != nil {
		return domain.EmailMessage{}, rerrors.Wrap(err, "read email")
	}
	return msg, nil
}

func (s *EmailServiceImpl) SendEmail(ctx context.Context, accountUuid uuid.UUID, to, subject, body string) error {
	account, err := s.accounts.GetByUuid(ctx, accountUuid)
	if err != nil {
		return rerrors.Wrap(err, "get email account")
	}

	c := smtp.New(account.SmtpHost, account.SmtpPort, account.Email, account.Password)
	if err := c.Send(ctx, to, subject, body); err != nil {
		return rerrors.Wrap(err, "send email")
	}
	return nil
}
