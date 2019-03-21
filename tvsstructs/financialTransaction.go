package tvsstructs
 
type FinancialTransaction struct {
				Text               string `xml:",chardata"`
				AmountExcludingTax string `xml:"AmountExcludingTax"`
				AmountIncludingTax string `xml:"AmountIncludingTax"`
				AppearedOnInvoice  string `xml:"AppearedOnInvoice"`
				BankDate           string `xml:"BankDate"`
				BaseAmount         string `xml:"BaseAmount"`
				BusinessUnitID     string `xml:"BusinessUnitId"`
				Comment            string `xml:"Comment"`
				CreateDate         string `xml:"CreateDate"`
				CreatedByEvent     string `xml:"CreatedByEvent"`
				CreatedByUserID    string `xml:"CreatedByUserId"`
				CurrencyID         string `xml:"CurrencyId"`
				CustomerID         string `xml:"CustomerId"`
				DebitOrCredit      string `xml:"DebitOrCredit"`
				DeviceID           string `xml:"DeviceId"`
				EntityID           string `xml:"EntityId"`
				EntityType         string `xml:"EntityType"`
				        
			}  
		 
