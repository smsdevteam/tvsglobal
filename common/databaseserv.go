package common

import (
	"fmt"
	"log"

	//"log"
	//"os"
	//"path/filepath"
	"database/sql"

	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	_ "gopkg.in/goracle.v2"
)

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
func ExecutetData(profilename string, sqlStmnt string) (string, error) {
	var dbname = " "
	var username = ""
	var password = ""
	var connectionstring = ""
	var jreSult = ""
	//var err error
	//var reSult=""
	//var db DB
	fmt.Print(connectionstring)
	dbname, username, password, connectionstring = readconfig(profilename)
	db, err := sql.Open("goracle", username+"/"+password+"@"+dbname)
	reSult, err := db.Exec(sqlStmnt)
	fmt.Print(reSult)
	jreSult = "-"
	if err != nil {
		jreSult = "error"
	} else {
		jreSult = "success"
	}
	defer db.Close()

	return jreSult, err
}
func main1() {
	var err error
	var result = " "
	result, err = ExecutetData("ICC", "INSERT INTO TEMP_OOD VALUES('111')")
	var dbname = " "
	var username = ""
	var password = ""
	var connectionstring = ""
	fmt.Println("start" + result)
	dbname, username, password, connectionstring = readconfig("ICC")

	fmt.Println(dbname + username + password + connectionstring)
	db, err := sql.Open("goracle", username+"/"+password+"@"+dbname)
	if err != nil {
		fmt.Println("Hell 77666, fail ")
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Hell 77666, complete ")

	config.Load(file.NewSource(
		file.WithPath("dbconfig.json"),
	))

}
