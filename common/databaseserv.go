package common

import (
	"database/sql"
	"fmt"

	config "github.com/micro/go-config"
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

//ExecutetData is function excute sql command
func ExecutetData(profilename string, sqlStmnt string) (string, error) {
	var dbname = " "
	var username = ""
	var password = ""
	var connectionstring = ""
	var jreSult = ""
	fmt.Print(connectionstring)
	dbname, username, password, connectionstring = readconfig(profilename)
	db, err := sql.Open("goracle", username+"/"+password+"@"+dbname)
	reSult, err := db.Exec(sqlStmnt)
	fmt.Print(dbname, username, password, connectionstring, reSult)
	jreSult = "-"
	if err != nil {
		jreSult = "error" + err.Error()
		Writelogfile()
	} else {
		jreSult = "success"
	}
	defer db.Close()

	return jreSult, err
}
