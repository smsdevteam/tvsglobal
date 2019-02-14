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
