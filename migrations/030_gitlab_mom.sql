-- +goose Up
INSERT INTO mcps (name, author, description, tools) VALUES (
    'gitlab',
    'Artel',
    'GitLab integration. Create/list merge requests and issues, comment on them.',
    '[
        {
            "api_description": {
                "name": "create_merge_request",
                "description": "Create a new merge request in a GitLab project",
                "properties": {
                    "project_id":    { "type": "string",  "description": "GitLab project ID or URL-encoded path (e.g. group%2Fproject)" },
                    "source_branch": { "type": "string",  "description": "Branch to merge from" },
                    "target_branch": { "type": "string",  "description": "Branch to merge into" },
                    "title":         { "type": "string",  "description": "Merge request title" },
                    "description":   { "type": "string",  "description": "Merge request description (optional)" }
                },
                "required": ["project_id", "source_branch", "target_branch", "title"]
            },
            "action": {
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
            }
        },
        {
            "api_description": {
                "name": "list_merge_requests",
                "description": "List merge requests in a GitLab project",
                "properties": {
                    "project_id": { "type": "string", "description": "GitLab project ID or URL-encoded path" },
                    "state":      { "type": "string", "description": "Filter by state", "enum": ["opened", "closed", "merged", "all"] }
                },
                "required": ["project_id"]
            },
            "action": {
                "http": {
                    "method": "GET",
                    "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/merge_requests",
                    "headers": { "PRIVATE-TOKEN": "__secrets.personal_access_token" },
                    "query": { "state": "${{params.state}}" },
                    "credentials": "gitlab"
                }
            }
        },
        {
            "api_description": {
                "name": "create_issue",
                "description": "Create a new issue in a GitLab project",
                "properties": {
                    "project_id": { "type": "string", "description": "GitLab project ID or URL-encoded path" },
                    "title":      { "type": "string", "description": "Issue title" },
                    "description":{ "type": "string", "description": "Issue description (optional)" }
                },
                "required": ["project_id", "title"]
            },
            "action": {
                "http": {
                    "method": "POST",
                    "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/issues",
                    "headers": { "PRIVATE-TOKEN": "__secrets.personal_access_token" },
                    "query": { "title": "${{params.title}}", "description": "${{params.description}}" },
                    "credentials": "gitlab"
                }
            }
        },
        {
            "api_description": {
                "name": "list_issues",
                "description": "List issues in a GitLab project",
                "properties": {
                    "project_id": { "type": "string", "description": "GitLab project ID or URL-encoded path" },
                    "state":      { "type": "string", "description": "Filter by state", "enum": ["opened", "closed", "all"] }
                },
                "required": ["project_id"]
            },
            "action": {
                "http": {
                    "method": "GET",
                    "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/issues",
                    "headers": { "PRIVATE-TOKEN": "__secrets.personal_access_token" },
                    "query": { "state": "${{params.state}}" },
                    "credentials": "gitlab"
                }
            }
        },
        {
            "api_description": {
                "name": "add_comment",
                "description": "Add a comment (note) to a GitLab merge request or issue",
                "properties": {
                    "project_id":     { "type": "string",  "description": "GitLab project ID or URL-encoded path" },
                    "noteable_type":  { "type": "string",  "description": "Type of item to comment on", "enum": ["merge_requests", "issues"] },
                    "noteable_iid":   { "type": "string",  "description": "Internal ID (iid) of the merge request or issue" },
                    "body":           { "type": "string",  "description": "Comment text" }
                },
                "required": ["project_id", "noteable_type", "noteable_iid", "body"]
            },
            "action": {
                "http": {
                    "method": "POST",
                    "url": "${{secrets.instance_url}}/api/v4/projects/${{params.project_id}}/${{params.noteable_type}}/${{params.noteable_iid}}/notes",
                    "headers": { "PRIVATE-TOKEN": "__secrets.personal_access_token" },
                    "query": { "body": "${{params.body}}" },
                    "credentials": "gitlab"
                }
            }
        }
    ]'
) ON CONFLICT (name) DO NOTHING;

-- +goose Down
DELETE FROM mcps WHERE name = 'gitlab';
