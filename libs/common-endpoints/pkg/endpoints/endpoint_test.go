package endpoints_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers"
	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/info"
	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/pkg/endpoints"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestHeartbeat(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	expectedCode := 200
	expectedHealth := []endpoints.Item{
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

func TestInfo(t *testing.T) {
	t.Parallel()

	expectedCode := 200
	expectedInfo := info.Report{
		Name:    "audit-app",
		Commit:  "963e91b",
		Version: "1.5.3",
	}
	expectedResult := handlers.Result{
		Success: true,
		Data:    expectedInfo,
	}
	givenInfo := endpoints.Info{
		Name:    "audit-app",
		Commit:  "963e91b",
		Version: "1.5.3",
	}
	httpHandler := httpHandlerMock{}
	mux := http.NewServeMux()
	mux.Handle("/", &httpHandler)

	endpoints.New(mux).WithInfo(givenInfo)
	code, result := hitInfo(t, mux)

	assert.Equal(t, expectedCode, code)
	assert.Equal(t, expectedResult, result)
}

func hitHeartbeat(t *testing.T, mux *http.ServeMux) (int, handlers.Result) {
	t.Helper()

	server := httptest.NewServer(mux)
	defer server.Close()

	response, err := doHTTPGet(t, server.URL+"/heartbeat")
	if err != nil {
		t.Errorf("unexpected error calling heartbeat: %s", err)
		t.FailNow()
	}
	defer response.Body.Close()

	var result handlers.Result

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	return response.StatusCode, result
}

func hitHealth(t *testing.T, mux *http.ServeMux) (int, handlers.Result) {
	t.Helper()

	server := httptest.NewServer(mux)
	defer server.Close()

	response, err := doHTTPGet(t, server.URL+"/health")
	if err != nil {
		t.Errorf("unexpected error calling health: %s", err)
		t.FailNow()
	}

	result := decodeHealthResponse(t, response)

	return response.StatusCode, result
}

func hitInfo(t *testing.T, mux *http.ServeMux) (int, handlers.Result) {
	t.Helper()

	server := httptest.NewServer(mux)
	defer server.Close()

	response, err := doHTTPGet(t, server.URL+"/info")
	if err != nil {
		t.Errorf("unexpected error calling health: %s", err)
		t.FailNow()
	}

	result := decodeInfoResponse(t, response)

	return response.StatusCode, result
}

type httpHandlerMock struct{}

func (h *httpHandlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func newMockedHealthChecker(report []endpoints.Item, healthy bool) func() endpoints.Report {
	return func() endpoints.Report {
		return endpoints.Report{
			Healthy: healthy,
			Data:    report,
		}
	}
}

func newHealthReportOK() []endpoints.Item {
	return []endpoints.Item{
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

	var items []endpoints.Item

	err = json.Unmarshal(resultData, &items)
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		t.FailNow()
	}

	result.Data = items

	return result
}

func decodeInfoResponse(t *testing.T, res *http.Response) handlers.Result {
	t.Helper()

	var result handlers.Result

	err := json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		t.FailNow()
	}

	bytes, err := json.Marshal(result.Data)
	if err != nil {
		t.Errorf("unexpected error mrshalling result.Data response: %s", err)
		t.FailNow()
	}

	var infoReport info.Report

	err = json.Unmarshal(bytes, &infoReport)
	if err != nil {
		t.Errorf("unexpected error unmarshalling info response: %s", err)
		t.FailNow()
	}

	result.Data = infoReport

	return result
}

func doHTTPGet(t *testing.T, url string) (*http.Response, error) {
	t.Helper()

	ctx := context.TODO()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("unexpected error createing heartbeat request: %s", err)
	}

	client := new(http.Client)

	response, err := client.Do(req)

	return response, errors.Wrap(err, "unexpected error")
}
