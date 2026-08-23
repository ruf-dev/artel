// Package openrouter is a thin wrapper around OpenRouter's own account-management REST API
// (https://openrouter.ai/api/v1), used to fetch the caller's key usage/limit/balance info
// (GetKeyInfo).
//
// This is deliberately NOT the OpenAI-compatible Chat Completions surface OpenRouter also exposes
// (see internal/clients/openai, reused for BYOK chat completions against OpenRouter via
// AddOpenAIConnection/CheckOpenAIConnection) — GetKeyInfo calls OpenRouter's own GET /key
// endpoint, which has no OpenAI-compatible equivalent and no official Go SDK, hence the plain
// net/http client here instead of an SDK wrapper.
package openrouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/utils"
)

// defaultBaseURL is OpenRouter's own account-management API host. GetKeyInfo always calls this
// directly — never a connection's configured BaseUrl override — because a BYOK connection's base
// URL may point at a proxy that doesn't implement OpenRouter's own /key endpoint.
const defaultBaseURL = "https://openrouter.ai/api/v1"

// Client wraps OpenRouter's own account-management REST API with the narrow surface the
// connections service needs (currently just GetKeyInfo).
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// New constructs a Client authenticated with apiKey. GetKeyInfo always targets OpenRouter's real
// account-management host (defaultBaseURL) — there is no baseUrl parameter to override it with,
// unlike internal/clients/openai, because a BYOK connection's base URL may point at a proxy that
// doesn't implement OpenRouter's own /key endpoint.
func New(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
	}
}

// KeyInfo is the subset of GET /key's response fields relevant to Artel's usage/balance display.
type KeyInfo struct {
	Label              string
	Limit              float64
	LimitUnlimited     bool // true when the response's "limit" field is JSON null (no cap set)
	LimitRemaining     float64
	IncludeByokInLimit bool
	Usage              float64
	UsageDaily         float64
	UsageWeekly        float64
	UsageMonthly       float64
	ByokUsage          float64
	ByokUsageDaily     float64
	ByokUsageWeekly    float64
	ByokUsageMonthly   float64
	IsFreeTier         bool
}

// keyInfoResponse mirrors GET /key's JSON body. Limit and LimitRemaining are pointers because
// OpenRouter returns JSON null for "limit" when the key has no spending cap.
type keyInfoResponse struct {
	Data struct {
		Label              string   `json:"label"`
		Limit              *float64 `json:"limit"`
		LimitReset         *string  `json:"limit_reset"`
		LimitRemaining     *float64 `json:"limit_remaining"`
		IncludeByokInLimit bool     `json:"include_byok_in_limit"`
		Usage              float64  `json:"usage"`
		UsageDaily         float64  `json:"usage_daily"`
		UsageWeekly        float64  `json:"usage_weekly"`
		UsageMonthly       float64  `json:"usage_monthly"`
		ByokUsage          float64  `json:"byok_usage"`
		ByokUsageDaily     float64  `json:"byok_usage_daily"`
		ByokUsageWeekly    float64  `json:"byok_usage_weekly"`
		ByokUsageMonthly   float64  `json:"byok_usage_monthly"`
		IsFreeTier         bool     `json:"is_free_tier"`
	} `json:"data"`
}

// APIError is returned when OpenRouter responds with a non-2xx status to a request made by this
// client. Use StatusCode to extract the code from an error returned by this client.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("openrouter api error: status %d: %s", e.StatusCode, e.Body)
}

// StatusCode extracts the HTTP status code from an error returned by this client, when the error
// originated from an actual API response (as opposed to a network-level failure like a timeout or
// DNS error, which this client does not wrap in *APIError). ok is false for those network-level
// cases.
func StatusCode(err error) (code int, ok bool) {
	var apiErr *APIError

	ok = errors.As(err, &apiErr)
	if !ok {
		return 0, false
	}

	return apiErr.StatusCode, true
}

// GetKeyInfo calls GET /key against OpenRouter's own account-management API — always
// c.baseURL (defaultBaseURL in production), never a caller-supplied override — returning the
// caller's current usage/limit/balance info. No caching, no retries: a plain synchronous call,
// live every time.
func (c *Client) GetKeyInfo(ctx context.Context) (KeyInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/key", nil)
	if err != nil {
		return KeyInfo{}, rerrors.Wrap(err, "error building openrouter key info request")
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return KeyInfo{}, rerrors.Wrap(err, "error calling openrouter key info endpoint")
	}

	defer utils.CloseWithLog(resp.Body, "openrouter key info response")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return KeyInfo{}, rerrors.Wrap(err, "error reading openrouter key info response")
	}

	if resp.StatusCode != http.StatusOK {
		apiErr := &APIError{StatusCode: resp.StatusCode, Body: string(body)}

		return KeyInfo{}, rerrors.Wrap(apiErr, "openrouter key info request failed")
	}

	var parsed keyInfoResponse

	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return KeyInfo{}, rerrors.Wrap(err, "error parsing openrouter key info response")
	}

	info := KeyInfo{
		Label:              parsed.Data.Label,
		LimitUnlimited:     parsed.Data.Limit == nil,
		IncludeByokInLimit: parsed.Data.IncludeByokInLimit,
		Usage:              parsed.Data.Usage,
		UsageDaily:         parsed.Data.UsageDaily,
		UsageWeekly:        parsed.Data.UsageWeekly,
		UsageMonthly:       parsed.Data.UsageMonthly,
		ByokUsage:          parsed.Data.ByokUsage,
		ByokUsageDaily:     parsed.Data.ByokUsageDaily,
		ByokUsageWeekly:    parsed.Data.ByokUsageWeekly,
		ByokUsageMonthly:   parsed.Data.ByokUsageMonthly,
		IsFreeTier:         parsed.Data.IsFreeTier,
	}

	if parsed.Data.Limit != nil {
		info.Limit = *parsed.Data.Limit
	}

	if parsed.Data.LimitRemaining != nil {
		info.LimitRemaining = *parsed.Data.LimitRemaining
	}

	return info, nil
}
