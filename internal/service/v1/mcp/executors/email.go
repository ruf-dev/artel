package executors

import (
	"context"
	"encoding/json"

	"github.com/ruf-dev/artel/internal/clients/imap"
	"github.com/ruf-dev/artel/internal/clients/smtp"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

const defaultEmailListLimit = 20

type EmailExecutor struct{}

func NewEmailExecutor() *EmailExecutor {
	return &EmailExecutor{}
}

func (e *EmailExecutor) Execute(
	ctx context.Context,
	action domain.ToolAction,
	creds domain.EmailCredentials,
	params map[string]interface{},
) (string, error) {
	switch {
	case action.Imap != nil:
		return e.executeImap(ctx, action.Imap.Operation, creds, params)
	case action.Smtp != nil:
		return e.executeSmtp(ctx, action.Smtp.Operation, creds, params)
	default:
		return "", user_errors.McpEmailActionMissing
	}
}

func (e *EmailExecutor) executeImap(
	ctx context.Context,
	op domain.ImapOperation,
	creds domain.EmailCredentials,
	params map[string]interface{},
) (string, error) {
	imapClient := imap.New(creds.ImapHost, creds.ImapPort, creds.Username, creds.Password)

	switch op {
	case domain.IMAP_OP_LIST_FOLDERS:
		folders, err := imapClient.ListFolders(ctx)
		if err != nil {
			return "", rerrors.Wrap(err, "list imap folders")
		}

		return marshalResult(folders)

	case domain.IMAP_OP_LIST_MESSAGES:
		limit := defaultEmailListLimit

		if raw, ok := params["limit"]; ok {
			switch v := raw.(type) {
			case float64:
				limit = int(v)
			case int:
				limit = v
			}
		}

		afterUid, err := optionalUidStringParam(params, "after_uid")
		if err != nil {
			return "", err
		}

		beforeUid, err := optionalUidStringParam(params, "before_uid")
		if err != nil {
			return "", err
		}

		if afterUid != nil && beforeUid != nil {
			return "", user_errors.EmailCursorConflict
		}

		opts := imap.ListEmailsOptions{
			Limit:     limit,
			AfterUid:  afterUid,
			BeforeUid: beforeUid,
		}

		emails, err := imapClient.ListEmails(ctx, opts)
		if err != nil {
			return "", rerrors.Wrap(err, "list emails")
		}

		return marshalResult(emails)

	case domain.IMAP_OP_FETCH_MESSAGE:
		id, ok := params["id"].(string)
		if !ok {
			return "", user_errors.McpIdRequired
		}

		msg, err := imapClient.ReadEmail(ctx, id)
		if err != nil {
			return "", rerrors.Wrap(err, "read email")
		}

		return marshalResult(msg)

	default:
		return "", user_errors.McpUnknownImapOperation
	}
}

func (e *EmailExecutor) executeSmtp(
	ctx context.Context,
	op domain.SmtpOperation,
	creds domain.EmailCredentials,
	params map[string]interface{},
) (string, error) {
	switch op {
	case domain.SMTP_OP_SEND:
		to, ok := params["to"].(string)
		if !ok {
			return "", user_errors.McpToRequired
		}

		subject, ok := params["subject"].(string)
		if !ok {
			return "", user_errors.McpSubjectRequired
		}

		body, ok := params["body"].(string)
		if !ok {
			return "", user_errors.McpBodyRequired
		}

		c := smtp.New(creds.SmtpHost, creds.SmtpPort, creds.Username, creds.Password)

		err := c.Send(ctx, to, subject, body)
		if err != nil {
			return "", rerrors.Wrap(err, "send email")
		}

		return "Email sent successfully", nil

	default:
		return "", user_errors.McpUnknownSmtpOperation
	}
}

// optionalUidStringParam reads an optional UID-string param (after_uid/before_uid), returning
// (nil, nil) if absent. A present-but-wrong-type value is rejected rather than silently ignored,
// since silently dropping a cursor would produce confusing pagination behavior. Numeric-format
// validation of the string itself happens inside imap.Client.ListEmails, matching how the id
// param is validated for IMAP_OP_FETCH_MESSAGE today.
func optionalUidStringParam(params map[string]interface{}, key string) (*string, error) {
	raw, ok := params[key]
	if !ok {
		return nil, nil
	}

	str, ok := raw.(string)
	if !ok {
		return nil, user_errors.InvalidEmailId
	}

	return &str, nil
}

func marshalResult(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", rerrors.Wrap(err, "marshal executor result")
	}

	return string(data), nil
}
