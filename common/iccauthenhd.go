package common

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"

	s "strings"

	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

// ICCAuthenHD struct
type ICCAuthenHD struct {
	ServiceDSN          string
	ServiceLocationURL  string
	ServiceUser         string
	ServiceUserIdentity string
}

type ServiceURL struct {
	AgreementURL     string `json:"agreementURL"`
	AuthenURL        string `json:"authenURL"`
	DeviceURL        string `json:"deviceURL"`
	ShippingOrderURL string `json:"shippingorderURL"`
}

// MyRespEnvelope
type MyRespEnvelope struct {
	XMLName xml.Name
	Body    body
}

type body struct {
	XMLName     xml.Name
	GetResponse completeResponse `xml:"AuthenticateByProofResponse"`
}

type completeResponse struct {
	XMLName  xml.Name `xml:"AuthenticateByProofResponse"`
	MyResult authenHD `xml:"AuthenticateByProofResult"`
}

type authenHD struct {
	Authenticated string `xml:"Authenticated"`
	Dsns          struct {
		Text      string `xml:",chardata"`
		DsnStatus struct {
			Text     string `xml:",chardata"`
			Extended struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"Extended"`
			IsValid string `xml:"IsValid"`
			Name    string `xml:"Name"`
			Status  string `xml:"Status"`
		} `xml:"DsnStatus"`
	} `xml:"Dsns"`
	Extended struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"Extended"`
	FailureReason struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"FailureReason"`
	Token    string `xml:"Token"`
	UserName string `xml:"UserName"`
}

// ICCReadConfig
func ICCReadConfig(profilename string) (ICCAuthenHD, ServiceURL) {
	var authenInfo ICCAuthenHD
	var serviceLnk ServiceURL
	config.Load(file.NewSource(
		file.WithPath("../iccconfig.json"),
	))

	authenInfo.ServiceDSN = config.Get(profilename, "authen", "ServiceDSN").String("")
	authenInfo.ServiceLocationURL = config.Get(profilename, "authen", "ServiceLocationURL").String("")
	authenInfo.ServiceUser = config.Get(profilename, "authen", "ServiceUser").String("")
	authenInfo.ServiceUserIdentity = config.Get(profilename, "authen", "ServiceUserIdentity").String("")

	serviceLnk.AgreementURL = config.Get(profilename, "service", "agreementURL").String("")
	serviceLnk.AuthenURL = config.Get(profilename, "service", "authenURL").String("")
	serviceLnk.DeviceURL = config.Get(profilename, "service", "deviceURL").String("")
	serviceLnk.ShippingOrderURL = config.Get(profilename, "service", "shippingorderURL").String("")

	return authenInfo, serviceLnk
}

// GetICCAuthenToken
func GetICCAuthenToken(profilename string) (string, error) {
	var authenInfo ICCAuthenHD
	var serviceLnk ServiceURL
	authenInfo, serviceLnk = ICCReadConfig(profilename)

	//fmt.Println(authenInfo)
	url := serviceLnk.AuthenURL
	client := &http.Client{}

	requestValue := s.Replace(getTemplate, "$DSN", authenInfo.ServiceDSN, -1)
	requestValue = s.Replace(requestValue, "$username", authenInfo.ServiceUser, -1)
	requestValue = s.Replace(requestValue, "$password", authenInfo.ServiceUserIdentity, -1)
	//fmt.Println(requestValue)
	requestContent := []byte(requestValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		return "", err
	}

	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/Authentication/IAuthenticationService/AuthenticateByProof"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("Error Respose " + resp.Status)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//fmt.Println(string(contents))

	authenHD := MyRespEnvelope{}
	xml.Unmarshal([]byte(contents), &authenHD)
	Token := authenHD.Body.GetResponse.MyResult.Token
	//ServerTime := authenHD.Body.GetResponse.MyResult.Token

	return Token, nil

}
