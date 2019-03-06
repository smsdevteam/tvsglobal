package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	cm "github.com/smsdevteam/tvsglobal/common"

	"github.com/gorilla/mux"
	//st "tvsglobal/tvsstructs"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Note Restful   ")
}

func main() {
	fmt.Println("Service Start...")
	//changepackage(108218909)
	//suspendsub(108218909)
	//restoresub(108218909)
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tgovsbn/ccbschangepackage/{customerid}", ccbschangepackage)
	mainRouter.HandleFunc("/tgovsbn/ccbssuspendsub/{customerid}", ccbssuspendsub)
	mainRouter.HandleFunc("/tgovsbn/ccbsrestoresub/{customerid}", ccbsrestoresub)

	log.Fatal(http.ListenAndServe(":8000", mainRouter))
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

func getft(w http.ResponseWriter, r *http.Request) {
	//var customerid int64
	//var err error
	params := mux.Vars(r)
	fmt.Println(params["ftid"])
	//customerid, err := strconv.ParseInt(params["ftid"], 10, 32)
	//customerid = 0
	//var noteResult st.FinancialTransaction

	//noteResult = GetNoteByNoteID(params["noteid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(changepackage(1))

}
