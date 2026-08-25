-- name: UpsertSimpleChatToolAllowance :exec
INSERT INTO simple_chat_tool_allowances (chat_id, tool_name, decision)
VALUES ($1, $2, $3)
ON CONFLICT (chat_id, tool_name) DO UPDATE SET decision = EXCLUDED.decision;

-- name: GetSimpleChatToolAllowance :one
SELECT decision
FROM simple_chat_tool_allowances
WHERE chat_id = $1
  AND tool_name = $2;
