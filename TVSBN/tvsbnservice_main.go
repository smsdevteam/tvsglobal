package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	cm "github.com/smsdevteam/tvsglobal/common"
	"github.com/gorilla/mux"
	st "github.com/smsdevteam/tvsglobal/TVSStructs"
)
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Note Restful   ")
}
func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tgovsbn/ccbschangepackage/{customerid}", ccbschangepackage)
	mainRouter.HandleFunc("/tgovsbn/ccbschangepackagep", ccbschangepackagep).Methods("POST")
	mainRouter.HandleFunc("/tgovsbn/ccbssuspendsub/{customerid}", ccbssuspendsub)
	mainRouter.HandleFunc("/tgovsbn/ccbsrestoresub/{customerid}", ccbsrestoresub)
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}
func ccbschangepackagep(w http.ResponseWriter, r *http.Request) {
	temp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	//Read Json Request
	var req st.CustomerInfo
	err = json.Unmarshal(temp, &req)
	if err != nil {
		fmt.Println("There was an error:", err)
		panic(err)
	}
	var oRes string
	oRes = changepackage(cm.StrToInt(req.ID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}

func ccbschangepackage(w http.ResponseWriter, r *http.Request) {
	//var customerid int64
	//var err error
	params := mux.Vars(r)
	//fmt.Println("start change package  ")
	fmt.Println("start change package  " + params["customerid"])
	customerid := cm.StrToInt(params["customerid"])
	//changepackage(customerid)
	//customerid = 0
	//var noteResult st.FinancialTransaction

	//noteResult = GetNoteByNoteID(params["noteid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(changepackage(customerid))

}
func ccbssuspendsub(w http.ResponseWriter, r *http.Request) {
	//var customerid int64
	//var err error
	params := mux.Vars(r)
	fmt.Println(params["customerid"])
	customerid := cm.StrToInt(params["customerid"])
	//changepackage(customerid)
	//customerid = 0
	//var noteResult st.FinancialTransaction

	//noteResult = GetNoteByNoteID(params["noteid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suspendsub(customerid))

}
func ccbsrestoresub(w http.ResponseWriter, r *http.Request) {
	//var customerid int64
	//var err error
	params := mux.Vars(r)
	fmt.Println(params["customerid"])
	customerid := cm.StrToInt(params["customerid"])
	//changepackage(customerid)
	//customerid = 0
	//var noteResult st.FinancialTransaction

	//noteResult = GetNoteByNoteID(params["noteid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restoresub(customerid))

}

