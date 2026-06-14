package external_connections_api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	callbackSuccessRedirect = "/settings/integrations?status=connected"
	callbackErrorRedirect   = "/settings/integrations?status=error"
)

func (e *ExternalConnectionsImpl) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	_, err := e.svc.HandleGoogleOAuthCallback(r.Context(), code, state)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("google oauth callback failed")
		http.Redirect(w, r, callbackErrorRedirect, http.StatusFound)
		return
	}

	http.Redirect(w, r, callbackSuccessRedirect, http.StatusFound)
}
