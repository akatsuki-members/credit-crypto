package health

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers"
)

const healthPattern = "/health"

// Checker function type to check service health.
type Checker func() Report

type Item struct {
	Name    string `json:"name"`    // integration or component name
	Healthy bool   `json:"healthy"` // health status of the component.
}

// Contains health report.
type Report struct {
	Healthy bool   `json:"healthy"`          // The service is healthy
	Data    []Item `json:"report,omitempty"` // health report
}

// Add use the given checker function to run the health check on your app.
func Add(router handlers.Router, checkHealthFunc Checker) {
	router.HandleFunc(healthPattern, newHealthHandler(checkHealthFunc))
}

func newHealthHandler(checkHealth Checker) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		healthStatus := checkHealth()
		res.WriteHeader(getHTTPStatus(healthStatus.Healthy))

		err := json.NewEncoder(res).Encode(newHealthResponse(healthStatus))
		if err != nil {
			log.Println("could not encode health response")
		}

		res.Header().Add("content-type", "application/json")
	}
}

func newHealthResponse(report Report) handlers.Result {
	return handlers.Result{
		Success: report.Healthy,
		Data:    report.Data,
	}
}

func getHTTPStatus(healthy bool) int {
	if healthy {
		return http.StatusOK
	}

	return http.StatusInternalServerError
}
