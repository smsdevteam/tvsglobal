package tvsstructs

import (
	"time"
)

//Keyword obj
type Keyword struct {
	AllowedKeywordID int64     `json:"allowed_keyword_id"`
	Attribute        string    `json:"attribute"`
	CountValue       int64     `json:"count_value"`
	CustomerID       int64     `json:"customer_id"`
	DateValue        time.Time `json:"date_value"`
	ID               int64     `json:"id"`
	KAAttribute      string    `json:"kaattribute"`
	KAKeyword        string    `json:"kakeyword"`
	KALongDescr      string    `json:"kalongdescr"`
	KeyTypesID       int64     `json:"keytypes_id"`
	KTName           string    `json:"ktname"`
	KTUserKey        string    `json:"ktuserkey"`
}

//GetKeywordResult obj
type GetKeywordResult struct {
	MyKeyword Keyword `json:"keyword"`
	ErrorCode int     `json:"errorcode"`
	ErrorDesc string  `json:"errordesc"`
}

// NewGetKeywordResult Obj
func NewGetKeywordResult() *GetKeywordResult {
	return &GetKeywordResult{
		ErrorCode: 1,
		ErrorDesc: "Unexpected Error",
	}
}

//ListKeyword obj
type ListKeyword struct {
	Keywords []Keyword `json:"keywords"`
}

//GetListKeywordResult obj
type GetListKeywordResult struct {
	MyListKeyword ListKeyword `json:"getlistkeywordresult"`
	ErrorCode     int         `json:"errorcode"`
	ErrorDesc     string      `json:"errordesc"`
}

// NewGetListKeywordResult Obj
func NewGetListKeywordResult() *GetListKeywordResult {
	return &GetListKeywordResult{
		ErrorCode: 1,
		ErrorDesc: "Unexpected Error",
	}
}

//CreateKeywordRequest obj
type CreateKeywordRequest struct {
	ByUser struct {
		ByChannel string `json:"bychannel"`
		ByHost    string `json:"byhost"`
		ByProject string `json:"byproject"`
		ByUser    string `json:"byuser"`
	} `json:"byuser"`
	InKeyword struct {
		Attribute     string    `json:"attribute"`
		Count         int64     `json:"count"`
		CustomerID    int64     `json:"customerid"`
		Date          time.Time `json:"date"`
		Extended      string    `json:"extended"`
		ID            int64     `json:"id"`
		KeywordID     int64     `json:"keywordid"`
		KeywordTypeID int64     `json:"keywordtypeid"`
		Name          string    `json:"name"`
	} `json:"inkeyword"`
	InReason int64 `json:"inreason"`
}

//CreateKeywordResponse obj
type CreateKeywordResponse struct {
	ErrorCode   int    `json:"errorcode"`
	ErrorDesc   string `json:"errordesc"`
	ResultValue string `json:"resultvalue"`
}

// NewCreateKeywordResponse Obj
func NewCreateKeywordResponse() *CreateKeywordResponse {
	return &CreateKeywordResponse{
		ErrorCode: 1,
		ErrorDesc: "Unexpected Error",
	}
}
