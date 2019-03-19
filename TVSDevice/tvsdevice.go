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
	"strings"
	"time"

	_ "gopkg.in/goracle.v2"

	s "strings"

	cm "github.com/smsdevteam/tvsglobal/common"     // db
	st "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
)

var p = log.Println

// MyRespEnvelope
type MyRespEnvelope struct {
	XMLName xml.Name
	Body    body
}

type body struct {
	XMLName                      xml.Name
	Fault                        fault
	ResponseGetDevice            responseGetDevice            `xml:"GetDeviceBySerialNumberResponse"`
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

type responseGetDevice struct {
	XMLName         xml.Name        `xml:"GetDeviceBySerialNumberResponse"`
	GetDeviceResult getDeviceResult `xml:"GetDeviceBySerialNumberResult"`
}

type getDeviceResult struct {
	CaReferenceNumber string `xml:"CaReferenceNumber"`
	CustomFields      struct {
		Text string `xml:",chardata"`
		A    string `xml:"a,attr"`
	} `xml:"CustomFields"`
	Extended struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"Extended"`
	ExternalId            string `xml:"ExternalId"`
	FinanceOptionId       string `xml:"FinanceOptionId"`
	FromStockHandlerId    string `xml:"FromStockHandlerId"`
	ID                    string `xml:"Id"`
	MACAddress1           string `xml:"MACAddress1"`
	MACAddress2           string `xml:"MACAddress2"`
	ModelId               string `xml:"ModelId"`
	OrderId               string `xml:"OrderId"`
	PalletId              string `xml:"PalletId"`
	Provisioned           string `xml:"Provisioned"`
	SerialNumber          string `xml:"SerialNumber"`
	ShipDate              string `xml:"ShipDate"`
	StatusId              string `xml:"StatusId"`
	StockHandlerId        string `xml:"StockHandlerId"`
	StockReceiveDetailsId string `xml:"StockReceiveDetailsId"`
	WarrantyEndDate       string `xml:"WarrantyEndDate"`
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

// GetDataSerialNumber
func GetDataSerialNumber(iSN string) st.DeviceData {
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

	var odv st.DeviceData

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

// IfThenElse evaluates a condition, if true returns the first parameter otherwise the second
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

// DefaultIfNil checks if the value is nil, if true returns the default value otherwise the original
func DefaultIfNil(value interface{}, defaultValue interface{}) interface{} {
	if value != nil {
		return value
	}
	return defaultValue
}

// GetDeviceViewBySerialNumber
func GetDeviceViewBySerialNumber(iSN string, iChipID string, iCustNr string) st.DeviceInfo {
	var odv st.DeviceInfo
	var resp string
	resp = "SUCCESS"

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
		statement = "begin tvsgo_product.GetDeviceBySerialNumber(:0,:1,:2,:3); end;"
		var resultC driver.Rows

		if _, err := db.Exec(statement, iSN, iChipID, iCustNr, sql.Out{Dest: &resultC}); err != nil {
			log.Fatal(err)
			resp = err.Error()
		}

		defer resultC.Close()
		values := make([]driver.Value, len(resultC.Columns()))

		colmap := cm.Createmapcol(resultC.Columns())

		for {
			err = resultC.Next(values)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("error:", err)
				resp = err.Error()
			}

			odv.ID = values[cm.Getcolindex(colmap, "ID")].(int64)
			odv.Serial_Number = values[cm.Getcolindex(colmap, "SERIAL_NUMBER")].(string)
			odv.Status_ID = values[cm.Getcolindex(colmap, "STATUS_ID")].(int64)
			odv.StatusDesc = values[cm.Getcolindex(colmap, "STATUSDESC")].(string)
			odv.Stock_HandlerID = values[cm.Getcolindex(colmap, "STOCK_HANDLERID")].(int64)
			odv.Stock_HandlerName = values[cm.Getcolindex(colmap, "STOCK_HANDLERNAME")].(string)
			odv.Model_ID = values[cm.Getcolindex(colmap, "MODEL_ID")].(int64)
			odv.Model_Desc = values[cm.Getcolindex(colmap, "MODEL_DESC")].(string)
			odv.Technical_Product_ID = values[cm.Getcolindex(colmap, "TECHNICAL_PRODUCT_ID")].(int64)
			odv.Technical_Product_Desc = values[cm.Getcolindex(colmap, "TECHNICAL_PRODUCT_DESC")].(string)
			odv.Technical_Product_Type = values[cm.Getcolindex(colmap, "TECHNICAL_PRODUCT_TYPE")].(string)
			odv.Names = values[cm.Getcolindex(colmap, "NAME")].(string)
			odv.Company = values[cm.Getcolindex(colmap, "COMPANY")].(string)
			odv.CustType = values[cm.Getcolindex(colmap, "CUSTTYPE")].(string)
			odv.SiliconFlag = values[cm.Getcolindex(colmap, "SILICONFLAG")].(string)
			odv.Duallnbf = values[cm.Getcolindex(colmap, "DUALLNBF")].(string)
			odv.Mac_Address1 = values[cm.Getcolindex(colmap, "MAC_ADDRESS1")].(string)
			odv.External_ID = values[cm.Getcolindex(colmap, "EXTERNAL_ID")].(string)
			odv.FinOption = values[cm.Getcolindex(colmap, "FINOPTION")].(string)
			odv.CustomerID = cm.StrToInt64(values[cm.Getcolindex(colmap, "CUSTOMER_ID")].(string))
			odv.DescLinkBasics = values[cm.Getcolindex(colmap, "DESCLINKBASICS")].(string)
			odv.Batch_number = values[cm.Getcolindex(colmap, "BATCH_NUMBER")].(string)
		}
	}
	p(resp)
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

// getAuthenHD : Replace header string
func getAuthenHD(iExtAgent string) (string, cm.ServiceURL) {
	var ICCAuthen cm.ICCAuthenHD
	var ServiceLnk cm.ServiceURL
	ICCAuthen, ServiceLnk = cm.ICCReadConfig("ICC")

	token, err := cm.GetICCAuthenToken("ICC")
	if err != nil {
		return "", ServiceLnk
	}

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)
	if len(strings.Trim(iExtAgent, " ")) != 0 {
		extAgentTag := "<h:ExternalAgent>" + iExtAgent + "</h:ExternalAgent>"
		requestHD = s.Replace(requestHD, `<h:ExternalAgent i:nil="true" />`, extAgentTag, -1)
	}

	return requestHD, ServiceLnk
}

// ----------------- GetDeviceBySerialNumber
const getTemplateforGetDeviceBySerialNumber = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
	<GetDeviceBySerialNumber xmlns="http://ibs.entriq.net/Devices">
		<serialNumber>$sn</serialNumber>
	</GetDeviceBySerialNumber>
</s:Body>
</s:Envelope>`

// GetDeviceBySerialNumber ---
func GetDeviceBySerialNumber(iSN string, iExtAgent string) st.GetDeviceResponse {
	var sn st.Device
	var oRes st.ResponseResult
	var result st.GetDeviceResponse

	requestHD, ServiceLnk := getAuthenHD(iExtAgent)
	if requestHD == "" {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = "Cannot get token from ICC API"
		result.ProcessResult = oRes
		return result
	}

	url := ServiceLnk.DeviceURL
	client := &http.Client{}

	requestValue := s.Replace(getTemplateforGetDeviceBySerialNumber, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$sn", iSN, -1)

	//p(requestValue)
	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		result.ProcessResult = oRes
		return result
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/GetDeviceBySerialNumber"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		oRes.ErrorCode = 300
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

	//p(string(contents))

	myRes := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myRes)
	sn.CaReferenceNumber = cm.StrToInt(myRes.Body.ResponseGetDevice.GetDeviceResult.CaReferenceNumber)
	sn.ExternalID = myRes.Body.ResponseGetDevice.GetDeviceResult.ExternalId
	sn.FinanceOptionID = cm.StrToInt(myRes.Body.ResponseGetDevice.GetDeviceResult.FinanceOptionId)
	sn.FromStockHandlerID = cm.StrToInt(myRes.Body.ResponseGetDevice.GetDeviceResult.FromStockHandlerId)
	sn.ID = cm.StrToInt(myRes.Body.ResponseGetDevice.GetDeviceResult.ID)
	sn.MACAddress1 = myRes.Body.ResponseGetDevice.GetDeviceResult.MACAddress1
	sn.MACAddress2 = myRes.Body.ResponseGetDevice.GetDeviceResult.MACAddress2
	sn.ModelID = cm.StrToInt(myRes.Body.ResponseGetDevice.GetDeviceResult.ModelId)
	sn.OrderID = cm.StrToInt(myRes.Body.ResponseGetDevice.GetDeviceResult.OrderId)
	sn.PalletID = cm.StrToInt(myRes.Body.ResponseGetDevice.GetDeviceResult.PalletId)
	sn.Provisioned = cm.StrToBool(myRes.Body.ResponseGetDevice.GetDeviceResult.Provisioned)
	sn.SerialNumber = myRes.Body.ResponseGetDevice.GetDeviceResult.SerialNumber
	sn.ShipDate = myRes.Body.ResponseGetDevice.GetDeviceResult.ShipDate
	sn.StatusID = cm.StrToInt(myRes.Body.ResponseGetDevice.GetDeviceResult.StatusId)
	sn.StockHandlerID = cm.StrToInt(myRes.Body.ResponseGetDevice.GetDeviceResult.StockHandlerId)
	sn.StockReceiveDetailsID = cm.StrToInt(myRes.Body.ResponseGetDevice.GetDeviceResult.StockReceiveDetailsId)
	sn.WarrantyEndDate = myRes.Body.ResponseGetDevice.GetDeviceResult.WarrantyEndDate
	p(sn)

	oRes.ErrorCode = 0
	oRes.ErrorDesc = "SUCCESS"
	result.ProcessResult = oRes
	result.DeviceResult = sn

	return result
}

// ----------------- GetDeviceBySerialNumber

// ----------------- PairOneDeviceToAnother
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

// PairOneDeviceToAnother
func PairOneDeviceToAnother(deviceFromid int64, deviceToId int64, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult

	requestHD, ServiceLnk := getAuthenHD(byusername)
	if requestHD == "" {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = "Cannot get token from ICC API"
		return oRes
	}

	url := ServiceLnk.DeviceURL
	client := &http.Client{}
	p(deviceFromid)
	p(deviceToId)

	requestValue := s.Replace(getTemplateforPairingDevice, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$devicefromid", cm.Int64ToStr(deviceFromid), -1)
	requestValue = s.Replace(requestValue, "$devicetoid", cm.Int64ToStr(deviceToId), -1)
	requestValue = s.Replace(requestValue, "$reason", cm.Int64ToStr(reasonnr), -1)

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
	oRes.ErrorCode = cm.StrToInt(myResult.Body.Fault.Code)
	oRes.ErrorDesc = myResult.Body.Fault.String

	return oRes
}

// ----------------- PairOneDeviceToAnother

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
func MoveDevice(ideviceid int64, idepotto int64, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult
	requestHD, ServiceLnk := getAuthenHD(byusername)
	if requestHD == "" {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = "Cannot get token from ICC API"
		return oRes
	}

	url := ServiceLnk.DeviceURL
	client := &http.Client{}

	requestValue := s.Replace(getTemplateforMoveDevice, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$deviceid", cm.Int64ToStr(ideviceid), -1)
	requestValue = s.Replace(requestValue, "$stockhandlerid", cm.Int64ToStr(idepotto), -1)
	requestValue = s.Replace(requestValue, "$reason", cm.Int64ToStr(reasonnr), -1)
	requestValue = s.Replace(requestValue, "$effectivedate", time.Now().Format("2006-01-02T15:04:05"), -1)

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
		oRes.ErrorCode = 0
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
func SendCommandToDevice(deviceid int64, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult
	requestHD, ServiceLnk := getAuthenHD(byusername)
	if requestHD == "" {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = "Cannot get token from ICC API"
		return oRes
	}

	url := ServiceLnk.DeviceURL
	client := &http.Client{}

	requestValue := s.Replace(getTemplateforsendcmd, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$deviceid", cm.Int64ToStr(deviceid), -1)
	requestValue = s.Replace(requestValue, "$reason", cm.Int64ToStr(reasonnr), -1)

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
		oRes.ErrorCode = 0
		oRes.ErrorDesc = "SUCCESS"
	}

	return oRes
}

// CreateNewSerialNumber
func CreateNewSerialNumber(iNewSNs st.NewDeviceReq) st.NewDeviceRes {
	var oRes st.ResponseResult
	var oSNRes []st.CreateSNRes
	var oSNReturn st.NewDeviceRes

	// 1. Get Token
	requestHD, ServiceLnk := getAuthenHD(iNewSNs.ByUser)
	if requestHD == "" {
		oRes.ErrorCode = 100
		oRes.ErrorDesc = "Cannot get token from ICC API"
		oSNReturn.ProcessRes = oRes
		return oSNReturn
	}

	url := ServiceLnk.DeviceURL

	// 2. createStockReceiveDetails
	oRes = CreateStockReceiveDetails(requestHD, url, iNewSNs.StockReceiveDetail, iNewSNs.Reason, iNewSNs.ByUser)
	if oRes.ErrorCode != 1 {
		oSNReturn.ProcessRes = oRes
		return oSNReturn
	}

	// 3. CreateBuildList
	iNewSNs.StockReceiveDetail.StockReceiveDetailsID = int64(oRes.CustomNum)
	oRes = CreateBuildList(requestHD, url, iNewSNs.StockReceiveDetail, iNewSNs.Reason, iNewSNs.ByUser)
	if oRes.ErrorCode != 1 {
		oSNReturn.ProcessRes = oRes
		return oSNReturn
	}

	// 4. AddDeviceToBuildListManually
	buildlstid := int64(oRes.CustomNum)
	for i := 0; i < len(iNewSNs.SerialNumber); i++ {
		oRes = AddDeviceToBuildListManually(requestHD, url, buildlstid, iNewSNs.SerialNumber[i], iNewSNs.ByUser)
		if oRes.ErrorCode == 1 { // Success
			var iSNRes st.CreateSNRes
			iSNRes.SerialNumber = iNewSNs.SerialNumber[i]
			iSNRes.ResultCode = oRes.ErrorCode
			iSNRes.ResultDesc = oRes.ErrorDesc
			oSNRes = append(oSNRes, iSNRes)
			p(iSNRes)
		} else {
			var iSNRes st.CreateSNRes
			iSNRes.SerialNumber = iNewSNs.SerialNumber[i]
			iSNRes.ResultCode = oRes.ErrorCode
			iSNRes.ResultDesc = oRes.ErrorDesc
			oSNRes = append(oSNRes, iSNRes)
			p(iSNRes)
		}
	}

	// 5. PerformBuildListAction
	oRes = PerformBuildListAction(requestHD, url, buildlstid, iNewSNs.ByUser)

	oSNReturn.ProcessRes = oRes
	oSNReturn.NewSNRes = oSNRes

	return oSNReturn
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
	requestValue = s.Replace(requestValue, "$fromdepotid", cm.Int64ToStr(iST.FromStockHanderID), -1)
	requestValue = s.Replace(requestValue, "$todepotid", cm.Int64ToStr(iST.ToStockHanderID), -1)
	requestValue = s.Replace(requestValue, "$wrenddate", iST.WarrantyEndDate, -1)
	requestValue = s.Replace(requestValue, "$reason", cm.Int64ToStr(reasonnr), -1)

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
		stockrecid = cm.StrToInt(myResult.Body.ResponseCreateStockRecv.CreateStockReceiveDetailsResult.ID)
		oRes.ErrorCode = 0
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
	requestValue = s.Replace(requestValue, "$modelid", cm.Int64ToStr(iST.IBSModelID), -1)
	requestValue = s.Replace(requestValue, "$stockreceiveid", cm.Int64ToStr(iST.StockReceiveDetailsID), -1)
	requestValue = s.Replace(requestValue, "$reason", cm.Int64ToStr(reasonnr), -1)

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
		buildlstid = cm.StrToInt(myResult.Body.ResponseCreateBuildList.CreateBuildListResult.ID)
		oRes.ErrorCode = 0
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
	requestValue = s.Replace(requestValue, "$blistid", cm.Int64ToStr(buildlstid), -1)
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
			oRes.ErrorCode = 0
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
	requestValue = s.Replace(requestValue, "$blistid", cm.Int64ToStr(buildlstid), -1)

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
		oRes.ErrorCode = 0
		oRes.ErrorDesc = "SUCCESS"
		oRes.CustomNum = cm.StrToInt(myResult.Body.ResponsePerformBuildList.PerformBuildListActionResult.ID)
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

const getTemplatereceiveexchangedevice = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
  <s:Body>
    <ReceiveExchangeDevice xmlns="http://ibs.entriq.net/Devices">
      <deviceId>$deviceid</deviceId>
      <stockHandlerId>$stockhandlerid</stockHandlerId>
      <palletId>$palletid</palletId>
      <reasonId>$reason</reasonId>
      <deviceExchangeReasonId>$deviceexreason</deviceExchangeReasonId>
      <shipDate>$shipdate</shipDate>
    </ReceiveExchangeDevice>
  </s:Body>
</s:Envelope>`

func ReceiveExchangeDevice(requestHD string, url string, iReq st.ReceiveExchangeDeviceReq, byusername string) st.ResponseResult {
	var oRes st.ResponseResult

	client := &http.Client{}

	requestValue := s.Replace(getTemplatereceiveexchangedevice, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$deviceid", cm.Int64ToStr(iReq.DeviceID), -1)
	requestValue = s.Replace(requestValue, "$stockhandlerid", cm.Int64ToStr(iReq.StockHandlerID), -1)
	requestValue = s.Replace(requestValue, "$palletid", cm.Int64ToStr(iReq.PalletID), -1)
	requestValue = s.Replace(requestValue, "$reason", cm.Int64ToStr(iReq.ReasonID), -1)
	requestValue = s.Replace(requestValue, "$deviceexreason", cm.Int64ToStr(iReq.DeviceExchangeReasonID), -1)
	requestValue = s.Replace(requestValue, "$shipdate", iReq.ShipDate, -1)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Devices/IDevicesService/ReceiveExchangeDevice"`)
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
		oRes.ErrorCode = 0
		oRes.ErrorDesc = "SUCCESS"
	}

	return oRes
}
