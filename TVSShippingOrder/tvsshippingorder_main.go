package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	so "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Shipping Order Restful")
}

func getOrderData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["soid"])
	var oSO so.ShippingOrderRes
	var soid int64
	soid, _ = strconv.ParseInt(params["soid"], 10, 64)
	oSO = GetShippingOrder(soid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oSO)
}

func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvsshippingorder", index)
	mainRouter.HandleFunc("/tvsshippingorder/getshippingorder/{soid}", getOrderData)

	log.Fatal(http.ListenAndServe(":8081", mainRouter))
}
