package common

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"

	"net/http"
)

//Envelope is strcut for get authen
type Envelope struct {
	//XMLName xml.Name `xml:"Envelope"`
	Text string `xml:",chardata"`
	S    string `xml:"xmlns,attr"`
	Body struct {
		Text                string `xml:",chardata"`
		AuthenticateByProof struct {
			Text         string `xml:",chardata"`
			Xmlns        string `xml:"xmlns,attr"`
			UserIdentity struct {
				Text string `xml:",chardata"`
				I    string `xml:"xmlns:i,attr"`
				Dsns struct {
					Text string `xml:",chardata"`
					Dsn  struct {
						Text     string `xml:",chardata"`
						Extended struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"i:nil,attr"`
						} `xml:"Extended"`
						Name string `xml:"Name"`
					} `xml:"Dsn"`
				} `xml:"Dsns"`
				Extended struct {
					Text string `xml:",chardata"`
					Nil  string `xml:"i:nil,attr"`
				} `xml:"Extended"`
				UserName string `xml:"UserName"`
			} `xml:"userIdentity"`
			Proof  string `xml:"proof"`
			Module struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"i:nil,attr"`
				I    string `xml:"xmlns:i,attr"`
			} `xml:"module"`
		} `xml:"AuthenticateByProof"`
	} `xml:"Body"`
}
type Enveloperes struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"s,attr"`
	Header  struct {
		Text                 string `xml:",chardata"`
		AuthenticationHeader struct {
			Text       string `xml:",chardata"`
			Xmlns      string `xml:"xmlns,attr"`
			I          string `xml:"i,attr"`
			ClientName struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"ClientName"`
			ClientProof struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"ClientProof"`
			Culture struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"Culture"`
			Dsn      string `xml:"Dsn"`
			Extended struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"Extended"`
			ExternalAgent struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"ExternalAgent"`
			Proof struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"Proof"`
			ServerTime string `xml:"ServerTime"`
			Token      string `xml:"Token"`
			UserName   string `xml:"UserName"`
		} `xml:"AuthenticationHeader"`
	} `xml:"Header"`
	Body struct {
		Text                        string `xml:",chardata"`
		AuthenticateByProofResponse struct {
			Text                      string `xml:",chardata"`
			Xmlns                     string `xml:"xmlns,attr"`
			AuthenticateByProofResult struct {
				Text          string `xml:",chardata"`
				I             string `xml:"i,attr"`
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
			} `xml:"AuthenticateByProofResult"`
		} `xml:"AuthenticateByProofResponse"`
	} `xml:"Body"`
}

const getTemplate = `<Envelope xmlns ="http://schemas.xmlsoap.org/soap/envelope/">
       <Body><AuthenticateByProof xmlns="http://ibs.entriq.net/Authentication">
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
       </Body>
	   </Envelope>`
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
		   <h:Token>6EE947AA505190DE38AD9A1088DB86AF0000DAE65A211E81D688</h:Token>
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
//Getauthentication is Getauthentication
func Createxmlstruct(filename str ){
	
}
func Getauthentication() (string, error) {
	url := "http://tv-uatibs62-w01.ubc.co.th/ASM/ALL/Authentication.svc"
	client := &http.Client{}

	var sreq Envelope
	sreq.S = "http://schemas.xmlsoap.org/soap/envelope/"
	sreq.Body.AuthenticateByProof.UserIdentity.Dsns.Dsn.Name = "UAT62"
	sreq.Body.AuthenticateByProof.UserIdentity.Dsns.Dsn.Extended.Nil = "true"
	sreq.Body.AuthenticateByProof.UserIdentity.Extended.Nil = "true"
	sreq.Body.AuthenticateByProof.UserIdentity.UserName = "ICCAPI"
	sreq.Body.AuthenticateByProof.UserIdentity.I = "http://www.w3.org/2001/XMLSchema-instance"
	sreq.Body.AuthenticateByProof.Proof = "iccdemon@1"
	sreq.Body.AuthenticateByProof.Xmlns = "http://ibs.entriq.net/Authentication"
	sreq.Body.AuthenticateByProof.Module.I = "http://www.w3.org/2001/XMLSchema-instance"
	sreq.Body.AuthenticateByProof.Module.Nil = "true"
	sreq.Body.AuthenticateByProof.UserIdentity.Extended.Nil = "true"
	output, err := xml.MarshalIndent(sreq, "  ", "    ")
	a := string(output)
	Writelogfile(a)
	//fmt.Print(a)
	requestContent := []byte(a) //[]byte(getTemplate) //
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
	v := Enveloperes{}
	err = xml.Unmarshal(contents, &v)
	//fmt.Printf(v.Body.AuthenticateByProofResponse.AuthenticateByProofResult.Token)
	if err != nil {
		fmt.Printf("error: %v", err)

	}
	return v.Body.AuthenticateByProofResponse.AuthenticateByProofResult.Token, nil
}

func Excuteserv(xmlrequest string) (string, error) {
	url := "http://172.22.247.125/ASM/ALL/AgreementManagement.svc"
	client := &http.Client{}

	requestContent :=[]byte(xmlrequest)// []byte(xmlrequest)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		return "", err
	}
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
