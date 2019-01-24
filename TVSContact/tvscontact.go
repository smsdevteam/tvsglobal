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

type MyRespEnvelopeGetContact struct {
	XMLName xml.Name    `xml:"Envelope"`
	Body    bodyContact `xml:"Body"`
}

type bodyContact struct {
	XMLName             xml.Name           `xml:"Body"`
	VGetContactResponse getContactResponse `xml:"GetContactResponse"`
}

type getContactResponse struct {
	XMLName xml.Name `xml:"GetContactResponse"`
	//VGetContactResult st.ContactXML `xml:"GetContactResult"`
	VGetContactResult getContactResult `xml:"GetContactResult"`
}

type getContactResult struct {
	XMLName                 xml.Name `xml:"GetContactResult"`
	ActionDate              string   `xml:"ActionDate"`
	ActionTaken             string   `xml:"ActionTaken" json:"ActionTaken"`
	AllocatedToUser         string   `xml:"AllocatedToUser" json:"AllocatedToUser"`
	CategoryKey             string   `xml:"CategoryKey" json:"CategoryKey"`
	CreatedByUserKey        string   `xml:"CreatedByUserKey" json:"CreatedByUserKey"`
	CreatedDate             string   `xml:"CreatedDate" json:"CreatedDate"`
	CustomerID              string   `xml:"CustomerId" json:"CustomerId"`
	CustomerProductID       string   `xml:"CustomerProductId" json:"CustomerProductId"`
	DeviceID                string   `xml:"DeviceId" json:"DeviceId"`
	ExternalReferenceID     string   `xml:"ExternalReferenceId" json:"ExternalReferenceId"`
	ExternalReferenceID1    string   `xml:"ExternalReferenceId1" json:"ExternalReferenceId1"`
	ExternalReferenceID2    string   `xml:"ExternalReferenceId2" json:"ExternalReferenceId2"`
	ExternalReferenceID3    string   `xml:"ExternalReferenceId3" json:"ExternalReferenceId3"`
	ExternalReferenceID4    string   `xml:"ExternalReferenceId4" json:"ExternalReferenceId4"`
	ExternalReferenceID5    string   `xml:"ExternalReferenceId5" json:"ExternalReferenceId5"`
	ExternalReferenceValue1 string   `xml:"ExternalReferenceValue1" json:"ExternalReferenceValue1"`
	ExternalReferenceValue2 string   `xml:"ExternalReferenceValue2" json:"ExternalReferenceValue2"`
	ExternalReferenceValue3 string   `xml:"ExternalReferenceValue3" json:"ExternalReferenceValue3"`
	ExternalReferenceValue4 string   `xml:"ExternalReferenceValue4" json:"ExternalReferenceValue4"`
	ExternalReferenceValue5 string   `xml:"ExternalReferenceValue5" json:"ExternalReferenceValue5"`
	ContactID               string   `xml:"Id" json:"Id"`
	InvoiceID               string   `xml:"InvoiceId" json:"InvoiceId"`
	LastUpdatedByUserID     string   `xml:"LastUpdatedByUserId" json:"LastUpdatedByUserId"`
	MethodKey               string   `xml:"MethodKey" json:"MethodKey"`
	OrderID                 string   `xml:"OrderId" json:"OrderId"`
	ProblemDescription      string   `xml:"ProblemDescription" json:"ProblemDescription"`
	ProductID               string   `xml:"ProductId" json:"ProductId"`
	StampDate               string   `xml:"StampDate" json:"StampDate"`
	StatusKey               string   `xml:"StatusKey" json:"StatusKey"`
	WorkOrderID             string   `xml:"WorkOrderId" json:"WorkOrderId"`
	Extended                string   `xml:"Extended" json:"Extended"`
}

const getTemplateforGetContact = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<GetContact xmlns="http://tempuri.org/">
		<inContactId>$inContactId</inContactId>
	</GetContact>
</s:Body>
</s:Envelope>`

// GetContactByContactID get info
func GetContactByContactID(iContactID string) st.Contact {
	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetContact"
	l.Request = "ContactID=" + iContactID
	l.Start = t0.String()
	l.InsertappLog("./log/tvscontactapplog.log", "GetContact")

	var oContact st.Contact
	var AppServiceLnk cm.AppServiceURL
	AppServiceLnk = cm.AppReadConfig("ENH")

	url := AppServiceLnk.ICCServiceURL
	client := &http.Client{}

	requestValue := s.Replace(getTemplateforGetContact, "$inContactId", iContactID, -1)

	//log.Println("requestValue: " + requestValue)
	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		return oContact
	}

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/GetContact"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	response, err := client.Do(req)
	if err != nil {
		return oContact
	}
	defer response.Body.Close()

	//log.Println(response.Body)

	if response.StatusCode != 200 {
		return oContact
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return oContact
	}

	//log.Println("contents : " + string(contents[:]))

	myResult := MyRespEnvelopeGetContact{}
	xml.Unmarshal([]byte(contents), &myResult)
	//log.Println(myResult)
	layoutForDatetime := "2006-01-02T15:04:05.000Z"
	oContact.ContactID, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.ContactID), 10, 64)
	oContact.CreatedByUser, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.CreatedByUserKey), 10, 64)
	oContact.CreatedDate, _ = time.Parse(layoutForDatetime, myResult.Body.VGetContactResponse.VGetContactResult.CreatedDate)
	oContact.CustomerID, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.CustomerID), 10, 64)
	oContact.CustomerProductID, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.CustomerProductID), 10, 64)
	oContact.DeviceID, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.DeviceID), 10, 64)
	oContact.ExternalReferenceID = myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceID
	oContact.ExternalReferenceID1, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceID1), 10, 64)
	oContact.ExternalReferenceID2, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceID2), 10, 64)
	oContact.ExternalReferenceID3, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceID3), 10, 64)
	oContact.ExternalReferenceID4, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceID4), 10, 64)
	oContact.ExternalReferenceID5, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceID5), 10, 64)
	oContact.ExternalReferenceValue1 = myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceValue1
	oContact.ExternalReferenceValue2 = myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceValue2
	oContact.ExternalReferenceValue3 = myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceValue3
	oContact.ExternalReferenceValue4 = myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceValue4
	oContact.ExternalReferenceValue5 = myResult.Body.VGetContactResponse.VGetContactResult.ExternalReferenceValue5
	oContact.InvoiceID, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.InvoiceID), 10, 64)
	oContact.LastUpdatedUserID, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.LastUpdatedByUserID), 10, 64)
	oContact.Method = myResult.Body.VGetContactResponse.VGetContactResult.MethodKey
	oContact.OrderID, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.OrderID), 10, 64)
	oContact.ProblemDesc = myResult.Body.VGetContactResponse.VGetContactResult.ProblemDescription
	oContact.ProductID, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.ProductID), 10, 64)
	oContact.StampDate, _ = time.Parse(layoutForDatetime, myResult.Body.VGetContactResponse.VGetContactResult.StampDate)
	oContact.Status = myResult.Body.VGetContactResponse.VGetContactResult.StatusKey
	oContact.WorkOrderID, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.WorkOrderID), 10, 64)

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetContact"
	l.Request = "ContactID=" + iContactID
	l.Response = resp
	l.Start = t0.String()
	l.End = t1.String()
	l.Duration = t2.String()
	l.InsertappLog("./log/tvscontactapplog.log", "GetContact")
	return oContact
}

//GetContactListByCustomerID for find Contacts by Customer
func GetContactListByCustomerID(iCustomerID string) st.ListContact {

	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetContactListByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Start = t0.String()
	l.InsertappLog("./log/tvscontactapplog.log", "GetContactListByCustomerID")

	//log.Println("getContactList")
	var oLContact st.ListContact
	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
		resp = err.Error()
	} else {
		defer db.Close()
		var statement string
		statement = "begin PK_ICC_CONTACT.GetContactListByCustomerId(:0,:1); end;"
		var resultC driver.Rows
		intCustomerID, err := strconv.Atoi(iCustomerID)
		if err != nil {
			log.Fatal(err)
			resp = err.Error()
		} else {
			if _, err := db.Exec(statement, intCustomerID, sql.Out{Dest: &resultC}); err != nil {
				log.Fatal(err)
				resp = err.Error()
			}
			defer resultC.Close()
			values := make([]driver.Value, len(resultC.Columns()))
			var oContacts []st.Contact
			for {
				err = resultC.Next(values)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Println("error:", err)
					resp = err.Error()
				}
				var oContact st.Contact

				if values[0] != nil {
					oContact.ContactID = values[0].(int64)
				}
				if values[1] != nil {
					oContact.ActionDate = values[1].(time.Time)
				}

				oContact.ActionTaken = values[2].(string)

				if values[3] != nil {
					oContact.AllocatedToUser = values[3].(int64)
				}
				if values[4] != nil {
					oContact.Category = values[4].(int64)
				}
				if values[5] != nil {
					oContact.CreatedByUser = values[5].(int64)
				}
				if values[6] != nil {
					oContact.CustomerID = values[6].(int64)
				}
				if values[7] != nil {
					oContact.CustomerProductID = values[7].(int64)
				}
				oContact.Method = values[8].(string)
				if values[9] != nil {
					oContact.OrderID = values[9].(int64)
				}
				oContact.ProblemDesc = values[10].(string)
				if values[11] != nil {
					oContact.ProductID = values[11].(int64)
				}
				if values[12] != nil {
					oContact.StampDate = values[12].(time.Time)
				}
				oContact.Status = values[13].(string)
				if values[14] != nil {
					oContact.WorkOrderID = values[14].(int64)
				}
				if values[15] != nil {
					oContact.CreatedDate = values[15].(time.Time)
				}
				oContact.ExternalReferenceID = values[16].(string)
				if values[17] != nil {
					oContact.DeviceID = values[16].(int64)
				}
				if values[18] != nil {
					oContact.InvoiceID = values[18].(int64)
				}
				if values[19] != nil {
					oContact.LastUpdatedUserID = values[19].(int64)
				}
				if values[20] != nil {
					oContact.ExternalReferenceID1 = values[20].(int64)
				}
				if values[21] != nil {
					oContact.ExternalReferenceID2 = values[21].(int64)
				}
				if values[22] != nil {
					oContact.ExternalReferenceID3 = values[22].(int64)
				}
				if values[23] != nil {
					oContact.ExternalReferenceID4 = values[23].(int64)
				}
				if values[24] != nil {
					oContact.ExternalReferenceID5 = values[24].(int64)
				}
				oContact.ExternalReferenceValue1 = values[25].(string)
				oContact.ExternalReferenceValue2 = values[26].(string)
				oContact.ExternalReferenceValue3 = values[27].(string)
				oContact.ExternalReferenceValue4 = values[28].(string)
				oContact.ExternalReferenceValue5 = values[29].(string)

				oContacts = append(oContacts, oContact)
			}
			oLContact.Contacts = oContacts

		}
	}

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetContactListByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Response = resp
	l.Start = t0.String()
	l.End = t1.String()
	l.Duration = t2.String()
	l.InsertappLog("./log/tvscontactapplog.log", "GetContactListByCustomerID")
	return oLContact
}

const getTemplateforCreateContact = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<CreateContact>
         <inContact>
            <ActionDate>2019-01-21T08:01:43+07:00</ActionDate>
            <ActionTaken>$actionTaken</ActionTaken>
            <AllocatedToUserKey>$allocatedToUserKey</AllocatedToUserKey>
            <CategoryKey>$categoryKey</CategoryKey>
            <CreatedByUserKey>$createdByUserKey</CreatedByUserKey>
            <CreatedDate>2019-01-21T08:01:43+07:00</CreatedDate>
            <CustomerId>$customerId</CustomerId>
            <CustomerProductId>$customerProductId</CustomerProductId>
            <DeviceId>$deviceId</DeviceId>
            <ExternalReferenceId>$externalReferenceId</ExternalReferenceId>
            <ExternalReferenceId1>$externalReferenceId1</ExternalReferenceId1>
            <ExternalReferenceId2>$externalReferenceId2</ExternalReferenceId2>
            <ExternalReferenceId3>$externalReferenceId3</ExternalReferenceId3>
            <ExternalReferenceId4>$externalReferenceId4</ExternalReferenceId4>
            <ExternalReferenceId5>$externalReferenceId5</ExternalReferenceId5>
            <ExternalReferenceValue1>$externalReferenceValue1</ExternalReferenceValue1>
            <ExternalReferenceValue2>$externalReferenceValue2</ExternalReferenceValue2>
            <ExternalReferenceValue3>$externalReferenceValue3</ExternalReferenceValue3>
            <ExternalReferenceValue4>$externalReferenceValue4</ExternalReferenceValue4>
            <ExternalReferenceValue5>$externalReferenceValue5</ExternalReferenceValue5>
            <Id>0</Id>
            <InvoiceId>$invoiceId</InvoiceId>
            <LastUpdatedByUserId>$lastUpdatedByUserId</LastUpdatedByUserId>
            <MethodKey>$methodKey</MethodKey>
            <OrderId>$orderId</OrderId>
            <ProblemDescription>$problemDescription</ProblemDescription>
            <ProductId>$productId</ProductId>
            <StampDate>2019-01-21T08:01:43+07:00</StampDate>
            <StatusKey>$statusKey</StatusKey>
            <WorkOrderId>$workOrderId</WorkOrderId>
            <Extended>$extended</Extended>
         </inContact>
         <inReason>$inReason</inReason>
         <byUser>
            <byUser>$byUser</byUser>
            <byChannel>$byChannel</byChannel>
            <byProject>$byProject</byProject>
            <byHost>$byHost</byHost>
         </byUser>
      </CreateContact>
</s:Body>
</s:Envelope>`

//CreateContact for icc microservice
func CreateContact(iReq st.CreateContactRequest) st.CreateContactResponse {

	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "CreateContact"
	l.Request = "ByUser=" + iReq.ByUser.ByUser + " ByChannel=" + iReq.ByUser.ByChannel
	l.Start = t0.String()
	l.InsertappLog("./log/tvscontactapplog.log", "CreateContact")

	var oRes st.CreateContactResponse
	var AppServiceLnk cm.AppServiceURL
	AppServiceLnk = cm.AppReadConfig("ENH")

	url := AppServiceLnk.ICCServiceURL
	client := &http.Client{}

	var sCustomerProductID string
	var sOrderID string

	sAllocatedToUser := strconv.FormatInt(iReq.InContact.AllocatedToUser, 10)
	sCategoryKey := strconv.FormatInt(iReq.InContact.CategoryKey, 10)
	sCreatedByUserKey := strconv.FormatInt(iReq.InContact.CreatedByUserKey, 10)
	sCustomerID := strconv.FormatInt(iReq.InContact.CustomerID, 10)

	if iReq.InContact.CustomerProductID == 0 {
		sCustomerProductID = ""
	} else {
		sCustomerProductID = strconv.FormatInt(iReq.InContact.CustomerProductID, 10)
	}

	if iReq.InContact.OrderID == 0 {
		sOrderID = ""
	} else {
		sOrderID = strconv.FormatInt(iReq.InContact.OrderID, 10)
	}

	requestValue := s.Replace(getTemplateforCreateContact, "$actionTaken", iReq.InContact.ActionTaken, -1)
	requestValue = s.Replace(requestValue, "$allocatedToUserKey", sAllocatedToUser, -1)
	requestValue = s.Replace(requestValue, "$categoryKey", sCategoryKey, -1)
	requestValue = s.Replace(requestValue, "$createdByUserKey", sCreatedByUserKey, -1)
	requestValue = s.Replace(requestValue, "$customerId", sCustomerID, -1)
	requestValue = s.Replace(requestValue, "$customerProductId", sCustomerProductID, -1)

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "CreateContact"
	l.Request = "ByUser=" + iReq.ByUser.ByUser + " ByChannel=" + iReq.ByUser.ByChannel
	l.Response = resp
	l.Start = t0.String()
	l.End = t1.String()
	l.Duration = t2.String()
	l.InsertappLog("./log/tvscontactapplog.log", "CreateContact")

	return oRes
}
