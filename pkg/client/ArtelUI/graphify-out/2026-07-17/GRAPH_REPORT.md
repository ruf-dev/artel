# Graph Report - ArtelUI  (2026-07-17)

## Corpus Check
- 445 files · ~123,152 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2013 nodes · 4807 edges · 94 communities (90 shown, 4 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 17 edges (avg confidence: 0.66)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `86c2f6f6`
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

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 152 edges
2. `cn` - 116 edges
3. `useBakeError()` - 94 edges
4. `useUser` - 88 edges
5. `useExternalConnections` - 41 edges
6. `MomCandidate` - 38 edges
7. `TractTool` - 38 edges
8. `ExternalConnectionInfo` - 33 edges
9. `useTracts` - 32 edges
10. `TractStep` - 31 edges

## Surprising Connections (you probably didn't know these)
- `NoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/notes/components/NoteViewer/NoteViewer.tsx → package.json
- `ConnectionOptionListProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/components/ConnectionOptionList/ConnectionOptionList.tsx → src/app/api/artel/mcp_keys.pb.ts
- `SelectConnectionScreenProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/screens/SelectConnectionScreen.tsx → src/app/api/artel/mcp_keys.pb.ts
- `Props` --references--> `VaultItem`  [EXTRACTED]
  src/components/VaultCard/VaultCardFront.tsx → src/app/api/artel/vaults.pb.ts
- `Props` --references--> `VaultItem`  [EXTRACTED]
  src/components/VaultCard/VaultCardHeader.tsx → src/app/api/artel/vaults.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`

## Communities (94 total, 4 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.03
Nodes (72): Absent, ActionStep, BaseTractStep, ConditionStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger (+64 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (29): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+21 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.05
Nodes (42): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+34 more)

### Community 3 - "useBakeError"
Cohesion: 0.22
Nodes (9): AdminPage(), Tab, AdminHero(), AdminHeroProps, InstancesActionBar(), InstancesActionBarProps, InstancesTab(), TabBar() (+1 more)

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.04
Nodes (48): Absent, AddEmailConnection, AddEmailConnectionResponse, AddGitlabConnection, AddGitlabConnectionResponse, AddSpreadsheet, AddSpreadsheetRequest, AddSpreadsheetResponse (+40 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (41): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+33 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.09
Nodes (34): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenProperty() (+26 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.09
Nodes (30): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), ParamInput(), Props, ParamRow() (+22 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.08
Nodes (9): AddEmailConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, ExternalConnectionsAPI, Spreadsheet, ExternalConnectionsState, GoogleOAuthCallbackPage(), ExternalConnectionsService (+1 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.08
Nodes (24): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, CouchUserEntry, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess (+16 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.13
Nodes (33): Props, ToolStep(), Props, Props, SchemaTree(), TractStepTreeProps, ActionBody(), Props (+25 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.06
Nodes (29): CouchInstancesAPI, DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceResponse, GetCouchInstanceStatus (+21 more)

### Community 12 - "TractsService"
Cohesion: 0.06
Nodes (33): TractsAPI, TractsState, triggerSourcesQueryKey, triggersQueryKey, Props, PresetDetails(), Props, Deps (+25 more)

### Community 13 - "Dialog.ts"
Cohesion: 0.16
Nodes (19): useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor(), providerLabel(), triggerChipLabel(), TriggerRow() (+11 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.21
Nodes (17): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+9 more)

### Community 15 - "cn"
Cohesion: 0.05
Nodes (48): cn, DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it (+40 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.15
Nodes (13): Input(), FormField(), Props, buildEmailRequest(), EmailAddDialog(), buildEmailRequest(), EmailEditDialog(), HostPortRowProps (+5 more)

### Community 17 - "compilerOptions"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.19
Nodes (8): McpToolInfo, ToolParamDef, ImapActionView(), ParamRow(), ParamsList(), coerceParams(), ToolDetail(), ToolRow()

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.18
Nodes (11): MAIL_DOMAIN_ICONS, mailProviderIcon(), ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), ProviderChip() (+3 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.07
Nodes (34): CheckEmailConnectionRequest, CheckGitlabConnectionRequest, CheckTrelloConnectionRequest, UserErrors, CheckStatus, EmailCheckButton(), EmailCheckButtonProps, CheckStatus (+26 more)

### Community 21 - "useDialog"
Cohesion: 0.11
Nodes (24): useDialog, useDialogKeyboard(), useTracts, StepPickerDialog(), Props, TriggerPanel(), AddTriggerDialog(), DialogHeaderWithClose() (+16 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.21
Nodes (9): McpKeyInfo, CardChips(), CardHeader(), Props, CardMeta(), formatDate(), Props, MainScreenProps (+1 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.13
Nodes (21): sleep(), ConnectorPath(), ConnectorPathProps, Props, TractCanvasArea(), TractCanvasBuilder(), useTractCanvasBuilderHandlers(), useTractRunTracking() (+13 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.09
Nodes (21): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+13 more)

### Community 25 - "VaultItem"
Cohesion: 0.24
Nodes (9): VaultItem, Props, VaultChip(), ExpertSettingsDrawer(), Props, Props, ExpertSettingsSection(), Props (+1 more)

### Community 26 - "devDependencies"
Cohesion: 0.10
Nodes (21): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+13 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.11
Nodes (9): CreateMcpKeyResponse, McpConnectorInfo, McpKeysAPI, McpKeysState, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps, IMcpKeysService (+1 more)

### Community 29 - "Router.tsx"
Cohesion: 0.09
Nodes (29): DialogManager, useMcpKeys, useServerStatus(), HomeLayout(), Path, Router(), routes, HeroSegment() (+21 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.21
Nodes (10): DialogHead(), DialogHeadProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), MainScreen(), SelectConnectionScreen(), SelectConnectionScreenProps (+2 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.16
Nodes (10): Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), CopyIcon(), ErrorDotIcon(), PencilIcon(), SpinnerIcon() (+2 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.21
Nodes (8): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, useNotesSearchQuery()

### Community 33 - "useExternalConnections"
Cohesion: 0.18
Nodes (18): useExternalConnections, useBakeError(), ManageEmailDialog(), ConnectForm(), tokenSettingsUrl(), DialogHead(), ManageGitlabDialog(), DialogHead() (+10 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.50
Nodes (6): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleConnectionContent()

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
Cohesion: 0.08
Nodes (18): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, ListS3Instances, ListS3InstancesRequest, ListS3InstancesResponse (+10 more)

### Community 39 - "useVaultMutations"
Cohesion: 0.15
Nodes (4): VaultsAPI, useVaultMutations(), IVaultService, VaultService

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.33
Nodes (3): Props, Props, TractRunTimeline()

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.11
Nodes (17): ArtelUI Frontend Rules, Async style, Buttons, Component hierarchy, Component Structure, CSS Modules, Dialog shells must scroll internally, Error and Confirmation Handling (+9 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.24
Nodes (13): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), NotesState, requireVaultId(), ConflictRow(), Props, ImportConflictsDialog() (+5 more)

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.31
Nodes (5): ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon(), TreeItemProps

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.14
Nodes (15): useNotes, CreateNoteDialog(), Props, ArtelLogoIcon(), ImportZipDialog(), Props, MobileDrawer(), MobileDrawerProps (+7 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.21
Nodes (9): ExternalConnectionInfo, AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps, AccountsSection(), AccountsSectionProps, TrelloConnectionRow() (+1 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.20
Nodes (11): useVaults(), vaultsQueryKey, VaultChipDisplayProps, VaultField(), VaultFieldProps, VaultOptionList(), VaultOptionListProps, ContentSegment() (+3 more)

### Community 48 - "dependencies"
Cohesion: 0.14
Nodes (14): dependencies, classnames, framer-motion, grpc-web, marked, react, react-dom, react-router-dom (+6 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.14
Nodes (11): EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), TrelloIcon(), TODO: placeholder glyph for providers without a dedicated brand icon yet - repla, UnknownProviderIcon(), ProviderIcon() (+3 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.38
Nodes (5): JsonBlock(), Props, Props, statusClass(), StepRow()

### Community 51 - "AuthMiddleware"
Cohesion: 0.11
Nodes (16): apiPrefix(), InitReq, Options, TelegramLoginResponse, UserState, LoginContent(), LoginContentProps, AuthService (+8 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.14
Nodes (23): Props, Step, StepDraft, InsertRow(), Props, ConditionCardProps, StepCard(), TractStepTree() (+15 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.10
Nodes (21): AdminUsersAPI, ArtelUserDetails, ArtelUserEntry, GetArtelUser, GetArtelUserRequest, GetArtelUserResponse, GetUserSessions, GetUserSessionsRequest (+13 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.07
Nodes (33): useIsMobileNav(), applyTheme(), Theme, useTheme(), BrandMarkIcon(), ConnectionsIcon(), base, NavIconProps (+25 more)

### Community 55 - "TractCanvasBuilder.tsx"
Cohesion: 0.12
Nodes (17): dompurify, NoteMode, BreadcrumbBarProps, DesktopNotesShellProps, VaultOption, MobileNotesShell(), MobileNotesShellProps, VaultOption (+9 more)

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
Cohesion: 0.18
Nodes (10): VaultInviteItem, VaultMemberInfo, Props, InviteRow(), Props, Props, Props, RoleBadge() (+2 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.19
Nodes (10): ExternalProvider, ConnectedContent(), ConnectedContentProps, DialogHead(), DialogHeadProps, NotConnectedContentProps, ConnectionDetailDialog(), PROVIDER_CONFIG (+2 more)

### Community 61 - "toTract"
Cohesion: 0.43
Nodes (4): VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.44
Nodes (6): usePortrait(), useAutosave(), NotesPage(), buildNotesUrl(), decodeNotePath(), encodeNotePath()

### Community 63 - "VaultCard.tsx"
Cohesion: 0.24
Nodes (7): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), Props, VaultCard()

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.06
Nodes (43): Props, RunButton(), Props, RunStatusBadge(), LogicCell(), LogicCellProps, OptIcon(), OptIconProps (+35 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.33
Nodes (4): ArtelAPI, Version, VersionRequest, VersionResponse

### Community 67 - "AuthAPI"
Cohesion: 0.15
Nodes (5): AuthAPI, pingServer(), BinaryStorageToggle(), Props, VaultSettingsSection()

### Community 68 - "connectionLabel"
Cohesion: 0.24
Nodes (9): connectionLabel(), SelectOption(), ConnectionStep(), Props, ConnectionOptionList(), ConnectionOptionListProps, ConnectionPicker(), ConnectionStep() (+1 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.18
Nodes (9): GetS3InstanceResponse, Props, S3ToggleFields(), Props, S3InstanceFormDialog(), S3InstanceRow(), TestStatus, S3InstancesActionBar() (+1 more)

### Community 70 - "compilerOptions"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 71 - "McpKeys.ts"
Cohesion: 0.26
Nodes (9): MomCandidate, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, AddConnectionScreen(), AddConnectionScreenProps, ConnectionFilterRow() (+1 more)

### Community 72 - "ResultView.tsx"
Cohesion: 0.19
Nodes (12): TaskTrackerTableBody(), TaskTrackerTableHead(), DisplayTaskTrackerTables(), buildTaskTrackerTable(), isPlainObject(), stringifyCell(), TASK_TRACKER_ADAPTERS, TaskTrackerAdapter (+4 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.17
Nodes (6): AdminCouchAPI, ChangePasswordDialog(), ChangePasswordDialogProps, ManageAccessDialog(), ManageAccessDialogProps, UserRow()

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.13
Nodes (14): name, private, scripts, build, build:ui, dev, gen, lint (+6 more)

### Community 75 - "dialog-scrollable.js"
Cohesion: 0.46
Nodes (7): allRules(), dialogScrollable(), directDeclsOf(), findScrollTarget(), isOverflowY(), messages, meta

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 77 - "RunTractDialog.tsx"
Cohesion: 0.38
Nodes (4): ParamsList(), Props, formatStartedAt(), RunTractDialog()

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.40
Nodes (4): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.50
Nodes (3): DangerZoneText(), Props, VaultDangerZone()

### Community 80 - "CardMeta.tsx"
Cohesion: 0.50
Nodes (3): ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps

## Knowledge Gaps
- **599 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+594 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `useDialog` connect `useDialog` to `TaskTrackersPage.tsx`, `useBakeError`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `TractIcons.tsx`, `ConnectorChip.tsx`, `tractCanvasLayout.ts`, `VaultItem`, `Router.tsx`, `TractStepTree.tsx`, `useExternalConnections`, `SchemaProperty`, `MobileNotesShell.tsx`, `useErrorToast.ts`, `AuthMiddleware`, `tractSteps.ts`, `admin_users.pb.ts`, `User.ts`, `TractBlockPicker.tsx`, `toTract`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `AuthAPI`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `MembersSection.tsx`, `RunTractDialog.tsx`?**
  _High betweenness centrality (0.099) - this node is a cross-community bridge._
- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `Tracts.ts`, `TractIcons.tsx`, `ConnectorChip.tsx`, `useDialog`, `tractCanvasLayout.ts`, `VaultItem`, `Router.tsx`, `NotesSidebar.tsx`, `ConnectionDetailDialog.tsx`, `index.ts`, `SchemaProperty`, `McpAuthPage.tsx`, `ProviderIcon.tsx`, `StepRow.tsx`, `tractSteps.ts`, `Topbar.tsx`, `TractCanvasBuilder.tsx`, `InviteLinksSection.tsx`, `TractCanvasTopBar.tsx`, `connectionLabel`, `McpKeys.ts`?**
  _High betweenness centrality (0.096) - this node is a cross-community bridge._
- **Why does `useUser` connect `Router.tsx` to `TaskTrackersPage.tsx`, `useBakeError`, `ExternalConnectionInfo`, `couch_instances.pb.ts`, `TractsService`, `Dialog.ts`, `TractIcons.tsx`, `grpcErrors.ts`, `VaultItem`, `McpKeysAPI`, `useExternalConnections`, `Notes.ts`, `index.ts`, `useVaultMutations`, `useErrorToast.ts`, `AuthMiddleware`, `admin_users.pb.ts`, `Topbar.tsx`, `User.ts`, `AuthAPI`, `S3InstanceFormDialog.tsx`, `MembersSection.tsx`?**
  _High betweenness centrality (0.080) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _602 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.0273972602739726 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07712765957446809 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.046511627906976744 - nodes in this community are weakly interconnected._