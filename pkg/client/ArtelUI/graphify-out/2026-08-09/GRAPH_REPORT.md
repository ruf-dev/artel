# Graph Report - ArtelUI  (2026-08-05)

## Corpus Check
- 595 files · ~153,773 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2733 nodes · 6686 edges · 126 communities (121 shown, 5 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 23 edges (avg confidence: 0.67)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e262c3dd`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- tracts.pb.ts
- TaskTrackersPage.tsx
- vaults.pb.ts
- useBakeError
- external_connections.pb.ts
- mcp_keys.pb.ts
- TemplateInput.tsx
- addTriggerDialogContext.ts
- ExternalConnectionInfo
- admin_couch.pb.ts
- TractCanvasInspectorBody.tsx
- couch_instances.pb.ts
- TractsService
- Dialog.ts
- Tracts.ts
- cn
- TractIcons.tsx
- compilerOptions
- ToolboxPage.tsx
- ConnectorChip.tsx
- grpcErrors.ts
- useDialog
- ManageKeyDialog.tsx
- tractCanvasLayout.ts
- auth.pb.ts
- VaultItem
- devDependencies
- fetch.pb.ts
- McpKeysAPI
- Router.tsx
- TractStepTree.tsx
- BreadcrumbBar.tsx
- NotesSidebar.tsx
- useExternalConnections
- ConnectionDetailDialog.tsx
- notes.pb.ts
- Notes.ts
- index.ts
- s3_instances.pb.ts
- useVaultMutations
- CreateNoteDialog.tsx
- ArtelUI Frontend Rules
- SchemaProperty
- NoteEditor.tsx
- McpAuthPage.tsx
- MobileNotesShell.tsx
- Tracts.ts
- useErrorToast.ts
- dependencies
- ProviderIcon.tsx
- StepRow.tsx
- AuthMiddleware
- tractSteps.ts
- admin_users.pb.ts
- Topbar.tsx
- TractCanvasBuilder.tsx
- User.ts
- Errors.ts
- NoteViewer.tsx
- InviteLinksSection.tsx
- TractBlockPicker.tsx
- toTract
- NotesPage.tsx
- VaultCard.tsx
- TractCanvasTopBar.tsx
- TractCanvasLogPanel.tsx
- scripts
- AuthAPI
- connectionLabel
- S3InstanceFormDialog.tsx
- compilerOptions
- McpKeys.ts
- ResultView.tsx
- MembersSection.tsx
- LinkScreen.tsx
- dialog-scrollable.js
- Handoff: lint/tooling parity gaps vs. ZpotifyUI
- RunTractDialog.tsx
- DbAccessList.tsx
- VaultDangerZone.tsx
- CardMeta.tsx
- EmptyState.tsx
- eslint.config.js
- GoogleSheetsSpreadsheetSection.tsx
- VaultCardHeader.tsx
- AuthMiddleware
- ConnectForm.tsx
- UsersTab.tsx
- AdminCouchAPI
- CouchInstancesAPI
- TractCanvasLogPanel.tsx
- RunLog.tsx
- package.json
- MobileDrawer.tsx
- CouchInstancesAPI
- .getRun
- UserList.tsx
- Vaults.ts
- AuthFetchInterceptor.ts
- EmailCheckButton.tsx
- GitlabCheckButton.tsx
- TrelloCheckButton.tsx
- TopbarDrawerCloseButton.tsx
- InsertConflictDialog.tsx
- TopbarMobileTrigger.tsx
- bun-test.d.ts
- requiredConnections.ts
- package.json
- ConnectedContent.tsx
- KebabMenu.tsx
- useTheme.ts
- RoadmapCanvasNode.tsx
- Tabs.tsx
- ParamsList.tsx
- ImportZipDialog.tsx

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 206 edges
2. `cn` - 162 edges
3. `useBakeError()` - 156 edges
4. `useUser` - 103 edges
5. `useExternalConnections` - 68 edges
6. `TractStep` - 46 edges
7. `TractTool` - 42 edges
8. `MomCandidate` - 41 edges
9. `ExternalConnectionInfo` - 40 edges
10. `useMcpKeys` - 39 edges

## Surprising Connections (you probably didn't know these)
- `NoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/notes/components/NoteViewer/NoteViewer.tsx → package.json
- `DocsNoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/docs/components/DocsNoteViewer/DocsNoteViewer.tsx → package.json
- `DialogHeadProps` --references--> `ExternalProvider`  [EXTRACTED]
  src/dialogs/ConnectionDetailDialog/components/DialogHead/DialogHead.tsx → src/app/api/artel/external_connections.pb.ts
- `AccountsSectionProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/AccountsSection.tsx → src/app/api/artel/external_connections.pb.ts
- `EmailConnectionRowProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/components/EmailConnectionRow/EmailConnectionRow.tsx → src/app/api/artel/external_connections.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/workbench/WorkbenchPage.tsx -> src/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/setup-wizard/SetupWizardPage.tsx -> src/pages/setup-wizard/screens/CreateAdminScreen.tsx -> src/app/routing/Router.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarMobileDrawer/TopbarMobileDrawer.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarUserMenu/TopbarUserMenu.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarBrand/TopbarBrand.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarNav/TopbarNav.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 5-file cycle: `src/app/routing/Router.tsx -> src/pages/tract-templates/TractTemplatesListPage.tsx -> src/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx -> src/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx -> src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx -> src/app/routing/Router.tsx`
- 5-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarMobileDrawer/TopbarMobileDrawer.tsx -> src/segments/Topbar/components/TopbarBrand/TopbarBrand.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`

## Communities (126 total, 5 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.02
Nodes (88): Absent, BaseTractStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger, CreateTriggerRequest, CreateTriggerResponse (+80 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (29): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+21 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.03
Nodes (63): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+55 more)

### Community 3 - "useBakeError"
Cohesion: 0.12
Nodes (15): GetVaultBySlug, GetVaultBySlugRequest, GetVaultBySlugResponse, PublicDocsGetNote, PublicDocsGetNoteRequest, PublicDocsGetNoteResponse, PublicDocsListFolders, PublicDocsListFoldersRequest (+7 more)

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.03
Nodes (57): Absent, AddAnthropicConnection, AddAnthropicConnectionResponse, AddEmailConnection, AddEmailConnectionResponse, AddGenericConnection, AddGenericConnectionResponse, AddGitlabConnection (+49 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (42): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+34 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.14
Nodes (23): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenArrayItems() (+15 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.11
Nodes (25): SetTractsState, triggerSourcesQueryKey, triggersQueryKey, useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor() (+17 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.08
Nodes (17): AddAnthropicConnectionRequest, AddEmailConnectionRequest, AddGenericConnectionRequest, AddGitlabConnectionRequest, AddTelegramConnectionRequest, AddTrelloConnectionRequest, CheckAnthropicConnectionRequest, CheckAnthropicConnectionResponse (+9 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.12
Nodes (44): MomCandidate, ScriptLanguage, ActionCard(), Props, buildSourcesFor(), collectIdsFromRoot(), ConditionCard(), ConditionCardProps (+36 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.09
Nodes (22): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+14 more)

### Community 12 - "TractsService"
Cohesion: 0.06
Nodes (18): TractsAPI, safeParseJson(), Props, TractsService, emptySchema, parseSchema(), toRun(), toRunStep() (+10 more)

### Community 13 - "Dialog.ts"
Cohesion: 0.22
Nodes (7): useServerStatus(), HomeLayout(), Router(), routes, GoogleOAuthCallbackPage(), ErrorPage(), ServiceUnavailablePage()

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.06
Nodes (37): cn, KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it, TODO: chures has no tab primitive yet, drop this wrapper once it does, Tabs(), TabsProps (+29 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.10
Nodes (19): useMcpKeys, Props, SearchInput(), ConnectionStep(), Props, AddTaskLinkDialog(), Props, RELATION_LABEL (+11 more)

### Community 17 - "compilerOptions"
Cohesion: 0.30
Nodes (12): TractsState, formatStartedAt(), Props, RunTractDialog(), Props, Deps, ITractsService, CreatedTrigger (+4 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.09
Nodes (25): DialogManager, useDialog, LlmConnectionStep(), Props, Props, ToolStep(), StepCard(), TractStepTree() (+17 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.10
Nodes (22): UserErrors, CreateVaultDialog(), HomePage(), BadRequestDetail, DetailType, DetailTypeName, ErrorInfoDetail, FieldViolation (+14 more)

### Community 21 - "useDialog"
Cohesion: 0.13
Nodes (17): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind, isNonEmptyBranch(), JsonBlock() (+9 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.11
Nodes (18): Absent, BaseCompleteSetupRequest, CompleteSetup, CompleteSetupRequest, CompleteSetupResponse, GetStatus, GetStatusRequest, GetStatusResponse (+10 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.12
Nodes (24): ConnectorPath(), ConnectorPathProps, ParallelBoxes(), Props, Props, TractCanvasArea(), useTractCanvasDrag(), cap() (+16 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (23): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+15 more)

### Community 25 - "VaultItem"
Cohesion: 0.14
Nodes (14): VaultItem, CardChips(), Props, Props, VaultCardHeader(), VaultChip(), ExpertSettingsDrawer(), Props (+6 more)

### Community 26 - "devDependencies"
Cohesion: 0.05
Nodes (42): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+34 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.07
Nodes (30): GetSettings, GetSettingsRequest, GetSettingsResponse, UpdateAuthMethods, UpdateAuthMethodsRequest, UpdateAuthMethodsResponse, UpdateRegistrationMode, UpdateRegistrationModeRequest (+22 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.11
Nodes (7): CommunityConnectorInfo, CreateMcpKeyResponse, McpKeyInfo, McpKeysAPI, McpKeysState, IMcpKeysService, McpKeysService

### Community 29 - "Router.tsx"
Cohesion: 0.19
Nodes (12): useNotes, ArtelLogoIcon(), MobileDrawer(), MobileDrawerProps, VaultOption, MobileNotesShell(), VaultOption, NotesSidebar() (+4 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.11
Nodes (21): McpConnectorInfo, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps (+13 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.12
Nodes (23): CheckEmailConnectionRequest, CheckGitlabConnectionRequest, CheckTelegramConnectionRequest, CheckTrelloConnectionRequest, Path, CheckStatus, EmailCheckButton(), EmailCheckButtonProps (+15 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.21
Nodes (8): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, useNotesSearchQuery()

### Community 33 - "useExternalConnections"
Cohesion: 0.10
Nodes (16): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, Props, SchemaFieldRow() (+8 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.13
Nodes (12): SuggestionList(), SuggestionListProps, CreateNoteDialog(), Props, PublicBadge(), SidebarTopBar(), SidebarTopBarProps, VaultOption (+4 more)

### Community 35 - "notes.pb.ts"
Cohesion: 0.06
Nodes (33): CheckImportConflicts, CheckImportConflictsRequest, CheckImportConflictsResponse, CommitImport, CommitImportRequest, CommitImportResponse, DeleteFolder, DeleteFolderRequest (+25 more)

### Community 36 - "Notes.ts"
Cohesion: 0.11
Nodes (4): b64Decode(), NotesAPI, INotesService, NotesService

### Community 37 - "index.ts"
Cohesion: 0.22
Nodes (8): ListPrompts, ListPromptsRequest, ListPromptsResponse, PromptId, PromptItem, PromptsAPI, FastSetupDialog(), Props

### Community 38 - "s3_instances.pb.ts"
Cohesion: 0.14
Nodes (18): paramCompletionSource(), scriptEditorTheme, scriptHighlightStyle, Props, ScriptCodeSection(), Props, addInputParam(), addOutputParam() (+10 more)

### Community 39 - "useVaultMutations"
Cohesion: 0.14
Nodes (6): VaultsAPI, useVaultMutations(), DangerZoneText(), JoinVaultPage(), Props, VaultDangerZone()

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.23
Nodes (15): SelectOption(), execute(), fetchBoards(), fetchCards(), fetchLists(), TrelloBoardLite, TrelloCardLite, TrelloListLite (+7 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.24
Nodes (8): MAIL_DOMAIN_ICONS, mailProviderIcon(), AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps, DialogHead(), DialogHeadProps

### Community 42 - "SchemaProperty"
Cohesion: 0.15
Nodes (13): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleSheetsSpreadsheetSection(), SpreadsheetRow(), ConnectedContent() (+5 more)

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.15
Nodes (14): Input(), FormField(), Props, Props, S3ToggleFields(), buildEmailRequest(), EmailAddDialog(), buildEmailRequest() (+6 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.12
Nodes (6): KNOWN_STATES, LoginRelayScreen(), LoginState, Props, IWorkbenchService, WorkbenchService

### Community 46 - "Tracts.ts"
Cohesion: 0.10
Nodes (21): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, framer-motion (+13 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.07
Nodes (25): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, GetS3InstanceResponse, ListS3Instances, ListS3InstancesRequest (+17 more)

### Community 48 - "dependencies"
Cohesion: 0.14
Nodes (13): useTracts, sleep(), Props, Step, StepDraft, StepPickerDialog(), Props, TractRunTimeline() (+5 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.31
Nodes (13): PublicDocsNoteItem, DocsFolderNode(), DocsFolderNodeProps, DocsSidebar(), DocsSidebarProps, DocsTreeItem(), buildFolderTree(), byName() (+5 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.16
Nodes (15): ActionStep, ConditionStep, GroupStep, LlmCallStep, ParallelStep, ScriptStep, TractCondition, TractDefinition (+7 more)

### Community 51 - "AuthMiddleware"
Cohesion: 0.09
Nodes (18): apiPrefix(), csrfHeader(), getCsrfToken(), InitReq, TelegramLoginResponse, AppConfigState, useAppConfig, pingServer() (+10 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.20
Nodes (23): appendStep(), branchArray(), buildStepFromDraft(), collapseThinParallels(), collectAllStepIds(), generateStepId(), hasChildren(), insertBlockAfter() (+15 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.07
Nodes (29): AdminUsersAPI, ArtelUserDetails, ArtelUserEntry, ChangeArtelUserPassword, ChangeArtelUserPasswordRequest, ChangeArtelUserPasswordResponse, CreateArtelUser, CreateArtelUserRequest (+21 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.08
Nodes (30): useIsMobileNav(), applyTheme(), Theme, useTheme(), BrandMarkIcon(), ConnectionsIcon(), base, NavIconProps (+22 more)

### Community 55 - "TractCanvasBuilder.tsx"
Cohesion: 0.11
Nodes (17): ArtelUI Frontend Rules, Async style, Buttons, Component hierarchy, Component Structure, CSS Modules, Dialog shells must scroll internally, Error and Confirmation Handling (+9 more)

### Community 56 - "User.ts"
Cohesion: 0.09
Nodes (27): AdminSubscriptionsAPI, EffectiveSubscriptionView, GetUserSubscription, GetUserSubscriptionRequest, GetUserSubscriptionResponse, ListSubscriptionPlans, ListSubscriptionPlansRequest, ListSubscriptionPlansResponse (+19 more)

### Community 57 - "Errors.ts"
Cohesion: 0.17
Nodes (6): ErrorReason, Errors, GrpcError, GrpcErrorDetail, ServiceError, ServiceErrorOption

### Community 58 - "NoteViewer.tsx"
Cohesion: 0.24
Nodes (8): ArtelMark(), ArtelMarkProps, ContentSegment, NoteContent(), NoteContentProps, parseWikiLinks(), WikiChip(), WikiChipProps

### Community 59 - "InviteLinksSection.tsx"
Cohesion: 0.24
Nodes (10): TaskTrackerTableHead(), DisplayTaskTrackerTables(), TrelloTableWidget(), RESULT_VIEW_WIDGETS, ResultViewWidgetEntry, ResultViewWidgetProps, isPlainObject(), TaskTrackerTable (+2 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.09
Nodes (19): ExternalProvider, ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), ProviderChip(), PROVIDER_CHIP_CLASS (+11 more)

### Community 61 - "toTract"
Cohesion: 0.26
Nodes (6): McpLogin(), McpLoginProps, VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.29
Nodes (7): connectionLabel(), LlmCallCard(), Props, ConnectionPicker(), NodeChips(), LlmConnectionStep(), Props

### Community 63 - "VaultCard.tsx"
Cohesion: 0.18
Nodes (12): AuthAPI, LoginRequest, PasswordLoginForm(), PasswordLoginFormProps, RegisterForm(), RegisterFormProps, LoginContentProps, Mode (+4 more)

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.05
Nodes (53): Props, RunButton(), Props, RunStatusBadge(), formatRelative(), Props, RunStatusDot(), LogicCell() (+45 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.16
Nodes (15): ResizeHandle(), ResizeHandleProps, buildLogLines(), formatDuration(), formatTime(), LogLineKind, stepMeta(), RunLog() (+7 more)

### Community 67 - "AuthAPI"
Cohesion: 0.21
Nodes (16): ImportConflictAction, ImportResolution, commitImportAndRefresh(), deleteFolderAndRefresh(), moveEntryAndRefresh(), NotesState, remapSelectedPath(), requireVaultId() (+8 more)

### Community 68 - "connectionLabel"
Cohesion: 0.23
Nodes (8): HeroSegment(), HeroSegmentProps, ToolboxPage(), NewTractButton(), Props, ContentSegment(), TractCanvasListPage(), TractCanvasCard()

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.27
Nodes (6): ByokTabIcon(), CommunityTabIcon(), ExternalConnectionsTabIcon(), ConnectionsPage(), ConnectionsTab, resolveTab()

### Community 70 - "compilerOptions"
Cohesion: 0.18
Nodes (4): CouchInstancesAPI, InstancesActionBar(), InstancesActionBarProps, InstancesTab()

### Community 71 - "McpKeys.ts"
Cohesion: 0.27
Nodes (7): GetCouchInstanceResponse, InstanceRow(), InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps

### Community 72 - "ResultView.tsx"
Cohesion: 0.13
Nodes (22): STEP_SCREENS, AddTriggerDialogContext, AddTriggerDialogState, AddTriggerStep, emptySchemaField(), FIELD_TYPES, fieldsToSchemaNode(), fieldsToSchemaProperties() (+14 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.09
Nodes (35): Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), boardListLabel(), Props, RoadmapCanvasNode() (+27 more)

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.22
Nodes (9): BinaryStorageToggle(), Props, Props, PublishSlugForm(), slugify(), validateSlug(), Props, PublishToggle() (+1 more)

### Community 75 - "dialog-scrollable.js"
Cohesion: 0.46
Nodes (7): allRules(), dialogScrollable(), directDeclsOf(), findScrollTarget(), isOverflowY(), messages, meta

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.14
Nodes (15): VaultInviteItem, VaultMemberInfo, vaultsQueryKey, CreateInviteLinkDialog(), Props, InviteRow(), Props, Props (+7 more)

### Community 77 - "RunTractDialog.tsx"
Cohesion: 0.24
Nodes (8): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), isJsonValue(), TaskTrackerCell(), stringifyCell()

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.27
Nodes (6): TemplateSource, Props, TemplateInput(), CONDITION_OPS, ConditionRowProps, Props

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.28
Nodes (4): AdminSystemSettingsAPI, AuthMethodsSection(), AuthMethodsSectionProps, SettingsTab()

### Community 80 - "CardMeta.tsx"
Cohesion: 0.15
Nodes (17): useDialogKeyboard(), TractTemplatesState, useTractTemplates, BrowseTemplatesDialog(), Props, TemplateRow(), Props, ListScreen() (+9 more)

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.35
Nodes (7): AdminPage(), resolveTab(), Tab, AdminHero(), AdminHeroProps, TabBar(), TabBarProps

### Community 93 - "VaultCardHeader.tsx"
Cohesion: 0.20
Nodes (7): DialogHead(), DialogHeadProps, NotConnectedContentProps, ConnectionDetailDialog(), PROVIDER_CONFIG, PROVIDER_KEY, ProviderConfig

### Community 94 - "AuthMiddleware"
Cohesion: 0.25
Nodes (6): CopyIcon(), ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon(), TreeItemProps

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.16
Nodes (10): BreadcrumbBarProps, Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), ErrorDotIcon(), PencilIcon(), SpinnerIcon() (+2 more)

### Community 96 - "UsersTab.tsx"
Cohesion: 0.22
Nodes (6): dompurify, PublicDocsAPI, DocsNoteViewer(), DocsNoteViewerProps, DocsPage(), usePublicDocs()

### Community 97 - "AdminCouchAPI"
Cohesion: 0.31
Nodes (5): DocsTreeItemProps, ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon()

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.13
Nodes (14): useVaults(), CardHeader(), Props, CardMeta(), formatDate(), Props, VaultChipDisplayProps, VaultField() (+6 more)

### Community 100 - "RunLog.tsx"
Cohesion: 0.23
Nodes (6): ToolDetail(), GenericToolIcon(), TODO: placeholder glyph for tool actions without a dedicated icon yet (non-smtp/, ImapIcon(), SmtpIcon(), ToolRow()

### Community 101 - "package.json"
Cohesion: 0.06
Nodes (31): DeleteDockerHost, DeleteDockerHostRequest, DeleteDockerHostResponse, DockerHostsAPI, GetDockerHost, GetDockerHostRequest, GetDockerHostResponse, ListDockerHosts (+23 more)

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.31
Nodes (6): useTemplateConnections(), UseTemplateConnectionsResult, InstantiateTemplateDialog(), ConnectionRequirement, ConnectionRequirementKind, requiredConnections()

### Community 103 - "CouchInstancesAPI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 104 - ".getRun"
Cohesion: 0.09
Nodes (32): useExternalConnections, useBakeError(), TriggerPanel(), DetailScreen(), AnthropicCheckButton(), AnthropicCheckButtonProps, CheckStatus, ConnectedContentProps (+24 more)

### Community 105 - "UserList.tsx"
Cohesion: 0.33
Nodes (4): CredentialRow(), CredentialRowProps, ConnectGenericDialog(), CredentialField

### Community 106 - "Vaults.ts"
Cohesion: 0.17
Nodes (5): AdminCouchAPI, ChangePasswordDialog(), ChangePasswordDialogProps, ManageAccessDialog(), UsersTab()

### Community 107 - "AuthFetchInterceptor.ts"
Cohesion: 0.43
Nodes (5): ResultViewMode, ViewModeToggle(), getResultViewWidget(), ResultView(), tryParseJson()

### Community 108 - "EmailCheckButton.tsx"
Cohesion: 0.32
Nodes (5): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps, ManageAccessDialogProps

### Community 109 - "GitlabCheckButton.tsx"
Cohesion: 0.50
Nodes (3): ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps

### Community 110 - "TrelloCheckButton.tsx"
Cohesion: 0.48
Nodes (4): CouchUserEntry, UserListProps, UserRow(), UserRowProps

### Community 111 - "TopbarDrawerCloseButton.tsx"
Cohesion: 0.19
Nodes (9): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), ContentSegment(), ContentSegmentProps, Props (+1 more)

### Community 112 - "InsertConflictDialog.tsx"
Cohesion: 0.50
Nodes (3): ImapOperation, ImapToolAction, ImapActionView()

### Community 113 - "TopbarMobileTrigger.tsx"
Cohesion: 0.50
Nodes (3): AccountsSection(), AccountsSectionProps, TrelloConnectionRow()

### Community 117 - "requiredConnections.ts"
Cohesion: 0.22
Nodes (12): useWorkbench(), useWorkbenchMutations(), workbenchQueryKey(), AuthMode, PickAuthModeScreen(), Props, Props, STATUS_CLASSES (+4 more)

### Community 119 - "ConnectedContent.tsx"
Cohesion: 0.22
Nodes (6): clearCsrfCookie(), queryClient, forceLogout(), originalFetch, refreshTokens(), SKIP_REFRESH_PATHS

### Community 120 - "KebabMenu.tsx"
Cohesion: 0.15
Nodes (14): RegistrationMode, SetupWizardAPI, RegistrationModeSection(), RegistrationModeSectionProps, AuthMethodsScreen(), OPTIONS, RegistrationModeScreen(), TokenEntryScreen() (+6 more)

### Community 122 - "RoadmapCanvasNode.tsx"
Cohesion: 0.22
Nodes (9): NoteMode, DesktopNotesShellProps, VaultOption, MobileNotesShellProps, NoteViewer(), NoteViewerProps, SaveStatus, UseAutosaveOptions (+1 more)

### Community 123 - "Tabs.tsx"
Cohesion: 0.20
Nodes (11): usePortrait(), Props, RenameDialog(), useAutosave(), NotesPage(), buildNotesUrl(), decodeNotePath(), encodeNotePath() (+3 more)

### Community 126 - "ParamsList.tsx"
Cohesion: 0.18
Nodes (9): McpToolInfo, SmtpOperation, SmtpToolAction, ToolParamDef, ParamRow(), ParamsList(), RunScreens(), SmtpActionView() (+1 more)

### Community 127 - "ImportZipDialog.tsx"
Cohesion: 0.28
Nodes (5): DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, ImportZipDialog(), Props

## Knowledge Gaps
- **802 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+797 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `cn` connect `cn` to `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Tracts.ts`, `TractIcons.tsx`, `useDialog`, `tractCanvasLayout.ts`, `VaultItem`, `NotesSidebar.tsx`, `ConnectionDetailDialog.tsx`, `index.ts`, `CreateNoteDialog.tsx`, `SchemaProperty`, `McpAuthPage.tsx`, `ProviderIcon.tsx`, `Topbar.tsx`, `TractBlockPicker.tsx`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `scripts`, `AuthAPI`, `connectionLabel`, `ResultView.tsx`, `MembersSection.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `RunTractDialog.tsx`, `DbAccessList.tsx`, `GoogleSheetsSpreadsheetSection.tsx`, `AuthMiddleware`, `AdminCouchAPI`, `package.json`, `.getRun`, `AuthFetchInterceptor.ts`, `requiredConnections.ts`, `KebabMenu.tsx`, `ParamsList.tsx`, `ImportZipDialog.tsx`?**
  _High betweenness centrality (0.110) - this node is a cross-community bridge._
- **Why does `useDialog` connect `ToolboxPage.tsx` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `TractsService`, `TractIcons.tsx`, `compilerOptions`, `grpcErrors.ts`, `VaultItem`, `Router.tsx`, `TractStepTree.tsx`, `BreadcrumbBar.tsx`, `ConnectionDetailDialog.tsx`, `CreateNoteDialog.tsx`, `SchemaProperty`, `McpAuthPage.tsx`, `useErrorToast.ts`, `dependencies`, `AuthMiddleware`, `admin_users.pb.ts`, `User.ts`, `TractBlockPicker.tsx`, `toTract`, `NotesPage.tsx`, `VaultCard.tsx`, `TractCanvasTopBar.tsx`, `AuthAPI`, `connectionLabel`, `compilerOptions`, `McpKeys.ts`, `ResultView.tsx`, `LinkScreen.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `CardMeta.tsx`, `VaultCardHeader.tsx`, `TractCanvasLogPanel.tsx`, `package.json`, `MobileDrawer.tsx`, `.getRun`, `UserList.tsx`, `Vaults.ts`, `EmailCheckButton.tsx`, `TrelloCheckButton.tsx`, `Tabs.tsx`, `ImportZipDialog.tsx`?**
  _High betweenness centrality (0.103) - this node is a cross-community bridge._
- **Why does `useUser` connect `BreadcrumbBar.tsx` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `ExternalConnectionInfo`, `couch_instances.pb.ts`, `TractsService`, `Dialog.ts`, `TractIcons.tsx`, `ToolboxPage.tsx`, `grpcErrors.ts`, `VaultItem`, `McpKeysAPI`, `index.ts`, `useVaultMutations`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `useErrorToast.ts`, `dependencies`, `AuthMiddleware`, `admin_users.pb.ts`, `Topbar.tsx`, `User.ts`, `toTract`, `VaultCard.tsx`, `AuthAPI`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `compilerOptions`, `McpKeys.ts`, `LinkScreen.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `VaultDangerZone.tsx`, `CardMeta.tsx`, `TractCanvasLogPanel.tsx`, `package.json`, `Vaults.ts`, `EmailCheckButton.tsx`, `TrelloCheckButton.tsx`, `requiredConnections.ts`, `ConnectedContent.tsx`?**
  _High betweenness centrality (0.078) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _808 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.02247191011235955 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0797872340425532 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.03125 - nodes in this community are weakly interconnected._