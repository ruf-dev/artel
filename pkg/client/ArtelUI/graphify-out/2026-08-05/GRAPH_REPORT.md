# Graph Report - ArtelUI  (2026-08-05)

## Corpus Check
- 589 files · ~153,090 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2709 nodes · 6613 edges · 126 communities (119 shown, 7 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 23 edges (avg confidence: 0.67)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e484ae25`
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

## God Nodes (most connected - your core abstractions)
1. `useDialog` - 202 edges
2. `cn` - 162 edges
3. `useBakeError()` - 152 edges
4. `useUser` - 101 edges
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
- `AccountsSectionProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/AccountsSection.tsx → src/app/api/artel/external_connections.pb.ts
- `EmailConnectionRowProps` --references--> `ExternalConnectionInfo`  [EXTRACTED]
  src/dialogs/ManageEmailDialog/components/AccountsSection/components/EmailConnectionRow/EmailConnectionRow.tsx → src/app/api/artel/external_connections.pb.ts

## Import Cycles
- 1-file cycle: `src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx -> src/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/setup-wizard/SetupWizardPage.tsx -> src/pages/setup-wizard/screens/CreateAdminScreen.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/init/InitPage.tsx -> src/pages/init/components/LoginContent/LoginContent.tsx -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/notes/NotesPage.tsx -> src/pages/notes/processes/notesUrl.ts -> src/app/routing/Router.tsx`
- 3-file cycle: `src/app/routing/Router.tsx -> src/pages/workbench/WorkbenchPage.tsx -> src/pages/workbench/components/PickAuthModeScreen/PickAuthModeScreen.tsx -> src/app/routing/Router.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarBrand/TopbarBrand.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarMobileDrawer/TopbarMobileDrawer.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarNav/TopbarNav.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 4-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarUserMenu/TopbarUserMenu.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 5-file cycle: `src/app/routing/HomeLayout.tsx -> src/segments/Topbar/Topbar.tsx -> src/segments/Topbar/components/TopbarMobileDrawer/TopbarMobileDrawer.tsx -> src/segments/Topbar/components/TopbarBrand/TopbarBrand.tsx -> src/app/routing/Router.tsx -> src/app/routing/HomeLayout.tsx`
- 5-file cycle: `src/app/routing/Router.tsx -> src/pages/tract-templates/TractTemplatesListPage.tsx -> src/pages/tract-templates/segments/ContentSegment/ContentSegment.tsx -> src/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx -> src/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx -> src/app/routing/Router.tsx`

## Communities (126 total, 7 thin omitted)

### Community 0 - "tracts.pb.ts"
Cohesion: 0.02
Nodes (91): Absent, ActionStep, BaseTractStep, ConditionStep, CreateTract, CreateTractRequest, CreateTractResponse, CreateTrigger (+83 more)

### Community 1 - "TaskTrackersPage.tsx"
Cohesion: 0.08
Nodes (28): AddTaskTracker, AddTaskTrackerRequest, AddTaskTrackerResponse, DeleteTaskTracker, DeleteTaskTrackerRequest, DeleteTaskTrackerResponse, ListTaskTrackers, ListTaskTrackersRequest (+20 more)

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
Cohesion: 0.04
Nodes (45): Absent, AddMcpConnector, AddMcpConnectorRequest, AddMcpConnectorResponse, BaseMcpToolInfo, BaseToolParamDef, CreateMcpKey, CreateMcpKeyRequest (+37 more)

### Community 6 - "TemplateInput.tsx"
Cohesion: 0.11
Nodes (28): Props, buildFlatList(), FilteredGroup, filterSources(), FlatEntry, FlatListItem, flattenArrayItems(), flattenProperty() (+20 more)

### Community 7 - "addTriggerDialogContext.ts"
Cohesion: 0.18
Nodes (18): useTrigger(), useTriggerSources(), categoryLabel(), PROVIDER_ENUM_BY_KEY, providerEnumFor(), providerLabel(), triggerChipLabel(), TriggerRow() (+10 more)

### Community 8 - "ExternalConnectionInfo"
Cohesion: 0.07
Nodes (5): CheckEmailConnectionRequest, ExternalConnectionsAPI, CheckStatus, EmailCheckButton(), EmailCheckButtonProps

### Community 9 - "admin_couch.pb.ts"
Cohesion: 0.09
Nodes (21): ChangeCouchUserPassword, ChangeCouchUserPasswordRequest, ChangeCouchUserPasswordResponse, DeleteCouchUser, DeleteCouchUserRequest, DeleteCouchUserResponse, GetUserDatabaseAccess, GetUserDatabaseAccessRequest (+13 more)

### Community 10 - "TractCanvasInspectorBody.tsx"
Cohesion: 0.12
Nodes (41): MomCandidate, ScriptLanguage, ActionCard(), Props, buildSourcesFor(), TractStepTreeProps, ActionBody(), Props (+33 more)

### Community 11 - "couch_instances.pb.ts"
Cohesion: 0.09
Nodes (22): DeleteCouchInstance, DeleteCouchInstanceRequest, DeleteCouchInstanceResponse, GetCouchInstance, GetCouchInstanceRequest, GetCouchInstanceStatus, GetCouchInstanceStatusRequest, GetCouchInstanceStatusResponse (+14 more)

### Community 12 - "TractsService"
Cohesion: 0.07
Nodes (6): TractsAPI, TractsService, toTract(), toTractTemplate(), toTractTemplateSummary(), definitionToProto()

### Community 13 - "Dialog.ts"
Cohesion: 0.13
Nodes (25): useMcpKeys, useServerStatus(), Path, Router(), routes, HeroSegment(), HeroSegmentProps, VaultCardHeader() (+17 more)

### Community 14 - "Tracts.ts"
Cohesion: 0.18
Nodes (21): NoteItem, PlusIcon(), UploadIcon(), FolderNodeItem(), FolderNodeItemProps, FolderSection(), FolderSectionProps, SearchResultsList() (+13 more)

### Community 15 - "cn"
Cohesion: 0.05
Nodes (44): cn, KebabMenu(), KebabMenuItem, Props, TODO: chures has no action/context-menu primitive yet, drop this wrapper once it, TODO: chures has no tab primitive yet, drop this wrapper once it does, Tabs(), TabsProps (+36 more)

### Community 16 - "TractIcons.tsx"
Cohesion: 0.14
Nodes (11): CreateInviteLinkDialog(), InviteLinksSection(), Props, ArtelUserDetailDialog(), ArtelUserDetailDialogProps, UserSessionsDialog(), UserSessionsDialogProps, ChangeArtelPasswordDialog() (+3 more)

### Community 17 - "compilerOptions"
Cohesion: 0.10
Nodes (24): SetTractsState, TractsState, triggerSourcesQueryKey, triggersQueryKey, sleep(), formatStartedAt(), Props, RunTractDialog() (+16 more)

### Community 18 - "ToolboxPage.tsx"
Cohesion: 0.18
Nodes (13): connectionLabel(), SelectOption(), ConnectionStep(), Props, LlmConnectionStep(), Props, ConnectionOptionList(), ConnectionOptionListProps (+5 more)

### Community 19 - "ConnectorChip.tsx"
Cohesion: 0.08
Nodes (24): compilerOptions, allowImportingTsExtensions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx (+16 more)

### Community 20 - "grpcErrors.ts"
Cohesion: 0.08
Nodes (31): CheckGitlabConnectionRequest, CheckTrelloConnectionRequest, UserErrors, CheckStatus, GitlabCheckButton(), GitlabCheckButtonProps, CheckStatus, TrelloCheckButton() (+23 more)

### Community 21 - "useDialog"
Cohesion: 0.11
Nodes (19): formatPrimitive(), JsonNode(), primitiveKind(), Props, tokenClass(), TokenKind, isNonEmptyBranch(), JsonBlock() (+11 more)

### Community 22 - "ManageKeyDialog.tsx"
Cohesion: 0.11
Nodes (18): Absent, BaseCompleteSetupRequest, CompleteSetup, CompleteSetupRequest, CompleteSetupResponse, GetStatus, GetStatusRequest, GetStatusResponse (+10 more)

### Community 23 - "tractCanvasLayout.ts"
Cohesion: 0.11
Nodes (26): NodeChips(), ConnectorPath(), ConnectorPathProps, ParallelBoxes(), Props, Props, TractCanvasArea(), useTractCanvasDrag() (+18 more)

### Community 24 - "auth.pb.ts"
Cohesion: 0.08
Nodes (24): Absent, BaseLoginRequest, GetConfig, GetConfigRequest, GetConfigResponse, GetMe, GetMeRequest, GetMeResponse (+16 more)

### Community 25 - "VaultItem"
Cohesion: 0.16
Nodes (12): VaultItem, Props, Props, VaultChip(), ExpertSettingsDrawer(), Props, Props, ExpertSettingsSection() (+4 more)

### Community 26 - "devDependencies"
Cohesion: 0.08
Nodes (26): devDependencies, eslint, eslint-import-resolver-typescript, @eslint/js, eslint-plugin-import-x, eslint-plugin-react, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+18 more)

### Community 27 - "fetch.pb.ts"
Cohesion: 0.13
Nodes (17): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+9 more)

### Community 28 - "McpKeysAPI"
Cohesion: 0.11
Nodes (10): CommunityConnectorInfo, CreateMcpKeyResponse, McpConnectorInfo, McpKeyInfo, McpKeysAPI, McpKeysState, ManageStep, MainScreenProps (+2 more)

### Community 29 - "Router.tsx"
Cohesion: 0.10
Nodes (17): useDialogKeyboard(), Props, PublishTemplateDialog(), SuggestionList(), SuggestionListProps, CreateNoteDialog(), Props, ArtelLogoIcon() (+9 more)

### Community 30 - "TractStepTree.tsx"
Cohesion: 0.12
Nodes (17): CandidateOptionList(), CandidateOptionListProps, MomCandidateCard(), MomCandidateCardProps, DialogHead(), DialogHeadProps, VaultField(), VaultFieldProps (+9 more)

### Community 31 - "BreadcrumbBar.tsx"
Cohesion: 0.20
Nodes (13): AddAnthropicConnectionRequest, AddEmailConnectionRequest, AddGenericConnectionRequest, AddGitlabConnectionRequest, AddTrelloConnectionRequest, CheckAnthropicConnectionRequest, CheckAnthropicConnectionResponse, ExternalConnectionInfo (+5 more)

### Community 32 - "NotesSidebar.tsx"
Cohesion: 0.14
Nodes (12): CloseIcon(), SearchIcon(), SearchIconProps, ListIcon(), TreeIcon(), NotesSearchBar(), NotesSearchBarProps, PublicBadge() (+4 more)

### Community 33 - "useExternalConnections"
Cohesion: 0.06
Nodes (34): ParamInput(), Props, ParamRow(), Props, ParamsList(), Props, Props, SchemaFieldRow() (+26 more)

### Community 34 - "ConnectionDetailDialog.tsx"
Cohesion: 0.13
Nodes (16): useNotes, DropZone(), Props, TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it d, ImportZipDialog(), Props, MobileNotesShell(), NotesSidebar() (+8 more)

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
Nodes (22): DeleteS3Instance, DeleteS3InstanceRequest, DeleteS3InstanceResponse, GetS3Instance, GetS3InstanceRequest, GetS3InstanceResponse, ListS3Instances, ListS3InstancesRequest (+14 more)

### Community 39 - "useVaultMutations"
Cohesion: 0.10
Nodes (4): VaultsAPI, useVaultMutations(), IWorkbenchService, WorkbenchService

### Community 40 - "CreateNoteDialog.tsx"
Cohesion: 0.08
Nodes (20): TractItem, TractRunItem, TractRunStepItem, TractTemplateItem, TractTemplateSummary, TractToolItem, TriggerItem, TriggerSourceItem (+12 more)

### Community 41 - "ArtelUI Frontend Rules"
Cohesion: 0.15
Nodes (14): MAIL_DOMAIN_ICONS, mailProviderIcon(), ConnectorChip(), EmailChip(), KNOWN_MAIL_DOMAIN_CLASSES, mailDomainAccent(), GenericChip(), PROVIDER_CHIP_CLASS (+6 more)

### Community 42 - "SchemaProperty"
Cohesion: 0.50
Nodes (6): parseScopeList(), SCOPE_INFO, trimScope(), buildScopeTooltipHtml(), GoogleSheetsInfoRows(), GoogleConnectionContent()

### Community 43 - "NoteEditor.tsx"
Cohesion: 0.16
Nodes (11): BoldIcon(), CodeIcon(), HeadingIcon(), ItalicIcon(), LinkIcon(), LineNumbers(), LineNumbersProps, NoteEditor() (+3 more)

### Community 44 - "McpAuthPage.tsx"
Cohesion: 0.09
Nodes (22): Input(), Props, RELATION_LABEL, RELATION_OPTIONS, RoadmapLinkTarget, WritableRelation, CredentialRow(), CredentialRowProps (+14 more)

### Community 45 - "MobileNotesShell.tsx"
Cohesion: 0.20
Nodes (9): GetSettings, GetSettingsRequest, GetSettingsResponse, UpdateAuthMethods, UpdateAuthMethodsRequest, UpdateAuthMethodsResponse, UpdateRegistrationMode, UpdateRegistrationModeRequest (+1 more)

### Community 46 - "Tracts.ts"
Cohesion: 0.10
Nodes (21): dependencies, classnames, @codemirror/autocomplete, @codemirror/lang-javascript, @codemirror/language, @codemirror/state, @codemirror/view, framer-motion (+13 more)

### Community 47 - "useErrorToast.ts"
Cohesion: 0.17
Nodes (10): FormField(), Props, Props, S3ToggleFields(), S3InstanceFormDialog(), S3InstancesActionBar(), InstanceFormDialog(), CreateVaultDialog() (+2 more)

### Community 48 - "dependencies"
Cohesion: 0.16
Nodes (18): LogicCell(), LogicCellProps, LogicSection(), Props, OptionCell(), ToolCell(), ToolCellProps, LogicOption (+10 more)

### Community 49 - "ProviderIcon.tsx"
Cohesion: 0.20
Nodes (9): LOGIC_OPTIONS, rank(), useTractBlockPickerData(), BranchIcon(), CodeIcon(), ForkIcon(), LayersIcon(), NodeIcon() (+1 more)

### Community 50 - "StepRow.tsx"
Cohesion: 0.43
Nodes (5): ResultViewMode, ViewModeToggle(), getResultViewWidget(), ResultView(), tryParseJson()

### Community 51 - "AuthMiddleware"
Cohesion: 0.12
Nodes (17): apiPrefix(), clearCsrfCookie(), csrfHeader(), getCsrfToken(), InitReq, TelegramLoginResponse, AppConfigState, useAppConfig (+9 more)

### Community 52 - "tractSteps.ts"
Cohesion: 0.20
Nodes (23): appendStep(), branchArray(), buildStepFromDraft(), collapseThinParallels(), collectAllStepIds(), generateStepId(), insertBlockAfter(), insertStepAt() (+15 more)

### Community 53 - "admin_users.pb.ts"
Cohesion: 0.06
Nodes (40): AdminUsersAPI, ArtelUserDetails, ArtelUserEntry, ChangeArtelUserPassword, ChangeArtelUserPasswordRequest, ChangeArtelUserPasswordResponse, CreateArtelUser, CreateArtelUserRequest (+32 more)

### Community 54 - "Topbar.tsx"
Cohesion: 0.15
Nodes (17): ConnectionsIcon(), base, NavIconProps, LogoutIcon(), NotesIcon(), ToolboxIcon(), TractsIcon(), VaultsIcon() (+9 more)

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
Nodes (14): DeleteDockerHost, DeleteDockerHostRequest, DeleteDockerHostResponse, GetDockerHost, GetDockerHostRequest, ListDockerHosts, ListDockerHostsRequest, ListDockerHostsResponse (+6 more)

### Community 60 - "TractBlockPicker.tsx"
Cohesion: 0.09
Nodes (20): ExternalProvider, ProviderChip(), AnthropicIcon(), EmailIcon(), GitlabIcon(), GoogleSheetsIcon(), MiroIcon(), TrelloIcon() (+12 more)

### Community 61 - "toTract"
Cohesion: 0.36
Nodes (4): VaultListProps, VaultSelect(), VaultSelectProps, Vault

### Community 62 - "NotesPage.tsx"
Cohesion: 0.15
Nodes (12): CardHeader(), Props, InsertRow(), Props, LlmCallCard(), Props, collectIdsFromRoot(), ConditionCard() (+4 more)

### Community 63 - "VaultCard.tsx"
Cohesion: 0.09
Nodes (17): AuthAPI, UserState, PasswordLoginForm(), PasswordLoginFormProps, RegisterForm(), RegisterFormProps, LoginContentProps, Mode (+9 more)

### Community 64 - "TractCanvasTopBar.tsx"
Cohesion: 0.21
Nodes (5): LogPanelBarProps, CollapseIcon(), ExpandIcon(), base, IconProps

### Community 65 - "TractCanvasLogPanel.tsx"
Cohesion: 0.36
Nodes (6): LocateIcon(), LocateIconProps, getNoteMeta(), getNoteTitle(), MobileTopBar(), MobileTopBarProps

### Community 66 - "scripts"
Cohesion: 0.27
Nodes (9): LogPanelBar(), ResizeHandle(), ResizeHandleProps, clampHeight(), dotClass(), formatDate(), loadStoredHeight(), Props (+1 more)

### Community 67 - "AuthAPI"
Cohesion: 0.22
Nodes (15): ImportConflictAction, commitImportAndRefresh(), deleteFolderAndRefresh(), moveEntryAndRefresh(), NotesState, remapSelectedPath(), requireVaultId(), ConflictRow() (+7 more)

### Community 68 - "connectionLabel"
Cohesion: 0.10
Nodes (27): DialogManager, useDialog, useTracts, Props, ToolStep(), StepPickerDialog(), Props, TriggerPanel() (+19 more)

### Community 69 - "S3InstanceFormDialog.tsx"
Cohesion: 0.10
Nodes (13): ComingSoonCardProps, COMING_SOON_CARDS, LLM_BYOK_PROVIDERS, ByokTabIcon(), CommunitySection(), CommunityTabIcon(), ExternalConnectionsTabIcon(), ConnectionsPage() (+5 more)

### Community 70 - "compilerOptions"
Cohesion: 0.18
Nodes (4): CouchInstancesAPI, InstancesActionBar(), InstancesActionBarProps, InstancesTab()

### Community 71 - "McpKeys.ts"
Cohesion: 0.27
Nodes (7): GetCouchInstanceResponse, InstanceRow(), InstanceRowProps, InstanceFormDialogProps, InstanceListProps, InstanceSelector(), InstanceSelectorProps

### Community 72 - "ResultView.tsx"
Cohesion: 0.08
Nodes (34): JSON_KIND_CLASS, JsonToken, JsonTokenKind, JsonView(), tokenizeJson(), AddTriggerDialog(), STEP_SCREENS, AddTriggerDialogContext (+26 more)

### Community 73 - "MembersSection.tsx"
Cohesion: 0.09
Nodes (35): Props, RELATION_CLASS, RoadmapConnectorPath(), Props, RoadmapCanvasArea(), boardListLabel(), Props, RoadmapCanvasNode() (+27 more)

### Community 74 - "LinkScreen.tsx"
Cohesion: 0.22
Nodes (9): BinaryStorageToggle(), Props, Props, PublishSlugForm(), slugify(), validateSlug(), Props, PublishToggle() (+1 more)

### Community 75 - "dialog-scrollable.js"
Cohesion: 0.46
Nodes (7): allRules(), dialogScrollable(), directDeclsOf(), findScrollTarget(), isOverflowY(), messages, meta

### Community 76 - "Handoff: lint/tooling parity gaps vs. ZpotifyUI"
Cohesion: 0.13
Nodes (14): VaultInviteItem, VaultMemberInfo, vaultsQueryKey, DangerZoneText(), Props, InviteRow(), Props, Props (+6 more)

### Community 77 - "RunTractDialog.tsx"
Cohesion: 0.23
Nodes (8): formatRelative(), Props, RunStatusDot(), ChevronRightIcon(), ManualTriggerIcon(), TrashIcon(), Props, TractLastRun

### Community 78 - "DbAccessList.tsx"
Cohesion: 0.16
Nodes (12): useIsMobileNav(), applyTheme(), Theme, useTheme(), HomeLayout(), BrandMarkIcon(), TopbarBrand(), TopbarMobileDrawer() (+4 more)

### Community 79 - "VaultDangerZone.tsx"
Cohesion: 0.20
Nodes (6): AdminSystemSettingsAPI, AuthMethodsSection(), AuthMethodsSectionProps, RegistrationModeSection(), RegistrationModeSectionProps, SettingsTab()

### Community 80 - "CardMeta.tsx"
Cohesion: 0.15
Nodes (21): CreateTriggerRequest, TractTemplatesState, useTractTemplates, PresetDetails(), SourcePicker(), BrowseTemplatesDialog(), Props, TemplateRow() (+13 more)

### Community 92 - "GoogleSheetsSpreadsheetSection.tsx"
Cohesion: 0.35
Nodes (7): AdminPage(), resolveTab(), Tab, AdminHero(), AdminHeroProps, TabBar(), TabBarProps

### Community 93 - "VaultCardHeader.tsx"
Cohesion: 0.28
Nodes (3): Spreadsheet, GoogleSheetsSpreadsheetSection(), SpreadsheetRow()

### Community 94 - "AuthMiddleware"
Cohesion: 0.25
Nodes (6): CopyIcon(), ArrowIcon(), ArrowIconProps, FileIcon(), FolderIcon(), TreeItemProps

### Community 95 - "ConnectForm.tsx"
Cohesion: 0.08
Nodes (25): dompurify, NoteMode, BreadcrumbBarProps, Mode, BreadcrumbPath(), BreadcrumbPathProps, DesktopNotesShellProps, VaultOption (+17 more)

### Community 96 - "UsersTab.tsx"
Cohesion: 0.36
Nodes (4): GenericToolIcon(), TODO: placeholder glyph for tool actions without a dedicated icon yet (non-smtp/, ImapIcon(), SmtpIcon()

### Community 97 - "AdminCouchAPI"
Cohesion: 0.24
Nodes (6): Props, RunButton(), Props, RunStatusBadge(), Props, PlayIcon()

### Community 98 - "CouchInstancesAPI"
Cohesion: 0.22
Nodes (8): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, skipLibCheck, strict, include

### Community 99 - "TractCanvasLogPanel.tsx"
Cohesion: 0.24
Nodes (7): CardChips(), CardHeader(), Props, CardMeta(), formatDate(), Props, Props

### Community 100 - "RunLog.tsx"
Cohesion: 0.15
Nodes (11): ImapOperation, ImapToolAction, McpToolInfo, ToolParamDef, ImapActionView(), ParamRow(), ParamsList(), RunScreens() (+3 more)

### Community 101 - "package.json"
Cohesion: 0.14
Nodes (11): GetDockerHostResponse, DockerHostFormDialog(), DockerHostFormDialogProps, DockerHostFormDialogSaveData, DockerHostListProps, DockerHostRow(), DockerHostRowProps, DockerHostsActionBar() (+3 more)

### Community 102 - "MobileDrawer.tsx"
Cohesion: 0.21
Nodes (9): ConnectionSection(), Props, useTemplateConnections(), UseTemplateConnectionsResult, InstantiateTemplateDialog(), Props, ConnectionRequirement, ConnectionRequirementKind (+1 more)

### Community 103 - "CouchInstancesAPI"
Cohesion: 0.33
Nodes (5): Handoff: lint/tooling parity gaps vs. ZpotifyUI, No CSS linter at all, No Prettier, Structural/style ESLint rules Zpotify enforces that Artel only documents, Suggested order of attack

### Community 104 - ".getRun"
Cohesion: 0.11
Nodes (29): useExternalConnections, useBakeError(), ConnectGenericDialog(), CredentialField, ConnectionDetailDialog(), PROVIDER_CONFIG, PROVIDER_KEY, ProviderConfig (+21 more)

### Community 105 - "UserList.tsx"
Cohesion: 0.18
Nodes (11): scripts, build, build:ui, dev, gen, lint, lint:css, lint:js (+3 more)

### Community 106 - "Vaults.ts"
Cohesion: 0.22
Nodes (3): AdminCouchAPI, ManageAccessDialog(), UsersTab()

### Community 107 - "AuthFetchInterceptor.ts"
Cohesion: 0.40
Nodes (3): ConnectorRow(), ConnectorRowProps, ConnectionsFieldProps

### Community 108 - "EmailCheckButton.tsx"
Cohesion: 0.32
Nodes (5): DbAccessList(), DbAccessListProps, DbAccessRow(), DbAccessRowProps, ManageAccessDialogProps

### Community 109 - "GitlabCheckButton.tsx"
Cohesion: 0.33
Nodes (4): ArtelAPI, Version, VersionRequest, VersionResponse

### Community 110 - "TrelloCheckButton.tsx"
Cohesion: 0.48
Nodes (4): CouchUserEntry, UserListProps, UserRow(), UserRowProps

### Community 111 - "TopbarDrawerCloseButton.tsx"
Cohesion: 0.27
Nodes (6): Props, VaultCardConnBar(), Props, VaultCardFront(), VaultCardStatus(), Props

### Community 112 - "InsertConflictDialog.tsx"
Cohesion: 0.36
Nodes (6): OptIcon(), OptIconProps, OptText(), OptTextProps, OptionCellProps, StepColor

### Community 113 - "TopbarMobileTrigger.tsx"
Cohesion: 0.40
Nodes (4): AccountsSection(), AccountsSectionProps, TrelloConnectionRow(), TrelloConnectionRowProps

### Community 117 - "requiredConnections.ts"
Cohesion: 0.18
Nodes (13): useVaults(), useWorkbench(), useWorkbenchMutations(), workbenchQueryKey(), VaultChipDisplayProps, VaultOptionList(), VaultOptionListProps, ContentSegment() (+5 more)

### Community 118 - "package.json"
Cohesion: 0.33
Nodes (5): name, private, trustedDependencies, type, version

### Community 120 - "KebabMenu.tsx"
Cohesion: 0.18
Nodes (12): RegistrationMode, SetupWizardAPI, AuthMethodsScreen(), OPTIONS, RegistrationModeScreen(), TokenEntryScreen(), SetupWizardContext, SetupWizardState (+4 more)

### Community 121 - "useTheme.ts"
Cohesion: 0.53
Nodes (5): buildLogLines(), formatDuration(), formatTime(), LogLineKind, stepMeta()

### Community 123 - "Tabs.tsx"
Cohesion: 0.33
Nodes (8): usePortrait(), useAutosave(), NotesPage(), buildNotesUrl(), decodeNotePath(), encodeNotePath(), buildRenameHandler(), useReadOnlyVaultGate()

## Knowledge Gaps
- **796 isolated node(s):** `localPlugin`, `name`, `private`, `version`, `type` (+791 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `useDialog` connect `connectionLabel` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `TractIcons.tsx`, `compilerOptions`, `ToolboxPage.tsx`, `grpcErrors.ts`, `VaultItem`, `McpKeysAPI`, `Router.tsx`, `TractStepTree.tsx`, `ConnectionDetailDialog.tsx`, `s3_instances.pb.ts`, `ArtelUI Frontend Rules`, `McpAuthPage.tsx`, `useErrorToast.ts`, `AuthMiddleware`, `admin_users.pb.ts`, `User.ts`, `TractBlockPicker.tsx`, `toTract`, `NotesPage.tsx`, `VaultCard.tsx`, `AuthAPI`, `S3InstanceFormDialog.tsx`, `compilerOptions`, `McpKeys.ts`, `LinkScreen.tsx`, `RunTractDialog.tsx`, `CardMeta.tsx`, `package.json`, `MobileDrawer.tsx`, `.getRun`, `Vaults.ts`, `EmailCheckButton.tsx`, `TrelloCheckButton.tsx`, `Tabs.tsx`?**
  _High betweenness centrality (0.115) - this node is a cross-community bridge._
- **Why does `cn` connect `cn` to `useBakeError`, `TemplateInput.tsx`, `addTriggerDialogContext.ts`, `TractCanvasInspectorBody.tsx`, `Dialog.ts`, `Tracts.ts`, `ToolboxPage.tsx`, `useDialog`, `tractCanvasLayout.ts`, `VaultItem`, `Router.tsx`, `NotesSidebar.tsx`, `ConnectionDetailDialog.tsx`, `index.ts`, `ArtelUI Frontend Rules`, `SchemaProperty`, `McpAuthPage.tsx`, `StepRow.tsx`, `admin_users.pb.ts`, `Topbar.tsx`, `TractBlockPicker.tsx`, `NotesPage.tsx`, `scripts`, `AuthAPI`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `ResultView.tsx`, `MembersSection.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `RunTractDialog.tsx`, `DbAccessList.tsx`, `VaultDangerZone.tsx`, `GoogleSheetsSpreadsheetSection.tsx`, `AuthMiddleware`, `ConnectForm.tsx`, `AdminCouchAPI`, `RunLog.tsx`, `KebabMenu.tsx`?**
  _High betweenness centrality (0.086) - this node is a cross-community bridge._
- **Why does `useUser` connect `Dialog.ts` to `TaskTrackersPage.tsx`, `addTriggerDialogContext.ts`, `ExternalConnectionInfo`, `couch_instances.pb.ts`, `TractIcons.tsx`, `compilerOptions`, `grpcErrors.ts`, `VaultItem`, `McpKeysAPI`, `BreadcrumbBar.tsx`, `Notes.ts`, `index.ts`, `s3_instances.pb.ts`, `useVaultMutations`, `McpAuthPage.tsx`, `useErrorToast.ts`, `AuthMiddleware`, `admin_users.pb.ts`, `Topbar.tsx`, `User.ts`, `toTract`, `VaultCard.tsx`, `connectionLabel`, `S3InstanceFormDialog.tsx`, `compilerOptions`, `McpKeys.ts`, `LinkScreen.tsx`, `Handoff: lint/tooling parity gaps vs. ZpotifyUI`, `DbAccessList.tsx`, `VaultDangerZone.tsx`, `CardMeta.tsx`, `package.json`, `Vaults.ts`, `EmailCheckButton.tsx`, `TrelloCheckButton.tsx`, `requiredConnections.ts`, `KebabMenu.tsx`?**
  _High betweenness centrality (0.076) - this node is a cross-community bridge._
- **What connects `localPlugin`, `name`, `private` to the rest of the system?**
  _802 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `tracts.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.02407667134174848 - nodes in this community are weakly interconnected._
- **Should `TaskTrackersPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07922705314009662 - nodes in this community are weakly interconnected._
- **Should `vaults.pb.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.03125 - nodes in this community are weakly interconnected._