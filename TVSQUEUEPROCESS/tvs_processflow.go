package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"log"
	st "tvsglobal/tvsstructs"

	//"github.com/streadway/amqp"
	_ "gopkg.in/goracle.v2"
)

func generatetasklist(Trackingno string, TVSOrdReq st.TVSQueuSubmitOrderRequest) (string, st.TVSQueueSubmitOrderReponse) {
	var queuename string
	var TVSOrdRes st.TVSQueueSubmitOrderReponse
	var resultI driver.Rows
	cm.ExcutestoreDS("ICC", `begin tvs_servorder.generatetasklist(:p_trackingno,:p_ordertype,:p_rs );  end;`,
	 Trackingno, TVSOrdReq.OrderType,sql.Out{Dest: &resultI})
	return queuename, TVSOrdRes
}
 