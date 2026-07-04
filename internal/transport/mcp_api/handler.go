package mcp_api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/requesthost"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/utils"

	"go.opentelemetry.io/otel/trace"
	"go.redsock.ru/rerrors"
)

const jsonrpcVersion = "2.0"

type rpcRequest struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type rpcResponse struct {
	Jsonrpc string    `json:"jsonrpc"`
	Id      any       `json:"id,omitempty"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceId string `json:"trace_id,omitempty"`
	Details string `json:"details,omitempty"`
}

type McpHandler struct {
	mcpSvc service.McpService
	momSvc service.MomService
}

func NewMcpHandler(mcpSvc service.McpService, momSvc service.MomService) *McpHandler {
	return &McpHandler{
		mcpSvc: mcpSvc,
		momSvc: momSvc,
	}
}

func (h *McpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		h.serveSSE(writer, request)

		return
	}

	writer.Header().Set("Content-Type", "application/json")

	ctx := request.Context()
	ctx = requesthost.WithHost(ctx, publicHost(request))

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	token := extractBearerToken(request)
	if token == "" {
		writer.Header().Set("WWW-Authenticate", `Bearer realm="`+publicHost(request)+`"`)
		writer.WriteHeader(http.StatusUnauthorized)

		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		writeErrorResponse(ctx, writer, nil, -32700, "parse error", err)

		return
	}
	defer utils.CloseWithLog(request.Body, "error closing request body (mcp api)")

	var req rpcRequest

	err = json.Unmarshal(body, &req)
	if err != nil {
		writeErrorResponse(ctx, writer, nil, -32700, "parse error", err)

		return
	}

	keyCtx, err := h.mcpSvc.ResolveKey(ctx, token)
	if err != nil {
		if h.isNotification(req) {
			writer.WriteHeader(http.StatusOK)

			return
		}

		writeErrorResponse(ctx, writer, req.Id, -32001, "unauthorized", nil)

		return
	}

	if h.isNotification(req) {
		switch req.Method {
		case "notifications/initialized":
			writer.WriteHeader(http.StatusOK)

			return
		default:
			writer.WriteHeader(http.StatusOK)

			return
		}
	}

	ctx = contextWithKeyCtx(ctx, keyCtx)

	h.dispatchMethod(ctx, writer, req)
}

func (h *McpHandler) serveSSE(writer http.ResponseWriter, request *http.Request) {
	if !strings.Contains(request.Header.Get("Accept"), "text/event-stream") {
		writer.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	token := extractBearerToken(request)
	if token == "" {
		writer.Header().Set("WWW-Authenticate", `Bearer realm="`+publicHost(request)+`"`)
		writer.WriteHeader(http.StatusUnauthorized)

		return
	}

	_, err := h.mcpSvc.ResolveKey(request.Context(), token)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)

		return
	}

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.WriteHeader(http.StatusOK)

	if f, ok := writer.(http.Flusher); ok {
		f.Flush()
	}

	<-request.Context().Done()
}

func (h *McpHandler) dispatchMethod(ctx context.Context, w http.ResponseWriter, req rpcRequest) {
	switch req.Method {
	case "initialize":
		h.handleInitialize(ctx, w, req)
	case "tools/list":
		h.handleToolsList(ctx, w, req)
	case "tools/call":
		h.handleToolsCall(ctx, w, req)
	default:
		writeErrorResponse(ctx, w, req.Id, -32601, "method not found", nil)
	}
}

func (h *McpHandler) handleInitialize(ctx context.Context, w http.ResponseWriter, req rpcRequest) {
	w.Header().Set("Mcp-Session-Id", uuid.New().String())

	result := map[string]any{
		"protocolVersion": "2025-03-26",
		"capabilities": map[string]any{
			"tools": map[string]any{},
		},
		"serverInfo": map[string]string{
			"name":    "artel-mcp",
			"version": "0.1.0",
		},
	}

	resp := rpcResponse{
		Jsonrpc: jsonrpcVersion,
		Id:      req.Id,
		Result:  result,
	}

	err := writeJsonResponse(w, resp)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to write initialize response")
		writeErrorResponse(ctx, w, req.Id, -32603, "internal error", err)
	}
}

func (h *McpHandler) handleToolsList(ctx context.Context, writer http.ResponseWriter, req rpcRequest) {
	keyCtx, ok := ctx.Value(keyContextKey).(domain.McpKeyContext)
	if !ok {
		log.Error().Msg("mcp: keyContext missing from context in tools/list")
		writeErrorResponse(ctx, writer, req.Id, -32603, "internal error", nil)

		return
	}

	builtinTools, err := h.mcpSvc.ListTools(ctx)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to list builtin tools")
		writeErrorResponse(ctx, writer, req.Id, -32603, "internal error", err)

		return
	}

	tools := make([]ToolDef, 0, len(builtinTools))
	for _, t := range builtinTools {
		tools = append(tools, momToolToToolDef(t))
	}

	momTools, err := h.momSvc.ListToolsForKey(ctx, keyCtx.KeyUuid)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to list mom tools")
	} else {
		for _, t := range momTools {
			tools = append(tools, momToolToToolDef(t))
		}
	}

	result := map[string]any{
		"tools": tools,
	}

	resp := rpcResponse{
		Jsonrpc: jsonrpcVersion,
		Id:      req.Id,
		Result:  result,
	}

	err = writeJsonResponse(writer, resp)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to write tools/list response")
		writeErrorResponse(ctx, writer, req.Id, -32603, "internal error", err)
	}
}

func (h *McpHandler) handleToolsCall(ctx context.Context, writer http.ResponseWriter, req rpcRequest) {
	var callReq struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	err := json.Unmarshal(req.Params, &callReq)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to unmarshal tools/call params")
		writeErrorResponse(ctx, writer, req.Id, -32602, "invalid params", err)

		return
	}

	keyCtx, ok := ctx.Value(keyContextKey).(domain.McpKeyContext)
	if !ok {
		log.Error().Msg("mcp: keyContext missing from context in tools/call")
		writeErrorResponse(ctx, writer, req.Id, -32603, "internal error", nil)

		return
	}

	var result interface{}

	if h.mcpSvc.IsBuiltinTool(callReq.Name) {
		var execRes domain.ToolExecResult
		execRes, err = h.mcpSvc.ExecuteTool(ctx, keyCtx, callReq.Name, callReq.Arguments)
		if err != nil {
			log.Error().Err(err).Str("tool", callReq.Name).Msg("mcp: tool execution failed")
			writeErrorResponse(ctx, writer, req.Id, -32603, "tool execution failed", err)

			return
		}

		result = toolResultFromExec(execRes)
	} else {
		var text string
		text, err = h.momSvc.ExecuteToolForKey(ctx, keyCtx.KeyUuid, callReq.Name, callReq.Arguments)
		if err != nil {
			log.Error().Err(err).Str("tool", callReq.Name).Msg("mcp: mom tool execution failed")
			writeErrorResponse(ctx, writer, req.Id, -32603, "tool execution failed", err)

			return
		}

		result = textResult(text)
	}

	resp := rpcResponse{
		Jsonrpc: jsonrpcVersion,
		Id:      req.Id,
		Result:  result,
	}

	err = writeJsonResponse(writer, resp)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to write tools/call response")
		writeErrorResponse(ctx, writer, req.Id, -32603, "internal error", err)
	}
}

func publicHost(r *http.Request) string {
	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		return h
	}

	return r.Host
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func writeJsonResponse(w http.ResponseWriter, response rpcResponse) error {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(response)
	if err != nil {
		return rerrors.Wrap(err, "failed to marshal response")
	}

	_, err = w.Write(data)
	if err != nil {
		return rerrors.Wrap(err, "failed to write response")
	}

	return nil
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, id any, code int, message string, cause error) {
	e := &rpcError{
		Code:    code,
		Message: message,
	}
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		e.TraceId = span.SpanContext().TraceID().String()
	}

	if cause != nil {
		e.Details = cause.Error()
	}

	resp := rpcResponse{
		Jsonrpc: jsonrpcVersion,
		Id:      id,
		Error:   e,
	}

	err := writeJsonResponse(w, resp)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to write error response")
	}
}

func (h *McpHandler) isNotification(req rpcRequest) bool {
	return req.Id == nil
}
