package tvsstructs

type GetCCBSSubscriberInfo struct {
	//XMLName           xml.Name `xml:"GetCCBSSubscriberInfo"`
	Text              string `xml:",chardata"`
	Xmlns             string `xml:"xmlns,attr"`
	ClientInformation struct {
		Text    string `xml:",chardata"`
		D4p1    string `xml:"d4p1,attr"`
		I       string `xml:"i,attr"`
		AppInfo struct {
			Text        string `xml:",chardata"`
			D5p1        string `xml:"d5p1,attr"`
			AppFunction struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"App_Function"`
			AppName struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"App_Name"`
			AppServerName struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"App_ServerName"`
			AppStartTime struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"App_StartTime"`
			AppTrnID struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"App_TrnID"`
			AppVersion struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"App_Version"`
		} `xml:"AppInfo"`
		ClientUserInfo struct {
			Text            string `xml:",chardata"`
			D5p1            string `xml:"d5p1,attr"`
			ExternalAgentId struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"ExternalAgentId"`
			UserSessionId struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"UserSessionId"`
			Userid struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"Userid"`
		} `xml:"ClientUserInfo"`
		CustomerNo struct {
			Text string `xml:",chardata"`
			Nil  string `xml:"nil,attr"`
		} `xml:"CustomerNo"`
	} `xml:"ClientInformation"`
	SubscriberId struct {
		Text string `xml:",chardata"`
	} `xml:"SubscriberId"`
}
