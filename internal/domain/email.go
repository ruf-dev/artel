package domain

type MailServerSuggestion struct {
	Domain         string
	Smtp           string
	SmtpPort       int
	Imap           string
	ImapPort       int
	AppPasswordUrl string
}

type EmailMeta struct {
	Id      string
	From    string
	Subject string
	Date    string
}

type EmailMessage struct {
	Id      string
	From    string
	To      string
	Subject string
	Date    string
	Body    string
}
