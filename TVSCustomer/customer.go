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

// GetCustomerByCustomerID get info
 
func CustomeGetDeviceInfo(iCustomerID string) c.Customerrespon {
	// Log#Start
	/*var l cm.Applog
	var trackingno string
	var resp string

	 t0 := time.Now()
	trackingno = t0.Format("20060102-150405")
	l.TrackingNo = trackingno
	l.ApplicationName = "TVScustomer"
	l.FunctionName = "Getcustomer"
	l.Request = "customerID=" + iCustomerID
	l.Start = t0.String()
	l.InsertappLog("./log/tvscustomerlog.log", "GetCustomer")
	*/
	//resp := "SUCCESS"
	var ocustomerInfo c.CustomerInfo
	var oDeviceinfo   c.DeviceInfo
	var oCustomerRespon c.Customerrespon
	var  oCustomerinfocolection   []c.CustomerInfo
	var  oDeviceinfocolection []c.DeviceInfo
	//var dbsource string 
	 
	dbsource :=  cm.GetDatasourceName("ICC") 
	  
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
		//resp = err.Error() //
	} else {
		defer db.Close()
		var statement string
		statement = "begin TVS_Go_Product.GetDeviceByCustomerID(:0,:1); end;"
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
		   colmap :=cm.Createmapcol(resultC.Columns())
		  
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
				if values[cm.Getcolindex(colmap, "CUSTOMER_ID")].(int64) != nil {
					ocustomerInfo.ID = values[cm.Getcolindex(colmap, "CUSTOMER_ID")].(int64)
				}
                oDeviceinfo.DeviceID     =  values[cm.Getcolindex(colmap,"DEVICE_ID")].(int64)
	            oDeviceinfo.SerialNumber     =   values[cm.Getcolindex(colmap, upper("SerialNumber"))].(string)
	            oDeviceinfo.StatusID           = values[cm.Getcolindex(colmap, upper("StatusID"))].(int64)
	            oDeviceinfo.StatusDesc          = values[cm.Getcolindex(colmap, upper("StatusDesc"))].(string)
	            oDeviceinfo.ModelID              =values[cm.Getcolindex(colmap, upper("ModelID"))].(int64)
	            oDeviceinfo.ModelDesc           =values[cm.Getcolindex(colmap, upper("ModelDesc")].(string)
				oDeviceinfo.ProductID          =values[cm.Getcolindex(colmap, upper("ProductID"))].(int64)
				oDeviceinfo.ProductDesc         =values[cm.Getcolindex(colmap, upper("ProductDesc"))].(string)
				oDeviceinfo.StockhandlerID       =values[cm.Getcolindex(colmap, upper("StockhandlerID"))].(int64)
				oDeviceinfo.StockhandlerDesc    =values[cm.Getcolindex(colmap, upper("StockhandlerDesc"))].(string)
				oDeviceinfo.AllowSystem          =values[cm.Getcolindex(colmap, upper("AllowSystem"))].(string)
				oDeviceinfo.FactoryWarrantyDate  =values[cm.Getcolindex(colmap, upper("FactoryWarrantyDate")].(string)
				oDeviceinfo.AgentKey            =values[cm.Getcolindex(colmap, upper("AgentKey"))].(string)
				oDeviceinfo.AgentName           =values[cm.Getcolindex(colmap, upper("AgentName"))].(string)
				oDeviceinfo.ReturnDate   =values[cm.Getcolindex(colmap, upper("ReturnDate"))].(string)
			
				oDeviceinfocolection =append(oDeviceinfocolection,oDeviceinfo)
				  
			   //	print(ocustomerInfo.BusinessUnitID)
			}
				  
			ocustomerInfo.DeviceList =append(ocustomerInfo.DeviceList,oDeviceinfocolection)
			oCustomerRespon.Customerrespon =ocustomerInfo
			//ocustomerInfo = oCustomer
         //log.Println(oCustomerinfocolection)
		}

	}
	    
        oCustomerRespon.CustomerInfocollection =oCustomerinfocolection
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
	return oCustomerRespon
}
 

