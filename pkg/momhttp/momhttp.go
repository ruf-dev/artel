// Package momhttp is the public surface for injecting per-mcp / per-tool HTTP middleware into
// artel's MoM http executor from outside the module. Pair it with pkg/app.WithToolHttpMiddleware
// / pkg/app.WithMcpHttpMiddleware and pkg/app.New.
//
// pkg/ importing internal/ within the same module is allowed (pkg/app already does it) — Go's
// internal rule only blocks importers outside the module.
package momhttp

import (
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
)

// Re-exported executor types, so embedders never import internal/ directly.
type (
	McpName    = executors.McpName
	McpTool    = executors.McpTool
	Middleware = executors.HttpMiddleware
	ToolIdent  = executors.ToolIdent
)

// Re-exported MoM tool-set identifiers (the Postgres mcp_name enum labels).
const (
	McpNameEmail    = executors.McpNameEmail
	McpNameGitlab   = executors.McpNameGitlab
	McpNameTelegram = executors.McpNameTelegram
	McpNameTrello   = executors.McpNameTrello
)

// Re-exported MoM tool identifiers (the Postgres mcp_tool enum labels, "<mcp>.<tool>").
const (
	McpToolEmailListEmailFolders = executors.McpToolEmailListEmailFolders
	McpToolEmailListEmails       = executors.McpToolEmailListEmails
	McpToolEmailReadEmail        = executors.McpToolEmailReadEmail
	McpToolEmailSendEmail        = executors.McpToolEmailSendEmail

	McpToolGitlabAddComment          = executors.McpToolGitlabAddComment
	McpToolGitlabAddMrComment        = executors.McpToolGitlabAddMrComment
	McpToolGitlabCreateIssue         = executors.McpToolGitlabCreateIssue
	McpToolGitlabCreateMergeRequest  = executors.McpToolGitlabCreateMergeRequest
	McpToolGitlabGetCurrentUser      = executors.McpToolGitlabGetCurrentUser
	McpToolGitlabGetMergeRequestDiff = executors.McpToolGitlabGetMergeRequestDiff
	McpToolGitlabListIssues          = executors.McpToolGitlabListIssues
	McpToolGitlabListMergeRequests   = executors.McpToolGitlabListMergeRequests
	McpToolGitlabUpdateMergeRequest  = executors.McpToolGitlabUpdateMergeRequest

	McpToolTelegramAnswerCallbackQuery = executors.McpToolTelegramAnswerCallbackQuery
	McpToolTelegramEditMessageText     = executors.McpToolTelegramEditMessageText
	McpToolTelegramGetMe               = executors.McpToolTelegramGetMe
	McpToolTelegramSendMessage         = executors.McpToolTelegramSendMessage
	McpToolTelegramSetWebhook          = executors.McpToolTelegramSetWebhook

	McpToolTrelloAddComment       = executors.McpToolTrelloAddComment
	McpToolTrelloAddLabelToCard   = executors.McpToolTrelloAddLabelToCard
	McpToolTrelloArchiveCard      = executors.McpToolTrelloArchiveCard
	McpToolTrelloCreateCard       = executors.McpToolTrelloCreateCard
	McpToolTrelloGetCard          = executors.McpToolTrelloGetCard
	McpToolTrelloGetCurrentUser   = executors.McpToolTrelloGetCurrentUser
	McpToolTrelloListBoards       = executors.McpToolTrelloListBoards
	McpToolTrelloListCardComments = executors.McpToolTrelloListCardComments
	McpToolTrelloListCards        = executors.McpToolTrelloListCards
	McpToolTrelloListLabels       = executors.McpToolTrelloListLabels
	McpToolTrelloListLists        = executors.McpToolTrelloListLists
	McpToolTrelloMoveCard         = executors.McpToolTrelloMoveCard
)
