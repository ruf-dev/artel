package external_connections_api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type exchangeRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func (e *ExternalConnectionsImpl) HandleGoogleExchange(w http.ResponseWriter, r *http.Request) {
	var req exchangeRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	_, err = e.svc.HandleGoogleOAuthCallback(r.Context(), req.Code, req.State)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("google oauth exchange failed")
		http.Error(w, "exchange failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
