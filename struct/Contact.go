package main

import (
	"fmt"
	"time"
)

type Contact struct {
	contactID               int
	actionDate              time.Time
	actionTaken             string
	allocatedToUser         int
	category                int
	createdByUser           int
	customerID              int
	customerProductID       int
	method                  string
	orderID                 int
	problemDesc             string
	productID               int
	stampDate               time.Time
	status                  string
	workOrderID             int
	createdDate             time.Time
	externalReferenceID     string
	deviceID                int
	invoiceID               int
	lastUpdatedUserID       int
	externalReferenceID1    int
	externalReferenceID2    int
	externalReferenceID3    int
	externalReferenceID4    int
	externalReferenceID5    int
	externalReferenceValue1 string
	externalReferenceValue2 string
	externalReferenceValue3 string
	externalReferenceValue4 string
	externalReferenceValue5 string
}

func main() {
	fmt.Println("test")
}
