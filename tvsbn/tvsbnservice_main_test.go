package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetccbschangepackage(t *testing.T) {

	req, err := http.NewRequest("GET", "/tvsbn/ccbschangepackage", nil)
	if err != nil {
		t.Error(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	resp := httptest.NewRecorder()
	vars := map[string]string{
		"customerid": "108218909",
	}
	req = mux.SetURLVars(req, vars)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler := http.HandlerFunc(ccbschangepackage)
	handler.ServeHTTP(resp, req)
	// Check the status code is what we expect.
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
