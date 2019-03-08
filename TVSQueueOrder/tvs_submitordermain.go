package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	// util
	st "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Device Restful")
}
func createNewSN(w http.ResponseWriter, r *http.Request) {
	temp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	//Read Json Request
	var req st.NewDeviceReq
	err = json.Unmarshal(temp, &req)
	if err != nil {
		fmt.Println("There was an error:", err)
		panic(err)
	}

	var oRes st.NewDeviceRes

	oRes = CreateNewSerialNumber(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}

func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvsqueueorder", index)
	mainRouter.HandleFunc("/tvsqueueorder/createnewserialnumber", createNewSN).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", mainRouter))

}
