# Graph Report - ArtelUI  (2026-07-17)

## Corpus Check
- 463 files · ~127,501 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2093 nodes · 5042 edges · 102 communities (99 shown, 3 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 17 edges (avg confidence: 0.66)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `672e73db`
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

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 160 edges
2. `cn` - 130 edges
3. `useBakeError()` - 104 edges
4. `useUser` - 90 edges
5. `useExternalConnections` - 41 edges
6. `MomCandidate` - 38 edges
7. `TractTool` - 38 edges
8. `useMcpKeys` - 34 edges
9. `ExternalConnectionInfo` - 33 edges
10. `useTracts` - 32 edges

## Surprising Connections (you probably didn't know these)
- `NoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/notes/components/NoteViewer/NoteViewer.tsx → package.json
- `AccountsSectionProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/AccountsSection.tsx → src/app/api/artel/external_connections.pb.ts
- `EmailConnectionRowProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/components/EmailConnectionRow/EmailConnectionRow.tsx → src/app/api/artel/external_connections.pb.ts
- `Props` --references--> `McpKeyInfo`  [EXTRACTED]
  src/widgets/McpKeyCard/McpKeyCard.tsx → src/app/api/artel/mcp_keys.pb.ts
- `SelectConnectionScreenProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/screens/SelectConnectionScreen.tsx → src/app/api/artel/mcp_keys.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`

## Communities (102 total, 3 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.03
Nodes (59): Absent, BaseTractStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger, CreateTriggerRequest, CreateTriggerResponse (+51 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.07
Nodes (31): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+23 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.05
Nodes (42): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+34 more)

### Community 3 - "useBakeError"
Cohesion: 0.33
Nodes (6): Tab, AdminHero(), AdminHeroProps, TabBar(), TabBarProps, UsersTab()

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.04
Nodes (48): Absent, AddEmailConnection, AddEmailConnectionResponse, AddGitlabConnection, AddGitlabConnectionResponse, AddSpreadsheet, AddSpreadsheetRequest, AddSpreadsheetResponse (+40 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (42): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+34 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.15
Nodes (19): Props, buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenProperty(), flattenSource() (+11 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.10
Nodes (26): Props, SchemaFieldRow(), AddTriggerDialog(), STEP_SCREENS, AddTriggerDialogContext, AddTriggerDialogState, AddTriggerStep, emptySchemaField() (+18 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.09
Nodes (9): AddEmailConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, ExternalConnectionsAPI, Spreadsheet, ExternalConnectionsState, GoogleOAuthCallbackPage(), ExternalConnectionsService (+1 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.15
Nodes (29): Props, Props, Props, SchemaTree(), TractStepTreeProps, Props, CONDITION_OPS, ConditionBody() (+21 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.06
Nodes (30): CouchInstancesAPI, DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceResponse, GetCouchInstanceStatus (+22 more)

### Community 12 - "TractsService"
Cohesion: 0.07
Nodes (28): TractsAPI, TractsState, sleep(), formatStartedAt(), Props, RunTractDialog(), Props, TractCanvasBuilder() (+20 more)

### Community 13 - "Dialog.ts"
Cohesion: 0.16
Nodes (19): useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor(), providerLabel(), triggerChipLabel(), TriggerRow() (+11 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.21
Nodes (17): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+9 more)

### Community 15 - "cn"
Cohesion: 0.06
Nodes (40): cn, KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it, SelectOption(), SourceGroups(), CardHeader() (+32 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.28
Nodes (7): Input(), buildEmailRequest(), EmailAddDialog(), buildEmailRequest(), EmailEditDialog(), HostPortRowProps, useMailServerSuggestion()

### Community 17 - "compilerOptions"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.21
Nodes (7): McpToolInfo, ToolParamDef, ParamRow(), ParamsList(), coerceParams(), ToolDetail(), ToolRow()

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.24
Nodes (8): MAIL_DOMAIN_ICONS, mailProviderIcon(), AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps, DialogHead(), DialogHeadProps

### Community 20 - "grpcErrors.ts"
Cohesion: 0.08
Nodes (30): CheckEmailConnectionRequest, CheckTrelloConnectionRequest, UserErrors, CheckStatus, EmailCheckButton(), EmailCheckButtonProps, CheckStatus, TrelloCheckButton() (+22 more)

### Community 21 - "useDialog"
Cohesion: 0.14
Nodes (14): triggerSourcesQueryKey, triggersQueryKey, useTracts, Props, Props, TractRunTimeline(), Props, TriggerPanel() (+6 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.21
Nodes (11): useMcpKeys, HeroSegment(), HeroSegmentProps, CardHeader(), Props, ContentSegment(), McpKeysPage(), ContentSegment() (+3 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.15
Nodes (18): Props, cap(), NodeStatus, Props, title(), TractCanvasNode(), typeLabel(), addEdges() (+10 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.09
Nodes (21): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+13 more)

### Community 25 - "VaultItem"
Cohesion: 0.18
Nodes (11): VaultItem, CardChips(), Props, Props, VaultChip(), ExpertSettingsDrawer(), Props, Props (+3 more)

### Community 26 - "devDependencies"
Cohesion: 0.10
Nodes (21): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+13 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.13
Nodes (6): CreateMcpKeyResponse, McpKeyInfo, McpKeysAPI, McpKeysState, IMcpKeysService, McpKeysService

### Community 29 - "Router.tsx"
Cohesion: 0.21
Nodes (16): Props, SearchInput(), execute(), fetchBoards(), fetchCards(), fetchLists(), TrelloBoardLite, TrelloCardLite (+8 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.19
Nodes (11): DialogHead(), DialogHeadProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), MainScreen(), MainScreenProps, SelectConnectionScreen() (+3 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.16
Nodes (10): Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), CopyIcon(), ErrorDotIcon(), PencilIcon(), SpinnerIcon() (+2 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.21
Nodes (8): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, useNotesSearchQuery()

### Community 33 - "useExternalConnections"
Cohesion: 0.17
Nodes (12): useExternalConnections, ManageEmailDialog(), ConnectedContent(), ConnectedContentProps, DialogHead(), ManageGitlabDialog(), ManageTrelloDialog(), ContentSegment() (+4 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.20
Nodes (11): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleSheetsSpreadsheetSection(), SpreadsheetRow(), GoogleConnectionContent() (+3 more)

### Community 35 - "notes.pb.ts"
Cohesion: 0.06
Nodes (30): CheckImportConflicts, CheckImportConflictsRequest, CheckImportConflictsResponse, CommitImport, CommitImportRequest, CommitImportResponse, DeleteFolder, DeleteFolderRequest (+22 more)

### Community 36 - "Notes.ts"
Cohesion: 0.11
Nodes (5): b64Decode(), ImportResolution, NotesAPI, INotesService, NotesService

### Community 37 - "index.ts"
Cohesion: 0.22
Nodes (8): ListPrompts, ListPromptsRequest, ListPromptsResponse, PromptId, PromptItem, PromptsAPI, FastSetupDialog(), Props

### Community 38 - "s3_instances.pb.ts"
Cohesion: 0.07
Nodes (27): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, GetS3InstanceResponse, ListS3Instances, ListS3InstancesRequest (+19 more)

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.12
Nodes (19): ActionStep, ConditionStep, GroupStep, ParallelStep, TractCondition, TractDefinition, TractItem, TractRunItem (+11 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.11
Nodes (17): ArtelUI Frontend Rules, Async style, Buttons, Component hierarchy, Component Structure, CSS Modules, Dialog shells must scroll internally, Error and Confirmation Handling (+9 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.10
Nodes (24): AdminUsersAPI, GetCouchInstanceStatusResponse, useBakeError(), FormField(), Props, AddTaskLinkDialog(), Props, RELATION_LABEL (+16 more)

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.09
Nodes (35): Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), boardListLabel(), Props, RoadmapCanvasNode() (+27 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.13
Nodes (16): useNotes, CreateNoteDialog(), Props, ArtelLogoIcon(), ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps, MobileDrawer() (+8 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.25
Nodes (9): ExternalConnectionInfo, AccountsSection(), AccountsSectionProps, TrelloConnectionRow(), TrelloConnectionRowProps, DialogHead(), DialogHeadProps, tokenAuthorizeUrl() (+1 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.20
Nodes (10): useVaults(), vaultsQueryKey, VaultChipDisplayProps, VaultField(), VaultFieldProps, VaultOptionList(), VaultOptionListProps, CreateKeyDialog() (+2 more)

### Community 48 - "dependencies"
Cohesion: 0.14
Nodes (14): dependencies, classnames, framer-motion, grpc-web, marked, react, react-dom, react-router-dom (+6 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.11
Nodes (19): ExternalProvider, ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), ProviderChip(), PROVIDER_CHIP_CLASS (+11 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.38
Nodes (5): JsonBlock(), Props, Props, statusClass(), StepRow()

### Community 51 - "AuthMiddleware"
Cohesion: 0.08
Nodes (31): apiPrefix(), InitReq, Options, TelegramLoginResponse, AuthAPI, pingServer(), useServerStatus(), HomeLayout() (+23 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.14
Nodes (21): InsertRow(), Props, collectIdsFromRoot(), ConditionCard(), ConditionCardProps, StepCard(), TractStepTree(), appendStep() (+13 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.13
Nodes (15): ArtelUserDetails, ArtelUserEntry, GetArtelUser, GetArtelUserRequest, GetArtelUserResponse, GetUserSessions, GetUserSessionsRequest, GetUserSessionsResponse (+7 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.09
Nodes (28): useIsMobileNav(), applyTheme(), Theme, useTheme(), BrandMarkIcon(), ConnectionsIcon(), base, NavIconProps (+20 more)

### Community 55 - "TractCanvasBuilder.tsx"
Cohesion: 0.17
Nodes (13): dompurify, NoteMode, BreadcrumbBarProps, DesktopNotesShellProps, VaultOption, MobileNotesShell(), MobileNotesShellProps, VaultOption (+5 more)

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
Cohesion: 0.23
Nodes (9): VaultInviteItem, CreateInviteLinkDialog(), Props, InviteRow(), Props, Props, RoleBadge(), InviteLinksSection() (+1 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.17
Nodes (7): ConnectedContent(), ConnectedContentProps, NotConnectedContentProps, ConnectionDetailDialog(), PROVIDER_CONFIG, PROVIDER_KEY, ProviderConfig

### Community 61 - "toTract"
Cohesion: 0.43
Nodes (4): VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.26
Nodes (12): buildSourcesFor(), buildSources(), computeVisibleStepIds(), extractRefs(), findStepById(), findVisible(), isEscapedOpen(), isOpen() (+4 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.19
Nodes (9): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), ContentSegment(), ContentSegmentProps, Props (+1 more)

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.05
Nodes (49): Props, RunButton(), Props, RunStatusBadge(), formatRelative(), Props, RunStatusDot(), LogicCell() (+41 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.33
Nodes (4): ArtelAPI, Version, VersionRequest, VersionResponse

### Community 67 - "AuthAPI"
Cohesion: 0.24
Nodes (13): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), NotesState, requireVaultId(), ConflictRow(), Props, ImportConflictsDialog() (+5 more)

### Community 68 - "connectionLabel"
Cohesion: 0.09
Nodes (24): DialogManager, useDialog, useDialogKeyboard(), ConnectionStep(), Props, ToolStep(), Props, Step (+16 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.44
Nodes (6): usePortrait(), useAutosave(), NotesPage(), buildNotesUrl(), decodeNotePath(), encodeNotePath()

### Community 70 - "compilerOptions"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 71 - "McpKeys.ts"
Cohesion: 0.23
Nodes (10): MomCandidate, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectionOptionListProps, AddConnectionScreen(), AddConnectionScreenProps (+2 more)

### Community 72 - "ResultView.tsx"
Cohesion: 0.13
Nodes (18): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), isJsonValue(), TaskTrackerCell(), TaskTrackerTableHead() (+10 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.39
Nodes (5): connectionLabel(), ActionCard(), ConnectionPicker(), ActionBody(), NodeChips()

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
Cohesion: 0.28
Nodes (6): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props

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
Cohesion: 0.48
Nodes (4): McpConnectorInfo, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps

### Community 93 - "VaultCardHeader.tsx"
Cohesion: 0.38
Nodes (4): VaultMemberInfo, Props, MembersSection(), Props

### Community 94 - "AuthMiddleware"
Cohesion: 0.14
Nodes (4): AuthMiddleware, clearLocalStorage(), fromLocalStorage(), saveToLocalStorage()

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.36
Nodes (6): CheckGitlabConnectionRequest, ConnectForm(), tokenSettingsUrl(), CheckStatus, GitlabCheckButton(), GitlabCheckButtonProps

### Community 96 - "UsersTab.tsx"
Cohesion: 0.33
Nodes (5): DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, ImportZipDialog(), Props

### Community 97 - "AdminCouchAPI"
Cohesion: 0.16
Nodes (8): AdminCouchAPI, CouchUserEntry, ChangePasswordDialog(), ChangePasswordDialogProps, ManageAccessDialog(), UserListProps, UserRow(), UserRowProps

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.40
Nodes (4): Mode, ModeBar(), ModeBarProps, MODES

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.67
Nodes (3): CardMeta(), formatDate(), Props

### Community 100 - "RunLog.tsx"
Cohesion: 0.43
Nodes (6): buildLogLines(), cap(), formatTime(), LogLine, RunLog(), RunLogProps

### Community 101 - "package.json"
Cohesion: 0.33
Nodes (5): name, private, trustedDependencies, type, version

## Knowledge Gaps
- **610 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+605 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `useDialog` connect `connectionLabel` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `couch_instances.pb.ts`, `TractsService`, `Dialog.ts`, `TractIcons.tsx`, `ConnectorChip.tsx`, `grpcErrors.ts`, `useDialog`, `ManageKeyDialog.tsx`, `VaultItem`, `Router.tsx`, `TractStepTree.tsx`, `useExternalConnections`, `ConnectionDetailDialog.tsx`, `s3_instances.pb.ts`, `SchemaProperty`, `MobileNotesShell.tsx`, `Tracts.ts`, `useErrorToast.ts`, `ProviderIcon.tsx`, `AuthMiddleware`, `tractSteps.ts`, `admin_users.pb.ts`, `User.ts`, `InviteLinksSection.tsx`, `TractBlockPicker.tsx`, `toTract`, `TractCanvasTopBar.tsx`, `AuthAPI`, `S3InstanceFormDialog.tsx`, `ConnectForm.tsx`, `UsersTab.tsx`, `AdminCouchAPI`?**
  _High betweenness centrality (0.115) - this node is a cross-community bridge._
- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `Tracts.ts`, `TractIcons.tsx`, `useDialog`, `ManageKeyDialog.tsx`, `tractCanvasLayout.ts`, `VaultItem`, `Router.tsx`, `NotesSidebar.tsx`, `useExternalConnections`, `ConnectionDetailDialog.tsx`, `index.ts`, `McpAuthPage.tsx`, `ProviderIcon.tsx`, `StepRow.tsx`, `tractSteps.ts`, `Topbar.tsx`, `InviteLinksSection.tsx`, `TractCanvasTopBar.tsx`, `AuthAPI`, `connectionLabel`, `McpKeys.ts`, `ResultView.tsx`, `MembersSection.tsx`, `CardMeta.tsx`, `UsersTab.tsx`, `CouchInstancesAPI`, `RunLog.tsx`?**
  _High betweenness centrality (0.107) - this node is a cross-community bridge._
- **Why does `useUser` connect `AuthMiddleware` to `TaskTrackersPage.tsx`, `useBakeError`, `ExternalConnectionInfo`, `couch_instances.pb.ts`, `TractsService`, `Dialog.ts`, `TractIcons.tsx`, `grpcErrors.ts`, `useDialog`, `ManageKeyDialog.tsx`, `VaultItem`, `McpKeysAPI`, `useExternalConnections`, `Notes.ts`, `index.ts`, `s3_instances.pb.ts`, `useVaultMutations`, `SchemaProperty`, `useErrorToast.ts`, `Topbar.tsx`, `User.ts`, `ConnectForm.tsx`, `AdminCouchAPI`?**
  _High betweenness centrality (0.076) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _613 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.03333333333333333 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07450980392156863 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.046511627906976744 - nodes in this community are weakly interconnected._