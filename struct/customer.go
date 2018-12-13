package customer_struct

import (
	"fmt"
)

// นำเข้า package fmt มาใช้งาน
 
  

type Customerinfo struct {
	Customerno              int
   BirthDate                string
  BusinessUnitId            int
  ClassId                    int
  CustomerSince         string
  EmailNotifyOptionKey string
  ExemptionCodeKey string
  ExemptionFrom string
  ExemptionSerialNumber string
  Extended string
  FiscalCode string
  FiscalNumber string
  InternetPassword string
  InternetUserId string
  IsDistributor string
  IsHeadend string
  IsProductProvider string
  IsServiceProvider string
  IsStockHandler string
  LanguageKey string
  Magazines string
  ParentId int
  PreferredContactMethodId int
  ReferenceNumber string
  ReferenceTypeKey string
  SegmentationKey string
  StatusKey string
  TypeKey string
  RefNo string
  ApplyID  string
  CustomerNo string
  NoOfPoint string
  BusinessRole string
  BirthDateEnc string
  BirthDateIv string
  BirthDatehash string
  InternetUserName string
  InternetPasswordEnc string
  InternetPasswordIv string
  ReferrnceNumberEnc string
  ReferrnceNumberIv string
  ReferrnceNumberHash string
  BirthDateKey string
  ENTITY_ID 
	 
}

func init() {
	fmt.Println("customer struct package initialized")

}
// func Area(len, wid float64) float64 {
// 	area := len * wid
// 	return area
// }
