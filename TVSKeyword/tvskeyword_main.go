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
	fmt.Fprintf(w, "Welcome to TVS Keyword Restful")
}

func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvskeyword", index)
	mainRouter.HandleFunc("/tvskeyword/getkeyword/{keywordid}", getKeyword)
	mainRouter.HandleFunc("/tvskeyword/getlistkeyword/{customerid}", getListKeyword)
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}

func getKeyword(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var keywordResult *st.GetKeywordResult

	keywordResult = GetKeywordByKeywordID(params["keywordid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keywordResult)
}

func getListKeyword(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var listKeywordResult *st.GetListKeywordResult

	listKeywordResult = GetListKeywordByCustomerID(params["customerid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listKeywordResult)
}
