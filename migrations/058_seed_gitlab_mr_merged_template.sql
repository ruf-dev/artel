-- +goose Up

-- Fixes get_merge_request_diff's output_schema (added in migration 057) to use the isArray/items
-- shape established by migration 051_tool_output_array_schemas.sql for every other array-returning
-- GitLab tool (list_merge_requests, list_issues) — GitLab's diffs endpoint returns an array of
-- per-file diff objects, not a single object, so it should have been declared the same way as
-- those siblings from the start. UPDATE-only, mirrors 051's style.
UPDATE mcp_tools
SET output_schema = '{
    "isArray": true,
    "items": {
        "type": "object",
        "properties": {
            "old_path":     { "type": "string" },
            "new_path":     { "type": "string" },
            "a_mode":       { "type": "string" },
            "b_mode":       { "type": "string" },
            "diff":         { "type": "string", "description": "unified diff text for this file" },
            "new_file":     { "type": "boolean" },
            "renamed_file": { "type": "boolean" },
            "deleted_file": { "type": "boolean" }
        },
        "required": ["old_path", "new_path", "diff"]
    }
}'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'get_merge_request_diff';

-- Seeds the "Describe MR on merge" built-in tract template (owner_id NULL, see migration 055) —
-- a working reference for the plan's second flow: gitlab_merge_request trigger (migration 057)
-- -> get_merge_request_diff -> a script step joining the per-file diffs into one text blob ->
-- an llm_call step summarizing the diff -> update_merge_request writing the summary back onto
-- the MR's description. Mirrors migration 056_seed_gitlab_mr_template.sql's shape: action steps'
-- connection_uuid and the llm_call step's llm_connection_uuid are the all-zero placeholder,
-- filled in by InstantiateTemplate from the caller's own connections (gitlab mcp name +
-- "anthropic" provider name, see templates.go's walkActions/walkLlmCallStepsMut). Fixed,
-- well-known id so this insert is idempotent across environments.
INSERT INTO tract_templates (id, owner_id, name, description, definition, category) VALUES (
    'c1d7b3b1-2d2b-4c0f-a05b-7f4b3c2d1e66',
    NULL,
    'Describe MR on merge',
    'On merge, summarizes the diff with an LLM and writes it back as the MR description.',
    '{"steps": [
        {"id": "get_diff", "mcp": "gitlab", "name": "get_merge_request_diff", "tool": "get_merge_request_diff", "type": "action", "params": {"mr_iid": "{{ trigger.mr_iid }}", "project_id": "{{ trigger.project.id }}"}, "connection_uuid": "00000000-0000-0000-0000-000000000000"},
        {"id": "format_diff", "code": "diff_text = diffs.map(function(d) { return \"--- \" + d.old_path + \"\\n+++ \" + d.new_path + \"\\n\" + d.diff; }).join(\"\\n\\n\")", "name": "format_diff", "type": "script", "params": {"diffs": "{{ get_diff }}"}, "language": "javascript", "input_params": [{"Name": "diffs", "Property": {"Enum": null, "Type": "array", "Items": null, "Required": null, "Properties": null, "Description": ""}}], "output_params": [{"Name": "diff_text", "Property": {"Enum": null, "Type": "string", "Items": null, "Required": null, "Properties": null, "Description": ""}}], "connection_uuid": "00000000-0000-0000-0000-000000000000"},
        {"id": "summarize", "name": "summarize", "type": "llm_call", "llm_model": "claude-opus-4-8", "llm_connection_uuid": "00000000-0000-0000-0000-000000000000", "system_prompt": "You are a helpful assistant that writes clear, concise merge request descriptions from code diffs. Respond with only the description text, no preamble.", "prompt": "Summarize the work done in this merge request based on the diff below. Write a short, professional description suitable for the merge request description field.\n\nTitle: {{ trigger.object_attributes.title }}\n\nDiff:\n{{ format_diff.diff_text }}"},
        {"id": "update_description", "mcp": "gitlab", "name": "update_merge_request", "tool": "update_merge_request", "type": "action", "params": {"mr_iid": "{{ trigger.mr_iid }}", "project_id": "{{ trigger.project.id }}", "description": "{{ summarize.text }}"}, "connection_uuid": "00000000-0000-0000-0000-000000000000"}
    ]}',
    'Mr on merge'
) ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM tract_templates WHERE id = 'c1d7b3b1-2d2b-4c0f-a05b-7f4b3c2d1e66';

UPDATE mcp_tools
SET output_schema = '{
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
}'::jsonb
WHERE mcp_name = 'gitlab' AND name = 'get_merge_request_diff';
