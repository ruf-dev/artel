# Graph Report - ArtelUI  (2026-07-25)

## Corpus Check
- 539 files · ~144,434 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2465 nodes · 6062 edges · 122 communities (113 shown, 9 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 22 edges (avg confidence: 0.68)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `be92ba1b`
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
- ConnectionFilterRow.tsx
- TrelloCheckButton.tsx
- requiredConnections.ts
- NotConnectedContent.tsx
- ConnectedContent.tsx

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 186 edges
2. `cn` - 150 edges
3. `useBakeError()` - 126 edges
4. `useUser` - 98 edges
5. `useExternalConnections` - 61 edges
6. `TractStep` - 46 edges
7. `TractTool` - 42 edges
8. `MomCandidate` - 41 edges
9. `ExternalConnectionInfo` - 38 edges
10. `useMcpKeys` - 36 edges

## Surprising Connections (you probably didn't know these)
- `NoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/notes/components/NoteViewer/NoteViewer.tsx → package.json
- `Props` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx → src/app/api/artel/external_connections.pb.ts
- `AccountsSectionProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/AccountsSection.tsx → src/app/api/artel/external_connections.pb.ts
- `EmailConnectionRowProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/components/EmailConnectionRow/EmailConnectionRow.tsx → src/app/api/artel/external_connections.pb.ts
- `Props` --references--> `McpKeyInfo`  [EXTRACTED]
  src/widgets/McpKeyCard/McpKeyCard.tsx → src/app/api/artel/mcp_keys.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/workbench/WorkbenchPage.tsx -> src/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx -> src/app/routing/Router.tsx`
- 5-file cycle: `src/app/routing/Router.tsx -> src/pages/tract-templates/TractTemplatesListPage.tsx -> src/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx -> src/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx -> src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx -> src/app/routing/Router.tsx`

## Communities (122 total, 9 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.02
Nodes (88): Absent, BaseTractStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger, CreateTriggerRequest, CreateTriggerResponse (+80 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (28): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+20 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.03
Nodes (57): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+49 more)

### Community 3 - "useBakeError"
Cohesion: 0.18
Nodes (11): useIsMobileNav(), applyTheme(), Theme, useTheme(), BrandMarkIcon(), TopbarBrand(), TopbarMobileDrawer(), TopbarThemeToggle() (+3 more)

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.04
Nodes (54): Absent, AddAnthropicConnection, AddAnthropicConnectionResponse, AddEmailConnection, AddEmailConnectionResponse, AddGenericConnection, AddGenericConnectionRequest, AddGenericConnectionResponse (+46 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.04
Nodes (46): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CommunityConnectorInfo, CreateMcpKey (+38 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.14
Nodes (21): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenArrayItems() (+13 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.11
Nodes (24): SetTractsState, triggerSourcesQueryKey, triggersQueryKey, useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor() (+16 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.13
Nodes (37): ScriptLanguage, Props, Props, TractStepTreeProps, ActionBody(), Props, CONDITION_OPS, ConditionBody() (+29 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.09
Nodes (22): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+14 more)

### Community 12 - "TractsService"
Cohesion: 0.06
Nodes (28): TractsAPI, TractsState, safeParseJson(), Props, Props, Deps, ITractsService, TractsService (+20 more)

### Community 13 - "Dialog.ts"
Cohesion: 0.12
Nodes (18): connectionLabel(), ConnectionStep(), Props, LlmConnectionStep(), Props, LlmCallCard(), ConnectionSection(), Props (+10 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.08
Nodes (26): cn, trimScope(), KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it, TODO: chures has no tab primitive yet, drop this wrapper once it does, Tabs() (+18 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.11
Nodes (24): useBakeError(), Input(), FormField(), Props, Props, S3ToggleFields(), S3InstanceFormDialog(), S3InstancesActionBar() (+16 more)

### Community 17 - "compilerOptions"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.18
Nodes (9): ImapOperation, ImapToolAction, McpToolInfo, ToolParamDef, ImapActionView(), ParamRow(), ParamsList(), RunScreens() (+1 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.15
Nodes (17): useWorkbench(), useWorkbenchMutations(), workbenchQueryKey(), KNOWN_STATES, LoginRelayScreen(), LoginState, Props, AuthMode (+9 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.13
Nodes (15): BadRequestDetail, DetailType, DetailTypeName, ErrorInfoDetail, FieldViolation, getDetail(), getFieldViolations(), GrpcErrorDetail (+7 more)

### Community 21 - "useDialog"
Cohesion: 0.20
Nodes (9): isNonEmptyBranch(), JsonBlock(), Props, Props, statusClass(), StepRow(), Props, Props (+1 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.14
Nodes (21): useServerStatus(), HomeLayout(), Path, Router(), routes, HeroSegment(), HeroSegmentProps, VaultCardHeader() (+13 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.10
Nodes (28): NodeChips(), ConnectorPath(), ConnectorPathProps, ParallelBoxes(), Props, Props, TractCanvasArea(), useTractCanvasDrag() (+20 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (24): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+16 more)

### Community 25 - "VaultItem"
Cohesion: 0.16
Nodes (13): VaultItem, CardChips(), Props, Props, VaultChip(), ExpertSettingsDrawer(), Props, Props (+5 more)

### Community 26 - "devDependencies"
Cohesion: 0.10
Nodes (21): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+13 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.12
Nodes (6): CreateMcpKeyResponse, McpKeyInfo, McpKeysAPI, McpKeysState, IMcpKeysService, McpKeysService

### Community 29 - "Router.tsx"
Cohesion: 0.23
Nodes (15): SelectOption(), execute(), fetchBoards(), fetchCards(), fetchLists(), TrelloBoardLite, TrelloCardLite, TrelloListLite (+7 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.18
Nodes (11): DialogHead(), DialogHeadProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), AddConnectionScreen(), MainScreen(), SelectConnectionScreen() (+3 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.12
Nodes (14): BreadcrumbBarProps, Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), ErrorDotIcon(), PencilIcon(), SpinnerIcon() (+6 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.21
Nodes (8): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, useNotesSearchQuery()

### Community 33 - "useExternalConnections"
Cohesion: 0.14
Nodes (18): paramCompletionSource(), scriptEditorTheme, scriptHighlightStyle, Props, ScriptCodeSection(), Props, addInputParam(), addOutputParam() (+10 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.14
Nodes (11): SuggestionList(), SuggestionListProps, CreateNoteDialog(), Props, ArtelLogoIcon(), ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps (+3 more)

### Community 35 - "notes.pb.ts"
Cohesion: 0.06
Nodes (33): CheckImportConflicts, CheckImportConflictsRequest, CheckImportConflictsResponse, CommitImport, CommitImportRequest, CommitImportResponse, DeleteFolder, DeleteFolderRequest (+25 more)

### Community 36 - "Notes.ts"
Cohesion: 0.11
Nodes (5): b64Decode(), ImportResolution, NotesAPI, INotesService, NotesService

### Community 37 - "index.ts"
Cohesion: 0.22
Nodes (8): ListPrompts, ListPromptsRequest, ListPromptsResponse, PromptId, PromptItem, PromptsAPI, FastSetupDialog(), Props

### Community 38 - "s3_instances.pb.ts"
Cohesion: 0.11
Nodes (17): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, ListS3Instances, ListS3InstancesRequest, ListS3InstancesResponse (+9 more)

### Community 39 - "useVaultMutations"
Cohesion: 0.09
Nodes (8): VaultsAPI, useVaultMutations(), DangerZoneText(), CreateVaultDialog(), IWorkbenchService, WorkbenchService, Props, VaultDangerZone()

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.16
Nodes (15): ActionStep, ConditionStep, GroupStep, LlmCallStep, ParallelStep, ScriptStep, TractCondition, TractDefinition (+7 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.11
Nodes (17): ArtelUI Frontend Rules, Async style, Buttons, Component hierarchy, Component Structure, CSS Modules, Dialog shells must scroll internally, Error and Confirmation Handling (+9 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.22
Nodes (8): parseScopeList(), SCOPE_INFO, GoogleSheetsSpreadsheetSection(), SpreadsheetRow(), GoogleConnectionContent(), GoogleSheetsConnectionContent(), GapiWindow, useSpreadsheetPicker()

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.29
Nodes (9): useMcpKeys, AddTaskLinkDialog(), Props, RELATION_LABEL, RELATION_OPTIONS, RoadmapLinkTarget, WritableRelation, createRoadmapGraph() (+1 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.18
Nodes (13): useNotes, SidebarTopBar(), SidebarTopBarProps, VaultOption, NotesSidebar(), NotesSidebarProps, VaultOption, useFolderActions() (+5 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.19
Nodes (6): AdminCouchAPI, CouchUserEntry, ManageAccessDialog(), UserListProps, UserRow(), UserRowProps

### Community 47 - "useErrorToast.ts"
Cohesion: 0.18
Nodes (9): useVaults(), vaultsQueryKey, VaultChipDisplayProps, VaultField(), VaultFieldProps, VaultOptionList(), VaultOptionListProps, IVaultService (+1 more)

### Community 48 - "dependencies"
Cohesion: 0.09
Nodes (22): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, dompurify (+14 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.15
Nodes (14): MAIL_DOMAIN_ICONS, mailProviderIcon(), ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), PROVIDER_CHIP_CLASS (+6 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.15
Nodes (19): AddTriggerDialog(), STEP_SCREENS, AddTriggerDialogContext, AddTriggerDialogState, AddTriggerStep, emptySchemaField(), FIELD_TYPES, fieldsToSchemaNode() (+11 more)

### Community 51 - "AuthMiddleware"
Cohesion: 0.20
Nodes (10): apiPrefix(), InitReq, Options, TelegramLoginResponse, AppConfigState, useAppConfig, pingServer(), VaultSettingsSection() (+2 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.16
Nodes (27): Props, InsertConflictDialog(), Props, appendStep(), branchArray(), buildStepFromDraft(), collapseThinParallels(), collectAllStepIds() (+19 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.09
Nodes (24): AdminUsersAPI, ArtelUserDetails, ArtelUserEntry, GetArtelUser, GetArtelUserRequest, GetArtelUserResponse, GetUserSessions, GetUserSessionsRequest (+16 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.09
Nodes (29): TractTemplatesState, useTractTemplates, BrowseTemplatesDialog(), Props, TemplateRow(), Props, ListScreen(), Props (+21 more)

### Community 55 - "TractCanvasBuilder.tsx"
Cohesion: 0.20
Nodes (11): NoteMode, DesktopNotesShellProps, VaultOption, MobileNotesShell(), MobileNotesShellProps, VaultOption, NoteViewer(), NoteViewerProps (+3 more)

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
Nodes (8): AdminPage(), Tab, AdminHero(), AdminHeroProps, DockerApiTab(), TabBar(), TabBarProps, UsersTab()

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.08
Nodes (24): ExternalProvider, ProviderChip(), AnthropicIcon(), EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), TrelloIcon() (+16 more)

### Community 61 - "toTract"
Cohesion: 0.29
Nodes (5): McpLoginProps, VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.12
Nodes (21): ActionCard(), Props, InsertRow(), Props, Props, SchemaTree(), buildSourcesFor(), collectIdsFromRoot() (+13 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.19
Nodes (9): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), ContentSegment(), ContentSegmentProps, Props (+1 more)

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.05
Nodes (53): Props, RunButton(), Props, RunStatusBadge(), formatRelative(), Props, RunStatusDot(), LogicCell() (+45 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.24
Nodes (11): cap(), LogEntryRow(), Props, buildLogLines(), formatDuration(), formatTime(), LogLine, LogLineKind (+3 more)

### Community 67 - "AuthAPI"
Cohesion: 0.22
Nodes (15): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), moveEntryAndRefresh(), NotesState, remapSelectedPath(), requireVaultId(), ConflictRow() (+7 more)

### Community 68 - "connectionLabel"
Cohesion: 0.06
Nodes (43): DialogManager, useDialog, useDialogKeyboard(), useTracts, sleep(), Props, ToolStep(), Props (+35 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.14
Nodes (19): AddAnthropicConnectionRequest, AddEmailConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, CheckAnthropicConnectionRequest, CheckAnthropicConnectionResponse, ExternalConnectionInfo, Spreadsheet (+11 more)

### Community 70 - "compilerOptions"
Cohesion: 0.18
Nodes (4): CouchInstancesAPI, InstancesActionBar(), InstancesActionBarProps, InstancesTab()

### Community 71 - "McpKeys.ts"
Cohesion: 0.27
Nodes (7): GetCouchInstanceResponse, InstanceRow(), InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps

### Community 72 - "ResultView.tsx"
Cohesion: 0.17
Nodes (13): isJsonValue(), TaskTrackerCell(), TaskTrackerTableHead(), DisplayTaskTrackerTables(), TrelloTableWidget(), RESULT_VIEW_WIDGETS, ResultViewWidgetEntry, ResultViewWidgetProps (+5 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.09
Nodes (35): Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), boardListLabel(), Props, RoadmapCanvasNode() (+27 more)

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.22
Nodes (9): scripts, build, build:ui, dev, gen, lint, lint:css, lint:js (+1 more)

### Community 75 - "dialog-scrollable.js"
Cohesion: 0.46
Nodes (7): allRules(), dialogScrollable(), directDeclsOf(), findScrollTarget(), isOverflowY(), messages, meta

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.27
Nodes (6): VaultMemberInfo, Props, Props, RoleBadge(), MembersSection(), Props

### Community 77 - "RunTractDialog.tsx"
Cohesion: 0.10
Nodes (17): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, formatStartedAt(), RunTractDialog() (+9 more)

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.23
Nodes (9): VaultInviteItem, CreateInviteLinkDialog(), Props, InviteRow(), Props, Props, RoleOption(), InviteLinksSection() (+1 more)

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.31
Nodes (8): ResizeHandle(), ResizeHandleProps, clampHeight(), dotClass(), formatDate(), loadStoredHeight(), Props, TractCanvasLogPanel()

### Community 80 - "CardMeta.tsx"
Cohesion: 0.39
Nodes (5): McpConnectorInfo, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps, MainScreenProps

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.33
Nodes (4): GetS3InstanceResponse, Props, S3InstanceRow(), TestStatus

### Community 94 - "AuthMiddleware"
Cohesion: 0.13
Nodes (11): UserState, LoginContentProps, AuthService, IAuthService, Session, AuthMiddleware, clearLocalStorage(), fromLocalStorage() (+3 more)

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.44
Nodes (6): usePortrait(), useAutosave(), NotesPage(), buildNotesUrl(), decodeNotePath(), encodeNotePath()

### Community 96 - "UsersTab.tsx"
Cohesion: 0.43
Nodes (6): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind

### Community 97 - "AdminCouchAPI"
Cohesion: 0.22
Nodes (10): MomCandidate, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectionOptionListProps, AddConnectionScreenProps, ConnectionFilterRow() (+2 more)

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.28
Nodes (5): DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, ImportZipDialog(), Props

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.24
Nodes (7): CardHeader(), Props, CardMeta(), formatDate(), Props, McpKeyCard(), Props

### Community 100 - "RunLog.tsx"
Cohesion: 0.25
Nodes (5): queryClient, originalFetch, RefreshResponseBody, refreshTokens(), SKIP_REFRESH_PATHS

### Community 101 - "package.json"
Cohesion: 0.33
Nodes (5): name, private, trustedDependencies, type, version

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.23
Nodes (6): ToolDetail(), GenericToolIcon(), TODO: placeholder glyph for tool actions without a dedicated icon yet (non-smtp/, ImapIcon(), SmtpIcon(), ToolRow()

### Community 103 - "CouchInstancesAPI"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 104 - ".getRun"
Cohesion: 0.15
Nodes (14): useExternalConnections, ConnectionDetailDialog(), DialogHead(), DialogHeadProps, tokenAuthorizeUrl(), TrelloAddDialog(), COMING_SOON_CARDS, LLM_BYOK_PROVIDERS (+6 more)

### Community 105 - "UserList.tsx"
Cohesion: 0.32
Nodes (5): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps, ManageAccessDialogProps

### Community 106 - "Vaults.ts"
Cohesion: 0.25
Nodes (6): CopyIcon(), ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon(), TreeItemProps

### Community 107 - "AuthFetchInterceptor.ts"
Cohesion: 0.33
Nodes (6): CheckEmailConnectionRequest, CheckStatus, EmailCheckButton(), EmailCheckButtonProps, isGrpcError(), notRetryOnStatus()

### Community 108 - "EmailCheckButton.tsx"
Cohesion: 0.22
Nodes (10): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), ResultViewMode, ViewModeToggle(), getResultViewWidget() (+2 more)

### Community 109 - "GitlabCheckButton.tsx"
Cohesion: 0.33
Nodes (4): ArtelAPI, Version, VersionRequest, VersionResponse

### Community 110 - "TrelloCheckButton.tsx"
Cohesion: 0.18
Nodes (6): AuthAPI, UserErrors, HomePage(), getPreconditionViolations(), GrpcStatusError, isMissingSubscription()

### Community 112 - "InsertConflictDialog.tsx"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 113 - "TopbarMobileTrigger.tsx"
Cohesion: 0.50
Nodes (3): TopbarHamburgerIcon(), TopbarMobileTrigger(), TopbarMobileTriggerProps

### Community 115 - "ConnectionFilterRow.tsx"
Cohesion: 0.50
Nodes (4): CheckGitlabConnectionRequest, CheckStatus, GitlabCheckButton(), GitlabCheckButtonProps

### Community 116 - "TrelloCheckButton.tsx"
Cohesion: 0.50
Nodes (4): CheckTrelloConnectionRequest, CheckStatus, TrelloCheckButton(), TrelloCheckButtonProps

## Knowledge Gaps
- **705 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+700 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Tracts.ts`, `TractIcons.tsx`, `ToolboxPage.tsx`, `ConnectorChip.tsx`, `useDialog`, `ManageKeyDialog.tsx`, `tractCanvasLayout.ts`, `VaultItem`, `Router.tsx`, `BreadcrumbBar.tsx`, `NotesSidebar.tsx`, `ConnectionDetailDialog.tsx`, `index.ts`, `SchemaProperty`, `ProviderIcon.tsx`, `StepRow.tsx`, `admin_users.pb.ts`, `Topbar.tsx`, `InviteLinksSection.tsx`, `TractBlockPicker.tsx`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `scripts`, `AuthAPI`, `connectionLabel`, `ResultView.tsx`, `MembersSection.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `DbAccessList.tsx`, `VaultDangerZone.tsx`, `UsersTab.tsx`, `AdminCouchAPI`, `CouchInstancesAPI`, `.getRun`, `Vaults.ts`, `EmailCheckButton.tsx`, `TopbarMobileTrigger.tsx`?**
  _High betweenness centrality (0.142) - this node is a cross-community bridge._
- **Why does `useDialog` connect `connectionLabel` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `TractsService`, `Dialog.ts`, `TractIcons.tsx`, `ManageKeyDialog.tsx`, `VaultItem`, `Router.tsx`, `TractStepTree.tsx`, `ConnectionDetailDialog.tsx`, `useVaultMutations`, `SchemaProperty`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `Tracts.ts`, `ProviderIcon.tsx`, `StepRow.tsx`, `AuthMiddleware`, `tractSteps.ts`, `admin_users.pb.ts`, `Topbar.tsx`, `User.ts`, `TractBlockPicker.tsx`, `toTract`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `AuthAPI`, `compilerOptions`, `McpKeys.ts`, `RunTractDialog.tsx`, `DbAccessList.tsx`, `GoogleSheetsSpreadsheetSection.tsx`, `AuthMiddleware`, `ConnectForm.tsx`, `AdminCouchAPI`, `CouchInstancesAPI`, `TractCanvasLogPanel.tsx`, `.getRun`, `UserList.tsx`, `TrelloCheckButton.tsx`?**
  _High betweenness centrality (0.108) - this node is a cross-community bridge._
- **Why does `useUser` connect `ManageKeyDialog.tsx` to `TaskTrackersPage.tsx`, `useBakeError`, `addTriggerDialogContext.ts`, `couch_instances.pb.ts`, `TractsService`, `TractIcons.tsx`, `ConnectorChip.tsx`, `VaultItem`, `McpKeysAPI`, `Notes.ts`, `index.ts`, `useVaultMutations`, `McpAuthPage.tsx`, `Tracts.ts`, `useErrorToast.ts`, `AuthMiddleware`, `admin_users.pb.ts`, `Topbar.tsx`, `User.ts`, `InviteLinksSection.tsx`, `toTract`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `compilerOptions`, `McpKeys.ts`, `GoogleSheetsSpreadsheetSection.tsx`, `AuthMiddleware`, `RunLog.tsx`, `.getRun`, `UserList.tsx`, `AuthFetchInterceptor.ts`, `TrelloCheckButton.tsx`, `ConnectionFilterRow.tsx`, `TrelloCheckButton.tsx`?**
  _High betweenness centrality (0.065) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _710 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.02247191011235955 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07922705314009662 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.034482758620689655 - nodes in this community are weakly interconnected._