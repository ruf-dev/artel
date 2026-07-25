package tract

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/anthropic"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// LlmCallRequest is a single-turn "Call LLM" step invocation — model/prompt/system prompt are
// already-rendered strings (template resolution happens in engine.go before Call is invoked).
type LlmCallRequest struct {
	Model        string
	Prompt       string
	SystemPrompt string
	MaxTokens    int
}

// LlmCallUsage mirrors internal/clients/anthropic.Usage — kept as its own type here (rather
// than reusing the anthropic package's directly) so the tract package's public step-output
// shape doesn't leak a client-package type, matching how ToolExecutor returns plain
// string/map results rather than MoM-client types.
type LlmCallUsage struct {
	InputTokens              int64
	OutputTokens             int64
	CacheCreationInputTokens int64
	CacheReadInputTokens     int64
}

// LlmCallResult is the text of the model's reply, the model actually used, and token usage.
type LlmCallResult struct {
	Text  string
	Model string
	Usage LlmCallUsage
}

// LlmExecutor is the narrow slice of LLM-call capability the engine needs for llm_call steps.
// Defined here — not imported from internal/clients/anthropic directly by the engine — to
// mirror ToolExecutor's narrow-interface pattern in toolexecutor.go, keeping the engine
// decoupled from the concrete client/connection wiring.
//
// Unlike ExecuteMomTool (whose contract requires the caller/engine to verify connection
// ownership itself before calling), Call re-verifies connection ownership on every invocation —
// callers must not skip a redundant check, none is needed. This mirrors 04_tract_llm_step.md's
// documented contract exactly (the opposite of ToolExecutor.ExecuteMomTool's contract).
type LlmExecutor interface {
	Call(ctx context.Context, userUuid uuid.UUID, connectionUuid uuid.UUID, req LlmCallRequest) (LlmCallResult, error)
}

// llmExecutorAdapter implements LlmExecutor over the tract package's own
// repository.ExternalConnectionRepo and internal/clients/anthropic.Client. No dependency on the
// externalconnections service is needed: the pg repo implementation already decrypts
// credentials_enc before returning domain.ExternalConnection.CredentialsJSON (see
// internal/repository/pg/repos/externalconnections), so GetByID alone is enough to recover the
// plaintext Anthropic API key.
type llmExecutorAdapter struct {
	externalConns repository.ExternalConnectionRepo
}

// NewLlmExecutor wires the shared external connections repo into an LlmExecutor for the tract
// engine — constructed in internal/app/custom.go alongside NewToolExecutor, before tract.New.
func NewLlmExecutor(externalConns repository.ExternalConnectionRepo) LlmExecutor {
	adapter := &llmExecutorAdapter{externalConns: externalConns}

	return adapter
}

func (a *llmExecutorAdapter) Call(
	ctx context.Context,
	userUuid uuid.UUID,
	connectionUuid uuid.UUID,
	req LlmCallRequest,
) (LlmCallResult, error) {
	conn, err := a.externalConns.GetByID(ctx, connectionUuid)
	if err != nil {
		return LlmCallResult{}, rerrors.Wrap(err, "error getting llm connection")
	}

	if conn.UserUuid != userUuid {
		return LlmCallResult{}, rerrors.Wrap(user_errors.TractConnectionNotOwned, connectionUuid.String())
	}

	if conn.Provider != domain.ProviderAnthropic {
		return LlmCallResult{}, rerrors.Wrap(user_errors.TractLlmConnectionProviderMismatch, connectionUuid.String())
	}

	var creds domain.AnthropicKeyCredentials

	err = json.Unmarshal(conn.CredentialsJSON, &creds)
	if err != nil {
		return LlmCallResult{}, rerrors.Wrap(err, "error unmarshaling anthropic credentials")
	}

	client := anthropic.New(creds.ApiKey, creds.BaseUrl)

	completeReq := anthropic.CompleteRequest{
		Model:        req.Model,
		SystemPrompt: req.SystemPrompt,
		Prompt:       req.Prompt,
		MaxTokens:    int64(req.MaxTokens),
	}

	result, err := client.Complete(ctx, completeReq)
	if err != nil {
		return LlmCallResult{}, rerrors.Wrap(user_errors.TractLlmCallFailed, err.Error())
	}

	callResult := LlmCallResult{
		Text:  result.Text,
		Model: req.Model,
		Usage: LlmCallUsage{
			InputTokens:              result.Usage.InputTokens,
			OutputTokens:             result.Usage.OutputTokens,
			CacheCreationInputTokens: result.Usage.CacheCreationInputTokens,
			CacheReadInputTokens:     result.Usage.CacheReadInputTokens,
		},
	}

	return callResult, nil
}

var _ LlmExecutor = (*llmExecutorAdapter)(nil)
