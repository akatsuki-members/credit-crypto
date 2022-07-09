package endpoints

import (
	"net/http"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/heartbeat"
)

type handler struct {
	router       *http.ServeMux
	hasEndpoints bool
}

func New(router *http.ServeMux) *handler {
	newHandler := handler{
		router: router,
	}
	return &newHandler
}

func (h *handler) WithHeartbeat() *handler {
	h.hasEndpoints = true
	heartbeat.Add(h.router)
	return h
}
