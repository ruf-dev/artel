# Graph Report - ArtelUI  (2026-08-05)

## Corpus Check
- 577 files · ~150,626 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2670 nodes · 6448 edges · 124 communities (117 shown, 7 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 23 edges (avg confidence: 0.67)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f4a3da35`
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
- NotConnectedContent.tsx
- ConnectedContent.tsx
- KebabMenu.tsx
- WorkbenchStatusBadge.tsx
- TopbarMobileTrigger.tsx
- Tabs.tsx

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 200 edges
2. `cn` - 156 edges
3. `useBakeError()` - 136 edges
4. `useUser` - 100 edges
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
- `MomCandidateCardProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/components/CandidateOptionList/components/MomCandidateCard/MomCandidateCard.tsx → src/app/api/artel/mcp_keys.pb.ts
- `ConnectionOptionListProps` --references--> `MomCandidate`  [EXTRACTED]
  src/dialogs/ManageKeyDialog/components/ConnectionOptionList/ConnectionOptionList.tsx → src/app/api/artel/mcp_keys.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/workbench/WorkbenchPage.tsx -> src/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx -> src/app/routing/Router.tsx`
- 5-file cycle: `src/app/routing/Router.tsx -> src/pages/tract-templates/TractTemplatesListPage.tsx -> src/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx -> src/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx -> src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx -> src/app/routing/Router.tsx`

## Communities (124 total, 7 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.02
Nodes (88): Absent, BaseTractStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger, CreateTriggerRequest, CreateTriggerResponse (+80 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.06
Nodes (44): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+36 more)

### Community 2 - "vaults.pb.ts"
Cohesion: 0.03
Nodes (63): AcceptInvite, AcceptInviteRequest, AcceptInviteResponse, AddMember, AddMemberRequest, AddMemberResponse, CreateInviteLink, CreateInviteLinkRequest (+55 more)

### Community 3 - "useBakeError"
Cohesion: 0.06
Nodes (38): GetVaultBySlug, GetVaultBySlugRequest, GetVaultBySlugResponse, PublicDocsAPI, PublicDocsGetNote, PublicDocsGetNoteRequest, PublicDocsGetNoteResponse, PublicDocsListFolders (+30 more)

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
Cohesion: 0.19
Nodes (17): useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor(), providerLabel(), triggerChipLabel(), TriggerRow() (+9 more)

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.14
Nodes (33): MomCandidate, Props, ToolStep(), Props, TractStepTreeProps, Props, Props, DangerZone() (+25 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.09
Nodes (22): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+14 more)

### Community 12 - "TractsService"
Cohesion: 0.06
Nodes (19): TractsAPI, safeParseJson(), PresetDetails(), TractsService, emptySchema, parseSchema(), toRun(), toRunStep() (+11 more)

### Community 13 - "Dialog.ts"
Cohesion: 0.12
Nodes (23): useServerStatus(), HomeLayout(), Path, Router(), routes, HeroSegment(), HeroSegmentProps, VaultCardHeader() (+15 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.10
Nodes (19): cn, ProviderChip(), CardHeader(), ConnBarRow(), RowProps, SchemaPropertyPreview(), Props, RoleOption() (+11 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.10
Nodes (22): useBakeError(), DetailScreen(), AnthropicCheckButton(), CheckStatus, CreateInviteLinkDialog(), InviteLinksSection(), Props, ArtelUserDetailDialog() (+14 more)

### Community 17 - "compilerOptions"
Cohesion: 0.10
Nodes (29): SetTractsState, TractsState, triggerSourcesQueryKey, triggersQueryKey, sleep(), Props, ParallelBoxes(), Props (+21 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.14
Nodes (12): ImapOperation, ImapToolAction, McpToolInfo, SmtpOperation, SmtpToolAction, ConnectionPicker(), ImapActionView(), RunScreens() (+4 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.08
Nodes (25): CheckGitlabConnectionRequest, CheckTrelloConnectionRequest, CheckStatus, GitlabCheckButton(), GitlabCheckButtonProps, CheckStatus, TrelloCheckButton(), TrelloCheckButtonProps (+17 more)

### Community 21 - "useDialog"
Cohesion: 0.20
Nodes (9): isNonEmptyBranch(), JsonBlock(), Props, Props, statusClass(), StepRow(), Props, Props (+1 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.08
Nodes (20): Absent, BaseCompleteSetupRequest, CompleteSetup, CompleteSetupRequest, CompleteSetupResponse, GetStatus, GetStatusRequest, GetStatusResponse (+12 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.21
Nodes (13): ConnectorPath(), ConnectorPathProps, Props, TractCanvasArea(), useTractCanvasDrag(), cap(), DragOverSide, NodeStatus (+5 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (24): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+16 more)

### Community 25 - "VaultItem"
Cohesion: 0.10
Nodes (20): VaultItem, CardChips(), Props, Props, VaultCardConnBar(), Props, VaultCardFront(), Props (+12 more)

### Community 26 - "devDependencies"
Cohesion: 0.08
Nodes (26): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+18 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.10
Nodes (8): CommunityConnectorInfo, CreateMcpKeyResponse, McpKeyInfo, McpKeysAPI, McpKeysState, IMcpKeysService, McpKeysService, Props

### Community 29 - "Router.tsx"
Cohesion: 0.17
Nodes (12): useWorkbench(), useWorkbenchMutations(), workbenchQueryKey(), KNOWN_STATES, LoginRelayScreen(), LoginState, Props, Stage (+4 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.16
Nodes (13): DialogHead(), DialogHeadProps, VaultField(), VaultFieldProps, ManageKeyDialog(), ManageStep, useManageKeyDialog(), MainScreen() (+5 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.11
Nodes (15): Mode, BreadcrumbPath(), BreadcrumbPathProps, CheckIcon(), CopyIcon(), ErrorDotIcon(), PencilIcon(), SpinnerIcon() (+7 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.14
Nodes (12): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, PublicBadge() (+4 more)

### Community 33 - "useExternalConnections"
Cohesion: 0.14
Nodes (18): paramCompletionSource(), scriptEditorTheme, scriptHighlightStyle, Props, ScriptCodeSection(), Props, addInputParam(), addOutputParam() (+10 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.09
Nodes (23): useNotes, SuggestionList(), SuggestionListProps, CreateNoteDialog(), Props, ArtelLogoIcon(), ChevronLeftIcon(), DrawerCloseButton() (+15 more)

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
Cohesion: 0.16
Nodes (15): ActionStep, ConditionStep, GroupStep, LlmCallStep, ParallelStep, ScriptStep, TractCondition, TractDefinition (+7 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.21
Nodes (9): ExternalConnectionInfo, AccountsSection(), AccountsSectionProps, EmailConnectionRow(), EmailConnectionRowProps, AccountsSection(), AccountsSectionProps, TrelloConnectionRow() (+1 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.50
Nodes (6): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleConnectionContent()

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.12
Nodes (16): CheckEmailConnectionRequest, Input(), FormField(), Props, Props, S3ToggleFields(), CredentialRowProps, buildEmailRequest() (+8 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.14
Nodes (10): AdminSystemSettingsAPI, GetSettings, GetSettingsRequest, GetSettingsResponse, UpdateAuthMethods, UpdateAuthMethodsRequest, UpdateAuthMethodsResponse, UpdateRegistrationMode (+2 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.22
Nodes (3): AdminCouchAPI, ManageAccessDialog(), UsersTab()

### Community 47 - "useErrorToast.ts"
Cohesion: 0.13
Nodes (14): DialogManager, useMcpKeys, CardHeader(), Props, SelectOption(), ConnectionStep(), Props, ContentSegment() (+6 more)

### Community 48 - "dependencies"
Cohesion: 0.10
Nodes (21): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, framer-motion (+13 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.24
Nodes (7): ByokTabIcon(), CommunityTabIcon(), ContentSegment(), ExternalConnectionsTabIcon(), ConnectionsPage(), ConnectionsTab, resolveTab()

### Community 51 - "AuthMiddleware"
Cohesion: 0.13
Nodes (13): apiPrefix(), csrfHeader(), getCsrfToken(), InitReq, TelegramLoginResponse, AppConfigState, useAppConfig, pingServer() (+5 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.15
Nodes (28): Props, Step, StepDraft, InsertConflictDialog(), Props, appendStep(), branchArray(), buildStepFromDraft() (+20 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.08
Nodes (25): AdminUsersAPI, ArtelUserDetails, ArtelUserEntry, ChangeArtelUserPassword, ChangeArtelUserPasswordRequest, ChangeArtelUserPasswordResponse, CreateArtelUser, CreateArtelUserRequest (+17 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.08
Nodes (30): useIsMobileNav(), applyTheme(), Theme, useTheme(), BrandMarkIcon(), ConnectionsIcon(), base, NavIconProps (+22 more)

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
Cohesion: 0.06
Nodes (30): DeleteDockerHost, DeleteDockerHostRequest, DeleteDockerHostResponse, DockerHostsAPI, GetDockerHost, GetDockerHostRequest, GetDockerHostResponse, ListDockerHosts (+22 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.10
Nodes (18): MAIL_DOMAIN_ICONS, mailProviderIcon(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), AnthropicIcon(), EmailIcon() (+10 more)

### Community 61 - "toTract"
Cohesion: 0.26
Nodes (6): McpLogin(), McpLoginProps, VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.08
Nodes (36): ScriptLanguage, connectionLabel(), TemplateSource, Props, ActionCard(), Props, CONDITION_OPS, ConditionRowProps (+28 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.26
Nodes (7): useVaults(), VaultChipDisplayProps, VaultOptionList(), VaultOptionListProps, ContentSegment(), ContentSegmentProps, HomePage()

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.05
Nodes (51): Props, RunButton(), Props, RunStatusBadge(), LogicCell(), LogicCellProps, LogicSection(), Props (+43 more)

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.24
Nodes (11): cap(), LogEntryRow(), Props, buildLogLines(), formatDuration(), formatTime(), LogLine, LogLineKind (+3 more)

### Community 67 - "AuthAPI"
Cohesion: 0.22
Nodes (15): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), moveEntryAndRefresh(), NotesState, remapSelectedPath(), requireVaultId(), ConflictRow() (+7 more)

### Community 68 - "connectionLabel"
Cohesion: 0.09
Nodes (32): useDialog, useDialogKeyboard(), useTracts, StepPickerDialog(), Props, TriggerPanel(), AddTaskLinkDialog(), Props (+24 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.17
Nodes (12): AddAnthropicConnectionRequest, AddEmailConnectionRequest, AddGenericConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, CheckAnthropicConnectionRequest, CheckAnthropicConnectionResponse, ExternalConnectionsState (+4 more)

### Community 70 - "compilerOptions"
Cohesion: 0.18
Nodes (4): CouchInstancesAPI, InstancesActionBar(), InstancesActionBarProps, InstancesTab()

### Community 71 - "McpKeys.ts"
Cohesion: 0.29
Nodes (6): GetCouchInstanceResponse, InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps

### Community 72 - "ResultView.tsx"
Cohesion: 0.08
Nodes (34): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), AddTriggerDialog(), STEP_SCREENS, AddTriggerDialogContext (+26 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.09
Nodes (38): RoadmapLinkTarget, Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), boardListLabel(), Props (+30 more)

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.18
Nodes (11): scripts, build, build:ui, dev, gen, lint, lint:css, lint:js (+3 more)

### Community 75 - "dialog-scrollable.js"
Cohesion: 0.46
Nodes (7): allRules(), dialogScrollable(), directDeclsOf(), findScrollTarget(), isOverflowY(), messages, meta

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.16
Nodes (13): VaultInviteItem, VaultMemberInfo, vaultsQueryKey, Props, InviteRow(), Props, Props, Props (+5 more)

### Community 77 - "RunTractDialog.tsx"
Cohesion: 0.10
Nodes (17): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, formatStartedAt(), RunTractDialog() (+9 more)

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.17
Nodes (11): BinaryStorageToggle(), Props, Props, PublishSlugForm(), slugify(), validateSlug(), Props, PublishToggle() (+3 more)

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.31
Nodes (8): ResizeHandle(), ResizeHandleProps, clampHeight(), dotClass(), formatDate(), loadStoredHeight(), Props, TractCanvasLogPanel()

### Community 80 - "CardMeta.tsx"
Cohesion: 0.23
Nodes (12): TractTemplatesState, useTractTemplates, BrowseTemplatesDialog(), Props, TemplateRow(), Props, ListScreen(), Props (+4 more)

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.35
Nodes (7): AdminPage(), resolveTab(), Tab, AdminHero(), AdminHeroProps, TabBar(), TabBarProps

### Community 93 - "VaultCardHeader.tsx"
Cohesion: 0.28
Nodes (3): Spreadsheet, GoogleSheetsSpreadsheetSection(), SpreadsheetRow()

### Community 94 - "AuthMiddleware"
Cohesion: 0.11
Nodes (12): AuthAPI, UserState, LoginContentProps, AuthService, IAuthService, Session, AuthMiddleware, clearLocalStorage() (+4 more)

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.13
Nodes (19): dompurify, NoteMode, usePortrait(), BreadcrumbBarProps, DesktopNotesShellProps, VaultOption, MobileNotesShellProps, NoteViewer() (+11 more)

### Community 96 - "UsersTab.tsx"
Cohesion: 0.36
Nodes (4): GenericToolIcon(), TODO: placeholder glyph for tool actions without a dedicated icon yet (non-smtp/, ImapIcon(), SmtpIcon()

### Community 97 - "AdminCouchAPI"
Cohesion: 0.20
Nodes (10): McpConnectorInfo, CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps (+2 more)

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.67
Nodes (3): CardMeta(), formatDate(), Props

### Community 100 - "RunLog.tsx"
Cohesion: 0.43
Nodes (3): ToolParamDef, ParamRow(), ParamsList()

### Community 101 - "package.json"
Cohesion: 0.33
Nodes (5): DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, ImportZipDialog(), Props

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.43
Nodes (5): useTemplateConnections(), UseTemplateConnectionsResult, InstantiateTemplateDialog(), Props, ConnectionRequirement

### Community 103 - "CouchInstancesAPI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 104 - ".getRun"
Cohesion: 0.06
Nodes (36): ExternalProvider, useExternalConnections, ComingSoonCardProps, LlmConnectionStep(), Props, CredentialRow(), ConnectGenericDialog(), CredentialField (+28 more)

### Community 105 - "UserList.tsx"
Cohesion: 0.32
Nodes (5): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps, ManageAccessDialogProps

### Community 106 - "Vaults.ts"
Cohesion: 0.43
Nodes (6): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind

### Community 107 - "AuthFetchInterceptor.ts"
Cohesion: 0.23
Nodes (7): ConnectedContent(), ConnectedContentProps, UseDefaultButton(), ConnectForm(), tokenSettingsUrl(), DialogHead(), ManageGitlabDialog()

### Community 108 - "EmailCheckButton.tsx"
Cohesion: 0.48
Nodes (4): CouchUserEntry, UserListProps, UserRow(), UserRowProps

### Community 109 - "GitlabCheckButton.tsx"
Cohesion: 0.21
Nodes (7): ArtelAPI, Version, VersionRequest, VersionResponse, UserErrors, getPreconditionViolations(), isMissingSubscription()

### Community 110 - "TrelloCheckButton.tsx"
Cohesion: 0.43
Nodes (5): ResultViewMode, ViewModeToggle(), getResultViewWidget(), ResultView(), tryParseJson()

### Community 111 - "TopbarDrawerCloseButton.tsx"
Cohesion: 0.13
Nodes (8): GetS3InstanceResponse, S3InstancesAPI, Props, S3InstanceFormDialog(), S3InstanceRow(), TestStatus, S3InstancesActionBar(), S3InstancesTab()

### Community 112 - "InsertConflictDialog.tsx"
Cohesion: 0.29
Nodes (5): TopbarUserMenuItem(), TopbarUserMenuItemProps, MenuRect, TopbarUserMenuList(), TopbarUserMenuListProps

### Community 113 - "TopbarMobileTrigger.tsx"
Cohesion: 0.60
Nodes (4): formatRelative(), Props, RunStatusDot(), TractLastRun

### Community 118 - "NotConnectedContent.tsx"
Cohesion: 0.33
Nodes (5): name, private, trustedDependencies, type, version

### Community 119 - "ConnectedContent.tsx"
Cohesion: 0.40
Nodes (3): Props, SearchInput(), ArtelUsersTab()

### Community 120 - "KebabMenu.tsx"
Cohesion: 0.40
Nodes (4): KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it

### Community 121 - "WorkbenchStatusBadge.tsx"
Cohesion: 0.40
Nodes (4): Props, STATUS_CLASSES, STATUS_LABELS, WorkbenchStatusBadge()

### Community 122 - "TopbarMobileTrigger.tsx"
Cohesion: 0.50
Nodes (3): TopbarHamburgerIcon(), TopbarMobileTrigger(), TopbarMobileTriggerProps

### Community 123 - "Tabs.tsx"
Cohesion: 0.50
Nodes (3): TODO: chures has no tab primitive yet, drop this wrapper once it does, Tabs(), TabsProps

## Knowledge Gaps
- **793 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+788 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `Tracts.ts`, `TractIcons.tsx`, `compilerOptions`, `ToolboxPage.tsx`, `useDialog`, `tractCanvasLayout.ts`, `VaultItem`, `BreadcrumbBar.tsx`, `NotesSidebar.tsx`, `ConnectionDetailDialog.tsx`, `index.ts`, `SchemaProperty`, `McpAuthPage.tsx`, `useErrorToast.ts`, `Topbar.tsx`, `InviteLinksSection.tsx`, `TractBlockPicker.tsx`, `NotesPage.tsx`, `TractCanvasTopBar.tsx`, `scripts`, `AuthAPI`, `connectionLabel`, `ResultView.tsx`, `MembersSection.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `VaultDangerZone.tsx`, `GoogleSheetsSpreadsheetSection.tsx`, `package.json`, `.getRun`, `Vaults.ts`, `TrelloCheckButton.tsx`, `InsertConflictDialog.tsx`, `TopbarMobileTrigger.tsx`, `ConnectedContent.tsx`, `KebabMenu.tsx`, `WorkbenchStatusBadge.tsx`, `TopbarMobileTrigger.tsx`, `Tabs.tsx`?**
  _High betweenness centrality (0.118) - this node is a cross-community bridge._
- **Why does `useDialog` connect `connectionLabel` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `TractIcons.tsx`, `compilerOptions`, `VaultItem`, `TractStepTree.tsx`, `ConnectionDetailDialog.tsx`, `McpAuthPage.tsx`, `Tracts.ts`, `useErrorToast.ts`, `StepRow.tsx`, `tractSteps.ts`, `admin_users.pb.ts`, `User.ts`, `InviteLinksSection.tsx`, `TractBlockPicker.tsx`, `toTract`, `NotesPage.tsx`, `VaultCard.tsx`, `TractCanvasTopBar.tsx`, `AuthAPI`, `compilerOptions`, `MembersSection.tsx`, `RunTractDialog.tsx`, `DbAccessList.tsx`, `CardMeta.tsx`, `AuthMiddleware`, `ConnectForm.tsx`, `package.json`, `MobileDrawer.tsx`, `.getRun`, `UserList.tsx`, `AuthFetchInterceptor.ts`, `EmailCheckButton.tsx`, `TopbarDrawerCloseButton.tsx`?**
  _High betweenness centrality (0.103) - this node is a cross-community bridge._
- **Why does `useUser` connect `Dialog.ts` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `couch_instances.pb.ts`, `TractsService`, `TractIcons.tsx`, `compilerOptions`, `grpcErrors.ts`, `VaultItem`, `McpKeysAPI`, `Router.tsx`, `Notes.ts`, `index.ts`, `useVaultMutations`, `McpAuthPage.tsx`, `Tracts.ts`, `useErrorToast.ts`, `StepRow.tsx`, `AuthMiddleware`, `admin_users.pb.ts`, `Topbar.tsx`, `User.ts`, `InviteLinksSection.tsx`, `toTract`, `VaultCard.tsx`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `compilerOptions`, `McpKeys.ts`, `MembersSection.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `DbAccessList.tsx`, `GoogleSheetsSpreadsheetSection.tsx`, `AuthMiddleware`, `UserList.tsx`, `EmailCheckButton.tsx`, `TopbarDrawerCloseButton.tsx`, `ConnectedContent.tsx`?**
  _High betweenness centrality (0.065) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _799 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.02247191011235955 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.058823529411764705 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.03125 - nodes in this community are weakly interconnected._