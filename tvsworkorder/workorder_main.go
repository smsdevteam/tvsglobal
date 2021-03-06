package main

// นำเข้า package fmt มาใช้งาน
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	c "github.com/smsdevteam/tvsglobal/tvsstructs" // referpath
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS WorkOrder Restful")
}
func main() {

	fmt.Println("start service Workorder...")
 
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/TVSWorkOrder", index)
	mainRouter.HandleFunc("/TVSWorkOrder/getworkorderbyworkorderid/{workorderid}", getworkorderbyworkorderid)
	mainRouter.HandleFunc("/TVSWorkOrder/Getworkorderservicebyid/{workorderid}", getworkorderservicebyid)


	
	//mainRouter.HandleFunc("/tvsnote/getcustomerinfo", getCustomer).Methods("POST")
	//GetWorkorderByCustomerID("111")
	//mainRouter := mux.NewRouter().StrictSlash(true)
	//mainRouter.HandleFunc("/tvsworkorder", index)
	//mainRouter.HandleFunc("/tvsworkorder/getworkorder/{customerid}", getWorkorder)
	log.Fatal(http.ListenAndServe(":8080", mainRouter))
}

func getWorkorder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var workorderResult c.WorkorderInfo

	workorderResult = GetWorkorderByCustomerID(params["customerid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workorderResult)
}
func getworkorderbyworkorderid(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	var workorderResult c.WorkorderInfo

	workorderResult = GetWorkorderByworkorderid(params["workorderid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workorderResult)
}
func getworkorderservicebyid(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	var workorderserviceResult []c.WorkorderServiceDTInfo

	workorderserviceResult = Getworkorderservicebyid(params["workorderid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workorderserviceResult)
}