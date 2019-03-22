//tvsbnservice_main_test.go
package main

import (
	"bytes"
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
		//"customerid": "11",
	}
	req = mux.SetURLVars(req, vars)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler := http.HandlerFunc(ccbschangepackage)
	handler.ServeHTTP(resp, req)
	// Check the status code is what we expect.
	checkResponseCode(t, http.StatusOK, resp.Code)
	if len(resp.Body.String()) < 1 {
		t.Errorf("Expected customerid to be ")
	}
}

func TestPostccbschangepackagep(t *testing.T) {
	payload := []byte(`{
		"Trackingno" : "",
		"TVSOrdReq" : {
		 "Orderid" : "1000",
		 "OrderType" : "1",
		 "ChannelCode" : "TEST",
		 "OrderDate" : "0001-01-01T00:00:00Z",
		 "TVSCustNo" : 108218909,
		 "Custinfo" : {
		  "Text" : "",
		  "BirthDate" : "",
		  "BusinessUnitID" : "",
		  "ClassID" : "",
		  "CustomFields" : {
		   "Text" : "",
		   "CustomFieldValue" : null
		  },
		  "DeviceList" : null,
		  "CustomerSince" : "",
		  "DefaultAddress" : {
		   "Text" : "",
		   "BigCity" : "",
		   "CareOfName" : "",
		   "CountryKey" : "",
		   "CustomerCaptureCategory" : "",
		   "Directions" : "",
		   "Email" : "",
		   "Extended" : {
			"Text" : "",
			"Nil" : ""
		   },
		   "Extra" : "",
		   "Extra1" : "",
		   "Extra2" : "",
		   "Extra3" : "",
		   "Extra4" : "",
		   "Extra5" : "",
		   "ExtraExtra" : "",
		   "Fax1" : "",
		   "Fax2" : "",
		   "FirstName" : "",
		   "FlatOrApartment" : "",
		   "GeoCodeID" : "",
		   "HomePhone" : "",
		   "HomePhoneExt" : "",
		   "HouseNumberAlpha" : "",
		   "HouseNumberNumeric" : "",
		   "CUSTOMERID" : "",
		   "LandMark" : "",
		   "MarketSegmentID" : "",
		   "MiddleName" : "",
		   "PostalCode" : "",
		   "ProvinceKey" : "",
		   "SmallCity" : "",
		   "Street" : "",
		   "Surname" : "",
		   "TitleKey" : "",
		   "ValidAddressID" : "",
		   "WorkPhone" : "",
		   "WorkPhoneExt" : ""
		  },
		  "EmailNotifyOptionKey" : "",
		  "ExemptionCodeKey" : "",
		  "ExemptionFrom" : "",
		  "ExemptionSerialNumber" : "",
		  "Extended" : {
		   "Text" : "",
		   "Nil" : ""
		  },
		  "FiscalCode" : "",
		  "FiscalNumber" : "",
		  "CUSTOMERId" : 108218909,
		  "InternetPassword" : "",
		  "InternetUserID" : "",
		  "IsDistributor" : "",
		  "IsHeadend" : "",
		  "IsProductProvider" : "",
		  "IsServiceProvider" : "",
		  "IsStockHandler" : "",
		  "LanguageKey" : "",
		  "Magazines" : "",
		  "ParentID" : "",
		  "PreferredContactMethodID" : "",
		  "ReferenceNumber" : "",
		  "ReferenceTypeKey" : "",
		  "SegmentationKey" : "",
		  "StatusKey" : "",
		  "TypeKey" : ""
		 }
		}
	   }`)

	req, err := http.NewRequest("POST", "/tvsbn/ccbschangepackagep/", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	if err != nil {
		t.Error(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(ccbschangepackagep)
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
