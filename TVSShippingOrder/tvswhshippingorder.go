package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"io"
	"io/ioutil"
	"log"
	"encoding/xml"
	"net/http"
	s "strings"
	"time"
	_ "gopkg.in/goracle.v2"

	cm "github.com/smsdevteam/tvsglobal/Common"     // db
	st "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
)

// GetWHOrder Method
func GetWHOrder(iOrderID int64) st.ShippingOrderRes {
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

// CancelWHOrder Method
func CancelWHOrder(soid int64, reasonnr int64, byusername string) st.ResponseResult {
	var oRes st.ResponseResult
	oRes = CancelShippingOrder(soid, reasonnr, byusername)
	return oRes
}

// CreateWHOrder Method
func CreateWHOrder(iSO st.ShippingOrderReq, reasonnr int64, byusername string) st.SOResult {
	var oRes st.SOResult
	//oRes = CreateShippingOrder(iSO)
	return oRes
}

const getTemplateforshipwhso = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
$TemplateHD
<s:Body>
	<ShipOrderBetweenStockHandlers xmlns="http://ibs.entriq.net/OrderManagement">
		<order xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
			$sodata
		</order>
      <reasonId>$reason</reasonId>
	  <printers i:nil="true" xmlns:i="http://www.w3.org/2001/XMLSchema-instance" />		
	  <reasonIdForDeviceTransfer>$reasondevice</reasonIdForDeviceTransfer>	
	</ShipOrderBetweenStockHandlers>
</s:Body>
</s:Envelope>`

// ShipWHOrder Method
func ShipWHOrder(soid int64, reasonnr int64, reasonnrfordevice int64,byusername string) st.ResponseResult {
	var oRes st.ResponseResult
	var SOData st.SOResult
	SOData = GetShippingOrder(soid, byusername)

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
	requestValue := s.Replace(getTemplateforshipwhso, "$TemplateHD", requestHD, -1)
	requestValue = s.Replace(requestValue, "$sodata", so, -1)
	requestValue = s.Replace(requestValue, "$reason", cm.Int64ToStr(reasonnr), -1)
	requestValue = s.Replace(requestValue, "$reasondevice", cm.Int64ToStr(reasonnrfordevice), -1)

	//p(requestValue)

	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/OrderManagement/IOrderManagementService/ShipOrderBetweenStockHandlers"`)
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
