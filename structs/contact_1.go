package structs

import "time"

//Contact is Contact
type Contact struct {
	ContactID               int
	ActionDate              time.Time
	ActionTaken             string
	AllocatedToUser         int
	Category                int
	CreatedByUser           int
	CustomerID              int
	CustomerProductID       int
	Method                  string
	OrderID                 int
	ProblemDesc             string
	ProductID               int
	StampDate               time.Time
	Status                  string
	WorkOrderID             int
	CreatedDate             time.Time
	ExternalReferenceID     string
	DeviceID                int
	InvoiceID               int
	LastUpdatedUserID       int
	ExternalReferenceID1    int
	ExternalReferenceID2    int
	ExternalReferenceID3    int
	ExternalReferenceID4    int
	ExternalReferenceID5    int
	ExternalReferenceValue1 string
	ExternalReferenceValue2 string
	ExternalReferenceValue3 string
	ExternalReferenceValue4 string
	ExternalReferenceValue5 string
}
