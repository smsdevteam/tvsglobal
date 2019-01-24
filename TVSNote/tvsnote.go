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
	s "strings"
	"time"

	_ "gopkg.in/goracle.v2"

	cm "github.com/smsdevteam/tvsglobal/common"     // db
	st "github.com/smsdevteam/tvsglobal/tvsstructs" // referpath
)

var p = fmt.Println

// MyRespEnvelope for CreateNote
type MyRespEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    body     `xml:"Body"`
}

type body struct {
	XMLName xml.Name `xml:"Body"`
	//Fault       *fault
	GetResponse completeResponse `xml:"CreateNoteResponse"`
}

// type fault struct {
// 	Code   string `xml:"faultcode"`
// 	String string `xml:"faultstring"`
// 	Detail string `xml:"detail"`
// }

type completeResponse struct {
	XMLName            xml.Name         `xml:"CreateNoteResponse"`
	MyCreateNoteResult createNoteResult `xml:"CreateNoteResult"`
	//	MyResult authenHD `xml:"AuthenticateByProofResult "`
}

type createNoteResult struct {
	XMLName     xml.Name `xml:"CreateNoteResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
}

// MyRespEnvelope for UpdateNote
type MyRespEnvelopeUpdateNote struct {
	XMLName xml.Name       `xml:"Envelope"`
	Body    bodyUpdateNote `xml:"Body"`
}

type bodyUpdateNote struct {
	XMLName xml.Name `xml:"Body"`
	//Fault       *fault
	GetResponse completeResponseUpdateNote `xml:"UpdateNoteResponse"`
}

type completeResponseUpdateNote struct {
	XMLName            xml.Name         `xml:"UpdateNoteResponse"`
	MyUpdateNoteResult updateNoteResult `xml:"UpdateNoteResult"`
	//	MyResult authenHD `xml:"AuthenticateByProofResult "`
}

type updateNoteResult struct {
	XMLName     xml.Name `xml:"UpdateNoteResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
}

// GetNoteByNoteID get info
func GetNoteByNoteID(iNoteID string) st.Note {
	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetNote"
	l.Request = "NoteID=" + iNoteID
	l.Start = t0.String()
	l.InsertappLog("./log/tvsnoteapplog.log", "GetNote")

	var oNote st.Note

	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
		resp = err.Error()
	} else {
		defer db.Close()
		var statement string
		statement = "begin PK_ICC_NOTE.GetNote(:0,:1); end;"
		var resultC driver.Rows
		intNoteID, err := strconv.Atoi(iNoteID)
		if err != nil {
			log.Fatal(err)
			resp = err.Error()
		} else {
			if _, err := db.Exec(statement, intNoteID, sql.Out{Dest: &resultC}); err != nil {
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

				if values[0] != nil {
					oNote.CustomerID = values[0].(int64)
				}
				if values[1] != nil {
					oNote.CreatedByUserID = values[1].(int64)
				}
				if values[2] != nil {
					oNote.ActionUserID = values[2].(int64)
				}

				oNote.CategoryID = values[3].(string)
				oNote.CompletionStageID = values[4].(string)
				oNote.Body = values[5].(string)

				if values[6] != nil {
					oNote.NoteID = values[6].(int64)
				}
				if values[7] != nil {
					oNote.CreateDateTime = values[7].(time.Time)
				}
			}

		}

	}

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetNote"
	l.Request = "NoteID=" + iNoteID
	l.Response = resp
	l.Start = t0.String()
	l.End = t1.String()
	l.Duration = t2.String()
	l.InsertappLog("./log/tvsnoteapplog.log", "GetNote")
	//test
	return oNote
}

// get list note by customer id
func GetListNoteByCustomerID(iCustomerID string) st.ListNote {
	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetListNoteByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Start = t0.String()
	l.InsertappLog("./log/tvsnoteapplog.log", "GetListNoteByCustomerID")

	var oListNote st.ListNote
	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
		resp = err.Error()
	} else {
		defer db.Close()
		var statement string
		statement = "begin PK_ICC_NOTE.GetNoteByCustomerID(:0,:1); end;"
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
			var oLNote []st.Note
			for {
				err = resultC.Next(values)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Println("error:", err)
					resp = err.Error()
				}
				var oNote st.Note

				if values[0] != nil {
					oNote.CustomerID = values[0].(int64)
				}
				if values[1] != nil {
					oNote.CreatedByUserID = values[1].(int64)
				}
				if values[2] != nil {
					oNote.ActionUserID = values[2].(int64)
				}

				oNote.CategoryID = values[3].(string)
				oNote.CompletionStageID = values[4].(string)
				oNote.Body = values[5].(string)

				if values[6] != nil {
					oNote.NoteID = values[6].(int64)
				}
				if values[7] != nil {
					oNote.CreateDateTime = values[7].(time.Time)
				}

				oLNote = append(oLNote, oNote)
			}
			oListNote.Notes = oLNote
		}
	}
	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetListNoteByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Response = resp
	l.Start = t0.String()
	l.End = t1.String()
	l.Duration = t2.String()
	l.InsertappLog("./log/tvsnoteapplog.log", "GetListNoteByCustomerID")
	//test
	return oListNote
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

const getTemplateforCreateNote = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<CreateNote xmlns="http://tempuri.org/">
		<inNote>
			<ActionUserKey>$actionUserKey</ActionUserKey>
			<Body>$body</Body>
			<CategoryKey>$categoryKey</CategoryKey>
			<CompletionStageKey>$completionStageKey</CompletionStageKey>
            <CreateDate>2019-01-21T08:01:43+07:00</CreateDate>
            <CreatedByUserId>$createdByUserId</CreatedByUserId>
            <CustomerId>$customerId</CustomerId>
			<Id>0</Id>
			<Extended>$extended</Extended>
		</inNote>
		<inReason>$inReason</inReason>
		<byUser>
			<byUser>$byUser</byUser>
            <byChannel>$byChannel</byChannel>
            <byProject>$byProject</byProject>
            <byHost>$byHost</byHost>
		</byUser>
	</CreateNote>
</s:Body>
</s:Envelope>`

//CreateNote for icc microservice
func CreateNote(iReq st.CreateNoteRequest) st.CreateNoteResponse {

	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "CreateNote"
	l.Request = "ByUser=" + iReq.ByUser.ByUser + " ByChannel=" + iReq.ByUser.ByChannel
	l.Start = t0.String()
	l.InsertappLog("./log/tvsnoteapplog.log", "CreateNote")

	var oRes st.CreateNoteResponse
	var AppServiceLnk cm.AppServiceURL
	AppServiceLnk = cm.AppReadConfig("ENH")

	url := AppServiceLnk.ICCServiceURL
	client := &http.Client{}

	sActionUsrKey := strconv.FormatInt(iReq.InNote.ActionUserKey, 10)
	//sCreateDate := ""
	//sCreateDate := (iReq.InNote.CreateDate).Format("2006-01-02T15:04:05")
	sCreatedByUserID := strconv.FormatInt(iReq.InNote.CreatedByUserID, 10)
	sCustomerID := strconv.FormatInt(iReq.InNote.CustomerID, 10)
	sInReason := strconv.FormatInt(iReq.InReason, 10)

	//log.Println("sCreateDate :" + sCreateDate)

	requestValue := s.Replace(getTemplateforCreateNote, "$actionUserKey", sActionUsrKey, -1)
	requestValue = s.Replace(requestValue, "$body", iReq.InNote.Body, -1)
	requestValue = s.Replace(requestValue, "$categoryKey", iReq.InNote.CategoryKey, -1)
	requestValue = s.Replace(requestValue, "$completionStageKey", iReq.InNote.CompletionStageKey, -1)
	//requestValue = s.Replace(requestValue, "$createDate", sCreateDate, -1)
	requestValue = s.Replace(requestValue, "$createdByUserId", sCreatedByUserID, -1)
	requestValue = s.Replace(requestValue, "$customerId", sCustomerID, -1)
	requestValue = s.Replace(requestValue, "$extended", iReq.InNote.Extended, -1)
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

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/CreateNote"`)
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

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	log.Println(myResult)
	oRes.ResultValue = myResult.Body.GetResponse.MyCreateNoteResult.ResultValue
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.GetResponse.MyCreateNoteResult.ErrorCode)

	//log.Println(oRes)

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "CreateNote"
	l.Request = "ByUser=" + iReq.ByUser.ByUser
	l.Response = resp
	l.Start = t0.String()
	l.End = t1.String()
	l.Duration = t2.String()
	l.InsertappLog("./log/tvsnoteapplog.log", "CreateNote")
	return oRes
}

const getTemplateforUpdateNote = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
	<UpdateNote xmlns="http://tempuri.org/">
		<inNote>
			<ActionUserKey>$actionUserKey</ActionUserKey>
			<Body>$body</Body>
			<CategoryKey>$categoryKey</CategoryKey>
			<CompletionStageKey>$completionStageKey</CompletionStageKey>
			<CreateDate>2019-01-21T08:01:43+07:00</CreateDate>
			<CreatedByUserId>$createdByUserId</CreatedByUserId>
			<CustomerId>$customerId</CustomerId>
			<Id>$id</Id>
			<Extended>$extended</Extended>
		</inNote>
		<inReason>$inReason</inReason>
		<byUser>
			<byUser>$byUser</byUser>
			<byChannel>$byChannel</byChannel>
			<byProject>$byProject</byProject>
			<byHost>$byHost</byHost>
		</byUser>
	</UpdateNote>
</s:Body>
</s:Envelope>`

func UpdateNote(iReq st.UpdateNoteRequest) st.UpdateNoteResponse {
	// Log#Start
	var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "UpdateNote"
	l.Request = "ByUser=" + iReq.ByUser.ByUser + " ByChannel=" + iReq.ByUser.ByChannel
	l.Start = t0.String()
	l.InsertappLog("./log/tvsnoteapplog.log", "UpdateNote")

	var oRes st.UpdateNoteResponse

	var AppServiceLnk cm.AppServiceURL
	AppServiceLnk = cm.AppReadConfig("ENH")

	url := AppServiceLnk.ICCServiceURL
	client := &http.Client{}

	sActionUsrKey := strconv.FormatInt(iReq.InNote.ActionUserKey, 10)
	sCreatedByUserID := strconv.FormatInt(iReq.InNote.CreatedByUserID, 10)
	sCustomerID := strconv.FormatInt(iReq.InNote.CustomerID, 10)
	sInReason := strconv.FormatInt(iReq.InReason, 10)
	sNoteID := strconv.FormatInt(iReq.InNote.NoteID, 10)

	requestValue := s.Replace(getTemplateforUpdateNote, "$actionUserKey", sActionUsrKey, -1)
	requestValue = s.Replace(requestValue, "$body", iReq.InNote.Body, -1)
	requestValue = s.Replace(requestValue, "$categoryKey", iReq.InNote.CategoryKey, -1)
	requestValue = s.Replace(requestValue, "$completionStageKey", iReq.InNote.CompletionStageKey, -1)
	requestValue = s.Replace(requestValue, "$createdByUserId", sCreatedByUserID, -1)
	requestValue = s.Replace(requestValue, "$customerId", sCustomerID, -1)
	requestValue = s.Replace(requestValue, "$id", sNoteID, -1)
	requestValue = s.Replace(requestValue, "$extended", iReq.InNote.Extended, -1)
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

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/UpdateNote"`)
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

	myResultUpdate := MyRespEnvelopeUpdateNote{}
	xml.Unmarshal([]byte(contents), &myResultUpdate)
	//log.Println(myResultUpdate)
	oRes.ResultValue = myResultUpdate.Body.GetResponse.MyUpdateNoteResult.ResultValue
	oRes.ErrorCode, _ = strconv.Atoi(myResultUpdate.Body.GetResponse.MyUpdateNoteResult.ErrorCode)
	oRes.ErrorDesc = myResultUpdate.Body.GetResponse.MyUpdateNoteResult.ErrorDesc

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "UpdateNote"
	l.Request = "ByUser=" + iReq.ByUser.ByUser
	l.Response = resp
	l.Start = t0.String()
	l.End = t1.String()
	l.Duration = t2.String()
	l.InsertappLog("./log/tvsnoteapplog.log", "UpdateNote")

	return oRes
}
