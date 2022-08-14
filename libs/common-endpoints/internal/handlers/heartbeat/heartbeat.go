package heartbeat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers"
)

const heartbeatPattern = "/heartbeat"

func Add(router handlers.Router) {
	router.HandleFunc(heartbeatPattern, newHeartbeatHandler())
}

func newHeartbeatHandler() func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, _ *http.Request) {
		err := json.NewEncoder(responseWriter).Encode(newOkResponse())
		if err != nil {
			log.Println("could not encode heartbeat", err)
		}

		responseWriter.Header().Add("content-type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
	}
}

func newOkResponse() handlers.Result {
	return handlers.Result{
		Success: true,
	}
}
