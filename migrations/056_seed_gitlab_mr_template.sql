-- +goose Up
-- Seeds the "Create MR on featureprep push" built-in tract template (owner_id NULL, see
-- migration 055) so it's available out-of-the-box on every fresh install, not just the account
-- it was originally published from. Fixed, well-known id so this insert is idempotent and the
-- row is addressable across environments.
INSERT INTO tract_templates (id, owner_id, name, description, definition, category) VALUES (
    'b8b6a2a0-1c1a-4b9e-9f4a-6f3a2b1c0d55',
    NULL,
    'Create MR on featurep push',
    '',
    '{"steps": [
        {"id": "list_merge_requests", "mcp": "gitlab", "name": "list_merge_requests", "tool": "list_merge_requests", "type": "action", "params": {"state": "opened", "labels": "Actual release", "project_id": "{{ trigger.project.id }}", "target_branch": "master"}, "connection_uuid": "00000000-0000-0000-0000-000000000000"},
        {"id": "script", "code": "target_branch = merge_requests[0].source_branch", "name": "script", "type": "script", "params": {"merge_requests": "{{ list_merge_requests }}"}, "language": "javascript", "input_params": [{"Name": "merge_requests", "Property": {"Enum": null, "Type": "array", "Items": null, "Required": null, "Properties": null, "Description": ""}}], "output_params": [{"Name": "target_branch", "Property": {"Enum": null, "Type": "string", "Items": null, "Required": null, "Properties": null, "Description": ""}}], "connection_uuid": "00000000-0000-0000-0000-000000000000"},
        {"id": "create_merge_request", "mcp": "gitlab", "name": "create_merge_request", "tool": "create_merge_request", "type": "action", "params": {"title": "{{ trigger.commits[0].message }}", "project_id": "{{ trigger.project.id }}", "source_branch": "{{ trigger.branch }}", "target_branch": "{{ script.target_branch }}"}, "connection_uuid": "00000000-0000-0000-0000-000000000000"}
    ]}',
    'Mr on push'
) ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM tract_templates WHERE id = 'b8b6a2a0-1c1a-4b9e-9f4a-6f3a2b1c0d55';
