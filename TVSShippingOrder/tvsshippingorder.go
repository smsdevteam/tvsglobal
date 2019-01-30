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

	cm "github.com/smsdevteam/tvsglobal/Common"     // db
	st "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
)

// MyRespEnvelope Object
type MyRespEnvelope struct {
	XMLName xml.Name
	Body    body
}

type body struct {
	XMLName xml.Name
	Fault   fault
	//	ResponseCreateStockRecv      responseCreateStockRecv      `xml:"CreateStockReceiveDetailsResponse"`
	//	ResponseCreateBuildList      responseCreateBuildList      `xml:"CreateBuildListResponse"`
	//	ResponseAddDeviceToBuildList responseAddDeviceToBuildList `xml:"AddDeviceToBuildListManuallyResponse"`
	//	ResponsePerformBuildList     responsePerformBuildList     `xml:"PerformBuildListActionResponse"`
}

type fault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
	Detail string `xml:"detail"`
}

// GetShippingOrder Method
func GetShippingOrder(iOrderID int64) st.ShippingOrderRes {
	//db, err := sql.Open("goracle", "bgweb/bgweb#1@//tv-uat62-dq.tvsit.co.th:1521/UAT62")
	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var statement string
	var resultC driver.Rows

	// Shipping Order Header
	statement = "begin redibsservice.getdatashippingorderheader(:0,:1); end;"

	if _, err := db.Exec(statement, iOrderID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
	}

	defer resultC.Close()
	values := make([]driver.Value, len(resultC.Columns()))

	var oSO st.ShippingOrderRes
	for {
		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
		}
		oSO.ID = values[0].(int64)
		oSO.DepotFrom = values[1].(int64)
		oSO.DepotTo = values[2].(int64)
		oSO.StatusID = values[3].(int64)
		oSO.StatusDesc = values[4].(string)
		oSO.TypeID = values[5].(int64)
		oSO.TypeDesc = values[6].(string)
		oSO.CreateComments = values[7].(string)
		oSO.CreateReference = values[8].(string)
		oSO.CreateDateTime = values[9].(string)
		oSO.CreateBy = values[10].(int64)
		oSO.CreateByName = values[11].(string)
	}

	// Shipping Order Line
	statement = "begin redibsservice.getdatashippingorderline(:0,:1); end;"

	if _, err := db.Exec(statement, iOrderID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
	}

	defer resultC.Close()
	values = make([]driver.Value, len(resultC.Columns()))

	var oSLList []st.ShippingOrderLineRes

	for {
		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
		}
		var oSL st.ShippingOrderLineRes
		oSL.LineID = values[0].(int64)
		oSL.ShippingOrderID = values[1].(int64)
		oSL.LineNr = values[2].(int64)
		oSL.ProductID = values[3].(int64)
		oSL.ProductKey = values[4].(string)
		oSL.ModelID = values[5].(int64)
		oSL.ModelKey = values[6].(string)
		oSL.Qty = values[7].(int64)
		oSLList = append(oSLList, oSL)
	}
	oSO.ShippingOrderLines = oSLList

	// Shipping Device
	statement = "begin redibsservice.getdatashippingordersn(:0,:1); end;"

	if _, err := db.Exec(statement, iOrderID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
	}

	defer resultC.Close()
	values = make([]driver.Value, len(resultC.Columns()))

	var oSDList []st.ShippingDeviceRes

	for {
		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
		}
		var oSD st.ShippingDeviceRes
		oSD.ShippingOrderID = values[0].(int64)
		oSD.LineID = values[1].(int64)
		oSD.SerialNumber = values[2].(string)
		oSD.StatusID = values[3].(int64)
		oSD.DVResult = values[4].(string)
		oSDList = append(oSDList, oSD)
	}
	oSO.ShippingDevices = oSDList

	return oSO
}

/*
func main() {
	var Val int64
	fmt.Printf("input : ")
	fmt.Scan(&Val)
	r := GetShippingOrder(Val)
	//r := GetShippingOrder(16301898)
	fmt.Println(r)

	json.NewEncoder(os.Stdout).Encode(r)
}
*/

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
  <AgreementId i:nil="true" />
  <Comment>$comment</>
  <CreateDateTime>$createdatetime<CreateDateTime>
  <CustomFields xmlns:d5p1="http://ibs.entriq.net/Core" />
  <CustomerId>$customerid</CustomerId>
  <Destination>Customer</Destination>
  <Extended i:nil="true" />
  <FinancialAccountId i:nil="true" />
  <FullyReceiveReturnedOrder i:nil="true" />
  <Id i:nil="true" />
  <IgnoreAgreementId i:nil="true" />
  <OldStatusId i:nil="true" />
  <ParentOrderId i:nil="true" />
  <ReceivedQuantity i:nil="true" />
  <Reference>$reference</Reference>
  <ReturnedQuantity i:nil="true" />
  <SandboxId i:nil="true" />
  <SandboxSkipValidation i:nil="true" />
  <ShipByDate>$shippeddate</ShipByDate>
  <ShipFromStockHandlerId>$shipfromstockhandler</ShipFromStockHandlerId>
  <ShipToAddressId i:nil="true" />
  <ShipToPartyId>$customerid</ShipToPartyId>
  <ShipToPostalCode i:nil="true" />
  <ShippedDate i:nil="true" />
  <ShippedQuantity i:nil="true" />
  <ShippingMethodId>$shippingmethod</ShippingMethodId>
  <ShippingOrderLines>
	<Items>
	  $itemline
	</Items>
	<More>false</More>
	<Page>0</Page>
	<TotalCount>0</TotalCount>
  </ShippingOrderLines>
  <StatusId i:nil="true" />
  <TotalQuantity i:nil="true" />
  <TrackingNumbers i:nil="true" />
  <TypeId>$ordertype</TypeId>
</order>
<reasonId>0</reasonId>
<printers i:nil="true" xmlns:i="http://www.w3.org/2001/XMLSchema-instance" />
</$method>
</s:Body>
</s:Envelope>`

const getTemplateforcreateslitem = `<ShippingOrderLine>
	<AgreeementDetailId i:nil="true" />
	<CorrelatedHardwareModelId i:nil="true" />
	<CustomFields xmlns:d8p1="http://ibs.entriq.net/Core" />
	<DevicePerAgreementDetailId i:nil="true" />
	<Extended i:nil="true" />
	<ExternalId i:nil="true" />
	<FinanceOptionId i:nil="true" />
	<HardwareModelId i:nil="true" />
	<Id i:nil="true" />
	<NonSubstitutableModel i:nil="true" />
	<OrderLineNumber i:nil="true" />
	<Quantity i:nil="true" />
	<ReceivedQuantity i:nil="true" />
	<ReturnedQuantity i:nil="true" />
	<SandboxId i:nil="true" />
	<SandboxSkipValidation i:nil="true" />
	<SerializedStock i:nil="true" />
	<ShippingOrderId i:nil="true" />
	<TechnicalProductId i:nil="true" />
	<TotalLinkedDevices i:nil="true" />
	<TotalUnlinkedDevices i:nil="true" />
</ShippingOrderLine> `

// CreateShippingOrder Obj
func CreateShippingOrder(iSO st.ShippingOrderReq) (st.ResponseResult, st.ShippingOrderRes) {
	var oRes st.ResponseResult
	var oSORes st.ShippingOrderRes

	// 1. Get Token
	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = err.Error()
		return oRes, oSORes
	}

	url := ServiceLnk.ShippingOrderURL
	client := &http.Client{}

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)

	if len(s.Trim(iSO.ExternalAgent, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + iSO.ExternalAgent + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}
	requestValue := s.Replace(getTemplatefortrxso, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$method", "CreateShippingOrder", -1)

	// Replace Value Parameters : SO Line
	var i int
	for (i=0;i<len(iSO.ShippingOrderLines),i++) {
		subrequest := s.Replace(getTemplateforcreateslitem, "","a",-1)
	}


	// Replace Value Parameters : SO Header
	requestValue = s.Replace(requestValue, "$comment", iSO.Comments, -1)
	requestValue = s.Replace(requestValue, "$createdatetime", iSO.CreateDateTime, -1)
	requestValue = s.Replace(requestValue, "$customerid", strconv.FormatInt(int64(iSO.CustomerId), 10), -1)
	requestValue = s.Replace(requestValue, "$reference", iSO.Reference, -1)
	requestValue = s.Replace(requestValue, "$shippeddate", iSO.CreateDateTime, -1)
	requestValue = s.Replace(requestValue, "$shipfromstockhandler", strconv.FormatInt(int64(iSO.CustomerId), 10), -1)
	requestValue = s.Replace(requestValue, "$shippingmethod", strconv.FormatInt(int64(iSO.CustomerId), 10), -1)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes, oSORes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/OrderManagement/IOrderManagementService/CreateShippingOrder"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes, oSORes
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		oRes.ErrorCode = resp.StatusCode
		oRes.ErrorDesc = resp.Status
		return oRes, oSORes
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes, oSORes
	}

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.Fault.Code)
	oRes.ErrorDesc = myResult.Body.Fault.String

	return oRes, oSORes
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
	requestValue = s.Replace(requestValue, "$soid", strconv.FormatInt(soid, 10), -1)
	requestValue = s.Replace(requestValue, "$reason", strconv.FormatInt(reasonnr, 10), -1)

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
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.Fault.Code)
	oRes.ErrorDesc = myResult.Body.Fault.String


	return oRes
}

=======
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

	cm "github.com/smsdevteam/tvsglobal/Common"     // db
	st "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
)

// MyRespEnvelope Object
type MyRespEnvelope struct {
	XMLName xml.Name
	Body    body
}

type body struct {
	XMLName xml.Name
	Fault   fault
	//	ResponseCreateStockRecv      responseCreateStockRecv      `xml:"CreateStockReceiveDetailsResponse"`
	//	ResponseCreateBuildList      responseCreateBuildList      `xml:"CreateBuildListResponse"`
	//	ResponseAddDeviceToBuildList responseAddDeviceToBuildList `xml:"AddDeviceToBuildListManuallyResponse"`
	//	ResponsePerformBuildList     responsePerformBuildList     `xml:"PerformBuildListActionResponse"`
}

type fault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
	Detail string `xml:"detail"`
}

// GetShippingOrder Method
func GetShippingOrder(iOrderID int64) st.ShippingOrderRes {
	//db, err := sql.Open("goracle", "bgweb/bgweb#1@//tv-uat62-dq.tvsit.co.th:1521/UAT62")
	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var statement string
	var resultC driver.Rows

	// Shipping Order Header
	statement = "begin redibsservice.getdatashippingorderheader(:0,:1); end;"

	if _, err := db.Exec(statement, iOrderID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
	}

	defer resultC.Close()
	values := make([]driver.Value, len(resultC.Columns()))

	var oSO st.ShippingOrderRes
	for {
		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
		}
		oSO.ID = values[0].(int64)
		oSO.DepotFrom = values[1].(int64)
		oSO.DepotTo = values[2].(int64)
		oSO.StatusID = values[3].(int64)
		oSO.StatusDesc = values[4].(string)
		oSO.TypeID = values[5].(int64)
		oSO.TypeDesc = values[6].(string)
		oSO.CreateComments = values[7].(string)
		oSO.CreateReference = values[8].(string)
		oSO.CreateDateTime = values[9].(string)
		oSO.CreateBy = values[10].(int64)
		oSO.CreateByName = values[11].(string)
	}

	// Shipping Order Line
	statement = "begin redibsservice.getdatashippingorderline(:0,:1); end;"

	if _, err := db.Exec(statement, iOrderID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
	}

	defer resultC.Close()
	values = make([]driver.Value, len(resultC.Columns()))

	var oSLList []st.ShippingOrderLineRes

	for {
		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
		}
		var oSL st.ShippingOrderLineRes
		oSL.LineID = values[0].(int64)
		oSL.ShippingOrderID = values[1].(int64)
		oSL.LineNr = values[2].(int64)
		oSL.ProductID = values[3].(int64)
		oSL.ProductKey = values[4].(string)
		oSL.ModelID = values[5].(int64)
		oSL.ModelKey = values[6].(string)
		oSL.Qty = values[7].(int64)
		oSLList = append(oSLList, oSL)
	}
	oSO.ShippingOrderLines = oSLList

	// Shipping Device
	statement = "begin redibsservice.getdatashippingordersn(:0,:1); end;"

	if _, err := db.Exec(statement, iOrderID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
	}

	defer resultC.Close()
	values = make([]driver.Value, len(resultC.Columns()))

	var oSDList []st.ShippingDeviceRes

	for {
		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
		}
		var oSD st.ShippingDeviceRes
		oSD.ShippingOrderID = values[0].(int64)
		oSD.LineID = values[1].(int64)
		oSD.SerialNumber = values[2].(string)
		oSD.StatusID = values[3].(int64)
		oSD.DVResult = values[4].(string)
		oSDList = append(oSDList, oSD)
	}
	oSO.ShippingDevices = oSDList

	return oSO
}

/*
func main() {
	var Val int64
	fmt.Printf("input : ")
	fmt.Scan(&Val)
	r := GetShippingOrder(Val)
	//r := GetShippingOrder(16301898)
	fmt.Println(r)

	json.NewEncoder(os.Stdout).Encode(r)
}
*/

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
  <AgreementId i:nil="true" />
  <Comment>$comment</>
  <CreateDateTime>$createdatetime<CreateDateTime>
  <CustomFields xmlns:d5p1="http://ibs.entriq.net/Core" />
  <CustomerId>$customerid</CustomerId>
  <Destination>Customer</Destination>
  <Extended i:nil="true" />
  <FinancialAccountId i:nil="true" />
  <FullyReceiveReturnedOrder i:nil="true" />
  <Id i:nil="true" />
  <IgnoreAgreementId i:nil="true" />
  <OldStatusId i:nil="true" />
  <ParentOrderId i:nil="true" />
  <ReceivedQuantity i:nil="true" />
  <Reference>$reference</Reference>
  <ReturnedQuantity i:nil="true" />
  <SandboxId i:nil="true" />
  <SandboxSkipValidation i:nil="true" />
  <ShipByDate>$shippeddate</ShipByDate>
  <ShipFromStockHandlerId>$shipfromstockhandler</ShipFromStockHandlerId>
  <ShipToAddressId i:nil="true" />
  <ShipToPartyId>$customerid</ShipToPartyId>
  <ShipToPostalCode i:nil="true" />
  <ShippedDate i:nil="true" />
  <ShippedQuantity i:nil="true" />
  <ShippingMethodId>$shippingmethod</ShippingMethodId>
  <ShippingOrderLines>
	<Items>
	  $itemline
	</Items>
	<More>false</More>
	<Page>0</Page>
	<TotalCount>0</TotalCount>
  </ShippingOrderLines>
  <StatusId i:nil="true" />
  <TotalQuantity i:nil="true" />
  <TrackingNumbers i:nil="true" />
  <TypeId>$ordertype</TypeId>
</order>
<reasonId>0</reasonId>
<printers i:nil="true" xmlns:i="http://www.w3.org/2001/XMLSchema-instance" />
</$method>
</s:Body>
</s:Envelope>`

const getTemplateforcreateslitem = `<ShippingOrderLine>
	<AgreeementDetailId i:nil="true" />
	<CorrelatedHardwareModelId i:nil="true" />
	<CustomFields xmlns:d8p1="http://ibs.entriq.net/Core" />
	<DevicePerAgreementDetailId i:nil="true" />
	<Extended i:nil="true" />
	<ExternalId i:nil="true" />
	<FinanceOptionId i:nil="true" />
	<HardwareModelId i:nil="true" />
	<Id i:nil="true" />
	<NonSubstitutableModel i:nil="true" />
	<OrderLineNumber i:nil="true" />
	<Quantity i:nil="true" />
	<ReceivedQuantity i:nil="true" />
	<ReturnedQuantity i:nil="true" />
	<SandboxId i:nil="true" />
	<SandboxSkipValidation i:nil="true" />
	<SerializedStock i:nil="true" />
	<ShippingOrderId i:nil="true" />
	<TechnicalProductId i:nil="true" />
	<TotalLinkedDevices i:nil="true" />
	<TotalUnlinkedDevices i:nil="true" />
</ShippingOrderLine> `

// CreateShippingOrder Obj
func CreateShippingOrder(iSO st.ShippingOrderReq) (st.ResponseResult, st.ShippingOrderRes) {
	var oRes st.ResponseResult
	var oSORes st.ShippingOrderRes

	// 1. Get Token
	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = err.Error()
		return oRes, oSORes
	}

	url := ServiceLnk.ShippingOrderURL
	client := &http.Client{}

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)

	if len(s.Trim(iSO.ExternalAgent, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + iSO.ExternalAgent + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}
	requestValue := s.Replace(getTemplatefortrxso, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$method", "CreateShippingOrder", -1)

	// Replace Value Parameters : SO Line
	var i int
	for (i=0;i<len(iSO.ShippingOrderLines),i++) {
		subrequest := s.Replace(getTemplateforcreateslitem, "","a",-1)
	}


	// Replace Value Parameters : SO Header
	requestValue = s.Replace(requestValue, "$comment", iSO.Comments, -1)
	requestValue = s.Replace(requestValue, "$createdatetime", iSO.CreateDateTime, -1)
	requestValue = s.Replace(requestValue, "$customerid", strconv.FormatInt(int64(iSO.CustomerId), 10), -1)
	requestValue = s.Replace(requestValue, "$reference", iSO.Reference, -1)
	requestValue = s.Replace(requestValue, "$shippeddate", iSO.CreateDateTime, -1)
	requestValue = s.Replace(requestValue, "$shipfromstockhandler", strconv.FormatInt(int64(iSO.CustomerId), 10), -1)
	requestValue = s.Replace(requestValue, "$shippingmethod", strconv.FormatInt(int64(iSO.CustomerId), 10), -1)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes, oSORes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/OrderManagement/IOrderManagementService/CreateShippingOrder"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes, oSORes
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		oRes.ErrorCode = resp.StatusCode
		oRes.ErrorDesc = resp.Status
		return oRes, oSORes
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes, oSORes
	}

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.Fault.Code)
	oRes.ErrorDesc = myResult.Body.Fault.String

	return oRes, oSORes
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
	requestValue = s.Replace(requestValue, "$soid", strconv.FormatInt(soid, 10), -1)
	requestValue = s.Replace(requestValue, "$reason", strconv.FormatInt(reasonnr, 10), -1)

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
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.Fault.Code)
	oRes.ErrorDesc = myResult.Body.Fault.String


	return oRes
}
