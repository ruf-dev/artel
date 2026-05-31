package mcp_api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service"
)

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
}

type McpHandler struct {
	mcpSvc   service.McpService
	emailSvc service.EmailService
}

func NewMcpHandler(mcpSvc service.McpService, emailSvc service.EmailService) *McpHandler {
	return &McpHandler{
		mcpSvc:   mcpSvc,
		emailSvc: emailSvc,
	}
}

func (h *McpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.serveSSE(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := extractBearerToken(r)
	if token == "" {
		w.Header().Set("WWW-Authenticate", `Bearer realm="`+publicHost(r)+`"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, nil, -32700, "parse error")
		return
	}
	defer r.Body.Close()

	var req rpcRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		writeErrorResponse(w, nil, -32700, "parse error")
		return
	}

	ctx := r.Context()

	keyCtx, err := h.mcpSvc.ResolveKey(ctx, token)
	if err != nil {
		if h.isNotification(req) {
			w.WriteHeader(http.StatusOK)
			return
		}
		writeErrorResponse(w, req.Id, -32001, "unauthorized")
		return
	}

	if h.isNotification(req) {
		switch req.Method {
		case "notifications/initialized":
			w.WriteHeader(http.StatusOK)
			return
		default:
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	ctx = contextWithKeyCtx(ctx, keyCtx)

	h.dispatchMethod(w, ctx, req)
}

func (h *McpHandler) serveSSE(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept"), "text/event-stream") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	token := extractBearerToken(r)
	if token == "" {
		w.Header().Set("WWW-Authenticate", `Bearer realm="`+publicHost(r)+`"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if _, err := h.mcpSvc.ResolveKey(r.Context(), token); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	<-r.Context().Done()
}

func (h *McpHandler) dispatchMethod(w http.ResponseWriter, ctx context.Context, req rpcRequest) {
	switch req.Method {
	case "initialize":
		h.handleInitialize(w, ctx, req)
	case "tools/list":
		h.handleToolsList(w, ctx, req)
	case "tools/call":
		h.handleToolsCall(w, ctx, req)
	default:
		writeErrorResponse(w, req.Id, -32601, "method not found")
	}
}

func (h *McpHandler) handleInitialize(w http.ResponseWriter, ctx context.Context, req rpcRequest) {
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
		Jsonrpc: "2.0",
		Id:      req.Id,
		Result:  result,
	}

	err := writeJsonResponse(w, resp)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to write initialize response")
		writeErrorResponse(w, req.Id, -32603, "internal error")
	}
}

func (h *McpHandler) handleToolsList(w http.ResponseWriter, ctx context.Context, req rpcRequest) {
	keyCtx, ok := ctx.Value(keyContextKey).(KeyContext)
	if !ok {
		log.Error().Msg("mcp: keyContext missing from context in tools/list")
		writeErrorResponse(w, req.Id, -32603, "internal error")
		return
	}

	tools := getToolDefinitions(keyCtx.HasEmails)

	result := map[string]any{
		"tools": tools,
	}

	resp := rpcResponse{
		Jsonrpc: "2.0",
		Id:      req.Id,
		Result:  result,
	}

	err := writeJsonResponse(w, resp)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to write tools/list response")
		writeErrorResponse(w, req.Id, -32603, "internal error")
	}
}

func (h *McpHandler) handleToolsCall(w http.ResponseWriter, ctx context.Context, req rpcRequest) {
	var callReq struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	err := json.Unmarshal(req.Params, &callReq)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to unmarshal tools/call params")
		writeErrorResponse(w, req.Id, -32602, "invalid params")
		return
	}

	keyCtx, ok := ctx.Value(keyContextKey).(KeyContext)
	if !ok {
		log.Error().Msg("mcp: keyContext missing from context in tools/call")
		writeErrorResponse(w, req.Id, -32603, "internal error")
		return
	}

	switch callReq.Name {
	case toolListEmailFolders, toolListEmailAccounts, toolListEmails, toolReadEmail, toolSendEmail:
		if !keyCtx.HasEmails {
			writeErrorResponse(w, req.Id, -32001, "permission denied")
			return
		}
	}

	result, err := dispatchToolCall(ctx, callReq.Name, callReq.Arguments, keyCtx, h.emailSvc)
	if err != nil {
		log.Error().Err(err).Str("tool", callReq.Name).Msg("mcp: tool execution failed")
		writeErrorResponse(w, req.Id, -32603, "tool execution failed")
		return
	}

	resp := rpcResponse{
		Jsonrpc: "2.0",
		Id:      req.Id,
		Result:  result,
	}

	err = writeJsonResponse(w, resp)
	if err != nil {
		log.Error().Err(err).Msg("mcp: failed to write tools/call response")
		writeErrorResponse(w, req.Id, -32603, "internal error")
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

func writeErrorResponse(w http.ResponseWriter, id any, code int, message string) error {
	resp := rpcResponse{
		Jsonrpc: "2.0",
		Id:      id,
		Error: &rpcError{
			Code:    code,
			Message: message,
		},
	}

	return writeJsonResponse(w, resp)
}

func (h *McpHandler) isNotification(req rpcRequest) bool {
	return req.Id == nil
}
