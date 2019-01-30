package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	st "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
	// db
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Shipping Order Restful")
}

func getOrderData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["soid"])
	var oSO st.ShippingOrderRes
	var soid int64
	soid, _ = strconv.ParseInt(params["soid"], 10, 64)
	oSO = GetShippingOrder(soid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oSO)
}

func cancelOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["soid"])
	var oSO st.ResponseResult
	var soid, reason int64
	soid, _ = strconv.ParseInt(params["soid"], 10, 64)
	reason, _ = strconv.ParseInt(params["reason"], 10, 64)
	oSO = CancelShippingOrder(soid, reason, params["by"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oSO)
}

func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvsshippingorder", index)
	mainRouter.HandleFunc("/tvsshippingorder/getshippingorder/{soid}", getOrderData)
	mainRouter.HandleFunc("/tvsshippingorder/cancelshippingorder/{soid}/{reason}/{by}", cancelOrder)

	log.Fatal(http.ListenAndServe(":8081", mainRouter))
}
