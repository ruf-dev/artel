package script

import (
	"context"
	"testing"
	"time"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func numericValue(t *testing.T, v interface{}) float64 {
	t.Helper()

	switch n := v.(type) {
	case float64:
		return n
	case int64:
		return float64(n)
	default:
		t.Fatalf("expected a number, got %T (%v)", v, v)

		return 0
	}
}

func TestJavaScriptEngine_HappyPath(t *testing.T) {
	engine := NewJavaScriptEngine()
	in := RunInput{
		Body: "total = a + b;",
		InputParams: []domain.ScriptParam{
			{Name: "a", Property: domain.ToolProperty{Type: "number"}},
			{Name: "b", Property: domain.ToolProperty{Type: "number"}},
		},
		OutputParams: []domain.ScriptParam{
			{Name: "total", Property: domain.ToolProperty{Type: "number"}},
		},
		Args: map[string]interface{}{"a": float64(2), "b": float64(3)},
	}

	out, err := engine.Run(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, float64(5), numericValue(t, out["total"]))
}

func TestJavaScriptEngine_ArrayInputAggregation(t *testing.T) {
	engine := NewJavaScriptEngine()
	in := RunInput{
		Body: "for (let i = 0; i < nums.length; i++) { total += nums[i]; }",
		InputParams: []domain.ScriptParam{
			{Name: "nums", Property: domain.ToolProperty{Type: "array"}},
		},
		OutputParams: []domain.ScriptParam{
			{Name: "total", Property: domain.ToolProperty{Type: "number"}},
		},
		Args: map[string]interface{}{"nums": []interface{}{float64(1), float64(2), float64(3)}},
	}

	out, err := engine.Run(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, float64(6), numericValue(t, out["total"]))
}

func TestJavaScriptEngine_OutputTypeMismatch(t *testing.T) {
	engine := NewJavaScriptEngine()
	in := RunInput{
		Body: "result = 42;",
		OutputParams: []domain.ScriptParam{
			{Name: "result", Property: domain.ToolProperty{Type: "string"}},
		},
		Args: map[string]interface{}{},
	}

	_, err := engine.Run(context.Background(), in)
	assert.Error(t, err)
}

func TestJavaScriptEngine_NoAmbientGlobals(t *testing.T) {
	engine := NewJavaScriptEngine()
	in := RunInput{
		Body: "result = typeof require + \"|\" + typeof fetch + \"|\" + typeof process;",
		OutputParams: []domain.ScriptParam{
			{Name: "result", Property: domain.ToolProperty{Type: "string"}},
		},
		Args: map[string]interface{}{},
	}

	out, err := engine.Run(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, "undefined|undefined|undefined", out["result"])
}

func TestJavaScriptEngine_Timeout(t *testing.T) {
	engine := &JavaScriptEngine{timeout: 50 * time.Millisecond}
	in := RunInput{Body: "while (true) {}", Args: map[string]interface{}{}}

	start := time.Now()
	_, err := engine.Run(context.Background(), in)
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Less(t, elapsed, 2*time.Second)
}

func TestJavaScriptEngine_ContextCancellation(t *testing.T) {
	engine := NewJavaScriptEngine()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	in := RunInput{Body: "while (true) {}", Args: map[string]interface{}{}}

	start := time.Now()
	_, err := engine.Run(ctx, in)
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Less(t, elapsed, 2*time.Second)
}

func TestJavaScriptEngine_Language(t *testing.T) {
	engine := NewJavaScriptEngine()
	assert.Equal(t, domain.ScriptLanguageJavaScript, engine.Language())
}
