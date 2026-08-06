// Package openai is a thin wrapper around the official OpenAI SDK
// (github.com/openai/openai-go), used to validate a user-supplied API key
// (via ListModels, a zero-token metadata call) and to run single-turn, non-streaming
// completions (via Complete) for the tract engine's "Call LLM" step.
//
// No streaming, no tool use, no multi-turn conversation — deliberately out of scope
// for this first pass. See docs/byok/04_tract_llm_step.md for the full design rationale.
package openai

import (
	"context"
	"errors"

	sdk "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"go.redsock.ru/rerrors"
)

// defaultMaxTokens is used for Complete when the caller doesn't set CompleteRequest.MaxTokens.
const defaultMaxTokens = 4096

// pingMaxTokens is used by Ping — a completion is only ever used to confirm the key/model
// authenticate, so the reply itself is discarded and kept as cheap as the API allows.
const pingMaxTokens = 1

// Client wraps the OpenAI SDK client with the narrow surface the connection service and
// tract engine actually need.
type Client struct {
	sdk sdk.Client
}

// New constructs a Client authenticated with apiKey. If baseUrl is non-empty, it overrides the
// SDK's default API host (e.g. for a proxy, regional endpoint, or OpenAI-compatible provider).
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

// ModelInfo is the subset of the Models-List response fields relevant to Artel. It is
// intentionally thinner than Anthropic's ModelInfo: OpenAI's GET /v1/models response only
// carries an id and an owning org — no display name or token-limit metadata is returned.
type ModelInfo struct {
	Id      string
	OwnedBy string
}

// ListModels wraps the SDK's models-list call (GET /v1/models). It is a metadata endpoint —
// no completion tokens billed — so it doubles as the key-validation call: an auth failure
// surfaces here as an error before any completion is ever attempted.
func (c *Client) ListModels(ctx context.Context) ([]ModelInfo, error) {
	page, err := c.sdk.Models.List(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing openai models")
	}

	models := make([]ModelInfo, 0, len(page.Data))

	for _, m := range page.Data {
		models = append(models, ModelInfo{
			Id:      m.ID,
			OwnedBy: m.OwnedBy,
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
// providers whose OpenAI-compatible endpoint doesn't implement GET /v1/models (ListModels) —
// common for third-party proxies that only mirror the Chat Completions API surface.
func (c *Client) Ping(ctx context.Context, model string) error {
	req := CompleteRequest{
		Model:     model,
		Prompt:    "hi",
		MaxTokens: pingMaxTokens,
	}

	_, err := c.Complete(ctx, req)
	if err != nil {
		return rerrors.Wrap(err, "error pinging openai-compatible endpoint")
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

// Usage mirrors the Chat Completions API response's usage object. Unlike Anthropic's Usage,
// there are no cache-related fields: OpenAI's prompt-caching accounting
// (prompt_tokens_details.cached_tokens) has no analogue in this narrow wrapper and is
// intentionally left unmapped.
type Usage struct {
	InputTokens  int64
	OutputTokens int64
}

// CompleteResult is the text of the model's reply plus its token usage.
type CompleteResult struct {
	Text  string
	Usage Usage
}

// Complete performs a single-turn, non-streaming call to the Chat Completions API.
func (c *Client) Complete(ctx context.Context, req CompleteRequest) (CompleteResult, error) {
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	messages := make([]sdk.ChatCompletionMessageParamUnion, 0, 2)

	if req.SystemPrompt != "" {
		messages = append(messages, sdk.SystemMessage(req.SystemPrompt))
	}

	messages = append(messages, sdk.UserMessage(req.Prompt))

	params := sdk.ChatCompletionNewParams{
		Model:    req.Model,
		Messages: messages,
		// MaxTokens is deprecated by OpenAI in favor of MaxCompletionTokens as of the current
		// Chat Completions API version; MaxTokens is additionally rejected outright by
		// o-series reasoning models.
		MaxCompletionTokens: sdk.Int(maxTokens),
	}

	resp, err := c.sdk.Chat.Completions.New(ctx, params)
	if err != nil {
		return CompleteResult{}, rerrors.Wrap(err, "error creating openai chat completion")
	}

	result := CompleteResult{
		Text: extractText(resp),
		Usage: Usage{
			InputTokens:  resp.Usage.PromptTokens,
			OutputTokens: resp.Usage.CompletionTokens,
		},
	}

	return result, nil
}

// extractText returns the first choice's message content. A single-turn completion is
// expected to contain exactly one choice (N is never set above its default of 1).
func extractText(resp *sdk.ChatCompletion) string {
	if len(resp.Choices) == 0 {
		return ""
	}

	return resp.Choices[0].Message.Content
}
