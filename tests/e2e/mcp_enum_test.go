//go:build e2e
// +build e2e

package e2e_test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
	"github.com/ruf-dev/artel/tests/harness"
)

// McpEnumSuite is the drift guard for the hand-written McpName / McpTool const block in
// internal/service/v1/mcp/executors/tool_enum.go: it reads the labels of the Postgres mcp_name /
// mcp_tool enums (created by migration 080_mcp_tool_enums.sql, extended by later tool-seed
// migrations) from the migrated e2e database and asserts they are exactly the Go const set, in
// both directions. Same "the real DB guards what Go believes" rule as the rest of tests/.
type McpEnumSuite struct {
	suite.Suite
}

func TestMcpEnum(t *testing.T) {
	suite.Run(t, new(McpEnumSuite))
}

func (s *McpEnumSuite) enumLabels(enumType string) []string {
	db := harness.OpenPostgres(s.T())

	rows, err := db.Query("SELECT unnest(enum_range(NULL::" + enumType + "))::text")
	s.Require().NoError(err)
	defer func() { _ = rows.Close() }()

	var labels []string

	for rows.Next() {
		var label string

		s.Require().NoError(rows.Scan(&label))

		labels = append(labels, label)
	}

	s.Require().NoError(rows.Err())

	return labels
}

func (s *McpEnumSuite) TestMcpNameEnumMatchesGoConsts() {
	want := make([]string, 0, len(executors.AllMcpNames()))
	for _, name := range executors.AllMcpNames() {
		want = append(want, string(name))
	}

	require.ElementsMatch(s.T(), want, s.enumLabels("mcp_name"),
		"executors.AllMcpNames() and the Postgres mcp_name enum must be the same set")
}

func (s *McpEnumSuite) TestMcpToolEnumMatchesGoConsts() {
	want := make([]string, 0, len(executors.AllMcpTools()))
	for _, tool := range executors.AllMcpTools() {
		want = append(want, string(tool))
	}

	require.ElementsMatch(s.T(), want, s.enumLabels("mcp_tool"),
		"executors.AllMcpTools() and the Postgres mcp_tool enum must be the same set")
}
