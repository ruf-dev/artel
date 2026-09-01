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
	"slices"
	"strconv"
	"strings"
	"time"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
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

const (
	interpolationNamespaceParams  = "params"
	interpolationNamespaceSecrets = "secrets"
)

var (
	// interpolationTokenRegexp matches one ${{ ... }} placeholder and captures its inner text.
	// The inner text is a single term ("params.foo") or a "|"-separated fallback chain
	// ("params.foo | secrets.foo") — resolveInterpolationToken interprets it.
	interpolationTokenRegexp = regexp.MustCompile(`\$\{\{([^{}]+)\}\}`)
	// interpolationTermWordRegexp is what each side of a term's "." must look like: a bare word.
	interpolationTermWordRegexp = regexp.MustCompile(`^\w+$`)
)

const secretValuePrefix = "__secrets."

type HttpExecutor struct {
	base   *http.Client
	byTool map[artel_q.McpToolName]*http.Client
	byMcp  map[artel_q.McpName]*http.Client
}

// NewHttpExecutor builds a HttpExecutor whose outbound requests all pass through the shared
// logging transport, plus any per-tool / per-mcp middleware wound from opts.
func NewHttpExecutor(opts ...HttpExecutorOption) *HttpExecutor {
	return newHttpExecutor(newLoggingTransport(), opts...)
}

// newHttpExecutor is the construction seam the white-box middleware test uses to inject a stub
// base transport in place of the real logging transport. It winds every per-tool and per-mcp
// http.Client once, here — nothing about middleware selection happens per request.
func newHttpExecutor(base http.RoundTripper, opts ...HttpExecutorOption) *HttpExecutor {
	cfg := &httpExecutorConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	baseClient, byTool, byMcp := buildHttpClients(cfg, base, httpExecutorTimeout)

	return &HttpExecutor{
		base:   baseClient,
		byTool: byTool,
		byMcp:  byMcp,
	}
}

// clientFor picks the pre-wound http.Client for a tool: its own per-tool client if one was
// registered, else its mcp's per-mcp client, else the shared base client.
func (e *HttpExecutor) clientFor(ident ToolIdent) *http.Client {
	tool := artel_q.McpToolName(ident.McpName + "." + ident.ToolName)

	if c, ok := e.byTool[tool]; ok {
		return c
	}

	if c, ok := e.byMcp[artel_q.McpName(ident.McpName)]; ok {
		return c
	}

	return e.base
}

func (e *HttpExecutor) Execute(ctx context.Context, ident ToolIdent, action domain.ToolAction,
	secrets map[string]interface{}, params map[string]interface{},
) (string, error) {
	if action.Http == nil {
		return "", user_errors.McpActionMissing
	}

	client := e.clientFor(ident)

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

	resp, err := client.Do(req)
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

// resolveActionValue resolves a single header/query/body value. If the value is exactly
// "__secrets.<field>", it is replaced wholesale with secrets[field] (hard error if missing).
// Otherwise the value is run through ${{params.*}} / ${{secrets.*}} interpolation, including
// "|" fallback chains such as ${{params.chat_id | secrets.chat_id}}.
func resolveActionValue(value string, params map[string]interface{}, secrets map[string]interface{}) (string, error) {
	field, ok := secretField(value)
	if !ok {
		resolved := interpolateTokens(value, params, secrets,
			interpolationNamespaceParams, interpolationNamespaceSecrets)

		return resolved, nil
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
	return interpolateTokens(value, params, nil, interpolationNamespaceParams)
}

// interpolateSecrets replaces ${{secrets.field}} tokens within a string (e.g. a URL's
// host) with the corresponding value from secrets. Unlike resolveActionValue's
// __secrets. convention, this allows secrets to be embedded inside a larger string
// rather than only resolving a header/query value wholesale.
func interpolateSecrets(value string, secrets map[string]interface{}) string {
	return interpolateTokens(value, nil, secrets, interpolationNamespaceSecrets)
}

// interpolateTokens rewrites every ${{...}} placeholder in value via resolveInterpolationToken,
// restricted to the namespaces listed in allowed. A token that names a namespace outside
// allowed, or whose inner text is not a well-formed term chain, is left verbatim so a later
// pass (e.g. interpolateSecrets after interpolateParams on a URL) can still handle it.
func interpolateTokens(
	value string,
	params map[string]interface{},
	secrets map[string]interface{},
	allowed ...string,
) string {
	return interpolationTokenRegexp.ReplaceAllStringFunc(value, func(match string) string {
		submatches := interpolationTokenRegexp.FindStringSubmatch(match)
		if len(submatches) < regexSubmatchWithGroup {
			return match
		}

		resolved, ok := resolveInterpolationToken(submatches[1], params, secrets, allowed...)
		if !ok {
			return match
		}

		return resolved
	})
}

// resolveInterpolationToken resolves one placeholder's inner text. inner is a single term
// ("params.foo") or a "|"-separated fallback chain ("params.foo | secrets.bar"): terms are
// tried left to right and the first that names a present key wins; every term missing yields
// "". A term is "<namespace>.<key>", namespace being params or secrets, with whitespace around
// "|" and "." tolerated. ok is false when inner is not a well-formed chain, or a term names a
// namespace not in allowed — the caller then leaves the raw token in place.
func resolveInterpolationToken(
	inner string,
	params map[string]interface{},
	secrets map[string]interface{},
	allowed ...string,
) (string, bool) {
	for _, term := range strings.Split(inner, "|") {
		namespace, key, ok := splitInterpolationTerm(term)
		if !ok {
			return "", false
		}

		if !slices.Contains(allowed, namespace) {
			return "", false
		}

		var source map[string]interface{}

		switch namespace {
		case interpolationNamespaceParams:
			source = params
		case interpolationNamespaceSecrets:
			source = secrets
		}

		resolved, found := source[key]
		if found {
			return stringifyValue(resolved), true
		}
	}

	return "", true
}

// splitInterpolationTerm splits "<namespace>.<key>" (optionally whitespace-padded) into its two
// identifiers. ok is false when there is no ".", or either side is not a bare word.
func splitInterpolationTerm(term string) (namespace, key string, ok bool) {
	dot := strings.IndexByte(term, '.')
	if dot < 0 {
		return "", "", false
	}

	namespace = strings.TrimSpace(term[:dot])
	key = strings.TrimSpace(term[dot+1:])

	if !interpolationTermWordRegexp.MatchString(namespace) {
		return "", "", false
	}

	if !interpolationTermWordRegexp.MatchString(key) {
		return "", "", false
	}

	return namespace, key, true
}
