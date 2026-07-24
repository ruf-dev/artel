# Graph Report - ArtelUI  (2026-07-21)

## Corpus Check
- 511 files · ~138,479 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2345 nodes · 5729 edges · 113 communities (110 shown, 3 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 22 edges (avg confidence: 0.68)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e38ad4cf`
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

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 178 edges
2. `cn` - 144 edges
3. `useBakeError()` - 116 edges
4. `useUser` - 92 edges
5. `useExternalConnections` - 43 edges
6. `TractTool` - 42 edges
7. `MomCandidate` - 41 edges
8. `TractStep` - 40 edges
9. `ExternalConnectionInfo` - 36 edges
10. `useMcpKeys` - 36 edges

## Surprising Connections (you probably didn't know these)
- `NoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/notes/components/NoteViewer/NoteViewer.tsx → package.json
- `Props` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx → src/app/api/artel/external_connections.pb.ts
- `Props` --references--> `McpKeyInfo`  [EXTRACTED]
  src/widgets/McpKeyCard/McpKeyCard.tsx → src/app/api/artel/mcp_keys.pb.ts
- `MomCandidateCardProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/components/CandidateOptionList/components/MomCandidateCard/MomCandidateCard.tsx → src/app/api/artel/mcp_keys.pb.ts
- `ConnectionOptionListProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/components/ConnectionOptionList/ConnectionOptionList.tsx → src/app/api/artel/mcp_keys.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 5-file cycle: `src/app/routing/Router.tsx -> src/pages/tract-templates/TractTemplatesListPage.tsx -> src/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx -> src/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx -> src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx -> src/app/routing/Router.tsx`

## Communities (113 total, 3 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.02
Nodes (80): Absent, BaseTractStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger, CreateTriggerRequest, CreateTriggerResponse (+72 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (28): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+20 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.05
Nodes (42): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+34 more)

### Community 3 - "useBakeError"
Cohesion: 0.18
Nodes (11): useIsMobileNav(), applyTheme(), Theme, useTheme(), BrandMarkIcon(), TopbarBrand(), TopbarMobileDrawer(), TopbarThemeToggle() (+3 more)

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.04
Nodes (51): Absent, AddAnthropicConnection, AddAnthropicConnectionResponse, AddEmailConnection, AddEmailConnectionResponse, AddGitlabConnection, AddGitlabConnectionResponse, AddSpreadsheet (+43 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (42): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+34 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.11
Nodes (24): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenArrayItems() (+16 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.13
Nodes (22): useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor(), providerLabel(), triggerChipLabel(), TriggerRow() (+14 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.09
Nodes (3): ExternalConnectionsAPI, GoogleOAuthCallbackPage(), ExternalConnectionsService

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.14
Nodes (34): MomCandidate, Props, ToolStep(), Props, TractStepTreeProps, Props, CONDITION_OPS, Props (+26 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.09
Nodes (22): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+14 more)

### Community 12 - "TractsService"
Cohesion: 0.07
Nodes (4): TractsAPI, TractsService, definitionToProto(), toTract()

### Community 13 - "Dialog.ts"
Cohesion: 0.15
Nodes (23): SetTractsState, TractsState, triggerSourcesQueryKey, triggersQueryKey, sleep(), formatStartedAt(), Props, RunTractDialog() (+15 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.05
Nodes (47): VaultInviteItem, cn, DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, KebabMenu(), KebabMenuItem, Props (+39 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.17
Nodes (12): Input(), AddTaskLinkDialog(), Props, RELATION_LABEL, RELATION_OPTIONS, WritableRelation, buildEmailRequest(), EmailAddDialog() (+4 more)

### Community 17 - "compilerOptions"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.11
Nodes (15): McpToolInfo, ToolParamDef, Props, SearchInput(), ContentSegment(), ParamRow(), ParamsList(), RunScreens() (+7 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.21
Nodes (9): ExternalConnectionInfo, AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps, AccountsSection(), AccountsSectionProps, TrelloConnectionRow() (+1 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.12
Nodes (18): BadRequestDetail, DetailType, DetailTypeName, ErrorInfoDetail, FieldViolation, getDetail(), getFieldViolations(), GrpcErrorDetail (+10 more)

### Community 21 - "useDialog"
Cohesion: 0.12
Nodes (18): isNonEmptyBranch(), JsonBlock(), Props, Props, statusClass(), StepRow(), Props, Props (+10 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.13
Nodes (24): useServerStatus(), HomeLayout(), Path, Router(), routes, HeroSegment(), HeroSegmentProps, useUser (+16 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.11
Nodes (25): NodeChips(), ConnectorPath(), ConnectorPathProps, ParallelBoxes(), Props, Props, TractCanvasArea(), useTractCanvasDrag() (+17 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (24): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+16 more)

### Community 25 - "VaultItem"
Cohesion: 0.16
Nodes (12): VaultItem, CardChips(), Props, Props, VaultChip(), ExpertSettingsDrawer(), Props, Props (+4 more)

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
Cohesion: 0.26
Nodes (14): execute(), fetchBoards(), fetchCards(), fetchLists(), TrelloBoardLite, TrelloCardLite, TrelloListLite, PickBoardStep() (+6 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.11
Nodes (21): McpConnectorInfo, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps (+13 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.11
Nodes (16): BreadcrumbBarProps, Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), CopyIcon(), ErrorDotIcon(), PencilIcon() (+8 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.21
Nodes (8): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, useNotesSearchQuery()

### Community 33 - "useExternalConnections"
Cohesion: 0.14
Nodes (17): paramCompletionSource(), scriptEditorTheme, scriptHighlightStyle, Props, ScriptCodeSection(), addInputParam(), addOutputParam(), uniqueParamName() (+9 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.50
Nodes (6): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleConnectionContent()

### Community 35 - "notes.pb.ts"
Cohesion: 0.06
Nodes (33): CheckImportConflicts, CheckImportConflictsRequest, CheckImportConflictsResponse, CommitImport, CommitImportRequest, CommitImportResponse, DeleteFolder, DeleteFolderRequest (+25 more)

### Community 36 - "Notes.ts"
Cohesion: 0.11
Nodes (5): b64Decode(), ImportResolution, NotesAPI, INotesService, NotesService

### Community 37 - "index.ts"
Cohesion: 0.11
Nodes (14): ArtelAPI, Version, VersionRequest, VersionResponse, ListPrompts, ListPromptsRequest, ListPromptsResponse, PromptId (+6 more)

### Community 38 - "s3_instances.pb.ts"
Cohesion: 0.07
Nodes (22): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, GetS3InstanceResponse, ListS3Instances, ListS3InstancesRequest (+14 more)

### Community 39 - "useVaultMutations"
Cohesion: 0.15
Nodes (5): VaultsAPI, useVaultMutations(), DangerZoneText(), Props, VaultDangerZone()

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.07
Nodes (32): ActionStep, ConditionStep, GroupStep, ParallelStep, ScriptStep, TractCondition, TractDefinition, TractItem (+24 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.11
Nodes (17): ArtelUI Frontend Rules, Async style, Buttons, Component hierarchy, Component Structure, CSS Modules, Dialog shells must scroll internally, Error and Confirmation Handling (+9 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.10
Nodes (24): useBakeError(), FormField(), Props, Props, S3ToggleFields(), S3InstanceFormDialog(), S3InstancesActionBar(), ManageEmailDialog() (+16 more)

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (10): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditorProps (+2 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.09
Nodes (38): RoadmapLinkTarget, Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), boardListLabel(), Props (+30 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.17
Nodes (13): useNotes, ImportZipDialog(), Props, MobileNotesShell(), NotesSidebar(), NotesSidebarProps, VaultOption, useFolderActions() (+5 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.16
Nodes (12): LogPanelBar(), LogPanelBarProps, ResizeHandle(), ResizeHandleProps, clampHeight(), dotClass(), formatDate(), loadStoredHeight() (+4 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.17
Nodes (11): VaultMemberInfo, useVaults(), vaultsQueryKey, VaultChipDisplayProps, VaultField(), VaultFieldProps, VaultOptionList(), VaultOptionListProps (+3 more)

### Community 48 - "dependencies"
Cohesion: 0.10
Nodes (21): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, framer-motion (+13 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.09
Nodes (20): MAIL_DOMAIN_ICONS, mailProviderIcon(), ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), ProviderChip() (+12 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.23
Nodes (9): TractTemplatesState, Props, TemplateRow(), Props, ListScreen(), Props, toTractTemplateSummary(), TractTemplate (+1 more)

### Community 51 - "AuthMiddleware"
Cohesion: 0.12
Nodes (18): apiPrefix(), InitReq, Options, TelegramLoginResponse, AuthAPI, AppConfigState, useAppConfig, pingServer() (+10 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.12
Nodes (31): Props, Step, StepDraft, InsertRow(), Props, collectIdsFromRoot(), ConditionCardProps, StepCard() (+23 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.09
Nodes (22): AdminUsersAPI, ArtelUserDetails, ArtelUserEntry, GetArtelUser, GetArtelUserRequest, GetArtelUserResponse, GetUserSessions, GetUserSessionsRequest (+14 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.19
Nodes (14): ConnectionsIcon(), base, NavIconProps, LogoutIcon(), NotesIcon(), ToolboxIcon(), TractsIcon(), VaultsIcon() (+6 more)

### Community 55 - "TractCanvasBuilder.tsx"
Cohesion: 0.14
Nodes (18): dompurify, NoteMode, usePortrait(), DesktopNotesShellProps, VaultOption, MobileNotesShellProps, VaultOption, NoteEditor() (+10 more)

### Community 56 - "User.ts"
Cohesion: 0.10
Nodes (25): AdminSubscriptionsAPI, EffectiveSubscriptionView, GetUserSubscription, GetUserSubscriptionRequest, GetUserSubscriptionResponse, ListSubscriptionPlans, ListSubscriptionPlansRequest, ListSubscriptionPlansResponse (+17 more)

### Community 57 - "Errors.ts"
Cohesion: 0.17
Nodes (6): ErrorReason, Errors, GrpcError, GrpcErrorDetail, ServiceError, ServiceErrorOption

### Community 58 - "NoteViewer.tsx"
Cohesion: 0.24
Nodes (8): ArtelMark(), ArtelMarkProps, ContentSegment, NoteContent(), NoteContentProps, parseWikiLinks(), WikiChip(), WikiChipProps

### Community 59 - "InviteLinksSection.tsx"
Cohesion: 0.23
Nodes (12): LogicCell(), LogicCellProps, OptionCell(), ToolCell(), ToolCellProps, LOGIC_OPTIONS, LogicOption, rank() (+4 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.08
Nodes (27): ExternalProvider, useExternalConnections, ComingSoonCardProps, ConnectedContent(), ConnectedContentProps, DialogHead(), DialogHeadProps, NotConnectedContentProps (+19 more)

### Community 61 - "toTract"
Cohesion: 0.36
Nodes (4): VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.14
Nodes (19): ScriptLanguage, ActionBody(), ConditionBody(), NameField(), Props, ScriptBody(), buildSources(), replaceStep() (+11 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.21
Nodes (8): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), ContentSegmentProps, Props, VaultCard()

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.17
Nodes (9): ChatIcon(), EditIcon(), EnvelopeIcon(), GlobeIcon(), base, PaperPlaneIcon(), PlusIcon(), SearchIcon() (+1 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.20
Nodes (9): formatRelative(), Props, RunStatusDot(), ChevronRightIcon(), ManualTriggerIcon(), TrashIcon(), WebhookIcon(), Props (+1 more)

### Community 67 - "AuthAPI"
Cohesion: 0.22
Nodes (15): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), moveEntryAndRefresh(), NotesState, remapSelectedPath(), requireVaultId(), ConflictRow() (+7 more)

### Community 68 - "connectionLabel"
Cohesion: 0.09
Nodes (30): DialogManager, useDialog, useTracts, useTractTemplates, StepPickerDialog(), Props, TriggerPanel(), AddTriggerDialog() (+22 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.45
Nodes (8): AddAnthropicConnectionRequest, AddEmailConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, CheckAnthropicConnectionRequest, CheckAnthropicConnectionResponse, ExternalConnectionsState, IExternalConnectionsService

### Community 70 - "compilerOptions"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 71 - "McpKeys.ts"
Cohesion: 0.24
Nodes (6): Props, RunButton(), Props, RunStatusBadge(), Props, PlayIcon()

### Community 72 - "ResultView.tsx"
Cohesion: 0.07
Nodes (39): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), STEP_SCREENS, AddTriggerDialogContext, AddTriggerDialogState (+31 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.20
Nodes (11): connectionLabel(), SelectOption(), ConnectionStep(), Props, ConnectionSection(), Props, ConnectionOptionList(), ConnectionOptionListProps (+3 more)

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
Cohesion: 0.16
Nodes (11): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, Props, SchemaFieldRow() (+3 more)

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.32
Nodes (5): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps, ManageAccessDialogProps

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.38
Nodes (7): OptIcon(), OptIconProps, OptText(), OptTextProps, OptionCellProps, IconProps, StepColor

### Community 80 - "CardMeta.tsx"
Cohesion: 0.28
Nodes (3): Spreadsheet, GoogleSheetsSpreadsheetSection(), SpreadsheetRow()

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.33
Nodes (4): BranchIcon(), CodeIcon(), ForkIcon(), LayersIcon()

### Community 93 - "VaultCardHeader.tsx"
Cohesion: 0.39
Nodes (5): UserErrors, HomePage(), getPreconditionViolations(), GrpcStatusError, isMissingSubscription()

### Community 94 - "AuthMiddleware"
Cohesion: 0.14
Nodes (4): AuthMiddleware, clearLocalStorage(), fromLocalStorage(), saveToLocalStorage()

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.21
Nodes (7): useDialogKeyboard(), Props, PublishTemplateDialog(), CreateNoteDialog(), Props, Props, RenameDialog()

### Community 96 - "UsersTab.tsx"
Cohesion: 0.43
Nodes (6): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind

### Community 97 - "AdminCouchAPI"
Cohesion: 0.22
Nodes (3): AdminCouchAPI, ManageAccessDialog(), UsersTab()

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.27
Nodes (7): GetCouchInstanceResponse, InstanceRow(), InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.17
Nodes (12): useMcpKeys, CardHeader(), Props, CardMeta(), formatDate(), Props, InstantiateTemplateDialog(), Props (+4 more)

### Community 100 - "RunLog.tsx"
Cohesion: 0.38
Nodes (5): ActionCard(), Props, SchemaTree(), buildSourcesFor(), ConditionCard()

### Community 101 - "package.json"
Cohesion: 0.33
Nodes (5): name, private, trustedDependencies, type, version

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.24
Nodes (7): ArtelLogoIcon(), ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps, MobileDrawer(), MobileDrawerProps, VaultOption

### Community 103 - "CouchInstancesAPI"
Cohesion: 0.18
Nodes (4): CouchInstancesAPI, InstancesActionBar(), InstancesActionBarProps, InstancesTab()

### Community 104 - ".getRun"
Cohesion: 0.48
Nodes (4): DialogHead(), DialogHeadProps, tokenAuthorizeUrl(), TrelloAddDialog()

### Community 105 - "UserList.tsx"
Cohesion: 0.48
Nodes (4): CouchUserEntry, UserListProps, UserRow(), UserRowProps

### Community 107 - "AuthFetchInterceptor.ts"
Cohesion: 0.33
Nodes (4): originalFetch, RefreshResponseBody, refreshTokens(), SKIP_REFRESH_PATHS

### Community 108 - "EmailCheckButton.tsx"
Cohesion: 0.50
Nodes (4): CheckEmailConnectionRequest, CheckStatus, EmailCheckButton(), EmailCheckButtonProps

### Community 109 - "GitlabCheckButton.tsx"
Cohesion: 0.50
Nodes (4): CheckGitlabConnectionRequest, CheckStatus, GitlabCheckButton(), GitlabCheckButtonProps

### Community 110 - "TrelloCheckButton.tsx"
Cohesion: 0.50
Nodes (4): CheckTrelloConnectionRequest, CheckStatus, TrelloCheckButton(), TrelloCheckButtonProps

### Community 111 - "TopbarDrawerCloseButton.tsx"
Cohesion: 0.50
Nodes (3): TopbarCloseIcon(), TopbarDrawerCloseButton(), TopbarDrawerCloseButtonProps

### Community 112 - "InsertConflictDialog.tsx"
Cohesion: 0.67
Nodes (3): InsertConflictDialog(), Props, MoveConflictChoice

## Knowledge Gaps
- **663 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+658 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Tracts.ts`, `TractIcons.tsx`, `ToolboxPage.tsx`, `useDialog`, `ManageKeyDialog.tsx`, `tractCanvasLayout.ts`, `VaultItem`, `BreadcrumbBar.tsx`, `NotesSidebar.tsx`, `ConnectionDetailDialog.tsx`, `index.ts`, `McpAuthPage.tsx`, `Tracts.ts`, `ProviderIcon.tsx`, `tractSteps.ts`, `admin_users.pb.ts`, `Topbar.tsx`, `TractBlockPicker.tsx`, `NotesPage.tsx`, `scripts`, `AuthAPI`, `connectionLabel`, `McpKeys.ts`, `ResultView.tsx`, `MembersSection.tsx`, `UsersTab.tsx`, `RunLog.tsx`?**
  _High betweenness centrality (0.120) - this node is a cross-community bridge._
- **Why does `useUser` connect `ManageKeyDialog.tsx` to `TaskTrackersPage.tsx`, `useBakeError`, `addTriggerDialogContext.ts`, `couch_instances.pb.ts`, `Dialog.ts`, `TractIcons.tsx`, `VaultItem`, `McpKeysAPI`, `Notes.ts`, `index.ts`, `s3_instances.pb.ts`, `useVaultMutations`, `SchemaProperty`, `McpAuthPage.tsx`, `useErrorToast.ts`, `AuthMiddleware`, `admin_users.pb.ts`, `Topbar.tsx`, `User.ts`, `toTract`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `DbAccessList.tsx`, `VaultCardHeader.tsx`, `AdminCouchAPI`, `CouchInstancesAPI`, `CouchInstancesAPI`, `UserList.tsx`, `Vaults.ts`, `AuthFetchInterceptor.ts`, `EmailCheckButton.tsx`, `GitlabCheckButton.tsx`, `TrelloCheckButton.tsx`?**
  _High betweenness centrality (0.087) - this node is a cross-community bridge._
- **Why does `useDialog` connect `connectionLabel` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `TractIcons.tsx`, `ManageKeyDialog.tsx`, `VaultItem`, `Router.tsx`, `TractStepTree.tsx`, `s3_instances.pb.ts`, `SchemaProperty`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `ProviderIcon.tsx`, `StepRow.tsx`, `AuthMiddleware`, `tractSteps.ts`, `admin_users.pb.ts`, `TractCanvasBuilder.tsx`, `User.ts`, `InviteLinksSection.tsx`, `TractBlockPicker.tsx`, `toTract`, `scripts`, `AuthAPI`, `MembersSection.tsx`, `DbAccessList.tsx`, `VaultCardHeader.tsx`, `ConnectForm.tsx`, `AdminCouchAPI`, `CouchInstancesAPI`, `TractCanvasLogPanel.tsx`, `MobileDrawer.tsx`, `CouchInstancesAPI`, `.getRun`, `UserList.tsx`, `InsertConflictDialog.tsx`?**
  _High betweenness centrality (0.084) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _668 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.024691358024691357 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07922705314009662 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.046511627906976744 - nodes in this community are weakly interconnected._