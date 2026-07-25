# Graph Report - ArtelUI  (2026-07-24)

## Corpus Check
- 521 files · ~140,190 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2384 nodes · 5814 edges · 116 communities (109 shown, 7 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 22 edges (avg confidence: 0.68)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b291830a`
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

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 184 edges
2. `cn` - 144 edges
3. `useBakeError()` - 120 edges
4. `useUser` - 92 edges
5. `useExternalConnections` - 49 edges
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
- `AccountsSectionProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/AccountsSection.tsx → src/app/api/artel/external_connections.pb.ts
- `EmailConnectionRowProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/components/EmailConnectionRow/EmailConnectionRow.tsx → src/app/api/artel/external_connections.pb.ts
- `AnthropicCheckButtonProps` --references--> `CheckAnthropicConnectionRequest`  [EXTRACTED]
  src/dialogs/ManageAnthropicDialog/components/AnthropicCheckButton/AnthropicCheckButton.tsx → src/app/api/artel/external_connections.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 5-file cycle: `src/app/routing/Router.tsx -> src/pages/tract-templates/TractTemplatesListPage.tsx -> src/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx -> src/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx -> src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx -> src/app/routing/Router.tsx`

## Communities (116 total, 7 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.02
Nodes (103): Absent, ActionStep, BaseTractStep, ConditionStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger (+95 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (28): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+20 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.04
Nodes (54): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+46 more)

### Community 3 - "useBakeError"
Cohesion: 0.15
Nodes (13): useIsMobileNav(), applyTheme(), Theme, useTheme(), BrandMarkIcon(), TopbarBrand(), TopbarMobileDrawer(), TopbarThemeToggle() (+5 more)

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.04
Nodes (51): Absent, AddAnthropicConnection, AddAnthropicConnectionResponse, AddEmailConnection, AddEmailConnectionResponse, AddGitlabConnection, AddGitlabConnectionResponse, AddSpreadsheet (+43 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (39): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+31 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.13
Nodes (22): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenArrayItems() (+14 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.16
Nodes (14): triggersQueryKey, useTrigger(), useTriggerSources(), providerEnumFor(), DialogHeaderWithClose(), CreateScreen(), KIND_OPTIONS, KindSettingsProps (+6 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.06
Nodes (32): AdminCouchAPI, ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, CouchUserEntry, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse (+24 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.12
Nodes (40): MomCandidate, ScriptLanguage, Props, ToolStep(), Props, TractStepTreeProps, ActionBody(), Props (+32 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.06
Nodes (34): CouchInstancesAPI, DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceResponse, GetCouchInstanceStatus (+26 more)

### Community 12 - "TractsService"
Cohesion: 0.05
Nodes (13): TractsAPI, TractsService, definitionFromProto(), definitionToProto(), parseSchema(), safeParseJson(), toRun(), toRunStep() (+5 more)

### Community 13 - "Dialog.ts"
Cohesion: 0.16
Nodes (21): SetTractsState, TractsState, triggerSourcesQueryKey, TractTemplatesState, Props, Props, TemplateRow(), Props (+13 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.07
Nodes (29): cn, DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it (+21 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.18
Nodes (11): Input(), buildEmailRequest(), EmailAddDialog(), buildEmailRequest(), EmailEditDialog(), HostPortRowProps, useMailServerSuggestion(), UseDefaultButton() (+3 more)

### Community 17 - "compilerOptions"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.18
Nodes (9): McpToolInfo, SmtpOperation, SmtpToolAction, ToolParamDef, ParamRow(), ParamsList(), RunScreens(), SmtpActionView() (+1 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.21
Nodes (5): ExternalConnectionInfo, AccountsSection(), AccountsSectionProps, TrelloConnectionRow(), TrelloConnectionRowProps

### Community 20 - "grpcErrors.ts"
Cohesion: 0.07
Nodes (32): CheckEmailConnectionRequest, CheckGitlabConnectionRequest, CheckTrelloConnectionRequest, UserErrors, CheckStatus, EmailCheckButton(), EmailCheckButtonProps, CheckStatus (+24 more)

### Community 21 - "useDialog"
Cohesion: 0.08
Nodes (28): isNonEmptyBranch(), JsonBlock(), Props, Props, statusClass(), StepRow(), Props, Props (+20 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.09
Nodes (27): useServerStatus(), HomeLayout(), Path, Router(), routes, HeroSegment(), HeroSegmentProps, useUser (+19 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.11
Nodes (25): NodeChips(), ConnectorPath(), ConnectorPathProps, ParallelBoxes(), Props, Props, TractCanvasArea(), useTractCanvasDrag() (+17 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (23): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+15 more)

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
Cohesion: 0.12
Nodes (7): CreateMcpKeyResponse, McpKeyInfo, McpKeysAPI, McpKeysState, IMcpKeysService, McpKeysService, Props

### Community 29 - "Router.tsx"
Cohesion: 0.12
Nodes (25): connectionLabel(), SelectOption(), ConnectionStep(), Props, ConnectionSection(), Props, ConnectionOptionList(), ConnectionOptionListProps (+17 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.16
Nodes (13): DialogHead(), DialogHeadProps, VaultField(), VaultFieldProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), MainScreen() (+5 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.16
Nodes (10): Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), CopyIcon(), ErrorDotIcon(), PencilIcon(), SpinnerIcon() (+2 more)

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
Cohesion: 0.22
Nodes (8): ListPrompts, ListPromptsRequest, ListPromptsResponse, PromptId, PromptItem, PromptsAPI, FastSetupDialog(), Props

### Community 38 - "s3_instances.pb.ts"
Cohesion: 0.07
Nodes (27): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, GetS3InstanceResponse, ListS3Instances, ListS3InstancesRequest (+19 more)

### Community 39 - "useVaultMutations"
Cohesion: 0.09
Nodes (9): VaultsAPI, useVaultMutations(), BinaryStorageToggle(), Props, VaultSettingsSection(), CreateVaultDialog(), JoinVaultPage(), IVaultService (+1 more)

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.19
Nodes (11): categoryLabel(), PROVIDER_ENUM_BY_KEY, providerLabel(), triggerChipLabel(), TriggerRow(), WebhookPicker(), WebhookPickerProps, PresetDetails() (+3 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.11
Nodes (17): ArtelUI Frontend Rules, Async style, Buttons, Component hierarchy, Component Structure, CSS Modules, Dialog shells must scroll internally, Error and Confirmation Handling (+9 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.18
Nodes (12): ArtelUserDetails, UserSession, useBakeError(), ArtelUserDetailDialog(), ArtelUserDetailDialogProps, UserSessionsDialog(), UserSessionsDialogProps, ChangePasswordDialog() (+4 more)

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (10): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditorProps (+2 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.11
Nodes (24): AddTaskLinkDialog(), Props, RELATION_LABEL, RELATION_OPTIONS, RoadmapLinkTarget, WritableRelation, Props, RELATION_CLASS (+16 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.11
Nodes (19): useNotes, CreateNoteDialog(), Props, ArtelLogoIcon(), ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps, MobileDrawer() (+11 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.15
Nodes (9): PRIMITIVE_TYPES, ITEM_TYPES, PARAM_TYPES, ParamType, LogPanelBar(), LogPanelBarProps, CloseIcon(), CollapseIcon() (+1 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.15
Nodes (14): useVaults(), vaultsQueryKey, FormField(), Props, VaultChipDisplayProps, VaultOptionList(), VaultOptionListProps, InstanceFormDialog() (+6 more)

### Community 48 - "dependencies"
Cohesion: 0.09
Nodes (22): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, dompurify (+14 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.16
Nodes (13): MAIL_DOMAIN_ICONS, mailProviderIcon(), ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), ProviderChip() (+5 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.17
Nodes (17): STEP_SCREENS, AddTriggerDialogContext, AddTriggerDialogState, AddTriggerStep, emptySchemaField(), FIELD_TYPES, fieldsToSchemaNode(), fieldsToSchemaProperties() (+9 more)

### Community 51 - "AuthMiddleware"
Cohesion: 0.22
Nodes (9): apiPrefix(), InitReq, Options, TelegramLoginResponse, AppConfigState, useAppConfig, pingServer(), LoginContent() (+1 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.12
Nodes (33): Props, Step, StepDraft, InsertRow(), Props, collectIdsFromRoot(), ConditionCard(), ConditionCardProps (+25 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.16
Nodes (11): AdminUsersAPI, Props, SearchInput(), AdminPage(), Tab, AdminHero(), AdminHeroProps, ArtelUsersTab() (+3 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.19
Nodes (14): ConnectionsIcon(), base, NavIconProps, LogoutIcon(), NotesIcon(), ToolboxIcon(), TractsIcon(), VaultsIcon() (+6 more)

### Community 55 - "TractCanvasBuilder.tsx"
Cohesion: 0.18
Nodes (13): NoteMode, BreadcrumbBarProps, DesktopNotesShellProps, VaultOption, MobileNotesShell(), MobileNotesShellProps, VaultOption, NoteEditor() (+5 more)

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
Cohesion: 0.17
Nodes (15): LogicCell(), LogicCellProps, OptionCell(), ToolCell(), ToolCellProps, Props, TractBlockPicker(), LOGIC_OPTIONS (+7 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.09
Nodes (21): ExternalProvider, ProviderIcon(), ConnectedContent(), ConnectedContentProps, DialogHead(), DialogHeadProps, NotConnectedContentProps, ConnectionDetailDialog() (+13 more)

### Community 61 - "toTract"
Cohesion: 0.43
Nodes (4): VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.14
Nodes (18): ActionCard(), Props, Props, SchemaFieldRow(), Props, SchemaTree(), buildSourcesFor(), buildSources() (+10 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.24
Nodes (7): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), Props, VaultCard()

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
Cohesion: 0.08
Nodes (34): DialogManager, useDialog, useTracts, useTractTemplates, StepPickerDialog(), StepCard(), TractStepTree(), Props (+26 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.23
Nodes (11): AddAnthropicConnectionRequest, AddEmailConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, CheckAnthropicConnectionRequest, CheckAnthropicConnectionResponse, Spreadsheet, ExternalConnectionsState (+3 more)

### Community 70 - "compilerOptions"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 71 - "McpKeys.ts"
Cohesion: 0.24
Nodes (6): Props, RunButton(), Props, RunStatusBadge(), Props, PlayIcon()

### Community 72 - "ResultView.tsx"
Cohesion: 0.13
Nodes (18): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), isJsonValue(), TaskTrackerCell(), TaskTrackerTableHead() (+10 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.16
Nodes (19): cardToNode(), expandNode(), findNodeByShortLink(), MomToolExecutor, RoadmapNode, TrelloCardDetail, TrelloCommentAction, extractTrelloShortLink() (+11 more)

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
Cohesion: 0.21
Nodes (10): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, formatStartedAt(), RunTractDialog() (+2 more)

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.13
Nodes (14): VaultInviteItem, VaultMemberInfo, CreateInviteLinkDialog(), Props, InviteRow(), Props, Props, Props (+6 more)

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.38
Nodes (7): OptIcon(), OptIconProps, OptText(), OptTextProps, OptionCellProps, IconProps, StepColor

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.43
Nodes (3): BranchIcon(), ForkIcon(), LayersIcon()

### Community 93 - "VaultCardHeader.tsx"
Cohesion: 0.14
Nodes (12): useExternalConnections, ComingSoonCardProps, AnthropicCheckButton(), AnthropicCheckButtonProps, CheckStatus, ConnectedContentProps, ConnectForm(), DialogHead() (+4 more)

### Community 94 - "AuthMiddleware"
Cohesion: 0.14
Nodes (4): AuthMiddleware, clearLocalStorage(), fromLocalStorage(), saveToLocalStorage()

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.25
Nodes (9): useDialogKeyboard(), usePortrait(), Props, RenameDialog(), useAutosave(), NotesPage(), buildNotesUrl(), decodeNotePath() (+1 more)

### Community 96 - "UsersTab.tsx"
Cohesion: 0.43
Nodes (6): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind

### Community 97 - "AdminCouchAPI"
Cohesion: 0.20
Nodes (10): McpConnectorInfo, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps (+2 more)

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.18
Nodes (8): AnthropicIcon(), EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), TrelloIcon(), TODO: placeholder glyph for providers without a dedicated brand icon yet - repla, UnknownProviderIcon()

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.22
Nodes (10): useMcpKeys, CardHeader(), Props, CardMeta(), formatDate(), Props, ContentSegment(), McpKeysPage() (+2 more)

### Community 100 - "RunLog.tsx"
Cohesion: 0.31
Nodes (8): LoginRequest, UserState, LoginContentProps, AuthService, IAuthService, Session, StoredAuth, UserInfo

### Community 101 - "package.json"
Cohesion: 0.33
Nodes (5): name, private, trustedDependencies, type, version

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.23
Nodes (6): ToolDetail(), GenericToolIcon(), TODO: placeholder glyph for tool actions without a dedicated icon yet (non-smtp/, ImapIcon(), SmtpIcon(), ToolRow()

### Community 103 - "CouchInstancesAPI"
Cohesion: 0.29
Nodes (7): sleep(), Props, PublishTemplateDialog(), TractCanvasBuilder(), useTractCanvasBuilderHandlers(), useTractRunTracking(), layoutTract()

### Community 104 - ".getRun"
Cohesion: 0.48
Nodes (4): DialogHead(), DialogHeadProps, tokenAuthorizeUrl(), TrelloAddDialog()

### Community 105 - "UserList.tsx"
Cohesion: 0.20
Nodes (9): GetArtelUser, GetArtelUserRequest, GetArtelUserResponse, GetUserSessions, GetUserSessionsRequest, GetUserSessionsResponse, ListArtelUsers, ListArtelUsersRequest (+1 more)

### Community 106 - "Vaults.ts"
Cohesion: 0.31
Nodes (5): ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon(), TreeItemProps

### Community 107 - "AuthFetchInterceptor.ts"
Cohesion: 0.43
Nodes (4): ArtelUserEntry, ArtelUserListProps, ArtelUserRow(), ArtelUserRowProps

### Community 108 - "EmailCheckButton.tsx"
Cohesion: 0.43
Nodes (5): ResultViewMode, ViewModeToggle(), getResultViewWidget(), ResultView(), tryParseJson()

### Community 109 - "GitlabCheckButton.tsx"
Cohesion: 0.33
Nodes (4): ArtelAPI, Version, VersionRequest, VersionResponse

### Community 111 - "TopbarDrawerCloseButton.tsx"
Cohesion: 0.50
Nodes (3): TopbarCloseIcon(), TopbarDrawerCloseButton(), TopbarDrawerCloseButtonProps

### Community 112 - "InsertConflictDialog.tsx"
Cohesion: 0.50
Nodes (3): DangerZoneText(), Props, VaultDangerZone()

### Community 113 - "TopbarMobileTrigger.tsx"
Cohesion: 0.50
Nodes (3): TopbarHamburgerIcon(), TopbarMobileTrigger(), TopbarMobileTriggerProps

## Knowledge Gaps
- **678 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+673 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Tracts.ts`, `TractIcons.tsx`, `ToolboxPage.tsx`, `useDialog`, `ManageKeyDialog.tsx`, `tractCanvasLayout.ts`, `VaultItem`, `Router.tsx`, `NotesSidebar.tsx`, `ConnectionDetailDialog.tsx`, `index.ts`, `McpAuthPage.tsx`, `ProviderIcon.tsx`, `StepRow.tsx`, `tractSteps.ts`, `admin_users.pb.ts`, `Topbar.tsx`, `TractBlockPicker.tsx`, `NotesPage.tsx`, `scripts`, `AuthAPI`, `connectionLabel`, `McpKeys.ts`, `ResultView.tsx`, `MembersSection.tsx`, `DbAccessList.tsx`, `UsersTab.tsx`, `Vaults.ts`, `EmailCheckButton.tsx`, `TopbarMobileTrigger.tsx`, `ConnectionFilterRow.tsx`?**
  _High betweenness centrality (0.147) - this node is a cross-community bridge._
- **Why does `useDialog` connect `connectionLabel` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `admin_couch.pb.ts`, `TractCanvasInspectorBody.tsx`, `couch_instances.pb.ts`, `Dialog.ts`, `TractIcons.tsx`, `ManageKeyDialog.tsx`, `VaultItem`, `Router.tsx`, `TractStepTree.tsx`, `s3_instances.pb.ts`, `useVaultMutations`, `CreateNoteDialog.tsx`, `SchemaProperty`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `useErrorToast.ts`, `ProviderIcon.tsx`, `StepRow.tsx`, `AuthMiddleware`, `tractSteps.ts`, `User.ts`, `InviteLinksSection.tsx`, `TractBlockPicker.tsx`, `toTract`, `scripts`, `AuthAPI`, `RunTractDialog.tsx`, `DbAccessList.tsx`, `VaultCardHeader.tsx`, `ConnectForm.tsx`, `TractCanvasLogPanel.tsx`, `RunLog.tsx`, `CouchInstancesAPI`, `.getRun`, `AuthFetchInterceptor.ts`?**
  _High betweenness centrality (0.103) - this node is a cross-community bridge._
- **Why does `useUser` connect `ManageKeyDialog.tsx` to `TaskTrackersPage.tsx`, `useBakeError`, `addTriggerDialogContext.ts`, `admin_couch.pb.ts`, `couch_instances.pb.ts`, `Dialog.ts`, `TractIcons.tsx`, `grpcErrors.ts`, `VaultItem`, `McpKeysAPI`, `Notes.ts`, `index.ts`, `s3_instances.pb.ts`, `useVaultMutations`, `SchemaProperty`, `McpAuthPage.tsx`, `useErrorToast.ts`, `admin_users.pb.ts`, `Topbar.tsx`, `User.ts`, `S3InstanceFormDialog.tsx`, `TractCanvasLogPanel.tsx`, `RunLog.tsx`?**
  _High betweenness centrality (0.067) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _683 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.02271062271062271 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07922705314009662 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.03636363636363636 - nodes in this community are weakly interconnected._