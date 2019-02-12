package main

import (
 	"database/sql"
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	"io"
 	"log"
 	"strconv"
 	_ "gopkg.in/goracle.v2"

	cm "github.com/smsdevteam/tvsglobal/common" //db
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
	XMLName xml.Name `xml:"CreateCustomerResponse"`
	//	MyCreateCustomerResult CreateCustomerResult `xml:"CreateCustomerResult"`
	//	MyResult authenHD `xml:"AuthenticateByProofResult "`
}
type bodyUpdateCustomer struct {
	XMLName xml.Name `xml:"Body"`
	//Fault       *fault
	GetResponse completeResponseUpdateCustomer `xml:"UpdateCustomerResponse"`
}

type completeResponseUpdateCustomer struct {
	XMLName                xml.Name             `xml:"UpdateCustomerResponse"`
	MyUpdateCustomerResult updateCustomerResult `xml:"UpdateCustomerResult"`
	//	MyResult authenHD `xml:"AuthenticateByProofResult "`
}

type updateCustomerResult struct {
	XMLName     xml.Name `xml:"UpdateCustomerResult"`
	ResultValue string   `xml:"ResultValue"`
	ErrorCode   string   `xml:"ErrorCode"`
	ErrorDesc   string   `xml:"ErrorDesc"`
}

// GetWorkorderByCustomerID get info
func GetWorkorderByCustomerID(iCustomerID string) c.WorkorderInfo {
 
	//resp := "SUCCESS"
	var oWorkorderinfo c.WorkorderInfo
	//var dbsource string 
	
	dbsource :=  cm.GetDatasourceName("ICC") 
	 
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
		//resp = err.Error() //
	} else {
		defer db.Close()
		var statement string
		statement = "begin tvs_manualupdate.getworkorderbycustomerid(:0,:1); end;"
		var resultC driver.Rows
		intCustomerID, err := strconv.Atoi(iCustomerID)
		if err != nil {
			log.Fatal(err)
			//resp = err.Error()
		} else {
			if _, err := db.Exec(statement, intCustomerID, sql.Out{Dest: &resultC}); err != nil {
				log.Fatal(err)
				//resp = err.Error()
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
					//resp = err.Error()
				}
				//var oCustomer c.CustomerInfo
				if values[0] != nil {
					oWorkorderinfo.ID = values[0].(string)
				} 
			}
		}
	}
	return oWorkorderinfo
}
