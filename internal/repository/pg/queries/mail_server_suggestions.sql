-- name: ListMailServerSuggestions :many
SELECT domain, smtp, smtp_port, imap, imap_port, app_password_url
FROM mail_server_suggestions
WHERE domain LIKE $1 || '%'
ORDER BY domain;
