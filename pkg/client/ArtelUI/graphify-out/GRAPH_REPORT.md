# Graph Report - ArtelUI  (2026-08-05)

## Corpus Check
- 586 files · ~152,288 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2700 nodes · 6563 edges · 124 communities (115 shown, 9 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 23 edges (avg confidence: 0.67)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `696f5966`
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
- ConnectedContent.tsx
- KebabMenu.tsx
- Tabs.tsx
- Textarea.tsx
- BYOKSection.tsx
- GitlabCheckButton.tsx

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 202 edges
2. `cn` - 158 edges
3. `useBakeError()` - 148 edges
4. `useUser` - 104 edges
5. `useExternalConnections` - 64 edges
6. `TractStep` - 46 edges
7. `TractTool` - 42 edges
8. `MomCandidate` - 41 edges
9. `ExternalConnectionInfo` - 39 edges
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
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/workbench/WorkbenchPage.tsx -> src/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx -> src/app/routing/Router.tsx`
- 5-file cycle: `src/app/routing/Router.tsx -> src/pages/tract-templates/TractTemplatesListPage.tsx -> src/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx -> src/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx -> src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx -> src/app/routing/Router.tsx`

## Communities (124 total, 9 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.03
Nodes (79): Absent, BaseTractStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger, CreateTriggerResponse, DeleteTract (+71 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (27): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+19 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.03
Nodes (63): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+55 more)

### Community 3 - "useBakeError"
Cohesion: 0.05
Nodes (45): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, dompurify (+37 more)

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.04
Nodes (53): Absent, AddAnthropicConnection, AddAnthropicConnectionResponse, AddEmailConnection, AddEmailConnectionResponse, AddGenericConnection, AddGenericConnectionResponse, AddGitlabConnection (+45 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.04
Nodes (48): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+40 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.14
Nodes (23): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenArrayItems() (+15 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.29
Nodes (7): categoryLabel(), PROVIDER_ENUM_BY_KEY, WebhookPicker(), WebhookPickerProps, PresetDetails(), SourcePicker(), TriggerSource

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.09
Nodes (3): ExternalConnectionsAPI, GoogleOAuthCallbackPage(), ExternalConnectionsService

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.14
Nodes (35): MomCandidate, ScriptLanguage, Props, TractStepTreeProps, Props, Props, DangerZone(), Props (+27 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.10
Nodes (20): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+12 more)

### Community 12 - "TractsService"
Cohesion: 0.07
Nodes (5): TractsAPI, TractsService, toTract(), toTractTemplate(), toTrigger()

### Community 13 - "Dialog.ts"
Cohesion: 0.08
Nodes (34): useServerStatus(), HomeLayout(), Path, Router(), routes, HeroSegment(), HeroSegmentProps, VaultCardHeader() (+26 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.07
Nodes (36): cn, KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it, TODO: chures has no tab primitive yet, drop this wrapper once it does, Tabs(), TabsProps (+28 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.10
Nodes (23): useBakeError(), ConnectionDetailDialog(), CreateInviteLinkDialog(), InviteLinksSection(), Props, InstanceRow(), AuthMethodsSection(), AuthMethodsSectionProps (+15 more)

### Community 17 - "compilerOptions"
Cohesion: 0.13
Nodes (18): SetTractsState, TractsState, triggerSourcesQueryKey, triggersQueryKey, sleep(), Props, Props, Deps (+10 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.15
Nodes (12): McpToolInfo, Props, SearchInput(), MomCandidateRow(), ContentSegment(), RunScreens(), ResultViewMode, ViewModeToggle() (+4 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.11
Nodes (20): UserErrors, BadRequestDetail, DetailType, DetailTypeName, ErrorInfoDetail, FieldViolation, getDetail(), getFieldViolations() (+12 more)

### Community 21 - "useDialog"
Cohesion: 0.13
Nodes (15): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind, isNonEmptyBranch(), JsonBlock() (+7 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.11
Nodes (18): Absent, BaseCompleteSetupRequest, CompleteSetup, CompleteSetupRequest, CompleteSetupResponse, GetStatus, GetStatusRequest, GetStatusResponse (+10 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.11
Nodes (25): ConnectorPath(), ConnectorPathProps, ParallelBoxes(), Props, Props, TractCanvasArea(), useTractCanvasDrag(), cap() (+17 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (24): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+16 more)

### Community 25 - "VaultItem"
Cohesion: 0.19
Nodes (9): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), ContentSegment(), ContentSegmentProps, Props (+1 more)

### Community 26 - "devDependencies"
Cohesion: 0.05
Nodes (42): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+34 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.11
Nodes (9): CommunityConnectorInfo, CreateMcpKeyResponse, McpConnectorInfo, McpKeyInfo, McpKeysAPI, McpKeysState, MainScreenProps, IMcpKeysService (+1 more)

### Community 29 - "Router.tsx"
Cohesion: 0.12
Nodes (6): KNOWN_STATES, LoginRelayScreen(), LoginState, Props, IWorkbenchService, WorkbenchService

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.15
Nodes (14): DialogHead(), DialogHeadProps, VaultField(), VaultFieldProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), AddConnectionScreen() (+6 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.25
Nodes (6): CopyIcon(), ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon(), TreeItemProps

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.14
Nodes (12): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, PublicBadge() (+4 more)

### Community 33 - "useExternalConnections"
Cohesion: 0.11
Nodes (21): paramCompletionSource(), scriptEditorTheme, scriptHighlightStyle, Props, ScriptCodeSection(), ITEM_TYPES, PARAM_TYPES, ParamType (+13 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.11
Nodes (19): useNotes, CreateNoteDialog(), Props, ArtelLogoIcon(), ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps, MobileDrawer() (+11 more)

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
Cohesion: 0.08
Nodes (18): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, ListS3Instances, ListS3InstancesRequest, ListS3InstancesResponse (+10 more)

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.08
Nodes (30): ActionStep, ConditionStep, GroupStep, LlmCallStep, ParallelStep, ScriptStep, TractCondition, TractDefinition (+22 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.24
Nodes (8): MAIL_DOMAIN_ICONS, mailProviderIcon(), AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps, DialogHead(), DialogHeadProps

### Community 42 - "SchemaProperty"
Cohesion: 0.50
Nodes (6): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleConnectionContent()

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.15
Nodes (13): Input(), FormField(), Props, CredentialRow(), CredentialRowProps, buildEmailRequest(), EmailAddDialog(), buildEmailRequest() (+5 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.14
Nodes (10): AdminSystemSettingsAPI, GetSettings, GetSettingsRequest, GetSettingsResponse, UpdateAuthMethods, UpdateAuthMethodsRequest, UpdateAuthMethodsResponse, UpdateRegistrationMode (+2 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.14
Nodes (16): useMcpKeys, ConnectionStep(), Props, AddTaskLinkDialog(), Props, RELATION_LABEL, RELATION_OPTIONS, RoadmapLinkTarget (+8 more)

### Community 48 - "dependencies"
Cohesion: 0.26
Nodes (14): execute(), fetchBoards(), fetchCards(), fetchLists(), TrelloBoardLite, TrelloCardLite, TrelloListLite, PickBoardStep() (+6 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.12
Nodes (15): GetVaultBySlug, GetVaultBySlugRequest, GetVaultBySlugResponse, PublicDocsGetNote, PublicDocsGetNoteRequest, PublicDocsGetNoteResponse, PublicDocsListFolders, PublicDocsListFoldersRequest (+7 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.17
Nodes (7): AdminUsersAPI, ArtelUserDetailDialog(), ArtelUserDetailDialogProps, UserSessionsDialog(), UserSessionsDialogProps, ChangeArtelPasswordDialog(), ChangeArtelPasswordDialogProps

### Community 51 - "AuthMiddleware"
Cohesion: 0.21
Nodes (10): apiPrefix(), csrfHeader(), getCsrfToken(), InitReq, TelegramLoginResponse, AppConfigState, useAppConfig, pingServer() (+2 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.17
Nodes (26): StepDraft, Props, InsertConflictDialog(), Props, appendStep(), branchArray(), buildStepFromDraft(), collapseThinParallels() (+18 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.10
Nodes (21): ArtelUserDetails, ArtelUserEntry, ChangeArtelUserPassword, ChangeArtelUserPasswordRequest, ChangeArtelUserPasswordResponse, CreateArtelUser, CreateArtelUserRequest, CreateArtelUserResponse (+13 more)

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
Cohesion: 0.13
Nodes (14): DeleteDockerHost, DeleteDockerHostRequest, DeleteDockerHostResponse, GetDockerHost, GetDockerHostRequest, ListDockerHosts, ListDockerHostsRequest, ListDockerHostsResponse (+6 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.08
Nodes (24): ExternalProvider, ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), ProviderChip(), PROVIDER_CHIP_CLASS (+16 more)

### Community 61 - "toTract"
Cohesion: 0.36
Nodes (4): VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.11
Nodes (19): TemplateSource, Props, ActionCard(), CardHeader(), Props, CONDITION_OPS, ConditionRowProps, InsertRow() (+11 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.17
Nodes (13): AuthAPI, UserState, PasswordLoginForm(), PasswordLoginFormProps, RegisterForm(), RegisterFormProps, LoginContentProps, Mode (+5 more)

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.05
Nodes (53): Props, RunButton(), Props, RunStatusBadge(), formatRelative(), Props, RunStatusDot(), LogicCell() (+45 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.13
Nodes (19): ResizeHandle(), ResizeHandleProps, cap(), LogEntryRow(), Props, buildLogLines(), formatDuration(), formatTime() (+11 more)

### Community 67 - "AuthAPI"
Cohesion: 0.13
Nodes (20): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), moveEntryAndRefresh(), NotesState, remapSelectedPath(), requireVaultId(), DropZone() (+12 more)

### Community 68 - "connectionLabel"
Cohesion: 0.07
Nodes (48): DialogManager, useDialog, useDialogKeyboard(), useTracts, useTrigger(), useTriggerSources(), Props, ToolStep() (+40 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.38
Nodes (10): AddAnthropicConnectionRequest, AddEmailConnectionRequest, AddGenericConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, CheckAnthropicConnectionRequest, CheckAnthropicConnectionResponse, ExternalConnectionsState (+2 more)

### Community 71 - "McpKeys.ts"
Cohesion: 0.29
Nodes (6): GetCouchInstanceResponse, InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps

### Community 72 - "ResultView.tsx"
Cohesion: 0.07
Nodes (37): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), STEP_SCREENS, AddTriggerDialogContext, AddTriggerDialogState (+29 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.09
Nodes (35): Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), boardListLabel(), Props, RoadmapCanvasNode() (+27 more)

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.18
Nodes (7): ExternalConnectionInfo, ConnectionSection(), Props, AccountsSection(), AccountsSectionProps, TrelloConnectionRow(), TrelloConnectionRowProps

### Community 75 - "dialog-scrollable.js"
Cohesion: 0.46
Nodes (7): allRules(), dialogScrollable(), directDeclsOf(), findScrollTarget(), isOverflowY(), messages, meta

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.18
Nodes (11): VaultInviteItem, VaultMemberInfo, vaultsQueryKey, Props, InviteRow(), Props, Props, Props (+3 more)

### Community 77 - "RunTractDialog.tsx"
Cohesion: 0.12
Nodes (15): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, formatStartedAt(), RunTractDialog() (+7 more)

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.22
Nodes (9): BinaryStorageToggle(), Props, Props, PublishSlugForm(), slugify(), validateSlug(), Props, PublishToggle() (+1 more)

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.22
Nodes (4): queryClient, originalFetch, refreshTokens(), SKIP_REFRESH_PATHS

### Community 80 - "CardMeta.tsx"
Cohesion: 0.12
Nodes (22): CreateTriggerRequest, TractTemplatesState, useTractTemplates, BrowseTemplatesDialog(), Props, TemplateRow(), Props, ListScreen() (+14 more)

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.39
Nodes (5): Tab, AdminHero(), AdminHeroProps, TabBar(), TabBarProps

### Community 93 - "VaultCardHeader.tsx"
Cohesion: 0.28
Nodes (3): Spreadsheet, GoogleSheetsSpreadsheetSection(), SpreadsheetRow()

### Community 94 - "AuthMiddleware"
Cohesion: 0.13
Nodes (4): AuthMiddleware, clearLocalStorage(), fromLocalStorage(), saveToLocalStorage()

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.09
Nodes (23): NoteMode, BreadcrumbBarProps, Mode, BreadcrumbPath(), BreadcrumbPathProps, DesktopNotesShellProps, VaultOption, CheckIcon() (+15 more)

### Community 96 - "UsersTab.tsx"
Cohesion: 0.31
Nodes (5): GenericToolIcon(), TODO: placeholder glyph for tool actions without a dedicated icon yet (non-smtp/, ImapIcon(), SmtpIcon(), ToolRow()

### Community 97 - "AdminCouchAPI"
Cohesion: 0.40
Nodes (4): CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.21
Nodes (8): CardChips(), Props, CardHeader(), Props, CardMeta(), formatDate(), Props, Props

### Community 100 - "RunLog.tsx"
Cohesion: 0.31
Nodes (4): ToolParamDef, ParamRow(), ParamsList(), coerceParams()

### Community 101 - "package.json"
Cohesion: 0.21
Nodes (6): DockerHostFormDialog(), DockerHostFormDialogSaveData, DockerHostsActionBar(), DockerHostsActionBarProps, DockerApiTab(), Props

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.24
Nodes (8): connectionLabel(), SelectOption(), LlmCallCard(), Props, ConnectionOptionList(), ConnectionOptionListProps, ConnectionPicker(), NodeChips()

### Community 103 - "CouchInstancesAPI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 104 - ".getRun"
Cohesion: 0.10
Nodes (25): useExternalConnections, LlmConnectionStep(), Props, ConnectGenericDialog(), CredentialField, AnthropicCheckButton(), CheckStatus, ConnectForm() (+17 more)

### Community 105 - "UserList.tsx"
Cohesion: 0.39
Nodes (5): GetDockerHostResponse, DockerHostFormDialogProps, DockerHostListProps, DockerHostRow(), DockerHostRowProps

### Community 107 - "AuthFetchInterceptor.ts"
Cohesion: 0.40
Nodes (3): ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps

### Community 108 - "EmailCheckButton.tsx"
Cohesion: 0.40
Nodes (4): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps

### Community 109 - "GitlabCheckButton.tsx"
Cohesion: 0.33
Nodes (4): ArtelAPI, Version, VersionRequest, VersionResponse

### Community 110 - "TrelloCheckButton.tsx"
Cohesion: 0.50
Nodes (3): CouchUserEntry, UserListProps, UserRowProps

### Community 111 - "TopbarDrawerCloseButton.tsx"
Cohesion: 0.18
Nodes (9): GetS3InstanceResponse, Props, S3ToggleFields(), Props, S3InstanceFormDialog(), S3InstanceRow(), TestStatus, S3InstancesActionBar() (+1 more)

### Community 112 - "InsertConflictDialog.tsx"
Cohesion: 0.15
Nodes (13): VaultItem, Props, DangerZoneText(), ExpertSettingsDrawer(), Props, Props, ExpertSettingsSection(), Props (+5 more)

### Community 117 - "requiredConnections.ts"
Cohesion: 0.14
Nodes (17): useVaults(), useWorkbench(), useWorkbenchMutations(), workbenchQueryKey(), VaultChipDisplayProps, VaultOptionList(), VaultOptionListProps, AuthMode (+9 more)

### Community 119 - "ConnectedContent.tsx"
Cohesion: 0.15
Nodes (8): ConnectedContent(), ConnectedContentProps, DialogHead(), DialogHeadProps, NotConnectedContentProps, PROVIDER_CONFIG, PROVIDER_KEY, ProviderConfig

### Community 120 - "KebabMenu.tsx"
Cohesion: 0.17
Nodes (11): RegistrationMode, SetupWizardAPI, RegistrationModeSection(), RegistrationModeSectionProps, TokenEntryScreen(), SetupWizardContext, SetupWizardState, SetupWizardStep (+3 more)

### Community 123 - "Tabs.tsx"
Cohesion: 0.16
Nodes (13): usePortrait(), Props, RenameDialog(), useAutosave(), UseAutosaveOptions, UseAutosaveResult, NotesPage(), buildNotesUrl() (+5 more)

### Community 125 - "Textarea.tsx"
Cohesion: 0.33
Nodes (5): Props, TODO: chures has no multiline variant yet, drop this wrapper once it does, Textarea(), DockerHostTlsFields(), Props

### Community 128 - "GitlabCheckButton.tsx"
Cohesion: 0.14
Nodes (16): CheckEmailConnectionRequest, CheckGitlabConnectionRequest, CheckTrelloConnectionRequest, CheckStatus, EmailCheckButton(), EmailCheckButtonProps, ConnectForm(), tokenSettingsUrl() (+8 more)

## Knowledge Gaps
- **796 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+791 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `Tracts.ts`, `TractIcons.tsx`, `ToolboxPage.tsx`, `useDialog`, `tractCanvasLayout.ts`, `BreadcrumbBar.tsx`, `NotesSidebar.tsx`, `index.ts`, `SchemaProperty`, `McpAuthPage.tsx`, `Topbar.tsx`, `TractBlockPicker.tsx`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `scripts`, `AuthAPI`, `connectionLabel`, `ResultView.tsx`, `MembersSection.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `GoogleSheetsSpreadsheetSection.tsx`, `ConnectForm.tsx`, `MobileDrawer.tsx`, `InsertConflictDialog.tsx`, `requiredConnections.ts`, `KebabMenu.tsx`, `Textarea.tsx`?**
  _High betweenness centrality (0.119) - this node is a cross-community bridge._
- **Why does `useDialog` connect `connectionLabel` to `GitlabCheckButton.tsx`, `TaskTrackersPage.tsx`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `TractIcons.tsx`, `compilerOptions`, `ToolboxPage.tsx`, `TractStepTree.tsx`, `ConnectionDetailDialog.tsx`, `McpAuthPage.tsx`, `useErrorToast.ts`, `dependencies`, `StepRow.tsx`, `AuthMiddleware`, `tractSteps.ts`, `admin_users.pb.ts`, `User.ts`, `TractBlockPicker.tsx`, `toTract`, `NotesPage.tsx`, `VaultCard.tsx`, `TractCanvasTopBar.tsx`, `AuthAPI`, `ResultView.tsx`, `RunTractDialog.tsx`, `DbAccessList.tsx`, `CardMeta.tsx`, `package.json`, `MobileDrawer.tsx`, `.getRun`, `UserList.tsx`, `TopbarDrawerCloseButton.tsx`, `ConnectedContent.tsx`, `Tabs.tsx`?**
  _High betweenness centrality (0.111) - this node is a cross-community bridge._
- **Why does `useUser` connect `Dialog.ts` to `GitlabCheckButton.tsx`, `TaskTrackersPage.tsx`, `TractIcons.tsx`, `compilerOptions`, `McpKeysAPI`, `Router.tsx`, `Notes.ts`, `index.ts`, `useVaultMutations`, `McpAuthPage.tsx`, `useErrorToast.ts`, `StepRow.tsx`, `Topbar.tsx`, `User.ts`, `toTract`, `VaultCard.tsx`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `McpKeys.ts`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `DbAccessList.tsx`, `VaultDangerZone.tsx`, `CardMeta.tsx`, `package.json`, `TopbarDrawerCloseButton.tsx`, `requiredConnections.ts`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _802 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.025 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08181818181818182 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.03125 - nodes in this community are weakly interconnected._