# Graph Report - ArtelUI  (2026-08-02)

## Corpus Check
- 575 files · ~150,041 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2622 nodes · 6394 edges · 120 communities (110 shown, 10 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 23 edges (avg confidence: 0.67)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ec00a453`
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
- NotConnectedContent.tsx
- ConnectedContent.tsx

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 200 edges
2. `cn` - 156 edges
3. `useBakeError()` - 136 edges
4. `useUser` - 100 edges
5. `useExternalConnections` - 64 edges
6. `TractStep` - 46 edges
7. `TractTool` - 42 edges
8. `MomCandidate` - 41 edges
9. `ExternalConnectionInfo` - 39 edges
10. `useMcpKeys` - 39 edges

## Surprising Connections (you probably didn't know these)
- `DocsNoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/docs/components/DocsNoteViewer/DocsNoteViewer.tsx → package.json
- `NoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/notes/components/NoteViewer/NoteViewer.tsx → package.json
- `EmailConnectionRowProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/components/EmailConnectionRow/EmailConnectionRow.tsx → src/app/api/artel/external_connections.pb.ts
- `AnthropicCheckButtonProps` --references--> `CheckAnthropicConnectionRequest`  [EXTRACTED]
  src/dialogs/ManageAnthropicDialog/components/AnthropicCheckButton/AnthropicCheckButton.tsx → src/app/api/artel/external_connections.pb.ts
- `MomCandidateCardProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/components/CandidateOptionList/components/MomCandidateCard/MomCandidateCard.tsx → src/app/api/artel/mcp_keys.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/workbench/WorkbenchPage.tsx -> src/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx -> src/app/routing/Router.tsx`
- 5-file cycle: `src/app/routing/Router.tsx -> src/pages/tract-templates/TractTemplatesListPage.tsx -> src/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx -> src/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx -> src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx -> src/app/routing/Router.tsx`

## Communities (120 total, 10 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.02
Nodes (80): Absent, BaseTractStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger, CreateTriggerRequest, CreateTriggerResponse (+72 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.06
Nodes (44): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+36 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.03
Nodes (63): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+55 more)

### Community 3 - "useBakeError"
Cohesion: 0.06
Nodes (38): GetVaultBySlug, GetVaultBySlugRequest, GetVaultBySlugResponse, PublicDocsAPI, PublicDocsGetNote, PublicDocsGetNoteRequest, PublicDocsGetNoteResponse, PublicDocsListFolders (+30 more)

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.03
Nodes (78): Absent, AddAnthropicConnection, AddAnthropicConnectionResponse, AddEmailConnection, AddEmailConnectionResponse, AddGenericConnection, AddGenericConnectionResponse, AddGitlabConnection (+70 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.04
Nodes (48): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+40 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.11
Nodes (28): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenArrayItems() (+20 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.16
Nodes (20): useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor(), providerLabel(), triggerChipLabel(), TriggerRow() (+12 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.12
Nodes (38): ScriptLanguage, Props, ToolStep(), Props, TractStepTreeProps, ActionBody(), Props, CONDITION_OPS (+30 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.10
Nodes (20): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+12 more)

### Community 12 - "TractsService"
Cohesion: 0.07
Nodes (4): TractsAPI, TractsService, toRun(), toRunStep()

### Community 13 - "Dialog.ts"
Cohesion: 0.22
Nodes (10): connectionLabel(), SelectOption(), ConnectionStep(), Props, ConnectionSection(), ConnectionOptionList(), ConnectionOptionListProps, ConnectionPicker() (+2 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.04
Nodes (63): VaultInviteItem, cn, DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, KebabMenu(), KebabMenuItem, Props (+55 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.06
Nodes (42): GetS3InstanceResponse, useBakeError(), FormField(), Props, Props, S3ToggleFields(), Props, S3InstanceFormDialog() (+34 more)

### Community 17 - "compilerOptions"
Cohesion: 0.16
Nodes (21): TractsState, sleep(), formatStartedAt(), Props, RunTractDialog(), Props, TractCanvasBuilder(), Deps (+13 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.13
Nodes (12): McpToolInfo, ToolParamDef, ParamRow(), ParamsList(), RunScreens(), coerceParams(), ToolDetail(), GenericToolIcon() (+4 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.11
Nodes (19): UserErrors, BadRequestDetail, DetailType, DetailTypeName, ErrorInfoDetail, FieldViolation, getDetail(), getFieldViolations() (+11 more)

### Community 21 - "useDialog"
Cohesion: 0.13
Nodes (17): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind, isNonEmptyBranch(), JsonBlock() (+9 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.12
Nodes (15): SetTractsState, triggerSourcesQueryKey, triggersQueryKey, useTracts, Props, Step, StepDraft, StepPickerDialog() (+7 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.09
Nodes (31): MomCandidate, ProviderIcon(), MomCandidateRow(), NodeChips(), ConnectionFilterRow(), ConnectionFilterRowProps, Args, ConnectorPath() (+23 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (24): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+16 more)

### Community 25 - "VaultItem"
Cohesion: 0.18
Nodes (11): VaultItem, Props, VaultCardHeader(), ExpertSettingsDrawer(), Props, Props, ExpertSettingsSection(), Props (+3 more)

### Community 26 - "devDependencies"
Cohesion: 0.05
Nodes (42): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+34 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.11
Nodes (6): CommunityConnectorInfo, CreateMcpKeyResponse, McpKeysAPI, McpKeysState, IMcpKeysService, McpKeysService

### Community 29 - "Router.tsx"
Cohesion: 0.12
Nodes (16): TractItem, TractRunItem, TractRunStepItem, TractTemplateItem, TractTemplateSummary, TractToolItem, TriggerItem, TriggerSourceItem (+8 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.18
Nodes (13): McpKeyInfo, DialogHead(), DialogHeadProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), MainScreen(), MainScreenProps (+5 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.12
Nodes (17): NoteMode, BreadcrumbBarProps, Mode, BreadcrumbPath(), BreadcrumbPathProps, DesktopNotesShellProps, VaultOption, CheckIcon() (+9 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.21
Nodes (8): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, useNotesSearchQuery()

### Community 33 - "useExternalConnections"
Cohesion: 0.22
Nodes (15): paramCompletionSource(), scriptEditorTheme, scriptHighlightStyle, Props, ScriptCodeSection(), Props, generatedFooter(), generatedHeader() (+7 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.20
Nodes (10): useNotes, ArtelLogoIcon(), ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps, MobileDrawer(), MobileDrawerProps, VaultOption (+2 more)

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
Cohesion: 0.10
Nodes (4): VaultsAPI, useVaultMutations(), IWorkbenchService, WorkbenchService

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.16
Nodes (15): ActionStep, ConditionStep, GroupStep, LlmCallStep, ParallelStep, ScriptStep, TractCondition, TractDefinition (+7 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.20
Nodes (14): LogicCell(), LogicCellProps, LogicSection(), Props, OptionCell(), ToolCell(), ToolCellProps, LOGIC_OPTIONS (+6 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.50
Nodes (6): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleConnectionContent()

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.11
Nodes (18): Input(), AddTaskLinkDialog(), Props, RELATION_LABEL, RELATION_OPTIONS, WritableRelation, CredentialRowProps, DialogHeadProps (+10 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.15
Nodes (13): PublicBadge(), SidebarTopBar(), SidebarTopBarProps, VaultOption, NotesSidebar(), NotesSidebarProps, VaultOption, useFolderActions() (+5 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.14
Nodes (14): useMcpKeys, useVaults(), CardHeader(), Props, Props, SearchInput(), VaultChipDisplayProps, VaultField() (+6 more)

### Community 48 - "dependencies"
Cohesion: 0.10
Nodes (21): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, framer-motion (+13 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.19
Nodes (11): MAIL_DOMAIN_ICONS, mailProviderIcon(), ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), PROVIDER_CHIP_CLASS (+3 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.16
Nodes (19): STEP_SCREENS, AddTriggerDialogContext, AddTriggerDialogState, AddTriggerStep, emptySchemaField(), FIELD_TYPES, fieldsToSchemaNode(), fieldsToSchemaProperties() (+11 more)

### Community 51 - "AuthMiddleware"
Cohesion: 0.13
Nodes (13): apiPrefix(), csrfHeader(), getCsrfToken(), InitReq, TelegramLoginResponse, AppConfigState, useAppConfig, pingServer() (+5 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.17
Nodes (25): InsertConflictDialog(), Props, appendStep(), branchArray(), buildStepFromDraft(), collapseThinParallels(), collectAllStepIds(), generateStepId() (+17 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.09
Nodes (24): AdminUsersAPI, ArtelUserDetails, ArtelUserEntry, GetArtelUser, GetArtelUserRequest, GetArtelUserResponse, GetUserSessions, GetUserSessionsRequest (+16 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.05
Nodes (48): useServerStatus(), useIsMobileNav(), applyTheme(), Theme, useTheme(), useWorkbench(), useWorkbenchMutations(), workbenchQueryKey() (+40 more)

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
Cohesion: 0.16
Nodes (11): dompurify, ArtelMark(), ArtelMarkProps, ContentSegment, NoteContent(), NoteContentProps, parseWikiLinks(), NoteViewer() (+3 more)

### Community 59 - "InviteLinksSection.tsx"
Cohesion: 0.06
Nodes (27): DeleteDockerHost, DeleteDockerHostRequest, DeleteDockerHostResponse, DockerHostsAPI, GetDockerHost, GetDockerHostRequest, GetDockerHostResponse, ListDockerHosts (+19 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.18
Nodes (8): AnthropicIcon(), EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), TrelloIcon(), TODO: placeholder glyph for providers without a dedicated brand icon yet - repla, UnknownProviderIcon()

### Community 61 - "toTract"
Cohesion: 0.29
Nodes (5): McpLoginProps, VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.11
Nodes (19): ActionCard(), CardHeader(), Props, InsertRow(), Props, LlmCallCard(), Props, Props (+11 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.19
Nodes (9): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), ContentSegment(), ContentSegmentProps, Props (+1 more)

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.18
Nodes (8): EditIcon(), EnvelopeIcon(), GlobeIcon(), base, PaperPlaneIcon(), PlusIcon(), SearchIcon(), WrenchIcon()

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.36
Nodes (7): buildLogLines(), formatDuration(), formatTime(), LogLineKind, stepMeta(), RunLog(), RunLogProps

### Community 67 - "AuthAPI"
Cohesion: 0.22
Nodes (15): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), moveEntryAndRefresh(), NotesState, remapSelectedPath(), requireVaultId(), ConflictRow() (+7 more)

### Community 68 - "connectionLabel"
Cohesion: 0.06
Nodes (40): DialogManager, useDialog, useDialogKeyboard(), useTractTemplates, HeroSegment(), HeroSegmentProps, LlmConnectionStep(), Props (+32 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.22
Nodes (15): AddAnthropicConnectionRequest, AddEmailConnectionRequest, AddGenericConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, CheckAnthropicConnectionRequest, CheckAnthropicConnectionResponse, Spreadsheet (+7 more)

### Community 71 - "McpKeys.ts"
Cohesion: 0.29
Nodes (6): GetCouchInstanceResponse, InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps

### Community 72 - "ResultView.tsx"
Cohesion: 0.13
Nodes (18): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), isJsonValue(), TaskTrackerCell(), TaskTrackerTableHead() (+10 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.09
Nodes (38): RoadmapLinkTarget, Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), boardListLabel(), Props (+30 more)

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.23
Nodes (8): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, Props, SchemaProperty

### Community 75 - "dialog-scrollable.js"
Cohesion: 0.46
Nodes (7): allRules(), dialogScrollable(), directDeclsOf(), findScrollTarget(), isOverflowY(), messages, meta

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.25
Nodes (5): VaultMemberInfo, vaultsQueryKey, Props, IVaultService, VaultService

### Community 77 - "RunTractDialog.tsx"
Cohesion: 0.14
Nodes (10): PRIMITIVE_TYPES, Props, ITEM_TYPES, PARAM_TYPES, ParamType, LogPanelBar(), LogPanelBarProps, CloseIcon() (+2 more)

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.22
Nodes (9): BinaryStorageToggle(), Props, Props, PublishSlugForm(), slugify(), validateSlug(), Props, PublishToggle() (+1 more)

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.31
Nodes (8): ResizeHandle(), ResizeHandleProps, clampHeight(), dotClass(), formatDate(), loadStoredHeight(), Props, TractCanvasLogPanel()

### Community 80 - "CardMeta.tsx"
Cohesion: 0.23
Nodes (8): TractTemplatesState, Props, TemplateRow(), Props, Props, Props, toTractTemplateSummary(), TractTemplateSummary

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.27
Nodes (5): BranchIcon(), ChatIcon(), CodeIcon(), ForkIcon(), LayersIcon()

### Community 93 - "VaultCardHeader.tsx"
Cohesion: 0.38
Nodes (7): OptIcon(), OptIconProps, OptText(), OptTextProps, OptionCellProps, IconProps, StepColor

### Community 94 - "AuthMiddleware"
Cohesion: 0.13
Nodes (4): AuthMiddleware, clearLocalStorage(), fromLocalStorage(), saveToLocalStorage()

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.20
Nodes (11): usePortrait(), Props, RenameDialog(), useAutosave(), NotesPage(), buildNotesUrl(), decodeNotePath(), encodeNotePath() (+3 more)

### Community 96 - "UsersTab.tsx"
Cohesion: 0.29
Nodes (5): ChevronRightIcon(), ManualTriggerIcon(), TrashIcon(), WebhookIcon(), Props

### Community 97 - "AdminCouchAPI"
Cohesion: 0.16
Nodes (12): McpConnectorInfo, CardChips(), Props, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectorRow() (+4 more)

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.67
Nodes (3): CardMeta(), formatDate(), Props

### Community 100 - "RunLog.tsx"
Cohesion: 0.32
Nodes (4): Props, RunButton(), Props, PlayIcon()

### Community 101 - "package.json"
Cohesion: 0.32
Nodes (3): addInputParam(), addOutputParam(), uniqueParamName()

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.43
Nodes (5): useTemplateConnections(), UseTemplateConnectionsResult, InstantiateTemplateDialog(), Props, ConnectionRequirement

### Community 103 - "CouchInstancesAPI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 104 - ".getRun"
Cohesion: 0.08
Nodes (24): useExternalConnections, ComingSoonCardProps, CredentialRow(), ConnectGenericDialog(), CredentialField, AnthropicCheckButton(), AnthropicCheckButtonProps, CheckStatus (+16 more)

### Community 105 - "UserList.tsx"
Cohesion: 0.40
Nodes (4): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps

### Community 106 - "Vaults.ts"
Cohesion: 0.25
Nodes (6): CopyIcon(), ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon(), TreeItemProps

### Community 107 - "AuthFetchInterceptor.ts"
Cohesion: 0.11
Nodes (17): CheckEmailConnectionRequest, CheckGitlabConnectionRequest, CheckStatus, EmailCheckButton(), EmailCheckButtonProps, ConnectedContent(), ConnectedContentProps, UseDefaultButton() (+9 more)

### Community 108 - "EmailCheckButton.tsx"
Cohesion: 0.50
Nodes (3): CouchUserEntry, UserListProps, UserRowProps

### Community 109 - "GitlabCheckButton.tsx"
Cohesion: 0.33
Nodes (4): ArtelAPI, Version, VersionRequest, VersionResponse

### Community 110 - "TrelloCheckButton.tsx"
Cohesion: 0.26
Nodes (8): AuthAPI, UserState, LoginContentProps, AuthService, IAuthService, Session, StoredAuth, UserInfo

### Community 112 - "InsertConflictDialog.tsx"
Cohesion: 0.50
Nodes (3): DangerZoneText(), Props, VaultDangerZone()

### Community 113 - "TopbarMobileTrigger.tsx"
Cohesion: 0.60
Nodes (4): formatRelative(), Props, RunStatusDot(), TractLastRun

## Knowledge Gaps
- **759 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+754 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `useDialog` connect `connectionLabel` to `TaskTrackersPage.tsx`, `external_connections.pb.ts`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `TractIcons.tsx`, `compilerOptions`, `ManageKeyDialog.tsx`, `tractCanvasLayout.ts`, `TractStepTree.tsx`, `ConnectionDetailDialog.tsx`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `useErrorToast.ts`, `ProviderIcon.tsx`, `StepRow.tsx`, `AuthMiddleware`, `tractSteps.ts`, `admin_users.pb.ts`, `User.ts`, `InviteLinksSection.tsx`, `toTract`, `NotesPage.tsx`, `AuthAPI`, `S3InstanceFormDialog.tsx`, `MembersSection.tsx`, `DbAccessList.tsx`, `CardMeta.tsx`, `ConnectForm.tsx`, `UsersTab.tsx`, `MobileDrawer.tsx`, `.getRun`, `AuthFetchInterceptor.ts`, `TrelloCheckButton.tsx`?**
  _High betweenness centrality (0.120) - this node is a cross-community bridge._
- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `Tracts.ts`, `TractIcons.tsx`, `ToolboxPage.tsx`, `useDialog`, `tractCanvasLayout.ts`, `VaultItem`, `NotesSidebar.tsx`, `index.ts`, `SchemaProperty`, `McpAuthPage.tsx`, `useErrorToast.ts`, `ProviderIcon.tsx`, `StepRow.tsx`, `admin_users.pb.ts`, `Topbar.tsx`, `NotesPage.tsx`, `AuthAPI`, `connectionLabel`, `ResultView.tsx`, `MembersSection.tsx`, `VaultDangerZone.tsx`, `UsersTab.tsx`, `Vaults.ts`, `TopbarMobileTrigger.tsx`?**
  _High betweenness centrality (0.111) - this node is a cross-community bridge._
- **Why does `useBakeError()` connect `TractIcons.tsx` to `TaskTrackersPage.tsx`, `external_connections.pb.ts`, `addTriggerDialogContext.ts`, `compilerOptions`, `ToolboxPage.tsx`, `ManageKeyDialog.tsx`, `tractCanvasLayout.ts`, `VaultItem`, `ConnectionDetailDialog.tsx`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `admin_users.pb.ts`, `Topbar.tsx`, `User.ts`, `InviteLinksSection.tsx`, `scripts`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `MembersSection.tsx`, `DbAccessList.tsx`, `ConnectForm.tsx`, `UsersTab.tsx`, `MobileDrawer.tsx`, `.getRun`, `AuthFetchInterceptor.ts`, `InsertConflictDialog.tsx`?**
  _High betweenness centrality (0.070) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _765 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.024691358024691357 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06010230179028133 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.03125 - nodes in this community are weakly interconnected._