-- +goose Up
-- domain.ToolSchema gained a real top-level array capability (IsArray + Items, mirroring how
-- ToolProperty already models a nested array via type=="array"+items — see internal/domain/mom.go)
-- so a tool's output_schema can now say "the whole response is an array of X" instead of the old
-- workaround of listing the array element's fields directly under "properties" as if they were
-- top-level object properties. That workaround meant the tract template autocomplete
-- (pkg/client/ArtelUI's TemplateInput) never offered the "[0]" index these tools' outputs actually
-- need.
--
-- This migration converts every array-shaped mcp_tools.output_schema seeded so far from that
-- flattened-properties shape to the new {"isArray": true, "items": {...}} shape, field lists
-- preserved verbatim:
--   - gitlab/list_merge_requests, gitlab/list_issues     (migration 034)
--   - email/list_emails                                  (migration 034)
--   - trello/list_boards, list_lists, list_cards, list_labels (migration 046)
--
-- It also declares output_schema for three tools that were previously left undeclared ('{}')
-- specifically because the old schema type had no way to express a primitive or nested-object
-- array at the root (see the comments at migration 034 lines 292-298, migration 046 lines
-- 217-219, and migration 049 lines 53-56) — that limitation no longer applies:
--   - email/list_email_folders   (bare []string of folder names)
--   - trello/add_label_to_card   (Trello's raw array of the card's label id strings)
--   - trello/list_card_comments  (array of Trello action objects)
--
-- UPDATE-only — only touches output_schema on existing rows, never edits description/
-- input_schema/action, never inserts/deletes rows. Does not touch create_issue,
-- get_current_user, read_email, get_card, move_card, create_card, archive_card — those return
-- single objects and are already correctly modeled.

UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": {
        "type": "object",
        "properties": {
            "iid":            { "type": "integer", "description": "Merge request internal ID" },
            "title":          { "type": "string" },
            "web_url":        { "type": "string" },
            "state":          { "type": "string", "enum": ["opened", "closed", "merged", "locked"] },
            "source_branch":  { "type": "string" },
            "target_branch":  { "type": "string" },
            "labels":         { "type": "array", "items": { "type": "string" } }
        },
        "required": ["iid", "title", "web_url", "state"]
    }
}'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'list_merge_requests';

UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": {
        "type": "object",
        "properties": {
            "iid":     { "type": "integer", "description": "Issue number within the project" },
            "id":      { "type": "integer" },
            "title":   { "type": "string" },
            "state":   { "type": "string", "enum": ["opened", "closed"] },
            "web_url": { "type": "string" },
            "labels":  { "type": "array", "items": { "type": "string" } }
        },
        "required": ["iid", "title", "state", "web_url"]
    }
}'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'list_issues';

UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": {
        "type": "object",
        "properties": {
            "Id":      { "type": "string" },
            "From":    { "type": "string" },
            "Subject": { "type": "string" },
            "Date":    { "type": "string" }
        },
        "required": ["Id", "From", "Subject", "Date"]
    }
}'::jsonb
WHERE mcp_name = 'email' AND name = 'list_emails';

UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": { "type": "string", "description": "Folder name" }
}'::jsonb
WHERE mcp_name = 'email' AND name = 'list_email_folders';

UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": {
        "type": "object",
        "properties": {
            "id":     { "type": "string" },
            "name":   { "type": "string" },
            "desc":   { "type": "string" },
            "closed": { "type": "boolean" },
            "url":    { "type": "string" }
        },
        "required": ["id", "name"]
    }
}'::jsonb
WHERE mcp_name = 'trello' AND name = 'list_boards';

UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": {
        "type": "object",
        "properties": {
            "id":      { "type": "string" },
            "name":    { "type": "string" },
            "closed":  { "type": "boolean" },
            "idBoard": { "type": "string" }
        },
        "required": ["id", "name"]
    }
}'::jsonb
WHERE mcp_name = 'trello' AND name = 'list_lists';

UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": {
        "type": "object",
        "properties": {
            "id":      { "type": "string" },
            "name":    { "type": "string" },
            "desc":    { "type": "string" },
            "idList":  { "type": "string", "description": "Id of the list (column) this card is in" },
            "idBoard": { "type": "string" },
            "closed":  { "type": "boolean" },
            "url":     { "type": "string" },
            "labels":  { "type": "array", "items": { "type": "string" } }
        },
        "required": ["id", "name", "idList"]
    }
}'::jsonb
WHERE mcp_name = 'trello' AND name = 'list_cards';

UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": {
        "type": "object",
        "properties": {
            "id":    { "type": "string" },
            "name":  { "type": "string" },
            "color": { "type": "string" }
        },
        "required": ["id", "color"]
    }
}'::jsonb
WHERE mcp_name = 'trello' AND name = 'list_labels';

UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": { "type": "string", "description": "Label id" }
}'::jsonb
WHERE mcp_name = 'trello' AND name = 'add_label_to_card';

UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": {
        "type": "object",
        "properties": {
            "id":            { "type": "string" },
            "data":          { "type": "object", "properties": { "text": { "type": "string" } } },
            "date":          { "type": "string" },
            "memberCreator": { "type": "object", "properties": { "fullName": { "type": "string" } } }
        },
        "required": ["id", "date"]
    }
}'::jsonb
WHERE mcp_name = 'trello' AND name = 'list_card_comments';

-- +goose Down
UPDATE mcp_tools
SET output_schema = '{
    "properties": {
        "iid":            { "type": "integer", "description": "Merge request internal ID" },
        "title":          { "type": "string" },
        "web_url":        { "type": "string" },
        "state":          { "type": "string", "enum": ["opened", "closed", "merged", "locked"] },
        "source_branch":  { "type": "string" },
        "target_branch":  { "type": "string" },
        "labels":         { "type": "array", "items": { "type": "string" } }
    },
    "required": ["iid", "title", "web_url", "state"]
}'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'list_merge_requests';

UPDATE mcp_tools
SET output_schema = '{
    "properties": {
        "iid":     { "type": "integer", "description": "Issue number within the project" },
        "id":      { "type": "integer" },
        "title":   { "type": "string" },
        "state":   { "type": "string", "enum": ["opened", "closed"] },
        "web_url": { "type": "string" },
        "labels":  { "type": "array", "items": { "type": "string" } }
    },
    "required": ["iid", "title", "state", "web_url"]
}'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'list_issues';

UPDATE mcp_tools
SET output_schema = '{
    "properties": {
        "Id":      { "type": "string" },
        "From":    { "type": "string" },
        "Subject": { "type": "string" },
        "Date":    { "type": "string" }
    },
    "required": ["Id", "From", "Subject", "Date"]
}'::jsonb
WHERE mcp_name = 'email' AND name = 'list_emails';

UPDATE mcp_tools SET output_schema = '{}'::jsonb
WHERE mcp_name = 'email' AND name = 'list_email_folders';

UPDATE mcp_tools
SET output_schema = '{
    "properties": {
        "id":     { "type": "string" },
        "name":   { "type": "string" },
        "desc":   { "type": "string" },
        "closed": { "type": "boolean" },
        "url":    { "type": "string" }
    },
    "required": ["id", "name"]
}'::jsonb
WHERE mcp_name = 'trello' AND name = 'list_boards';

UPDATE mcp_tools
SET output_schema = '{
    "properties": {
        "id":      { "type": "string" },
        "name":    { "type": "string" },
        "closed":  { "type": "boolean" },
        "idBoard": { "type": "string" }
    },
    "required": ["id", "name"]
}'::jsonb
WHERE mcp_name = 'trello' AND name = 'list_lists';

UPDATE mcp_tools
SET output_schema = '{
    "properties": {
        "id":      { "type": "string" },
        "name":    { "type": "string" },
        "desc":    { "type": "string" },
        "idList":  { "type": "string", "description": "Id of the list (column) this card is in" },
        "idBoard": { "type": "string" },
        "closed":  { "type": "boolean" },
        "url":     { "type": "string" },
        "labels":  { "type": "array", "items": { "type": "string" } }
    },
    "required": ["id", "name", "idList"]
}'::jsonb
WHERE mcp_name = 'trello' AND name = 'list_cards';

UPDATE mcp_tools
SET output_schema = '{
    "properties": {
        "id":    { "type": "string" },
        "name":  { "type": "string" },
        "color": { "type": "string" }
    },
    "required": ["id", "color"]
}'::jsonb
WHERE mcp_name = 'trello' AND name = 'list_labels';

UPDATE mcp_tools SET output_schema = '{}'::jsonb
WHERE mcp_name = 'trello' AND name = 'add_label_to_card';

UPDATE mcp_tools SET output_schema = '{}'::jsonb
WHERE mcp_name = 'trello' AND name = 'list_card_comments';
