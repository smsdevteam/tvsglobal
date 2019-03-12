package main

//"github.com/jmoiron/sqlx"

//Getccbssubinfo is function excute sql command
import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
	cm "tvsglobal/common"
	st "tvsglobal/tvsstructs"

	"github.com/streadway/amqp"
	_ "gopkg.in/goracle.v2"
)

func savereq(TVSOrdReq st.TVSSubmitOrdReqData) (string, st.TVSSubmitOrdResData) {
	var queuename string
	var TVSOrdRes st.TVSSubmitOrdResData
	cm.ExcutestoreDS("ICC", `begin
	-- Call the procedure
	tvs_servorder.save_requestservorder(:p_orderid,:p_ordertype,:p_channelcode,
										:p_orderdate,:p_tvscustno,:p_trackingno,
										:p_queuename,:p_errorcode,:p_errordesc);  end;`, TVSOrdReq.Orderid, TVSOrdReq.OrderType, TVSOrdReq.ChannelCode,
		TVSOrdReq.OrderDate, TVSOrdReq.TVSCustNo, sql.Out{Dest: &TVSOrdRes.Trackingno},
		sql.Out{Dest: &queuename}, sql.Out{Dest: &TVSOrdRes.ResponseResultobj.ErrorCode}, sql.Out{Dest: &TVSOrdRes.ResponseResultobj.ErrorDesc})
	return queuename, TVSOrdRes
}
func submitorder(TVSSubmitOrderRequest st.TVSSubmitOrdReqData) st.TVSSubmitOrdResData {
	// save to request log and put to queue
	var TVSOrdRes st.TVSSubmitOrdResData
	var queuename string
	queuename, TVSOrdRes = savereq(TVSSubmitOrderRequest)
	sendtoqueue(queuename, TVSSubmitOrderRequest, &TVSOrdRes)
	return TVSOrdRes
}
func sendtoqueue(queuename string, TVSOrdReq st.TVSSubmitOrdReqData, TVSOrdRes *st.TVSSubmitOrdResData) {
	var TVSOrdReqtoQueue st.TVSSubmitOrderData
	conn, err := amqp.Dial("amqp://admin:admin@172.19.218.104:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		queuename, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	TVSOrdReqtoQueue.TVSOrdReq = TVSOrdReq
	TVSOrdReqtoQueue.Trackingno = TVSOrdRes.Trackingno
	body, err := json.Marshal(TVSOrdReqtoQueue)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	var req st.TVSSubmitOrdReqData
	req.ChannelCode = "ood"
	req.OrderDate = time.Now()
	req.OrderType = "1"
	req.Orderid = "TEST001"
	req.TVSCustNo = 0
	submitorder(req)
}
