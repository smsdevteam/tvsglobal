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

	cm "github.com/pimpina/tvsglobalb/Common"     // db
	st "github.com/pimpina/tvsglobalb/TVSStructs" // referpath
)

var p = fmt.Println

// MyRespEnvelope
type MyRespEnvelope struct {
	XMLName xml.Name
	Body    body
}

type body struct {
	XMLName     xml.Name
	Fault       *fault
	GetResponse completeResponse `xml:"AuthenticateByProofResponse"`
}

type fault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
	Detail string `xml:"detail"`
}

type completeResponse struct {
	XMLName xml.Name `xml:"AuthenticateByProofResponse"`
	//	MyResult authenHD `xml:"AuthenticateByProofResult "`
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
	l.Start = t0.String()
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
		l.Start = t0.String()
		l.End = t1.String()
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

	p(requestValue)
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
	//p(token)
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

	p(requestValue)
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
		//return oRes
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
	var p = fmt.Println
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
	p(url)
	client := &http.Client{}

	requestHD := s.Replace(getTemplateAuthenHD, "$username", ICCAuthen.ServiceUser, -1)
	requestHD = s.Replace(requestHD, "$password", ICCAuthen.ServiceUserIdentity, -1)
	requestHD = s.Replace(requestHD, "$dsn", ICCAuthen.ServiceDSN, -1)
	requestHD = s.Replace(requestHD, "$servicetime", time.Now().Format("2006-01-02T15:04:05"), -1)
	requestHD = s.Replace(requestHD, "$token", token, -1)

	requestValue := s.Replace(getTemplateforsendcmd, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$deviceid", strconv.FormatInt(int64(dv.DeviceID), 10), -1)
	requestValue = s.Replace(requestValue, "$reason", strconv.FormatInt(int64(reasonnr), 10), -1)

	//p(requestValue)

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
	//p(string(contents))

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
