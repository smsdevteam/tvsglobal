package main

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"log"
	"strconv"

	cm "github.com/smsdevteam/tvsglobal/common"
	c "github.com/smsdevteam/tvsglobal/tvsstructs" // referpath
	_ "gopkg.in/goracle.v2"
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
				colmap := cm.Createmapcol(resultC.Columns())
				print(colmap)
				err = resultC.Next(values)
				if err == nil {
					if err == io.EOF {
						break
					}
				} else {
					break
				}
				//oWorkorderinfo.Id = cm.StrToInt64(values[1].(string))
				//print(values[cm.Getcolindex(colmap, "PROBLEM_DESCRIPTION")].(string))
				oWorkorderinfo.ID = values[cm.Getcolindex(colmap, "ID")].(int64)
				//oWorkorderinfo.ProblemDesc = values[cm.Getcolindex(colmap, "Problem_Description")].(string)
			}
		}
	}
	return oWorkorderinfo
}
// GetWorkorderByworkorderid get info
func GetWorkorderByworkorderid(iworkorderID string) c.WorkorderInfo {

	//resp := "SUCCESS"
	var oWorkorderinfo c.WorkorderInfo
	var sworkoderid string 
	//var dbsource string

	dbsource := cm.GetDatasourceName("ICC")

	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
		//resp = err.Error() //
	} else {
		defer db.Close()
		var statement string
		statement = "begin tvs_workorder.getworkorderbyid(:0,:1); end;"
		var resultC driver.Rows
		intworkorderID, err := strconv.Atoi(iworkorderID)
		if err != nil {
			log.Fatal(err)
		} else {
			if _, err := db.Exec(statement, intworkorderID, sql.Out{Dest: &resultC}); err != nil {
				log.Fatal(err)
			}

			defer resultC.Close()
			values := make([]driver.Value, len(resultC.Columns()))

			for {
				colmap := cm.Createmapcol(resultC.Columns())
				print(colmap)
				err = resultC.Next(values)
				if err == nil {
					if err == io.EOF {
						break
					}
				} else {
					break
				}
				sworkoderid  = cm.Int64ToStr( values[cm.Getcolindex(colmap, "ID")].(int64))
				//oWorkorderinfo.Id = cm.StrToInt64(values[1].(string))
				//print(values[cm.Getcolindex(colmap, "PROBLEM_DESCRIPTION")].(string))
				oWorkorderinfo.ID = cm.StrToInt64( sworkoderid   )
				oWorkorderinfo.WorkorderServiceDTlist =Getworkorderservicebyid(sworkoderid)
				//oWorkorderinfo.ProblemDesc = values[cm.Getcolindex(colmap, "Problem_Description")].(string)
			}
		}
	}
	return oWorkorderinfo
}
// Getworkorderservicebyid get info
func Getworkorderservicebyid(iworkorderID string) []c.WorkorderServiceDTInfo {

	//resp := "SUCCESS"
	var oWorkorderServiceinfo c.WorkorderServiceDTInfo
	var oWorkorderServicecolection  []c.WorkorderServiceDTInfo
	//var dbsource string

	dbsource := cm.GetDatasourceName("ICC")

	db, err := sql.Open("goracle", dbsource)
	if err != nil {
		log.Fatal(err)
		//resp = err.Error() //
	} else {
		defer db.Close()
		var statement string
		statement = "begin tvs_workorder.getworkorderservicebyid(:0,:1); end;"
		var resultC driver.Rows
		intworkorderID, err := strconv.Atoi(iworkorderID)
		if err != nil {
			log.Fatal(err)
		} else {
			if _, err := db.Exec(statement, intworkorderID, sql.Out{Dest: &resultC}); err != nil {
				log.Fatal(err)
			}

			defer resultC.Close()
			values := make([]driver.Value, len(resultC.Columns()))

			for {
				colmap := cm.Createmapcol(resultC.Columns())
				print(colmap)
				err = resultC.Next(values)
				if err == nil {
					if err == io.EOF {
						break
					}
				} else {
					break
				}
				//oWorkorderinfo.Id = cm.StrToInt64(values[1].(string))
				//print(values[cm.Getcolindex(colmap, "PROBLEM_DESCRIPTION")].(string))
				oWorkorderServiceinfo.ServiceId = values[cm.Getcolindex(colmap, "SERVICE_ID")].(int64)
		 	    oWorkorderServiceinfo.ServiceDescription = values[cm.Getcolindex(colmap, "DESCRIPTION")].(string)
			 
				oWorkorderServicecolection = append(oWorkorderServicecolection, oWorkorderServiceinfo)

				//oWorkorderinfo.ProblemDesc = values[cm.Getcolindex(colmap, "Problem_Description")].(string)
			}

		}
	}
	return oWorkorderServicecolection
}
