# Graph Report - ArtelUI  (2026-08-09)

## Corpus Check
- 610 files · ~157,069 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2833 nodes · 6990 edges · 128 communities (118 shown, 10 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 23 edges (avg confidence: 0.67)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `40c463fe`
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
- DockerHostsAPI
- PostgresActions.tsx
- ParamsList.tsx
- ImportZipDialog.tsx

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 222 edges
2. `useBakeError()` - 172 edges
3. `cn` - 162 edges
4. `useUser` - 103 edges
5. `useExternalConnections` - 72 edges
6. `TractStep` - 46 edges
7. `ExternalConnectionInfo` - 44 edges
8. `TractTool` - 42 edges
9. `MomCandidate` - 41 edges
10. `ExternalProvider` - 39 edges

## Surprising Connections (you probably didn't know these)
- `DocsNoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/docs/components/DocsNoteViewer/DocsNoteViewer.tsx → package.json
- `NoteViewer()` --references--> `dompurify`  [EXTRACTED]
  src/pages/notes/components/NoteViewer/NoteViewer.tsx → package.json
- `AccountsSectionProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/AccountsSection.tsx → src/app/api/artel/external_connections.pb.ts
- `EmailConnectionRowProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/components/EmailConnectionRow/EmailConnectionRow.tsx → src/app/api/artel/external_connections.pb.ts
- `AccountsSectionProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageTrelloDialog/components/AccountsSection/AccountsSection.tsx → src/app/api/artel/external_connections.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/setup-wizard/SetupWizardPage.tsx -> src/pages/setup-wizard/screens/CreateAdminScreen.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/workbench/WorkbenchPage.tsx -> src/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx -> src/app/routing/Router.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarMobileDrawer/TopbarMobileDrawer.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarBrand/TopbarBrand.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarNav/TopbarNav.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarUserMenu/TopbarUserMenu.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 5-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarMobileDrawer/TopbarMobileDrawer.tsx -> src/segments/Topbar/components/TopbarBrand/TopbarBrand.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 5-file cycle: `src/app/routing/Router.tsx -> src/pages/tract-templates/TractTemplatesListPage.tsx -> src/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx -> src/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx -> src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx -> src/app/routing/Router.tsx`

## Communities (128 total, 10 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.02
Nodes (88): Absent, BaseTractStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger, CreateTriggerRequest, CreateTriggerResponse (+80 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (27): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+19 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.03
Nodes (70): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+62 more)

### Community 3 - "useBakeError"
Cohesion: 0.06
Nodes (38): GetVaultBySlug, GetVaultBySlugRequest, GetVaultBySlugResponse, PublicDocsAPI, PublicDocsGetNote, PublicDocsGetNoteRequest, PublicDocsGetNoteResponse, PublicDocsListFolders (+30 more)

### Community 4 - "external_connections.pb.ts"
Cohesion: 0.03
Nodes (69): Absent, AddAnthropicConnection, AddAnthropicConnectionResponse, AddCouchDBConnection, AddCouchDBConnectionResponse, AddEmailConnection, AddEmailConnectionResponse, AddGenericConnection (+61 more)

### Community 5 - "mcp_keys.pb.ts"
Cohesion: 0.05
Nodes (42): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+34 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.08
Nodes (35): Props, SourceGroups(), buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenArrayItems() (+27 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.13
Nodes (20): SetTractsState, triggerSourcesQueryKey, triggersQueryKey, useTrigger(), useTriggerSources(), categoryLabel(), providerEnumFor(), providerLabel() (+12 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.06
Nodes (6): ExternalConnectionInfo, ExternalConnectionsAPI, ConnectionSection(), Props, GoogleOAuthCallbackPage(), ExternalConnectionsService

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.13
Nodes (40): MomCandidate, ScriptLanguage, Props, TractStepTreeProps, ActionBody(), Props, CONDITION_OPS, ConditionBody() (+32 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.10
Nodes (20): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+12 more)

### Community 12 - "TractsService"
Cohesion: 0.06
Nodes (5): TractsAPI, TractsService, toTract(), toTractTemplate(), definitionFromProto()

### Community 13 - "Dialog.ts"
Cohesion: 0.07
Nodes (34): AdminUsersAPI, useServerStatus(), Path, Router(), routes, HeroSegment(), HeroSegmentProps, VaultCardHeader() (+26 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.05
Nodes (50): cn, KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it, TODO: chures has no tab primitive yet, drop this wrapper once it does, Tabs(), TabsProps (+42 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.10
Nodes (21): useMcpKeys, ConnectionStep(), Props, AddTaskLinkDialog(), Props, RELATION_LABEL, RELATION_OPTIONS, RoadmapLinkTarget (+13 more)

### Community 17 - "compilerOptions"
Cohesion: 0.14
Nodes (19): TractsState, sleep(), formatStartedAt(), Props, RunTractDialog(), Props, TractCanvasBuilder(), Deps (+11 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.07
Nodes (35): DialogManager, useDialog, useDialogKeyboard(), LlmConnectionStep(), Props, Props, TriggerPanel(), VaultChip() (+27 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.07
Nodes (36): CheckEmailConnectionRequest, CheckGitlabConnectionRequest, CheckTelegramConnectionRequest, CheckTrelloConnectionRequest, UserErrors, CheckStatus, EmailCheckButton(), EmailCheckButtonProps (+28 more)

### Community 21 - "useDialog"
Cohesion: 0.15
Nodes (14): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind, isNonEmptyBranch(), JsonBlock() (+6 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.09
Nodes (20): Absent, BaseCompleteSetupRequest, CompleteSetup, CompleteSetupRequest, CompleteSetupResponse, GetStatus, GetStatusRequest, GetStatusResponse (+12 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.12
Nodes (23): ConnectorPath(), ConnectorPathProps, ParallelBoxes(), Props, Props, TractCanvasArea(), useTractCanvasDrag(), cap() (+15 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (24): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+16 more)

### Community 25 - "VaultItem"
Cohesion: 0.14
Nodes (13): VaultItem, Props, DangerZoneText(), ExpertSettingsDrawer(), Props, Props, ExpertSettingsSection(), Props (+5 more)

### Community 26 - "devDependencies"
Cohesion: 0.05
Nodes (43): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+35 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.11
Nodes (7): CommunityConnectorInfo, CreateMcpKeyResponse, McpKeyInfo, McpKeysAPI, McpKeysState, IMcpKeysService, McpKeysService

### Community 29 - "Router.tsx"
Cohesion: 0.13
Nodes (16): useNotes, CreateNoteDialog(), Props, ArtelLogoIcon(), MobileDrawer(), MobileDrawerProps, VaultOption, NotesSidebar() (+8 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.19
Nodes (11): DialogHead(), DialogHeadProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), MainScreen(), MainScreenProps, SelectConnectionScreen() (+3 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.10
Nodes (15): GetS3InstanceResponse, FormField(), Props, Props, S3InstanceFormDialog(), S3InstanceRow(), TestStatus, S3InstancesActionBar() (+7 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.21
Nodes (8): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, useNotesSearchQuery()

### Community 33 - "useExternalConnections"
Cohesion: 0.28
Nodes (6): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.40
Nodes (4): PublicBadge(), SidebarTopBar(), SidebarTopBarProps, VaultOption

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
Cohesion: 0.14
Nodes (17): paramCompletionSource(), scriptEditorTheme, scriptHighlightStyle, Props, ScriptCodeSection(), addInputParam(), addOutputParam(), uniqueParamName() (+9 more)

### Community 39 - "useVaultMutations"
Cohesion: 0.09
Nodes (4): VaultsAPI, useVaultMutations(), IWorkbenchService, WorkbenchService

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.11
Nodes (23): Props, SearchInput(), SelectOption(), execute(), fetchBoards(), fetchCards(), fetchLists(), TrelloBoardLite (+15 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.27
Nodes (7): MAIL_DOMAIN_ICONS, mailProviderIcon(), AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps, DialogHead()

### Community 42 - "SchemaProperty"
Cohesion: 0.50
Nodes (6): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleConnectionContent()

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.14
Nodes (13): Input(), CredentialRow(), CredentialRowProps, CredentialField, DialogHeadProps, buildEmailRequest(), EmailAddDialog(), buildEmailRequest() (+5 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.12
Nodes (25): TractTemplatesState, safeParseJson(), Props, statusClass(), StepRow(), PROVIDER_ENUM_BY_KEY, Props, TemplateRow() (+17 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.10
Nodes (21): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, framer-motion (+13 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.11
Nodes (17): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, ListS3Instances, ListS3InstancesRequest, ListS3InstancesResponse (+9 more)

### Community 48 - "dependencies"
Cohesion: 0.12
Nodes (18): useTracts, Props, ToolStep(), Props, Step, StepDraft, StepPickerDialog(), Props (+10 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.23
Nodes (24): AddAnthropicConnectionRequest, AddCouchDBConnectionRequest, AddEmailConnectionRequest, AddGenericConnectionRequest, AddGitlabConnectionRequest, AddOpenAIConnectionRequest, AddPostgresConnectionRequest, AddS3ConnectionRequest (+16 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.16
Nodes (15): ActionStep, ConditionStep, GroupStep, LlmCallStep, ParallelStep, ScriptStep, TractCondition, TractDefinition (+7 more)

### Community 51 - "AuthMiddleware"
Cohesion: 0.15
Nodes (11): apiPrefix(), AuthAPI, AppConfigState, useAppConfig, pingServer(), LoginContent(), CreateAdminScreen(), refreshTokens() (+3 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.14
Nodes (29): InsertRow(), Props, collectIdsFromRoot(), ConditionCardProps, StepCard(), TractStepTree(), appendStep(), branchArray() (+21 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.11
Nodes (17): ArtelUserDetails, ChangeArtelUserPassword, ChangeArtelUserPasswordRequest, ChangeArtelUserPasswordResponse, CreateArtelUser, CreateArtelUserRequest, CreateArtelUserResponse, GetArtelUser (+9 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.08
Nodes (29): useIsMobileNav(), applyTheme(), Theme, useTheme(), HomeLayout(), BrandMarkIcon(), ConnectionsIcon(), base (+21 more)

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
Nodes (18): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), isJsonValue(), TaskTrackerCell(), TaskTrackerTableHead() (+10 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.12
Nodes (10): AnthropicIcon(), EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), OpenAIIcon(), TelegramIcon(), TrelloIcon() (+2 more)

### Community 61 - "toTract"
Cohesion: 0.29
Nodes (5): McpLoginProps, VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.15
Nodes (14): connectionLabel(), ActionCard(), CardHeader(), Props, LlmCallCard(), Props, Props, SchemaTree() (+6 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.22
Nodes (11): UserState, PasswordLoginForm(), PasswordLoginFormProps, RegisterForm(), RegisterFormProps, LoginContentProps, Mode, AuthService (+3 more)

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.05
Nodes (51): Props, RunButton(), formatRelative(), Props, RunStatusDot(), LogicCell(), LogicCellProps, LogicSection() (+43 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.16
Nodes (13): dompurify, NoteMode, LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps (+5 more)

### Community 66 - "scripts"
Cohesion: 0.31
Nodes (8): ResizeHandle(), ResizeHandleProps, clampHeight(), dotClass(), formatDate(), loadStoredHeight(), Props, TractCanvasLogPanel()

### Community 67 - "AuthAPI"
Cohesion: 0.22
Nodes (15): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), moveEntryAndRefresh(), NotesState, remapSelectedPath(), requireVaultId(), ConflictRow() (+7 more)

### Community 68 - "connectionLabel"
Cohesion: 0.15
Nodes (12): CouchDBConnectForm(), PostgresConnectForm(), PostgresConnectFormFieldsProps, SSL_MODE_OPTIONS, S3ConnectForm(), S3ToggleFields(), S3ToggleFieldsProps, StorageCheckButton() (+4 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.27
Nodes (6): ByokTabIcon(), CommunityTabIcon(), ExternalConnectionsTabIcon(), ConnectionsPage(), ConnectionsTab, resolveTab()

### Community 71 - "McpKeys.ts"
Cohesion: 0.29
Nodes (6): GetCouchInstanceResponse, InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps

### Community 72 - "ResultView.tsx"
Cohesion: 0.09
Nodes (27): Props, SchemaFieldRow(), STEP_SCREENS, AddTriggerDialogContext, AddTriggerDialogState, AddTriggerStep, emptySchemaField(), FIELD_TYPES (+19 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.09
Nodes (35): Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), boardListLabel(), Props, RoadmapCanvasNode() (+27 more)

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.15
Nodes (12): BinaryStorageToggle(), Props, PostgresSection(), PostgresSectionProps, StatusType, Props, PublishSlugForm(), slugify() (+4 more)

### Community 75 - "dialog-scrollable.js"
Cohesion: 0.46
Nodes (7): allRules(), dialogScrollable(), directDeclsOf(), findScrollTarget(), isOverflowY(), messages, meta

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.12
Nodes (16): EnablePostgresDatabaseResponse, GetPostgresDatabaseResponse, VaultInviteItem, VaultMemberInfo, useVaults(), vaultsQueryKey, VaultChipDisplayProps, VaultField() (+8 more)

### Community 77 - "RunTractDialog.tsx"
Cohesion: 0.20
Nodes (10): McpConnectorInfo, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps (+2 more)

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.13
Nodes (14): DeleteDockerHost, DeleteDockerHostRequest, DeleteDockerHostResponse, GetDockerHost, GetDockerHostRequest, ListDockerHosts, ListDockerHostsRequest, ListDockerHostsResponse (+6 more)

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.14
Nodes (10): AdminSystemSettingsAPI, GetSettings, GetSettingsRequest, GetSettingsResponse, UpdateAuthMethods, UpdateAuthMethodsRequest, UpdateAuthMethodsResponse, UpdateRegistrationMode (+2 more)

### Community 80 - "CardMeta.tsx"
Cohesion: 0.38
Nodes (6): useTractTemplates, BrowseTemplatesDialog(), ListScreen(), ContentSegment(), TractTemplatesListPage(), TractTemplateCard()

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.35
Nodes (7): AdminPage(), resolveTab(), Tab, AdminHero(), AdminHeroProps, TabBar(), TabBarProps

### Community 94 - "AuthMiddleware"
Cohesion: 0.31
Nodes (5): ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon(), TreeItemProps

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.11
Nodes (16): BreadcrumbBarProps, Mode, BreadcrumbPath(), BreadcrumbPathProps, DesktopNotesShellProps, VaultOption, CheckIcon(), CopyIcon() (+8 more)

### Community 96 - "UsersTab.tsx"
Cohesion: 0.17
Nodes (3): AuthMiddleware, fromLocalStorage(), saveToLocalStorage()

### Community 97 - "AdminCouchAPI"
Cohesion: 0.33
Nodes (6): ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), PROVIDER_CHIP_CLASS

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.21
Nodes (8): CardChips(), Props, CardHeader(), Props, CardMeta(), formatDate(), Props, Props

### Community 100 - "RunLog.tsx"
Cohesion: 0.28
Nodes (3): Spreadsheet, GoogleSheetsSpreadsheetSection(), SpreadsheetRow()

### Community 101 - "package.json"
Cohesion: 0.15
Nodes (11): GetDockerHostResponse, DockerHostFormDialog(), DockerHostFormDialogProps, DockerHostFormDialogSaveData, DockerHostListProps, DockerHostRow(), DockerHostRowProps, DockerHostsActionBar() (+3 more)

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.39
Nodes (4): UseTemplateConnectionsResult, ConnectionRequirement, ConnectionRequirementKind, requiredConnections()

### Community 103 - "CouchInstancesAPI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 104 - ".getRun"
Cohesion: 0.06
Nodes (55): ExternalProvider, useExternalConnections, useBakeError(), LlmKeyConnectedContentProps, CheckRequest, LlmKeyConnectForm(), LlmKeyDialogHead(), LlmKeyDialogHeadProps (+47 more)

### Community 105 - "UserList.tsx"
Cohesion: 0.29
Nodes (5): CheckAnthropicConnectionResponse, CheckStatus, LlmKeyCheckButton(), LlmKeyCheckButtonProps, LlmKeyConnectFormProps

### Community 107 - "AuthFetchInterceptor.ts"
Cohesion: 0.43
Nodes (4): ArtelUserEntry, ArtelUserListProps, ArtelUserRow(), ArtelUserRowProps

### Community 108 - "EmailCheckButton.tsx"
Cohesion: 0.40
Nodes (4): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps

### Community 109 - "GitlabCheckButton.tsx"
Cohesion: 0.50
Nodes (3): ChevronLeftIcon(), DrawerCloseButton(), DrawerCloseButtonProps

### Community 110 - "TrelloCheckButton.tsx"
Cohesion: 0.50
Nodes (3): CouchUserEntry, UserListProps, UserRowProps

### Community 111 - "TopbarDrawerCloseButton.tsx"
Cohesion: 0.19
Nodes (9): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), ContentSegment(), ContentSegmentProps, Props (+1 more)

### Community 113 - "TopbarMobileTrigger.tsx"
Cohesion: 0.40
Nodes (4): AccountsSection(), AccountsSectionProps, TrelloConnectionRow(), TrelloConnectionRowProps

### Community 117 - "requiredConnections.ts"
Cohesion: 0.15
Nodes (17): useWorkbench(), useWorkbenchMutations(), workbenchQueryKey(), KNOWN_STATES, LoginRelayScreen(), LoginState, Props, AuthMode (+9 more)

### Community 119 - "ConnectedContent.tsx"
Cohesion: 0.18
Nodes (9): clearCsrfCookie(), csrfHeader(), getCsrfToken(), InitReq, TelegramLoginResponse, queryClient, forceLogout(), originalFetch (+1 more)

### Community 120 - "KebabMenu.tsx"
Cohesion: 0.18
Nodes (12): RegistrationMode, SetupWizardAPI, AuthMethodsScreen(), OPTIONS, RegistrationModeScreen(), TokenEntryScreen(), SetupWizardContext, SetupWizardState (+4 more)

### Community 121 - "useTheme.ts"
Cohesion: 0.33
Nodes (5): Props, TODO: chures has no multiline variant yet, drop this wrapper once it does, Textarea(), DockerHostTlsFields(), Props

### Community 122 - "RoadmapCanvasNode.tsx"
Cohesion: 0.33
Nodes (4): ArtelAPI, Version, VersionRequest, VersionResponse

### Community 123 - "Tabs.tsx"
Cohesion: 0.20
Nodes (11): usePortrait(), Props, RenameDialog(), useAutosave(), NotesPage(), buildNotesUrl(), decodeNotePath(), encodeNotePath() (+3 more)

### Community 126 - "ParamsList.tsx"
Cohesion: 0.14
Nodes (12): ImapOperation, ImapToolAction, McpToolInfo, SmtpOperation, SmtpToolAction, ToolParamDef, ImapActionView(), ParamRow() (+4 more)

### Community 127 - "ImportZipDialog.tsx"
Cohesion: 0.33
Nodes (5): DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, ImportZipDialog(), Props

## Knowledge Gaps
- **834 isolated node(s):** `localPlugin`, `name`, `private`, `license`, `version` (+829 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `Tracts.ts`, `TractIcons.tsx`, `ToolboxPage.tsx`, `useDialog`, `ManageKeyDialog.tsx`, `tractCanvasLayout.ts`, `VaultItem`, `NotesSidebar.tsx`, `index.ts`, `CreateNoteDialog.tsx`, `SchemaProperty`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `tractSteps.ts`, `Topbar.tsx`, `InviteLinksSection.tsx`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `scripts`, `AuthAPI`, `ResultView.tsx`, `MembersSection.tsx`, `GoogleSheetsSpreadsheetSection.tsx`, `AuthMiddleware`, `AdminCouchAPI`, `.getRun`, `requiredConnections.ts`, `KebabMenu.tsx`, `useTheme.ts`, `ParamsList.tsx`, `ImportZipDialog.tsx`?**
  _High betweenness centrality (0.103) - this node is a cross-community bridge._
- **Why does `useDialog` connect `ToolboxPage.tsx` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `TractIcons.tsx`, `compilerOptions`, `Router.tsx`, `TractStepTree.tsx`, `BreadcrumbBar.tsx`, `CreateNoteDialog.tsx`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `dependencies`, `AuthMiddleware`, `tractSteps.ts`, `User.ts`, `toTract`, `NotesPage.tsx`, `VaultCard.tsx`, `TractCanvasTopBar.tsx`, `AuthAPI`, `connectionLabel`, `ResultView.tsx`, `LinkScreen.tsx`, `CardMeta.tsx`, `AdminCouchAPI`, `package.json`, `.getRun`, `AuthFetchInterceptor.ts`, `Tabs.tsx`, `ImportZipDialog.tsx`?**
  _High betweenness centrality (0.100) - this node is a cross-community bridge._
- **Why does `useUser` connect `Dialog.ts` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `TractIcons.tsx`, `ToolboxPage.tsx`, `grpcErrors.ts`, `VaultItem`, `McpKeysAPI`, `BreadcrumbBar.tsx`, `Notes.ts`, `index.ts`, `useVaultMutations`, `McpAuthPage.tsx`, `MobileNotesShell.tsx`, `ProviderIcon.tsx`, `Topbar.tsx`, `User.ts`, `toTract`, `VaultCard.tsx`, `S3InstanceFormDialog.tsx`, `McpKeys.ts`, `LinkScreen.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `CardMeta.tsx`, `package.json`, `.getRun`, `Vaults.ts`, `requiredConnections.ts`, `ConnectedContent.tsx`, `KebabMenu.tsx`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _840 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.02247191011235955 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08181818181818182 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.028169014084507043 - nodes in this community are weakly interconnected._