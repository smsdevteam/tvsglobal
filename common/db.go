package common

import (
	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	_ "gopkg.in/goracle.v2"
)

// DatabaseInfo struct
type DatabaseInfo struct {
	Username string
	Password string
	DBName   string
	HostIP   string
}

// DBReadConfig
func DBReadConfig(profilename string) DatabaseInfo {
	var dbInfo DatabaseInfo
	config.Load(file.NewSource(
		file.WithPath("../common/dbconfig.json"),
	))

	dbInfo.DBName = config.Get("hosts", profilename, "dbname").String("")
	dbInfo.Username = config.Get("hosts", profilename, "username").String("")
	dbInfo.Password = config.Get("hosts", profilename, "password").String("")
	dbInfo.HostIP = config.Get("hosts", profilename, "hostip").String("")

	return dbInfo
}

// GetDatasourceName
func GetDatasourceName(profilename string) string {
	var dbInfo DatabaseInfo
	dbInfo = DBReadConfig(profilename)

	var constr string
	constr = dbInfo.Username + "/" + dbInfo.Password + "@" + dbInfo.HostIP + dbInfo.DBName

	return constr
}
