package main

import (
	"database/sql"
	"database/sql/driver"

	//"encoding/json"
	//"log"
	cm "tvsglobal/common"
	st "tvsglobal/tvsstructs"

	//"github.com/streadway/amqp"
	_ "gopkg.in/goracle.v2"
)

func generatetasklist(Trackingno string, TVSOrdReq  st.TVSSubmitOrderToQueue) (string, st.TVSQueueSubmitOrderReponse) {

	var TVSOrdRes st.TVSQueueSubmitOrderReponse
	var resultI driver.Rows
	cm.ExcutestoreDS("ICC", `begin tvs_servorder.generatetasklist(:p_trackingno,:p_ordertype,:p_rs );  end;`,
		Trackingno, TVSOrdReq.TVSOrdReq.OrderType, sql.Out{Dest: &resultI})
	return "", TVSOrdRes
}
