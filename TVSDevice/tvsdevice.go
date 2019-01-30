package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "gopkg.in/goracle.v2"

	s "strings"

	cm "github.com/smsdevteam/tvsglobal/Common"     // db
	st "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
)

var p = fmt.Println

// MyRespEnvelope
type MyRespEnvelope struct {
	XMLName xml.Name
	Body    body
}

type body struct {
	XMLName                      xml.Name
	Fault                        fault
	ResponseCreateStockRecv      responseCreateStockRecv      `xml:"CreateStockReceiveDetailsResponse"`
	ResponseCreateBuildList      responseCreateBuildList      `xml:"CreateBuildListResponse"`
	ResponseAddDeviceToBuildList responseAddDeviceToBuildList `xml:"AddDeviceToBuildListManuallyResponse"`
	ResponsePerformBuildList     responsePerformBuildList     `xml:"PerformBuildListActionResponse"`
}

type fault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
	Detail string `xml:"detail"`
}

type responseCreateStockRecv struct {
	XMLName                         xml.Name                        `xml:"CreateStockReceiveDetailsResponse"`
	CreateStockReceiveDetailsResult createStockReceiveDetailsResult `xml:"CreateStockReceiveDetailsResult"`
}

type createStockReceiveDetailsResult struct {
	BatchComment         string `xml:"BatchComment"`
	BatchNumber          string `xml:"BatchNumber"`
	BatchReferenceNumber string `xml:"BatchReferenceNumber"`
	DepreciationDetail   struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"DepreciationDetail"`
	Extended struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"Extended"`
	FromStockHanderId        string `xml:"FromStockHanderId"`
	ID                       string `xml:"Id"`
	MACAddress1              string `xml:"MACAddress1"`
	MACAddress2              string `xml:"MACAddress2"`
	OrderId                  string `xml:"OrderId"`
	PalletId                 string `xml:"PalletId"`
	Reason                   string `xml:"Reason"`
	ReorderReason            string `xml:"ReorderReason"`
	ToStockHanderId          string `xml:"ToStockHanderId"`
	UseRangeToDetermineModel string `xml:"UseRangeToDetermineModel"`
	WarrantyEndDate          string `xml:"WarrantyEndDate"`
}

type responseCreateBuildList struct {
	XMLName               xml.Name              `xml:"CreateBuildListResponse"`
	CreateBuildListResult createBuildListResult `xml:"CreateBuildListResult"`
}

type createBuildListResult struct {
	Extended struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"Extended"`
	FinanceOptionId               string `xml:"FinanceOptionId"`
	ID                            string `xml:"Id"`
	ModelId                       string `xml:"ModelId"`
	OrderId                       string `xml:"OrderId"`
	OrderLineId                   string `xml:"OrderLineId"`
	PalletId                      string `xml:"PalletId"`
	Reason                        string `xml:"Reason"`
	ReceiveExchangeDeviceDetailId string `xml:"ReceiveExchangeDeviceDetailId"`
	ReceiveReturnedDeviceDetailId string `xml:"ReceiveReturnedDeviceDetailId"`
	ShipDate                      string `xml:"ShipDate"`
	StockHandlerId                string `xml:"StockHandlerId"`
	StockReceiveDetailsId         string `xml:"StockReceiveDetailsId"`
	StockTakeHeaderId             string `xml:"StockTakeHeaderId"`
	TotalAccepted                 string `xml:"TotalAccepted"`
	TotalFailed                   string `xml:"TotalFailed"`
	TransactionType               string `xml:"TransactionType"`
	UseRange                      string `xml:"UseRange"`
}

type responseAddDeviceToBuildList struct {
	XMLName                            xml.Name                           `xml:"AddDeviceToBuildListManuallyResponse"`
	AddDeviceToBuildListManuallyResult addDeviceToBuildListManuallyResult `xml:"AddDeviceToBuildListManuallyResult"`
}

type addDeviceToBuildListManuallyResult struct {
	Accepted    string `xml:"Accepted"`
	BuildListId string `xml:"BuildListId"`
	DeviceId    string `xml:"DeviceId"`
	Error       string `xml:"Error"`
	Extended    struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"Extended"`
	FinanceOptionId    string `xml:"FinanceOptionId"`
	FromStockHandlerId string `xml:"FromStockHandlerId"`
	ID                 string `xml:"Id"`
	MACAddress1        string `xml:"MACAddress1"`
	MACAddress2        string `xml:"MACAddress2"`
	ModelId            string `xml:"ModelId"`
	OrderId            string `xml:"OrderId"`
	PalletId           string `xml:"PalletId"`
	Reason             string `xml:"Reason"`
	ReceiveId          string `xml:"ReceiveId"`
	SerialNumber       string `xml:"SerialNumber"`
	StatusId           string `xml:"StatusId"`
	StockHandlerId     string `xml:"StockHandlerId"`
}

type responsePerformBuildList struct {
	XMLName                      xml.Name                     `xml:"PerformBuildListActionResponse"`
	PerformBuildListActionResult performBuildListActionResult `xml:"PerformBuildListActionResult"`
}

type performBuildListActionResult struct {
	Extended struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"Extended"`
	FinanceOptionId               string `xml:"FinanceOptionId"`
	ID                            string `xml:"Id"`
	ModelId                       string `xml:"ModelId"`
	OrderId                       string `xml:"OrderId"`
	OrderLineId                   string `xml:"OrderLineId"`
	PalletId                      string `xml:"PalletId"`
	Reason                        string `xml:"Reason"`
	ReceiveExchangeDeviceDetailId string `xml:"ReceiveExchangeDeviceDetailId"`
	ReceiveReturnedDeviceDetailId string `xml:"ReceiveReturnedDeviceDetailId"`
	ShipDate                      string `xml:"ShipDate"`
	StockHandlerId                string `xml:"StockHandlerId"`
	StockReceiveDetailsId         string `xml:"StockReceiveDetailsId"`
	StockTakeHeaderId             string `xml:"StockTakeHeaderId"`
	TotalAccepted                 string `xml:"TotalAccepted"`
	TotalFailed                   string `xml:"TotalFailed"`
	TransactionType               string `xml:"TransactionType"`
	UseRange                      string `xml:"UseRange"`
}

// Get Device By SN
func GetDeviceBySerialNumber(iSN string) st.Device {
	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSDevice"
	l.FunctionName = "GetDeviceBySerialNumber"
	l.Request = "SN=" + iSN
	l.Start = t0.Format(time.RFC3339Nano)
	l.InsertappLog("./log/tvsdeviceapplog.log", "GetDeviceBySerialNumber")

	var odv st.Device

	// Database
	//db, err := sql.Open("goracle", "bgweb/bgweb#1@//tv-uat62-dq.tvsit.co.th:1521/UAT62")
	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
		resp = err.Error()
	} else {

		defer db.Close()
		var statement string
		statement = "begin redibsservice.getdatadevice62(:0,:1,:2,:3,:4); end;"
		var resultC driver.Rows
		var resultAgentKey, resultAgentName, resultReturnDate string

		if _, err := db.Exec(statement, iSN, sql.Out{Dest: &resultC}, sql.Out{Dest: &resultAgentKey}, sql.Out{Dest: &resultAgentName}, sql.Out{Dest: &resultReturnDate}); err != nil {
			log.Fatal(err)
			resp = err.Error()
		}

		defer resultC.Close()
		values := make([]driver.Value, len(resultC.Columns()))

		for {
			err = resultC.Next(values)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("error:", err)
				resp = err.Error()
			}

			odv.DeviceID = values[0].(int64)
			odv.SerialNumber = values[1].(string)
			odv.StatusID = values[2].(int64)
			odv.StatusDesc = values[3].(string)
			odv.ModelID = values[4].(int64)
			odv.ModelDesc = values[5].(string)
			odv.ProductID = values[6].(int64)
			odv.ProductDesc = values[7].(string)
			odv.StockhandlerID = values[8].(int64)
			odv.StockhandlerDesc = values[9].(string)
			odv.AllowSystem = values[10].(string)
			odv.FactoryWarrantyDate = values[11].(string)
			odv.AgentKey = resultAgentKey
			odv.AgentName = resultAgentName
			odv.ReturnDate = resultReturnDate
		}

		// Log#Stop
		t1 := time.Now()
		t2 := t1.Sub(t0)
		l.TrackingNo = trackingno
		l.ApplicationName = "TVSDevice"
		l.FunctionName = "GetDeviceBySerialNumber"
		l.Request = "SN=" + iSN
		l.Response = resp
		l.Start = t0.Format(time.RFC3339Nano)
		l.End = t1.Format(time.RFC3339Nano)
		l.Duration = t2.String()
		l.InsertappLog("./log/tvsdeviceapplog.log", "GetDeviceBySerialNumber")
	}
	return odv
}

// Device : ICC API
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

const getTemplateforPairingDevice = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
  <PairOneDeviceToAnother xmlns="http://ibs.entriq.net/Devices">
	<deviceFromId>$devicefromid</deviceFromId>
	<deviceToId>$devicetoid</deviceToId>
	<reason>$reason</reason>
  </PairOneDeviceToAnother>
</s:Body>
</s:Envelope>`

func PairOneDeviceToAnother(deviceFromid int64, deviceToId int64, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult

	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	url := ServiceLnk.DeviceURL
	client := &http.Client{}
	p(deviceFromid)
	p(deviceToId)

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)
	if len(strings.Trim(byusername, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + byusername + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}
	requestValue := s.Replace(getTemplateforPairingDevice, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$devicefromid", strconv.FormatInt(int64(deviceFromid), 10), -1)
	requestValue = s.Replace(requestValue, "$devicetoid", strconv.FormatInt(int64(deviceToId), 10), -1)
	requestValue = s.Replace(requestValue, "$reason", strconv.FormatInt(int64(reasonnr), 10), -1)

	//p(requestValue)
	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/PairOneDeviceToAnother"`)
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

/*
func main() {
	var strVal string
	fmt.Printf("input : ")
	fmt.Scan(&strVal)
	r := GetDeviceBySerialNumber(strVal)
	fmt.Println(r)

	json.NewEncoder(os.Stdout).Encode(r)
}
*/

const getTemplateforMoveDevice = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
  <MoveDevice xmlns="http://ibs.entriq.net/Devices">
	<deviceId>$deviceid</deviceId>
	<stockhandlerId>$stockhandlerid</stockhandlerId>
	<palletId>0</palletId>
	<reasonId>$reason</reasonId>
	<effectiveDate>$effectivedate</effectiveDate>
  </MoveDevice>
</s:Body>
</s:Envelope>`

// MoveDevice
func MoveDevice(iSN string, iDepotTo int64, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult
	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")
	//p(ICCAuthen)
	//fmt.Println(ServiceLnk)

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	p(token)
	//fmt.Scanln()

	dv := GetDeviceBySerialNumber(iSN)
	//p(dv)

	url := ServiceLnk.DeviceURL
	client := &http.Client{}

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)

	if len(strings.Trim(byusername, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + byusername + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}

	requestValue := s.Replace(getTemplateforMoveDevice, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$deviceid", strconv.FormatInt(int64(dv.DeviceID), 10), -1)
	requestValue = s.Replace(requestValue, "$stockhandlerid", strconv.FormatInt(int64(iDepotTo), 10), -1)
	requestValue = s.Replace(requestValue, "$reason", strconv.FormatInt(int64(reasonnr), 10), -1)
	requestValue = s.Replace(requestValue, "$effectivedate", time.Now().Format("2006-01-02T15:04:05"), -1)

	//p(requestValue)
	//fmt.Scanln()
	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/MoveDevice"`)
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
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	if resp.StatusCode != 200 {
		oRes.ErrorCode = 600
		oRes.ErrorDesc = myResult.Body.Fault.String
	} else {
		oRes.ErrorCode = 1
		oRes.ErrorDesc = "SUCCESS"
	}

	return oRes
}

const getTemplateforsendcmd = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
	<SendCommandToDevice xmlns="http://ibs.entriq.net/Devices">
		<deviceId>$deviceid</deviceId>
		<reason>$reason</reason>
	</SendCommandToDevice>
</s:Body>
</s:Envelope>`

// SendCommandToDevice
func SendCommandToDevice(iSN string, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult
	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	dv := GetDeviceBySerialNumber(iSN)

	url := ServiceLnk.DeviceURL
	client := &http.Client{}

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)

	if len(strings.Trim(byusername, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + byusername + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}

	requestValue := s.Replace(getTemplateforsendcmd, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$deviceid", strconv.FormatInt(int64(dv.DeviceID), 10), -1)
	requestValue = s.Replace(requestValue, "$reason", strconv.FormatInt(int64(reasonnr), 10), -1)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/SendCommandToDevice"`)
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
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	if resp.StatusCode != 200 {
		oRes.ErrorCode = 600
		oRes.ErrorDesc = myResult.Body.Fault.String
	} else {
		oRes.ErrorCode = 1
		oRes.ErrorDesc = "SUCCESS"
	}

	return oRes
}

// CreateNewSerialNumber
func CreateNewSerialNumber(iNewSNs st.NewDeviceReq) (st.ResponseResult, []st.NewDeviceRes) {
	var oRes st.ResponseResult
	var oSNRes []st.NewDeviceRes

	// 1. Get Token
	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = err.Error()
		return oRes, oSNRes
	}

	url := ServiceLnk.DeviceURL

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)

	if len(strings.Trim(iNewSNs.ByUser, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + iNewSNs.ByUser + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}

	// 2. createStockReceiveDetails
	oRes = CreateStockReceiveDetails(requestHD, url, iNewSNs.StockReceiveDetail, iNewSNs.Reason, iNewSNs.ByUser)
	if oRes.ErrorCode != 1 {
		return oRes, oSNRes
	}

	// 3. CreateBuildList
	iNewSNs.StockReceiveDetail.StockReceiveDetailsID = int64(oRes.CustomNum)
	oRes = CreateBuildList(requestHD, url, iNewSNs.StockReceiveDetail, iNewSNs.Reason, iNewSNs.ByUser)
	if oRes.ErrorCode != 1 {
		return oRes, oSNRes
	}

	// 4. AddDeviceToBuildListManually
	buildlstid := int64(oRes.CustomNum)
	for i := 0; i < len(iNewSNs.SerialNumber); i++ {
		oRes = AddDeviceToBuildListManually(requestHD, url, buildlstid, iNewSNs.SerialNumber[i], iNewSNs.ByUser)
		if oRes.ErrorCode == 1 { // Success
			var iSNRes st.NewDeviceRes
			iSNRes.SerialNumber = iNewSNs.SerialNumber[i]
			iSNRes.ResultCode = oRes.ErrorCode
			iSNRes.ResultDesc = oRes.ErrorDesc
			oSNRes = append(oSNRes, iSNRes)
			p(iSNRes)
		} else {
			var iSNRes st.NewDeviceRes
			iSNRes.SerialNumber = iNewSNs.SerialNumber[i]
			iSNRes.ResultCode = oRes.ErrorCode
			iSNRes.ResultDesc = oRes.ErrorDesc
			oSNRes = append(oSNRes, iSNRes)
			p(iSNRes)
		}
	}

	// 5. PerformBuildListAction
	oRes = PerformBuildListAction(requestHD, url, buildlstid, iNewSNs.ByUser)

	return oRes, oSNRes
}

const getTemplateforcreatestock = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
	<CreateStockReceiveDetails xmlns="http://ibs.entriq.net/Devices">
		<stockReceiveDetails xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
			<BatchComment i:nil="true" />
			<BatchNumber>$batchnumber</BatchNumber>
			<BatchReferenceNumber i:nil="true" />
			<DepreciationDetail i:nil="true" />
			<Extended i:nil="true" />
			<FromStockHanderId>$fromdepotid</FromStockHanderId>
			<Id i:nil="true" />
			<MACAddress1 i:nil="true" />
			<MACAddress2 i:nil="true" />
			<OrderId i:nil="true" />
			<PalletId i:nil="true" />
			<Reason>$reason</Reason>
			<ReorderReason i:nil="true" />
			<ToStockHanderId>$todepotid</ToStockHanderId>
			<UseRangeToDetermineModel i:nil="true" />
			<WarrantyEndDate>$wrenddate</WarrantyEndDate>
		</stockReceiveDetails>
	</CreateStockReceiveDetails>
</s:Body>
</s:Envelope>`

func CreateStockReceiveDetails(requestHD string, url string, iST st.StockReceiveDetails, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult
	var stockrecid int

	client := &http.Client{}

	p(reasonnr)
	// 2. XML for CreateStockReceiveDetails
	requestValue := s.Replace(getTemplateforcreatestock, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$batchnumber", iST.BatchNumber, -1)
	requestValue = s.Replace(requestValue, "$fromdepotid", strconv.FormatInt(int64(iST.FromStockHanderId), 10), -1)
	requestValue = s.Replace(requestValue, "$todepotid", strconv.FormatInt(int64(iST.ToStockHanderId), 10), -1)
	requestValue = s.Replace(requestValue, "$wrenddate", iST.WarrantyEndDate, -1)
	requestValue = s.Replace(requestValue, "$reason", strconv.FormatInt(int64(reasonnr), 10), -1)

	requestContent := []byte(requestValue)
	//p(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/CreateStockReceiveDetails"`)
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
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	//p(string(contents))

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	if resp.StatusCode != 200 {
		oRes.ErrorCode = 600
		oRes.ErrorDesc = myResult.Body.Fault.String
	} else {
		stockrecid, _ = strconv.Atoi(myResult.Body.ResponseCreateStockRecv.CreateStockReceiveDetailsResult.ID)
		oRes.ErrorCode = 1
		oRes.ErrorDesc = "SUCCESS"
		oRes.CustomNum = stockrecid
	}

	return oRes
}

const getTemplateforcreatebuildlist = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
	<CreateBuildList xmlns="http://ibs.entriq.net/Devices">
		<buildList xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
			<Extended i:nil="true" />
			<FinanceOptionId i:nil="true" />
			<Id i:nil="true" />
			<ModelId>$modelid</ModelId>
			<OrderId i:nil="true" />
			<OrderLineId i:nil="true" />
			<PalletId i:nil="true" />
			<Reason>$reason</Reason>
			<ReceiveExchangeDeviceDetailId i:nil="true" />
			<ReceiveReturnedDeviceDetailId i:nil="true" />
			<ShipDate i:nil="true" />
			<StockHandlerId i:nil="true" />
			<StockReceiveDetailsId>$stockreceiveid</StockReceiveDetailsId>
			<StockTakeHeaderId i:nil="true" />
			<TotalAccepted i:nil="true" />
			<TotalFailed i:nil="true" />
			<TransactionType>ReceiveNewStock</TransactionType>
			<UseRange i:nil="true" />
		</buildList>
	</CreateBuildList>
</s:Body>
</s:Envelope>`

func CreateBuildList(requestHD string, url string, iST st.StockReceiveDetails, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult
	var buildlstid int

	client := &http.Client{}

	// 3. CreateBuildList
	requestValue := s.Replace(getTemplateforcreatebuildlist, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$modelid", strconv.FormatInt(int64(iST.IBSModelId), 10), -1)
	requestValue = s.Replace(requestValue, "$stockreceiveid", strconv.FormatInt(int64(iST.StockReceiveDetailsID), 10), -1)
	requestValue = s.Replace(requestValue, "$reason", strconv.FormatInt(int64(reasonnr), 10), -1)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/CreateBuildList"`)
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
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	//p(string(contents))

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	if resp.StatusCode != 200 {
		oRes.ErrorCode = 600
		oRes.ErrorDesc = myResult.Body.Fault.String
	} else {
		buildlstid, _ = strconv.Atoi(myResult.Body.ResponseCreateBuildList.CreateBuildListResult.ID)
		oRes.ErrorCode = 1
		oRes.ErrorDesc = "SUCCESS"
		oRes.CustomNum = buildlstid
	}

	return oRes
}

const getTemplateforadddevicetobuildlist = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
	<AddDeviceToBuildListManually xmlns="http://ibs.entriq.net/Devices">
		<buildListId>$blistid</buildListId>
		<serialNumber>$sn</serialNumber>
	</AddDeviceToBuildListManually>
</s:Body>
</s:Envelope>`

func AddDeviceToBuildListManually(requestHD string, url string, buildlstid int64, newsn string, byusername string) st.ResponseResult {
	var oRes st.ResponseResult

	client := &http.Client{}

	// 4. AddDeviceToBuildListManually
	requestValue := s.Replace(getTemplateforadddevicetobuildlist, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$blistid", strconv.FormatInt(buildlstid, 10), -1)
	requestValue = s.Replace(requestValue, "$sn", newsn, -1)

	requestContent := []byte(requestValue)
	//p(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/AddDeviceToBuildListManually"`)
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
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	//p(string(contents))

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	if resp.StatusCode != 200 {
		oRes.ErrorCode = 600
		oRes.ErrorDesc = myResult.Body.Fault.String
	} else {
		if myResult.Body.ResponseAddDeviceToBuildList.AddDeviceToBuildListManuallyResult.Accepted != "true" {
			oRes.ErrorCode = 600
			oRes.ErrorDesc = myResult.Body.ResponseAddDeviceToBuildList.AddDeviceToBuildListManuallyResult.Error
			return oRes
		} else {
			oRes.ErrorCode = 1
			oRes.ErrorDesc = "SUCCESS"
		}
	}
	return oRes
}

const getTemplateforperformbuildlist = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
  <s:Body>
    <PerformBuildListAction xmlns="http://ibs.entriq.net/Devices">
      <buildListId>$blistid</buildListId>
    </PerformBuildListAction>
  </s:Body>
</s:Envelope>`

func PerformBuildListAction(requestHD string, url string, buildlstid int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult

	client := &http.Client{}

	// 5. PerformBuildListAction
	requestValue := s.Replace(getTemplateforperformbuildlist, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$blistid", strconv.FormatInt(buildlstid, 10), -1)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/PerformBuildListAction"`)
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
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	//p(string(contents))

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	if resp.StatusCode != 200 {
		oRes.ErrorCode = 600
		oRes.ErrorDesc = myResult.Body.Fault.String
	} else {
		oRes.ErrorCode = 1
		oRes.ErrorDesc = "SUCCESS"
		oRes.CustomNum, _ = strconv.Atoi(myResult.Body.ResponsePerformBuildList.PerformBuildListActionResult.ID)
	}

	return oRes
}

/*
func createStockReceiveDetailsA(requestHD string, url string, iST st.StockReceiveDetails, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult
	var stockrecid int64
	var buildlstid int64

	client := &http.Client{}

	// 2. XML for CreateStockReceiveDetails
	requestValue := s.Replace(getTemplateforcreatestock, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$batchnumber", iST.BatchNumber, -1)
	requestValue = s.Replace(requestValue, "$fromdepotid", strconv.FormatInt(int64(iST.FromStockHanderId), 10), -1)
	requestValue = s.Replace(requestValue, "$todepotid", strconv.FormatInt(int64(iST.ToStockHanderId), 10), -1)
	requestValue = s.Replace(requestValue, "$wrenddate", iST.WarrantyEndDate, -1)
	requestValue = s.Replace(requestValue, "$reason", strconv.FormatInt(int64(reasonnr), 10), -1)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/CreateStockReceiveDetails"`)
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
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	//p(string(contents))

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	if resp.StatusCode != 200 {
		oRes.ErrorCode = 600
		oRes.ErrorDesc = myResult.Body.Fault.String
	} else {
		stockrecid, _ = strconv.ParseInt(myResult.Body.ResponseCreateStockRecv.CreateStockReceiveDetailsResult.ID, 10, 64)
	}

	// 3. CreateBuildList
	requestValue = s.Replace(getTemplateforcreatebuildlist, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$modelid", strconv.FormatInt(int64(iST.IBSModelId), 10), -1)
	requestValue = s.Replace(requestValue, "$stockreceiveid", strconv.FormatInt(stockrecid, 10), -1)
	requestValue = s.Replace(requestValue, "$reason", strconv.FormatInt(int64(reasonnr), 10), -1)

	requestContent = []byte(requestValue)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/CreateBuildList"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err = client.Do(req)
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		oRes.ErrorCode = resp.StatusCode
		oRes.ErrorDesc = resp.Status
	}
	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	//p(string(contents))

	myResult = MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	if resp.StatusCode != 200 {
		oRes.ErrorCode = 600
		oRes.ErrorDesc = myResult.Body.Fault.String
	} else {
		buildlstid, _ = strconv.ParseInt(myResult.Body.ResponseCreateBuildList.CreateBuildListResult.ID, 10, 64)
	}

	// 4. AddDeviceToBuildListManually
	requestValue = s.Replace(getTemplateforadddevicetobuildlist, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$blistid", strconv.FormatInt(buildlstid, 10), -1)
	requestValue = s.Replace(requestValue, "$sn", iST.SerialNumber, -1)

	requestContent = []byte(requestValue)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/AddDeviceToBuildListManually"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err = client.Do(req)
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		oRes.ErrorCode = resp.StatusCode
		oRes.ErrorDesc = resp.Status
	}
	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	//p(string(contents))

	myResult = MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	if resp.StatusCode != 200 {
		oRes.ErrorCode = 600
		oRes.ErrorDesc = myResult.Body.Fault.String
	} else {
		if myResult.Body.ResponseAddDeviceToBuildList.AddDeviceToBuildListManuallyResult.Accepted != "true" {
			oRes.ErrorCode = 600
			oRes.ErrorDesc = myResult.Body.ResponseAddDeviceToBuildList.AddDeviceToBuildListManuallyResult.Error
			return oRes
		}
	}

	// 5. PerformBuildListAction
	requestValue = s.Replace(getTemplateforperformbuildlist, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$blistid", strconv.FormatInt(buildlstid, 10), -1)

	requestContent = []byte(requestValue)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/PerformBuildListAction"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err = client.Do(req)
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		oRes.ErrorCode = resp.StatusCode
		oRes.ErrorDesc = resp.Status
	}
	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	//p(string(contents))

	myResult = MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	if resp.StatusCode != 200 {
		oRes.ErrorCode = 600
		oRes.ErrorDesc = myResult.Body.Fault.String
	} else {
		oRes.ErrorCode = 1
		oRes.ErrorDesc = "SUCCESS"
		oRes.CustomNum, _ = strconv.Atoi(myResult.Body.ResponsePerformBuildList.PerformBuildListActionResult.ID)
	}
	return oRes
}
*/
