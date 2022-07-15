package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers"
	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/health"
	"github.com/stretchr/testify/assert"
)

func TestAddHealth(t *testing.T) {
	// GIVEN
	expectedHealth := []health.Item{
		{"Database", true},
		{"Cache", true},
		{"OtherService", true},
	}
	healthChecker := newMockedHealthChecker(newHealthReportOK(), true)
	expectedCode := 200
	expectedResponse := handlers.Result{
		Success: true,
		Data:    expectedHealth,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleFuncOk)
	health.Add(mux, healthChecker)
	server := httptest.NewServer(mux)
	defer server.Close()

	// WHEN
	res, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		t.FailNow()
	}
	result := decodeHealthResponse(t, res)
	// THEN
	assert.Equal(t, expectedCode, res.StatusCode)
	assert.Equal(t, expectedResponse, result)

}

func TestAddHealthThatFails(t *testing.T) {
	// GIVEN
	expectedHealth := []health.Item{
		{"Database", false},
		{"Cache", true},
		{"OtherService", true},
	}
	healthChecker := newMockedHealthChecker(newHealthReportFailed(), false)
	expectedCode := 500
	expectedResponse := handlers.Result{
		Success: false,
		Data:    expectedHealth,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleFuncOk)
	health.Add(mux, healthChecker)
	server := httptest.NewServer(mux)
	defer server.Close()

	// WHEN
	res, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		t.FailNow()
	}
	result := decodeHealthResponse(t, res)
	// THEN
	assert.Equal(t, expectedCode, res.StatusCode)
	assert.Equal(t, expectedResponse, result)

}

func decodeHealthResponse(t *testing.T, res *http.Response) handlers.Result {
	t.Helper()
	defer res.Body.Close()
	var result handlers.Result
	err := json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		t.FailNow()
	}
	resultData, err := json.Marshal(result.Data)
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		t.FailNow()
	}
	var items []health.Item
	err = json.Unmarshal(resultData, &items)
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		t.FailNow()
	}
	result.Data = items
	return result
}

func handleFuncOk(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func newMockedHealthChecker(report []health.Item, healthy bool) func() health.Report {
	return func() health.Report {
		return health.Report{
			Healthy: healthy,
			Data:    report,
		}
	}
}

func newHealthReportOK() []health.Item {
	return []health.Item{
		{"Database", true},
		{"Cache", true},
		{"OtherService", true},
	}
}

func newHealthReportFailed() []health.Item {
	return []health.Item{
		{"Database", false},
		{"Cache", true},
		{"OtherService", true},
	}
}
