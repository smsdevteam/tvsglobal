package tvsstructs

import (
	"encoding/xml"
	"time"
)

type SubmitOrderOpRequest struct {
	XMLName xml.Name `xml:"SubmitOrderOpRequest"`
	Text    string   `xml:",chardata"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Order   struct {
		Text    string `xml:",chardata"`
		Channel struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"channel"`
		OrderId struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"orderId"`
		OrderType struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"orderType"`
		EffectiveDate struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"effectiveDate"`
		DealerCode struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"dealerCode"`
	} `xml:"Order"`
	Customer struct {
		Text       string `xml:",chardata"`
		CustomerId struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"customerId"`
		Account struct {
			Text      string `xml:",chardata"`
			Xmlns     string `xml:"xmlns,attr"`
			AccountId string `xml:"accountId"`
		} `xml:"Account"`
		OU struct {
			Text       string `xml:",chardata"`
			Xmlns      string `xml:"xmlns,attr"`
			OuId       string `xml:"ouId"`
			Subscriber struct {
				Text         string `xml:",chardata"`
				Status       string `xml:"status"`
				ActivityInfo struct {
					Text           string `xml:",chardata"`
					ActivityReason string `xml:"activityReason"`
					UserText       string `xml:"userText"`
				} `xml:"activityInfo"`
				Offers []struct {
					Text               string `xml:",chardata"`
					Action             string `xml:"action"`
					EffectiveDate      string `xml:"effectiveDate"`
					ExpirationDate     string `xml:"expirationDate"`
					OfferName          string `xml:"offerName"`
					OfferInstanceId    string `xml:"offerInstanceId"`
					ServiceType        string `xml:"serviceType"`
					TargetPayChannelId string `xml:"targetPayChannelId"`
				} `xml:"offers"`
				PayChannelIdPrimary   string `xml:"payChannelIdPrimary"`
				PayChannelIdSecondary string `xml:"payChannelIdSecondary"`
				ResourceInfo          struct {
					Text             string `xml:",chardata"`
					ResourceCategory string `xml:"resourceCategory"`
					ResourceName     string `xml:"resourceName"`
					ValuesArray      string `xml:"valuesArray"`
				} `xml:"resourceInfo"`
				SubscriberAddress struct {
					Text               string `xml:",chardata"`
					Amphur             string `xml:"amphur"`
					BuildingName       string `xml:"buildingName"`
					City               string `xml:"city"`
					Country            string `xml:"country"`
					Floor              string `xml:"floor"`
					HouseNo            string `xml:"houseNo"`
					Moo                string `xml:"moo"`
					RoomNo             string `xml:"roomNo"`
					Soi                string `xml:"soi"`
					SubSoi             string `xml:"subSoi"`
					StreetName         string `xml:"streetName"`
					TimeAtAddress      string `xml:"timeAtAddress"`
					Tumbon             string `xml:"tumbon"`
					TypeOfAccomodation string `xml:"typeOfAccomodation"`
					Zip                string `xml:"zip"`
				} `xml:"subscriberAddress"`
				SubscriberGeneralInfo struct {
					Text     string `xml:",chardata"`
					Language string `xml:"language"`
					SmsLang  string `xml:"smsLang"`
				} `xml:"subscriberGeneralInfo"`
				SubscriberId       string `xml:"subscriberId"`
				SubscriberIndyName struct {
					Text               string `xml:",chardata"`
					Identification     string `xml:"identification"`
					IdentificationType string `xml:"identificationType"`
					Language           string `xml:"language"`
					HomePhone          string `xml:"homePhone"`
					Title              string `xml:"title"`
					FirstName          string `xml:"firstName"`
					LastName           string `xml:"lastName"`
					Gender             string `xml:"gender"`
				} `xml:"subscriberIndyName"`
				SubscriberNumber string `xml:"subscriberNumber"`
				SubscriberType   string `xml:"subscriberType"`
				ExtendedInfo     struct {
					Text  string `xml:",chardata"`
					Name  string `xml:"name"`
					Value string `xml:"value"`
				} `xml:"ExtendedInfo"`
			} `xml:"subscriber"`
		} `xml:"OU"`
	} `xml:"Customer"`
}
type SubmitOrderOpResponse struct {
	XMLName       xml.Name `xml:"SubmitOrderOpResponse"`
	Text          string   `xml:",chardata"`
	Xsi           string   `xml:"xsi,attr"`
	Xsd           string   `xml:"xsd,attr"`
	RespMsg       string   `xml:"respMsg"`
	RespCode      string   `xml:"respCode"`
	OMXTrackingId string   `xml:"OMXTrackingId"`
}
type TVSBNProperty struct {
	TRNSEQNO         string
	TVSCUSTOMERNO    int
	TVSAccountno     string
	Reftype          string
	Refvalue         string
	TVSCustomerType  string
	CCBSCustomerno   string
	CCBSAccountno    string
	CCBSOUNo         string
	CCBSSubNo        string
	CCBSAGREEMENTID  string
	OLDCCBSACCOUNTNO string
	CCBSorderno      string
	CCBSFN           string
	CCBSORDERTYPEID  string
	CCBSSUBFN        string
	CCBSACTIVITYREON string
	CCBSUSERTEXT     string
	CCBSURLSERVICE   string
	CCBSarURLSERVICE string
	SHNO             string
	CREATEDON        string
	PROCESSFLAG      string
	ERRORCODE        string
	ERRORDESC        string
	REQUESTDATA      string
	RESPONSEDATA     string
	STARTCALL        time.Time
	STOPCALL         time.Time
	activityReon     string
	usertext         string
	OMXTrackingID    string
	OMXOrderType     string
	omxrespMsg       string
	omxrespCode      string
	HAVEOCCHARE      string
	RECALLOFFER      string
	//TVS_BN_OChargePropertylist  List(Of TVS_BN_OChargeProperty)
	TVSBNCCBSOfferPropertylist []TVSBNCCBSOfferProperty
	//TVS_BN_CCBSOffer_OU  List(Of TVS_BN_CCBSOfferProperty)
	//TVS_BN_CCBSOffer_Prepaid  List(Of TVS_BN_CCBSOfferProperty)
	ExternalRefno    string
	FINDEXISTSCUST   string
	LEGACYBAN        string
	INITREON         string
	OLDBANDATE       string
	FoundTVSonCCBS   bool
	ToCCBSCustomerno string
	Oldccbssubno     string
	AddSOCLevelOU    string
}

type TVSBNCCBSOfferProperty struct {
	action                   string
	ccbsoffername            string
	ccbssocid                string
	offerInstanceID          string
	effectivedate            time.Time
	effectiveDateSpecified   int
	expirationdate           time.Time
	processtype              string
	targetPayChannelID       string
	OverrideRCAmount         float32
	OverrideRCDescription    string
	OverrideRCDescriptionEng string
	OverrideRCSpecified      float32
	OverrideOCAmount         float32
	OverrideOCDescription    string
	OverrideOCDescriptionEng string
	OverrideOCSpecified      int
	Newperiodind             string
	extendedinfoname         string
	extendedinfovalue        string
	ccbsservicetype          string
}
