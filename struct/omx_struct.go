package omx_struct

import "encoding/xml"

type SubmitOrderOpRequest struct {
	XMLName xml.Name `xml:"submitOrderRequest"`
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
				Offers struct {
					Text               string `xml:",chardata"`
					Action             string `xml:"action"`
					EffectiveDate      string `xml:"effectiveDate"`
					OfferName          string `xml:"offerName"`
					TargetPayChannelId string `xml:"targetPayChannelId"`
					ServiceType        string `xml:"serviceType"`
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
				SubscriberIndyName struct {
					Text               string `xml:",chardata"`
					Identification     string `xml:"identification"`
					IdentificationType string `xml:"identificationType"`
					Language           string `xml:"language"`
					Title              string `xml:"title"`
					FirstName          string `xml:"firstName"`
					LastName           string `xml:"lastName"`
					Gender             string `xml:"gender"`
				} `xml:"subscriberIndyName"`
				SubscriberNumber string `xml:"subscriberNumber"`
				SubscriberType   string `xml:"subscriberType"`
				ExtendedInfo     []struct {
					Text  string `xml:",chardata"`
					Name  string `xml:"name"`
					Value string `xml:"value"`
				} `xml:"ExtendedInfo"`
			} `xml:"subscriber"`
		} `xml:"OU"`
	} `xml:"Customer"`
}
