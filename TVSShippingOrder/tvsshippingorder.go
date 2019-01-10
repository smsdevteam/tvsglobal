package main

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"log"

	_ "gopkg.in/goracle.v2"

	cm "github.com/pimpina/tvsglobalb/Common"     // db
	so "github.com/pimpina/tvsglobalb/TVSStructs" // referpath
)

// GetShippingOrder
func GetShippingOrder(iOrderID int64) so.ShippingOrderRes {
	//db, err := sql.Open("goracle", "bgweb/bgweb#1@//tv-uat62-dq.tvsit.co.th:1521/UAT62")
	var dbsource string
	dbsource = cm.GetDatasourceName("ICC")
	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var statement string
	var resultC driver.Rows

	// Shipping Order Header
	statement = "begin redibsservice.getdatashippingorderheader(:0,:1); end;"

	if _, err := db.Exec(statement, iOrderID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
	}

	defer resultC.Close()
	values := make([]driver.Value, len(resultC.Columns()))

	var oSO so.ShippingOrderRes
	for {
		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
		}
		oSO.ID = values[0].(int64)
		oSO.DepotFrom = values[1].(int64)
		oSO.DepotTo = values[2].(int64)
		oSO.StatusID = values[3].(int64)
		oSO.StatusDesc = values[4].(string)
		oSO.TypeID = values[5].(int64)
		oSO.TypeDesc = values[6].(string)
		oSO.CreateComments = values[7].(string)
		oSO.CreateReference = values[8].(string)
		oSO.CreateDateTime = values[9].(string)
		oSO.CreateBy = values[10].(int64)
		oSO.CreateByName = values[11].(string)
	}

	// Shipping Order Line
	statement = "begin redibsservice.getdatashippingorderline(:0,:1); end;"

	if _, err := db.Exec(statement, iOrderID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
	}

	defer resultC.Close()
	values = make([]driver.Value, len(resultC.Columns()))

	var oSLList []so.ShippingOrderLineRes

	for {
		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
		}
		var oSL so.ShippingOrderLineRes
		oSL.LineID = values[0].(int64)
		oSL.ShippingOrderID = values[1].(int64)
		oSL.LineNr = values[2].(int64)
		oSL.ProductID = values[3].(int64)
		oSL.ProductKey = values[4].(string)
		oSL.ModelID = values[5].(int64)
		oSL.ModelKey = values[6].(string)
		oSL.Qty = values[7].(int64)
		oSLList = append(oSLList, oSL)
	}
	oSO.ShippingOrderLines = oSLList

	// Shipping Device
	statement = "begin redibsservice.getdatashippingordersn(:0,:1); end;"

	if _, err := db.Exec(statement, iOrderID, sql.Out{Dest: &resultC}); err != nil {
		log.Fatal(err)
	}

	defer resultC.Close()
	values = make([]driver.Value, len(resultC.Columns()))

	var oSDList []so.ShippingDeviceRes

	for {
		err = resultC.Next(values)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error:", err)
		}
		var oSD so.ShippingDeviceRes
		oSD.ShippingOrderID = values[0].(int64)
		oSD.LineID = values[1].(int64)
		oSD.SerialNumber = values[2].(string)
		oSD.StatusID = values[3].(int64)
		oSD.DVResult = values[4].(string)
		oSDList = append(oSDList, oSD)
	}
	oSO.ShippingDevices = oSDList

	return oSO
}

/*
func main() {
	var Val int64
	fmt.Printf("input : ")
	fmt.Scan(&Val)
	r := GetShippingOrder(Val)
	//r := GetShippingOrder(16301898)
	fmt.Println(r)

	json.NewEncoder(os.Stdout).Encode(r)
}
*/
