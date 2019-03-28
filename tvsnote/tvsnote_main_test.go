//tvsnote_main_test.go
package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetgetnote(t *testing.T) {
	req, err := http.NewRequest("GET", "/tvsnote/getnote", nil)
	if err != nil {
		t.Error(err)
	}
	resp := httptest.NewRecorder()
	vars := map[string]string{
		"noteid": "14789084",
	}
	req = mux.SetURLVars(req, vars)
	handler := http.HandlerFunc(getnote)
	handler.ServeHTTP(resp, req)

	checkResponseCode(t, http.StatusOK, resp.Code)
	if len(resp.Body.String()) < 1 {
		t.Errorf("Expected noteid to be ")
	}
}

func TestGetgetlistnote(t *testing.T) {
	req, err := http.NewRequest("GET", "/tvsnote/getlistnote", nil)
	if err != nil {
		t.Error(err)
	}
	resp := httptest.NewRecorder()
	vars := map[string]string{
		"customerid": "20202070",
	}
	req = mux.SetURLVars(req, vars)
	handler := http.HandlerFunc(getnote)
	handler.ServeHTTP(resp, req)

	checkResponseCode(t, http.StatusOK, resp.Code)
	if len(resp.Body.String()) < 1 {
		t.Errorf("Expected customerid to be ")
	}
}

func TestPostcreatenote(t *testing.T) {

	payload := []byte(`{
		"inNote": 
		  {
			"CustomerId": 20202070,
			"CreatedByUserId": 797,
			"ActionUserKey": 797,
			"CategoryKey": "I",
			"CompletionStageKey": "N",
			"Body": "Test Create from Arnon #11",
			"CreateDate":"2012-04-23T18:25:43.511Z",
			"Id":0
		  }
		,
		"InReason": 1152,
		
		"byUser":	
			{
				"byUser":"arnon12",
				"byChannel":"",
				"byProject":"",
				"byHost":""
			}
		
		
	}`)

	req, err := http.NewRequest("POST", "/tvsnote/createnote/", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	if err != nil {
		t.Error(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(createnote)
	handler.ServeHTTP(resp, req)

	checkResponseCode(t, http.StatusOK, resp.Code)
	if len(resp.Body.String()) < 1 {
		t.Errorf("Expected customerid to be ")
	}

}

func TestPostupdatenote(t *testing.T) {

	payload := []byte(`{
		"inNote": 
		  {
			"CustomerId": 20202070,
			"CreatedByUserId": 797,
			"ActionUserKey": 797,
			"CategoryKey": "I",
			"CompletionStageKey": "Z",
			"Body": "Test update from Arnon#1",
			"CreateDate":"2012-04-23T18:25:43.511Z",
			"Id":14789322,
			"Extended": ""
		  }
		,
		"InReason": 0,
		
		"byUser":	
			{
				"byUser":"arnon12",
				"byChannel":"",
				"byProject":"",
				"byHost":""
			}
		
		
	}`)

	req, err := http.NewRequest("POST", "/tvsnote/updatenote/", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	if err != nil {
		t.Error(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(updatenote)
	handler.ServeHTTP(resp, req)
	// Check the status code is what we expect.
	checkResponseCode(t, http.StatusOK, resp.Code)
	if len(resp.Body.String()) < 1 {
		t.Errorf("Expected customerid to be ")
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
