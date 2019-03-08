package main

import (
	"database/sql"
	"database/sql/driver"
	"io"

	//"encoding/json"
	//"log"
	cm "tvsglobal/common"
	st "tvsglobal/tvsstructs"

	//"github.com/streadway/amqp"
	_ "gopkg.in/goracle.v2"
)

func generatetasklist(Trackingno string, TVSOrdReq st.TVSSubmitOrderToQueue) (string, st.TVSSubmitOrderToQueue) {

	var resultI driver.Rows
	var err error
	var tvstask st.TVSTaskinfo
	cm.ExcutestoreDS("ICC", `begin tvs_servorder.generatetasklist(:p_trackingno,:p_ordertype,:p_rs );  end;`,
		Trackingno, TVSOrdReq.TVSOrdReq.OrderType, sql.Out{Dest: &resultI})
	defer resultI.Close()
	values := make([]driver.Value, len(resultI.Columns()))
	colmap := cm.Createmapcol(resultI.Columns())
	for {
		print(colmap)
		err = resultI.Next(values)
		if err == nil {
			if err == io.EOF {
				break
			}
		} else {
			break
		}
		tvstask.Taskid = values[colmap["TASKID"]].(int64)
		tvstask.Taskname = values[colmap["TASKNAME"]].(string)
		tvstask.MSname = values[colmap["MSNAME"]].(string)

		TVSOrdReq.TVSTaskList = append(TVSOrdReq.TVSTaskList, tvstask)
	}
	return "", TVSOrdReq
}
