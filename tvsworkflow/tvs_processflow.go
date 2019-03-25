package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"net/http"

	cm "github.com/smsdevteam/tvsglobal/common"
	st "github.com/smsdevteam/tvsglobal/tvsstructs"

	_ "gopkg.in/goracle.v2"
)

func generatetasklist(Trackingno string, TVSOrdprocess st.TVSSubmitOrderProcess) st.TVSSubmitOrderProcess {

	var resultI driver.Rows
	var err error
	var tvstask st.TVSTaskinfo
	var dataprocess st.TVSSubmitOrderProcess
	cm.ExcutestoreDS("ICC", `begin tvs_servorder.generatetasklist(:p_trackingno,:p_ordertype,:p_rs );  end;`,
		Trackingno, TVSOrdprocess.Orderdata.TVSOrdReq.OrderType, sql.Out{Dest: &resultI})
	defer resultI.Close()
	values := make([]driver.Value, len(resultI.Columns()))
	colmap := cm.Createmapcol(resultI.Columns())
	for {
		//print(colmap)
		err = resultI.Next(values)
		if err == nil {
			if err == io.EOF {
				break
			}
		} else {
			break
		}
		tvstask.Taskid = values[colmap["TASKID"]].(string)
		tvstask.Seqno = values[colmap["SEQNO"]].(int64)
		tvstask.Taskname = values[colmap["TASKNAME"]].(string)
		tvstask.MSname = values[colmap["MSNAME"]].(string)
		tvstask.Servurl = values[colmap["SERVURL"]].(string)
		tvstask.Responseobjname = values[colmap["RESPONSEOBJNAME"]].(string)
		dataprocess.TVSTaskList = append(dataprocess.TVSTaskList, tvstask)
	}
	TVSOrdprocess.TVSTaskList = dataprocess.TVSTaskList
	return TVSOrdprocess
}
func savelogtask(Trackingno string, seqno int64, response st.ResponseResult) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Error func savelogtask .. %s\n", err)

		}
	}()

	cm.ExcutestoreDS("ICC", `begin tvs_servorder.savelogtask(:p_trackingno,:p_seqno,:p_errorcode,:p_errordesc );  end;`,
		Trackingno, seqno, response.ErrorCode, response.ErrorDesc)

	return "success"
}

func callsendcommand(tvssubmitdata st.TVSSubmitOrderData, taskobj st.TVSTaskinfo) {
	var msresponce st.TVSBN_Responseresult
	url := taskobj.Servurl //"http://restapi3.apiary.io/notes"

	b, _ := json.Marshal(tvssubmitdata)
	s := string(b)
	var jsonStr = []byte(s)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	tempbody := string(body)
	fmt.Println("response Body:", tempbody)
	tempbody = strings.Replace(tempbody, taskobj.Responseobjname, "TVSBN_RESPONSERESULT", -1)

	mySlice := []byte(tempbody)
	err = json.Unmarshal(mySlice, &msresponce)

	fmt.Println("response json:", msresponce)
	fmt.Println("*********************************************************")
}
