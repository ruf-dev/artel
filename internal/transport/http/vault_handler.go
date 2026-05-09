package http

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service"
)

type VaultHandler struct {
	vs service.VaultService
}

func NewVaultHandler(vs service.VaultService) *VaultHandler {
	return &VaultHandler{vs: vs}
}

type createVaultRequest struct {
	Name string `json:"name"`
}

func (h *VaultHandler) CreateVault(w http.ResponseWriter, r *http.Request) {
	var req createVaultRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.vs.CreateVault(r.Context(), req.Name)
	if err != nil {
		log.Error().Err(err).Msg("failed to create vault")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
