package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
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

// MyRespEnvelopeUpdateNote for UpdateNote
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
func GetNoteByNoteID(iNoteID string) *st.GetNoteResult {
	// Log#Start
	l := cm.NewApplog()
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetNote"
	l.Request = "NoteID=" + iNoteID
	l.Start = t0.Format(time.RFC3339Nano)
	var tags []string
	tags = append(tags, "env7")
	tags = append(tags, "TVSNote")
	tags = append(tags, "applogs")
	l.Tags = tags
	l.InsertappLog("./log/tvsnoteapplog.log", "GetNote")
	l.PrintJSONLog()

	oRes := st.NewGetNoteResult()
	var oNote st.Note

	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		//log.Println(err)
		resp = err.Error()
		ts := time.Now()
		l.Timestamp = ts.Format(time.RFC3339Nano)
		l.Response = resp
		l.PrintJSONLog()
		oRes.ErrorCode = 2
		oRes.ErrorDesc = err.Error()
		return oRes
	} else {
		defer db.Close()
		var statement string
		statement = "begin PK_ICC_NOTE.GetNote(:0,:1); end;"
		var resultC driver.Rows
		intNoteID, err := strconv.Atoi(iNoteID)
		if err != nil {
			//log.Println(err)
			resp = err.Error()
			ts := time.Now()
			l.Timestamp = ts.Format(time.RFC3339Nano)
			l.Response = resp
			l.PrintJSONLog()
			oRes.ErrorCode = 3
			oRes.ErrorDesc = err.Error()
			return oRes
		} else {
			if _, err := db.Exec(statement, intNoteID, sql.Out{Dest: &resultC}); err != nil {
				//log.Println(err)
				resp = err.Error()
				ts := time.Now()
				l.Timestamp = ts.Format(time.RFC3339Nano)
				l.Response = resp
				l.PrintJSONLog()
				oRes.ErrorCode = 4
				oRes.ErrorDesc = err.Error()
				return oRes
			}

			defer resultC.Close()
			values := make([]driver.Value, len(resultC.Columns()))
			for {

				colmap := cm.Createmapcol(resultC.Columns())
				//log.Println(colmap)

				err = resultC.Next(values)
				if err != nil {
					if err == io.EOF {
						break
					}
					//log.Println("error:", err)
					resp = err.Error()
					ts := time.Now()
					l.Timestamp = ts.Format(time.RFC3339Nano)
					l.Response = resp
					l.PrintJSONLog()
					oRes.ErrorCode = 5
					oRes.ErrorDesc = err.Error()
					return oRes
				}

				if values[cm.Getcolindex(colmap, "CUSTOMER_ID")] != nil {
					oNote.CustomerID = values[cm.Getcolindex(colmap, "CUSTOMER_ID")].(int64)
				}
				if values[cm.Getcolindex(colmap, "CREATED_BY_USER_ID")] != nil {
					oNote.CreatedByUserID = values[cm.Getcolindex(colmap, "CREATED_BY_USER_ID")].(int64)
				}
				if values[cm.Getcolindex(colmap, "ACTION_USER_ID")] != nil {
					oNote.ActionUserID = values[cm.Getcolindex(colmap, "ACTION_USER_ID")].(int64)
				}

				oNote.CategoryID = values[cm.Getcolindex(colmap, "CATEGORY_ID")].(string)
				oNote.CompletionStageID = values[cm.Getcolindex(colmap, "COMPLETION_STAGE_ID")].(string)
				oNote.Body = values[cm.Getcolindex(colmap, "BODY")].(string)

				if values[cm.Getcolindex(colmap, "ID")] != nil {
					oNote.NoteID = values[cm.Getcolindex(colmap, "ID")].(int64)
				}
				if values[cm.Getcolindex(colmap, "CREATE_DATETIME")] != nil {
					oNote.CreateDateTime = values[cm.Getcolindex(colmap, "CREATE_DATETIME")].(time.Time)
				}
			}

		}

	}

	oRes.MyNote = oNote
	if oRes.ErrorCode == 1 {
		oRes.ErrorCode = 0
		oRes.ErrorDesc = "Success"
	}

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetNote"
	l.Request = "NoteID=" + iNoteID
	l.Response = resp
	l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	ts := time.Now()
	l.Timestamp = ts.Format(time.RFC3339Nano)
	l.InsertappLog("./log/tvsnoteapplog.log", "GetNote")
	l.PrintJSONLog()
	//test
	return oRes
}

//GetListNoteByCustomerID get list note by customer id
func GetListNoteByCustomerID(iCustomerID string) *st.GetListNoteResult {
	// Log#Start
	l := cm.NewApplog()
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetListNoteByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Start = t0.Format(time.RFC3339Nano)
	var tags []string
	tags = append(tags, "env7")
	tags = append(tags, "TVSNote")
	tags = append(tags, "applogs")
	l.Tags = tags
	l.InsertappLog("./log/tvsnoteapplog.log", "GetListNoteByCustomerID")
	l.PrintJSONLog()

	oRes := st.NewGetListNoteResult()
	var oListNote st.ListNote
	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		//log.Println("error:", err)
		resp = err.Error()
		ts := time.Now()
		l.Timestamp = ts.Format(time.RFC3339Nano)
		l.Response = resp
		l.PrintJSONLog()
		oRes.ErrorCode = 2
		oRes.ErrorDesc = err.Error()
		return oRes
	} else {
		defer db.Close()
		var statement string
		statement = "begin PK_ICC_NOTE.GetNoteByCustomerID(:0,:1); end;"
		var resultC driver.Rows
		intCustomerID, err := strconv.Atoi(iCustomerID)
		if err != nil {
			//log.Println("error:", err)
			resp = err.Error()
			ts := time.Now()
			l.Timestamp = ts.Format(time.RFC3339Nano)
			l.Response = resp
			l.PrintJSONLog()
			oRes.ErrorCode = 3
			oRes.ErrorDesc = err.Error()
			return oRes
		} else {
			if _, err := db.Exec(statement, intCustomerID, sql.Out{Dest: &resultC}); err != nil {
				//log.Println("error:", err)
				resp = err.Error()
				ts := time.Now()
				l.Timestamp = ts.Format(time.RFC3339Nano)
				l.Response = resp
				l.PrintJSONLog()
				oRes.ErrorCode = 4
				oRes.ErrorDesc = err.Error()
				return oRes
			}
			defer resultC.Close()
			values := make([]driver.Value, len(resultC.Columns()))
			var oLNote []st.Note
			for {

				colmap := CreateMapCol(resultC.Columns())

				err = resultC.Next(values)
				if err != nil {
					if err == io.EOF {
						break
					}
					//log.Println("error:", err)
					resp = err.Error()
					ts := time.Now()
					l.Timestamp = ts.Format(time.RFC3339Nano)
					l.Response = resp
					l.PrintJSONLog()
					oRes.ErrorCode = 5
					oRes.ErrorDesc = err.Error()
					return oRes
				}
				var oNote st.Note

				if values[colmap["CUSTOMER_ID"]] != nil {
					oNote.CustomerID = values[colmap["CUSTOMER_ID"]].(int64)
				}
				if values[colmap["CREATED_BY_USER_ID"]] != nil {
					oNote.CreatedByUserID = values[colmap["CREATED_BY_USER_ID"]].(int64)
				}
				if values[colmap["ACTION_USER_ID"]] != nil {
					oNote.ActionUserID = values[colmap["ACTION_USER_ID"]].(int64)
				}

				oNote.CategoryID = values[colmap["CATEGORY_ID"]].(string)
				oNote.CompletionStageID = values[colmap["COMPLETION_STAGE_ID"]].(string)
				oNote.Body = values[colmap["BODY"]].(string)

				if values[6] != nil {
					oNote.NoteID = values[colmap["ID"]].(int64)
				}
				if values[7] != nil {
					oNote.CreateDateTime = values[colmap["CREATE_DATETIME"]].(time.Time)
				}

				oLNote = append(oLNote, oNote)
			}
			oListNote.Notes = oLNote
		}
	}
	oRes.MyListNote = oListNote
	if oRes.ErrorCode == 1 {
		oRes.ErrorCode = 0
		oRes.ErrorDesc = "Success"
	}

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "GetListNoteByCustomerID"
	l.Request = "CustomerID=" + iCustomerID
	l.Response = resp
	l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	ts := time.Now()
	l.Timestamp = ts.Format(time.RFC3339Nano)
	l.PrintJSONLog()
	l.InsertappLog("./log/tvsnoteapplog.log", "GetListNoteByCustomerID")
	//test
	return oRes
}

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
func CreateNote(iReq st.CreateNoteRequest) *st.CreateNoteResponse {

	// Log#Start
	l := cm.NewApplog()
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "CreateNote"
	l.Request = "ByUser=" + iReq.ByUser.ByUser + " ByChannel=" + iReq.ByUser.ByChannel
	l.Start = t0.Format(time.RFC3339Nano)
	var tags []string
	tags = append(tags, "env7")
	tags = append(tags, "TVSNote")
	tags = append(tags, "applogs")
	l.Tags = tags
	l.PrintJSONLog()
	l.InsertappLog("./log/tvsnoteapplog.log", "CreateNote")

	oRes := st.NewCreateNoteResponse()
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
		resp = err.Error()
		ts := time.Now()
		l.Timestamp = ts.Format(time.RFC3339Nano)
		l.Response = resp
		l.PrintJSONLog()
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/CreateNote"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	response, err := client.Do(req)
	if err != nil {
		resp = err.Error()
		ts := time.Now()
		l.Timestamp = ts.Format(time.RFC3339Nano)
		l.Response = resp
		l.PrintJSONLog()
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer response.Body.Close()

	//log.Println(response.Body)

	if response.StatusCode != 200 {
		resp = err.Error()
		ts := time.Now()
		l.Timestamp = ts.Format(time.RFC3339Nano)
		l.Response = resp
		l.PrintJSONLog()
		oRes.ErrorCode = response.StatusCode
		oRes.ErrorDesc = response.Status
		return oRes
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		resp = err.Error()
		ts := time.Now()
		l.Timestamp = ts.Format(time.RFC3339Nano)
		l.Response = resp
		l.PrintJSONLog()
		oRes.ErrorCode = 400
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	//log.Println("contents : " + string(contents[:]))

	myResult := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &myResult)
	//log.Println(myResult)
	oRes.ResultValue = myResult.Body.GetResponse.MyCreateNoteResult.ResultValue
	oRes.ErrorCode, _ = strconv.Atoi(myResult.Body.GetResponse.MyCreateNoteResult.ErrorCode)
	oRes.ErrorDesc = myResult.Body.GetResponse.MyCreateNoteResult.ErrorDesc

	//log.Println(oRes)

	// Log#Stop
	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "CreateNote"
	l.Request = "ByUser=" + iReq.ByUser.ByUser
	l.Response = resp
	l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	ts := time.Now()
	l.Timestamp = ts.Format(time.RFC3339Nano)
	l.PrintJSONLog()
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

// UpdateNote to
func UpdateNote(iReq st.UpdateNoteRequest) st.UpdateNoteResponse {
	// Log#Start
	l := cm.NewApplog()
	var trackingno string
	var resp string
	resp = "SUCCESS"
	t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSNote"
	l.FunctionName = "UpdateNote"
	l.Request = "ByUser=" + iReq.ByUser.ByUser + " ByChannel=" + iReq.ByUser.ByChannel
	l.Start = t0.Format(time.RFC3339Nano)
	var tags []string
	tags = append(tags, "env7")
	tags = append(tags, "TVSNote")
	tags = append(tags, "applogs")
	l.Tags = tags
	l.PrintJSONLog()
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
		resp = err.Error()
		ts := time.Now()
		l.Timestamp = ts.Format(time.RFC3339Nano)
		l.Response = resp
		l.PrintJSONLog()
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}

	req.Header.Add("SOAPAction", `"http://tempuri.org/IICCServiceInterface/UpdateNote"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	response, err := client.Do(req)
	if err != nil {
		resp = err.Error()
		ts := time.Now()
		l.Timestamp = ts.Format(time.RFC3339Nano)
		l.Response = resp
		l.PrintJSONLog()
		oRes.ErrorCode = 200
		oRes.ErrorDesc = err.Error()
		return oRes
	}
	defer response.Body.Close()

	//log.Println(response.Body)

	if response.StatusCode != 200 {
		resp = err.Error()
		ts := time.Now()
		l.Timestamp = ts.Format(time.RFC3339Nano)
		l.Response = resp
		l.PrintJSONLog()
		oRes.ErrorCode = response.StatusCode
		oRes.ErrorDesc = response.Status
		return oRes
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		resp = err.Error()
		ts := time.Now()
		l.Timestamp = ts.Format(time.RFC3339Nano)
		l.Response = resp
		l.PrintJSONLog()
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
	l.Start = t0.Format(time.RFC3339Nano)
	l.End = t1.Format(time.RFC3339Nano)
	l.Duration = t2.String()
	ts := time.Now()
	l.Timestamp = ts.Format(time.RFC3339Nano)
	l.PrintJSONLog()
	l.InsertappLog("./log/tvsnoteapplog.log", "UpdateNote")

	return oRes
}

//CreateMapCol for use column name to point
func CreateMapCol(data []string) map[string]int {
	var colmap = map[string]int{}

	for k, v := range data {
		colmap[v] = k
	}
	return colmap
}
