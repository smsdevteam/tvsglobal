package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	s "strings"
	"time"

	_ "gopkg.in/goracle.v2"

	cm "github.com/smsdevteam/tvsglobal/common"     // db
	st "github.com/smsdevteam/tvsglobal/tvsstructs" // referpath
)

// MyRespEnvelopeGetOffer for GetOffer
type MyRespEnvelopeGetOffer struct {
	XMLName xml.Name  `xml:"Envelope"`
	Body    bodyOffer `xml:"Body"`
}

type bodyOffer struct {
	XMLName           xml.Name         `xml:"Body"`
	VGetOfferResponse getOfferResponse `xml:"GetOfferResponse"`
}

type getOfferResponse struct {
	XMLName         xml.Name       `xml:"GetOfferResponse"`
	VGetOfferResult getOfferResult `xml:"GetOfferResult"`
}

type getOfferResult struct {
	XMLName               xml.Name `xml:"GetOfferResult"`
	Active                string   `xml:"Active"`
	AgreementDetailID     string   `xml:"AgreementDetailId"`
	AgreementID           string   `xml:"AgreementId"`
	ApplyToLevel          string   `xml:"ApplyToLevel"`
	CustomerID            string   `xml:"CustomerId"`
	EndDate               string   `xml:"EndDate"`
	FinancialAccountID    string   `xml:"FinancialAccountId"`
	ID                    string   `xml:"Id"`
	OfferDefinitionID     string   `xml:"OfferDefinitionId"`
	SandboxID             string   `xml:"SandboxId"`
	SandboxSkipValidation string   `xml:"SandboxSkipValidation"`
	StartDate             string   `xml:"StartDate"`
}

const getTemplateforGetOffer = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<GetOffer xmlns="http://tempuri.org/">
		<inOfferId>$inOfferId</inOfferId>
	</GetOffer>
</s:Body>
</s:Envelope>`

//GetOfferByOfferID function
func GetOfferByOfferID(offerID string) *st.GetOfferResponse {
	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSOffer"
	l.FunctionName = "GetOffer"
	l.Request = "OfferID=" + offerID
	l.Start = t0.Format(time.RFC3339Nano)
	l.InsertappLog("./log/tvsofferapplog.log", "GetOffer")

	oRes := st.NewGetOfferResponse()
	var offer st.Offer

	_, err := strconv.Atoi(offerID)
	if err != nil {
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 2
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	var AppServiceLnk cm.AppServiceURL
	AppServiceLnk = cm.AppReadConfig("ENH")

	url := AppServiceLnk.ICCServiceURL
	client := &http.Client{}

	requestValue := s.Replace(getTemplateforGetOffer, "$inOfferId", offerID, -1)

	//log.Println("requestValue: " + requestValue)
	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/GetOffer"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	response, err := client.Do(req)
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer response.Body.Close()

	//log.Println(response.Body)

	if response.StatusCode != 200 {
		oRes.ErrorCode = response.StatusCode
		oRes.ErrorDesc = response.Status
		return oRes
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	//log.Println("contents : " + string(contents[:]))
	myResult := MyRespEnvelopeGetOffer{}
	xml.Unmarshal([]byte(contents), &myResult)
	//log.Println(myResult)
	layoutForDatetime := "2006-01-02T15:04:05"

	offer.Active = myResult.Body.VGetOfferResponse.VGetOfferResult.Active
	offer.AgreementDetailID, _ = strconv.ParseInt((myResult.Body.VGetOfferResponse.VGetOfferResult.AgreementDetailID), 10, 64)
	offer.AgreementID, _ = strconv.ParseInt((myResult.Body.VGetOfferResponse.VGetOfferResult.AgreementID), 10, 64)
	offer.ApplyToLevel = myResult.Body.VGetOfferResponse.VGetOfferResult.ApplyToLevel
	offer.CustomerID, _ = strconv.ParseInt((myResult.Body.VGetOfferResponse.VGetOfferResult.CustomerID), 10, 64)
	offer.EndDate, _ = time.Parse(layoutForDatetime, myResult.Body.VGetOfferResponse.VGetOfferResult.EndDate)
	offer.FinancialAccountID, _ = strconv.ParseInt((myResult.Body.VGetOfferResponse.VGetOfferResult.FinancialAccountID), 10, 64)
	offer.ID, _ = strconv.ParseInt((myResult.Body.VGetOfferResponse.VGetOfferResult.ID), 10, 64)
	offer.OfferDefinitionID, _ = strconv.ParseInt((myResult.Body.VGetOfferResponse.VGetOfferResult.OfferDefinitionID), 10, 64)
	offer.SandboxID, _ = strconv.ParseInt((myResult.Body.VGetOfferResponse.VGetOfferResult.SandboxID), 10, 64)
	offer.SandboxSkipValidation = myResult.Body.VGetOfferResponse.VGetOfferResult.SandboxSkipValidation
	offer.StartDate, _ = time.Parse(layoutForDatetime, myResult.Body.VGetOfferResponse.VGetOfferResult.StartDate)

	oRes.GetOfferResult = offer
	oRes.ErrorCode = 0
	oRes.ErrorDesc = ""

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSOffer"
	l.FunctionName = "GetOffer"
	l.Request = "OfferID=" + offerID
	l.Response = resp
	l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	l.InsertappLog("./log/tvsofferapplog.log", "GetOffer")
	//test
	return oRes
}

//GetListOfferByCustomerID get list offer by customer id
func GetListOfferByCustomerID(iCustomerID string) *st.GetListOfferResult {

	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSOffer"
	l.FunctionName = "GetListOfferByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Start = t0.Format(time.RFC3339Nano)
	l.InsertappLog("./log/tvsofferapplog.log", "GetListOfferByCustomerID")

	oRes := st.NewGetListOfferResult()
	//var oListOffer st.ListOffer

	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 2
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer db.Close()
	var statement string
	statement = "begin PK_ICC_OFFER.GetOfferByCustomerID(:0,:1); end;"
	var resultC driver.Rows
	intCustomerID, err := strconv.Atoi(iCustomerID)
	if err != nil {
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 3
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	if _, err := db.Exec(statement, intCustomerID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
		resp = err.Error()
		oRes.ErrorCode = 4
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer resultC.Close()
	values := make([]driver.Value, len(resultC.Columns()))
	//var oLOffer []st.Offer
	for {
		colmap := cm.Createmapcol(resultC.Columns())

		log.Println(colmap)

		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
			resp = err.Error()
			oRes.ErrorCode = 5
			oRes.ErrorDesc = err.Error()
			return oRes
		}
		//var oOffer st.Offer
	}

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSOffer"
	l.FunctionName = "GetListOfferByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Response = resp
	l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	l.InsertappLog("./log/tvsofferapplog.log", "GetListOfferByCustomerID")
	//test
	return oRes
}
