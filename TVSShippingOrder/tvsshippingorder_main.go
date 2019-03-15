package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	cm "github.com/smsdevteam/tvsglobal/Common"     // db
	st "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Shipping Order Restful")
}

// Warehouse Shipping Order
func cancelWHOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//p(params["soid"])
	var oSO st.ResponseResult
	var soid, reason int64
	soid = cm.StrToInt64(params["soid"])
	reason = cm.StrToInt64(params["reason"])
	oSO = CancelShippingOrder(soid, reason, params["by"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oSO)
}

func createWHOrder(w http.ResponseWriter, r *http.Request) {
	temp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	//Read Json Request
	var req st.ShippingOrderReq
	err = json.Unmarshal(temp, &req)
	if err != nil {
		p("There was an error:", err)
		p(err)
		panic(err)
	}
	p(req)

	var oRes st.SOResult

	oRes = CreateWHOrder(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}

func getWHOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var oSO st.ShippingOrderRes

	oSO = GetWHOrder(cm.StrToInt64(params["soid"]))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oSO)
}

func shipWHOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var oSO st.ResponseResult
	var soid, reason, reasondv int64
	soid = cm.StrToInt64(params["soid"])
	reason = cm.StrToInt64(params["reason"])
	reasondv = cm.StrToInt64(params["reasondv"])

	oSO = ShipWHOrder(soid, reason, reasondv, params["by"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oSO)
}

// Customer Shipping Order
func cancelShippingOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//p(params["soid"])
	var oSO st.ResponseResult
	var soid, reason int64
	soid = cm.StrToInt64(params["soid"])
	reason = cm.StrToInt64(params["reason"])
	oSO = CancelShippingOrder(soid, reason, params["by"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oSO)
}

func createShippingOrder(w http.ResponseWriter, r *http.Request) {
	temp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}


	//Read Json Request
	var req st.ShippingOrderDataReq
	err = json.Unmarshal(temp, &req)
	if err != nil {
		p("There was an error:", err)
		panic(err)
	}
	p(req)

	var oRes st.SOResult

	oRes = CreateShippingOrder(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}

func getShippingOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//p(params["soid"])
	var oSO st.SOResult

	oSO = GetShippingOrder(cm.StrToInt64(params["soid"]), params["by"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oSO)
}

func shipShippingOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var oSO st.ResponseResult
	var soid, reason, reasondv int64
	soid = cm.StrToInt64(params["soid"])
	reason = cm.StrToInt64(params["reason"])
	reasondv = cm.StrToInt64(params["reasondv"])

	oSO = ShipWHOrder(soid, reason, reasondv, params["by"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oSO)
}

func main() {
	p("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvsshippingorder", index)

	mainRouter.HandleFunc("/tvsshippingorder/cancelwhshippingorder/{soid}/{reason}/{by}", cancelWHOrder)
	mainRouter.HandleFunc("/tvsshippingorder/createwhshippingorder", createWHOrder).Methods("POST")
	mainRouter.HandleFunc("/tvsshippingorder/getwhshippingorder/{soid}", getWHOrder)
	mainRouter.HandleFunc("/tvsshippingorder/shipwhshippingorder/{soid}/{reason}/{reasondv}/{by}", shipWHOrder)

	mainRouter.HandleFunc("/tvsshippingorder/cancelshippingorder/{soid}/{reason}/{by}", cancelShippingOrder)
	mainRouter.HandleFunc("/tvsshippingorder/createshippingorder", createShippingOrder).Methods("POST")
	mainRouter.HandleFunc("/tvsshippingorder/getshippingorder/{soid}", getShippingOrder)
	mainRouter.HandleFunc("/tvsshippingorder/getshippingorder/{soid}/{by}", getShippingOrder)
	mainRouter.HandleFunc("/tvsshippingorder/shipshippingorder/{soid}/{reason}/{reasondv}/{by}", shipShippingOrder)

	log.Fatal(http.ListenAndServe(":8081", mainRouter))
}
