package main

import (
	"encoding/json"
	"log"

	st "tvsglobal/tvsstructs"

	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/streadway/amqp"
	_ "gopkg.in/goracle.v2"
)

func initailreceiver() st.Tvsqueueinfo {

	var Queueinfo st.Tvsqueueinfo
	config.Load(file.NewSource(
		file.WithPath("queueconfig.json"),
	))
	Queueinfo.Queuename = config.Get("RBQUEUE", "Queueinfo", "Queuename").String("") //"queue1"
	Queueinfo.Address = config.Get("RBQUEUE", "Queueinfo", "Address").String("")
	return Queueinfo
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func main() {
	var Queueinfo st.Tvsqueueinfo
	Queueinfo = initailreceiver()
	conn, err := amqp.Dial(Queueinfo.Address)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		Queueinfo.Queuename, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	var tvssubmitdata st.TVSSubmitOrderData
	go func() {
		for d := range msgs {
			err := json.Unmarshal(d.Body, &tvssubmitdata)
			if err != nil {
				print(err.Error())
			}
			initialtask(tvssubmitdata)
			//log.Printf("Received a message: %s",d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func initialtask(tvssubmitdata st.TVSSubmitOrderData) {
	var resultcode string
	var Processdata st.TVSSubmitOrderProcess
	Processdata.Orderdata = tvssubmitdata
	print("Get Task Config For Order Type " + tvssubmitdata.TVSOrdReq.OrderType + " Tracking no " + tvssubmitdata.Trackingno)
	resultcode, Processdata = generatetasklist(tvssubmitdata.Trackingno, tvssubmitdata)
	if resultcode == "success" {
		for i := 0; i < len(Processdata.TVSTaskList); i++ {
			taskid := Processdata.TVSTaskList[i].Taskid
			msname := Processdata.TVSTaskList[i].MSname
			switch taskid {
			case "1":
				log.Printf(" Start procee number " + msname)

			}
		}

	}
}
