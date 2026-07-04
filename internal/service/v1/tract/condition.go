package tract

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// evaluate renders both sides of cond through r and compares them per cond.Op. It returns the
// boolean result plus the rendered left/right values (for callers that also want to record
// what was actually compared, e.g. the engine's run-step input log).
func evaluate(cond domain.TractCondition, r *resolver) (bool, interface{}, interface{}, error) {
	left, err := r.render(cond.Left)
	if err != nil {
		return false, nil, nil, err
	}

	right, err := r.render(cond.Right)
	if err != nil {
		return false, nil, nil, err
	}

	result, err := compare(cond.Op, left, right)
	if err != nil {
		return false, left, right, err
	}

	return result, left, right, nil
}

// EvaluateTriggerFilters reports whether every filter in filters holds against payload (AND
// semantics), rendering template refs with payload as the `trigger.*` context — no step
// outputs exist yet since filters run before any run is created. Exported so
// internal/transport/tract_webhook can reuse the condition evaluator for trigger_tracts link
// filters instead of duplicating it.
func EvaluateTriggerFilters(payload json.RawMessage, filters []domain.TractCondition) (bool, error) {
	if len(filters) == 0 {
		return true, nil
	}

	var trigger interface{}

	err := json.Unmarshal(payload, &trigger)
	if err != nil {
		trigger = map[string]interface{}{}
	}

	r := &resolver{trigger: trigger, outputs: map[string]interface{}{}}

	return evaluateAll(filters, r)
}

// evaluateAll AND's a list of conditions (short-circuits on the first false/error) — shared by
// condition steps and trigger_tracts link filters.
func evaluateAll(conditions []domain.TractCondition, r *resolver) (bool, error) {
	for _, cond := range conditions {
		result, _, _, err := evaluate(cond, r)
		if err != nil {
			return false, err
		}

		if !result {
			return false, nil
		}
	}

	return true, nil
}

func compare(op string, left, right interface{}) (bool, error) {
	switch op {
	case "==":
		return compareEqual(left, right), nil
	case "!=":
		return !compareEqual(left, right), nil
	case ">", "<", ">=", "<=":
		return compareNumeric(op, left, right)
	case "contains":
		return compareContains(left, right), nil
	case "glob":
		return compareGlob(left, right)
	case "regex":
		return compareRegex(left, right)
	default:
		return false, rerrors.Wrap(user_errors.TractUnknownStepType, "unknown condition operator: "+op)
	}
}

func numericValue(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0, false
		}

		return f, true
	default:
		return 0, false
	}
}

func compareEqual(a, b interface{}) bool {
	af, aok := numericValue(a)
	bf, bok := numericValue(b)

	if aok && bok {
		return af == bf
	}

	return stringify(a) == stringify(b)
}

func compareNumeric(op string, a, b interface{}) (bool, error) {
	af, aok := numericValue(a)
	bf, bok := numericValue(b)

	if !aok || !bok {
		return false, rerrors.Wrap(user_errors.TractMalformedTemplate, "numeric comparison requires numeric operands")
	}

	switch op {
	case ">":
		return af > bf, nil
	case "<":
		return af < bf, nil
	case ">=":
		return af >= bf, nil
	case "<=":
		return af <= bf, nil
	default:
		return false, rerrors.Wrap(user_errors.TractUnknownStepType, "unknown numeric operator: "+op)
	}
}

func compareContains(a, b interface{}) bool {
	if arr, ok := a.([]interface{}); ok {
		needle := stringify(b)
		for _, item := range arr {
			if stringify(item) == needle {
				return true
			}
		}

		return false
	}

	return strings.Contains(stringify(a), stringify(b))
}

func compareGlob(a, b interface{}) (bool, error) {
	matched, err := filepath.Match(stringify(b), stringify(a))
	if err != nil {
		return false, rerrors.Wrap(user_errors.TractMalformedTemplate, "invalid glob pattern: "+err.Error())
	}

	return matched, nil
}

func compareRegex(a, b interface{}) (bool, error) {
	matched, err := regexp.MatchString(stringify(b), stringify(a))
	if err != nil {
		return false, rerrors.Wrap(user_errors.TractMalformedTemplate, "invalid regex: "+err.Error())
	}

	return matched, nil
}
