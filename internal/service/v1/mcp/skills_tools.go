package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// Static builtin tool names for skill management — every authenticated MCP key may call these
// (no admin gate); skills.Service itself enforces ownership/quota. Not admin-gated, unlike
// create_community_connector.
const (
	toolListSkills      = "list_skills"
	toolGetSkillBody    = "get_skill_body"
	toolCreateSkill     = "create_skill"
	toolUpdateSkill     = "update_skill"
	toolDeleteSkill     = "delete_skill"
	toolSetSkillHotPlug = "set_skill_hot_plug"
)

// Field names used across this file's tool schemas/params/results.
const (
	fieldSlug        = "slug"
	fieldStorageMode = "storage_mode"
	fieldBody        = "body"
	fieldHotPlug     = "hot_plug"
	fieldIsHotPlug   = "is_hot_plug"
	fieldIsSystem    = "is_system"
)

// skillStorageModeEnum is the wire-visible enum for the storage_mode param/field, shared by
// every skill tool schema that carries it.
var skillStorageModeEnum = []string{
	string(domain.SkillStorageNone),
	string(domain.SkillStorageFreeform),
	string(domain.SkillStorageStructured),
}

// isSkillManagementTool reports whether name is one of the six static skill_* builtin tools —
// distinct from the dynamic per-hot-plug-skill skill_<slug> tools, which are not part of this
// static set (see ListHotPlugSkillTools/ExecuteSkillTool).
func isSkillManagementTool(name string) bool {
	switch name {
	case toolListSkills, toolGetSkillBody, toolCreateSkill, toolUpdateSkill, toolDeleteSkill, toolSetSkillHotPlug:
		return true
	}

	return false
}

// skillsToolDefinitions describes the six static skill-management builtin tools, appended to
// ListTools's static catalog.
func skillsToolDefinitions() []domain.McpToolDef {
	listSkillsTool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name: toolListSkills,
			Description: "List every skill visible in this vault (the built-in skill-creator plus every " +
				"skill under the active/library folders), without body text — use get_skill_body to read " +
				"one's full instructions.",
			Properties: map[string]domain.ToolProperty{},
			Required:   []string{},
		},
		OutputSchema: domain.ToolSchema{
			IsArray: true,
			Items: &domain.ToolProperty{
				Type: schemaTypeObject,
				Properties: map[string]domain.ToolProperty{
					fieldSlug:        {Type: schemaTypeString, Description: "Skill slug, used to refer to it in other skill tools"},
					fieldName:        {Type: schemaTypeString, Description: "Human-readable skill name"},
					fieldDescription: {Type: schemaTypeString, Description: "Keyword-trigger description"},
					fieldStorageMode: {
						Type:        schemaTypeString,
						Description: "none | freeform_notes | structured",
						Enum:        skillStorageModeEnum,
					},
					fieldIsHotPlug: {
						Type:        schemaTypeBoolean,
						Description: "Whether this skill is hot-plugged as its own skill_<slug> tool",
					},
					fieldIsSystem: {Type: schemaTypeBoolean, Description: "Whether this is the built-in skill-creator skill"},
				},
				Required: []string{fieldSlug, fieldName, fieldDescription, fieldStorageMode, fieldIsHotPlug, fieldIsSystem},
			},
		},
	}

	getSkillBodyTool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:        toolGetSkillBody,
			Description: "Read a skill's full instructional body by slug (see list_skills for available slugs).",
			Properties: map[string]domain.ToolProperty{
				fieldSlug: {Type: schemaTypeString},
			},
			Required: []string{fieldSlug},
		},
		OutputSchema: domain.ToolSchema{
			Properties: map[string]domain.ToolProperty{
				fieldContent: {Type: schemaTypeString, Description: "The skill's instructional markdown body"},
			},
			Required: []string{fieldContent},
		},
	}

	createSkillTool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name: toolCreateSkill,
			Description: "Create a new skill in this vault. name is slugified into a unique slug " +
				"(returned in the result); hot_plug=true synthesizes it immediately as its own skill_<slug> " +
				"tool, subject to the plan's hot-plug cap — leave it false to add to the library instead.",
			Properties: map[string]domain.ToolProperty{
				fieldName:        {Type: schemaTypeString, Description: "Human-readable skill name; slugified to form the skill's slug"},
				fieldDescription: {Type: schemaTypeString, Description: "Keyword-trigger description shown to other agents once hot-plugged"},
				fieldStorageMode: {
					Type:        schemaTypeString,
					Description: "none | freeform_notes | structured",
					Enum:        skillStorageModeEnum,
				},
				fieldBody:    {Type: schemaTypeString, Description: "Instructional markdown body, returned verbatim when the skill is invoked"},
				fieldHotPlug: {Type: schemaTypeBoolean, Description: "Hot-plug immediately as its own skill_<slug> tool; defaults to false"},
			},
			Required: []string{fieldName, fieldDescription, fieldStorageMode, fieldBody},
		},
		OutputSchema: domain.ToolSchema{
			Properties: map[string]domain.ToolProperty{
				fieldSlug: {Type: schemaTypeString, Description: "The created skill's slug"},
			},
			Required: []string{fieldSlug},
		},
	}

	updateSkillTool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name: toolUpdateSkill,
			Description: "Rewrite an existing skill's name/description/storage_mode/body in place. The " +
				"slug never changes — use set_skill_hot_plug to move a skill between the library and the " +
				"hot-plug set.",
			Properties: map[string]domain.ToolProperty{
				fieldSlug:        {Type: schemaTypeString},
				fieldName:        {Type: schemaTypeString},
				fieldDescription: {Type: schemaTypeString},
				fieldStorageMode: {Type: schemaTypeString, Enum: skillStorageModeEnum},
				fieldBody:        {Type: schemaTypeString},
			},
			Required: []string{fieldSlug, fieldName, fieldDescription, fieldStorageMode, fieldBody},
		},
		OutputSchema: domain.ToolSchema{
			Properties: map[string]domain.ToolProperty{
				fieldSlug: {Type: schemaTypeString},
			},
			Required: []string{fieldSlug},
		},
	}

	deleteSkillTool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:        toolDeleteSkill,
			Description: "Delete a skill by slug from whichever folder (active or library) it currently lives in.",
			Properties: map[string]domain.ToolProperty{
				fieldSlug: {Type: schemaTypeString},
			},
			Required: []string{fieldSlug},
		},
		OutputSchema: domain.ToolSchema{
			Properties: map[string]domain.ToolProperty{
				fieldMessage: {Type: schemaTypeString, Description: msgConfirmation},
			},
			Required: []string{fieldMessage},
		},
	}

	setSkillHotPlugTool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name: toolSetSkillHotPlug,
			Description: "Move a skill between the hot-plug set (synthesized as its own skill_<slug> " +
				"tool, subject to the plan's hot-plug cap) and the library (available but not synthesized " +
				"into a standalone tool).",
			Properties: map[string]domain.ToolProperty{
				fieldSlug:    {Type: schemaTypeString},
				fieldHotPlug: {Type: schemaTypeBoolean},
			},
			Required: []string{fieldSlug, fieldHotPlug},
		},
		OutputSchema: domain.ToolSchema{
			Properties: map[string]domain.ToolProperty{
				fieldSlug:      {Type: schemaTypeString},
				fieldIsHotPlug: {Type: schemaTypeBoolean},
			},
			Required: []string{fieldSlug, fieldIsHotPlug},
		},
	}

	return []domain.McpToolDef{
		listSkillsTool, getSkillBodyTool, createSkillTool, updateSkillTool, deleteSkillTool, setSkillHotPlugTool,
	}
}

// withSkillCallerContext embeds keyCtx's owning user into ctx as a user_context.UserContext —
// skills.Service.resolveVault (mirroring notes.Service) authenticates/authorizes off
// user_context rather than off keyCtx directly, since it's shared with the dashboard's
// session-authenticated skills_api transport. Mirrors executeTractTool's identical embedding for
// the tract-authoring builtin tools.
func withSkillCallerContext(ctx context.Context, keyCtx domain.McpKeyContext) context.Context {
	uc := user_context.UserContext{UserUuid: keyCtx.UserUuid}

	return user_context.WithUserContext(ctx, uc)
}

// skillToolCollides reports whether hot-plugging a skill with this slug would synthesize a
// dynamic skill_<slug> tool name that collides with an existing builtin (or tract/postgres) tool
// name — isBuiltin is passed in (rather than always calling (*ServiceImpl).IsBuiltinTool
// directly) so this pure decision can be unit-tested against a fake predicate without needing an
// actual colliding builtin tool name to exist.
func skillToolCollides(isBuiltin func(name string) bool, slug string) bool {
	return isBuiltin(domain.SkillToolPrefix + slug)
}

func isValidSkillStorageMode(mode domain.SkillStorageMode) bool {
	switch mode {
	case domain.SkillStorageNone, domain.SkillStorageFreeform, domain.SkillStorageStructured:
		return true
	}

	return false
}

// executeSkillManagementTool dispatches one of the six static skill_* builtin tool names to its
// handler.
func (s *ServiceImpl) executeSkillManagementTool(
	ctx context.Context,
	keyCtx domain.McpKeyContext,
	toolName string,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	switch toolName {
	case toolListSkills:
		return s.listSkillsTool(ctx, keyCtx)
	case toolGetSkillBody:
		return s.getSkillBodyTool(ctx, keyCtx, params)
	case toolCreateSkill:
		return s.createSkillTool(ctx, keyCtx, params)
	case toolUpdateSkill:
		return s.updateSkillTool(ctx, keyCtx, params)
	case toolDeleteSkill:
		return s.deleteSkillTool(ctx, keyCtx, params)
	case toolSetSkillHotPlug:
		return s.setSkillHotPlugTool(ctx, keyCtx, params)
	default:
		return domain.ToolExecResult{}, rerrors.New(fmt.Sprintf("unknown skill tool: %s", toolName))
	}
}

func (s *ServiceImpl) listSkillsTool(ctx context.Context, keyCtx domain.McpKeyContext) (domain.ToolExecResult, error) {
	ctx = withSkillCallerContext(ctx, keyCtx)

	skillList, err := s.skillsSvc.ListSkills(ctx, keyCtx.VaultUuid)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error listing skills")
	}

	rows := make([]map[string]interface{}, len(skillList))
	for i, sk := range skillList {
		rows[i] = map[string]interface{}{
			fieldSlug:        sk.Slug,
			fieldName:        sk.Name,
			fieldDescription: sk.Description,
			fieldStorageMode: string(sk.StorageMode),
			fieldIsHotPlug:   sk.IsHotPlug,
			fieldIsSystem:    sk.IsSystem,
		}
	}

	text, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error marshaling skills")
	}

	return domain.ToolExecResult{Text: string(text)}, nil
}

// getSkillBodyTool mirrors executors.VaultExecutor.readFile's plain-text result shape (rather
// than a JSON wrapper) — this is the single-item read counterpart to list_skills, and the caller
// already has slug/name/description/storage_mode from that listing.
func (s *ServiceImpl) getSkillBodyTool(
	ctx context.Context,
	keyCtx domain.McpKeyContext,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	slug, ok := params[fieldSlug].(string)
	if !ok || slug == "" {
		return domain.ToolExecResult{}, user_errors.SkillSlugRequired
	}

	ctx = withSkillCallerContext(ctx, keyCtx)

	skill, err := s.skillsSvc.GetSkillBody(ctx, keyCtx.VaultUuid, slug)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error getting skill body")
	}

	return domain.ToolExecResult{Text: skill.Body}, nil
}

func (s *ServiceImpl) createSkillTool(
	ctx context.Context,
	keyCtx domain.McpKeyContext,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	name, ok := params[fieldName].(string)
	if !ok || name == "" {
		return domain.ToolExecResult{}, user_errors.SkillNameRequired
	}

	description, ok := params[fieldDescription].(string)
	if !ok || description == "" {
		return domain.ToolExecResult{}, user_errors.SkillDescriptionRequired
	}

	storageModeStr, ok := params[fieldStorageMode].(string)
	if !ok || storageModeStr == "" {
		return domain.ToolExecResult{}, user_errors.SkillStorageModeRequired
	}

	storageMode := domain.SkillStorageMode(storageModeStr)
	if !isValidSkillStorageMode(storageMode) {
		return domain.ToolExecResult{}, user_errors.SkillStorageModeInvalid
	}

	body, ok := params[fieldBody].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.SkillBodyRequired
	}

	// hot_plug is optional, defaulting to false — the type assertion's zero value on a missing
	// or wrong-typed param is exactly that default, so its ok is intentionally ignored.
	hotPlug, _ := params[fieldHotPlug].(bool)

	ctx = withSkillCallerContext(ctx, keyCtx)

	skill, err := s.skillsSvc.CreateSkill(ctx, keyCtx.VaultUuid, name, description, storageMode, body, hotPlug)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error creating skill")
	}

	if hotPlug && skillToolCollides(s.IsBuiltinTool, skill.Slug) {
		delErr := s.skillsSvc.DeleteSkill(ctx, keyCtx.VaultUuid, skill.Slug)
		if delErr != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(delErr, "error rolling back colliding skill")
		}

		return domain.ToolExecResult{}, user_errors.SkillNameCollidesWithBuiltinTool
	}

	text := fmt.Sprintf("Skill %q created with slug %q.", skill.Name, skill.Slug)

	return domain.ToolExecResult{Text: text}, nil
}

func (s *ServiceImpl) updateSkillTool(
	ctx context.Context,
	keyCtx domain.McpKeyContext,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	slug, ok := params[fieldSlug].(string)
	if !ok || slug == "" {
		return domain.ToolExecResult{}, user_errors.SkillSlugRequired
	}

	name, ok := params[fieldName].(string)
	if !ok || name == "" {
		return domain.ToolExecResult{}, user_errors.SkillNameRequired
	}

	description, ok := params[fieldDescription].(string)
	if !ok || description == "" {
		return domain.ToolExecResult{}, user_errors.SkillDescriptionRequired
	}

	storageModeStr, ok := params[fieldStorageMode].(string)
	if !ok || storageModeStr == "" {
		return domain.ToolExecResult{}, user_errors.SkillStorageModeRequired
	}

	storageMode := domain.SkillStorageMode(storageModeStr)
	if !isValidSkillStorageMode(storageMode) {
		return domain.ToolExecResult{}, user_errors.SkillStorageModeInvalid
	}

	body, ok := params[fieldBody].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.SkillBodyRequired
	}

	ctx = withSkillCallerContext(ctx, keyCtx)

	skill, err := s.skillsSvc.UpdateSkill(ctx, keyCtx.VaultUuid, slug, name, description, storageMode, body)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error updating skill")
	}

	text := fmt.Sprintf("Skill %q updated.", skill.Slug)

	return domain.ToolExecResult{Text: text}, nil
}

func (s *ServiceImpl) deleteSkillTool(
	ctx context.Context,
	keyCtx domain.McpKeyContext,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	slug, ok := params[fieldSlug].(string)
	if !ok || slug == "" {
		return domain.ToolExecResult{}, user_errors.SkillSlugRequired
	}

	ctx = withSkillCallerContext(ctx, keyCtx)

	err := s.skillsSvc.DeleteSkill(ctx, keyCtx.VaultUuid, slug)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error deleting skill")
	}

	return domain.ToolExecResult{Text: "Skill deleted successfully"}, nil
}

func (s *ServiceImpl) setSkillHotPlugTool(
	ctx context.Context,
	keyCtx domain.McpKeyContext,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	slug, ok := params[fieldSlug].(string)
	if !ok || slug == "" {
		return domain.ToolExecResult{}, user_errors.SkillSlugRequired
	}

	hotPlug, ok := params[fieldHotPlug].(bool)
	if !ok {
		return domain.ToolExecResult{}, user_errors.SkillHotPlugRequired
	}

	ctx = withSkillCallerContext(ctx, keyCtx)

	skill, err := s.skillsSvc.SetSkillHotPlug(ctx, keyCtx.VaultUuid, slug, hotPlug)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error setting skill hot-plug state")
	}

	if hotPlug && skillToolCollides(s.IsBuiltinTool, skill.Slug) {
		_, revertErr := s.skillsSvc.SetSkillHotPlug(ctx, keyCtx.VaultUuid, slug, false)
		if revertErr != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(revertErr, "error rolling back colliding hot-plug toggle")
		}

		return domain.ToolExecResult{}, user_errors.SkillNameCollidesWithBuiltinTool
	}

	text := fmt.Sprintf("Skill %q hot-plug set to %t.", skill.Slug, skill.IsHotPlug)

	return domain.ToolExecResult{Text: text}, nil
}

// ListHotPlugSkillTools returns one dynamic skill_<slug> tool per hot-plug skill visible in
// keyCtx's vault (the always-hot-plug system skill-creator skill included). See the doc comment
// on service.McpService for why this is kept separate from ListTools.
func (s *ServiceImpl) ListHotPlugSkillTools(
	ctx context.Context,
	keyCtx domain.McpKeyContext,
) ([]domain.McpToolDef, error) {
	ctx = withSkillCallerContext(ctx, keyCtx)

	skillList, err := s.skillsSvc.ListSkills(ctx, keyCtx.VaultUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing skills for hot-plug tool synthesis")
	}

	tools := make([]domain.McpToolDef, 0, len(skillList))

	for _, sk := range skillList {
		if !sk.IsHotPlug {
			continue
		}

		tool := domain.McpToolDef{
			ApiDescription: domain.ToolApiDescription{
				Name:        domain.SkillToolPrefix + sk.Slug,
				Description: sk.Description,
				Properties:  map[string]domain.ToolProperty{},
				Required:    []string{},
			},
		}

		tools = append(tools, tool)
	}

	return tools, nil
}

// ExecuteSkillTool executes a dynamic skill_<slug> tool: returns the skill's body verbatim,
// regardless of its current hot-plug state (a caller that already knows the slug — e.g. from a
// tool list resolved moments earlier — can still fetch it even if it was just demoted).
func (s *ServiceImpl) ExecuteSkillTool(ctx context.Context, keyCtx domain.McpKeyContext, slug string) (string, error) {
	ctx = withSkillCallerContext(ctx, keyCtx)

	skill, err := s.skillsSvc.GetSkillBody(ctx, keyCtx.VaultUuid, slug)
	if err != nil {
		return "", rerrors.Wrap(err, "error executing skill tool")
	}

	return skill.Body, nil
}
