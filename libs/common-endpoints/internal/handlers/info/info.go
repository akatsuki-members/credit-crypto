package info

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers"
)

const infoPattern = "/info"

// Report contains metadata related to the service.
type Report struct {
	Name    string
	Commit  string
	Version string
}

func Add(router handlers.Router, infoData Report) {
	router.HandleFunc(infoPattern, newInfoHandler(infoData))
}

func newInfoHandler(infoData Report) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("content-type", "application/json")

		err := json.NewEncoder(rw).Encode(newResult(infoData))
		if err != nil {
			log.Println("info report cannot be generated", err.Error(), "info", infoData)
		}
	}
}

func newResult(data Report) handlers.Result {
	return handlers.Result{
		Success: true,
		Data:    data,
	}
}
