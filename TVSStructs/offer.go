package tvsstructs

import (
	"time"
)

type Offer struct {
	Active                string    `json:"Active"`
	AgreementDetailID     int64     `json:"AgreementDetailId"`
	AgreementID           int64     `json:"AgreementId"`
	ApplyToLevel          string    `json:"ApplyToLevel"`
	CustomerID            int64     `json:"CustomerId"`
	EndDate               time.Time `json:"EndDate"`
	FinancialAccountID    int64     `json:"FinancialAccountId"`
	ID                    int64     `json:"Id"`
	OfferDefinitionID     int64     `json:"OfferDefinitionId"`
	SandboxID             int64     `json:"SandboxId"`
	SandboxSkipValidation string    `json:"SandboxSkipValidation"`
	StartDate             time.Time `json:"StartDate"`
}

type GetOfferResponse struct {
	GetOfferResult Offer  `xml:"GetOfferResult" json:"GetOfferResult"`
	ErrorCode      int    `json:"errorcode"`
	ErrorDesc      string `json:"errordesc"`
}