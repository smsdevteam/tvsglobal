package main

import (
	"encoding/json"
	"log"
	st "tvsglobal/tvsstructs"

	"github.com/streadway/amqp"
	_ "gopkg.in/goracle.v2"
)

func initailreceiver() st.Tvsqueueinfo {

	var Queueinfo st.Tvsqueueinfo
	Queueinfo.Queuename = "queue1"
	Queueinfo.Address = "amqp://admin:admin@172.19.218.104:5672/"
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
	var TVSOrdReqtoQueue st.TVSSubmitOrderToQueue
	go func() {
		for d := range msgs {
			err := json.Unmarshal(d.Body, &TVSOrdReqtoQueue)
			if err != nil {
				print(err.Error())
			}
			initialtask(TVSOrdReqtoQueue)
			//log.Printf("Received a message: %s",d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func initialtask(TVSOrdReqtoQueue st.TVSSubmitOrderToQueue) {
	var resultcode string
	print("Get Task Config For Order Type " +TVSOrdReqtoQueue.TVSOrdReq.OrderType +" Tracking no " + TVSOrdReqtoQueue.Trackingno)
	resultcode ,TVSOrdReqtoQueue= generatetasklist( TVSOrdReqtoQueue.Trackingno,TVSOrdReqtoQueue)
	if resultcode =="success"{
		
		select {
        case msg1 := <-c1:
            fmt.Println("received", msg1)
        case msg2 := <-c2:
            fmt.Println("received", msg2)
        }
	}
}
