-- name: GetUserSettings :one
SELECT u.id AS user_id,
       COALESCE(us.user_prompt, ''::text) AS user_prompt,
       NOW() AS created_at,
       NOW() AS updated_at
FROM users u
LEFT JOIN user_settings us ON us.user_id = u.id
WHERE u.id = $1;
