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

// MyRespEnvelopeCreateOffer for GetOffer
type MyRespEnvelopeCreateOffer struct {
	XMLName xml.Name        `xml:"Envelope"`
	Body    bodyCreateOffer `xml:"Body"`
}

type bodyCreateOffer struct {
	XMLName              xml.Name            `xml:"Body"`
	VCreateOfferResponse createOfferResponse `xml:"CreateOfferResponse"`
}

type createOfferResponse struct {
	XMLName            xml.Name          `xml:"CreateOfferResponse"`
	VCreateOfferResult createOfferResult `xml:"CreateOfferResult"`
}

type createOfferResult struct {
	XMLName     xml.Name `xml:"CreateOfferResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
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
	var oListOffer st.ListOffer

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
	var oLOffer []st.Offer
	for {
		colmap := cm.Createmapcol(resultC.Columns())

		//log.Println(colmap)

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

		var oOffer st.Offer
		if values[cm.Getcolindex(colmap, "ACTIVE")] != nil {
			iActive := values[cm.Getcolindex(colmap, "ACTIVE")].(int64)
			if iActive == 1 {
				oOffer.Active = "true"
			} else {
				oOffer.Active = "false"
			}
		}
		if values[cm.Getcolindex(colmap, "AGREEMENT_DETAIL_ID")] != nil {
			oOffer.AgreementDetailID = values[cm.Getcolindex(colmap, "AGREEMENT_DETAIL_ID")].(int64)
		}
		if values[cm.Getcolindex(colmap, "AGREEMENT_ID")] != nil {
			oOffer.AgreementID = values[cm.Getcolindex(colmap, "AGREEMENT_ID")].(int64)
		}
		//oOffer.ApplyToLevel = values[cm.Getcolindex(colmap, "AGREEMENT_ID")].(int64)
		if values[cm.Getcolindex(colmap, "CUSTOMER_ID")] != nil {
			oOffer.CustomerID = values[cm.Getcolindex(colmap, "CUSTOMER_ID")].(int64)
		}
		if values[cm.Getcolindex(colmap, "END_DATE")] != nil {
			oOffer.EndDate = values[cm.Getcolindex(colmap, "END_DATE")].(time.Time)
		}
		if values[cm.Getcolindex(colmap, "FINANCIAL_ACCOUNT_ID")] != nil {
			oOffer.FinancialAccountID = values[cm.Getcolindex(colmap, "FINANCIAL_ACCOUNT_ID")].(int64)
		}
		if values[cm.Getcolindex(colmap, "ID")] != nil {
			oOffer.ID = values[cm.Getcolindex(colmap, "ID")].(int64)
		}
		if values[cm.Getcolindex(colmap, "OFFER_DEFINITION_ID")] != nil {
			oOffer.OfferDefinitionID = values[cm.Getcolindex(colmap, "OFFER_DEFINITION_ID")].(int64)
		}
		if values[cm.Getcolindex(colmap, "SANDBOX_ID")] != nil {
			oOffer.SandboxID = values[cm.Getcolindex(colmap, "SANDBOX_ID")].(int64)
		}
		if values[cm.Getcolindex(colmap, "START_DATE")] != nil {
			oOffer.StartDate = values[cm.Getcolindex(colmap, "START_DATE")].(time.Time)
		}

		oLOffer = append(oLOffer, oOffer)
	}
	oListOffer.Offers = oLOffer
	oRes.MyListOffer = oListOffer
	if oRes.ErrorCode == 1 {
		oRes.ErrorCode = 0
		oRes.ErrorDesc = "Success"
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

const getTemplateforCreateOffer = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<CreateOffer xmlns="http://tempuri.org/">
		<inOffer>
			<Active>$active</Active>
			<AgreementDetailId>$agreementDetailId</AgreementDetailId>
			<AgreementId>$agreementId</AgreementId>
			<CustomerId>$customerId</CustomerId>
			$endDate
			<FinancialAccountId>$financialAccountId</FinancialAccountId>
			<Id>0</Id>
			<OfferDefinitionId>$offerDefinitionId</OfferDefinitionId>
			<SandboxId>$sandboxId</SandboxId>
			<StartDate>$startDate</StartDate>
			<Extended>$extended</Extended>
		</inOffer>
		<inReason>$inReason</inReason>
		<byUser>
			<byUser>$byUser</byUser>
            <byChannel>$byChannel</byChannel>
            <byProject>$byProject</byProject>
            <byHost>$byHost</byHost>
		</byUser>
	</CreateOffer>
</s:Body>
</s:Envelope>`

//CreateOffer for icc microservice
func CreateOffer(iReq st.CreateOfferRequest) *st.CreateOfferResponse {

	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSOffer"
	l.FunctionName = "CreateOffer"
	l.Request = "ByUser=" + iReq.ByUser.ByUser + " ByChannel=" + iReq.ByUser.ByChannel
	l.Start = t0.Format(time.RFC3339Nano)
	l.InsertappLog("./log/tvsofferapplog.log", "CreateOffer")

	oRes := st.NewCreateOfferResponse()

	var AppServiceLnk cm.AppServiceURL
	AppServiceLnk = cm.AppReadConfig("ENH")

	url := AppServiceLnk.ICCServiceURL
	client := &http.Client{}

	sAgreementDetailID := strconv.FormatInt(iReq.InOffer.AgreementDetailID, 10)
	sAgreementID := strconv.FormatInt(iReq.InOffer.AgreementID, 10)
	sCustomerID := strconv.FormatInt(iReq.InOffer.CustomerID, 10)
	//checkTime
	var sEndDate string
	if iReq.InOffer.EndDate != "" {
		layoutForDatetime := "2006-01-02T15:04:05Z"
		tEndDate, err := time.Parse(layoutForDatetime, iReq.InOffer.EndDate)
		if err != nil {
			oRes.ErrorCode = 2
			oRes.ErrorDesc = err.Error()
			return oRes
		}
		sEndDate = "<EndDate>" + (tEndDate).Format("2006-01-02T15:04:05") + "</EndDate>"
	}

	sFinancialAccountID := strconv.FormatInt(iReq.InOffer.FinancialAccountID, 10)
	sOfferDefinitionID := strconv.FormatInt(iReq.InOffer.OfferDefinitionID, 10)
	sSandboxID := strconv.FormatInt(iReq.InOffer.SandboxID, 10)
	//sSandboxSkipValidation := strconv.FormatInt(iReq.InOffer.SandboxSkipValidation, 10)
	sStartDate := (iReq.InOffer.StartDate).Format("2006-01-02T15:04:05")
	sinReason := strconv.FormatInt(iReq.InReason, 10)

	requestValue := s.Replace(getTemplateforCreateOffer, "$active", iReq.InOffer.Active, -1)
	requestValue = s.Replace(requestValue, "$agreementDetailId", sAgreementDetailID, -1)
	requestValue = s.Replace(requestValue, "$agreementId", sAgreementID, -1)
	//requestValue = s.Replace(requestValue, "$applyToLevel", iReq.InOffer.ApplyToLevel, -1)
	requestValue = s.Replace(requestValue, "$customerId", sCustomerID, -1)
	requestValue = s.Replace(requestValue, "$endDate", sEndDate, -1)
	requestValue = s.Replace(requestValue, "$financialAccountId", sFinancialAccountID, -1)
	requestValue = s.Replace(requestValue, "$offerDefinitionId", sOfferDefinitionID, -1)
	requestValue = s.Replace(requestValue, "$sandboxId", sSandboxID, -1)
	//requestValue = s.Replace(requestValue, "$sandboxSkipValidation", iReq.InOffer.SandboxSkipValidation, -1)
	requestValue = s.Replace(requestValue, "$startDate", sStartDate, -1)
	requestValue = s.Replace(requestValue, "$extended", iReq.InOffer.Extended, -1)
	requestValue = s.Replace(requestValue, "$inReason", sinReason, -1)
	requestValue = s.Replace(requestValue, "$byUser", iReq.ByUser.ByUser, -1)
	requestValue = s.Replace(requestValue, "$byChannel", iReq.ByUser.ByChannel, -1)
	requestValue = s.Replace(requestValue, "$byProject", iReq.ByUser.ByProject, -1)
	requestValue = s.Replace(requestValue, "$byHost", iReq.ByUser.ByHost, -1)

	//log.Println("requestValue: " + requestValue)
	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/CreateOffer"`)
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

	myResult := MyRespEnvelopeCreateOffer{}
	xml.Unmarshal([]byte(contents), &myResult)
	//log.Println(myResult)
	oRes.ResultValue = myResult.Body.VCreateOfferResponse.VCreateOfferResult.ResultValue
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.VCreateOfferResponse.VCreateOfferResult.ErrorCode)
	oRes.ErrorDesc = myResult.Body.VCreateOfferResponse.VCreateOfferResult.ErrorDesc

	//log.Println(oRes)

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSOffer"
	l.FunctionName = "CreateOffer"
	l.Request = "ByUser=" + iReq.ByUser.ByUser
	l.Response = resp
	l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	l.InsertappLog("./log/tvsofferapplog.log", "CreateOffer")
	return oRes
}
