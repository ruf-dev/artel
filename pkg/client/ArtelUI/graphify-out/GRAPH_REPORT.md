# Graph Report - ArtelUI  (2026-07-09)

## Corpus Check
- 239 files · ~110,317 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1538 nodes · 3399 edges · 92 communities (85 shown, 7 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 13 edges (avg confidence: 0.57)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7edcce74`
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
- TopbarUserMenu.tsx
- compilerOptions
- McpKeys.ts
- S3InstancesAPI
- MembersSection.tsx
- LinkScreen.tsx
- package.json
- Handoff: lint/tooling parity gaps vs. ZpotifyUI
- Vaults.ts
- RunStatusDot.tsx
- AdminUsersAPI
- CardMeta.tsx
- EmptyState.tsx
- eslint.config.js

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 119 edges
2. `cn` - 87 edges
3. `useBakeError()` - 72 edges
4. `useUser` - 72 edges
5. `useTracts` - 34 edges
6. `useExternalConnections` - 31 edges
7. `TractStep` - 30 edges
8. `TractTool` - 30 edges
9. `SchemaNode` - 27 edges
10. `useMcpKeys` - 24 edges

## Surprising Connections (you probably didn't know these)
- `NoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/notes/components/NoteViewer/NoteViewer.tsx → package.json
- `Props` --references--> `VaultItem`  [EXTRACTED]
  src/widgets/VaultCard/VaultCard.tsx → src/app/api/artel/vaults.pb.ts
- `StepCard()` --calls--> `useDialog`  [EXTRACTED]
  src/components/TractStepTree/TractStepTree.tsx → src/app/hooks/Dialog.ts
- `TractStepTree()` --calls--> `useDialog`  [EXTRACTED]
  src/components/TractStepTree/TractStepTree.tsx → src/app/hooks/Dialog.ts
- `DialogHead()` --calls--> `useDialog`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/ManageKeyDialog.tsx → src/app/hooks/Dialog.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/mcp-keys/McpKeysPage.tsx -> src/dialogs/ManageKeyDialog/ManageKeyDialog.tsx -> src/app/routing/Router.tsx`

## Communities (92 total, 7 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.03
Nodes (58): Absent, BaseTractStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger, CreateTriggerResponse, DeleteTract (+50 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (24): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+16 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.05
Nodes (39): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+31 more)

### Community 3 - "useBakeError"
Cohesion: 0.09
Nodes (25): useBakeError(), LinkedView(), S3InstanceFormDialog(), useUser, AdminPage(), ArtelUserDetailDialog(), ArtelUserRow(), ArtelUsersTab() (+17 more)

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.05
Nodes (37): Absent, AddEmailConnection, AddEmailConnectionResponse, AddGitlabConnection, AddGitlabConnectionResponse, AddSpreadsheet, AddSpreadsheetRequest, AddSpreadsheetResponse (+29 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (36): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+28 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.12
Nodes (26): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenProperty() (+18 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.12
Nodes (24): AddTriggerDialog(), STEP_SCREENS, AddTriggerDialogContext, AddTriggerDialogState, AddTriggerStep, emptySchemaField(), FIELD_TYPES, fieldsToSchemaNode() (+16 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.13
Nodes (8): AddEmailConnectionRequest, AddGitlabConnectionRequest, ExternalConnectionInfo, ExternalConnectionsAPI, Spreadsheet, ExternalConnectionsState, ExternalConnectionsService, IExternalConnectionsService

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.06
Nodes (23): AdminCouchAPI, ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, CouchUserEntry, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse (+15 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.21
Nodes (24): MomCandidate, Props, TractStepTreeProps, ActionBody(), Props, Props, DangerZone(), Props (+16 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.07
Nodes (22): CouchInstancesAPI, DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceResponse, GetCouchInstanceStatus (+14 more)

### Community 12 - "TractsService"
Cohesion: 0.10
Nodes (4): TractsAPI, ITractsService, toTrigger(), TractsService

### Community 13 - "Dialog.ts"
Cohesion: 0.15
Nodes (20): useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor(), providerLabel(), triggerChipLabel(), Props, TriggerRow() (+12 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.09
Nodes (24): ActionStep, ConditionStep, CreateTriggerRequest, GroupStep, ParallelStep, TractCondition, TractDefinition, TractItem (+16 more)

### Community 15 - "cn"
Cohesion: 0.12
Nodes (18): cn, SelectOption(), CardHeader(), Props, ConnBarRow(), RowProps, InviteRow(), Props (+10 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.14
Nodes (19): StepColor, base, BranchIcon(), ChatIcon(), ChevronRightIcon(), EditIcon(), EnvelopeIcon(), ForkIcon() (+11 more)

### Community 17 - "compilerOptions"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.10
Nodes (14): ImapOperation, ImapToolAction, McpToolInfo, SmtpOperation, SmtpToolAction, ToolParamDef, useMcpKeys, HeroSegment() (+6 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.14
Nodes (15): ExternalProvider, MAIL_DOMAIN_ICONS, mailProviderIcon(), ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip() (+7 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.11
Nodes (20): UserErrors, BadRequestDetail, DetailType, DetailTypeName, ErrorInfoDetail, FieldViolation, getDetail(), getFieldViolations() (+12 more)

### Community 21 - "useDialog"
Cohesion: 0.14
Nodes (16): DialogManager, useDialog, Props, ToolStep(), TriggerPanel(), ConnectionOptionList(), HomePage(), InitPage() (+8 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.14
Nodes (15): useTrigger(), useVaults(), vaultsQueryKey, CardHeader(), Props, CandidateOptionList(), DialogHead(), ManageKeyDialog() (+7 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.15
Nodes (19): ConnectorPath(), Props, TractCanvasArea(), cap(), NodeStatus, title(), TractCanvasNode(), typeLabel() (+11 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.09
Nodes (21): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+13 more)

### Community 25 - "VaultItem"
Cohesion: 0.17
Nodes (14): GetS3InstanceResponse, VaultItem, Props, Props, ManageVaultS3Section(), Props, S3LinkPatch, CardChips() (+6 more)

### Community 26 - "devDependencies"
Cohesion: 0.10
Nodes (21): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+13 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.12
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 29 - "Router.tsx"
Cohesion: 0.15
Nodes (12): useServerStatus(), Path, Router(), routes, queryClient, ClosedAlphaPage(), GoogleOAuthCallbackPage(), ErrorPage() (+4 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.15
Nodes (16): InsertRow(), Props, buildSourcesFor(), collectIdsFromRoot(), CONDITION_OPS, ConditionCard(), ConditionCardProps, StepCard() (+8 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.12
Nodes (11): BreadcrumbBar(), BreadcrumbBarProps, Mode, Mode, ModeBar(), ModeBarProps, MODES, SaveStatus (+3 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.12
Nodes (14): PlusIcon(), SearchIcon(), SearchIconProps, buildFolderTree(), FolderNode, FolderNodeItem(), FolderNodeItemProps, FolderSection() (+6 more)

### Community 33 - "useExternalConnections"
Cohesion: 0.16
Nodes (13): useExternalConnections, ConnectionDetailDialog(), EmailAddDialog(), EmailEditDialog(), ManageEmailDialog(), ConnectForm(), ManageGitlabDialog(), WebhookSecretSection() (+5 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.16
Nodes (9): parseScopeList(), SCOPE_INFO, trimScope(), PROVIDER_CONFIG, PROVIDER_KEY, ProviderConfig, GoogleConnectionContent(), GapiWindow (+1 more)

### Community 35 - "notes.pb.ts"
Cohesion: 0.11
Nodes (18): GetNote, GetNoteRequest, GetNoteResponse, ListFolders, ListFoldersRequest, ListFoldersResponse, ListNotes, ListNotesRequest (+10 more)

### Community 36 - "Notes.ts"
Cohesion: 0.18
Nodes (5): NoteItem, NotesAPI, NotesState, INotesService, NotesService

### Community 37 - "index.ts"
Cohesion: 0.14
Nodes (12): ArtelAPI, Version, VersionRequest, VersionResponse, ListPrompts, ListPromptsRequest, ListPromptsResponse, PromptId (+4 more)

### Community 38 - "s3_instances.pb.ts"
Cohesion: 0.11
Nodes (17): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, ListS3Instances, ListS3InstancesRequest, ListS3InstancesResponse (+9 more)

### Community 39 - "useVaultMutations"
Cohesion: 0.16
Nodes (4): VaultsAPI, useVaultMutations(), Props, VaultDangerZone()

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.17
Nodes (11): useDialogKeyboard(), CreateNoteDialog(), Props, SuggestionList(), SuggestionListProps, NewTractButton(), Props, NewTractDialog() (+3 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.12
Nodes (15): ArtelUI Frontend Rules, Async style, Buttons, Component hierarchy, Component Structure, CSS Modules, Error and Confirmation Handling, Known debt (documented, not migrated) (+7 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.18
Nodes (11): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, Props, SchemaFieldRow() (+3 more)

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.17
Nodes (9): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+1 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.14
Nodes (9): apiPrefix(), InitReq, Options, TelegramLoginResponse, LoginContent(), McpLoginProps, Vault, VaultSelect() (+1 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.18
Nodes (11): NoteMode, ArtelLogoIcon(), CloseIcon(), getNoteMeta(), getNoteTitle(), MobileDrawerProps, MobileNotesShell(), MobileNotesShellProps (+3 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.24
Nodes (12): TractsState, triggerSourcesQueryKey, triggersQueryKey, formatStartedAt(), Props, RunTractDialog(), Props, CreatedTrigger (+4 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.19
Nodes (9): FormField(), Props, Props, S3ToggleFields(), Input, InputProps, TODO: predates the components/atoms/ convention (see pkg/client/ArtelUI/CLAUDE.m, ContentSegment() (+1 more)

### Community 48 - "dependencies"
Cohesion: 0.14
Nodes (14): dependencies, classnames, framer-motion, grpc-web, marked, react, react-dom, react-router-dom (+6 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.21
Nodes (7): EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), TrelloIcon(), TODO: placeholder glyph for providers without a dedicated brand icon yet - repla, UnknownProviderIcon()

### Community 50 - "StepRow.tsx"
Cohesion: 0.18
Nodes (9): JsonBlock(), Props, Props, statusClass(), StepRow(), Props, Props, TractRunTimeline() (+1 more)

### Community 51 - "AuthMiddleware"
Cohesion: 0.15
Nodes (4): AuthMiddleware, clearLocalStorage(), fromLocalStorage(), saveToLocalStorage()

### Community 52 - "tractSteps.ts"
Cohesion: 0.26
Nodes (13): appendStep(), branchArray(), buildStepFromDraft(), collapseThinParallels(), generateStepId(), hasChildren(), insertBlockAfter(), insertStepAt() (+5 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.15
Nodes (12): ArtelUserDetails, ArtelUserEntry, GetArtelUser, GetArtelUserRequest, GetArtelUserResponse, GetUserSessions, GetUserSessionsRequest, GetUserSessionsResponse (+4 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.23
Nodes (7): applyTheme(), Theme, useTheme(), HomeLayout(), TopbarBrand(), TopbarThemeToggle(), Topbar()

### Community 55 - "TractCanvasBuilder.tsx"
Cohesion: 0.29
Nodes (8): useTracts, sleep(), Props, Step, StepPickerDialog(), TractCanvasBuilder(), useTractRunTracking(), collectAllStepIds()

### Community 56 - "User.ts"
Cohesion: 0.42
Nodes (7): UserState, LoginContentProps, AuthService, IAuthService, Session, StoredAuth, UserInfo

### Community 57 - "Errors.ts"
Cohesion: 0.17
Nodes (6): ErrorReason, Errors, GrpcError, GrpcErrorDetail, ServiceError, ServiceErrorOption

### Community 58 - "NoteViewer.tsx"
Cohesion: 0.22
Nodes (8): dompurify, ContentSegment, NoteContent(), NoteViewer(), NoteViewerProps, parseWikiLinks(), WikiChip(), WikiChipProps

### Community 59 - "InviteLinksSection.tsx"
Cohesion: 0.24
Nodes (8): VaultInviteItem, CreateInviteLinkDialog(), Props, Props, Props, RoleOption(), InviteLinksSection(), Props

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.25
Nodes (10): StepDraft, ConnectionStep(), LOGIC_OPTIONS, LogicCell(), LogicOption, Props, ToolCell(), colorForKind() (+2 more)

### Community 61 - "toTract"
Cohesion: 0.22
Nodes (4): conditionToProto(), definitionToProto(), stepToProto(), toTract()

### Community 62 - "NotesPage.tsx"
Cohesion: 0.31
Nodes (6): useNotes, usePortrait(), Props, RenameDialog(), useAutosave(), NotesPage()

### Community 63 - "VaultCard.tsx"
Cohesion: 0.24
Nodes (6): Props, VaultCardConnBar(), VaultCardStatus(), Props, VaultCard(), VaultCardBack()

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.24
Nodes (6): Props, RunButton(), Props, RunStatusBadge(), Props, PlayIcon()

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.31
Nodes (9): buildLogLines(), cap(), dotClass(), formatDate(), formatTime(), LogLine, Props, RunLog() (+1 more)

### Community 66 - "scripts"
Cohesion: 0.22
Nodes (9): scripts, build, build:ui, dev, gen, lint, lint:css, lint:js (+1 more)

### Community 68 - "connectionLabel"
Cohesion: 0.39
Nodes (5): connectionLabel(), ConnectionStep(), Props, ActionCard(), NodeChips()

### Community 69 - "TopbarUserMenu.tsx"
Cohesion: 0.25
Nodes (5): MenuRect, TopbarUserMenu(), TopbarUserMenuItem(), TopbarUserMenuItemProps, TopbarUserMenuListProps

### Community 70 - "compilerOptions"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 71 - "McpKeys.ts"
Cohesion: 0.57
Nodes (6): CreateMcpKeyResponse, McpConnectorInfo, McpKeyInfo, McpKeysState, IMcpKeysService, Props

### Community 73 - "MembersSection.tsx"
Cohesion: 0.38
Nodes (4): VaultMemberInfo, Props, MembersSection(), Props

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.43
Nodes (4): useAddTriggerDialog(), DialogHeaderWithClose(), LinkScreen(), CloseIcon()

### Community 75 - "package.json"
Cohesion: 0.33
Nodes (5): name, private, trustedDependencies, type, version

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 78 - "RunStatusDot.tsx"
Cohesion: 0.60
Nodes (4): formatRelative(), Props, RunStatusDot(), TractLastRun

### Community 80 - "CardMeta.tsx"
Cohesion: 0.67
Nodes (3): CardMeta(), formatDate(), Props

## Knowledge Gaps
- **517 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+512 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `useUser` connect `useBakeError` to `TaskTrackersPage.tsx`, `ExternalConnectionInfo`, `Dialog.ts`, `Tracts.ts`, `ToolboxPage.tsx`, `useDialog`, `ManageKeyDialog.tsx`, `VaultItem`, `Router.tsx`, `useExternalConnections`, `Notes.ts`, `index.ts`, `useVaultMutations`, `CreateNoteDialog.tsx`, `McpAuthPage.tsx`, `Tracts.ts`, `useErrorToast.ts`, `Topbar.tsx`, `User.ts`, `AuthAPI`, `TopbarUserMenu.tsx`, `McpKeys.ts`, `S3InstancesAPI`, `Vaults.ts`?**
  _High betweenness centrality (0.109) - this node is a cross-community bridge._
- **Why does `useDialog` connect `useDialog` to `TaskTrackersPage.tsx`, `useBakeError`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `TractIcons.tsx`, `ToolboxPage.tsx`, `ConnectorChip.tsx`, `ManageKeyDialog.tsx`, `VaultItem`, `TractStepTree.tsx`, `NotesSidebar.tsx`, `useExternalConnections`, `ConnectionDetailDialog.tsx`, `CreateNoteDialog.tsx`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `Tracts.ts`, `useErrorToast.ts`, `TractCanvasBuilder.tsx`, `User.ts`, `InviteLinksSection.tsx`, `TractBlockPicker.tsx`, `NotesPage.tsx`, `VaultCard.tsx`, `connectionLabel`, `LinkScreen.tsx`?**
  _High betweenness centrality (0.101) - this node is a cross-community bridge._
- **Why does `cn` connect `cn` to `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `TractIcons.tsx`, `ConnectorChip.tsx`, `useDialog`, `tractCanvasLayout.ts`, `VaultItem`, `TractStepTree.tsx`, `BreadcrumbBar.tsx`, `ConnectionDetailDialog.tsx`, `index.ts`, `CreateNoteDialog.tsx`, `useErrorToast.ts`, `StepRow.tsx`, `InviteLinksSection.tsx`, `TractBlockPicker.tsx`, `VaultCard.tsx`, `TractCanvasTopBar.tsx`, `TractCanvasLogPanel.tsx`, `connectionLabel`, `LinkScreen.tsx`, `RunStatusDot.tsx`?**
  _High betweenness centrality (0.056) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _519 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.03389830508474576 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08414634146341464 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.05 - nodes in this community are weakly interconnected._