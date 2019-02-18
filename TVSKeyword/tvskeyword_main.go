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
	fmt.Fprintf(w, "Welcome to TVS Keyword Restful")
}

func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvskeyword", index)
	mainRouter.HandleFunc("/tvskeyword/getkeyword/{keywordid}", getKeyword)
	mainRouter.HandleFunc("/tvskeyword/getlistkeyword/{customerid}", getListKeyword)
	mainRouter.HandleFunc("/tvskeyword/createkeyword", createKeyword).Methods("POST")
	mainRouter.HandleFunc("/tvskeyword/deletekeyword", deleteKeyword).Methods("POST")
	mainRouter.HandleFunc("/tvskeyword/updatekeyword", updateKeyword).Methods("POST")
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

func createKeyword(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var req st.CreateKeywordRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	var res *st.CreateKeywordResponse
	res = CreateKeyword(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func deleteKeyword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var req st.DeleteKeywordRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	var res *st.DeleteKeywordResponse
	res = DeleteKeyword(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func updateKeyword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var req st.UpdateKeywordRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	var res *st.UpdateKeywordResponse
	res = UpdateKeyword(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
