package info_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers"
	"github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/info"
	"github.com/stretchr/testify/assert"
)

func TestAddInfo(t *testing.T) {
	// GIVEN
	expectedCode := http.StatusOK
	infoReport := info.Report{
		Name:    "audit-app",
		Commit:  "963e91b",
		Version: "1.5.8",
	}
	expectedResult := handlers.Result{
		Success: true,
		Data:    infoReport,
	}
	mux := http.NewServeMux()

	// WHEN
	info.Add(mux, infoReport)
	code, result := serve(t, mux)

	// THEN
	assert.Equal(t, expectedCode, code)
	assert.Equal(t, expectedResult, result)
}

func serve(t *testing.T, mux *http.ServeMux) (int, handlers.Result) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/info", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)
	got, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("unexpected error reading info body: %s", err)
		t.FailNow()
	}
	var result handlers.Result
	err = json.Unmarshal(got, &result)
	if err != nil {
		t.Errorf("unexpected error unmarshalling info response: %s", err)
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

	return response.Code, result
}
