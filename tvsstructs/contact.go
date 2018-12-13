package tvsstructs

import (
	"time"
)

type Contact struct {
	ContactID               int       `json:"contactid"`
	ActionDate              time.Time `json:"actiondate"`
	ActionTaken             string    `json:"actiontaken"`
	AllocatedToUser         int       `json:"allocatedtouser"`
	Category                int       `json:"category"`
	CreatedByUser           int       `json:"createdbyuser"`
	CustomerID              int       `json:"customerid"`
	CustomerProductID       int       `json:"customerproductid"`
	Method                  string    `json:"method"`
	OrderID                 int       `json:"orderid"`
	ProblemDesc             string    `json:"problemdesc"`
	ProductID               int       `json:"productid"`
	StampDate               time.Time `json:"stampdate"`
	Status                  string    `json:"status"`
	WorkOrderID             int       `json:"workorderid"`
	CreatedDate             time.Time `json:"createddate"`
	ExternalReferenceID     string    `json:"externalreferenceid"`
	DeviceID                int       `json:"deviceid"`
	InvoiceID               int       `json:"invoiceid"`
	LastUpdatedUserID       int       `json:"lastupdateduserid"`
	ExternalReferenceID1    int       `json:"externalreferenceid1"`
	ExternalReferenceID2    int       `json:"externalreferenceid2"`
	ExternalReferenceID3    int       `json:"externalreferenceid3"`
	ExternalReferenceID4    int       `json:"externalreferenceid4"`
	ExternalReferenceID5    int       `json:"externalreferenceid5"`
	ExternalReferenceValue1 string    `json:"externalreferencevalue1"`
	ExternalReferenceValue2 string    `json:"externalreferencevalue2"`
	ExternalReferenceValue3 string    `json:"externalreferencevalue3"`
	ExternalReferenceValue4 string    `json:"externalreferencevalue4"`
	ExternalReferenceValue5 string    `json:"externalreferencevalue5"`
}
