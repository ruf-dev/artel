-- name: InsertSimpleChatMessage :one
INSERT INTO simple_chat_messages (chat_id, role, content, tool_call_id, tool_name, tool_input, is_error, model, seq)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, chat_id, role, content, tool_call_id, tool_name, tool_input, is_error, model, seq, created_at;

-- name: ListSimpleChatMessagesByChatID :many
SELECT id, chat_id, role, content, tool_call_id, tool_name, tool_input, is_error, model, seq, created_at
FROM simple_chat_messages
WHERE chat_id = $1
ORDER BY seq ASC;

-- name: GetMaxSeqForChat :one
SELECT COALESCE(MAX(seq), 0)::bigint AS max_seq
FROM simple_chat_messages
WHERE chat_id = $1;
