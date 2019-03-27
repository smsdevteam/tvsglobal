package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	s "strings"
	"time"

	_ "gopkg.in/goracle.v2"

	cm "github.com/smsdevteam/tvsglobal/common"     // db
	st "github.com/smsdevteam/tvsglobal/tvsstructs" // referpath
)

const applicationname string = "tvs-keyword"
const tagappname string = "icc-tvskeyword"
const taglogtype string = "applogs"

var tagenv = os.Getenv("ENVAPP")

// MyRespEnvelopeCreateKeyword obj
type MyRespEnvelopeCreateKeyword struct {
	XMLName xml.Name          `xml:"Envelope"`
	Body    bodyCreateKeyword `xml:"Body"`
}

//bodyCreateKeyword obj
type bodyCreateKeyword struct {
	XMLName                xml.Name              `xml:"Body"`
	VCreateKeywordResponse createKeywordResponse `xml:"CreateKeywordResponse"`
}

//createKeywordResponse obj
type createKeywordResponse struct {
	XMLName              xml.Name            `xml:"CreateKeywordResponse"`
	VCreateKeywordResult createKeywordResult `xml:"CreateKeywordResult"`
}

//createKeywordResult obj
type createKeywordResult struct {
	XMLName     xml.Name `xml:"CreateKeywordResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
}

// MyRespEnvelopeDeleteKeyword obj
type MyRespEnvelopeDeleteKeyword struct {
	XMLName xml.Name          `xml:"Envelope"`
	Body    bodyDeleteKeyword `xml:"Body"`
}

//bodyDeleteKeyword obj
type bodyDeleteKeyword struct {
	XMLName                xml.Name              `xml:"Body"`
	VDeleteKeywordResponse deleteKeywordResponse `xml:"DeleteKeywordResponse"`
}

//deleteKeywordResponse obj
type deleteKeywordResponse struct {
	XMLName              xml.Name            `xml:"DeleteKeywordResponse"`
	VDeleteKeywordResult deleteKeywordResult `xml:"DeleteKeywordResult"`
}

//deleteKeywordResult obj
type deleteKeywordResult struct {
	XMLName     xml.Name `xml:"DeleteKeywordResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
}

// MyRespEnvelopeUpdateKeyword obj
type MyRespEnvelopeUpdateKeyword struct {
	XMLName xml.Name          `xml:"Envelope"`
	Body    bodyUpdateKeyword `xml:"Body"`
}

//bodyUpdateKeyword obj
type bodyUpdateKeyword struct {
	XMLName                xml.Name              `xml:"Body"`
	VUpdateKeywordResponse updateKeywordResponse `xml:"UpdateKeywordResponse"`
}

//updateKeywordResponse obj
type updateKeywordResponse struct {
	XMLName              xml.Name            `xml:"UpdateKeywordResponse"`
	VUpdateKeywordResult updateKeywordResult `xml:"UpdateKeywordResult"`
}

//updateKeywordResult obj
type updateKeywordResult struct {
	XMLName     xml.Name `xml:"UpdateKeywordResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
}

//GetKeywordByKeywordID function
func GetKeywordByKeywordID(iKeywordID string) *st.GetKeywordResult {

	// Log#Start
	l := cm.NewApplog()
	defer l.PrintJSONLog()

	defer func() {
		if err := recover(); err != nil {
			error := fmt.Sprint(err)
			l.Response = error
			//fmt.Printf("Error func GetNoteByNoteID .. %s\n", err)
		}
	}()

	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format(time.RFC3339Nano)
	l.TrackingNo = trackingno
	l.ApplicationName = applicationname
	l.FunctionName = "GetKeyword"
	l.Request = "KeywordID=" + iKeywordID
	l.Start = t0.Format(time.RFC3339Nano)
	var tags []string
	tags = append(tags, tagenv)
	tags = append(tags, tagappname)
	tags = append(tags, taglogtype)
	l.Tags = tags
	//l.InsertappLog("./log/tvskeywordapplog.log", "GetKeyword")

	oRes := st.NewGetKeywordResult()
	var oKeyword st.Keyword

	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		//log.Println(err)
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 2
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer db.Close()
	var statement string
	statement = "begin PK_ICC_KEYWORD.GetKeyword(:0,:1); end;"
	var resultC driver.Rows
	intKeywordID, err := strconv.Atoi(iKeywordID)
	if err != nil {
		//log.Println(err)
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 3
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	if _, err := db.Exec(statement, intKeywordID, sql.Out{Dest: &resultC}); err != nil {
		//log.Println(err)
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 4
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer resultC.Close()
	values := make([]driver.Value, len(resultC.Columns()))
	for {
		colmap := cm.Createmapcol(resultC.Columns())
		log.Println(colmap)

		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			//log.Println(err)
			resp = err.Error()
			l.Response = resp
			oRes.ErrorCode = 5
			oRes.ErrorDesc = err.Error()
			return oRes
		}

		if values[cm.Getcolindex(colmap, "ALLOWED_KEYWORD_ID")] != nil {
			oKeyword.AllowedKeywordID = values[cm.Getcolindex(colmap, "ALLOWED_KEYWORD_ID")].(int64)
		}
		oKeyword.Attribute = values[cm.Getcolindex(colmap, "ATTRIBUTE")].(string)
		if values[cm.Getcolindex(colmap, "COUNT_VALUE")] != nil {
			oKeyword.CountValue = values[cm.Getcolindex(colmap, "COUNT_VALUE")].(int64)
		}
		if values[cm.Getcolindex(colmap, "CUSTOMER_ID")] != nil {
			oKeyword.CustomerID = values[cm.Getcolindex(colmap, "CUSTOMER_ID")].(int64)
		}
		if values[cm.Getcolindex(colmap, "DATE_VALUE")] != nil {
			oKeyword.DateValue = values[cm.Getcolindex(colmap, "DATE_VALUE")].(time.Time)
		}
		if values[cm.Getcolindex(colmap, "ID")] != nil {
			oKeyword.ID = values[cm.Getcolindex(colmap, "ID")].(int64)
		}
		oKeyword.KAAttribute = values[cm.Getcolindex(colmap, "KAATTRIBUTE")].(string)
		oKeyword.KAKeyword = values[cm.Getcolindex(colmap, "KAKEYWORD")].(string)
		oKeyword.KALongDescr = values[cm.Getcolindex(colmap, "KALONGDESCR")].(string)
		oKeyword.KTName = values[cm.Getcolindex(colmap, "KTNAME")].(string)
		oKeyword.KTUserKey = values[cm.Getcolindex(colmap, "KTUSERKEY")].(string)
		if values[cm.Getcolindex(colmap, "KEYTYPES_ID")] != nil {
			oKeyword.KeyTypesID = values[cm.Getcolindex(colmap, "KEYTYPES_ID")].(int64)
		}

	}
	oRes.MyKeyword = oKeyword
	if oRes.ErrorCode == 1 {
		oRes.ErrorCode = 0
		oRes.ErrorDesc = "Success"
	}
	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	//l.TrackingNo = trackingno
	//l.ApplicationName = "TVSKeyword"
	//l.FunctionName = "GetKeyword"
	//l.Request = "KeywordID=" + iKeywordID
	jSRes, _ := json.Marshal(oRes)
	sJSRes := string(jSRes)

	l.Response = sJSRes
	//l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	//l.InsertappLog("./log/tvskeywordapplog.log", "GetKeyword")
	//test

	return oRes
}

//GetListKeywordByCustomerID function
func GetListKeywordByCustomerID(iCustomerID string) *st.GetListKeywordResult {
	// Log#Start

	l := cm.NewApplog()
	defer l.PrintJSONLog()

	defer func() {
		if err := recover(); err != nil {
			error := fmt.Sprint(err)
			l.Response = error
			//fmt.Printf("Error func GetNoteByNoteID .. %s\n", err)
		}
	}()

	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format(time.RFC3339Nano)
	l.TrackingNo = trackingno
	l.ApplicationName = applicationname
	l.FunctionName = "GetListKeywordByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Start = t0.Format(time.RFC3339Nano)
	var tags []string
	tags = append(tags, tagenv)
	tags = append(tags, tagappname)
	tags = append(tags, taglogtype)
	l.Tags = tags

	oRes := st.NewGetListKeywordResult()
	var oListKeyword st.ListKeyword

	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		//log.Println(err)
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 2
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer db.Close()
	var statement string
	statement = "begin PK_ICC_KEYWORD.GetListKeywordByCustomerId(:0,:1); end;"
	var resultC driver.Rows
	intCustomerID, err := strconv.Atoi(iCustomerID)
	if err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 3
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	if _, err := db.Exec(statement, intCustomerID, sql.Out{Dest: &resultC}); err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 4
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer resultC.Close()
	values := make([]driver.Value, len(resultC.Columns()))
	var oLKeyword []st.Keyword
	for {
		colmap := cm.Createmapcol(resultC.Columns())

		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			resp = err.Error()
			l.Response = resp
			oRes.ErrorCode = 5
			oRes.ErrorDesc = err.Error()
			return oRes
		}

		var oKeyword st.Keyword
		if values[cm.Getcolindex(colmap, "ALLOWED_KEYWORD_ID")] != nil {
			oKeyword.AllowedKeywordID = values[cm.Getcolindex(colmap, "ALLOWED_KEYWORD_ID")].(int64)
		}
		oKeyword.Attribute = values[cm.Getcolindex(colmap, "ATTRIBUTE")].(string)
		if values[cm.Getcolindex(colmap, "COUNT_VALUE")] != nil {
			oKeyword.CountValue = values[cm.Getcolindex(colmap, "COUNT_VALUE")].(int64)
		}
		if values[cm.Getcolindex(colmap, "CUSTOMER_ID")] != nil {
			oKeyword.CustomerID = values[cm.Getcolindex(colmap, "CUSTOMER_ID")].(int64)
		}
		if values[cm.Getcolindex(colmap, "DATE_VALUE")] != nil {
			oKeyword.DateValue = values[cm.Getcolindex(colmap, "DATE_VALUE")].(time.Time)
		}
		if values[cm.Getcolindex(colmap, "ID")] != nil {
			oKeyword.ID = values[cm.Getcolindex(colmap, "ID")].(int64)
		}
		oKeyword.KAAttribute = values[cm.Getcolindex(colmap, "KAATTRIBUTE")].(string)
		oKeyword.KAKeyword = values[cm.Getcolindex(colmap, "KAKEYWORD")].(string)
		oKeyword.KALongDescr = values[cm.Getcolindex(colmap, "KALONGDESCR")].(string)
		oKeyword.KTName = values[cm.Getcolindex(colmap, "KTNAME")].(string)
		oKeyword.KTUserKey = values[cm.Getcolindex(colmap, "KTUSERKEY")].(string)
		if values[cm.Getcolindex(colmap, "KEYTYPES_ID")] != nil {
			oKeyword.KeyTypesID = values[cm.Getcolindex(colmap, "KEYTYPES_ID")].(int64)
		}

		oLKeyword = append(oLKeyword, oKeyword)

	}

	oListKeyword.Keywords = oLKeyword
	oRes.MyListKeyword = oListKeyword
	if oRes.ErrorCode == 1 {
		oRes.ErrorCode = 0
		oRes.ErrorDesc = "Success"
	}
	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	//l.TrackingNo = trackingno
	//l.ApplicationName = "TVSKeyword"
	//l.FunctionName = "GetListKeywordByCustomerID"
	//l.Request = "CustomerID=" + iCustomerID
	jSRes, _ := json.Marshal(oRes)
	sJSRes := string(jSRes)

	l.Response = sJSRes
	//l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	//l.InsertappLog("./log/tvskeywordapplog.log", "GetListKeywordByCustomerID")
	//test
	return oRes
}

//getTemplateforCreateKeyword is xmltemplate for post to ICC service
const getTemplateforCreateKeyword = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<CreateKeyword xmlns="http://tempuri.org/">
		<inKeyword>
			<Attribute>$attribute</Attribute>
			<Count>$count</Count>
			<CustomerId>$customerId</CustomerId>
			<Date>$date</Date>
            <Id>0</Id>
            <KeywordId>$keywordId</KeywordId>
            <KeywordTypeId>$keywordTypeId</KeywordTypeId>
			<Name>$name</Name>
			<Extended>$extended</Extended>
		</inKeyword>
		<inReason>$inReason</inReason>
		<byUser>
			<byUser>$byUser</byUser>
            <byChannel>$byChannel</byChannel>
            <byProject>$byProject</byProject>
            <byHost>$byHost</byHost>
		</byUser>
	</CreateKeyword>
</s:Body>
</s:Envelope>`

//CreateKeyword function
func CreateKeyword(iReq st.CreateKeywordRequest) *st.CreateKeywordResponse {

	// Log#Start

	l := cm.NewApplog()
	defer l.PrintJSONLog()

	defer func() {
		if err := recover(); err != nil {
			error := fmt.Sprint(err)
			l.Response = error
			//fmt.Printf("Error func GetNoteByNoteID .. %s\n", err)
		}
	}()

	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format(time.RFC3339Nano)
	l.TrackingNo = trackingno
	l.ApplicationName = applicationname
	l.FunctionName = "CreateKeyword"

	jSReq, _ := json.Marshal(iReq)
	sJSReq := string(jSReq)

	l.Request = sJSReq

	l.Start = t0.Format(time.RFC3339Nano)
	//l.InsertappLog("./log/tvskeywordapplog.log", "CreateKeyword")

	oRes := st.NewCreateKeywordResponse()

	var AppServiceLnk cm.AppServiceURL
	AppServiceLnk = cm.AppReadConfig("ENH")

	url := AppServiceLnk.ICCServiceURL
	client := &http.Client{}

	sCount := strconv.FormatInt(iReq.InKeyword.Count, 10)
	sCustomerID := strconv.FormatInt(iReq.InKeyword.CustomerID, 10)
	sDate := (iReq.InKeyword.Date).Format("2006-01-02T15:04:05")
	sKeywordID := strconv.FormatInt(iReq.InKeyword.KeywordID, 10)
	sKeywordTypeID := strconv.FormatInt(iReq.InKeyword.KeywordTypeID, 10)
	sReason := strconv.FormatInt(iReq.InReason, 10)

	requestValue := s.Replace(getTemplateforCreateKeyword, "$attribute", iReq.InKeyword.Attribute, -1)
	requestValue = s.Replace(requestValue, "$count", sCount, -1)
	requestValue = s.Replace(requestValue, "$customerId", sCustomerID, -1)
	requestValue = s.Replace(requestValue, "$date", sDate, -1)
	requestValue = s.Replace(requestValue, "$keywordId", sKeywordID, -1)
	requestValue = s.Replace(requestValue, "$keywordTypeId", sKeywordTypeID, -1)
	requestValue = s.Replace(requestValue, "$name", iReq.InKeyword.Name, -1)
	requestValue = s.Replace(requestValue, "$extended", iReq.InKeyword.Extended, -1)
	requestValue = s.Replace(requestValue, "$inReason", sReason, -1)
	requestValue = s.Replace(requestValue, "$byUser", iReq.ByUser.ByUser, -1)
	requestValue = s.Replace(requestValue, "$byChannel", iReq.ByUser.ByChannel, -1)
	requestValue = s.Replace(requestValue, "$byProject", iReq.ByUser.ByProject, -1)
	requestValue = s.Replace(requestValue, "$byHost", iReq.ByUser.ByHost, -1)

	//log.Println("requestValue: " + requestValue)
	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 200
		oRes.ErrorDesc = resp
		return oRes
	}

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/CreateKeyword"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	response, err := client.Do(req)
	if err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 200
		oRes.ErrorDesc = resp
		return oRes
	}
	defer response.Body.Close()

	//log.Println(response.Body)

	if response.StatusCode != 200 {
		resp = "error " + response.Status
		l.Response = resp
		oRes.ErrorCode = response.StatusCode
		oRes.ErrorDesc = response.Status
		return oRes
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 400
		oRes.ErrorDesc = resp
		return oRes
	}

	//log.Println("contents : " + string(contents[:]))

	myResult := MyRespEnvelopeCreateKeyword{}
	xml.Unmarshal([]byte(contents), &myResult)
	//log.Println(myResult)
	oRes.ResultValue = myResult.Body.VCreateKeywordResponse.VCreateKeywordResult.ResultValue
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.VCreateKeywordResponse.VCreateKeywordResult.ErrorCode)
	oRes.ErrorDesc = myResult.Body.VCreateKeywordResponse.VCreateKeywordResult.ErrorDesc

	//log.Println(oRes)
	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	//l.TrackingNo = trackingno
	//l.ApplicationName = "TVSKeyword"
	//l.FunctionName = "CreateKeyword"
	//l.Request = "ByUser=" + iReq.ByUser.ByUser
	jSRes, _ := json.Marshal(oRes)
	sJSRes := string(jSRes)

	l.Response = sJSRes
	//l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	//l.InsertappLog("./log/tvskeywordapplog.log", "CreateKeyword")
	return oRes
}

//getTemplateforDeleteKeyword is xmltemplate for post to ICC service
const getTemplateforDeleteKeyword = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<DeleteKeyword xmlns="http://tempuri.org/">
		<inKeywordId>$inKeywordId</inKeywordId>
		<inReason>$inReason</inReason>
		<byUser>
			<byUser>$byUser</byUser>
            <byChannel>$byChannel</byChannel>
            <byProject>$byProject</byProject>
            <byHost>$byHost</byHost>
		</byUser>
	</DeleteKeyword>
</s:Body>
</s:Envelope>`

//DeleteKeyword function
func DeleteKeyword(iReq st.DeleteKeywordRequest) *st.DeleteKeywordResponse {

	// Log#Start

	l := cm.NewApplog()
	defer l.PrintJSONLog()

	defer func() {
		if err := recover(); err != nil {
			error := fmt.Sprint(err)
			l.Response = error
			//fmt.Printf("Error func GetNoteByNoteID .. %s\n", err)
		}
	}()

	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format(time.RFC3339Nano)
	l.TrackingNo = trackingno
	l.ApplicationName = applicationname
	l.FunctionName = "DeleteKeyword"

	jSReq, _ := json.Marshal(iReq)
	sJSReq := string(jSReq)

	l.Request = sJSReq

	l.Start = t0.Format(time.RFC3339Nano)
	//l.InsertappLog("./log/tvskeywordapplog.log", "DeleteKeyword")

	oRes := st.NewDeleteKeywordResponse()
	var AppServiceLnk cm.AppServiceURL
	AppServiceLnk = cm.AppReadConfig("ENH")

	url := AppServiceLnk.ICCServiceURL
	client := &http.Client{}

	sInKeywordID := strconv.FormatInt(iReq.InKeywordID, 10)
	sInReason := strconv.FormatInt(iReq.InReason, 10)

	requestValue := s.Replace(getTemplateforDeleteKeyword, "$inKeywordId", sInKeywordID, -1)
	requestValue = s.Replace(requestValue, "$inReason", sInReason, -1)
	requestValue = s.Replace(requestValue, "$byUser", iReq.ByUser.ByUser, -1)
	requestValue = s.Replace(requestValue, "$byChannel", iReq.ByUser.ByChannel, -1)
	requestValue = s.Replace(requestValue, "$byProject", iReq.ByUser.ByProject, -1)
	requestValue = s.Replace(requestValue, "$byHost", iReq.ByUser.ByHost, -1)

	//log.Println("requestValue: " + requestValue)
	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 200
		oRes.ErrorDesc = resp
		return oRes
	}

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/DeleteKeyword"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	response, err := client.Do(req)
	if err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 200
		oRes.ErrorDesc = resp
		return oRes
	}
	defer response.Body.Close()

	//log.Println(response.Body)

	if response.StatusCode != 200 {
		resp = "error " + response.Status
		l.Response = resp
		oRes.ErrorCode = response.StatusCode
		oRes.ErrorDesc = response.Status
		return oRes
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 400
		oRes.ErrorDesc = resp
		return oRes
	}

	log.Println("contents : " + string(contents[:]))

	myResult := MyRespEnvelopeDeleteKeyword{}
	xml.Unmarshal([]byte(contents), &myResult)
	//log.Println(myResult)
	oRes.ResultValue = myResult.Body.VDeleteKeywordResponse.VDeleteKeywordResult.ResultValue
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.VDeleteKeywordResponse.VDeleteKeywordResult.ErrorCode)
	oRes.ErrorDesc = myResult.Body.VDeleteKeywordResponse.VDeleteKeywordResult.ErrorDesc

	//log.Println(oRes)

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	//l.TrackingNo = trackingno
	//l.ApplicationName = "TVSKeyword"
	//l.FunctionName = "DeleteKeyword"
	//l.Request = "ByUser=" + iReq.ByUser.ByUser
	l.Response = resp
	//l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	//l.InsertappLog("./log/tvskeywordapplog.log", "DeleteKeyword")

	return oRes
}

//getTemplateforUpdateKeyword is xmltemplate for post to ICC service
const getTemplateforUpdateKeyword = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<UpdateKeyword xmlns="http://tempuri.org/">
		<inKeyword>
			<Attribute>$attribute</Attribute>
			<Count>$count</Count>
			<CustomerId>$customerId</CustomerId>
			<Date>$date</Date>
			<Id>$id</Id>
			<KeywordId>$keywordId</KeywordId>
			<KeywordTypeId>$keywordTypeId</KeywordTypeId>
			<Name>$name</Name>
			<Extended>$extended</Extended>
		</inKeyword>
		<inReason>$inReason</inReason>
		<byUser>
			<byUser>$byUser</byUser>
			<byChannel>$byChannel</byChannel>
			<byProject>$byProject</byProject>
			<byHost>$byHost</byHost>
		</byUser>
	</UpdateKeyword>
</s:Body>
</s:Envelope>`

//UpdateKeyword function
func UpdateKeyword(iReq st.UpdateKeywordRequest) *st.UpdateKeywordResponse {

	// Log#Start

	l := cm.NewApplog()
	defer l.PrintJSONLog()

	defer func() {
		if err := recover(); err != nil {
			error := fmt.Sprint(err)
			l.Response = error
			//fmt.Printf("Error func GetNoteByNoteID .. %s\n", err)
		}
	}()

	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format(time.RFC3339Nano)
	l.TrackingNo = trackingno
	l.ApplicationName = applicationname
	l.FunctionName = "UpdateKeyword"

	jSReq, _ := json.Marshal(iReq)
	sJSReq := string(jSReq)

	l.Request = sJSReq

	l.Start = t0.Format(time.RFC3339Nano)
	//l.InsertappLog("./log/tvskeywordapplog.log", "UpdateKeyword")

	oRes := st.NewUpdateKeywordResponse()

	var AppServiceLnk cm.AppServiceURL
	AppServiceLnk = cm.AppReadConfig("ENH")

	url := AppServiceLnk.ICCServiceURL
	client := &http.Client{}

	sCount := strconv.FormatInt(iReq.InKeyword.Count, 10)
	sCustomerID := strconv.FormatInt(iReq.InKeyword.CustomerID, 10)
	sDate := (iReq.InKeyword.Date).Format("2006-01-02T15:04:05")
	sKeywordID := strconv.FormatInt(iReq.InKeyword.KeywordID, 10)
	sKeywordTypeID := strconv.FormatInt(iReq.InKeyword.KeywordTypeID, 10)
	sID := strconv.FormatInt(iReq.InKeyword.ID, 10)
	sReason := strconv.FormatInt(iReq.InReason, 10)

	requestValue := s.Replace(getTemplateforUpdateKeyword, "$attribute", iReq.InKeyword.Attribute, -1)
	requestValue = s.Replace(requestValue, "$count", sCount, -1)
	requestValue = s.Replace(requestValue, "$customerId", sCustomerID, -1)
	requestValue = s.Replace(requestValue, "$date", sDate, -1)
	requestValue = s.Replace(requestValue, "$keywordId", sKeywordID, -1)
	requestValue = s.Replace(requestValue, "$keywordTypeId", sKeywordTypeID, -1)
	requestValue = s.Replace(requestValue, "$name", iReq.InKeyword.Name, -1)
	requestValue = s.Replace(requestValue, "$extended", iReq.InKeyword.Extended, -1)
	requestValue = s.Replace(requestValue, "$id", sID, -1)
	requestValue = s.Replace(requestValue, "$inReason", sReason, -1)
	requestValue = s.Replace(requestValue, "$byUser", iReq.ByUser.ByUser, -1)
	requestValue = s.Replace(requestValue, "$byChannel", iReq.ByUser.ByChannel, -1)
	requestValue = s.Replace(requestValue, "$byProject", iReq.ByUser.ByProject, -1)
	requestValue = s.Replace(requestValue, "$byHost", iReq.ByUser.ByHost, -1)

	//log.Println("requestValue: " + requestValue)
	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 200
		oRes.ErrorDesc = resp
		return oRes
	}

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/UpdateKeyword"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	response, err := client.Do(req)
	if err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 200
		oRes.ErrorDesc = resp
		return oRes
	}
	defer response.Body.Close()

	//log.Println(response.Body)

	if response.StatusCode != 200 {
		resp = "error " + response.Status
		l.Response = resp
		oRes.ErrorCode = response.StatusCode
		oRes.ErrorDesc = response.Status
		return oRes
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		resp = err.Error()
		l.Response = resp
		oRes.ErrorCode = 400
		oRes.ErrorDesc = resp
		return oRes
	}

	//log.Println("contents : " + string(contents[:]))

	myResult := MyRespEnvelopeUpdateKeyword{}
	xml.Unmarshal([]byte(contents), &myResult)
	//log.Println(myResult)
	oRes.ResultValue = myResult.Body.VUpdateKeywordResponse.VUpdateKeywordResult.ResultValue
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.VUpdateKeywordResponse.VUpdateKeywordResult.ErrorCode)
	oRes.ErrorDesc = myResult.Body.VUpdateKeywordResponse.VUpdateKeywordResult.ErrorDesc

	//log.Println(oRes)

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	//l.TrackingNo = trackingno
	//l.ApplicationName = "TVSKeyword"
	//l.FunctionName = "UpdateKeyword"
	//l.Request = "ByUser=" + iReq.ByUser.ByUser
	jSRes, _ := json.Marshal(oRes)
	sJSRes := string(jSRes)

	l.Response = sJSRes
	//l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	//l.InsertappLog("./log/tvskeywordapplog.log", "UpdateKeyword")
	return oRes
}
