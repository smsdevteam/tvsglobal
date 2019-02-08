package common

import (
	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

type AppServiceURL struct {
	ICCServiceURL string `json:"iccserviceURL"`
}

func AppReadConfig(profilename string) AppServiceURL {

	var appServiceURL AppServiceURL
	config.Load(file.NewSource(
		file.WithPath("../common/appconfig.json"),
	))

	appServiceURL.ICCServiceURL = config.Get(profilename, "service", "iccServiceURL").String("")

	return appServiceURL
}
