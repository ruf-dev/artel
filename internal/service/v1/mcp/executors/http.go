package executors

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/utils"
)

const (
	httpExecutorTimeout  = 15 * time.Second
	httpMaxResponseBytes = 1 << 20 // 1MB
)

var paramInterpolationRegexp = regexp.MustCompile(`\$\{\{params\.(\w+)\}\}`)
var secretInterpolationRegexp = regexp.MustCompile(`\$\{\{secrets\.(\w+)\}\}`)

const secretValuePrefix = "__secrets."

type HttpExecutor struct {
	client *http.Client
}

func NewHttpExecutor() *HttpExecutor {
	c := &http.Client{
		Timeout: httpExecutorTimeout,
	}

	return &HttpExecutor{
		client: c,
	}
}

func (e *HttpExecutor) Execute(ctx context.Context, action domain.ToolAction,
	secrets map[string]interface{}, params map[string]interface{}) (string, error) {
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

	req, err := http.NewRequestWithContext(ctx, httpAction.Method, reqUrl, nil)
	if err != nil {
		return "", rerrors.Wrap(err, "error building http request")
	}

	for key, value := range httpAction.Headers {
		resolved, err := resolveActionValue(value, params, secrets)
		if err != nil {
			return "", err
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

	return fmt.Sprint(secretValue), nil
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
		if len(submatches) < 2 {
			return match
		}
		paramValue, ok := params[submatches[1]]
		if !ok {
			return ""
		}
		return fmt.Sprint(paramValue)
	})
}

// interpolateSecrets replaces ${{secrets.field}} tokens within a string (e.g. a URL's
// host) with the corresponding value from secrets. Unlike resolveActionValue's
// __secrets. convention, this allows secrets to be embedded inside a larger string
// rather than only resolving a header/query value wholesale.
func interpolateSecrets(value string, secrets map[string]interface{}) string {
	return secretInterpolationRegexp.ReplaceAllStringFunc(value, func(match string) string {
		submatches := secretInterpolationRegexp.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		secretValue, ok := secrets[submatches[1]]
		if !ok {
			return ""
		}
		return fmt.Sprint(secretValue)
	})
}
