package endpoints

import (
	"net/http"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/health"
	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/heartbeat"
)

// HealthChecker function type to check service health
type HealthChecker func() Report

// Item contains information about components and their status.
type Item struct {
	Name    string `json:"name"`    // integration or component name
	Healthy bool   `json:"healthy"` // health status of the component.
}

// Contains health report
type Report struct {
	Healthy bool   `json:"healthy"`          // The service is healthy
	Data    []Item `json:"report,omitempty"` // health report
}

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

func (h *handler) WithHealth(checker HealthChecker) *handler {
	h.hasEndpoints = true
	health.Add(h.router, h.newHealthChecker(checker))
	return h
}

func (h *handler) newHealthChecker(checker HealthChecker) health.HealthChecker {
	return func() health.Report {
		result := checker()
		report := make([]health.Item, len(result.Data))
		for idx, v := range result.Data {
			report[idx] = health.Item{
				Name:    v.Name,
				Healthy: v.Healthy,
			}
		}
		return health.Report{
			Healthy: result.Healthy,
			Data:    report,
		}
	}
}
