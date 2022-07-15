package heartbeat

import (
	"encoding/json"
	"net/http"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers"
)

const heartbeatPattern = "/heartbeat"

func Add(router handlers.Router) {
	router.HandleFunc(heartbeatPattern, newHeartbeatHandler())
}

func newHeartbeatHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(newOkResponse())
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func newOkResponse() handlers.Result {
	return handlers.Result{
		Success: true,
	}
}
