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

//DeleteKeywordRequest obj
type DeleteKeywordRequest struct {
	ByUser struct {
		ByChannel string `json:"bychannel"`
		ByHost    string `json:"byhost"`
		ByProject string `json:"byproject"`
		ByUser    string `json:"byuser"`
	} `json:"byuser"`
	InKeywordID int64 `json:"inkeywordid"`
	InReason    int64 `json:"inreason"`
}

//DeleteKeywordResponse obj
type DeleteKeywordResponse struct {
	ResultValue string `json:"resultvalue"`
	ErrorCode   int    `json:"errorcode"`
	ErrorDesc   string `json:"errordesc"`
}

// NewDeleteKeywordResponse Obj
func NewDeleteKeywordResponse() *DeleteKeywordResponse {
	return &DeleteKeywordResponse{
		ErrorCode: 1,
		ErrorDesc: "Unexpected Error",
	}
}

//UpdateKeywordRequest obj
type UpdateKeywordRequest struct {
	ByUser struct {
		ByChannel string `json:"byChannel"`
		ByHost    string `json:"byHost"`
		ByProject string `json:"byProject"`
		ByUser    string `json:"byUser"`
	} `json:"byUser"`
	InKeyword struct {
		Attribute     string    `json:"Attribute"`
		Count         int64     `json:"Count"`
		CustomerID    int64     `json:"CustomerId"`
		Date          time.Time `json:"Date"`
		Extended      string    `json:"Extended"`
		ID            int64     `json:"Id"`
		KeywordID     int64     `json:"KeywordId"`
		KeywordTypeID int64     `json:"KeywordTypeId"`
		Name          string    `json:"Name"`
	} `json:"inKeyword"`
	InReason int64 `json:"inReason"`
}

//UpdateKeywordResponse obj
type UpdateKeywordResponse struct {
	ResultValue string `json:"resultvalue"`
	ErrorCode   int    `json:"errorcode"`
	ErrorDesc   string `json:"errordesc"`
}

// NewUpdateKeywordResponse Obj
func NewUpdateKeywordResponse() *UpdateKeywordResponse {
	return &UpdateKeywordResponse{
		ErrorCode: 1,
		ErrorDesc: "Unexpected Error",
	}
}
