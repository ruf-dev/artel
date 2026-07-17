# Graph Report - ArtelUI  (2026-07-17)

## Corpus Check
- 461 files · ~127,162 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2087 nodes · 5024 edges · 104 communities (97 shown, 7 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 17 edges (avg confidence: 0.66)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b3ebdc09`
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
- RunStatusDot.tsx
- TopbarMobileTrigger.tsx

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 160 edges
2. `cn` - 126 edges
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
- `Props` --references--> `McpKeyInfo`  [EXTRACTED]
  src/widgets/McpKeyCard/McpKeyCard.tsx → src/app/api/artel/mcp_keys.pb.ts
- `ConnectionOptionListProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/components/ConnectionOptionList/ConnectionOptionList.tsx → src/app/api/artel/mcp_keys.pb.ts
- `SelectConnectionScreenProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/screens/SelectConnectionScreen.tsx → src/app/api/artel/mcp_keys.pb.ts
- `ConnectionFilterRowProps` --references--> `MomCandidate`  [EXTRACTED]
  src/pages/tract-canvas/components/TractBlockPicker/components/ConnectionFilterRow/ConnectionFilterRow.tsx → src/app/api/artel/mcp_keys.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`

## Communities (104 total, 7 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.03
Nodes (78): Absent, ActionStep, BaseTractStep, ConditionStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger (+70 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.06
Nodes (42): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+34 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.05
Nodes (42): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+34 more)

### Community 3 - "useBakeError"
Cohesion: 0.38
Nodes (6): AdminPage(), Tab, AdminHero(), AdminHeroProps, TabBar(), TabBarProps

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.04
Nodes (48): Absent, AddEmailConnection, AddEmailConnectionResponse, AddGitlabConnection, AddGitlabConnectionResponse, AddSpreadsheet, AddSpreadsheetRequest, AddSpreadsheetResponse (+40 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (42): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+34 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.14
Nodes (20): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenProperty() (+12 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.08
Nodes (31): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), ParamInput(), Props, ParamRow() (+23 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.08
Nodes (9): AddEmailConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, ExternalConnectionsAPI, Spreadsheet, ExternalConnectionsState, GoogleOAuthCallbackPage(), ExternalConnectionsService (+1 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.14
Nodes (27): Props, TractStepTreeProps, Props, CONDITION_OPS, Props, DangerZone(), Props, GroupBody() (+19 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.10
Nodes (20): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+12 more)

### Community 12 - "TractsService"
Cohesion: 0.06
Nodes (8): TractsAPI, parseSchema(), safeParseJson(), toRun(), toRunStep(), toTool(), toTrigger(), toTriggerSource()

### Community 13 - "Dialog.ts"
Cohesion: 0.21
Nodes (15): categoryLabel(), PROVIDER_ENUM_BY_KEY, providerLabel(), triggerChipLabel(), TriggerRow(), WebhookPicker(), WebhookPickerProps, PresetDetails() (+7 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.21
Nodes (17): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+9 more)

### Community 15 - "cn"
Cohesion: 0.06
Nodes (40): cn, DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it (+32 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.16
Nodes (13): Input(), AddTaskLinkDialog(), Props, RELATION_LABEL, RELATION_OPTIONS, RoadmapLinkTarget, WritableRelation, buildEmailRequest() (+5 more)

### Community 17 - "compilerOptions"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.21
Nodes (7): McpToolInfo, ToolParamDef, ParamRow(), ParamsList(), coerceParams(), ToolDetail(), ToolRow()

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.18
Nodes (11): MAIL_DOMAIN_ICONS, mailProviderIcon(), ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), ProviderChip() (+3 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.11
Nodes (20): UserErrors, BadRequestDetail, DetailType, DetailTypeName, ErrorInfoDetail, FieldViolation, getDetail(), getFieldViolations() (+12 more)

### Community 21 - "useDialog"
Cohesion: 0.14
Nodes (15): DialogManager, useDialogKeyboard(), useTracts, StepPickerDialog(), Props, TriggerPanel(), Props, RenameDialog() (+7 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.21
Nodes (8): CardChips(), Props, CardHeader(), Props, CardMeta(), formatDate(), Props, Props

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.16
Nodes (18): Props, cap(), NodeStatus, Props, title(), TractCanvasNode(), typeLabel(), addEdges() (+10 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.09
Nodes (21): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+13 more)

### Community 25 - "VaultItem"
Cohesion: 0.16
Nodes (11): VaultItem, Props, ExpertSettingsDrawer(), Props, Props, BinaryStorageToggle(), Props, Props (+3 more)

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
Cohesion: 0.17
Nodes (11): useServerStatus(), HomeLayout(), Path, Router(), routes, queryClient, ClosedAlphaPage(), ErrorPage() (+3 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.14
Nodes (14): DialogHead(), DialogHeadProps, VaultOptionList(), VaultOptionListProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), AddConnectionScreen() (+6 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.11
Nodes (15): Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), CopyIcon(), ErrorDotIcon(), PencilIcon(), SpinnerIcon() (+7 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.21
Nodes (8): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, useNotesSearchQuery()

### Community 33 - "useExternalConnections"
Cohesion: 0.13
Nodes (23): useDialog, useExternalConnections, VaultChip(), ManageEmailDialog(), WebhookSecretSection(), ConnectedContent(), ConnectedContentProps, ManageGitlabDialog() (+15 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.50
Nodes (6): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleConnectionContent()

### Community 35 - "notes.pb.ts"
Cohesion: 0.06
Nodes (30): CheckImportConflicts, CheckImportConflictsRequest, CheckImportConflictsResponse, CommitImport, CommitImportRequest, CommitImportResponse, DeleteFolder, DeleteFolderRequest (+22 more)

### Community 36 - "Notes.ts"
Cohesion: 0.08
Nodes (18): b64Decode(), ImportConflictAction, ImportResolution, NotesAPI, commitImportAndRefresh(), deleteFolderAndRefresh(), NotesState, requireVaultId() (+10 more)

### Community 37 - "index.ts"
Cohesion: 0.22
Nodes (8): ListPrompts, ListPromptsRequest, ListPromptsResponse, PromptId, PromptItem, PromptsAPI, FastSetupDialog(), Props

### Community 38 - "s3_instances.pb.ts"
Cohesion: 0.08
Nodes (18): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, ListS3Instances, ListS3InstancesRequest, ListS3InstancesResponse (+10 more)

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.33
Nodes (3): Props, Props, TractRunTimeline()

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.11
Nodes (17): ArtelUI Frontend Rules, Async style, Buttons, Component hierarchy, Component Structure, CSS Modules, Dialog shells must scroll internally, Error and Confirmation Handling (+9 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.11
Nodes (23): CheckTrelloConnectionRequest, useVaults(), FormField(), Props, HeroSegment(), HeroSegmentProps, CheckStatus, TrelloCheckButton() (+15 more)

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.11
Nodes (30): Props, Props, cardToNode(), createRoadmapGraph(), expandNode(), findNodeByShortLink(), MomToolExecutor, RoadmapEdge (+22 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.10
Nodes (20): useNotes, SuggestionList(), SuggestionListProps, CreateNoteDialog(), Props, ArtelLogoIcon(), ChevronLeftIcon(), ImportZipDialog() (+12 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.26
Nodes (9): ExternalConnectionInfo, AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps, AccountsSection(), AccountsSectionProps, TrelloConnectionRow() (+1 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.33
Nodes (3): VaultChipDisplayProps, VaultField(), VaultFieldProps

### Community 48 - "dependencies"
Cohesion: 0.14
Nodes (14): dependencies, classnames, framer-motion, grpc-web, marked, react, react-dom, react-router-dom (+6 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.12
Nodes (13): EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), TrelloIcon(), TODO: placeholder glyph for providers without a dedicated brand icon yet - repla, UnknownProviderIcon(), ProviderIcon() (+5 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.38
Nodes (5): JsonBlock(), Props, Props, statusClass(), StepRow()

### Community 51 - "AuthMiddleware"
Cohesion: 0.20
Nodes (12): apiPrefix(), InitReq, Options, TelegramLoginResponse, UserState, LoginContent(), LoginContentProps, AuthService (+4 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.26
Nodes (13): appendStep(), branchArray(), buildStepFromDraft(), collapseThinParallels(), collectAllStepIds(), generateStepId(), insertBlockAfter(), insertStepAt() (+5 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.09
Nodes (24): AdminUsersAPI, ArtelUserDetails, ArtelUserEntry, GetArtelUser, GetArtelUserRequest, GetArtelUserResponse, GetUserSessions, GetUserSessionsRequest (+16 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.09
Nodes (28): useIsMobileNav(), applyTheme(), Theme, useTheme(), BrandMarkIcon(), ConnectionsIcon(), base, NavIconProps (+20 more)

### Community 55 - "TractCanvasBuilder.tsx"
Cohesion: 0.13
Nodes (19): dompurify, NoteMode, usePortrait(), BreadcrumbBarProps, DesktopNotesShellProps, VaultOption, MobileNotesShell(), MobileNotesShellProps (+11 more)

### Community 56 - "User.ts"
Cohesion: 0.09
Nodes (24): AdminSubscriptionsAPI, EffectiveSubscriptionView, GetUserSubscription, GetUserSubscriptionRequest, GetUserSubscriptionResponse, ListSubscriptionPlans, ListSubscriptionPlansRequest, ListSubscriptionPlansResponse (+16 more)

### Community 57 - "Errors.ts"
Cohesion: 0.17
Nodes (6): ErrorReason, Errors, GrpcError, GrpcErrorDetail, ServiceError, ServiceErrorOption

### Community 58 - "NoteViewer.tsx"
Cohesion: 0.24
Nodes (8): ArtelMark(), ArtelMarkProps, ContentSegment, NoteContent(), NoteContentProps, parseWikiLinks(), WikiChip(), WikiChipProps

### Community 59 - "InviteLinksSection.tsx"
Cohesion: 0.21
Nodes (10): VaultInviteItem, VaultMemberInfo, vaultsQueryKey, Props, InviteRow(), Props, Props, Props (+2 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.17
Nodes (9): ExternalProvider, ConnectedContent(), ConnectedContentProps, DialogHead(), DialogHeadProps, NotConnectedContentProps, PROVIDER_CONFIG, PROVIDER_KEY (+1 more)

### Community 61 - "toTract"
Cohesion: 0.29
Nodes (5): McpLoginProps, VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.09
Nodes (27): ActionCard(), CardHeader(), Props, InsertRow(), Props, Props, SchemaTree(), buildSourcesFor() (+19 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.21
Nodes (8): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), ContentSegmentProps, Props, VaultCard()

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.07
Nodes (37): Props, RunButton(), LogicCell(), LogicCellProps, OptIcon(), OptIconProps, OptText(), OptTextProps (+29 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.33
Nodes (4): ArtelAPI, Version, VersionRequest, VersionResponse

### Community 68 - "connectionLabel"
Cohesion: 0.18
Nodes (13): useMcpKeys, connectionLabel(), SelectOption(), ConnectionStep(), Props, ConnectionOptionList(), ConnectionOptionListProps, ContentSegment() (+5 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.18
Nodes (9): GetS3InstanceResponse, Props, S3ToggleFields(), Props, S3InstanceFormDialog(), S3InstanceRow(), TestStatus, S3InstancesActionBar() (+1 more)

### Community 70 - "compilerOptions"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 71 - "McpKeys.ts"
Cohesion: 0.24
Nodes (10): McpConnectorInfo, MomCandidate, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectorRow(), ConnectorRowProps (+2 more)

### Community 72 - "ResultView.tsx"
Cohesion: 0.19
Nodes (12): TaskTrackerTableBody(), TaskTrackerTableHead(), DisplayTaskTrackerTables(), buildTaskTrackerTable(), isPlainObject(), stringifyCell(), TASK_TRACKER_ADAPTERS, TaskTrackerAdapter (+4 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.16
Nodes (13): useBakeError(), ConnectionDetailDialog(), CreateInviteLinkDialog(), InviteLinksSection(), InstanceRow(), ChangePasswordDialog(), ChangePasswordDialogProps, ManageAccessDialog() (+5 more)

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
Cohesion: 0.13
Nodes (19): TractsState, triggerSourcesQueryKey, triggersQueryKey, sleep(), formatStartedAt(), Props, RunTractDialog(), Props (+11 more)

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.40
Nodes (4): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.21
Nodes (7): DangerZoneText(), Props, ExpertSettingsSection(), MembersSection(), Props, Props, VaultDangerZone()

### Community 80 - "CardMeta.tsx"
Cohesion: 0.16
Nodes (15): useTrigger(), useTriggerSources(), providerEnumFor(), useAddTriggerDialog(), DialogHeaderWithClose(), CreateScreen(), KIND_OPTIONS, KindSettingsProps (+7 more)

### Community 93 - "VaultCardHeader.tsx"
Cohesion: 0.21
Nodes (12): Props, ToolStep(), Props, Step, StepDraft, Props, TractBlockPicker(), Args (+4 more)

### Community 94 - "AuthMiddleware"
Cohesion: 0.14
Nodes (4): AuthMiddleware, clearLocalStorage(), fromLocalStorage(), saveToLocalStorage()

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.18
Nodes (12): CheckEmailConnectionRequest, CheckGitlabConnectionRequest, CheckStatus, EmailCheckButton(), EmailCheckButtonProps, ConnectForm(), tokenSettingsUrl(), CheckStatus (+4 more)

### Community 96 - "UsersTab.tsx"
Cohesion: 0.22
Nodes (9): CouchUserEntry, GetCouchInstanceResponse, InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps, UserListProps (+1 more)

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LogPanelBar(), LogPanelBarProps, dotClass(), formatDate(), Props, TractCanvasLogPanel()

### Community 100 - "RunLog.tsx"
Cohesion: 0.43
Nodes (6): buildLogLines(), cap(), formatTime(), LogLine, RunLog(), RunLogProps

### Community 101 - "package.json"
Cohesion: 0.33
Nodes (5): name, private, trustedDependencies, type, version

### Community 102 - "RunStatusDot.tsx"
Cohesion: 0.60
Nodes (4): formatRelative(), Props, RunStatusDot(), TractLastRun

### Community 103 - "TopbarMobileTrigger.tsx"
Cohesion: 0.50
Nodes (3): TopbarHamburgerIcon(), TopbarMobileTrigger(), TopbarMobileTriggerProps

## Knowledge Gaps
- **610 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+605 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `useDialog` connect `useExternalConnections` to `TaskTrackersPage.tsx`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `TractIcons.tsx`, `ConnectorChip.tsx`, `useDialog`, `TractStepTree.tsx`, `Notes.ts`, `SchemaProperty`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `AuthMiddleware`, `admin_users.pb.ts`, `TractCanvasBuilder.tsx`, `User.ts`, `InviteLinksSection.tsx`, `TractBlockPicker.tsx`, `toTract`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `MembersSection.tsx`, `RunTractDialog.tsx`, `CardMeta.tsx`, `VaultCardHeader.tsx`, `ConnectForm.tsx`?**
  _High betweenness centrality (0.121) - this node is a cross-community bridge._
- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Tracts.ts`, `TractIcons.tsx`, `ConnectorChip.tsx`, `useDialog`, `tractCanvasLayout.ts`, `VaultItem`, `BreadcrumbBar.tsx`, `NotesSidebar.tsx`, `useExternalConnections`, `ConnectionDetailDialog.tsx`, `Notes.ts`, `index.ts`, `SchemaProperty`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `ProviderIcon.tsx`, `StepRow.tsx`, `Topbar.tsx`, `InviteLinksSection.tsx`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `connectionLabel`, `CardMeta.tsx`, `TractCanvasLogPanel.tsx`, `RunLog.tsx`, `RunStatusDot.tsx`, `TopbarMobileTrigger.tsx`?**
  _High betweenness centrality (0.116) - this node is a cross-community bridge._
- **Why does `useUser` connect `SchemaProperty` to `TaskTrackersPage.tsx`, `useBakeError`, `ExternalConnectionInfo`, `Dialog.ts`, `TractIcons.tsx`, `useDialog`, `VaultItem`, `McpKeysAPI`, `Router.tsx`, `useExternalConnections`, `Notes.ts`, `index.ts`, `useVaultMutations`, `McpAuthPage.tsx`, `AuthMiddleware`, `admin_users.pb.ts`, `Topbar.tsx`, `User.ts`, `InviteLinksSection.tsx`, `toTract`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `MembersSection.tsx`, `RunTractDialog.tsx`, `CardMeta.tsx`, `ConnectForm.tsx`, `UsersTab.tsx`, `AdminCouchAPI`?**
  _High betweenness centrality (0.071) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _613 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.030063291139240507 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06153846153846154 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.046511627906976744 - nodes in this community are weakly interconnected._