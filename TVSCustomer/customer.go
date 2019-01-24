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
	c "github.com/smsdevteam/tvsglobal/tvsstructs" // referpath
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
	GetResponse completeResponse `xml:"CreateCustomerResponse"`
}

// type fault struct {
// 	Code   string `xml:"faultcode"`
// 	String string `xml:"faultstring"`
// 	Detail string `xml:"detail"`
// }

type completeResponse struct {
	XMLName            xml.Name         `xml:"CreateCustomerResponse"`
	MyCreateCustomerResult CreateCustomerResult `xml:"CreateCustomerResult"`
	//	MyResult authenHD `xml:"AuthenticateByProofResult "`
}

type CreateCustomerResult struct {
	XMLName     xml.Name `xml:"CreateCustomerResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
}

// MyRespEnvelope for UpdateNote
type MyRespEnvelopeUpdateCustomer struct {
	XMLName xml.Name       `xml:"Envelope"`
	Body    bodyUpdateNote `xml:"Body"`
}

type bodyUpdateCustomer struct {
	XMLName xml.Name `xml:"Body"`
	//Fault       *fault
	GetResponse completeResponseUpdateCustomer `xml:"UpdateCustomerResponse"`
}

type completeResponseUpdateCustomer struct {
	XMLName            xml.Name         `xml:"UpdateCustomerResponse"`
	MyUpdateCustomerResult updateCustomerResult `xml:"UpdateCustomerResult"`
	//	MyResult authenHD `xml:"AuthenticateByProofResult "`
}

type updateCustomerResult struct {
	XMLName     xml.Name `xml:"UpdateCustomerResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
}

// GetCustomerByCustomerID get info
func GetCustomerByCustomerID(iCustomerID string) c.customer {
	// Log#Start
	/*var l cm.Applog
	var trackingno string
	var resp string
	resp = "SUCCESS"
	 t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVScustomer"
	l.FunctionName = "Getcustomer"
	l.Request = "customerID=" + iCustomerID
	l.Start = t0.String()
	l.InsertappLog("./log/tvscustomerlog.log", "GetCustomer")
 */
	var ocustomerInfo c.customerInfo

	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
		resp = err.Error()
	} else {
		defer db.Close()
		var statement string
		statement = "begin TVS_customer.getCustomerINFO(:0,:1); end;"
		var resultC driver.Rows
		intCustomerID, err := strconv.Atoi(icustomerID)
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
			for {
				err = resultC.Next(values)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Println("error:", err)
					resp = err.Error()
				}
		var oCustomer c.customer
				if values[0] != nil {
					oCustomer.CustomerID = values[0].(string)
				}
 

				oCustomer.BusinessUnitId = values[0].(string)
 
			}

		}

	}

	// Log#Stop
/* 	t1 := time.Now()
	t2 := t1.Sub(t0)
	l.TrackingNo = trackingno
	l.ApplicationName = "TVSCustomer"
	l.FunctionName = "GetCustomer"
	l.Request = "customerID=" + iCustomerID
	l.Response = resp
	l.Start = t0.String()
	l.End = t1.String()
	l.Duration = t2.String()
	l.InsertappLog("./log/tvscustomerapplog.log", "GetCustomer") */
	//test
	return oCustomer
}