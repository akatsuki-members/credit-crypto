package endpoints

import (
	"net/http"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/health"
	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/heartbeat"
	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/info"
)

// HealthChecker function type to check service health.
type HealthChecker func() Report

// Info contains service information data.
type Info struct {
	Name    string `json:"name"`
	Commit  string `json:"commit"` // commit hash
	Version string `json:"version"`
}

// Item contains information about components and their status.
type Item struct {
	Name    string `json:"name"`    // integration or component name
	Healthy bool   `json:"healthy"` // health status of the component.
}

// Contains health report.
type Report struct {
	Healthy bool   `json:"healthy"`          // The service is healthy
	Data    []Item `json:"report,omitempty"` // health report
}

type Handler struct {
	router       *http.ServeMux
	hasEndpoints bool
}

func New(router *http.ServeMux) *Handler {
	newHandler := Handler{
		router: router,
	}

	return &newHandler
}

// WithHeartbeat add heartbeat endpoint to the service.
func (h *Handler) WithHeartbeat() *Handler {
	h.hasEndpoints = true
	heartbeat.Add(h.router)

	return h
}

// WithHealth add health endpoint to the service.
func (h *Handler) WithHealth(checker HealthChecker) *Handler {
	h.hasEndpoints = true
	health.Add(h.router, h.newHealthChecker(checker))

	return h
}

func (h *Handler) WithInfo(data Info) *Handler {
	h.hasEndpoints = true
	infoReport := info.Report{
		Name:    data.Name,
		Commit:  data.Commit,
		Version: data.Version,
	}
	info.Add(h.router, infoReport)

	return h
}

func (h *Handler) newHealthChecker(checker HealthChecker) health.Checker {
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
