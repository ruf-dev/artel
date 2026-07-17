package executors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/utils"
)

const (
	httpExecutorTimeout  = 15 * time.Second
	httpMaxResponseBytes = 1 << 20 // 1MB

	// regexSubmatchWithGroup is len(FindStringSubmatch(...)) for a pattern with one capture
	// group: index 0 is the full match, index 1 is the group.
	regexSubmatchWithGroup = 2
)

var (
	paramInterpolationRegexp  = regexp.MustCompile(`\$\{\{params\.(\w+)\}\}`)
	secretInterpolationRegexp = regexp.MustCompile(`\$\{\{secrets\.(\w+)\}\}`)
)

const secretValuePrefix = "__secrets."

type HttpExecutor struct {
	client *http.Client
}

func NewHttpExecutor() *HttpExecutor {
	c := &http.Client{
		Timeout:   httpExecutorTimeout,
		Transport: newLoggingTransport(),
	}

	return &HttpExecutor{
		client: c,
	}
}

func (e *HttpExecutor) Execute(ctx context.Context, action domain.ToolAction,
	secrets map[string]interface{}, params map[string]interface{},
) (string, error) {
	if action.Http == nil {
		return "", user_errors.McpActionMissing
	}

	httpAction := action.Http

	reqUrl := interpolateParams(httpAction.Url, params)
	reqUrl = interpolateSecrets(reqUrl, secrets)

	query := url.Values{}

	for key, value := range httpAction.Query {
		resolved, err := resolveActionValue(value, params, secrets)
		if err != nil {
			return "", err
		}

		if resolved == "" {
			continue
		}

		query.Set(key, resolved)
	}

	if len(query) > 0 {
		parsedUrl, err := url.Parse(reqUrl)
		if err != nil {
			return "", rerrors.Wrap(err, "error parsing http action url")
		}

		existing := parsedUrl.Query()

		for key, values := range query {
			for _, v := range values {
				existing.Add(key, v)
			}
		}

		parsedUrl.RawQuery = existing.Encode()
		reqUrl = parsedUrl.String()
	}

	var bodyReader io.Reader

	if len(httpAction.Body) > 0 {
		renderedBody, err := renderHttpBody(httpAction.Body, params, secrets)
		if err != nil {
			return "", err
		}

		bodyReader = bytes.NewReader(renderedBody)
	}

	req, err := http.NewRequestWithContext(ctx, httpAction.Method, reqUrl, bodyReader)
	if err != nil {
		return "", rerrors.Wrap(err, "error building http request")
	}

	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range httpAction.Headers {
		var resolved string
		resolved, err = resolveActionValue(value, params, secrets)
		if err != nil {
			return "", err
		}

		if resolved == "" {
			continue
		}

		req.Header.Set(key, resolved)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return "", rerrors.Wrap(err, "error executing http request")
	}
	defer utils.CloseWithLog(resp.Body, "http executor response body")

	limitedReader := io.LimitReader(resp.Body, httpMaxResponseBytes)

	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading http response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMessage := fmt.Sprintf("status %d: %s", resp.StatusCode, body)

		return "", rerrors.Wrap(user_errors.McpHttpRequestFailed, errMessage)
	}

	return string(body), nil
}

// renderHttpBody walks the HttpAction.Body JSON template and interpolates every string leaf
// via resolveActionValue (same ${{params.*}} / __secrets.* rules as headers and query values),
// then re-marshals the result to send as the request body.
func renderHttpBody(
	body json.RawMessage,
	params map[string]interface{},
	secrets map[string]interface{},
) (json.RawMessage, error) {
	var parsed interface{}

	err := json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing http action body template")
	}

	rendered, err := renderBodyValue(parsed, params, secrets)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(rendered)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshaling rendered http body")
	}

	return data, nil
}

func renderBodyValue(
	value interface{},
	params map[string]interface{},
	secrets map[string]interface{},
) (interface{}, error) {
	switch typedValue := value.(type) {
	case string:
		resolved, err := resolveActionValue(typedValue, params, secrets)
		if err != nil {
			return nil, err
		}

		return resolved, nil

	case map[string]interface{}:
		rendered := make(map[string]interface{}, len(typedValue))

		for key, val := range typedValue {
			renderedVal, err := renderBodyValue(val, params, secrets)
			if err != nil {
				return nil, err
			}

			if renderedStr, ok := renderedVal.(string); ok && renderedStr == "" {
				continue
			}

			rendered[key] = renderedVal
		}

		return rendered, nil

	case []interface{}:
		rendered := make([]interface{}, len(typedValue))

		for i, val := range typedValue {
			renderedVal, err := renderBodyValue(val, params, secrets)
			if err != nil {
				return nil, err
			}

			rendered[i] = renderedVal
		}

		return rendered, nil

	default:
		return typedValue, nil
	}
}

// resolveActionValue resolves a single header/query value. If the value is exactly
// "__secrets.<field>", it is replaced wholesale with secrets[field] (hard error if missing).
// Otherwise the value is run through ${{params.*}} interpolation.
func resolveActionValue(value string, params map[string]interface{}, secrets map[string]interface{}) (string, error) {
	field, ok := secretField(value)
	if !ok {
		return interpolateParams(value, params), nil
	}

	secretValue, ok := secrets[field]
	if !ok {
		return "", user_errors.McpSecretFieldMissing
	}

	return stringifyValue(secretValue), nil
}

// stringifyValue formats a param/secret value for embedding in a URL, header, or query
// string. float64 needs special handling because encoding/json decodes all bare JSON numbers
// into float64, and fmt.Sprint's default formatting switches large values (e.g. numeric IDs)
// to exponential notation, which breaks when the result is used as an identifier.
func stringifyValue(value interface{}) string {
	floatValue, ok := value.(float64)
	if !ok {
		return fmt.Sprint(value)
	}

	return strconv.FormatFloat(floatValue, 'f', -1, 64)
}

func secretField(value string) (string, bool) {
	if len(value) <= len(secretValuePrefix) {
		return "", false
	}

	if value[:len(secretValuePrefix)] != secretValuePrefix {
		return "", false
	}

	return value[len(secretValuePrefix):], true
}

func interpolateParams(value string, params map[string]interface{}) string {
	return paramInterpolationRegexp.ReplaceAllStringFunc(value, func(match string) string {
		submatches := paramInterpolationRegexp.FindStringSubmatch(match)
		if len(submatches) < regexSubmatchWithGroup {
			return match
		}

		paramValue, ok := params[submatches[1]]
		if !ok {
			return ""
		}

		return stringifyValue(paramValue)
	})
}

// interpolateSecrets replaces ${{secrets.field}} tokens within a string (e.g. a URL's
// host) with the corresponding value from secrets. Unlike resolveActionValue's
// __secrets. convention, this allows secrets to be embedded inside a larger string
// rather than only resolving a header/query value wholesale.
func interpolateSecrets(value string, secrets map[string]interface{}) string {
	return secretInterpolationRegexp.ReplaceAllStringFunc(value, func(match string) string {
		submatches := secretInterpolationRegexp.FindStringSubmatch(match)
		if len(submatches) < regexSubmatchWithGroup {
			return match
		}

		secretValue, ok := secrets[submatches[1]]
		if !ok {
			return ""
		}

		return stringifyValue(secretValue)
	})
}
