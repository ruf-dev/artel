package tract

import (
	"encoding/json"
	"fmt"
	"sort"
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
			msg := fmt.Sprintf(
				"template {{ %s }} references step %q, but no step with that id has produced output yet in this run (available so far: %s) — check the step id, and that it runs before this one",
				ref, base, availableStepIds(r.outputs),
			)

			return nil, rerrors.Wrap(user_errors.TractStepNotFound, msg)
		}

		root = value
	}

	return navigate(root, base, segments, ref)
}

func availableStepIds(outputs map[string]interface{}) string {
	if len(outputs) == 0 {
		return "none yet"
	}

	ids := make([]string, 0, len(outputs))
	for id := range outputs {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	return strings.Join(ids, ", ")
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

// navigate walks segments starting from value (the already-resolved output of base). ref is
// the full original reference expression (e.g. "list_prs[0].source_branch"), used only to
// build precise error messages that quote the exact template the caller wrote; path tracks how
// much of it has been consumed so far so an error can point at the specific step that broke,
// not just the final segment.
func navigate(value interface{}, base string, segments []pathSegment, ref string) (interface{}, error) {
	cur := value
	path := base

	for _, seg := range segments {
		if seg.isIndex {
			arr, ok := cur.([]interface{})
			if !ok {
				msg := fmt.Sprintf(
					"template {{ %s }}: %s is %s, not a list, so [%d] cannot be applied",
					ref, path, describeValue(cur), seg.index,
				)

				return nil, rerrors.Wrap(user_errors.TractTemplateIndexOnNonArray, msg)
			}

			if seg.index < 0 || seg.index >= len(arr) {
				msg := fmt.Sprintf(
					"template {{ %s }}: %s has %d item(s), so index [%d] is out of range",
					ref, path, len(arr), seg.index,
				)

				return nil, rerrors.Wrap(user_errors.TractTemplateIndexOutOfRange, msg)
			}

			cur = arr[seg.index]
			path = fmt.Sprintf("%s[%d]", path, seg.index)

			continue
		}

		m, ok := cur.(map[string]interface{})
		if !ok {
			msg := fmt.Sprintf(
				"template {{ %s }}: %s is %s, not an object, so field %q cannot be accessed",
				ref, path, describeValue(cur), seg.field,
			)

			return nil, rerrors.Wrap(user_errors.TractTemplateFieldOnNonObject, msg)
		}

		fieldVal, ok := m[seg.field]
		if !ok {
			msg := fmt.Sprintf(
				"template {{ %s }}: %s has no field %q (available fields: %s)",
				ref, path, seg.field, availableFields(m),
			)

			return nil, rerrors.Wrap(user_errors.TractTemplateFieldNotFound, msg)
		}

		cur = fieldVal
		path = path + "." + seg.field
	}

	return cur, nil
}

func availableFields(m map[string]interface{}) string {
	if len(m) == 0 {
		return "none"
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return strings.Join(keys, ", ")
}

// describeValue renders a short, human-readable description of a resolved template value for
// use in error messages — e.g. "a string (\"main\")" rather than a raw Go %T/%v dump.
func describeValue(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case string:
		return "a string (" + truncateForError(v) + ")"
	case float64:
		return "a number (" + strconv.FormatFloat(v, 'f', -1, 64) + ")"
	case bool:
		return "a boolean (" + strconv.FormatBool(v) + ")"
	case []interface{}:
		return fmt.Sprintf("a list with %d item(s)", len(v))
	case map[string]interface{}:
		return "an object"
	default:
		return fmt.Sprintf("a %T", v)
	}
}

func truncateForError(s string) string {
	const maxLen = 40

	if len(s) <= maxLen {
		return strconv.Quote(s)
	}

	return strconv.Quote(s[:maxLen]) + "…"
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
