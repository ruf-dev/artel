-- +goose Up
-- Extends the gitlab MoM (seeded into mcp_tools by migration 030/032) so its tools are useful as
-- tract action steps: list_merge_requests gains source_branch/target_branch/labels filters,
-- create_merge_request switches from query params to a JSON body (plus assignee_id/
-- remove_source_branch), update_merge_request is added (title/assignee_id/reviewer_ids/
-- add_labels/remove_labels/state_event), a new MR-specific add_mr_comment tool is added
-- alongside the existing generic add_comment (left untouched), and
-- add_comment/list_merge_requests/create_merge_request/update_merge_request/add_mr_comment/
-- create_issue/list_issues/get_current_user/list_emails/read_email all declare an output_schema
-- so list_tract_actions can show step authors what a step produces.
--
-- UPSERT only (ON CONFLICT (mcp_name, name) DO UPDATE) — never deletes an existing tool row.

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'gitlab',
    'list_merge_requests',
    'List merge requests in a GitLab project',
    '{
        "properties": {
            "project_id":    { "type": "string", "description": "GitLab project ID or URL-encoded path" },
            "state":         { "type": "string", "description": "Filter by state", "enum": ["opened", "closed", "merged", "all"] },
            "source_branch": { "type": "string", "description": "Filter by source branch name (optional)" },
            "target_branch": { "type": "string", "description": "Only MRs targeting this branch" },
            "labels":        { "type": "string", "description": "Comma-separated label names to filter by (optional)" }
        },
        "required": ["project_id"]
    }'::jsonb,
    '{
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
    }'::jsonb,
    '{
        "http": {
            "method": "GET",
            "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/merge_requests",
            "headers": { "PRIVATE-TOKEN": "__secrets.personal_access_token" },
            "query": {
                "state": "${{params.state}}",
                "source_branch": "${{params.source_branch}}",
                "target_branch": "${{params.target_branch}}",
                "labels": "${{params.labels}}"
            },
            "credentials": "gitlab"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'gitlab',
    'create_merge_request',
    'Create a new merge request in a GitLab project',
    '{
        "properties": {
            "project_id":    { "type": "string", "description": "GitLab project ID or URL-encoded path (e.g. group%2Fproject)" },
            "source_branch": { "type": "string", "description": "Branch to merge from" },
            "target_branch": { "type": "string", "description": "Branch to merge into" },
            "title":         { "type": "string", "description": "Merge request title" },
            "description":   { "type": "string", "description": "Merge request description (optional)" },
            "assignee_id":          { "type": "integer", "description": "GitLab user id to assign" },
            "remove_source_branch": { "type": "boolean", "description": "Delete source branch on merge" }
        },
        "required": ["project_id", "source_branch", "target_branch", "title"]
    }'::jsonb,
    '{
        "properties": {
            "iid":     { "type": "integer", "description": "Merge request internal ID" },
            "web_url": { "type": "string" },
            "title":   { "type": "string" },
            "state":   { "type": "string", "enum": ["opened", "closed", "merged", "locked"] }
        },
        "required": ["iid", "web_url", "title", "state"]
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/merge_requests",
            "headers": {
                "PRIVATE-TOKEN": "__secrets.personal_access_token",
                "Content-Type": "application/json"
            },
            "body": {
                "source_branch": "${{params.source_branch}}",
                "target_branch": "${{params.target_branch}}",
                "title": "${{params.title}}",
                "description": "${{params.description}}",
                "assignee_id": "${{params.assignee_id}}",
                "remove_source_branch": "${{params.remove_source_branch}}"
            },
            "credentials": "gitlab"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'gitlab',
    'update_merge_request',
    'Update an existing merge request: assignee, reviewers, labels, title, state',
    '{
        "properties": {
            "project_id":    { "type": "string",  "description": "GitLab project ID or URL-encoded path" },
            "mr_iid":        { "type": "integer", "description": "Internal ID (iid) of the merge request" },
            "title":         { "type": "string",  "description": "New merge request title (optional)" },
            "assignee_id":   { "type": "integer", "description": "GitLab user id to assign (optional)" },
            "reviewer_ids":  { "type": "string",  "description": "Comma-separated GitLab user ids to set as reviewers (optional)" },
            "add_labels":    { "type": "string",  "description": "Comma-separated labels to add (optional)" },
            "remove_labels": { "type": "string",  "description": "Comma-separated labels to remove (optional)" },
            "state_event":   { "type": "string",  "description": "Close or reopen the MR (optional)", "enum": ["close", "reopen"] }
        },
        "required": ["project_id", "mr_iid"]
    }'::jsonb,
    '{
        "properties": {
            "iid":   { "type": "integer", "description": "Merge request internal ID" },
            "state": { "type": "string", "enum": ["opened", "closed", "merged", "locked"] }
        },
        "required": ["iid", "state"]
    }'::jsonb,
    '{
        "http": {
            "method": "PUT",
            "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/merge_requests/${{params.mr_iid}}",
            "headers": {
                "PRIVATE-TOKEN": "__secrets.personal_access_token",
                "Content-Type": "application/json"
            },
            "body": {
                "title":         "${{params.title}}",
                "assignee_id":   "${{params.assignee_id}}",
                "reviewer_ids":  "${{params.reviewer_ids}}",
                "add_labels":    "${{params.add_labels}}",
                "remove_labels": "${{params.remove_labels}}",
                "state_event":   "${{params.state_event}}"
            },
            "credentials": "gitlab"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'gitlab',
    'add_comment',
    'Add a comment (note) to a GitLab merge request or issue',
    '{
        "properties": {
            "project_id":    { "type": "string", "description": "GitLab project ID or URL-encoded path" },
            "noteable_type": { "type": "string", "description": "Type of item to comment on", "enum": ["merge_requests", "issues"] },
            "noteable_iid":  { "type": "string", "description": "Internal ID (iid) of the merge request or issue" },
            "body":          { "type": "string", "description": "Comment text" }
        },
        "required": ["project_id", "noteable_type", "noteable_iid", "body"]
    }'::jsonb,
    '{
        "properties": {
            "id":   { "type": "integer", "description": "Comment (note) ID" },
            "body": { "type": "string" }
        },
        "required": ["id", "body"]
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/${{params.noteable_type}}/${{params.noteable_iid}}/notes",
            "headers": { "PRIVATE-TOKEN": "__secrets.personal_access_token" },
            "query": { "body": "${{params.body}}" },
            "credentials": "gitlab"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

-- add_mr_comment is a NEW, MR-specific tool, separate from the generic add_comment above (which
-- covers both MRs and issues via noteable_type/noteable_iid and stays untouched — other
-- sessions/existing tracts may already reference it by name).
INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'gitlab',
    'add_mr_comment',
    'Post a comment (note) on a merge request',
    '{
        "properties": {
            "project_id": { "type": "string",  "description": "GitLab project ID or URL-encoded path" },
            "mr_iid":     { "type": "integer", "description": "MR number within the project" },
            "body":       { "type": "string",  "description": "Comment text (GitLab Markdown)" }
        },
        "required": ["project_id", "mr_iid", "body"]
    }'::jsonb,
    '{
        "properties": {
            "id":   { "type": "integer", "description": "Note id" },
            "body": { "type": "string" }
        },
        "required": ["id", "body"]
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/merge_requests/${{params.mr_iid}}/notes",
            "headers": {
                "PRIVATE-TOKEN": "__secrets.personal_access_token",
                "Content-Type": "application/json"
            },
            "body": { "body": "${{params.body}}" },
            "credentials": "gitlab"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

-- output_schema backfill for tools not already touched above (migration 032 backfilled every
-- tool's output_schema to '{}' — these upserts ONLY set output_schema, description/input_schema/
-- action are untouched since they're already correct as seeded by migration 030/031/027).
--
-- Array-shaped outputs follow list_merge_requests' own convention above: top-level
-- output_schema is domain.ToolSchema, which is only {properties, required} (see
-- internal/domain/mom.go) — there is no top-level "type"/"items" field, so an array response is
-- described with "properties" listing one element's fields directly, NOT wrapped in "items".
UPDATE mcp_tools
SET output_schema = '{
    "properties": {
        "iid":         { "type": "integer", "description": "Issue number within the project" },
        "id":          { "type": "integer" },
        "title":       { "type": "string" },
        "description": { "type": "string" },
        "state":       { "type": "string", "enum": ["opened", "closed"] },
        "web_url":     { "type": "string" }
    },
    "required": ["iid", "title", "state", "web_url"]
}'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'create_issue';

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
        "id":       { "type": "integer" },
        "username": { "type": "string", "description": "GitLab handle" },
        "name":     { "type": "string" },
        "email":    { "type": "string" }
    },
    "required": ["id", "username"]
}'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'get_current_user';

-- email/list_email_folders is intentionally SKIPPED here (no UPDATE statement for it): it
-- returns a bare []string (executeImap's IMAP_OP_LIST_FOLDERS case, marshaled directly) — a
-- top-level array of plain strings. domain.ToolSchema has no way to represent that (no top-level
-- "type"/"items", only "properties"/"required", which describe named object fields) so its
-- output_schema is left undeclared ('{}', its current value from migration 032's backfill)
-- rather than write a schema shape the code can't actually express.

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

UPDATE mcp_tools
SET output_schema = '{
    "properties": {
        "Id":      { "type": "string" },
        "From":    { "type": "string" },
        "To":      { "type": "string" },
        "Subject": { "type": "string" },
        "Date":    { "type": "string" },
        "Body":    { "type": "string" }
    },
    "required": ["Id", "From", "To", "Subject", "Date", "Body"]
}'::jsonb
WHERE mcp_name = 'email' AND name = 'read_email';

-- send_email (executeSmtp's SMTP_OP_SEND case) returns the literal string "Email sent
-- successfully", not JSON — describing a schema for a non-JSON string result would misrepresent
-- what the tool actually produces, so output_schema is left undeclared ('{}') here too.

-- +goose Down
DELETE FROM mcp_tools WHERE mcp_name = 'gitlab' AND name = 'update_merge_request';
DELETE FROM mcp_tools WHERE mcp_name = 'gitlab' AND name = 'add_mr_comment';

-- revert the Fix5 output_schema-only backfills to '{}' (their state after migration 032, before
-- this migration declared them) — description/input_schema/action for these tools were never
-- touched by this migration so nothing else needs reverting.
UPDATE mcp_tools SET output_schema = '{}'::jsonb
WHERE mcp_name = 'gitlab' AND name IN ('create_issue', 'list_issues', 'get_current_user');

UPDATE mcp_tools SET output_schema = '{}'::jsonb
WHERE mcp_name = 'email' AND name IN ('list_emails', 'read_email');

UPDATE mcp_tools
SET description   = 'List merge requests in a GitLab project',
    input_schema  = '{
        "properties": {
            "project_id": { "type": "string", "description": "GitLab project ID or URL-encoded path" },
            "state":      { "type": "string", "description": "Filter by state", "enum": ["opened", "closed", "merged", "all"] }
        },
        "required": ["project_id"]
    }'::jsonb,
    output_schema = '{}'::jsonb,
    action        = '{
        "http": {
            "method": "GET",
            "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/merge_requests",
            "headers": { "PRIVATE-TOKEN": "__secrets.personal_access_token" },
            "query": { "state": "${{params.state}}" },
            "credentials": "gitlab"
        }
    }'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'list_merge_requests';

UPDATE mcp_tools
SET description   = 'Create a new merge request in a GitLab project',
    input_schema  = '{
        "properties": {
            "project_id":    { "type": "string",  "description": "GitLab project ID or URL-encoded path (e.g. group%2Fproject)" },
            "source_branch": { "type": "string",  "description": "Branch to merge from" },
            "target_branch": { "type": "string",  "description": "Branch to merge into" },
            "title":         { "type": "string",  "description": "Merge request title" },
            "description":   { "type": "string",  "description": "Merge request description (optional)" }
        },
        "required": ["project_id", "source_branch", "target_branch", "title"]
    }'::jsonb,
    output_schema = '{}'::jsonb,
    action        = '{
        "http": {
            "method": "POST",
            "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/merge_requests",
            "headers": {
                "PRIVATE-TOKEN": "__secrets.personal_access_token",
                "Content-Type": "application/json"
            },
            "query": {
                "source_branch": "${{params.source_branch}}",
                "target_branch": "${{params.target_branch}}",
                "title": "${{params.title}}",
                "description": "${{params.description}}"
            },
            "credentials": "gitlab"
        }
    }'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'create_merge_request';

UPDATE mcp_tools
SET output_schema = '{}'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'add_comment';
