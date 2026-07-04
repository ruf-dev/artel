package tract

import (
	"testing"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluate_NumericVsString(t *testing.T) {
	r := &resolver{trigger: map[string]interface{}{}, outputs: map[string]interface{}{
		"count": float64(3),
	}}

	tests := []struct {
		name string
		cond domain.TractCondition
		want bool
	}{
		{"numeric equal, numeric literal", domain.TractCondition{Left: "{{ count }}", Op: "==", Right: "3"}, true},
		{"numeric equal, string vs number", domain.TractCondition{Left: "{{ count }}", Op: "==", Right: "3.0"}, true},
		{"numeric not equal", domain.TractCondition{Left: "{{ count }}", Op: "!=", Right: "4"}, true},
		{"string equal, non-numeric", domain.TractCondition{Left: "abc", Op: "==", Right: "abc"}, true},
		{"string not equal, non-numeric", domain.TractCondition{Left: "abc", Op: "==", Right: "xyz"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _, _, err := evaluate(tt.cond, r)
			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestEvaluate_AllOperators(t *testing.T) {
	r := &resolver{trigger: map[string]interface{}{}, outputs: map[string]interface{}{}}

	tests := []struct {
		name string
		cond domain.TractCondition
		want bool
	}{
		{"gt true", domain.TractCondition{Left: "5", Op: ">", Right: "3"}, true},
		{"gt false", domain.TractCondition{Left: "2", Op: ">", Right: "3"}, false},
		{"lt true", domain.TractCondition{Left: "2", Op: "<", Right: "3"}, true},
		{"gte equal", domain.TractCondition{Left: "3", Op: ">=", Right: "3"}, true},
		{"lte equal", domain.TractCondition{Left: "3", Op: "<=", Right: "3"}, true},
		{"contains substring", domain.TractCondition{Left: "feature-123", Op: "contains", Right: "feat"}, true},
		{"contains miss", domain.TractCondition{Left: "feature-123", Op: "contains", Right: "bugfix"}, false},
		{"glob match", domain.TractCondition{Left: "feature-123", Op: "glob", Right: "feature-*"}, true},
		{"glob miss", domain.TractCondition{Left: "bugfix-123", Op: "glob", Right: "feature-*"}, false},
		{"regex match", domain.TractCondition{Left: "feature-123", Op: "regex", Right: `^feature-\d+$`}, true},
		{"regex miss", domain.TractCondition{Left: "feature-abc", Op: "regex", Right: `^feature-\d+$`}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _, _, err := evaluate(tt.cond, r)
			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestEvaluate_NumericOperatorsRequireNumbers(t *testing.T) {
	r := &resolver{trigger: map[string]interface{}{}, outputs: map[string]interface{}{}}

	cond := domain.TractCondition{Left: "abc", Op: ">", Right: "3"}
	_, _, _, err := evaluate(cond, r)
	assert.Error(t, err)
}

func TestEvaluate_ContainsOnArray(t *testing.T) {
	r := &resolver{trigger: map[string]interface{}{}, outputs: map[string]interface{}{
		"labels": []interface{}{"actual release", "bug"},
	}}

	cond := domain.TractCondition{Left: "{{ labels }}", Op: "contains", Right: "actual release"}
	result, _, _, err := evaluate(cond, r)
	require.NoError(t, err)
	assert.True(t, result)

	missCond := domain.TractCondition{Left: "{{ labels }}", Op: "contains", Right: "nope"}
	result, _, _, err = evaluate(missCond, r)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestEvaluateAll_ANDSemantics(t *testing.T) {
	r := &resolver{trigger: map[string]interface{}{}, outputs: map[string]interface{}{}}

	allTrue := []domain.TractCondition{
		{Left: "1", Op: "==", Right: "1"},
		{Left: "2", Op: "==", Right: "2"},
	}
	result, err := evaluateAll(allTrue, r)
	require.NoError(t, err)
	assert.True(t, result)

	oneFalse := []domain.TractCondition{
		{Left: "1", Op: "==", Right: "1"},
		{Left: "2", Op: "==", Right: "3"},
	}
	result, err = evaluateAll(oneFalse, r)
	require.NoError(t, err)
	assert.False(t, result)

	empty := []domain.TractCondition{}
	result, err = evaluateAll(empty, r)
	require.NoError(t, err)
	assert.True(t, result)
}

func TestEvaluate_UnknownOperator(t *testing.T) {
	r := &resolver{trigger: map[string]interface{}{}, outputs: map[string]interface{}{}}

	cond := domain.TractCondition{Left: "a", Op: "bogus", Right: "b"}
	_, _, _, err := evaluate(cond, r)
	assert.Error(t, err)
}
