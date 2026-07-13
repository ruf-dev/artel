# Graph Report - ArtelUI  (2026-07-13)

## Corpus Check
- 430 files · ~119,126 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1920 nodes · 4580 edges · 86 communities (84 shown, 2 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 13 edges (avg confidence: 0.71)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6062818e`
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
- compilerOptions
- McpKeys.ts
- MembersSection.tsx
- LinkScreen.tsx
- Handoff: lint/tooling parity gaps vs. ZpotifyUI
- CardMeta.tsx
- EmptyState.tsx
- eslint.config.js

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 144 edges
2. `cn` - 114 edges
3. `useBakeError()` - 88 edges
4. `useUser` - 82 edges
5. `useExternalConnections` - 39 edges
6. `TractTool` - 38 edges
7. `MomCandidate` - 37 edges
8. `useTracts` - 34 edges
9. `useMcpKeys` - 33 edges
10. `TractStep` - 31 edges

## Surprising Connections (you probably didn't know these)
- `NoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/notes/components/NoteViewer/NoteViewer.tsx → package.json
- `Props` --references--> `McpKeyInfo`  [EXTRACTED]
  src/widgets/McpKeyCard/McpKeyCard.tsx → src/app/api/artel/mcp_keys.pb.ts
- `MomCandidateCardProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/components/CandidateOptionList/components/MomCandidateCard/MomCandidateCard.tsx → src/app/api/artel/mcp_keys.pb.ts
- `ConnectionOptionListProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/components/ConnectionOptionList/ConnectionOptionList.tsx → src/app/api/artel/mcp_keys.pb.ts
- `SelectConnectionScreenProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/screens/SelectConnectionScreen.tsx → src/app/api/artel/mcp_keys.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`

## Communities (86 total, 2 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.03
Nodes (72): Absent, ActionStep, BaseTractStep, ConditionStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger (+64 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.07
Nodes (30): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+22 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.05
Nodes (42): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+34 more)

### Community 3 - "useBakeError"
Cohesion: 0.44
Nodes (5): Tab, AdminHero(), AdminHeroProps, TabBar(), TabBarProps

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.05
Nodes (39): Absent, AddEmailConnection, AddEmailConnectionResponse, AddGitlabConnection, AddGitlabConnectionResponse, AddSpreadsheet, AddSpreadsheetRequest, AddSpreadsheetResponse (+31 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (42): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+34 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.14
Nodes (23): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenProperty() (+15 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.09
Nodes (27): ParamInput(), Props, Props, Props, SchemaFieldRow(), Props, SchemaTree(), AddTriggerDialogContext (+19 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.09
Nodes (10): AddEmailConnectionRequest, AddGitlabConnectionRequest, ExternalConnectionsAPI, Spreadsheet, ExternalConnectionsState, GoogleSheetsSpreadsheetSection(), SpreadsheetRow(), GoogleOAuthCallbackPage() (+2 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.06
Nodes (32): AdminCouchAPI, ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, CouchUserEntry, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse (+24 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.14
Nodes (34): ActionCard(), Props, InsertRow(), Props, buildSourcesFor(), collectIdsFromRoot(), ConditionCard(), ConditionCardProps (+26 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.06
Nodes (35): CouchInstancesAPI, DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceResponse, GetCouchInstanceStatus (+27 more)

### Community 12 - "TractsService"
Cohesion: 0.06
Nodes (34): TractsAPI, TractsState, ParamRow(), ParamsList(), Props, formatStartedAt(), Props, RunTractDialog() (+26 more)

### Community 13 - "Dialog.ts"
Cohesion: 0.17
Nodes (17): triggersQueryKey, useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor(), providerLabel(), WebhookPicker() (+9 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.21
Nodes (17): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+9 more)

### Community 15 - "cn"
Cohesion: 0.08
Nodes (35): MomCandidate, cn, CardHeader(), Props, ConnBarRow(), RowProps, NameField(), Props (+27 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.16
Nodes (13): CheckEmailConnectionRequest, Input(), buildEmailRequest(), EmailAddDialog(), CheckStatus, EmailCheckButton(), EmailCheckButtonProps, buildEmailRequest() (+5 more)

### Community 17 - "compilerOptions"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.18
Nodes (8): McpToolInfo, ToolParamDef, ParamRow(), ParamsList(), ResultPanel(), coerceParams(), ToolDetail(), ToolRow()

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.21
Nodes (10): MAIL_DOMAIN_ICONS, mailProviderIcon(), ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), PROVIDER_CHIP_CLASS (+2 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.12
Nodes (18): BadRequestDetail, DetailType, DetailTypeName, ErrorInfoDetail, FieldViolation, getDetail(), getFieldViolations(), GrpcErrorDetail (+10 more)

### Community 21 - "useDialog"
Cohesion: 0.09
Nodes (24): DialogManager, useDialog, useDialogKeyboard(), ToolStep(), StepCard(), TractStepTree(), useAddTriggerDialog(), DialogHeaderWithClose() (+16 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.20
Nodes (11): useMcpKeys, CardHeader(), Props, ContentSegment(), HeroSegment(), HeroSegmentProps, ContentSegment(), HeroSegment() (+3 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.15
Nodes (16): sleep(), TractCanvasBuilder(), useTractCanvasBuilderHandlers(), useTractRunTracking(), addEdges(), CanvasEdge, CanvasNodeKind, Extent (+8 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.09
Nodes (21): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+13 more)

### Community 25 - "VaultItem"
Cohesion: 0.14
Nodes (13): VaultItem, CardChips(), Props, Props, VaultChip(), ExpertSettingsDrawer(), Props, Props (+5 more)

### Community 26 - "devDependencies"
Cohesion: 0.06
Nodes (35): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+27 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.13
Nodes (6): CreateMcpKeyResponse, McpKeyInfo, McpKeysAPI, McpKeysState, IMcpKeysService, McpKeysService

### Community 29 - "Router.tsx"
Cohesion: 0.14
Nodes (19): useServerStatus(), HomeLayout(), Path, Router(), routes, useUser, queryClient, AdminPage() (+11 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.19
Nodes (11): DialogHead(), DialogHeadProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), MainScreen(), MainScreenProps, SelectConnectionScreen() (+3 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.12
Nodes (14): Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), CopyIcon(), ErrorDotIcon(), PencilIcon(), SpinnerIcon() (+6 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.14
Nodes (14): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, NotesSidebar() (+6 more)

### Community 33 - "useExternalConnections"
Cohesion: 0.20
Nodes (16): useExternalConnections, useBakeError(), ManageEmailDialog(), ConnectedContent(), ConnectedContentProps, ConnectForm(), tokenSettingsUrl(), ManageGitlabDialog() (+8 more)

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
Cohesion: 0.07
Nodes (27): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, GetS3InstanceResponse, ListS3Instances, ListS3InstancesRequest (+19 more)

### Community 39 - "useVaultMutations"
Cohesion: 0.15
Nodes (5): VaultsAPI, useVaultMutations(), DangerZoneText(), Props, VaultDangerZone()

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.09
Nodes (23): triggerSourcesQueryKey, useTracts, Props, Props, TractRunTimeline(), triggerChipLabel(), Props, TriggerPanel() (+15 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.12
Nodes (16): ArtelUI Frontend Rules, Async style, Buttons, Component hierarchy, Component Structure, CSS Modules, Error and Confirmation Handling, Known debt (documented, not migrated) (+8 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.24
Nodes (13): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), NotesState, requireVaultId(), ConflictRow(), Props, ImportConflictsDialog() (+5 more)

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.18
Nodes (9): KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it, ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon() (+1 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.20
Nodes (10): useNotes, ArtelLogoIcon(), ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps, MobileDrawer(), MobileDrawerProps, VaultOption (+2 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.26
Nodes (7): ExternalConnectionInfo, ConnectedContent(), ConnectedContentProps, AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps

### Community 47 - "useErrorToast.ts"
Cohesion: 0.11
Nodes (20): UserErrors, useVaults(), vaultsQueryKey, FormField(), Props, VaultChipDisplayProps, VaultField(), VaultFieldProps (+12 more)

### Community 48 - "dependencies"
Cohesion: 0.14
Nodes (14): dependencies, classnames, framer-motion, grpc-web, marked, react, react-dom, react-router-dom (+6 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.13
Nodes (14): ExternalProvider, ProviderChip(), EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), TrelloIcon(), TODO: placeholder glyph for providers without a dedicated brand icon yet - repla (+6 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.38
Nodes (5): JsonBlock(), Props, Props, statusClass(), StepRow()

### Community 51 - "AuthMiddleware"
Cohesion: 0.11
Nodes (16): apiPrefix(), InitReq, Options, TelegramLoginResponse, UserState, LoginContent(), LoginContentProps, AuthService (+8 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.11
Nodes (26): Props, SearchInput(), Props, Props, Step, StepDraft, StepPickerDialog(), Props (+18 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.09
Nodes (21): AdminUsersAPI, ArtelUserDetails, ArtelUserEntry, GetArtelUser, GetArtelUserRequest, GetArtelUserResponse, GetUserSessions, GetUserSessionsRequest (+13 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.07
Nodes (33): useIsMobileNav(), applyTheme(), Theme, useTheme(), BrandMarkIcon(), ConnectionsIcon(), base, NavIconProps (+25 more)

### Community 55 - "TractCanvasBuilder.tsx"
Cohesion: 0.25
Nodes (8): NoteMode, BreadcrumbBarProps, DesktopNotesShellProps, VaultOption, MobileNotesShellProps, SaveStatus, UseAutosaveOptions, UseAutosaveResult

### Community 56 - "User.ts"
Cohesion: 0.28
Nodes (5): DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, ImportZipDialog(), Props

### Community 57 - "Errors.ts"
Cohesion: 0.17
Nodes (6): ErrorReason, Errors, GrpcError, GrpcErrorDetail, ServiceError, ServiceErrorOption

### Community 58 - "NoteViewer.tsx"
Cohesion: 0.16
Nodes (11): dompurify, ArtelMark(), ArtelMarkProps, ContentSegment, NoteContent(), NoteContentProps, parseWikiLinks(), NoteViewer() (+3 more)

### Community 59 - "InviteLinksSection.tsx"
Cohesion: 0.18
Nodes (11): VaultInviteItem, CreateInviteLinkDialog(), Props, InviteRow(), Props, Props, RoleBadge(), Props (+3 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.22
Nodes (6): DialogHead(), NotConnectedContentProps, ConnectionDetailDialog(), PROVIDER_CONFIG, PROVIDER_KEY, ProviderConfig

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
Nodes (45): Props, RunButton(), Props, RunStatusBadge(), formatRelative(), Props, RunStatusDot(), LogicCell() (+37 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.33
Nodes (4): ArtelAPI, Version, VersionRequest, VersionResponse

### Community 67 - "AuthAPI"
Cohesion: 0.18
Nodes (5): AuthAPI, pingServer(), BinaryStorageToggle(), Props, VaultSettingsSection()

### Community 68 - "connectionLabel"
Cohesion: 0.24
Nodes (9): connectionLabel(), SelectOption(), ConnectionStep(), Props, ConnectionOptionList(), ConnectionOptionListProps, ConnectionPicker(), ConnectionStep() (+1 more)

### Community 70 - "compilerOptions"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 71 - "McpKeys.ts"
Cohesion: 0.20
Nodes (10): McpConnectorInfo, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps (+2 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.38
Nodes (4): VaultMemberInfo, Props, MembersSection(), Props

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.22
Nodes (6): TemplateSource, Props, CONDITION_OPS, ConditionRowProps, Props, TractCondition

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 80 - "CardMeta.tsx"
Cohesion: 0.67
Nodes (3): CardMeta(), formatDate(), Props

## Knowledge Gaps
- **568 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+563 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `Tracts.ts`, `TractIcons.tsx`, `ConnectorChip.tsx`, `useDialog`, `tractCanvasLayout.ts`, `VaultItem`, `BreadcrumbBar.tsx`, `NotesSidebar.tsx`, `ConnectionDetailDialog.tsx`, `index.ts`, `CreateNoteDialog.tsx`, `SchemaProperty`, `McpAuthPage.tsx`, `ProviderIcon.tsx`, `StepRow.tsx`, `tractSteps.ts`, `Topbar.tsx`, `User.ts`, `InviteLinksSection.tsx`, `TractCanvasTopBar.tsx`, `connectionLabel`, `LinkScreen.tsx`?**
  _High betweenness centrality (0.132) - this node is a cross-community bridge._
- **Why does `useDialog` connect `useDialog` to `TaskTrackersPage.tsx`, `admin_couch.pb.ts`, `TractCanvasInspectorBody.tsx`, `couch_instances.pb.ts`, `TractsService`, `Dialog.ts`, `TractIcons.tsx`, `ConnectorChip.tsx`, `ManageKeyDialog.tsx`, `tractCanvasLayout.ts`, `VaultItem`, `Router.tsx`, `TractStepTree.tsx`, `NotesSidebar.tsx`, `useExternalConnections`, `s3_instances.pb.ts`, `CreateNoteDialog.tsx`, `SchemaProperty`, `MobileNotesShell.tsx`, `useErrorToast.ts`, `AuthMiddleware`, `tractSteps.ts`, `admin_users.pb.ts`, `User.ts`, `InviteLinksSection.tsx`, `TractBlockPicker.tsx`, `toTract`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `connectionLabel`, `LinkScreen.tsx`?**
  _High betweenness centrality (0.120) - this node is a cross-community bridge._
- **Why does `useUser` connect `Router.tsx` to `TaskTrackersPage.tsx`, `useBakeError`, `ExternalConnectionInfo`, `admin_couch.pb.ts`, `couch_instances.pb.ts`, `TractsService`, `Dialog.ts`, `TractIcons.tsx`, `ManageKeyDialog.tsx`, `VaultItem`, `McpKeysAPI`, `useExternalConnections`, `Notes.ts`, `index.ts`, `s3_instances.pb.ts`, `useVaultMutations`, `CreateNoteDialog.tsx`, `useErrorToast.ts`, `AuthMiddleware`, `admin_users.pb.ts`, `Topbar.tsx`, `AuthAPI`?**
  _High betweenness centrality (0.069) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _571 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.0273972602739726 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07428571428571429 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.046511627906976744 - nodes in this community are weakly interconnected._