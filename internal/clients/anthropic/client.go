// Package anthropic is a thin wrapper around the official Anthropic SDK
// (github.com/anthropics/anthropic-sdk-go), used to validate a user-supplied API key
// (via ListModels, a zero-token metadata call) and to run single-turn, non-streaming
// completions (via Complete) for the tract engine's "Call LLM" step.
//
// No streaming, no tool use, no multi-turn conversation — deliberately out of scope
// for this first pass. See docs/byok/04_tract_llm_step.md for the full design rationale.
package anthropic

import (
	"context"
	"errors"

	sdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"go.redsock.ru/rerrors"
)

// defaultMaxTokens is used for Complete when the caller doesn't set CompleteRequest.MaxTokens.
const defaultMaxTokens = 4096

// pingMaxTokens is used by Ping — a completion is only ever used to confirm the key/model
// authenticate, so the reply itself is discarded and kept as cheap as the API allows.
const pingMaxTokens = 1

// Client wraps the Anthropic SDK client with the narrow surface the connection service and
// tract engine actually need.
type Client struct {
	sdk sdk.Client
}

// New constructs a Client authenticated with apiKey. If baseUrl is non-empty, it overrides the
// SDK's default API host (e.g. for a proxy or regional endpoint).
func New(apiKey, baseUrl string) *Client {
	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	if baseUrl != "" {
		opts = append(opts, option.WithBaseURL(baseUrl))
	}

	return &Client{
		sdk: sdk.NewClient(opts...),
	}
}

// ModelInfo is the subset of the Models-List response fields relevant to Artel: a model
// catalog entry used both to confirm a key is valid and to populate a connection's
// available-models cache / the Call LLM step's model picker.
type ModelInfo struct {
	Id             string
	DisplayName    string
	MaxInputTokens int64
	MaxTokens      int64
}

// ListModels wraps the SDK's models-list call (GET /v1/models). It is a metadata endpoint —
// no beta header, zero tokens billed — so it doubles as the key-validation call: an auth
// failure surfaces here as an error before any completion is ever attempted.
func (c *Client) ListModels(ctx context.Context) ([]ModelInfo, error) {
	page, err := c.sdk.Models.List(ctx, sdk.ModelListParams{})
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing anthropic models")
	}

	models := make([]ModelInfo, 0, len(page.Data))

	for _, m := range page.Data {
		models = append(models, ModelInfo{
			Id:             m.ID,
			DisplayName:    m.DisplayName,
			MaxInputTokens: m.MaxInputTokens,
			MaxTokens:      m.MaxTokens,
		})
	}

	return models, nil
}

// StatusCode extracts the HTTP status code from an error returned by this client, when the
// error originated from an actual API response (as opposed to a network-level failure like a
// timeout or DNS error, which the SDK does not wrap in *sdk.Error). ok is false for those
// network-level cases.
func StatusCode(err error) (code int, ok bool) {
	var apiErr *sdk.Error

	ok = errors.As(err, &apiErr)
	if !ok {
		return 0, false
	}

	return apiErr.StatusCode, true
}

// Ping confirms apiKey authenticates and model is usable by sending the cheapest possible
// completion request, discarding the reply. It exists as a fallback key-validation path for
// providers whose Anthropic-compatible endpoint doesn't implement GET /v1/models (ListModels) —
// common for third-party proxies that only mirror the Messages API surface.
func (c *Client) Ping(ctx context.Context, model string) error {
	req := CompleteRequest{
		Model:     model,
		Prompt:    "hi",
		MaxTokens: pingMaxTokens,
	}

	_, err := c.Complete(ctx, req)
	if err != nil {
		return rerrors.Wrap(err, "error pinging anthropic-compatible endpoint")
	}

	return nil
}

// CompleteRequest is a single-turn, non-streaming completion request: one user message
// (Prompt), an optional SystemPrompt, and a model. No tool use, no thinking, no history.
type CompleteRequest struct {
	Model        string
	SystemPrompt string
	Prompt       string
	MaxTokens    int64
}

// Usage mirrors the Messages API response's usage object.
type Usage struct {
	InputTokens              int64
	OutputTokens             int64
	CacheCreationInputTokens int64
	CacheReadInputTokens     int64
}

// CompleteResult is the text of the model's reply plus its token usage.
type CompleteResult struct {
	Text  string
	Usage Usage
}

// Complete performs a single-turn, non-streaming call to the Messages API.
func (c *Client) Complete(ctx context.Context, req CompleteRequest) (CompleteResult, error) {
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	params := sdk.MessageNewParams{
		Model:     req.Model,
		MaxTokens: maxTokens,
		Messages: []sdk.MessageParam{
			sdk.NewUserMessage(sdk.NewTextBlock(req.Prompt)),
		},
	}

	if req.SystemPrompt != "" {
		params.System = []sdk.TextBlockParam{
			{Text: req.SystemPrompt},
		}
	}

	resp, err := c.sdk.Messages.New(ctx, params)
	if err != nil {
		return CompleteResult{}, rerrors.Wrap(err, "error creating anthropic message")
	}

	result := CompleteResult{
		Text: extractText(resp),
		Usage: Usage{
			InputTokens:              resp.Usage.InputTokens,
			OutputTokens:             resp.Usage.OutputTokens,
			CacheCreationInputTokens: resp.Usage.CacheCreationInputTokens,
			CacheReadInputTokens:     resp.Usage.CacheReadInputTokens,
		},
	}

	return result, nil
}

// extractText concatenates every text block in the response content, in order. A
// single-turn completion is expected to contain exactly one text block, but concatenating
// is harmless if the model ever returns more than one.
func extractText(resp *sdk.Message) string {
	var text string

	for _, block := range resp.Content {
		switch variant := block.AsAny().(type) {
		case sdk.TextBlock:
			text += variant.Text
		}
	}

	return text
}
