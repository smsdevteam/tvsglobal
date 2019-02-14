package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"

	//"log"
	"net/http"

	"github.com/gorilla/mux"
	//st "tvsglobal/tvsstructs"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Note Restful")
}

func main() {
	fmt.Println("Service Start...")
	changepackage(108218909)
	//mainRouter := mux.NewRouter().StrictSlash(true)
	//mainRouter.HandleFunc("/tgovsbn/getft/{ftid}", getft)
	//log.Fatal(http.ListenAndServe(":8000", mainRouter))
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
