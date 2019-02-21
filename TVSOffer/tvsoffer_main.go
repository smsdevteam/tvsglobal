package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	st "github.com/smsdevteam/tvsglobal/tvsstructs"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Offer Restful")
}

func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvsoffer", index)
	mainRouter.HandleFunc("/tvsoffer/getoffer/{offerid}", getOffer)
	mainRouter.HandleFunc("/tvsoffer/getlistoffer/{customerid}", getListOffer)
	mainRouter.HandleFunc("/tvsoffer/createoffer", createOffer).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}

func getOffer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//fmt.Println(params["offerid"])

	var offerResult *st.GetOfferResponse
	offerResult = GetOfferByOfferID(params["offerid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offerResult)
}

func getListOffer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//fmt.Println(params["customerid"])

	var listOfferResult *st.GetListOfferResult

	listOfferResult = GetListOfferByCustomerID(params["customerid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listOfferResult)
}

func createOffer(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var req st.CreateOfferRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	//log.Println(req)

	var res *st.CreateOfferResponse
	res = CreateOffer(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
