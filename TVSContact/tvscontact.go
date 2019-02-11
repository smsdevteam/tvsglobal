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

//MyRespEnvelopeGetContact is for get contact
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

//MyRespEnvelopeCreateContact is for Create Contact
type MyRespEnvelopeCreateContact struct {
	XMLName xml.Name          `xml:"Envelope"`
	Body    bodyCreateContact `xml:"Body"`
}

type bodyCreateContact struct {
	XMLName                xml.Name              `xml:"Body"`
	VCreateContactResponse createContactResponse `xml:"CreateContactResponse"`
}

type createContactResponse struct {
	XMLName xml.Name `xml:"CreateContactResponse"`
	//VGetContactResult st.ContactXML `xml:"GetContactResult"`
	VCreateContactResult createContactResult `xml:"CreateContactResult"`
}

type createContactResult struct {
	XMLName     xml.Name `xml:"CreateContactResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
}

//MyRespEnvelopeUpdateContact is for update Contact
type MyRespEnvelopeUpdateContact struct {
	XMLName xml.Name          `xml:"Envelope"`
	Body    bodyUpdateContact `xml:"Body"`
}

type bodyUpdateContact struct {
	XMLName                xml.Name              `xml:"Body"`
	VUpdateContactResponse updateContactResponse `xml:"UpdateContactResponse"`
}

type updateContactResponse struct {
	XMLName xml.Name `xml:"UpdateContactResponse"`
	//VGetContactResult st.ContactXML `xml:"GetContactResult"`
	VUpdateContactResult updateContactResult `xml:"UpdateContactResult"`
}

type updateContactResult struct {
	XMLName     xml.Name `xml:"UpdateContactResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
}

const getTemplateforGetContact = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<GetContact xmlns="http://tempuri.org/">
		<inContactId>$inContactId</inContactId>
	</GetContact>
</s:Body>
</s:Envelope>`

// GetContactByContactID get info
func GetContactByContactID(iContactID string) *st.GetContactResponse {
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

	oRes := st.NewGetContactResponse()
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
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 2
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/GetContact"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 3
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer response.Body.Close()

	//log.Println(response.Body)

	if response.StatusCode != 200 {
		//log.Println(err)
		resp = "response status code :" + strconv.Itoa(response.StatusCode)
		oRes.ErrorCode = 4
		oRes.ErrorDesc = "response status code :" + strconv.Itoa(response.StatusCode)
		return oRes
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		resp = err.Error()
		oRes.ErrorCode = 5
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	//log.Println("contents : " + string(contents[:]))

	myResult := MyRespEnvelopeGetContact{}
	xml.Unmarshal([]byte(contents), &myResult)
	//log.Println(myResult)
	layoutForDatetime := "2006-01-02T15:04:05"
	oContact.ActionDate, _ = time.Parse(layoutForDatetime, myResult.Body.VGetContactResponse.VGetContactResult.ActionDate)
	oContact.ActionTaken = myResult.Body.VGetContactResponse.VGetContactResult.ActionTaken
	oContact.AllocatedToUser, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.AllocatedToUser), 10, 64)
	oContact.Category, _ = strconv.ParseInt((myResult.Body.VGetContactResponse.VGetContactResult.AllocatedToUser), 10, 64)
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

	oRes.GetContactResult = oContact
	oRes.ErrorCode = 0
	oRes.ErrorDesc = "Success"

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
	return oRes
}

//GetContactListByCustomerID for find Contacts by Customer
func GetContactListByCustomerID(iCustomerID string) *st.ListContact {

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
	oLContact := st.NewListContactResponse()
	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {

		log.Println(err)
		resp = err.Error()
		oLContact.ErrorCode = 2
		oLContact.ErrorDesc = err.Error()
		return oLContact
	} else {
		defer db.Close()
		var statement string
		statement = "begin PK_ICC_CONTACT.GetContactListByCustomerId(:0,:1); end;"
		var resultC driver.Rows
		intCustomerID, err := strconv.Atoi(iCustomerID)
		if err != nil {
			log.Println(err)
			resp = err.Error()
			oLContact.ErrorCode = 3
			oLContact.ErrorDesc = err.Error()
			return oLContact
		} else {
			if _, err := db.Exec(statement, intCustomerID, sql.Out{Dest: &resultC}); err != nil {
				log.Println(err)
				resp = err.Error()
				oLContact.ErrorCode = 4
				oLContact.ErrorDesc = err.Error()
				return oLContact
			}
			defer resultC.Close()
			values := make([]driver.Value, len(resultC.Columns()))
			var oContacts []st.Contact
			for {

				colmap := cm.Createmapcol(resultC.Columns())
				//log.Println(colmap)

				err = resultC.Next(values)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Println(err)
					resp = err.Error()
					oLContact.ErrorCode = 5
					oLContact.ErrorDesc = err.Error()
					return oLContact
				}
				var oContact st.Contact

				if values[cm.Getcolindex(colmap, "ID")] != nil {
					oContact.ContactID = values[cm.Getcolindex(colmap, "ID")].(int64)
				}
				if values[cm.Getcolindex(colmap, "ACTION_DATE")] != nil {
					oContact.ActionDate = values[cm.Getcolindex(colmap, "ACTION_DATE")].(time.Time)
				}

				oContact.ActionTaken = values[cm.Getcolindex(colmap, "ACTION_TAKEN")].(string)

				if values[cm.Getcolindex(colmap, "ALLOCATED_TO_USER")] != nil {
					oContact.AllocatedToUser = values[3].(int64)
				}
				if values[cm.Getcolindex(colmap, "CATEGORY")] != nil {
					oContact.Category = values[cm.Getcolindex(colmap, "CATEGORY")].(int64)
				}
				if values[cm.Getcolindex(colmap, "CREATED_BY_USER")] != nil {
					oContact.CreatedByUser = values[cm.Getcolindex(colmap, "CREATED_BY_USER")].(int64)
				}
				if values[cm.Getcolindex(colmap, "CUSTOMER_ID")] != nil {
					oContact.CustomerID = values[cm.Getcolindex(colmap, "CUSTOMER_ID")].(int64)
				}
				if values[cm.Getcolindex(colmap, "CUSTOMER_PRODUCT_ID")] != nil {
					oContact.CustomerProductID = values[cm.Getcolindex(colmap, "CUSTOMER_PRODUCT_ID")].(int64)
				}
				oContact.Method = values[cm.Getcolindex(colmap, "METHOD")].(string)
				if values[cm.Getcolindex(colmap, "ORDER_ID")] != nil {
					oContact.OrderID = values[cm.Getcolindex(colmap, "ORDER_ID")].(int64)
				}
				oContact.ProblemDesc = values[cm.Getcolindex(colmap, "PROBLEM_DESC")].(string)
				if values[cm.Getcolindex(colmap, "PRODUCT_ID")] != nil {
					oContact.ProductID = values[cm.Getcolindex(colmap, "PRODUCT_ID")].(int64)
				}
				if values[cm.Getcolindex(colmap, "STAMP_DATE")] != nil {
					oContact.StampDate = values[cm.Getcolindex(colmap, "STAMP_DATE")].(time.Time)
				}
				oContact.Status = values[cm.Getcolindex(colmap, "STATUS")].(string)
				if values[cm.Getcolindex(colmap, "WORK_ORDER_ID")] != nil {
					oContact.WorkOrderID = values[cm.Getcolindex(colmap, "WORK_ORDER_ID")].(int64)
				}
				if values[cm.Getcolindex(colmap, "CREATED_DATE")] != nil {
					oContact.CreatedDate = values[cm.Getcolindex(colmap, "CREATED_DATE")].(time.Time)
				}
				oContact.ExternalReferenceID = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID")].(string)
				if values[cm.Getcolindex(colmap, "DEVICE_ID")] != nil {
					oContact.DeviceID = values[cm.Getcolindex(colmap, "DEVICE_ID")].(int64)
				}
				if values[cm.Getcolindex(colmap, "INVOICE_ID")] != nil {
					oContact.InvoiceID = values[cm.Getcolindex(colmap, "INVOICE_ID")].(int64)
				}
				if values[cm.Getcolindex(colmap, "LAST_UPDATED_USER_ID")] != nil {
					oContact.LastUpdatedUserID = values[cm.Getcolindex(colmap, "LAST_UPDATED_USER_ID")].(int64)
				}
				if values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID1")] != nil {
					oContact.ExternalReferenceID1 = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID1")].(int64)
				}
				if values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID2")] != nil {
					oContact.ExternalReferenceID2 = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID2")].(int64)
				}
				if values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID3")] != nil {
					oContact.ExternalReferenceID3 = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID3")].(int64)
				}
				if values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID4")] != nil {
					oContact.ExternalReferenceID4 = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID4")].(int64)
				}
				if values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID5")] != nil {
					oContact.ExternalReferenceID5 = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_ID5")].(int64)
				}
				oContact.ExternalReferenceValue1 = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_VALUE1")].(string)
				oContact.ExternalReferenceValue2 = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_VALUE2")].(string)
				oContact.ExternalReferenceValue3 = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_VALUE3")].(string)
				oContact.ExternalReferenceValue4 = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_VALUE4")].(string)
				oContact.ExternalReferenceValue5 = values[cm.Getcolindex(colmap, "EXTERNAL_REFERENCE_VALUE5")].(string)

				oContacts = append(oContacts, oContact)
			}
			oLContact.Contacts = oContacts
			oLContact.ErrorCode = 0
			oLContact.ErrorDesc = "Success"
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
	<CreateContact xmlns="http://tempuri.org/">
         <inContact>
            <ActionDate>$actionDate</ActionDate>
            <ActionTaken>$actionTaken</ActionTaken>
            <AllocatedToUserKey>$allocatedToUserKey</AllocatedToUserKey>
            <CategoryKey>$categoryKey</CategoryKey>
            <CreatedByUserKey>$createdByUserKey</CreatedByUserKey>
            
            <CustomerId>$customerId</CustomerId>
            <CustomerProductId>$customerProductId</CustomerProductId>
            <DeviceId>$deviceId</DeviceId>
            <ExternalReferenceId></ExternalReferenceId>
            <ExternalReferenceId1>0</ExternalReferenceId1>
            <ExternalReferenceId2>0</ExternalReferenceId2>
            <ExternalReferenceId3>0</ExternalReferenceId3>
            <ExternalReferenceId4>0</ExternalReferenceId4>
            <ExternalReferenceId5>0</ExternalReferenceId5>
            <ExternalReferenceValue1></ExternalReferenceValue1>
            <ExternalReferenceValue2></ExternalReferenceValue2>
            <ExternalReferenceValue3></ExternalReferenceValue3>
            <ExternalReferenceValue4></ExternalReferenceValue4>
            <ExternalReferenceValue5></ExternalReferenceValue5>
            <Id>0</Id>
            <InvoiceId>$invoiceId</InvoiceId>
            <LastUpdatedByUserId>$lastUpdatedByUserId</LastUpdatedByUserId>
            <MethodKey>$methodKey</MethodKey>
            <OrderId>$orderId</OrderId>
            <ProblemDescription>$problemDescription</ProblemDescription>
            <ProductId>$productId</ProductId>
            
            <StatusKey>$statusKey</StatusKey>
            <WorkOrderId>$workOrderId</WorkOrderId>
            <Extended></Extended>
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
	var sProductID string
	var sWorkOrderID string
	var sDeviceID string
	var sInvoiceID string

	layoutForDatetime := "2006-01-02T15:04:05.000Z"

	sActionDate := iReq.InContact.ActionDate.Format(layoutForDatetime)
	sAllocatedToUser := strconv.FormatInt(iReq.InContact.AllocatedToUser, 10)
	sCategoryKey := strconv.FormatInt(iReq.InContact.CategoryKey, 10)
	sCreatedByUserKey := strconv.FormatInt(iReq.InContact.CreatedByUserKey, 10)
	sCustomerID := strconv.FormatInt(iReq.InContact.CustomerID, 10)
	sLastUpdateUserID := strconv.FormatInt(iReq.InContact.LastUpdatedByUserID, 10)
	sInReason := strconv.FormatInt(iReq.InReason, 10)
	sCustomerProductID = strconv.FormatInt(iReq.InContact.CustomerProductID, 10)
	sOrderID = strconv.FormatInt(iReq.InContact.OrderID, 10)
	sProductID = strconv.FormatInt(iReq.InContact.ProductID, 10)
	sWorkOrderID = strconv.FormatInt(iReq.InContact.WorkOrderID, 10)
	sDeviceID = strconv.FormatInt(iReq.InContact.DeviceID, 10)
	sInvoiceID = strconv.FormatInt(iReq.InContact.InvoiceID, 10)

	requestValue := s.Replace(getTemplateforCreateContact, "$actionTaken", iReq.InContact.ActionTaken, -1)
	requestValue = s.Replace(requestValue, "$actionDate", sActionDate, -1)
	requestValue = s.Replace(requestValue, "$allocatedToUserKey", sAllocatedToUser, -1)
	requestValue = s.Replace(requestValue, "$categoryKey", sCategoryKey, -1)
	requestValue = s.Replace(requestValue, "$createdByUserKey", sCreatedByUserKey, -1)
	requestValue = s.Replace(requestValue, "$customerId", sCustomerID, -1)
	requestValue = s.Replace(requestValue, "$customerProductId", sCustomerProductID, -1)
	requestValue = s.Replace(requestValue, "$methodKey", iReq.InContact.MethodKey, -1)
	requestValue = s.Replace(requestValue, "$orderId", sOrderID, -1)
	requestValue = s.Replace(requestValue, "$problemDescription", iReq.InContact.ProblemDescription, -1)
	requestValue = s.Replace(requestValue, "$productId", sProductID, -1)
	requestValue = s.Replace(requestValue, "$statusKey", iReq.InContact.StatusKey, -1)
	requestValue = s.Replace(requestValue, "$workOrderId", sWorkOrderID, -1)
	requestValue = s.Replace(requestValue, "$externalReferenceId", iReq.InContact.ExternalReferenceID, -1)
	requestValue = s.Replace(requestValue, "$deviceId", sDeviceID, -1)
	requestValue = s.Replace(requestValue, "$invoiceId", sInvoiceID, -1)
	requestValue = s.Replace(requestValue, "$lastUpdatedByUserId", sLastUpdateUserID, -1)
	requestValue = s.Replace(requestValue, "$inReason", sInReason, -1)
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

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/CreateContact"`)
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

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	//log.Println("contents : " + string(contents[:]))
	myResult := MyRespEnvelopeCreateContact{}
	xml.Unmarshal([]byte(contents), &myResult)
	log.Println(myResult)
	oRes.ResultValue = myResult.Body.VCreateContactResponse.VCreateContactResult.ResultValue
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.VCreateContactResponse.VCreateContactResult.ErrorCode)
	oRes.ErrorDesc = myResult.Body.VCreateContactResponse.VCreateContactResult.ErrorDesc

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

const getTemplateforUpdateContact = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<UpdateContact xmlns="http://tempuri.org/">
         <inContact>
            <ActionDate>$actionDate</ActionDate>
            <ActionTaken>$actionTaken</ActionTaken>
            <AllocatedToUserKey>$allocatedToUserKey</AllocatedToUserKey>
            <CategoryKey>$categoryKey</CategoryKey>
            <CreatedByUserKey>$createdByUserKey</CreatedByUserKey>
            
            <CustomerId>$customerId</CustomerId>
            <CustomerProductId>$customerProductId</CustomerProductId>
            <DeviceId>$deviceId</DeviceId>
            <ExternalReferenceId></ExternalReferenceId>
            <ExternalReferenceId1>0</ExternalReferenceId1>
            <ExternalReferenceId2>0</ExternalReferenceId2>
            <ExternalReferenceId3>0</ExternalReferenceId3>
            <ExternalReferenceId4>0</ExternalReferenceId4>
            <ExternalReferenceId5>0</ExternalReferenceId5>
            <ExternalReferenceValue1></ExternalReferenceValue1>
            <ExternalReferenceValue2></ExternalReferenceValue2>
            <ExternalReferenceValue3></ExternalReferenceValue3>
            <ExternalReferenceValue4></ExternalReferenceValue4>
            <ExternalReferenceValue5></ExternalReferenceValue5>
            <Id>$id</Id>
            <InvoiceId>$invoiceId</InvoiceId>
            <LastUpdatedByUserId>$lastUpdatedByUserId</LastUpdatedByUserId>
            <MethodKey>$methodKey</MethodKey>
            <OrderId>$orderId</OrderId>
            <ProblemDescription>$problemDescription</ProblemDescription>
            <ProductId>$productId</ProductId>
            
            <StatusKey>$statusKey</StatusKey>
            <WorkOrderId>$workOrderId</WorkOrderId>
            <Extended></Extended>
         </inContact>
         <inReason>$inReason</inReason>
         <byUser>
            <byUser>$byUser</byUser>
            <byChannel>$byChannel</byChannel>
            <byProject>$byProject</byProject>
            <byHost>$byHost</byHost>
         </byUser>
      </UpdateContact>
</s:Body>
</s:Envelope>`

func UpdateContact(iReq st.UpdateContactRequest) st.UpdateContactResponse {

	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "UpdateContact"
	l.Request = "ByUser=" + iReq.ByUser.ByUser + " ByChannel=" + iReq.ByUser.ByChannel
	l.Start = t0.String()
	l.InsertappLog("./log/tvscontactapplog.log", "UpdateContact")

	var oRes st.UpdateContactResponse
	var AppServiceLnk cm.AppServiceURL
	AppServiceLnk = cm.AppReadConfig("ENH")

	url := AppServiceLnk.ICCServiceURL
	client := &http.Client{}

	var sCustomerProductID string
	var sOrderID string
	var sProductID string
	var sWorkOrderID string
	var sDeviceID string
	var sInvoiceID string

	layoutForDatetime := "2006-01-02T15:04:05.000Z"

	sActionDate := iReq.InContact.ActionDate.Format(layoutForDatetime)
	sAllocatedToUser := strconv.FormatInt(iReq.InContact.AllocatedToUser, 10)
	sCategoryKey := strconv.FormatInt(iReq.InContact.CategoryKey, 10)
	sCreatedByUserKey := strconv.FormatInt(iReq.InContact.CreatedByUserKey, 10)
	sCustomerID := strconv.FormatInt(iReq.InContact.CustomerID, 10)
	sLastUpdateUserID := strconv.FormatInt(iReq.InContact.LastUpdatedByUserID, 10)
	sInReason := strconv.FormatInt(iReq.InReason, 10)
	sCustomerProductID = strconv.FormatInt(iReq.InContact.CustomerProductID, 10)
	sOrderID = strconv.FormatInt(iReq.InContact.OrderID, 10)
	sProductID = strconv.FormatInt(iReq.InContact.ProductID, 10)
	sWorkOrderID = strconv.FormatInt(iReq.InContact.WorkOrderID, 10)
	sDeviceID = strconv.FormatInt(iReq.InContact.DeviceID, 10)
	sInvoiceID = strconv.FormatInt(iReq.InContact.InvoiceID, 10)
	sID := strconv.FormatInt(iReq.InContact.ContactID, 10)

	requestValue := s.Replace(getTemplateforUpdateContact, "$actionTaken", iReq.InContact.ActionTaken, -1)
	requestValue = s.Replace(requestValue, "$actionDate", sActionDate, -1)
	requestValue = s.Replace(requestValue, "$allocatedToUserKey", sAllocatedToUser, -1)
	requestValue = s.Replace(requestValue, "$categoryKey", sCategoryKey, -1)
	requestValue = s.Replace(requestValue, "$createdByUserKey", sCreatedByUserKey, -1)
	requestValue = s.Replace(requestValue, "$customerId", sCustomerID, -1)
	requestValue = s.Replace(requestValue, "$customerProductId", sCustomerProductID, -1)
	requestValue = s.Replace(requestValue, "$methodKey", iReq.InContact.MethodKey, -1)
	requestValue = s.Replace(requestValue, "$orderId", sOrderID, -1)
	requestValue = s.Replace(requestValue, "$problemDescription", iReq.InContact.ProblemDescription, -1)
	requestValue = s.Replace(requestValue, "$productId", sProductID, -1)
	requestValue = s.Replace(requestValue, "$statusKey", iReq.InContact.StatusKey, -1)
	requestValue = s.Replace(requestValue, "$workOrderId", sWorkOrderID, -1)
	requestValue = s.Replace(requestValue, "$externalReferenceId", iReq.InContact.ExternalReferenceID, -1)
	requestValue = s.Replace(requestValue, "$deviceId", sDeviceID, -1)
	requestValue = s.Replace(requestValue, "$invoiceId", sInvoiceID, -1)
	requestValue = s.Replace(requestValue, "$id", sID, -1)
	requestValue = s.Replace(requestValue, "$lastUpdatedByUserId", sLastUpdateUserID, -1)
	requestValue = s.Replace(requestValue, "$inReason", sInReason, -1)
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

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/UpdateContact"`)
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

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	//log.Println("contents : " + string(contents[:]))
	myResult := MyRespEnvelopeUpdateContact{}
	xml.Unmarshal([]byte(contents), &myResult)
	//log.Println(myResult)
	oRes.ResultValue = myResult.Body.VUpdateContactResponse.VUpdateContactResult.ResultValue
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.VUpdateContactResponse.VUpdateContactResult.ErrorCode)
	oRes.ErrorDesc = myResult.Body.VUpdateContactResponse.VUpdateContactResult.ErrorDesc

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "UpdateContact"
	l.Request = "ByUser=" + iReq.ByUser.ByUser + " ByChannel=" + iReq.ByUser.ByChannel
	l.Response = resp
	l.Start = t0.String()
	l.End = t1.String()
	l.Duration = t2.String()
	l.InsertappLog("./log/tvscontactapplog.log", "UpdateContact")
	return oRes
}
