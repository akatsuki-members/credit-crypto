package endpoints_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers"
	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/health"
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

func TestHealth(t *testing.T) {
	expectedCode := 200
	expectedHealth := []health.Item{
		{Name: "Database", Healthy: true},
		{Name: "Cache", Healthy: true},
		{Name: "OtherService", Healthy: true},
	}
	expectedResult := handlers.Result{
		Success: true,
		Data:    expectedHealth,
	}
	httpHandler := httpHandlerMock{}
	mux := http.NewServeMux()
	mux.Handle("/", &httpHandler)
	healthChecker := newMockedHealthChecker(newHealthReportOK(), true)

	endpoints.New(mux).WithHealth(healthChecker)

	code, result := hitHealth(t, mux)

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

func hitHealth(t *testing.T, mux *http.ServeMux) (int, handlers.Result) {
	t.Helper()
	server := httptest.NewServer(mux)
	defer server.Close()
	response, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Errorf("unexpected error calling health: %s", err)
		t.FailNow()
	}
	result := decodeHealthResponse(t, response)
	return response.StatusCode, result
}

type httpHandlerMock struct{}

func (h *httpHandlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
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
		{Name: "Database", Healthy: true},
		{Name: "Cache", Healthy: true},
		{Name: "OtherService", Healthy: true},
	}
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
