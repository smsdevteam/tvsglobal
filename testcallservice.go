package main
//test 2 times arm
import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	auth, err := queryAuthentication()
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(auth)
	}
	fmt.Println("*********************************************************************************")
	agreement, err_ar := queryagreementdetail()
	if err_ar != nil {
		log.Println(err_ar)
	} else {
		fmt.Println(agreement)
	}
}

const getTemplate = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
       <s:Body><AuthenticateByProof xmlns="http://ibs.entriq.net/Authentication">
        <userIdentity xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
        <Dsns>
         <Dsn>
          <Extended i:nil="true"/>
          <Name>UAT62</Name>
         </Dsn>
        </Dsns>
        <Extended i:nil="true"/>
        <UserName>ICCAPI</UserName>
        </userIdentity>
        <proof>iccdemon@1</proof>
        <module i:nil="true" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"/>
        </AuthenticateByProof>
       </s:Body>
       </s:Envelope>`

const getTemplateforagreementdetail = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Header>
  <h:CacheControlHeader i:nil="true" xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns:h="http://ibs.entriq.net/Core" />
  <h:AuthenticationHeader xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns:h="http://ibs.entriq.net/Security">
	<h:ClientName i:nil="true" />
	<h:ClientProof i:nil="true" />
	<h:Culture i:nil="true" />
	<h:Dsn>UAT62</h:Dsn>
	<h:Extended i:nil="true" />
	<h:ExternalAgent i:nil="true" />
	<h:Proof>iccdemon@1</h:Proof>
	<h:ServerTime>2018-11-29T16:33:00</h:ServerTime>
	<h:Token>886326529B40BED841425E59C7E75B6B000070CFD37ADD55D688</h:Token>
	<h:UserName>ICCAPI</h:UserName>
  </h:AuthenticationHeader>
</s:Header>
<s:Body>
  <GetAgreementDetails xmlns="http://ibs.entriq.net/AgreementManagement">
	<agreementId>102314173</agreementId>
	<page>1</page>
  </GetAgreementDetails>
</s:Body>
</s:Envelope>`

func queryAuthentication() (string, error) {
	url := "http://tv-uatibs62-w01.ubc.co.th/ASM/ALL/Authentication.svc"
	client := &http.Client{}

	requestContent := []byte(getTemplate)
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

	return string(contents), nil
}
func queryagreementdetail() (string, error) {
	url := "http://172.22.247.125/ASM/ALL/AgreementManagement.svc"
	client := &http.Client{}

	requestContent := []byte(getTemplateforagreementdetail)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		return "", err
	}
	///                            http://ibs.entriq.net/Authentication/IAuthenticationService/AuthenticateByProof
	//http://ibs.entriq.net/AgreementManagement/IAgreementManagementService.GetAgreementDetail
	req.Header.Add("SOAPAction", `"http://ibs.entriq.net/AgreementManagement/IAgreementManagementService/GetAgreementDetails"`)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("sssss")
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("200")
		return "", errors.New("Error Respose " + resp.Status)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
// edit by nattachais
