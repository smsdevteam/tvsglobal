package tvsstructs

type CustomerInfo struct {
	Text           string `xml:",chardata"`
	BirthDate      string `xml:"BirthDate"`
	BusinessUnitID string `xml:"BusinessUnitId"`
	ClassID        string `xml:"ClassId"`
	CustomFields   struct {
		Text             string `xml:",chardata"`
		CustomFieldValue []struct {
			Text     string `xml:",chardata"`
			Extended struct {
				Text string `xml:",chardata"`
				Nil  string `xml:"nil,attr"`
			} `xml:"Extended"`
			ID       string `xml:"Id"`
			Name     string `xml:"Name"`
			Sequence string `xml:"Sequence"`
			Value    string `xml:"Value"`
		} `xml:"CustomFieldValue"`
	} `xml:"CustomFields"`
	CustomerSince  string `xml:"CustomerSince"`
	DefaultAddress struct {
		Text                    string `xml:",chardata"`
		BigCity                 string `xml:"BigCity"`
		CareOfName              string `xml:"CareOfName"`
		CountryKey              string `xml:"CountryKey"`
		CustomerCaptureCategory string `xml:"CustomerCaptureCategory"`
		Directions              string `xml:"Directions"`
		Email                   string `xml:"Email"`
		Extended                struct {
			Text string `xml:",chardata"`
			Nil  string `xml:"nil,attr"`
		} `xml:"Extended"`
		Extra              string `xml:"Extra"`
		Extra1             string `xml:"Extra1"`
		Extra2             string `xml:"Extra2"`
		Extra3             string `xml:"Extra3"`
		Extra4             string `xml:"Extra4"`
		Extra5             string `xml:"Extra5"`
		ExtraExtra         string `xml:"ExtraExtra"`
		Fax1               string `xml:"Fax1"`
		Fax2               string `xml:"Fax2"`
		FirstName          string `xml:"FirstName"`
		FlatOrApartment    string `xml:"FlatOrApartment"`
		GeoCodeID          string `xml:"GeoCodeId"`
		HomePhone          string `xml:"HomePhone"`
		HomePhoneExt       string `xml:"HomePhoneExt"`
		HouseNumberAlpha   string `xml:"HouseNumberAlpha"`
		HouseNumberNumeric string `xml:"HouseNumberNumeric"`
		ID                 string `xml:"Id"`
		LandMark           string `xml:"LandMark"`
		MarketSegmentID    string `xml:"MarketSegmentId"`
		MiddleName         string `xml:"MiddleName"`
		PostalCode         string `xml:"PostalCode"`
		ProvinceKey        string `xml:"ProvinceKey"`
		SmallCity          string `xml:"SmallCity"`
		Street             string `xml:"Street"`
		Surname            string `xml:"Surname"`
		TitleKey           string `xml:"TitleKey"`
		ValidAddressID     string `xml:"ValidAddressId"`
		WorkPhone          string `xml:"WorkPhone"`
		WorkPhoneExt       string `xml:"WorkPhoneExt"`
	} `xml:"DefaultAddress"`
	EmailNotifyOptionKey  string `xml:"EmailNotifyOptionKey"`
	ExemptionCodeKey      string `xml:"ExemptionCodeKey"`
	ExemptionFrom         string `xml:"ExemptionFrom"`
	ExemptionSerialNumber string `xml:"ExemptionSerialNumber"`
	Extended              struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"Extended"`
	FiscalCode               string `xml:"FiscalCode"`
	FiscalNumber             string `xml:"FiscalNumber"`
	ID                       string `xml:"Id"`
	InternetPassword         string `xml:"InternetPassword"`
	InternetUserID           string `xml:"InternetUserId"`
	IsDistributor            string `xml:"IsDistributor"`
	IsHeadend                string `xml:"IsHeadend"`
	IsProductProvider        string `xml:"IsProductProvider"`
	IsServiceProvider        string `xml:"IsServiceProvider"`
	IsStockHandler           string `xml:"IsStockHandler"`
	LanguageKey              string `xml:"LanguageKey"`
	Magazines                string `xml:"Magazines"`
	ParentID                 string `xml:"ParentId"`
	PreferredContactMethodID string `xml:"PreferredContactMethodId"`
	ReferenceNumber          string `xml:"ReferenceNumber"`
	ReferenceTypeKey         string `xml:"ReferenceTypeKey"`
	SegmentationKey          string `xml:"SegmentationKey"`
	StatusKey                string `xml:"StatusKey"`
	TypeKey                  string `xml:"TypeKey"`
}

/* type Customerrespon  struct {
   CustomerInfoobj  CustomerInfo
   ResponResult  ResponseResult
} */
// func Area(len, wid float64) float64 {
// 	area := len * wid
// 	return area
// }
