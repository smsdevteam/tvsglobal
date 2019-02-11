package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"

	//"fmt"

	//"strconv"
	cm "tvsglobal/common"
	st "tvsglobal/tvsstructs"

	//"github.com/jmoiron/sqlx"
	"net/http"

	_ "gopkg.in/goracle.v2"
)

func osbgetsubinfo(ccbssubno int) (string, error) {
	url := "http://172.22.203.63/TVS_GlobalWCFuat/CCBS_OSB_FinanceService.svc"
	client := &http.Client{}
	var subinfo st.GetCCBSSubscriberInfo
	subinfo.SubscriberId.Text = cm.IntToStr(ccbssubno)
	output, err := xml.MarshalIndent(subinfo, "  ", "    ")
	a := string(output)
	a = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/" xmlns:tvs="http://schemas.datacontract.org/2004/07/TVS_GlobalProperty.CommonProperty" xmlns:tvs1="http://schemas.datacontract.org/2004/07/TVS_Public">
	<soapenv:Header/>
	<soapenv:Body>
	   <tem:GetCCBSSubscriberInfo>
		  <tem:ClientInformation>
			 <tvs:AppInfo>
				<tvs1:App_Function>go-getsub</tvs1:App_Function>
				<tvs1:App_Name>go</tvs1:App_Name>
				<tvs1:App_ServerName>go</tvs1:App_ServerName>
				<tvs1:App_StartTime>?</tvs1:App_StartTime>
				<tvs1:App_TrnID>?</tvs1:App_TrnID>
				<tvs1:App_Version>?</tvs1:App_Version>
			 </tvs:AppInfo>
			 <tvs:ClientUserInfo>
				<tvs1:ExternalAgentId>1</tvs1:ExternalAgentId>
				<tvs1:UserSessionId>?</tvs1:UserSessionId>
				<tvs1:Userid>?</tvs1:Userid>
			 </tvs:ClientUserInfo>
			 <tvs:CustomerNo>?</tvs:CustomerNo>
		  </tem:ClientInformation>
		  <tem:SubscriberId>` + cm.IntToStr(ccbssubno) + `</tem:SubscriberId>
		  </tem:GetCCBSSubscriberInfo>
	   </soapenv:Body>
	</soapenv:Envelope>`
	//fmt.Println(a)
	//fmt.Println("********************************************************")
	requestContent := []byte(a)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("SOAPAction", `http://tempuri.org/ICCBS_OSB_FinanceService/GetCCBSSubscriberInfo`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(errors.New("Error Respose " + resp.Status))
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	//m, _ := mxj.NewMapXml(contents, true)
	//fmt.Println(&m)
	return string(contents), nil
}
func main() {
	osbgetsubinfo(1337981)
	print("aaa")
}
