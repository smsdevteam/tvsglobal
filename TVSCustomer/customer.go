package main

import (
 	"database/sql"
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	"io"
	"runtime"
 	"log"
 	"strconv"
 	_ "gopkg.in/goracle.v2"
    en  "OS"
	cm "github.com/smsdevteam/tvsglobal/common" //db
	c "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
 
)
const applicationname string = "tvscustomer"
const tagappname string = "icc-tvscustomer"
const taglogtype string = "applogs"
const tagenv string = "set02"
var p = fmt.Println

//MyRespEnvelope for CreateNote
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

 
 //CustomeGetDeviceInfo 
func CustomeGetDeviceInfo(iCustomerID string) c.Customerrespon {
	// Log#Start
	/* var TVSCUSTRes c.TVSCustomerOrdResData
	 var queuename string
	var applog cm.Applog
	defer applog.PrintJSONLog()
	applog = cm.NewApploginfo("", applicationname, "CustomeGetDeviceInfo",
	tagenv,  tagappname, taglogtype)
	b, _ := json.Marshal(iCustomerID)
	// Convert bytes to string.
	s := string(b)
	applog.Request = s
	queuename, TVSCUSTRes = savereq(iCustomerID)
	 
	applog.TrackingNo = TVSCUSTRes.Trackingno
 */

	var ocustomerInfo c.CustomerInfo
 	var oDeviceinfo   c.DeviceInfo
	var oCustomerRespon c.Customerrespon
//	var  oCustomerinfocolection   []c.CustomerInfo
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
				if values[0]!= nil {
					ocustomerInfo.CUSTOMERId = values[cm.Getcolindex(colmap, "CUSTOMERID")].(int64)
				}
		  	 	defer func() {
              if err := recover(); err != nil {
                   fmt.Println("HERE")
                   fmt.Println(err)
                 _, fn, line, _ := runtime.Caller(1)
		         log.Printf("[error] %s:%d %v", fn, line, err) 
                  }
				}()  
				
				 oDeviceinfo.ID               =  values[cm.Getcolindex(colmap,"DEVICEID")].(int64)
	            oDeviceinfo.Serial_Number     =   values[cm.Getcolindex(colmap,  "SERIAL_NUMBER")].(string)
			 	oDeviceinfo.Status_ID         = values[cm.Getcolindex(colmap,  "STATUS_ID")].(int64)
				oDeviceinfo.StatusDesc        = values[cm.Getcolindex(colmap,  "STATUSDESC")].(string)
                 oDeviceinfo.Stock_HandlerID       =values[cm.Getcolindex(colmap,  "STOCK_HANDLERID")].(int64)
				oDeviceinfo.Stock_HandlerName    =values[cm.Getcolindex(colmap,  "STOCK_HANDLERNAME")].(string)
			    oDeviceinfo.Model_ID              =values[cm.Getcolindex(colmap,  "MODEL_ID")].(int64)
	            oDeviceinfo.Model_Desc           =values[cm.Getcolindex(colmap,  "MODEL_DESC")].(string)
				oDeviceinfo.Technical_Product_ID  =values[cm.Getcolindex(colmap,  "TECHNICAL_PRODUCT_ID")].(int64)
				oDeviceinfo.Technical_Product_Desc   =values[cm.Getcolindex(colmap,  "TECHNICAL_PRODUCT_DESC")].(string)
				oDeviceinfo.Technical_Product_Type          =values[cm.Getcolindex(colmap,  "TECHNICAL_PRODUCT_TYPE")].(string)
				oDeviceinfo.Names               = values[cm.Getcolindex(colmap,  "NAMES")].(string)
				oDeviceinfo.Company               = values[cm.Getcolindex(colmap,  "COMPANY")].(string)
                oDeviceinfo.CustType   =values[cm.Getcolindex(colmap,  "CUSTTYPE")].(string)
			    oDeviceinfo.SiliconFlag    =values[cm.Getcolindex(colmap,  "SILICONFLAG")].(string)
                oDeviceinfo.Duallnbf            =values[cm.Getcolindex(colmap,  "DUALLNBF")].(string)
                oDeviceinfo.Mac_Address1   =values[cm.Getcolindex(colmap,  "MAC_ADDRESS1")].(string)
				oDeviceinfo.External_ID   =values[cm.Getcolindex(colmap, "EXTERNAL_ID")].(string)
				oDeviceinfo.CustomerID   =values[cm.Getcolindex(colmap, "CUSTOMERID")].(int64)
				oDeviceinfo.FinOption    =values[cm.Getcolindex(colmap, "FINOPTION")].(string)
				oDeviceinfo.DescLinkBasics   =values[cm.Getcolindex(colmap, "DESCLINKBASICS")].(string)
				oDeviceinfo.Batch_number   =values[cm.Getcolindex(colmap, "BATCH_NUMBER")].(string)
				oDeviceinfo.HardwareType  =values[cm.Getcolindex(colmap, "HARDWARETYPE")].(string)
				
				oDeviceinfocolection =append(oDeviceinfocolection,oDeviceinfo)
				 
			    	//print(err)
			}
				  ocustomerInfo.DeviceList = oDeviceinfocolection 
		//	ocustomerInfo.DeviceList =append(ocustomerInfo.DeviceList,oDeviceinfocolection)
			 
		
			//ocustomerInfo = oCustomer
         //log.Println(oCustomerinfocolection)
		}
	 
	}
	    
        oCustomerRespon.CustomerInfocollection =append(oCustomerRespon.CustomerInfocollection,  ocustomerInfo)
	    
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
 

