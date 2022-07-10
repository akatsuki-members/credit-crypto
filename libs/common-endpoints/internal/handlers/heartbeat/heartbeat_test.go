package heartbeat_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers"
	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/heartbeat"
	"github.com/stretchr/testify/assert"
)

func TestAddHeartbeat(t *testing.T) {
	// GIVEN
	expectedCode := http.StatusOK
	expectedResult := handlers.Result{
		Success: true,
	}
	mux := http.NewServeMux()

	// WHEN
	heartbeat.Add(mux)
	code, result := serve(t, mux)

	// THEN
	assert.Equal(t, expectedCode, code)
	assert.Equal(t, expectedResult, result)
}

func serve(t *testing.T, mux *http.ServeMux) (int, handlers.Result) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/heartbeat", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)
	got, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("unexpected error reading heartbeat body: %s", err)
		t.FailNow()
	}
	var result handlers.Result
	err = json.Unmarshal(got, &result)
	if err != nil {
		t.Errorf("unexpected error unmarshalling heartbeat response: %s", err)
		t.FailNow()
	}
	return response.Code, result
}
