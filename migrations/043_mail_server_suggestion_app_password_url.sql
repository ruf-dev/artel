-- +goose Up

ALTER TABLE mail_server_suggestions ADD COLUMN app_password_url TEXT NULL;

UPDATE mail_server_suggestions
SET app_password_url = 'https://id.yandex.ru/security/app-passwords'
WHERE domain = 'yandex.ru';

INSERT INTO mail_server_suggestions (domain, smtp, smtp_port, imap, imap_port, app_password_url)
SELECT 'yandex.com', smtp, smtp_port, imap, imap_port, app_password_url
FROM mail_server_suggestions WHERE domain = 'yandex.ru'
UNION ALL
SELECT 'ya.ru', smtp, smtp_port, imap, imap_port, app_password_url
FROM mail_server_suggestions WHERE domain = 'yandex.ru';

-- +goose Down

DELETE FROM mail_server_suggestions WHERE domain IN ('yandex.com', 'ya.ru');
ALTER TABLE mail_server_suggestions DROP COLUMN app_password_url;
