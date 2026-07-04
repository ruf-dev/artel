package tract

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// resolver renders template strings against one run's context: the trigger payload and the
// outputs recorded by steps that have already executed, keyed by step id. Values are
// already-unmarshaled JSON (interface{}), not raw bytes.
type resolver struct {
	trigger interface{}
	outputs map[string]interface{}
	// now overrides $now for deterministic tests; zero value means time.Now().
	now time.Time
}

// render renders a single param/condition-side template string. Grammar: `{{ <expr> }}`
// where expr is `trigger.<path>`, `<step_id>.<path>`, `$now`, or `length(<ref>)`; path is dot
// fields plus `[N]` indexing. If the whole (trimmed) string is exactly one such token, the
// resolved value is returned as its raw typed value; otherwise every embedded token is
// stringified and substituted into the surrounding text. Literal `${{ ... }}` sequences
// (MoM's own params/secrets interpolation, which runs later) pass through untouched. Resolved
// values are never re-scanned for further tokens — this is a single pass over the input.
func (r *resolver) render(tmpl string) (interface{}, error) {
	if expr, ok := singleToken(strings.TrimSpace(tmpl)); ok {
		return r.eval(expr)
	}

	var sb strings.Builder

	pos := 0
	for pos < len(tmpl) {
		if isEscapedOpen(tmpl, pos) {
			end := strings.Index(tmpl[pos:], "}}")
			if end == -1 {
				return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "unterminated ${{ ... }} sequence")
			}

			sb.WriteString(tmpl[pos : pos+end+2])
			pos += end + len("}}")

			continue
		}

		if isOpen(tmpl, pos) {
			end := strings.Index(tmpl[pos:], "}}")
			if end == -1 {
				return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "unterminated {{ ... }} token")
			}

			expr := strings.TrimSpace(tmpl[pos+2 : pos+end])

			value, err := r.eval(expr)
			if err != nil {
				return nil, err
			}

			sb.WriteString(stringify(value))

			pos += end + len("}}")

			continue
		}

		sb.WriteByte(tmpl[pos])

		pos++
	}

	return sb.String(), nil
}

// extractRefs scans tmpl for `{{ ... }}` tract tokens (skipping literal `${{ ... }}` MoM
// escapes) and returns each token'tmpl raw reference expression — `length(...)` wrappers are
// unwrapped to their inner reference, and `$`-variables (e.g. `$now`) are omitted since they
// have no visibility/schema implications. Used by validation, not execution.
func extractRefs(tmpl string) ([]string, error) {
	var refs []string

	pos := 0
	for pos < len(tmpl) {
		if isEscapedOpen(tmpl, pos) {
			end := strings.Index(tmpl[pos:], "}}")
			if end == -1 {
				return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "unterminated ${{ ... }} sequence")
			}

			pos += end + len("}}")

			continue
		}

		if isOpen(tmpl, pos) {
			end := strings.Index(tmpl[pos:], "}}")
			if end == -1 {
				return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "unterminated {{ ... }} token")
			}

			expr := strings.TrimSpace(tmpl[pos+2 : pos+end])
			if !strings.HasPrefix(expr, "$") {
				ref := expr
				if strings.HasPrefix(expr, "length(") && strings.HasSuffix(expr, ")") {
					ref = strings.TrimSpace(expr[len("length(") : len(expr)-1])
				}

				refs = append(refs, ref)
			}

			pos += end + len("}}")

			continue
		}

		pos++
	}

	return refs, nil
}

func isOpen(s string, i int) bool {
	return s[i] == '{' && i+1 < len(s) && s[i+1] == '{'
}

func isEscapedOpen(s string, i int) bool {
	return s[i] == '$' && i+2 < len(s) && s[i+1] == '{' && s[i+2] == '{'
}

// singleToken reports whether trimmed is exactly one non-escaped `{{ ... }}` token with
// nothing else around it, returning its inner expression.
func singleToken(trimmed string) (string, bool) {
	if !strings.HasPrefix(trimmed, "{{") || !strings.HasSuffix(trimmed, "}}") {
		return "", false
	}

	inner := trimmed[2 : len(trimmed)-2]
	if strings.Contains(inner, "{{") || strings.Contains(inner, "}}") {
		return "", false
	}

	return strings.TrimSpace(inner), true
}

func (r *resolver) eval(expr string) (interface{}, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "empty expression")
	}

	if strings.HasPrefix(expr, "$") {
		return r.evalVar(expr)
	}

	if strings.HasPrefix(expr, "length(") && strings.HasSuffix(expr, ")") {
		inner := strings.TrimSpace(expr[len("length(") : len(expr)-1])

		value, err := r.resolveRef(inner)
		if err != nil {
			return nil, err
		}

		return lengthOf(value)
	}

	return r.resolveRef(expr)
}

func (r *resolver) evalVar(expr string) (interface{}, error) {
	switch expr {
	case "$now":
		now := r.now
		if now.IsZero() {
			now = time.Now()
		}

		return now.Format(time.RFC3339), nil
	default:
		return nil, rerrors.Wrap(user_errors.TractUnknownTemplateVar, expr)
	}
}

func (r *resolver) resolveRef(ref string) (interface{}, error) {
	base, segments, err := splitRef(ref)
	if err != nil {
		return nil, err
	}

	var root interface{}
	if base == "trigger" {
		root = r.trigger
	} else {
		value, ok := r.outputs[base]
		if !ok {
			return nil, rerrors.Wrap(user_errors.TractStepNotFound, base)
		}

		root = value
	}

	return navigate(root, segments)
}

type pathSegment struct {
	field   string
	index   int
	isIndex bool
}

// splitRef splits a reference expression (e.g. "list_prs[0].target_branch") into its base
// identifier ("list_prs") and the remaining path segments ([0], .target_branch).
func splitRef(ref string) (string, []pathSegment, error) {
	if ref == "" {
		return "", nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "empty reference")
	}

	i := 0
	for i < len(ref) && ref[i] != '.' && ref[i] != '[' {
		i++
	}

	base := ref[:i]
	if base == "" {
		return "", nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "reference missing base identifier: "+ref)
	}

	segments, err := parsePathSegments(ref[i:])
	if err != nil {
		return "", nil, err
	}

	return base, segments, nil
}

func parsePathSegments(rest string) ([]pathSegment, error) {
	var segments []pathSegment

	pos := 0
	for pos < len(rest) {
		switch rest[pos] {
		case '.':
			pos++
			start := pos

			for pos < len(rest) && rest[pos] != '.' && rest[pos] != '[' {
				pos++
			}

			name := rest[start:pos]
			if name == "" {
				return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "empty path segment in: "+rest)
			}

			segments = append(segments, pathSegment{field: name})

		case '[':
			end := strings.IndexByte(rest[pos:], ']')
			if end == -1 {
				return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "unterminated [ index in: "+rest)
			}

			numStr := rest[pos+1 : pos+end]

			idx, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "non-numeric index in: "+rest)
			}

			segments = append(segments, pathSegment{index: idx, isIndex: true})
			pos += end + 1

		default:
			return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "unexpected character in path: "+rest)
		}
	}

	return segments, nil
}

func navigate(value interface{}, segments []pathSegment) (interface{}, error) {
	cur := value

	for _, seg := range segments {
		if seg.isIndex {
			arr, ok := cur.([]interface{})
			if !ok {
				return nil, rerrors.Wrap(user_errors.TractStepNotFound, "index access on a non-array value")
			}

			if seg.index < 0 || seg.index >= len(arr) {
				return nil, rerrors.Wrap(user_errors.TractStepNotFound, "index out of range")
			}

			cur = arr[seg.index]

			continue
		}

		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil, rerrors.Wrap(user_errors.TractStepNotFound, "field access on a non-object value: "+seg.field)
		}

		fieldVal, ok := m[seg.field]
		if !ok {
			return nil, rerrors.Wrap(user_errors.TractStepNotFound, "field not found: "+seg.field)
		}

		cur = fieldVal
	}

	return cur, nil
}

func lengthOf(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case []interface{}:
		return len(v), nil
	case string:
		return len(v), nil
	case map[string]interface{}:
		return len(v), nil
	case nil:
		return 0, nil
	default:
		return nil, rerrors.Wrap(user_errors.TractMalformedTemplate, "length() requires an array, string, or object")
	}
}

func stringify(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}

		return string(data)
	}
}
