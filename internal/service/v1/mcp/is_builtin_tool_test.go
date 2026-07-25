package mcp_test

import (
	"testing"

	"github.com/ruf-dev/artel/internal/service/v1/mcp"
	"github.com/ruf-dev/artel/internal/service/v1/subscription"
)

func TestIsBuiltinTool(t *testing.T) {
	svc := mcp.New(nil, nil, nil, nil, nil, nil, nil, nil, subscription.NewFree(), nil)

	cases := []struct {
		name string
		want bool
	}{
		{"connections", true},
		{"list_connections_for_tracts", true},
		{"create_community_connector", true},
		{"not_a_real_tool", false},
	}

	for _, c := range cases {
		got := svc.IsBuiltinTool(c.name)
		if got != c.want {
			t.Errorf("IsBuiltinTool(%q) = %v, want %v", c.name, got, c.want)
		}
	}
}
