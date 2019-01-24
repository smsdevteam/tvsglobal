package tvsstructs

import (
	"encoding/xml"
	"time"
)

type Contact struct {
	ContactID               int64     `json:"contactid"`
	ActionDate              time.Time `json:"actiondate"`
	ActionTaken             string    `json:"actiontaken"`
	AllocatedToUser         int64     `json:"allocatedtouser"`
	Category                int64     `json:"category"`
	CreatedByUser           int64     `json:"createdbyuser"`
	CustomerID              int64     `json:"customerid"`
	CustomerProductID       int64     `json:"customerproductid"`
	Method                  string    `json:"method"`
	OrderID                 int64     `json:"orderid"`
	ProblemDesc             string    `json:"problemdesc"`
	ProductID               int64     `json:"productid"`
	StampDate               time.Time `json:"stampdate"`
	Status                  string    `json:"status"`
	WorkOrderID             int64     `json:"workorderid"`
	CreatedDate             time.Time `json:"createddate"`
	ExternalReferenceID     string    `json:"externalreferenceid"`
	DeviceID                int64     `json:"deviceid"`
	InvoiceID               int64     `json:"invoiceid"`
	LastUpdatedUserID       int64     `json:"lastupdateduserid"`
	ExternalReferenceID1    int64     `json:"externalreferenceid1"`
	ExternalReferenceID2    int64     `json:"externalreferenceid2"`
	ExternalReferenceID3    int64     `json:"externalreferenceid3"`
	ExternalReferenceID4    int64     `json:"externalreferenceid4"`
	ExternalReferenceID5    int64     `json:"externalreferenceid5"`
	ExternalReferenceValue1 string    `json:"externalreferencevalue1"`
	ExternalReferenceValue2 string    `json:"externalreferencevalue2"`
	ExternalReferenceValue3 string    `json:"externalreferencevalue3"`
	ExternalReferenceValue4 string    `json:"externalreferencevalue4"`
	ExternalReferenceValue5 string    `json:"externalreferencevalue5"`
}

type ListContact struct {
	Contacts []Contact `json:"contacts"`
}

type GetContactResponse struct {
	GetContactResult ContactXML `xml:"GetContactResult" json:"GetContactResult"`
}

type ContactXML struct {
	XMLName                 xml.Name  `xml:"GetContactResult"`
	ActionDate              time.Time `xml:"ActionDate" json:"ActionDate"`
	ActionTaken             string    `xml:"ActionTaken" json:"ActionTaken"`
	AllocatedToUser         int64     `xml:"AllocatedToUser" json:"AllocatedToUser"`
	CategoryKey             int64     `xml:"CategoryKey" json:"CategoryKey"`
	CreatedByUserKey        int64     `xml:"CreatedByUserKey" json:"CreatedByUserKey"`
	CreatedDate             time.Time `xml:"CreatedDate" json:"CreatedDate"`
	CustomerID              int64     `xml:"CustomerId" json:"CustomerId"`
	CustomerProductID       int64     `xml:"CustomerProductId" json:"CustomerProductId"`
	DeviceID                int64     `xml:"DeviceId" json:"DeviceId"`
	ExternalReferenceID     string    `xml:"ExternalReferenceId" json:"ExternalReferenceId"`
	ExternalReferenceID1    int64     `xml:"ExternalReferenceId1" json:"ExternalReferenceId1"`
	ExternalReferenceID2    int64     `xml:"ExternalReferenceId2" json:"ExternalReferenceId2"`
	ExternalReferenceID3    int64     `xml:"ExternalReferenceId3" json:"ExternalReferenceId3"`
	ExternalReferenceID4    int64     `xml:"ExternalReferenceId4" json:"ExternalReferenceId4"`
	ExternalReferenceID5    int64     `xml:"ExternalReferenceId5" json:"ExternalReferenceId5"`
	ExternalReferenceValue1 string    `xml:"ExternalReferenceValue1" json:"ExternalReferenceValue1"`
	ExternalReferenceValue2 string    `xml:"ExternalReferenceValue2" json:"ExternalReferenceValue2"`
	ExternalReferenceValue3 string    `xml:"ExternalReferenceValue3" json:"ExternalReferenceValue3"`
	ExternalReferenceValue4 string    `xml:"ExternalReferenceValue4" json:"ExternalReferenceValue4"`
	ExternalReferenceValue5 string    `xml:"ExternalReferenceValue5" json:"ExternalReferenceValue5"`
	ContactID               int64     `xml:"Id" json:"Id"`
	InvoiceID               int64     `xml:"InvoiceId" json:"InvoiceId"`
	LastUpdatedByUserID     int64     `xml:"LastUpdatedByUserId" json:"LastUpdatedByUserId"`
	MethodKey               string    `xml:"MethodKey" json:"MethodKey"`
	OrderID                 int64     `xml:"OrderId" json:"OrderId"`
	ProblemDescription      string    `xml:"ProblemDescription" json:"ProblemDescription"`
	ProductID               int64     `xml:"ProductId" json:"ProductId"`
	StampDate               time.Time `xml:"StampDate" json:"StampDate"`
	StatusKey               string    `xml:"StatusKey" json:"StatusKey"`
	WorkOrderID             int64     `xml:"WorkOrderId" json:"WorkOrderId"`
	Extended                string    `xml:"Extended" json:"Extended"`
}

type CreateContactRequest struct {
	InContact struct {
		ActionDate              time.Time `xml:"ActionDate" json:"ActionDate"`
		ActionTaken             string    `xml:"ActionTaken" json:"ActionTaken"`
		AllocatedToUser         int64     `xml:"AllocatedToUserKey" json:"AllocatedToUserKey"`
		CategoryKey             int64     `xml:"CategoryKey" json:"CategoryKey"`
		CreatedByUserKey        int64     `xml:"CreatedByUserKey" json:"CreatedByUserKey"`
		CreatedDate             time.Time `xml:"CreatedDate" json:"CreatedDate"`
		CustomerID              int64     `xml:"CustomerId" json:"CustomerId"`
		CustomerProductID       int64     `xml:"CustomerProductId" json:"CustomerProductId"`
		DeviceID                int64     `xml:"DeviceId" json:"DeviceId"`
		ExternalReferenceID     string    `xml:"ExternalReferenceId" json:"ExternalReferenceId"`
		ExternalReferenceID1    int64     `xml:"ExternalReferenceId1" json:"ExternalReferenceId1"`
		ExternalReferenceID2    int64     `xml:"ExternalReferenceId2" json:"ExternalReferenceId2"`
		ExternalReferenceID3    int64     `xml:"ExternalReferenceId3" json:"ExternalReferenceId3"`
		ExternalReferenceID4    int64     `xml:"ExternalReferenceId4" json:"ExternalReferenceId4"`
		ExternalReferenceID5    int64     `xml:"ExternalReferenceId5" json:"ExternalReferenceId5"`
		ExternalReferenceValue1 string    `xml:"ExternalReferenceValue1" json:"ExternalReferenceValue1"`
		ExternalReferenceValue2 string    `xml:"ExternalReferenceValue2" json:"ExternalReferenceValue2"`
		ExternalReferenceValue3 string    `xml:"ExternalReferenceValue3" json:"ExternalReferenceValue3"`
		ExternalReferenceValue4 string    `xml:"ExternalReferenceValue4" json:"ExternalReferenceValue4"`
		ExternalReferenceValue5 string    `xml:"ExternalReferenceValue5" json:"ExternalReferenceValue5"`
		ContactID               int64     `xml:"Id" json:"Id"`
		InvoiceID               int64     `xml:"InvoiceId" json:"InvoiceId"`
		LastUpdatedByUserID     int64     `xml:"LastUpdatedByUserId" json:"LastUpdatedByUserId"`
		MethodKey               string    `xml:"MethodKey" json:"MethodKey"`
		OrderID                 int64     `xml:"OrderId" json:"OrderId"`
		ProblemDescription      string    `xml:"ProblemDescription" json:"ProblemDescription"`
		ProductID               int64     `xml:"ProductId" json:"ProductId"`
		StampDate               time.Time `xml:"StampDate" json:"StampDate"`
		StatusKey               string    `xml:"StatusKey" json:"StatusKey"`
		WorkOrderID             int64     `xml:"WorkOrderId" json:"WorkOrderId"`
		Extended                string    `xml:"Extended" json:"Extended"`
	} `json:"inContact"`
	InReason int64 `json:"InReason"`
	ByUser   struct {
		ByUser    string `json:"byUser"`
		ByChannel string `json:"byChannel"`
		ByProject string `json:"byProject"`
		ByHost    string `json:"byHost"`
	} `json:"byUser"`
}

type CreateContactResponse struct {
	ErrorCode   int    `json:"errorcode"`
	ErrorDesc   string `json:"errordesc"`
	ResultValue string `json:"resultvalue"`
}
