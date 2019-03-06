package tvsstructs

import (
	"encoding/xml"
	"time"
)

type SubmitOrderOpRequest struct {
	XMLName xml.Name `xml:"submitOrderRequest"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"xmlns,attr"`
	SE      string   `xml:"xmlns:SOAP-ENV,attr"`
	XSD     string   `xml:"xmlns:xsd,attr"`

	Order struct {
		Text    string `xml:",chardata"`
		Channel struct {
			Text string `xml:",chardata"`
		} `xml:"channel"`
		OrderId struct {
			Text string `xml:",chardata"`
		} `xml:"orderId"`
		OrderType struct {
			Text string `xml:",chardata"`
		} `xml:"orderType"`
		EffectiveDate struct {
			Text string `xml:",chardata"`
		} `xml:"effectiveDate"`
		DealerCode struct {
			Text string `xml:",chardata"`
		} `xml:"dealerCode"`
	} `xml:"Order"`
	Customer struct {
		Text       string `xml:",chardata"`
		CustomerId struct {
			Text string `xml:",chardata"`
		} `xml:"customerId"`
		Account struct {
			Text      string `xml:",chardata"`
			AccountId string `xml:"accountId"`
		} `xml:"Account"`
		OU struct {
			Text       string `xml:",chardata"`
			OuId       string `xml:"ouId"`
			Subscriber struct {
				Text         string `xml:",chardata"`
				Status       string `xml:"status"`
				ActivityInfo struct {
					Text           string `xml:",chardata"`
					ActivityReason string `xml:"activityReason"`
					UserText       string `xml:"userText"`
				} `xml:"activityInfo"`
				Offers []Omxccbsoffer `xml:"offers"`
				/*Offers []struct {
					Text               string `xml:",chardata"`
					Action             string `xml:"action"`
					EffectiveDate      string `xml:"effectiveDate"`
					ExpirationDate     string `xml:"expirationDate"`
					OfferName          string `xml:"offerName"`
					OfferInstanceId    string `xml:"offerInstanceId"`
					ServiceType        string `xml:"serviceType"`
					TargetPayChannelId string `xml:"targetPayChannelId"`
				} `xml:"offers"`
				*/
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
				SubscriberNumber string        `xml:"subscriberNumber"`
				SubscriberType   string        `xml:"subscriberType"`
				ExtendedInfoobj  *ExtendedInfo `xml:"ExtendedInfo,omitempty"`
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
	ExternalRefno       string
	FINDEXISTSCUST      string
	LEGACYBAN           string
	INITREON            string
	OLDBANDATE          string
	FoundTVSonCCBS      bool
	ToCCBSCustomerno    string
	Oldccbssubno        string
	AddSOCLevelOU       string
	TVSBNOMXPropertyobj TVSBNOMXProperty
}

type TVSBNCCBSOfferProperty struct {
	Action                   string
	Ccbsoffername            string
	Ccbssocid                string
	OfferInstanceID          string
	Effectivedate            string
	EffectiveDateSpecified   int64
	Expirationdate           string
	Processtype              string
	TargetPayChannelID       int64
	OverrideRCAmount         float64
	OverrideRCDescription    string
	OverrideRCDescriptionEng string
	OverrideRCSpecified      float64
	OverrideOCAmount         float64
	OverrideOCDescription    string
	OverrideOCDescriptionEng string
	OverrideOCSpecified      int64
	Newperiodind             string
	Extendedinfoname         string
	Extendedinfovalue        string
	Ccbsservicetype          string
}
type TVSBNOMXProperty struct {
	Channel                string
	DealerCode             string
	EffectiveDateSpecified int
	EffectiveDate          string
}
type Omxccbsoffer struct {
	Text            string             `xml:",chardata"`
	Action          string             `xml:"action,omitempty"`
	EffectiveDate   string             `xml:"effectiveDate,omitempty"` //iso 8601
	ExpirationDate  string             `xml:"expirationDate,omitempty"`
	OfferName       string             `xml:"offerName"`
	OfferInstanceId string             `xml:"offerInstanceId,omitempty"`
	ServiceType     string             `xml:"serviceType"`
	Offerparas      []Omxccbsofferpara `xml:"omxParameterInfo"`
}
type Omxccbsofferpara struct {
	Text          string `xml:",chardata"`
	ParamName     string `xml:"paramName,omitempty"`
	ValuesArray   string `xml:"valuesArray,omitempty"`
	EffectiveDate string `xml:"effectiveDate,omitempty"` //iso 8601
}
type ExtendedInfo struct {
	Text  string `xml:",omitempty"`
	Name  string `xml:"name,omitempty"`
	Value string `xml:"value,omitempty"`
}
