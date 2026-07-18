package script

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

const defaultTimeout = 5 * time.Second

// JavaScriptEngine runs script steps through goja, a pure-Go JS interpreter (no cgo). Each
// Run gets a fresh *goja.Runtime with nothing set on it besides __args — no require,
// fetch, filesystem, or process access exists to sandbox against, because nothing exposes
// it in the first place.
type JavaScriptEngine struct {
	timeout time.Duration
}

func NewJavaScriptEngine() *JavaScriptEngine {
	engine := &JavaScriptEngine{timeout: defaultTimeout}

	return engine
}

func (e *JavaScriptEngine) Language() domain.ScriptLanguage {
	return domain.ScriptLanguageJavaScript
}

func (e *JavaScriptEngine) Run(ctx context.Context, in RunInput) (map[string]interface{}, error) {
	vm := goja.New()

	argv := make([]interface{}, len(in.InputParams))
	for i, p := range in.InputParams {
		argv[i] = in.Args[p.Name]
	}

	vm.Set("__args", argv)

	done := make(chan struct{})
	defer close(done)

	timer := time.AfterFunc(e.timeout, func() {
		vm.Interrupt("script execution timed out")
	})
	defer timer.Stop()

	go watchCancellation(ctx, vm, done)

	source := buildSource(in)

	value, err := vm.RunString(source)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.TractScriptRuntimeError, err.Error())
	}

	exported, ok := value.Export().(map[string]interface{})
	if !ok {
		return nil, rerrors.Wrap(user_errors.TractScriptRuntimeError, "script did not return an object")
	}

	output, err := coerceOutput(exported, in.OutputParams)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func watchCancellation(ctx context.Context, vm *goja.Runtime, done chan struct{}) {
	select {
	case <-ctx.Done():
		vm.Interrupt("context canceled")
	case <-done:
	}
}

// buildSource wraps the user-authored body with the generated signature, a zero-valued
// declaration for every output param, and a return statement listing them — the only parts
// of the executed source that ever came from persisted step data are in.Body and the
// param names/types, all of which passed validation before the step could be saved.
func buildSource(in RunInput) string {
	var sb strings.Builder

	sb.WriteString("(function(")
	writeNames(&sb, in.InputParams)
	sb.WriteString(") {\n")

	if len(in.OutputParams) > 0 {
		sb.WriteString("  let ")

		for i, p := range in.OutputParams {
			if i > 0 {
				sb.WriteString(", ")
			}

			fmt.Fprintf(&sb, "%s = %s", p.Name, zeroValueLiteral(p.Property.Type))
		}

		sb.WriteString(";\n\n")
	}

	sb.WriteString(in.Body)
	sb.WriteString("\n\n  return {")

	for i, p := range in.OutputParams {
		if i > 0 {
			sb.WriteString(", ")
		}

		fmt.Fprintf(&sb, "%s: %s", p.Name, p.Name)
	}

	sb.WriteString("};\n})(")

	for i := range in.InputParams {
		if i > 0 {
			sb.WriteString(", ")
		}

		fmt.Fprintf(&sb, "__args[%d]", i)
	}

	sb.WriteString(")")

	return sb.String()
}

func writeNames(sb *strings.Builder, params []domain.ScriptParam) {
	for i, p := range params {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(p.Name)
	}
}

func zeroValueLiteral(propType string) string {
	switch propType {
	case "string":
		return `""`
	case "integer", "number":
		return "0"
	case "boolean":
		return "false"
	case "array":
		return "[]"
	case "object":
		return "{}"
	default:
		return "null"
	}
}

// coerceOutput picks exactly the declared output names out of raw (dropping anything else
// the script may have smuggled onto the returned object — unreachable via the generated
// return statement, but defensive against a future wrapper change) and validates each
// value's JS-inferred type against its declared Property.Type.
func coerceOutput(raw map[string]interface{}, outputParams []domain.ScriptParam) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(outputParams))

	for _, p := range outputParams {
		value := raw[p.Name]

		if !matchesType(value, p.Property.Type) {
			msg := fmt.Sprintf(
				"output %q: expected %s, got %s",
				p.Name, p.Property.Type, describeValue(value),
			)

			return nil, rerrors.Wrap(user_errors.TractScriptOutputTypeMismatch, msg)
		}

		result[p.Name] = value
	}

	return result, nil
}

func matchesType(value interface{}, propType string) bool {
	switch propType {
	case "string":
		_, ok := value.(string)
		return ok
	case "integer", "number":
		switch value.(type) {
		case int64, float64:
			return true
		default:
			return false
		}
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "array":
		_, ok := value.([]interface{})
		return ok
	case "object":
		_, ok := value.(map[string]interface{})
		return ok
	default:
		return true
	}
}

func describeValue(value interface{}) string {
	if value == nil {
		return "undefined"
	}

	return fmt.Sprintf("%T", value)
}
