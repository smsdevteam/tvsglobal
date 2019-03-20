package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	st "github.com/smsdevteam/tvsglobal/TVSStructs"
)

func index(w http.ResponseWriter, r *http.Request) {
	var req st.TVSSubmitOrdReqData
	fmt.Println(w, "Welcome to TVS Device Restful")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvssubmitorder/", index)
	mainRouter.HandleFunc("/tvssubmitorder/submitorder/", submitorder).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", mainRouter))

}

func submitorder(w http.ResponseWriter, r *http.Request) {

	fmt.Println("start call submitorder")
	fmt.Println("************************************************************************")
	temp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	//Read Json Request
	var req st.TVSSubmitOrdReqData
	err = json.Unmarshal(temp, &req)
	if err != nil {
		fmt.Println("There was an error:", err)
		panic(err)
	}

	oRes := tvssubmitorder(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}
