package tract

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/subscription"
	"github.com/ruf-dev/artel/internal/service/v1/tract/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestServiceWithTemplates builds a Service backed by a fakeTractTemplatesRepo, for tests
// that only exercise the template-related methods (PublishTemplate/UnpublishTemplate/
// ListTemplates/GetTemplate).
func newTestServiceWithTemplates(templates *fakeTractTemplatesRepo) *Service {
	svc := New(nil, templates, nil, nil, nil, nil, nil, nil, nil, nil)

	return svc
}

// newTestServiceForTemplateConnections builds a fully-wired Service (real tracts/external-conns/
// mcp-defs/script-engine fakes, mirroring newEngineTestServiceWithLlm's setup) for
// PublishTemplate/InstantiateTemplate tests that need a source tract to snapshot, an owned
// connection to validate against, or a template whose steps go through full createTractInternal
// validation (validateTools/validateScriptEngines/validateLlmConnections) — unlike
// newTestServiceWithTemplates's templates-only stub.
func newTestServiceForTemplateConnections() (*Service, *fakeTractsRepo, *fakeTractTemplatesRepo, *fakeExternalConnsRepo, *fakeMcpDefsRepo) {
	tracts := newFakeTractsRepo()
	templates := newFakeTractTemplatesRepo()
	externalConns := newFakeExternalConnsRepo()
	mcpDefs := newFakeMcpDefsRepo()
	toolExecutor := newFakeToolExecutor()
	scriptEngines := script.NewRegistry(script.NewJavaScriptEngine())

	svc := New(
		tracts, templates, nil, nil, externalConns, mcpDefs, toolExecutor, subscription.NewFree(), scriptEngines,
		newFakeLlmExecutor(),
	)

	return svc, tracts, templates, externalConns, mcpDefs
}

func TestPublishTemplate_StripsLlmConnectionUuid(t *testing.T) {
	svc, tracts, _, _, _ := newTestServiceForTemplateConnections()

	callerUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: callerUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	sourceUuid := uuid.New()
	source := domain.Tract{
		Uuid:     sourceUuid,
		UserUuid: callerUuid,
		Name:     "Describe MR on merge",
		Definition: domain.TractDefinition{
			Steps: []domain.TractStep{
				{
					Id:                "summarize",
					Type:              stepTypeLlmCall,
					LlmConnectionUuid: uuid.New(),
					LlmModel:          "claude-opus-4-8",
					Prompt:            "Summarize: {{ trigger.mr_iid }}",
				},
			},
		},
	}
	tracts.seedTract(source)

	published, err := svc.PublishTemplate(ctx, sourceUuid, "gitlab")
	require.NoError(t, err)
	require.Len(t, published.Definition.Steps, 1)
	assert.Equal(t, uuid.Nil, published.Definition.Steps[0].LlmConnectionUuid,
		"LlmConnectionUuid must be stripped on publish, same as an action step's ConnectionUuid")
	assert.Equal(t, "claude-opus-4-8", published.Definition.Steps[0].LlmModel,
		"non-connection llm_call fields must survive publish unchanged")
}

func TestInstantiateTemplate_WiresLlmConnectionUuidFromProviderKey(t *testing.T) {
	svc, _, templates, externalConns, _ := newTestServiceForTemplateConnections()

	callerUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: callerUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	ownedConn := domain.ExternalConnection{
		Uuid:     uuid.New(),
		UserUuid: callerUuid,
		Provider: domain.ProviderAnthropic,
	}
	_, err := externalConns.Insert(ctx, ownedConn)
	require.NoError(t, err)

	templateUuid := uuid.New()
	templates.templates[templateUuid] = domain.TractTemplate{
		Uuid: templateUuid,
		Definition: domain.TractDefinition{
			Steps: []domain.TractStep{
				{Id: "summarize", Type: stepTypeLlmCall, LlmModel: "claude-opus-4-8", Prompt: "hi"},
			},
		},
	}

	connections := map[string]uuid.UUID{domain.ProviderAnthropic: ownedConn.Uuid}

	created, _, err := svc.InstantiateTemplate(ctx, templateUuid, "My tract", "", connections)
	require.NoError(t, err)
	require.Len(t, created.Definition.Steps, 1)
	assert.Equal(t, ownedConn.Uuid, created.Definition.Steps[0].LlmConnectionUuid)
}

// describeMrOnMergeTemplateDefinition mirrors migrations/058_seed_gitlab_mr_merged_template.sql's
// "Describe MR on merge" built-in template definition verbatim — a regression guard that the
// seeded JSON actually passes save-time validation (shape, dataflow visibility, llm connection
// checks) end to end, the same path a real install goes through via InstantiateTemplate. Keep in
// sync with the migration if either changes.
const describeMrOnMergeTemplateDefinition = `{"steps": [
    {"id": "get_diff", "mcp": "gitlab", "name": "get_merge_request_diff", "tool": "get_merge_request_diff", "type": "action", "params": {"mr_iid": "{{ trigger.mr_iid }}", "project_id": "{{ trigger.project.id }}"}, "connection_uuid": "00000000-0000-0000-0000-000000000000"},
    {"id": "format_diff", "code": "diff_text = diffs.map(function(d) { return \"--- \" + d.old_path + \"\\n+++ \" + d.new_path + \"\\n\" + d.diff; }).join(\"\\n\\n\")", "name": "format_diff", "type": "script", "params": {"diffs": "{{ get_diff }}"}, "language": "javascript", "input_params": [{"Name": "diffs", "Property": {"Enum": null, "Type": "array", "Items": null, "Required": null, "Properties": null, "Description": ""}}], "output_params": [{"Name": "diff_text", "Property": {"Enum": null, "Type": "string", "Items": null, "Required": null, "Properties": null, "Description": ""}}], "connection_uuid": "00000000-0000-0000-0000-000000000000"},
    {"id": "summarize", "name": "summarize", "type": "llm_call", "llm_model": "claude-opus-4-8", "llm_connection_uuid": "00000000-0000-0000-0000-000000000000", "system_prompt": "You are a helpful assistant that writes clear, concise merge request descriptions from code diffs. Respond with only the description text, no preamble.", "prompt": "Summarize the work done in this merge request based on the diff below. Write a short, professional description suitable for the merge request description field.\n\nTitle: {{ trigger.object_attributes.title }}\n\nDiff:\n{{ format_diff.diff_text }}"},
    {"id": "update_description", "mcp": "gitlab", "name": "update_merge_request", "tool": "update_merge_request", "type": "action", "params": {"mr_iid": "{{ trigger.mr_iid }}", "project_id": "{{ trigger.project.id }}", "description": "{{ summarize.text }}"}, "connection_uuid": "00000000-0000-0000-0000-000000000000"}
]}`

func TestInstantiateTemplate_DescribeMrOnMergeSeedPassesValidation(t *testing.T) {
	svc, _, templates, externalConns, mcpDefs := newTestServiceForTemplateConnections()
	mcpDefs.addTool("gitlab", "get_merge_request_diff")
	mcpDefs.addTool("gitlab", "update_merge_request")

	callerUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: callerUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	anthropicConn := domain.ExternalConnection{Uuid: uuid.New(), UserUuid: callerUuid, Provider: domain.ProviderAnthropic}
	_, err := externalConns.Insert(ctx, anthropicConn)
	require.NoError(t, err)

	gitlabConn := domain.ExternalConnection{Uuid: uuid.New(), UserUuid: callerUuid, Provider: "gitlab"}
	_, err = externalConns.Insert(ctx, gitlabConn)
	require.NoError(t, err)

	var def domain.TractDefinition

	err = json.Unmarshal([]byte(describeMrOnMergeTemplateDefinition), &def)
	require.NoError(t, err)

	templateUuid := uuid.New()
	templates.templates[templateUuid] = domain.TractTemplate{Uuid: templateUuid, Definition: def}

	connections := map[string]uuid.UUID{
		"gitlab":                 gitlabConn.Uuid,
		domain.ProviderAnthropic: anthropicConn.Uuid,
	}

	created, warnings, err := svc.InstantiateTemplate(ctx, templateUuid, "Describe MR on merge", "", connections)
	require.NoError(t, err)
	assert.Empty(t, warnings)
	require.Len(t, created.Definition.Steps, 4)
	assert.Equal(t, anthropicConn.Uuid, created.Definition.Steps[2].LlmConnectionUuid)
	assert.Equal(t, gitlabConn.Uuid, created.Definition.Steps[0].ConnectionUuid)
	assert.Equal(t, gitlabConn.Uuid, created.Definition.Steps[3].ConnectionUuid)
}

func TestInstantiateTemplate_MissingLlmConnectionRejected(t *testing.T) {
	svc, _, templates, _, _ := newTestServiceForTemplateConnections()

	callerUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: callerUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	templateUuid := uuid.New()
	templates.templates[templateUuid] = domain.TractTemplate{
		Uuid: templateUuid,
		Definition: domain.TractDefinition{
			Steps: []domain.TractStep{
				{Id: "summarize", Type: stepTypeLlmCall, LlmModel: "claude-opus-4-8", Prompt: "hi"},
			},
		},
	}

	_, _, err := svc.InstantiateTemplate(ctx, templateUuid, "My tract", "", nil)
	require.Error(t, err, "instantiate must reject a template whose llm_call step has no matching provider connection supplied — all-or-nothing, same as an unwired action step")
}

func TestUnpublishTemplate_BuiltinNeverOwned(t *testing.T) {
	templates := newFakeTractTemplatesRepo()

	builtin := domain.TractTemplate{
		Uuid:      uuid.New(),
		OwnerUuid: uuid.Nil,
		Name:      "Create MR on featurep push",
		Category:  "Mr on push",
	}
	templates.templates[builtin.Uuid] = builtin

	svc := newTestServiceWithTemplates(templates)

	callerUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: callerUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	err := svc.UnpublishTemplate(ctx, builtin.Uuid)
	require.Error(t, err)
	assert.True(t, errors.Is(err, user_errors.TractTemplateNotOwned))

	// The template must still be present — the delete was never reached.
	_, ok := templates.templates[builtin.Uuid]
	assert.True(t, ok)
}

func TestListTemplates_IncludesBuiltins(t *testing.T) {
	templates := newFakeTractTemplatesRepo()

	builtin := domain.TractTemplate{
		Uuid:      uuid.New(),
		OwnerUuid: uuid.Nil,
		Name:      "Create MR on featurep push",
		Category:  "Mr on push",
	}
	templates.templates[builtin.Uuid] = builtin

	svc := newTestServiceWithTemplates(templates)

	callerUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: callerUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	got, err := svc.ListTemplates(ctx, "", false)
	require.NoError(t, err)

	var found bool

	for _, template := range got {
		if template.Uuid == builtin.Uuid {
			found = true
		}
	}

	assert.True(t, found, "expected ListTemplates(mineOnly=false) to include the built-in template")
}
