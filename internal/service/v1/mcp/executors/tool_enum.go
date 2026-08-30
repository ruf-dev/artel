package executors

// McpName identifies a MoM tool-set (one per MoM record, e.g. "trello"). It is the hand-written
// mirror of the Postgres `mcp_name` ENUM created by migration 080_mcp_tool_enums.sql;
// tests/e2e/mcp_enum_test.go asserts AllMcpNames and enum_range(NULL::mcp_name) are the same set.
type McpName string

// McpTool identifies a single MoM tool, formatted "<mcp>.<tool>" (e.g. "trello.create_card"). It
// is the hand-written mirror of the Postgres `mcp_tool` ENUM created by migration
// 080_mcp_tool_enums.sql (and extended by later migrations that seed new tools);
// tests/e2e/mcp_enum_test.go asserts AllMcpTools and enum_range(NULL::mcp_tool) are the same set.
type McpTool string

const (
	McpNameEmail    McpName = "email"
	McpNameGitlab   McpName = "gitlab"
	McpNameTelegram McpName = "telegram"
	McpNameTrello   McpName = "trello"
)

const (
	McpToolEmailListEmailFolders McpTool = "email.list_email_folders"
	McpToolEmailListEmails       McpTool = "email.list_emails"
	McpToolEmailReadEmail        McpTool = "email.read_email"
	McpToolEmailSendEmail        McpTool = "email.send_email"

	McpToolGitlabAddComment          McpTool = "gitlab.add_comment"
	McpToolGitlabAddMrComment        McpTool = "gitlab.add_mr_comment"
	McpToolGitlabCreateIssue         McpTool = "gitlab.create_issue"
	McpToolGitlabCreateMergeRequest  McpTool = "gitlab.create_merge_request"
	McpToolGitlabGetCurrentUser      McpTool = "gitlab.get_current_user"
	McpToolGitlabGetMergeRequestDiff McpTool = "gitlab.get_merge_request_diff"
	McpToolGitlabListIssues          McpTool = "gitlab.list_issues"
	McpToolGitlabListMergeRequests   McpTool = "gitlab.list_merge_requests"
	McpToolGitlabUpdateMergeRequest  McpTool = "gitlab.update_merge_request"

	McpToolTelegramAnswerCallbackQuery McpTool = "telegram.answer_callback_query"
	McpToolTelegramEditMessageText     McpTool = "telegram.edit_message_text"
	McpToolTelegramGetMe               McpTool = "telegram.get_me"
	McpToolTelegramSendMessage         McpTool = "telegram.send_message"
	McpToolTelegramSetWebhook          McpTool = "telegram.set_webhook"

	McpToolTrelloAddComment       McpTool = "trello.add_comment"
	McpToolTrelloAddLabelToCard   McpTool = "trello.add_label_to_card"
	McpToolTrelloArchiveCard      McpTool = "trello.archive_card"
	McpToolTrelloCreateCard       McpTool = "trello.create_card"
	McpToolTrelloGetCard          McpTool = "trello.get_card"
	McpToolTrelloGetCurrentUser   McpTool = "trello.get_current_user"
	McpToolTrelloListBoards       McpTool = "trello.list_boards"
	McpToolTrelloListCardComments McpTool = "trello.list_card_comments"
	McpToolTrelloListCards        McpTool = "trello.list_cards"
	McpToolTrelloListLabels       McpTool = "trello.list_labels"
	McpToolTrelloListLists        McpTool = "trello.list_lists"
	McpToolTrelloMoveCard         McpTool = "trello.move_card"
)

// AllMcpNames returns every McpName const. It backs the enum drift guard in
// tests/e2e/mcp_enum_test.go (which can only see exported identifiers) and lets callers
// enumerate the set. Keep in sync with the const block above and migration 080.
func AllMcpNames() []McpName {
	return []McpName{
		McpNameEmail,
		McpNameGitlab,
		McpNameTelegram,
		McpNameTrello,
	}
}

// AllMcpTools returns every McpTool const. It backs the enum drift guard in
// tests/e2e/mcp_enum_test.go and lets callers enumerate the set. Keep in sync with the const
// block above and the migrations that seed MoM tools.
func AllMcpTools() []McpTool {
	return []McpTool{
		McpToolEmailListEmailFolders,
		McpToolEmailListEmails,
		McpToolEmailReadEmail,
		McpToolEmailSendEmail,
		McpToolGitlabAddComment,
		McpToolGitlabAddMrComment,
		McpToolGitlabCreateIssue,
		McpToolGitlabCreateMergeRequest,
		McpToolGitlabGetCurrentUser,
		McpToolGitlabGetMergeRequestDiff,
		McpToolGitlabListIssues,
		McpToolGitlabListMergeRequests,
		McpToolGitlabUpdateMergeRequest,
		McpToolTelegramAnswerCallbackQuery,
		McpToolTelegramEditMessageText,
		McpToolTelegramGetMe,
		McpToolTelegramSendMessage,
		McpToolTelegramSetWebhook,
		McpToolTrelloAddComment,
		McpToolTrelloAddLabelToCard,
		McpToolTrelloArchiveCard,
		McpToolTrelloCreateCard,
		McpToolTrelloGetCard,
		McpToolTrelloGetCurrentUser,
		McpToolTrelloListBoards,
		McpToolTrelloListCardComments,
		McpToolTrelloListCards,
		McpToolTrelloListLabels,
		McpToolTrelloListLists,
		McpToolTrelloMoveCard,
	}
}

// mcpOf returns the McpName half of a "<mcp>.<tool>" McpTool, or "" if it has no ".".
func mcpOf(tool McpTool) McpName {
	s := string(tool)

	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			return McpName(s[:i])
		}
	}

	return ""
}
