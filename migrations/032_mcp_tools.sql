-- +goose Up
-- Lifts mcps.tools (one JSONB array per MoM) into one row per tool, so tools can carry an
-- output_schema and be looked up individually (McpDefinitionsRepo.GetTool/ListAllTools).
--
-- Rollback path (see Down below): mcp_tools rows are re-aggregated with jsonb_agg back into the
-- {api_description:{name,description,properties,required}, action} shape mcps.tools used to
-- hold, grouped by mcp_name, then the column is restored and the table dropped.
CREATE TABLE mcp_tools (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mcp_name      TEXT NOT NULL REFERENCES mcps(name) ON DELETE CASCADE,
    name          TEXT NOT NULL,
    description   TEXT NOT NULL DEFAULT '',
    input_schema  JSONB NOT NULL DEFAULT '{}',   -- {properties:{...}, required:[...]}
    output_schema JSONB NOT NULL DEFAULT '{}',   -- same shape; {} = undeclared
    action        JSONB NOT NULL,                -- {imap|smtp|http: {...}} unchanged
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (mcp_name, name)
);

-- name is extracted via ->>'name' into a NOT NULL column: a tool document missing
-- api_description.name fails the migration loudly instead of silently dropping the tool.
INSERT INTO mcp_tools (mcp_name, name, description, input_schema, action)
SELECT m.name,
       t->'api_description'->>'name',
       COALESCE(t->'api_description'->>'description', ''),
       jsonb_build_object('properties', COALESCE(t->'api_description'->'properties', '{}'::jsonb),
                           'required', COALESCE(t->'api_description'->'required', '[]'::jsonb)),
       t->'action'
FROM mcps m, jsonb_array_elements(m.tools) t;

ALTER TABLE mcps DROP COLUMN tools;

-- +goose Down
ALTER TABLE mcps ADD COLUMN tools JSONB NOT NULL DEFAULT '[]';

UPDATE mcps m
SET tools = COALESCE(agg.tools, '[]'::jsonb)
FROM (
    SELECT mcp_name,
           jsonb_agg(
               jsonb_build_object(
                   'api_description', jsonb_build_object(
                       'name', name,
                       'description', description,
                       'properties', input_schema->'properties',
                       'required', input_schema->'required'
                   ),
                   'action', action
               )
           ) AS tools
    FROM mcp_tools
    GROUP BY mcp_name
) agg
WHERE agg.mcp_name = m.name;

DROP TABLE mcp_tools;
