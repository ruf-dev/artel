-- name: GetUserSettings :one
SELECT u.id AS user_id,
       COALESCE(us.user_prompt, ''::text) AS user_prompt,
       COALESCE(us.liked_openrouter_models, '{}'::text[]) AS liked_openrouter_models,
       COALESCE(us.last_used_model, ''::text) AS last_used_model,
       NOW() AS created_at,
       NOW() AS updated_at
FROM users u
LEFT JOIN user_settings us ON us.user_id = u.id
WHERE u.id = $1;

-- name: UpsertLikedModels :exec
INSERT INTO user_settings (user_id, liked_openrouter_models)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
    SET liked_openrouter_models = EXCLUDED.liked_openrouter_models,
        updated_at              = NOW();

-- name: UpsertLastUsedModel :exec
INSERT INTO user_settings (user_id, last_used_model)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
    SET last_used_model = EXCLUDED.last_used_model,
        updated_at       = NOW();
