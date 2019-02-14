package common

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"

	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

type product struct {
	ID      int64  `json:"id"`
	Eventid int64  `json:"event_id"`
	Descrip string `json:"description"`
}

func readconfig(profilename string) (string, string, string, string) {
	var dbname = " "
	var username = ""
	var password = ""
	var connectionstring = ""
	config.Load(file.NewSource(
		file.WithPath("dbconfig.json"),
	))
	dbname = config.Get("hosts", profilename, "dbname").String("")
	username = config.Get("hosts", profilename, "username").String("")
	password = config.Get("hosts", profilename, "password").String("")
	connectionstring = config.Get("hosts", profilename, "connectionstring").String("")
	return dbname, username, password, connectionstring
}

//Connecttodb is function excute sql command
func Connecttodb(profilename string) (*sql.DB, error) {
	var dbname = " "
	var username = ""
	var password = ""
	var connectionstring = ""
	fmt.Print(connectionstring)
	dbname, username, password, connectionstring = readconfig(profilename)
	db, err := sql.Open("goracle", username+"/"+password+"@"+dbname)

	return db, err
}

//ExcutestoreDS is function excute sql command
func ExcutestoreDS(profilename string, sqlStmnt string, args ...interface{}) (driver.Rows, error) {
	var resultI driver.Rows
	db, err := Connecttodb(profilename)
	defer db.Close()
	//var resultI driver.Rows
	if err != nil {
		Writelogfile("connect db Error : |command :" + sqlStmnt + "|Error : " + err.Error() +
			"|dbname: " + profilename + "\n")
	} else {
		if _, err := db.Exec(sqlStmnt, args...); err != nil {
			Writelogfile("Found Error Execute " + err.Error())
		}
	}
	return resultI, err
}

//ExecuteCMD is function excute sql command
func ExecuteCMD(profilename string, sqlStmnt string) (string, error) {
	var dbname = " "
	var jreSult = ""
	db, err := Connecttodb(profilename)
	defer db.Close()
	if err != nil {
		Writelogfile("connect db Error : |command :" + sqlStmnt + "|Error : " + err.Error() +
			"|dbname: " + profilename + "\n")
	} else {
		result, err := db.Exec(sqlStmnt)
		result.RowsAffected()
		result = nil
		jreSult = "success"
		if err != nil {
			Writelogfile("ExecutetData Error : |command :" + sqlStmnt + "|Error : " + err.Error() +
				"|dbname: " + dbname + "\n")
			jreSult = "fail"
		}
	}
	return jreSult, err
}

//Getdatalist is function excute sql command
func Getdatalist(profilename string, sqlStmnt string) (*sql.Rows, error) {
	db, err := Connecttodb(profilename)

	defer db.Close()
	if err != nil {
		Writelogfile("connect db Error : |command :" + sqlStmnt + "|Error : " + err.Error() +
			"|dbname: " + profilename + "\n")
	} else {
		rows, err := db.Query(sqlStmnt)
		if err != nil {
			Writelogfile("Query  : |command :" + sqlStmnt + "|Error : " + err.Error() +
				"|profilename: " + profilename + "\n")
		}
		for rows.Next() {
			var intCol, strCol string

			if err := rows.Scan(&intCol, &strCol); err != nil {
				Writelogfile(err.Error())
				break
			} else {
				fmt.Printf("IntCol=%s, StrCol=%s\n", intCol, strCol)
			}

		}
		return rows, err
	}
	//values := make([]driver.Value, len(rows.Columns()))

	return nil, err
}

//Excutestore is function excute sql command
func Excutestore(profilename string, sqlStmnt string) (string, error) {
	var oProduct []product
	db, err := Connecttodb(profilename)
	defer db.Close()
	//var resultI driver.Rows
	if err != nil {
		Writelogfile("connect db Error : |command :" + sqlStmnt + "|Error : " + err.Error() +
			"|dbname: " + profilename + "\n")
	} else {
		Writelogfile("connect db complete")
		var resultI driver.Rows
		if _, err := db.Exec(sqlStmnt, 142, sql.Out{Dest: &resultI}); err != nil {
			fmt.Println(err)
		}
		defer resultI.Close()
		values := make([]driver.Value, len(resultI.Columns()))
		for {
			err = resultI.Next(values)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("error:", err)
			}
			var lProduct product
			lProduct.ID = values[0].(int64)
			lProduct.Descrip = values[2].(string)
			oProduct = append(oProduct, lProduct)

		}
		fmt.Println(oProduct)
		fmt.Println("End..")
	}
	return "!", err
}
func Getvalue(value []driver.Value, colmap map[string]int, colname string) interface{} {

	i, err := colmap[colname]
	if err == true {
		return value[i]
	} else {
		return ""
	}

}
func Getcolindex(colmap map[string]int, colname string) int {

	i, err := colmap[colname]
	if err == true {
		return i
	} else {
		return -1
	}

}
func Createmapcol(data []string) map[string]int {
	var colmap = map[string]int{}

	for k, v := range data {
		colmap[v] = k
	}
	return colmap
}
