package endpoints_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers"
	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/pkg/endpoints"
	"github.com/stretchr/testify/assert"
)

func TestHeartbeat(t *testing.T) {
	expectedCode := 200
	expectedResult := handlers.Result{
		Success: true,
	}
	httpHandler := httpHandlerMock{}
	mux := http.NewServeMux()
	mux.Handle("/", &httpHandler)

	endpoints.New(mux).WithHeartbeat()

	code, result := hitHeartbeat(t, mux)

	assert.Equal(t, expectedCode, code)
	assert.Equal(t, expectedResult, result)
}

func hitHeartbeat(t *testing.T, mux *http.ServeMux) (int, handlers.Result) {
	t.Helper()
	server := httptest.NewServer(mux)
	defer server.Close()
	response, err := http.Get(server.URL + "/heartbeat")
	if err != nil {
		t.Errorf("unexpected error calling heartbeat: %s", err)
		t.FailNow()
	}
	defer response.Body.Close()
	var result handlers.Result
	json.NewDecoder(response.Body).Decode(&result)
	return response.StatusCode, result
}

type httpHandlerMock struct{}

func (h *httpHandlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}
