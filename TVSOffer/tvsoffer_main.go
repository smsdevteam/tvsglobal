package main

import (
	"encoding/json"
	"fmt"
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
	/* mainRouter.HandleFunc("/tvsnote/getnote/{noteid}", getNote)
	mainRouter.HandleFunc("/tvsnote/getlistnote/{customerid}", getListNote)
	mainRouter.HandleFunc("/tvsnote/createnote", createNote).Methods("POST")
	mainRouter.HandleFunc("/tvsnote/updatenote", updateNote).Methods("POST") */
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}

func getOffer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["offerid"])

	var offerResult st.GetOfferResponse
	offerResult = GetOfferByOfferID(params["offerid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offerResult)
}
