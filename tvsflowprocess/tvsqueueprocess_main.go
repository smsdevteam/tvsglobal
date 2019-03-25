package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"net/http"

	cm "github.com/smsdevteam/tvsglobal/common"
	st "github.com/smsdevteam/tvsglobal/tvsstructs"

	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/streadway/amqp"
	_ "gopkg.in/goracle.v2"
)

const applicationname string = "tvsqueueprocess"
const tagappname string = "icc-tvsqueueprocess"
const taglogtype string = "info"

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
				print("found error " + err.Error())
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Error func initialtask .. %s\n", err)
		}
	}()
	var applog cm.Applog
	defer applog.PrintJSONLog()
	applog = cm.NewApploginfo("", applicationname, "initialtask",
		tagappname, taglogtype)

	var resultcode string
	var Processdata st.TVSSubmitOrderProcess
	Processdata.Orderdata = tvssubmitdata
	//print("Get Task Config For Order Type " + tvssubmitdata.TVSOrdReq.OrderType + " Tracking no " + tvssubmitdata.Trackingno)
	Processdata = generatetasklist(tvssubmitdata.Trackingno, Processdata)
	resultcode = "success"
	if resultcode == "success" {
		for i := 0; i < len(Processdata.TVSTaskList); i++ {
			taskid := Processdata.TVSTaskList[i].Taskid
			//msname := Processdata.TVSTaskList[i].MSname
			switch taskid {
			case "1": // change package
				log.Printf(" Start process number " + taskid)
				callserv(Processdata.Orderdata, Processdata.TVSTaskList[i])
			case "2": // refresh signal
				log.Printf(" Start procee number " + taskid)
				 callservrefreshsignal(Processdata.Orderdata, Processdata.TVSTaskList[i])

			}
		}

	}
}
func callserv(tvssubmitdata st.TVSSubmitOrderData, taskobj st.TVSTaskinfo) {
	var msresponce st.TVSBN_Responseresult
	url := taskobj.Servurl //"http://restapi3.apiary.io/notes"
	fmt.Println("URL:>", url)
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
func callservrefreshsignal(tvssubmitdata st.TVSSubmitOrderData, taskobj st.TVSTaskinfo) {
	url := "http://localhost:8081/tvsdevice/sendcmdtodevice/deviceid=1/reason=1/by=1"
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "2d4ad3e9-71a4-4620-aac9-2b4771dc4d7b")
	res, _ := http.DefaultClient.Do(req)
	//defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res)
	fmt.Println(string(body))

	/*var msresponce st.TVSBN_Responseresult
	url := taskobj.Servurl //"http://restapi3.apiary.io/notes"
	fmt.Println("URL:>", url)
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
	*/
}
