-- +goose Up
UPDATE mcps SET tools = tools || '[
    {
        "api_description": {
            "name": "get_current_user",
            "description": "Get the current authenticated GitLab user. The returned username is your GitLab handle.",
            "properties": {},
            "required": []
        },
        "action": {
            "http": {
                "method": "GET",
                "url": "${{secrets.instance_url}}/api/v4/user",
                "headers": { "PRIVATE-TOKEN": "__secrets.personal_access_token" },
                "credentials": "gitlab"
            }
        }
    }
]'::jsonb
WHERE name = 'gitlab';

-- +goose Down
UPDATE mcps SET tools = (
    SELECT COALESCE(jsonb_agg(t), '[]'::jsonb)
    FROM jsonb_array_elements(tools) t
    WHERE t->'api_description'->>'name' != 'get_current_user'
)
WHERE name = 'gitlab';
