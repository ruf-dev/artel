-- +goose Up

-- Adds the second GitLab trigger preset ("MR merged") plus the two GitLab tool gaps needed for
-- a tract to read an MR's diff and write a description back onto it (see plan Part A):
--   1. trigger_presets row: gitlab_merge_request.
--   2. mcp_tools row: get_merge_request_diff (new).
--   3. mcp_tools row: update_merge_request (extended with a "description" property).
--
-- GitLab sends X-Gitlab-Event: "Merge Request Hook" for opened/updated/merged/closed alike — the
-- header alone can't distinguish "merged" from other MR actions. Rather than push that
-- distinction onto every tract author via a condition step, this preset uses the new
-- domain.TriggerMatchers.CheckBody matcher (mirrors CheckHeaders' shape, added alongside this
-- migration in internal/service/v1/tract/matchers.go) to also require
-- object_attributes.action == "merge" at dispatch time, so the preset is self-describing and
-- correct out of the box.

INSERT INTO trigger_presets (key, category, label, description, provider, payload_schema, default_matchers) VALUES
('gitlab_merge_request', 'gitlab', 'Merge request merged',
 'Fires when a merge request is merged',
 'gitlab',
 '{"properties": {
      "object_attributes": {"type": "object", "properties": {
          "iid":            {"type": "integer", "description": "merge request internal ID (project-scoped)"},
          "title":           {"type": "string"},
          "description":     {"type": "string"},
          "source_branch":   {"type": "string"},
          "target_branch":   {"type": "string"},
          "action":          {"type": "string", "description": "open | close | reopen | update | approved | unapproved | merge", "enum": ["open", "close", "reopen", "update", "approved", "unapproved", "merge"]},
          "merge_status":    {"type": "string"},
          "url":             {"type": "string"}
      }},
      "mr_iid": {"type": "integer", "description": "object_attributes.iid, surfaced at the top level by the normalizer"},
      "action": {"type": "string", "description": "object_attributes.action, surfaced at the top level by the normalizer"},
      "project": {"type": "object", "properties": {
          "id": {"type": "integer"}, "name": {"type": "string"}, "path": {"type": "string"}
      }},
      "user": {"type": "object", "properties": {
          "name": {"type": "string"}, "username": {"type": "string"}
      }}
  }, "required": ["object_attributes", "mr_iid", "action"]}'::jsonb,
 '{"check_headers": [{"header": "X-Gitlab-Event", "equals": "Merge Request Hook"}],
   "check_body": [{"path": "object_attributes.action", "equals": "merge"}]}'::jsonb
);

-- get_merge_request_diff — GitLab REST API v4's dedicated MR diffs endpoint (GET
-- .../merge_requests/:merge_request_iid/diffs), returns the array of per-file diffs.
INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'gitlab',
    'get_merge_request_diff',
    'Get the diff (per-file changes) of a merge request',
    '{
        "properties": {
            "project_id": { "type": "string",  "description": "GitLab project ID or URL-encoded path" },
            "mr_iid":     { "type": "integer", "description": "Internal ID (iid) of the merge request" }
        },
        "required": ["project_id", "mr_iid"]
    }'::jsonb,
    '{
        "properties": {
            "old_path":         { "type": "string" },
            "new_path":         { "type": "string" },
            "a_mode":           { "type": "string" },
            "b_mode":           { "type": "string" },
            "diff":             { "type": "string", "description": "unified diff text for this file" },
            "new_file":         { "type": "boolean" },
            "renamed_file":     { "type": "boolean" },
            "deleted_file":     { "type": "boolean" }
        },
        "required": ["old_path", "new_path", "diff"]
    }'::jsonb,
    '{
        "http": {
            "method": "GET",
            "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/merge_requests/${{params.mr_iid}}/diffs",
            "headers": { "PRIVATE-TOKEN": "__secrets.personal_access_token" },
            "credentials": "gitlab"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

-- update_merge_request: extend with "description" so a tract can write a diff summary back onto
-- the MR. Every other property (title, assignee_id, reviewer_ids, add_labels, remove_labels,
-- state_event) is carried over unchanged from migration 034_gitlab_tract_tools.sql's row.
INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'gitlab',
    'update_merge_request',
    'Update an existing merge request: assignee, reviewers, labels, title, description, state',
    '{
        "properties": {
            "project_id":    { "type": "string",  "description": "GitLab project ID or URL-encoded path" },
            "mr_iid":        { "type": "integer", "description": "Internal ID (iid) of the merge request" },
            "title":         { "type": "string",  "description": "New merge request title (optional)" },
            "description":   { "type": "string",  "description": "New merge request description (optional)" },
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
                "description":   "${{params.description}}",
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

-- +goose Down
DELETE FROM mcp_tools WHERE mcp_name = 'gitlab' AND name = 'get_merge_request_diff';

-- revert update_merge_request to its pre-057 shape (migration 034_gitlab_tract_tools.sql's row,
-- no "description" property).
UPDATE mcp_tools
SET description   = 'Update an existing merge request: assignee, reviewers, labels, title, state',
    input_schema  = '{
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
    output_schema = '{
        "properties": {
            "iid":   { "type": "integer", "description": "Merge request internal ID" },
            "state": { "type": "string", "enum": ["opened", "closed", "merged", "locked"] }
        },
        "required": ["iid", "state"]
    }'::jsonb,
    action        = '{
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
WHERE mcp_name = 'gitlab' AND name = 'update_merge_request';

DELETE FROM trigger_presets WHERE key = 'gitlab_merge_request';
