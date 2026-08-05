# Graph Report - ArtelUI  (2026-08-05)

## Corpus Check
- 588 files · ~152,544 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2705 nodes · 6589 edges · 128 communities (120 shown, 8 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 23 edges (avg confidence: 0.67)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `701af161`
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
- package.json
- ConnectedContent.tsx
- KebabMenu.tsx
- useTheme.ts
- RoadmapCanvasNode.tsx
- Tabs.tsx
- TopbarDrawerCloseButton.tsx
- Textarea.tsx
- BYOKSection.tsx
- TopbarMobileTrigger.tsx

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 202 edges
2. `cn` - 162 edges
3. `useBakeError()` - 150 edges
4. `useUser` - 104 edges
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
- `Props` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx → src/app/api/artel/external_connections.pb.ts
- `Props` --references--> `McpKeyInfo`  [EXTRACTED]
  src/widgets/McpKeyCard/McpKeyCard.tsx → src/app/api/artel/mcp_keys.pb.ts
- `ConnectionOptionListProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/components/ConnectionOptionList/ConnectionOptionList.tsx → src/app/api/artel/mcp_keys.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/workbench/WorkbenchPage.tsx -> src/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx -> src/app/routing/Router.tsx`
- 5-file cycle: `src/app/routing/Router.tsx -> src/pages/tract-templates/TractTemplatesListPage.tsx -> src/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx -> src/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx -> src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx -> src/app/routing/Router.tsx`

## Communities (128 total, 8 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.02
Nodes (98): Absent, ActionStep, BaseTractStep, ConditionStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger (+90 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (29): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+21 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.03
Nodes (63): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+55 more)

### Community 3 - "useBakeError"
Cohesion: 0.17
Nodes (18): PublicDocsNoteItem, DocsFolderNode(), DocsFolderNodeProps, DocsSidebar(), DocsSidebarProps, DocsTreeItem(), DocsTreeItemProps, ArrowIcon() (+10 more)

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.04
Nodes (53): Absent, AddAnthropicConnection, AddAnthropicConnectionResponse, AddEmailConnection, AddEmailConnectionResponse, AddGenericConnection, AddGenericConnectionResponse, AddGitlabConnection (+45 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (42): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+34 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.14
Nodes (24): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenArrayItems() (+16 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.14
Nodes (21): useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor(), providerLabel(), triggerChipLabel(), TriggerRow() (+13 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.07
Nodes (14): AddAnthropicConnectionRequest, AddEmailConnectionRequest, AddGenericConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, CheckAnthropicConnectionRequest, CheckAnthropicConnectionResponse, ExternalConnectionsAPI (+6 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.16
Nodes (29): Props, Props, TractStepTreeProps, Props, Props, Props, GroupBody(), Props (+21 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.10
Nodes (20): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+12 more)

### Community 13 - "Dialog.ts"
Cohesion: 0.15
Nodes (11): useServerStatus(), Path, Router(), routes, ConnectionSection(), Props, ConnectionOptionListProps, GoogleOAuthCallbackPage() (+3 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.06
Nodes (31): cn, KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it, TODO: chures has no tab primitive yet, drop this wrapper once it does, Tabs(), TabsProps (+23 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.07
Nodes (38): AdminCouchAPI, AdminUsersAPI, useBakeError(), HomeLayout(), S3InstanceFormDialog(), S3InstancesActionBar(), BinaryStorageToggle(), Props (+30 more)

### Community 17 - "compilerOptions"
Cohesion: 0.10
Nodes (22): SetTractsState, triggerSourcesQueryKey, triggersQueryKey, useTracts, sleep(), ToolStep(), Props, Step (+14 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.19
Nodes (9): useMcpKeys, Props, SearchInput(), ConnectionStep(), Props, ContentSegment(), ToolDetail(), ConnectionStep() (+1 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.08
Nodes (31): CheckEmailConnectionRequest, CheckGitlabConnectionRequest, CheckTrelloConnectionRequest, UserErrors, CheckStatus, EmailCheckButton(), EmailCheckButtonProps, CheckStatus (+23 more)

### Community 21 - "useDialog"
Cohesion: 0.18
Nodes (11): isNonEmptyBranch(), JsonBlock(), Props, Props, statusClass(), StepRow(), Props, cap() (+3 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.08
Nodes (20): Absent, BaseCompleteSetupRequest, CompleteSetup, CompleteSetupRequest, CompleteSetupResponse, GetStatus, GetStatusRequest, GetStatusResponse (+12 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.12
Nodes (24): ConnectorPath(), ConnectorPathProps, ParallelBoxes(), Props, Props, TractCanvasArea(), useTractCanvasDrag(), cap() (+16 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (24): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+16 more)

### Community 25 - "VaultItem"
Cohesion: 0.13
Nodes (14): VaultItem, Props, VaultCardConnBar(), Props, VaultCardFront(), Props, VaultCardHeader(), VaultCardStatus() (+6 more)

### Community 26 - "devDependencies"
Cohesion: 0.08
Nodes (26): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+18 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.11
Nodes (7): CommunityConnectorInfo, CreateMcpKeyResponse, McpKeyInfo, McpKeysAPI, McpKeysState, IMcpKeysService, McpKeysService

### Community 29 - "Router.tsx"
Cohesion: 0.50
Nodes (4): KNOWN_STATES, LoginRelayScreen(), LoginState, Props

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.16
Nodes (13): DialogHead(), DialogHeadProps, VaultField(), VaultFieldProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), MainScreen() (+5 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.11
Nodes (21): ScriptLanguage, TemplateSource, Props, CONDITION_OPS, ConditionRowProps, Props, ActionBody(), CONDITION_OPS (+13 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.14
Nodes (12): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, PublicBadge() (+4 more)

### Community 33 - "useExternalConnections"
Cohesion: 0.14
Nodes (18): paramCompletionSource(), scriptEditorTheme, scriptHighlightStyle, Props, ScriptCodeSection(), Props, addInputParam(), addOutputParam() (+10 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.12
Nodes (19): useNotes, ArtelLogoIcon(), ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps, MobileDrawer(), MobileDrawerProps, VaultOption (+11 more)

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

### Community 39 - "useVaultMutations"
Cohesion: 0.10
Nodes (4): VaultsAPI, useVaultMutations(), IWorkbenchService, WorkbenchService

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.11
Nodes (27): CreateTriggerRequest, TractRunItem, TractTemplateItem, TriggerItem, TriggerSourceItem, TractsState, safeParseJson(), PresetDetails() (+19 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.15
Nodes (13): ExternalConnectionInfo, MAIL_DOMAIN_ICONS, mailProviderIcon(), AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps, DialogHead() (+5 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.50
Nodes (6): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleConnectionContent()

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.09
Nodes (22): Input(), AddTaskLinkDialog(), Props, RELATION_LABEL, RELATION_OPTIONS, RoadmapLinkTarget, WritableRelation, CredentialRow() (+14 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.14
Nodes (10): AdminSystemSettingsAPI, GetSettings, GetSettingsRequest, GetSettingsResponse, UpdateAuthMethods, UpdateAuthMethodsRequest, UpdateAuthMethodsResponse, UpdateRegistrationMode (+2 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.10
Nodes (21): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, framer-motion (+13 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.19
Nodes (9): FormField(), Props, HeroSegment(), HeroSegmentProps, CreateVaultDialog(), HomePage(), CreateKeyDialog(), ToolboxPage() (+1 more)

### Community 48 - "dependencies"
Cohesion: 0.23
Nodes (15): SelectOption(), execute(), fetchBoards(), fetchCards(), fetchLists(), TrelloBoardLite, TrelloCardLite, TrelloListLite (+7 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.12
Nodes (15): GetVaultBySlug, GetVaultBySlugRequest, GetVaultBySlugResponse, PublicDocsGetNote, PublicDocsGetNoteRequest, PublicDocsGetNoteResponse, PublicDocsListFolders, PublicDocsListFoldersRequest (+7 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.16
Nodes (13): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), isJsonValue(), TaskTrackerCell(), ResultViewMode (+5 more)

### Community 51 - "AuthMiddleware"
Cohesion: 0.36
Nodes (5): AppConfigState, useAppConfig, pingServer(), LoginContent(), UnsecureBanner()

### Community 52 - "tractSteps.ts"
Cohesion: 0.15
Nodes (26): DangerZone(), InsertConflictDialog(), Props, appendStep(), branchArray(), buildStepFromDraft(), collapseThinParallels(), collectAllStepIds() (+18 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.10
Nodes (21): ArtelUserDetails, ArtelUserEntry, ChangeArtelUserPassword, ChangeArtelUserPasswordRequest, ChangeArtelUserPasswordResponse, CreateArtelUser, CreateArtelUserRequest, CreateArtelUserResponse (+13 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.19
Nodes (14): ConnectionsIcon(), base, NavIconProps, LogoutIcon(), NotesIcon(), ToolboxIcon(), TractsIcon(), VaultsIcon() (+6 more)

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
Cohesion: 0.10
Nodes (15): DeleteDockerHost, DeleteDockerHostRequest, DeleteDockerHostResponse, DockerHostsAPI, GetDockerHost, GetDockerHostRequest, ListDockerHosts, ListDockerHostsRequest (+7 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.18
Nodes (8): AnthropicIcon(), EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), TrelloIcon(), TODO: placeholder glyph for providers without a dedicated brand icon yet - repla, UnknownProviderIcon()

### Community 61 - "toTract"
Cohesion: 0.36
Nodes (4): VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.13
Nodes (18): connectionLabel(), ActionCard(), Props, InsertRow(), Props, LlmCallCard(), Props, Props (+10 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.22
Nodes (11): UserState, PasswordLoginForm(), PasswordLoginFormProps, RegisterForm(), RegisterFormProps, LoginContentProps, Mode, AuthService (+3 more)

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.05
Nodes (51): Props, RunButton(), Props, RunStatusBadge(), formatRelative(), Props, RunStatusDot(), LogicCell() (+43 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.16
Nodes (16): ResizeHandle(), ResizeHandleProps, buildLogLines(), formatDuration(), formatTime(), LogLineKind, stepMeta(), RunLog() (+8 more)

### Community 67 - "AuthAPI"
Cohesion: 0.13
Nodes (20): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), moveEntryAndRefresh(), NotesState, remapSelectedPath(), requireVaultId(), DropZone() (+12 more)

### Community 68 - "connectionLabel"
Cohesion: 0.07
Nodes (35): DialogManager, useDialog, useDialogKeyboard(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), Props, TriggerPanel() (+27 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.15
Nodes (10): ByokTabIcon(), CommunitySection(), CommunityTabIcon(), ExternalConnectionsTabIcon(), ConnectionsPage(), ConnectionsTab, resolveTab(), CommunityConnectorCard() (+2 more)

### Community 71 - "McpKeys.ts"
Cohesion: 0.29
Nodes (6): GetCouchInstanceResponse, InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps

### Community 72 - "ResultView.tsx"
Cohesion: 0.15
Nodes (20): AddTriggerDialog(), STEP_SCREENS, AddTriggerDialogContext, AddTriggerDialogState, AddTriggerStep, emptySchemaField(), FIELD_TYPES, fieldsToSchemaNode() (+12 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.10
Nodes (33): Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), cardToNode(), createRoadmapGraph(), expandNode() (+25 more)

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.24
Nodes (10): TaskTrackerTableHead(), DisplayTaskTrackerTables(), TrelloTableWidget(), RESULT_VIEW_WIDGETS, ResultViewWidgetEntry, ResultViewWidgetProps, isPlainObject(), TaskTrackerTable (+2 more)

### Community 75 - "dialog-scrollable.js"
Cohesion: 0.46
Nodes (7): allRules(), dialogScrollable(), directDeclsOf(), findScrollTarget(), isOverflowY(), messages, meta

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.18
Nodes (11): VaultInviteItem, VaultMemberInfo, vaultsQueryKey, Props, InviteRow(), Props, Props, Props (+3 more)

### Community 77 - "RunTractDialog.tsx"
Cohesion: 0.12
Nodes (15): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, formatStartedAt(), Props (+7 more)

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.21
Nodes (9): useIsMobileNav(), BrandMarkIcon(), TopbarBrand(), TopbarMobileDrawer(), MenuRect, TopbarUserMenu(), TopbarUserMenuPill(), TopbarUserMenuPillProps (+1 more)

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.17
Nodes (5): AuthAPI, queryClient, originalFetch, refreshTokens(), SKIP_REFRESH_PATHS

### Community 80 - "CardMeta.tsx"
Cohesion: 0.20
Nodes (13): TractTemplatesState, useTractTemplates, BrowseTemplatesDialog(), Props, TemplateRow(), Props, ListScreen(), Props (+5 more)

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.39
Nodes (5): Tab, AdminHero(), AdminHeroProps, TabBar(), TabBarProps

### Community 94 - "AuthMiddleware"
Cohesion: 0.12
Nodes (10): apiPrefix(), csrfHeader(), getCsrfToken(), InitReq, TelegramLoginResponse, AuthMiddleware, clearLocalStorage(), fromLocalStorage() (+2 more)

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.11
Nodes (16): BreadcrumbBarProps, Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), CopyIcon(), ErrorDotIcon(), PencilIcon() (+8 more)

### Community 96 - "UsersTab.tsx"
Cohesion: 0.31
Nodes (5): GenericToolIcon(), TODO: placeholder glyph for tool actions without a dedicated icon yet (non-smtp/, ImapIcon(), SmtpIcon(), ToolRow()

### Community 97 - "AdminCouchAPI"
Cohesion: 0.23
Nodes (10): MomCandidate, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, AddConnectionScreen(), AddConnectionScreenProps, ConnectionFilterRow() (+2 more)

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.21
Nodes (8): CardChips(), Props, CardHeader(), Props, CardMeta(), formatDate(), Props, Props

### Community 100 - "RunLog.tsx"
Cohesion: 0.13
Nodes (13): ImapOperation, ImapToolAction, McpToolInfo, SmtpOperation, SmtpToolAction, ToolParamDef, ConnectionPicker(), ImapActionView() (+5 more)

### Community 101 - "package.json"
Cohesion: 0.15
Nodes (11): GetDockerHostResponse, DockerHostFormDialog(), DockerHostFormDialogProps, DockerHostFormDialogSaveData, DockerHostListProps, DockerHostRow(), DockerHostRowProps, DockerHostsActionBar() (+3 more)

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.27
Nodes (7): useTemplateConnections(), UseTemplateConnectionsResult, InstantiateTemplateDialog(), Props, ConnectionRequirement, ConnectionRequirementKind, requiredConnections()

### Community 103 - "CouchInstancesAPI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 104 - ".getRun"
Cohesion: 0.06
Nodes (45): ExternalProvider, useExternalConnections, ProviderIcon(), LlmConnectionStep(), Props, ConnectGenericDialog(), CredentialField, ConnectedContent() (+37 more)

### Community 105 - "UserList.tsx"
Cohesion: 0.18
Nodes (11): scripts, build, build:ui, dev, gen, lint, lint:css, lint:js (+3 more)

### Community 106 - "Vaults.ts"
Cohesion: 0.36
Nodes (4): ConnectorChip(), GenericChip(), ProviderChip(), PROVIDER_CHIP_CLASS

### Community 107 - "AuthFetchInterceptor.ts"
Cohesion: 0.48
Nodes (4): McpConnectorInfo, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps

### Community 108 - "EmailCheckButton.tsx"
Cohesion: 0.40
Nodes (4): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps

### Community 109 - "GitlabCheckButton.tsx"
Cohesion: 0.15
Nodes (9): ArtelAPI, Version, VersionRequest, VersionResponse, PublicDocsAPI, DocsNoteViewer(), DocsNoteViewerProps, DocsPage() (+1 more)

### Community 110 - "TrelloCheckButton.tsx"
Cohesion: 0.50
Nodes (3): CouchUserEntry, UserListProps, UserRowProps

### Community 111 - "TopbarDrawerCloseButton.tsx"
Cohesion: 0.23
Nodes (6): GetS3InstanceResponse, Props, S3ToggleFields(), Props, S3InstanceRow(), TestStatus

### Community 112 - "InsertConflictDialog.tsx"
Cohesion: 0.24
Nodes (6): DangerZoneText(), Props, MembersSection(), Props, Props, VaultDangerZone()

### Community 113 - "TopbarMobileTrigger.tsx"
Cohesion: 0.43
Nodes (6): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind

### Community 117 - "requiredConnections.ts"
Cohesion: 0.12
Nodes (20): useVaults(), useWorkbench(), useWorkbenchMutations(), workbenchQueryKey(), VaultChipDisplayProps, VaultOptionList(), VaultOptionListProps, ContentSegment() (+12 more)

### Community 118 - "package.json"
Cohesion: 0.33
Nodes (5): name, private, trustedDependencies, type, version

### Community 120 - "KebabMenu.tsx"
Cohesion: 0.24
Nodes (11): RegistrationMode, RegistrationModeSection(), RegistrationModeSectionProps, AuthMethodsScreen(), OPTIONS, RegistrationModeScreen(), SetupWizardContext, SetupWizardState (+3 more)

### Community 121 - "useTheme.ts"
Cohesion: 0.53
Nodes (4): applyTheme(), Theme, useTheme(), TopbarThemeToggle()

### Community 122 - "RoadmapCanvasNode.tsx"
Cohesion: 0.60
Nodes (4): boardListLabel(), Props, RoadmapCanvasNode(), RoadmapLayoutNode

### Community 123 - "Tabs.tsx"
Cohesion: 0.11
Nodes (21): dompurify, NoteMode, usePortrait(), DesktopNotesShellProps, VaultOption, MobileNotesShellProps, NoteViewer(), NoteViewerProps (+13 more)

### Community 124 - "TopbarDrawerCloseButton.tsx"
Cohesion: 0.50
Nodes (3): TopbarCloseIcon(), TopbarDrawerCloseButton(), TopbarDrawerCloseButtonProps

### Community 125 - "Textarea.tsx"
Cohesion: 0.33
Nodes (5): Props, TODO: chures has no multiline variant yet, drop this wrapper once it does, Textarea(), DockerHostTlsFields(), Props

### Community 127 - "TopbarMobileTrigger.tsx"
Cohesion: 0.50
Nodes (3): TopbarHamburgerIcon(), TopbarMobileTrigger(), TopbarMobileTriggerProps

## Knowledge Gaps
- **797 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+792 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Tracts.ts`, `ToolboxPage.tsx`, `useDialog`, `tractCanvasLayout.ts`, `VaultItem`, `BreadcrumbBar.tsx`, `NotesSidebar.tsx`, `index.ts`, `SchemaProperty`, `McpAuthPage.tsx`, `useErrorToast.ts`, `dependencies`, `StepRow.tsx`, `Topbar.tsx`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `scripts`, `AuthAPI`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `ResultView.tsx`, `MembersSection.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `DbAccessList.tsx`, `GoogleSheetsSpreadsheetSection.tsx`, `ConnectForm.tsx`, `AdminCouchAPI`, `RunLog.tsx`, `.getRun`, `Vaults.ts`, `TopbarMobileTrigger.tsx`, `requiredConnections.ts`, `KebabMenu.tsx`, `RoadmapCanvasNode.tsx`, `Textarea.tsx`, `TopbarMobileTrigger.tsx`?**
  _High betweenness centrality (0.117) - this node is a cross-community bridge._
- **Why does `useDialog` connect `connectionLabel` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `Dialog.ts`, `TractIcons.tsx`, `compilerOptions`, `ToolboxPage.tsx`, `TractStepTree.tsx`, `ConnectionDetailDialog.tsx`, `McpAuthPage.tsx`, `useErrorToast.ts`, `dependencies`, `AuthMiddleware`, `tractSteps.ts`, `admin_users.pb.ts`, `User.ts`, `toTract`, `NotesPage.tsx`, `VaultCard.tsx`, `TractCanvasTopBar.tsx`, `AuthAPI`, `S3InstanceFormDialog.tsx`, `ResultView.tsx`, `MembersSection.tsx`, `RunTractDialog.tsx`, `CardMeta.tsx`, `package.json`, `MobileDrawer.tsx`, `.getRun`, `TopbarDrawerCloseButton.tsx`, `Tabs.tsx`?**
  _High betweenness centrality (0.108) - this node is a cross-community bridge._
- **Why does `useUser` connect `TractIcons.tsx` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `ExternalConnectionInfo`, `Dialog.ts`, `compilerOptions`, `grpcErrors.ts`, `VaultItem`, `McpKeysAPI`, `Notes.ts`, `index.ts`, `useVaultMutations`, `CreateNoteDialog.tsx`, `McpAuthPage.tsx`, `useErrorToast.ts`, `Topbar.tsx`, `User.ts`, `toTract`, `VaultCard.tsx`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `McpKeys.ts`, `MembersSection.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `DbAccessList.tsx`, `VaultDangerZone.tsx`, `CardMeta.tsx`, `package.json`, `.getRun`, `TopbarDrawerCloseButton.tsx`, `requiredConnections.ts`?**
  _High betweenness centrality (0.067) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _803 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.022626262626262626 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0797872340425532 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.03125 - nodes in this community are weakly interconnected._