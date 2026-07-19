# Graph Report - ArtelUI  (2026-07-19)

## Corpus Check
- 493 files · ~135,504 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2254 nodes · 5462 edges · 107 communities (101 shown, 6 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 21 edges (avg confidence: 0.67)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `13a9ce84`
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

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 162 edges
2. `cn` - 142 edges
3. `useBakeError()` - 108 edges
4. `useUser` - 90 edges
5. `TractTool` - 42 edges
6. `MomCandidate` - 41 edges
7. `useExternalConnections` - 41 edges
8. `TractStep` - 39 edges
9. `useMcpKeys` - 34 edges
10. `ExternalConnectionInfo` - 33 edges

## Surprising Connections (you probably didn't know these)
- `NoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/notes/components/NoteViewer/NoteViewer.tsx → package.json
- `AccountsSectionProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/AccountsSection.tsx → src/app/api/artel/external_connections.pb.ts
- `EmailConnectionRowProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/components/EmailConnectionRow/EmailConnectionRow.tsx → src/app/api/artel/external_connections.pb.ts
- `AccountsSectionProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageTrelloDialog/components/AccountsSection/AccountsSection.tsx → src/app/api/artel/external_connections.pb.ts
- `TrelloConnectionRowProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageTrelloDialog/components/AccountsSection/components/TrelloConnectionRow/TrelloConnectionRow.tsx → src/app/api/artel/external_connections.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`

## Communities (107 total, 6 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.03
Nodes (65): Absent, BaseTractStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger, CreateTriggerRequest, CreateTriggerResponse (+57 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (27): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+19 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.05
Nodes (42): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+34 more)

### Community 3 - "useBakeError"
Cohesion: 0.39
Nodes (5): Tab, AdminHero(), AdminHeroProps, TabBar(), TabBarProps

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.04
Nodes (48): Absent, AddEmailConnection, AddEmailConnectionResponse, AddGitlabConnection, AddGitlabConnectionResponse, AddSpreadsheet, AddSpreadsheetRequest, AddSpreadsheetResponse (+40 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (36): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+28 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.14
Nodes (20): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenArrayItems() (+12 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.08
Nodes (41): triggersQueryKey, useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor(), providerLabel(), triggerChipLabel() (+33 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.09
Nodes (10): AddEmailConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, ExternalConnectionInfo, ExternalConnectionsAPI, Spreadsheet, ExternalConnectionsState, GoogleOAuthCallbackPage() (+2 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.20
Nodes (26): MomCandidate, Props, TractStepTreeProps, Props, Props, Props, GroupBody(), Props (+18 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.10
Nodes (20): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+12 more)

### Community 12 - "TractsService"
Cohesion: 0.08
Nodes (4): TractsAPI, TractsService, toTrigger(), toTriggerSource()

### Community 13 - "Dialog.ts"
Cohesion: 0.16
Nodes (16): SetTractsState, TractsState, triggerSourcesQueryKey, formatStartedAt(), Props, RunTractDialog(), Props, Deps (+8 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.08
Nodes (26): cn, KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it, ConnBarRow(), RowProps, Props (+18 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.15
Nodes (13): Input(), FormField(), Props, buildEmailRequest(), EmailAddDialog(), buildEmailRequest(), EmailEditDialog(), HostPortRowProps (+5 more)

### Community 17 - "compilerOptions"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.15
Nodes (11): ImapOperation, ImapToolAction, SmtpOperation, SmtpToolAction, ToolParamDef, ImapActionView(), ParamRow(), ParamsList() (+3 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.24
Nodes (8): MAIL_DOMAIN_ICONS, mailProviderIcon(), AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps, DialogHead(), DialogHeadProps

### Community 20 - "grpcErrors.ts"
Cohesion: 0.07
Nodes (36): CheckEmailConnectionRequest, CheckGitlabConnectionRequest, CheckTrelloConnectionRequest, UserErrors, CheckStatus, EmailCheckButton(), EmailCheckButtonProps, CheckStatus (+28 more)

### Community 21 - "useDialog"
Cohesion: 0.07
Nodes (34): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind, isNonEmptyBranch(), JsonBlock() (+26 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.09
Nodes (30): useServerStatus(), HomeLayout(), Path, Router(), routes, HeroSegment(), HeroSegmentProps, useUser (+22 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.11
Nodes (26): sleep(), ConnectorPath(), ConnectorPathProps, ParallelBoxes(), Props, Props, TractCanvasArea(), useTractCanvasDrag() (+18 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (24): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+16 more)

### Community 25 - "VaultItem"
Cohesion: 0.22
Nodes (10): VaultItem, ExpertSettingsDrawer(), Props, Props, ExpertSettingsSection(), Props, BinaryStorageToggle(), Props (+2 more)

### Community 26 - "devDependencies"
Cohesion: 0.10
Nodes (21): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+13 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.12
Nodes (7): CreateMcpKeyResponse, McpKeyInfo, McpKeysAPI, McpKeysState, IMcpKeysService, McpKeysService, Props

### Community 29 - "Router.tsx"
Cohesion: 0.10
Nodes (25): McpToolInfo, Props, SearchInput(), SelectOption(), execute(), fetchBoards(), fetchCards(), fetchLists() (+17 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.18
Nodes (12): DialogHead(), DialogHeadProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), AddConnectionScreen(), MainScreen(), MainScreenProps (+4 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.15
Nodes (11): BreadcrumbBarProps, Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), CopyIcon(), ErrorDotIcon(), PencilIcon() (+3 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.21
Nodes (8): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, useNotesSearchQuery()

### Community 33 - "useExternalConnections"
Cohesion: 0.14
Nodes (17): paramCompletionSource(), scriptEditorTheme, scriptHighlightStyle, Props, ScriptCodeSection(), addInputParam(), addOutputParam(), uniqueParamName() (+9 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.20
Nodes (11): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleSheetsSpreadsheetSection(), SpreadsheetRow(), GoogleConnectionContent() (+3 more)

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

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.09
Nodes (27): ActionStep, ConditionStep, GroupStep, ParallelStep, ScriptStep, TractCondition, TractDefinition, TractItem (+19 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.11
Nodes (17): ArtelUI Frontend Rules, Async style, Buttons, Component hierarchy, Component Structure, CSS Modules, Dialog shells must scroll internally, Error and Confirmation Handling (+9 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.12
Nodes (23): AdminUsersAPI, useExternalConnections, useBakeError(), ConnectionDetailDialog(), ManageEmailDialog(), WebhookSecretSection(), ConnectForm(), tokenSettingsUrl() (+15 more)

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.09
Nodes (35): Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), boardListLabel(), Props, RoadmapCanvasNode() (+27 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.25
Nodes (10): useNotes, NotesSidebar(), NotesSidebarProps, VaultOption, useFolderActions(), useHighlightNote(), getParentFolder(), isInvalidDrop() (+2 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.40
Nodes (4): AccountsSection(), AccountsSectionProps, TrelloConnectionRow(), TrelloConnectionRowProps

### Community 47 - "useErrorToast.ts"
Cohesion: 0.10
Nodes (23): VaultInviteItem, useMcpKeys, useVaults(), vaultsQueryKey, CardHeader(), Props, VaultChipDisplayProps, VaultField() (+15 more)

### Community 48 - "dependencies"
Cohesion: 0.09
Nodes (22): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, dompurify (+14 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.13
Nodes (12): EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), TrelloIcon(), TODO: placeholder glyph for providers without a dedicated brand icon yet - repla, UnknownProviderIcon(), ProviderIcon() (+4 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.19
Nodes (8): GetS3InstanceResponse, Props, S3ToggleFields(), Props, S3InstanceFormDialog(), S3InstanceRow(), TestStatus, S3InstancesActionBar()

### Community 51 - "AuthMiddleware"
Cohesion: 0.14
Nodes (13): apiPrefix(), InitReq, Options, TelegramLoginResponse, AuthAPI, pingServer(), UserState, LoginContentProps (+5 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.15
Nodes (27): StepDraft, DangerZone(), Props, InsertConflictDialog(), Props, branchArray(), buildStepFromDraft(), collapseThinParallels() (+19 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.13
Nodes (15): ArtelUserDetails, ArtelUserEntry, GetArtelUser, GetArtelUserRequest, GetArtelUserResponse, GetUserSessions, GetUserSessionsRequest, GetUserSessionsResponse (+7 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.08
Nodes (30): useIsMobileNav(), applyTheme(), Theme, useTheme(), BrandMarkIcon(), ConnectionsIcon(), base, NavIconProps (+22 more)

### Community 55 - "TractCanvasBuilder.tsx"
Cohesion: 0.13
Nodes (16): NoteMode, DesktopNotesShellProps, VaultOption, MobileDrawer(), MobileNotesShell(), MobileNotesShellProps, VaultOption, Mode (+8 more)

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
Cohesion: 0.25
Nodes (10): ScriptLanguage, ActionBody(), CONDITION_OPS, ConditionBody(), Props, ScriptBody(), Props, Section() (+2 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.17
Nodes (9): ExternalProvider, ConnectedContent(), ConnectedContentProps, DialogHead(), DialogHeadProps, NotConnectedContentProps, PROVIDER_CONFIG, PROVIDER_KEY (+1 more)

### Community 61 - "toTract"
Cohesion: 0.26
Nodes (6): McpLogin(), McpLoginProps, VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.10
Nodes (25): ActionCard(), CardHeader(), Props, InsertRow(), Props, Props, SchemaTree(), buildSourcesFor() (+17 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.16
Nodes (9): Props, VaultCardConnBar(), Props, VaultCardFront(), Props, VaultCardStatus(), ContentSegmentProps, Props (+1 more)

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.05
Nodes (50): Props, RunButton(), Props, RunStatusBadge(), formatRelative(), Props, RunStatusDot(), LogicCell() (+42 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.15
Nodes (5): ArtelAPI, Version, VersionRequest, VersionResponse, S3InstancesAPI

### Community 67 - "AuthAPI"
Cohesion: 0.22
Nodes (15): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), moveEntryAndRefresh(), NotesState, remapSelectedPath(), requireVaultId(), ConflictRow() (+7 more)

### Community 68 - "connectionLabel"
Cohesion: 0.10
Nodes (27): DialogManager, useDialog, useTracts, Props, ToolStep(), Props, Step, StepPickerDialog() (+19 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.44
Nodes (6): usePortrait(), useAutosave(), NotesPage(), buildNotesUrl(), decodeNotePath(), encodeNotePath()

### Community 70 - "compilerOptions"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 71 - "McpKeys.ts"
Cohesion: 0.22
Nodes (9): McpConnectorInfo, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps (+1 more)

### Community 72 - "ResultView.tsx"
Cohesion: 0.10
Nodes (23): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), isJsonValue(), TaskTrackerCell(), TaskTrackerTableHead() (+15 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.29
Nodes (7): connectionLabel(), ConnectorChip(), ConnectionStep(), Props, ConnectionOptionList(), ConnectionOptionListProps, ConnectionPicker()

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.22
Nodes (9): scripts, build, build:ui, dev, gen, lint, lint:css, lint:js (+1 more)

### Community 75 - "dialog-scrollable.js"
Cohesion: 0.46
Nodes (7): allRules(), dialogScrollable(), directDeclsOf(), findScrollTarget(), isOverflowY(), messages, meta

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 77 - "RunTractDialog.tsx"
Cohesion: 0.10
Nodes (20): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, Props, SchemaFieldRow() (+12 more)

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.40
Nodes (4): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.50
Nodes (3): DangerZoneText(), Props, VaultDangerZone()

### Community 80 - "CardMeta.tsx"
Cohesion: 0.31
Nodes (5): ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon(), TreeItemProps

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.22
Nodes (8): EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), CardChips(), Props, ProviderChip(), PROVIDER_CHIP_CLASS

### Community 93 - "VaultCardHeader.tsx"
Cohesion: 0.27
Nodes (6): VaultMemberInfo, Props, Props, RoleBadge(), MembersSection(), Props

### Community 94 - "AuthMiddleware"
Cohesion: 0.14
Nodes (4): AuthMiddleware, clearLocalStorage(), fromLocalStorage(), saveToLocalStorage()

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.21
Nodes (7): useDialogKeyboard(), SuggestionList(), SuggestionListProps, CreateNoteDialog(), Props, Props, RenameDialog()

### Community 96 - "UsersTab.tsx"
Cohesion: 0.28
Nodes (5): DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, ImportZipDialog(), Props

### Community 97 - "AdminCouchAPI"
Cohesion: 0.17
Nodes (6): AdminCouchAPI, ChangePasswordDialog(), ChangePasswordDialogProps, ManageAccessDialog(), ManageAccessDialogProps, UserRow()

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.29
Nodes (6): GetCouchInstanceResponse, InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.67
Nodes (3): CardMeta(), formatDate(), Props

### Community 100 - "RunLog.tsx"
Cohesion: 0.27
Nodes (8): AddTaskLinkDialog(), Props, RELATION_LABEL, RELATION_OPTIONS, RoadmapLinkTarget, WritableRelation, createRoadmapGraph(), RoadmapPage()

### Community 101 - "package.json"
Cohesion: 0.33
Nodes (5): name, private, trustedDependencies, type, version

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.27
Nodes (6): ArtelLogoIcon(), ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps, MobileDrawerProps, VaultOption

### Community 105 - "UserList.tsx"
Cohesion: 0.50
Nodes (3): CouchUserEntry, UserListProps, UserRowProps

## Knowledge Gaps
- **638 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+633 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Tracts.ts`, `TractIcons.tsx`, `ToolboxPage.tsx`, `useDialog`, `ManageKeyDialog.tsx`, `tractCanvasLayout.ts`, `VaultItem`, `Router.tsx`, `NotesSidebar.tsx`, `ConnectionDetailDialog.tsx`, `index.ts`, `McpAuthPage.tsx`, `useErrorToast.ts`, `ProviderIcon.tsx`, `Topbar.tsx`, `TractCanvasBuilder.tsx`, `InviteLinksSection.tsx`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `AuthAPI`, `connectionLabel`, `ResultView.tsx`, `RunTractDialog.tsx`, `CardMeta.tsx`, `GoogleSheetsSpreadsheetSection.tsx`, `VaultCardHeader.tsx`, `ConnectForm.tsx`, `UsersTab.tsx`?**
  _High betweenness centrality (0.135) - this node is a cross-community bridge._
- **Why does `useDialog` connect `connectionLabel` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `Dialog.ts`, `TractIcons.tsx`, `grpcErrors.ts`, `ManageKeyDialog.tsx`, `tractCanvasLayout.ts`, `Router.tsx`, `TractStepTree.tsx`, `ConnectionDetailDialog.tsx`, `SchemaProperty`, `MobileNotesShell.tsx`, `useErrorToast.ts`, `StepRow.tsx`, `AuthMiddleware`, `tractSteps.ts`, `admin_users.pb.ts`, `TractCanvasBuilder.tsx`, `User.ts`, `TractBlockPicker.tsx`, `toTract`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `AuthAPI`, `S3InstanceFormDialog.tsx`, `MembersSection.tsx`, `GoogleSheetsSpreadsheetSection.tsx`, `ConnectForm.tsx`, `UsersTab.tsx`, `AdminCouchAPI`, `RunLog.tsx`, `MobileDrawer.tsx`?**
  _High betweenness centrality (0.104) - this node is a cross-community bridge._
- **Why does `useUser` connect `ManageKeyDialog.tsx` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `ExternalConnectionInfo`, `Dialog.ts`, `TractIcons.tsx`, `grpcErrors.ts`, `VaultItem`, `McpKeysAPI`, `Notes.ts`, `index.ts`, `useVaultMutations`, `SchemaProperty`, `useErrorToast.ts`, `StepRow.tsx`, `AuthMiddleware`, `Topbar.tsx`, `User.ts`, `toTract`, `connectionLabel`, `RunTractDialog.tsx`, `AdminCouchAPI`, `CouchInstancesAPI`, `RunLog.tsx`, `Vaults.ts`?**
  _High betweenness centrality (0.067) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _642 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.030303030303030304 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08181818181818182 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.046511627906976744 - nodes in this community are weakly interconnected._