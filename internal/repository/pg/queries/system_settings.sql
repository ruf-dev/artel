-- name: GetSystemSettings :one
SELECT id, setup_completed, password_auth_enabled, telegram_auth_enabled, registration_mode,
       setup_token_hash, setup_token_issued_at, created_at, updated_at, default_docs_vault_id,
       default_docs_source, system_prompt
FROM system_settings
WHERE id = 1;

-- name: GetSystemSettingsForUpdate :one
SELECT id, setup_completed, password_auth_enabled, telegram_auth_enabled, registration_mode,
       setup_token_hash, setup_token_issued_at, created_at, updated_at, default_docs_vault_id,
       default_docs_source, system_prompt
FROM system_settings
WHERE id = 1
FOR UPDATE;

-- name: UpdateSystemSettingsAuthMethods :exec
UPDATE system_settings
SET password_auth_enabled = $1,
    telegram_auth_enabled = $2,
    updated_at            = NOW()
WHERE id = 1;

-- name: UpdateSystemSettingsRegistrationMode :exec
UPDATE system_settings
SET registration_mode = $1,
    updated_at         = NOW()
WHERE id = 1;

-- name: SetSetupToken :exec
UPDATE system_settings
SET setup_token_hash      = $1,
    setup_token_issued_at = $2,
    updated_at             = NOW()
WHERE id = 1;

-- name: CompleteSetup :exec
UPDATE system_settings
SET setup_completed  = TRUE,
    setup_token_hash = NULL,
    updated_at       = NOW()
WHERE id = 1;

-- name: UpdateSystemSettingsDefaultDocsVault :exec
UPDATE system_settings
SET default_docs_vault_id = $1,
    updated_at             = NOW()
WHERE id = 1;

-- name: UpdateSystemSettingsDefaultDocsSource :exec
UPDATE system_settings SET default_docs_source = $1, updated_at = NOW() WHERE id = 1;

-- name: UpdateSystemPrompt :exec
UPDATE system_settings SET system_prompt = $1, updated_at = NOW() WHERE id = 1;
