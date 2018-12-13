package main

// นำเข้า package fmt มาใช้งาน
import (
 c	"github.com/smsdevteam/tvsglobal/tvsstructs" // referpath
	"fmt"
)

type TVS_Customer_request struct {
	customer_obj    c.Customerinfo
	Orderno          string 
	
}
type TVS_Customer_response struct {
	customer_obj    c.Customerinfo 
	Orderno          string 
	Resultcode       string
}
func add(tvscustreq TVS_Customer_request) TVS_Customer_response  {
	  resulttvsresponse := TVS_Customer_response{}
	  resulttvsresponse.Orderno =tvscustreq.Orderno
	return resulttvsresponse
}
func main() {
	tvscustreq := TVS_Customer_request{}
	tvscustreq.Orderno="100"
	//tvscustreq.customer_obj.
	Res:= add(tvscustreq)
	fmt.Println(Res)
	//	var Customer_requestobj c.Customerinfo[]
//	Customer_requestobj.customerno   = 1
//	Customer_requestobj.BirthDate = "06/10/23"
//	Customer_requestobj.EmailNotifyOptionKey = "nattachai.son@truecorp.co.th"
//	fmt.Println(Customer_requestobj)
 

}
