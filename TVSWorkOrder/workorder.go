package main

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"log"
	"strconv" 
	_ "gopkg.in/goracle.v2" 
	cm "github.com/smsdevteam/tvsglobal/common"
	c "github.com/smsdevteam/tvsglobal/tvsstructs" // referpath
)

// GetWorkorderByCustomerID get info
func GetWorkorderByCustomerID(iCustomerID string) c.WorkorderInfo {

	//resp := "SUCCESS"
	var oWorkorderinfo c.WorkorderInfo
	//var dbsource string

	dbsource := cm.GetDatasourceName("ICC")

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
		} else {
			if _, err := db.Exec(statement, intCustomerID, sql.Out{Dest: &resultC}); err != nil {
				log.Fatal(err)
			}

			defer resultC.Close()
			values := make([]driver.Value, len(resultC.Columns()))

			for {
				 //colmap := cm.createmapcol(resultC.Columns())
				//print(colmap)
				err = resultC.Next(values)
				if err == nil {
					if err == io.EOF {
						break
					}
				} else {
					break
				}
				oWorkorderinfo.Id = cm.StrToInt64(values[1].(string))
			}
		}
	}
	return oWorkorderinfo
}
