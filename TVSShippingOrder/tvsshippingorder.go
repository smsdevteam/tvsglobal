package main

import (
	"bytes"
	// "database/sql"
	// "database/sql/driver"
	"encoding/xml"
	// "io"
	"io/ioutil"
	"log"
	"net/http"
	s "strings"
	"time"

	_ "gopkg.in/goracle.v2"

	cm "github.com/smsdevteam/tvsglobal/Common"     // db
	st "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
)

var p = log.Println

// MyRespEnvelope Object
type MyRespEnvelope struct {
	XMLName xml.Name
	Body    body
}

type body struct {
	XMLName                     xml.Name
	Fault                       fault
	ResponseGetShippingOrder    responseGetShippingOrder    `xml:"GetShippingOrderResponse"`
	ResponseCreateShippingOrder responseCreateShippingOrder `xml:"CreateShippingOrderResponse"`
}

type fault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
	Detail string `xml:"detail"`
}

type responseGetShippingOrder struct {
	XMLName                xml.Name               `xml:"GetShippingOrderResponse"`
	GetShippingOrderResult getShippingOrderResult `xml:"GetShippingOrderResult"`
}

type responseCreateShippingOrder struct {
	XMLName                xml.Name               `xml:"CreateShippingOrderResponse"`
	GetShippingOrderResult getShippingOrderResult `xml:"CreateShippingOrderResult"`
}

type getShippingOrderResult struct {
	AgreementId    int64  `xml: agreementId"`
	Comment        string `xml:"Comment"`
	CreateDateTime string `xml:"CreateDateTime"`
	CustomFields   struct {
		CustomFieldValue []struct {
			Extended string `xml:"Extended"`
			Id       int64  `xml:"Id"`
			Name     string `xml:"Name"`
			Sequence int64  `xml:"Sequence"`
			Value    string `xml:"Value"`
		} `xml:"CustomFieldValue"`
	} `xml:"CustomFields"`
	CustomerId                int64  `xml:"CustomerId"`
	Destination               string `xml:"Destination"`
	Extended                  string `xml:"Extended"`
	FinancialAccountId        int64  `xml:"FinancialAccountId"`
	FullyReceiveReturnedOrder bool   `xml:"FullyReceiveReturnedOrder"`
	ID                        int64  `xml:"Id"`
	IgnoreAgreementId         bool   `xml:"IgnoreAgreementId"`
	OldStatusId               int64  `xml:"OldStatusId"`
	ParentOrderId             int64  `xml:"ParentOrderId"`
	ReceivedQuantity          int64  `xml:"ReceivedQuantity"`
	Reference                 string `xml:"Reference"`
	ReturnedQuantity          int64  `xml:"ReturnedQuantity"`
	SandboxId                 int64  `xml:"SandboxId"`
	SandboxSkipValidation     bool   `xml:"SandboxSkipValidation"`
	ShipByDate                string `xml:"ShipByDate"`
	ShipFromStockHandlerId    int64  `xml:"ShipFromStockHandlerId"`
	ShipToAddressId           int64  `xml:"ShipToAddressId"`
	ShipToPartyId             int64  `xml:"ShipToPartyId"`
	ShipToPostalCode          string `xml:"ShipToPostalCode"`
	ShippedDate               string `xml:"ShippedDate"`
	ShippedQuantity           int64  `xml:"ShippedQuantity"`
	ShippingMethodId          int64  `xml:"ShippingMethodId"`
	ShippingOrderLines        struct {
		Items struct {
			ShippingOrderLine []struct {
				AgreeementDetailId        int64 `xml:"AgreeementDetailId"`
				CorrelatedHardwareModelId int64 `xml:"CorrelatedHardwareModelId"`
				CustomFields              struct {
					CustomFieldValue []struct {
						Extended string `xml:"Extended"`
						Id       int64  `xml:"Id"`
						Name     string `xml:"Name"`
						Sequence int64  `xml:"Sequence"`
						Value    string `xml:"Value"`
					} `xml:"CustomFieldValue"`
				} `xml:"CustomFields"`
				DevicePerAgreementDetailId int64  `xml:"DevicePerAgreementDetailId"`
				Extended                   string `xml:"Extended"`
				ExternalId                 string `xml:"ExternalId"`
				FinanceOptionId            int64  `xml:"FinanceOptionId"`
				HardwareModelId            int64  `xml:"HardwareModelId"`
				ID                         int64  `xml:"Id"`
				NonSubstitutableModel      bool   `xml:"NonSubstitutableModel"`
				OrderLineNumber            int64  `xml:"OrderLineNumber"`
				Quantity                   int64  `xml:"Quantity"`
				ReceivedQuantity           int64  `xml:"ReceivedQuantity"`
				ReturnedQuantity           int64  `xml:"ReturnedQuantity"`
				SandboxId                  int64  `xml:"SandboxId"`
				SandboxSkipValidation      bool   `xml:"SandboxSkipValidation"`
				SerializedStock            bool   `xml:"SerializedStock"`
				ShippingOrderId            int64  `xml:"ShippingOrderId"`
				TechnicalProductId         int64  `xml:"TechnicalProductId"`
				TotalLinkedDevices         int64  `xml:"TotalLinkedDevices"`
				TotalUnlinkedDevices       int64  `xml:"TotalUnlinkedDevices"`
			} `xml:"ShippingOrderLine"`
		} `xml:"Items"`
		More       bool  `xml:"More"`
		Page       int64 `xml:"Page"`
		TotalCount int64 `xml:"TotalCount"`
	} `xml:"ShippingOrderLines"`
	StatusId        int64 `xml:"StatusId"`
	TotalQuantity   int64 `xml:"TotalQuantity"`
	TrackingNumbers struct {
		Items struct {
			TrackingNumber []struct {
				Extended        string `xml:"Extended"`
				Id              int64  `xml:"Id"`
				Number          string `xml:"Number"`
				ShippingOrderId int64  `xml:"ShippingOrderId"`
			} `xml:"TrackingNumber"`
		} `xml:"Items"`
		More       bool  `xml:"More"`
		Page       int64 `xml:"Page"`
		TotalCount int64 `xml:"TotalCount"`
	} `xml:"TrackingNumbers"`
	TypeId int64 `xml:"TypeId"`
}

// ShippingOrder : ICC API
const getTemplateAuthenHD = `<s:Header>
  <h:CacheControlHeader i:nil="true" xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns:h="http://ibs.entriq.net/Core" />
  <h:AuthenticationHeader xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns:h="http://ibs.entriq.net/Security">
	<h:ClientName i:nil="true" />
	<h:ClientProof i:nil="true" />
	<h:Culture i:nil="true" />
	<h:Dsn>$dsn</h:Dsn>
	<h:Extended i:nil="true" />
	<h:ExternalAgent i:nil="true" />
	<h:Proof>$password</h:Proof>
	<h:ServerTime>$servicetime</h:ServerTime>
	<h:Token>$token</h:Token>
	<h:UserName>$username</h:UserName>
  </h:AuthenticationHeader>
</s:Header>`

const getTemplatefortrxso = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
	<$method xmlns="http://ibs.entriq.net/OrderManagement">
		<order xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
			$order
		</order>
	  <reasonId>$reason</reasonId>
  	<printers i:nil="true" xmlns:i="http://www.w3.org/2001/XMLSchema-instance" />
	</$method>
</s:Body>
</s:Envelope>`

// CreateShippingOrder Obj
func CreateShippingOrder(SORequest st.ShippingOrderDataReq) st.SOResult {
	var oRes st.ResponseResult
	var oSO st.ShippingOrderData
	var result st.SOResult

	// 1. Get Token
	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = err.Error()
		result.ProcessResult = oRes
		return result
	}

	url := ServiceLnk.ShippingOrderURL
	client := &http.Client{}

	if len(s.Trim(SORequest.SODetail.CreateDateTime, " ")) == 0 {
		SORequest.SODetail.CreateDateTime = time.Now().Format("2006-01-02T15:04:05")
	}
	if len(s.Trim(SORequest.SODetail.ShipByDate, " ")) == 0 {
		SORequest.SODetail.ShipByDate = time.Now().Format("2006-01-02T15:04:05")
	}
	if len(s.Trim(SORequest.SODetail.ShippedDate, " ")) == 0 {
		SORequest.SODetail.ShippedDate = time.Now().Format("2006-01-02T15:04:05")
	}

	output, err := xml.MarshalIndent(SORequest.SODetail, "  ", "    ")
	var so string
	so = string(output)

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)

	if len(s.Trim(SORequest.ByUsername, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + SORequest.ByUsername + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}
	requestValue := s.Replace(getTemplatefortrxso, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$method", "CreateShippingOrder", -1)
	requestValue = s.Replace(requestValue, "$order", so, -1)
	requestValue = s.Replace(requestValue, "$reason", cm.Int64ToStr(SORequest.Reasonnr), -1)
	requestValue = s.Replace(requestValue, "<ShippingOrderData>", "", -1)
	requestValue = s.Replace(requestValue, "</ShippingOrderData>", "", -1)

	// Test Fix
	requestValue = s.Replace(requestValue, "<CustomFields></CustomFields>", `<CustomFields xmlns:d5p1="http://ibs.entriq.net/Core" />`, -1)

	p(requestValue)
	p("----------------------------------------------")

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		result.ProcessResult = oRes
		return result
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/OrderManagement/IOrderManagementService/CreateShippingOrder"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		result.ProcessResult = oRes
		return result
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		oRes.ErrorCode = resp.StatusCode
		oRes.ErrorDesc = resp.Status
		result.ProcessResult = oRes
		return result
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		result.ProcessResult = oRes
		return result
	}

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	oSO.ID = myResult.Body.ResponseCreateShippingOrder.GetShippingOrderResult.ID
	p("Shipping Order Id : ",oSO.ID)

	var gRes st.SOResult
	gRes = GetShippingOrder(oSO.ID, SORequest.ByUsername)

	oRes.ErrorCode = 0
	oRes.ErrorDesc = "SUCCESS"
	result.ProcessResult = oRes
	result.SODetail = gRes.SODetail

	return result
}

const getTemplateforcancelso = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
<CancelShippingOrder xmlns="http://ibs.entriq.net/OrderManagement">
  <orderId>$soid</orderId>
  <reasonId>$reason</reasonId>
</CancelShippingOrder>
</s:Body>
</s:Envelope>`

// CancelShippingOrder Method
func CancelShippingOrder(soid int64, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult

	// 1. Get Token
	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	url := ServiceLnk.ShippingOrderURL
	client := &http.Client{}

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)

	if len(s.Trim(byusername, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + byusername + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}
	requestValue := s.Replace(getTemplateforcancelso, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$soid", cm.Int64ToStr(soid), -1)
	requestValue = s.Replace(requestValue, "$reason", cm.Int64ToStr(reasonnr), -1)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/OrderManagement/IOrderManagementService/CancelShippingOrder"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		oRes.ErrorCode = resp.StatusCode
		oRes.ErrorDesc = resp.Status
		return oRes
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	oRes.ErrorCode = cm.StrToInt(myResult.Body.Fault.Code)
	oRes.ErrorDesc = myResult.Body.Fault.String

	return oRes
}

const getTemplateforgetso = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
    <GetShippingOrder xmlns="http://ibs.entriq.net/OrderManagement">
      <orderId>$soid</orderId>
    </GetShippingOrder>
</s:Body>
</s:Envelope>`

// GetShippingOrder Method
func GetShippingOrder(soid int64, byusername string) st.SOResult {
	var oRes st.ResponseResult
	var oSO st.ShippingOrderData
	var result st.SOResult

	// 1. Get Token
	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = err.Error()
		result.ProcessResult = oRes
		return result
	}

	url := ServiceLnk.ShippingOrderURL
	client := &http.Client{}

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)

	if len(s.Trim(byusername, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + byusername + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}
	requestValue := s.Replace(getTemplateforgetso, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$soid", cm.Int64ToStr(soid), -1)

	//p(requestValue)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		result.ProcessResult = oRes
		return result
	}

	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/OrderManagement/IOrderManagementService/GetShippingOrder"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		result.ProcessResult = oRes
		return result
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		oRes.ErrorCode = resp.StatusCode
		oRes.ErrorDesc = resp.Status
		result.ProcessResult = oRes
		return result
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		result.ProcessResult = oRes
		return result
	}

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	oSO.AgreementID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.AgreementId
	oSO.Comment = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.Comment
	oSO.CreateDateTime = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.CreateDateTime
	oSO.CustomerID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.CustomerId
	oSO.Destination = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.Destination
	oSO.FinancialAccountID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.FinancialAccountId
	oSO.FullyReceiveReturnedOrder = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.FullyReceiveReturnedOrder
	oSO.ID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ID
	oSO.IgnoreAgreementID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.IgnoreAgreementId
	oSO.OldStatusID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.OldStatusId
	oSO.ParentOrderID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ParentOrderId
	oSO.ReceivedQuantity = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ReceivedQuantity
	oSO.Reference = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.Reference
	oSO.ReturnedQuantity = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ReturnedQuantity
	oSO.SandboxID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.SandboxId
	oSO.SandboxSkipValidation = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.SandboxSkipValidation
	oSO.ShipByDate = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShipByDate
	oSO.ShipFromStockHandlerID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShipFromStockHandlerId
	oSO.ShipToAddressID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShipToAddressId
	oSO.ShipToPartyID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShipToPartyId
	oSO.ShipToPostalCode = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShipToPostalCode
	oSO.ShippedDate = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippedDate
	oSO.ShippedQuantity = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippedQuantity
	oSO.ShippingMethodID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingMethodId

	// ShippingOrderLine
	qty := len(myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine)
	for i := 0; i < qty; i++ {
		//p(i)
		var oSL st.ShippingOrderLineStruct

		oSL.AgreeementDetailID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].AgreeementDetailId
		oSL.CorrelatedHardwareModelID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].CorrelatedHardwareModelId
		oSL.DevicePerAgreementDetailID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].DevicePerAgreementDetailId

		oSL.ExternalID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].ExternalId
		oSL.FinanceOptionID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].FinanceOptionId
		oSL.HardwareModelID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].HardwareModelId
		oSL.ID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].ID
		oSL.NonSubstitutableModel = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].NonSubstitutableModel
		oSL.OrderLineNumber = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].OrderLineNumber
		oSL.Quantity = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].Quantity
		oSL.ReceivedQuantity = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].ReceivedQuantity
		oSL.ReturnedQuantity = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].ReturnedQuantity
		oSL.SandboxID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].SandboxId
		oSL.SandboxSkipValidation = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].SandboxSkipValidation
		oSL.SerializedStock = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].SerializedStock
		oSL.ShippingOrderID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].ShippingOrderId
		oSL.TechnicalProductID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].TechnicalProductId
		oSL.TotalLinkedDevices = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].TotalLinkedDevices
		oSL.TotalUnlinkedDevices = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].TotalUnlinkedDevices
		oSO.ShippingOrderLines.Items.ShippingOrderLine = append(oSO.ShippingOrderLines.Items.ShippingOrderLine, oSL)

		// CustomField
		sqty := len(myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].CustomFields.CustomFieldValue)
		for j := 0; j < sqty; j++ {
			var oSCT st.CustomFieldValue

			oSCT.Extended = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].CustomFields.CustomFieldValue[i].Extended
			oSCT.ID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].CustomFields.CustomFieldValue[i].Id
			oSCT.Name = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].CustomFields.CustomFieldValue[i].Name
			oSCT.Sequence = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].CustomFields.CustomFieldValue[i].Sequence
			oSCT.Value = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Items.ShippingOrderLine[i].CustomFields.CustomFieldValue[i].Value
			oSL.CustomFields.CustomFields = append(oSL.CustomFields.CustomFields, oSCT)

		}
	}
	oSO.ShippingOrderLines.More = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.More
	oSO.ShippingOrderLines.Page = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.Page
	oSO.ShippingOrderLines.TotalCount = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.ShippingOrderLines.TotalCount
	oSO.StatusID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.StatusId
	oSO.TotalQuantity = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.TotalQuantity
	oSO.TypeID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.TypeId

	// CustomField
	qty = len(myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.CustomFields.CustomFieldValue)
	for i := 0; i < qty; i++ {
		var oCT st.CustomFieldValue

		oCT.Extended = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.CustomFields.CustomFieldValue[i].Extended
		oCT.ID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.CustomFields.CustomFieldValue[i].Id
		oCT.Name = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.CustomFields.CustomFieldValue[i].Name
		oCT.Sequence = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.CustomFields.CustomFieldValue[i].Sequence
		oCT.Value = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.CustomFields.CustomFieldValue[i].Value
		oSO.CustomFields.CustomFields = append(oSO.CustomFields.CustomFields, oCT)

	}

	// TrackingNumber
	qty = len(myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.TrackingNumbers.Items.TrackingNumber)
	for i := 0; i < qty; i++ {
		var oTR st.TrackingNumber

		oTR.Extended = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.TrackingNumbers.Items.TrackingNumber[i].Extended
		oTR.ID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.TrackingNumbers.Items.TrackingNumber[i].Id
		oTR.Number = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.TrackingNumbers.Items.TrackingNumber[i].Number
		oTR.ShippingOrderID = myResult.Body.ResponseGetShippingOrder.GetShippingOrderResult.TrackingNumbers.Items.TrackingNumber[i].ShippingOrderId

		oSO.TrackingNumbers.Items.TrackingNumbers = append(oSO.TrackingNumbers.Items.TrackingNumbers, oTR)
	}

	oRes.ErrorCode = 0
	oRes.ErrorDesc = "SUCCESS"
	result.ProcessResult = oRes
	result.SODetail = oSO

	return result
}

const getTemplateforshipso = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
		<ShipOrder xmlns="http://ibs.entriq.net/OrderManagement">
			<order xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
				$sodata
			</order>
      <reasonId>$reason</reasonId>
      <printers i:nil="true" xmlns:i="http://www.w3.org/2001/XMLSchema-instance" />			
		</ShipOrder>
</s:Body>
</s:Envelope>`

// ShipOrder Method
func ShipOrder(SOData st.SOResult, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult
	// var SOData st.SOResult
	// SOData = GetShippingOrder(soid, byusername)

	output, err := xml.MarshalIndent(SOData.SODetail, "  ", "    ")
	var so string
	so = string(output)

	// 1. Get Token
	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	url := ServiceLnk.ShippingOrderURL
	client := &http.Client{}

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)

	if len(s.Trim(byusername, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + byusername + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}
	requestValue := s.Replace(getTemplateforshipso, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$sodata", so, -1)
	requestValue = s.Replace(requestValue, "$reason", cm.Int64ToStr(reasonnr), -1)

	//p(requestValue)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/OrderManagement/IOrderManagementService/ShipOrder"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		oRes.ErrorCode = resp.StatusCode
		oRes.ErrorDesc = resp.Status
		return oRes
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)

	return oRes
}
