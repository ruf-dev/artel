-- +goose Up

CREATE TABLE mail_server_suggestions (
    domain    TEXT PRIMARY KEY,
    smtp      TEXT NOT NULL,
    smtp_port INT  NOT NULL,
    imap      TEXT NOT NULL,
    imap_port INT  NOT NULL
);

INSERT INTO mail_server_suggestions (domain, smtp, smtp_port, imap, imap_port) VALUES
    ('gmail.com',   'smtp.gmail.com',           587, 'imap.gmail.com',          993),
    ('yahoo.com',   'smtp.mail.yahoo.com',       587, 'imap.mail.yahoo.com',     993),
    ('outlook.com', 'smtp-mail.outlook.com',     587, 'outlook.office365.com',   993),
    ('hotmail.com', 'smtp-mail.outlook.com',     587, 'outlook.office365.com',   993),
    ('icloud.com',  'smtp.mail.me.com',          587, 'imap.mail.me.com',        993),
    ('mail.ru',     'smtp.mail.ru',              465, 'imap.mail.ru',            993),
    ('yandex.ru',   'smtp.yandex.ru',            465, 'imap.yandex.ru',          993);

-- +goose Down

DROP TABLE IF EXISTS mail_server_suggestions;
