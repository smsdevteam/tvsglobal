package main

// นำเข้า package fmt มาใช้งาน
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	c "github.com/smsdevteam/tvsglobal/tvsstructs" // referpath
)

type TVS_Customer_request struct {
	customer_obj c.Customerinfo
	Orderno      string
}
type TVS_Customer_response struct {
	customer_obj c.Customerinfo
	Orderno      string
	Resultcode   string
}

func getcustomerinfo(tvscustreq TVS_Customer_request) TVS_Customer_response {
	resulttvsresponse := TVS_Customer_response{}
	resulttvsresponse.Orderno = tvscustreq.Orderno
	return resulttvsresponse
}
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Customer Restful")
}
func main() {
 
//	fmt.Println("start service...")
  // 	mainRouter.HandleFunc("/tvscustomer", index)
//	mainRouter.HandleFunc("/tvscustomer/getcustomerinfo/{customerid}", getCustomer)
//	mainRouter.HandleFunc("/tvscustomer/getlistcustomer/{customerid}", getListcustomer)
//	mainRouter.HandleFunc("/tvsnote/createnote", createNote).Methods("POST")
//	mainRouter.HandleFunc("/tvsnote/updatenote", updateNote).Methods("POST")

//	log.Fatal(http.ListenAndServe(":8080", router))
	//	tvscustreq := TVS_Customer_request{}
	//	tvscustreq.Orderno="100"
	// tvscustreq.customer_obj.
	//	Res:= add(tvscustreq)
	//	fmt.Println(Res)
	//	var Customer_requestobj c.Customerinfo[]
	//	Customer_requestobj.customerno   = 1
	//	Customer_requestobj.BirthDate = "06/10/23"
	//	Customer_requestobj.EmailNotifyOptionKey = "nattachai.son@truecorp.co.th"
	//	fmt.Println(Customer_requestobj)

}

func getCustomer(w http.ResponseWriter, r *http.Request) {
params := mux.Vars(r)
	 
	var cusotmerResult  c.Customer

	customerResult = GetCustomerByCustomerID(params["customerid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cusotmerResult)
}
/* func addTVS(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	//Read Json Request
	var req TVS_Customer_request
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	log.Println("get tvs")
	log.Println(req)

	//call recon api
	var res TVS_Customer_response
	// assign orderno
	res.Orderno = req.Orderno

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
} */
