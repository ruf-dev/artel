package tract

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testResolver() *resolver {
	trigger := map[string]interface{}{
		"branch": "feature-123",
		"project": map[string]interface{}{
			"id":   float64(42),
			"name": "artel",
		},
		"commits": []interface{}{
			map[string]interface{}{"message": "first"},
			map[string]interface{}{"message": "second"},
		},
	}

	outputs := map[string]interface{}{
		"list_prs": []interface{}{
			map[string]interface{}{"target_branch": "main", "count": float64(3)},
		},
		"create_pr": map[string]interface{}{
			"url": "https://example.com/mr/1",
		},
		"count_step": float64(5),
	}

	fixedNow := time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC)

	r := &resolver{trigger: trigger, outputs: outputs, now: fixedNow}

	return r
}

func TestResolver_Render_TypedSingleToken(t *testing.T) {
	r := testResolver()

	tests := []struct {
		name string
		expr string
		want interface{}
	}{
		{"string field", "{{ trigger.branch }}", "feature-123"},
		{"numeric field", "{{ trigger.project.id }}", float64(42)},
		{"nested field", "{{ trigger.project.name }}", "artel"},
		{"array index then field", "{{ list_prs[0].target_branch }}", "main"},
		{"array index numeric field", "{{ list_prs[0].count }}", float64(3)},
		{"whole-output numeric", "{{ count_step }}", float64(5)},
		{"length of array", "{{ length(trigger.commits) }}", 2},
		{"length of string", "{{ length(trigger.branch) }}", len("feature-123")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.render(tt.expr)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResolver_Render_Stringification(t *testing.T) {
	r := testResolver()

	got, err := r.render("branch is {{ trigger.branch }} for project {{ trigger.project.id }}")
	require.NoError(t, err)
	assert.Equal(t, "branch is feature-123 for project 42", got)
}

func TestResolver_Render_Now(t *testing.T) {
	r := testResolver()

	got, err := r.render("{{ $now }}")
	require.NoError(t, err)
	assert.Equal(t, "2026-07-03T12:00:00Z", got)

	embedded, err := r.render("at {{ $now }} exactly")
	require.NoError(t, err)
	assert.Equal(t, "at 2026-07-03T12:00:00Z exactly", embedded)
}

func TestResolver_Render_UnknownVar(t *testing.T) {
	r := testResolver()

	_, err := r.render("{{ $bogus }}")
	require.Error(t, err)
}

func TestResolver_Render_MissingRef(t *testing.T) {
	r := testResolver()

	tests := []string{
		"{{ nonexistent_step.field }}",
		"{{ trigger.does.not.exist }}",
		"{{ list_prs[5].target_branch }}",
		"{{ list_prs[0].nonexistent }}",
	}

	for _, expr := range tests {
		t.Run(expr, func(t *testing.T) {
			_, err := r.render(expr)
			assert.Error(t, err)
		})
	}
}

func TestResolver_Render_MalformedTokens(t *testing.T) {
	r := testResolver()

	tests := []string{
		"{{ trigger.branch",   // unterminated
		"{{ }}",               // empty expression
		"{{ .leadingdot }}",   // missing base identifier
		"{{ list_prs[abc] }}", // non-numeric index
		"{{ list_prs[0 }}",    // unterminated bracket
	}

	for _, expr := range tests {
		t.Run(expr, func(t *testing.T) {
			_, err := r.render(expr)
			assert.Error(t, err)
		})
	}
}

func TestResolver_Render_MomEscapePassthrough(t *testing.T) {
	r := testResolver()

	got, err := r.render("${{params.project_id}}")
	require.NoError(t, err)
	assert.Equal(t, "${{params.project_id}}", got)

	mixed, err := r.render("prefix ${{params.project_id}} and {{ trigger.branch }}")
	require.NoError(t, err)
	assert.Equal(t, "prefix ${{params.project_id}} and feature-123", mixed)
}

func TestResolver_Render_NoReentry(t *testing.T) {
	r := &resolver{
		trigger: map[string]interface{}{
			"template_lookalike": "{{ trigger.branch }}",
		},
		outputs: map[string]interface{}{},
	}

	// The resolved value of trigger.template_lookalike happens to look like another token —
	// it must be returned/embedded verbatim, not re-scanned for further substitution.
	got, err := r.render("{{ trigger.template_lookalike }}")
	require.NoError(t, err)
	assert.Equal(t, "{{ trigger.branch }}", got)

	embedded, err := r.render("value: {{ trigger.template_lookalike }}")
	require.NoError(t, err)
	assert.Equal(t, "value: {{ trigger.branch }}", embedded)
}

func TestExtractRefs(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{"single ref", "{{ trigger.branch }}", []string{"trigger.branch"}},
		{
			"embedded refs",
			"{{ trigger.branch }} and {{ list_prs[0].id }}",
			[]string{"trigger.branch", "list_prs[0].id"},
		},
		{"length unwrapped", "{{ length(trigger.commits) }}", []string{"trigger.commits"}},
		{"dollar var omitted", "{{ $now }}", nil},
		{"mom escape skipped", "${{params.x}} {{ trigger.branch }}", []string{"trigger.branch"}},
		{"no tokens", "just plain text", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractRefs(tt.s)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSplitRef(t *testing.T) {
	base, segments, err := splitRef("list_prs[0].target_branch")
	require.NoError(t, err)
	assert.Equal(t, "list_prs", base)
	require.Len(t, segments, 2)
	assert.True(t, segments[0].isIndex)
	assert.Equal(t, 0, segments[0].index)
	assert.Equal(t, "target_branch", segments[1].field)

	base, segments, err = splitRef("trigger")
	require.NoError(t, err)
	assert.Equal(t, "trigger", base)
	assert.Empty(t, segments)
}
