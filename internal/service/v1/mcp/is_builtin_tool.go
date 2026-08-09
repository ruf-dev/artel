package mcp

import "github.com/ruf-dev/artel/internal/service/v1/mcp/executors"

func (s *ServiceImpl) IsBuiltinTool(name string) bool {
	switch name {
	case executors.ToolListFiles, executors.ToolReadFile, executors.ToolWriteFile, executors.ToolDeleteFile,
		executors.ToolMoveFile, executors.ToolListFolders, executors.ToolListTags, executors.ToolGetFileMetadata,
		executors.ToolPgListTables, executors.ToolPgDescribeTable, executors.ToolPgQuery, executors.ToolPgExecute,
		toolConnections, toolConnectionsForTracts, toolCreateCommunityConnector:
		return true
	}

	return executors.IsTractTool(name)
}
