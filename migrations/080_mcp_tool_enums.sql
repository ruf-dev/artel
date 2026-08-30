-- +goose Up
-- Two Postgres ENUM types that name the currently-seeded MoM tool-sets (mcp_name) and their
-- individual tools (mcp_tool_name, formatted "<mcp>.<tool>"). Neither enum is attached to any
-- column: the existing mcps.name PK and mcp_tools.mcp_name FK stay TEXT (retyping them is
-- invasive and fights migration reversibility). The types exist to be consumed through their
-- sqlc-generated Go types (artel_q.McpName / artel_q.McpToolName) — the values a caller names
-- when scoping per-mcp / per-tool HTTP middleware. sqlc emits an enum type from a bare CREATE
-- TYPE with no backing column; mcp_tool (singular) would collide with the mcp_tools table's row
-- struct, hence mcp_tool_name.
--
-- CREATE TYPE is transaction-safe, so this migration stays transactional. A later migration that
-- needs to add a label uses ALTER TYPE ... ADD VALUE, which is NOT transaction-safe (see the
-- trello write-tools migration).
CREATE TYPE mcp_name AS ENUM (
    'email',
    'gitlab',
    'telegram',
    'trello'
);

CREATE TYPE mcp_tool_name AS ENUM (
    'email.list_email_folders',
    'email.list_emails',
    'email.read_email',
    'email.send_email',
    'gitlab.add_comment',
    'gitlab.add_mr_comment',
    'gitlab.create_issue',
    'gitlab.create_merge_request',
    'gitlab.get_current_user',
    'gitlab.get_merge_request_diff',
    'gitlab.list_issues',
    'gitlab.list_merge_requests',
    'gitlab.update_merge_request',
    'telegram.answer_callback_query',
    'telegram.edit_message_text',
    'telegram.get_me',
    'telegram.send_message',
    'telegram.set_webhook',
    'trello.add_comment',
    'trello.add_label_to_card',
    'trello.archive_card',
    'trello.create_card',
    'trello.get_card',
    'trello.get_current_user',
    'trello.list_boards',
    'trello.list_card_comments',
    'trello.list_cards',
    'trello.list_labels',
    'trello.list_lists',
    'trello.move_card'
);

-- +goose Down
DROP TYPE mcp_tool_name;
DROP TYPE mcp_name;
