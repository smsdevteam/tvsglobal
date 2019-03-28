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
	fmt.Fprintf(w, "Welcome to TVS Note Restful")
}

func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvsnote", index)
	mainRouter.HandleFunc("/tvsnote/getnote/{noteid}", getnote)
	mainRouter.HandleFunc("/tvsnote/getlistnote/{customerid}", getlistnote)
	mainRouter.HandleFunc("/tvsnote/createnote", createnote).Methods("POST")
	mainRouter.HandleFunc("/tvsnote/updatenote", updatenote).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}

func getnote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//fmt.Println(params["noteid"])

	var noteResult *st.GetNoteResult

	noteResult = GetNoteByNoteID(params["noteid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(noteResult)
}

func getlistnote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//fmt.Println(params["noteid"])

	var listNoteResult *st.GetListNoteResult

	listNoteResult = GetListNoteByCustomerID(params["customerid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listNoteResult)
}

func createnote(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var req st.CreateNoteRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	//log.Println(req)

	var res *st.CreateNoteResponse
	res = CreateNote(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func updatenote(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var req st.UpdateNoteRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	//log.Println(req)

	var res st.UpdateNoteResponse
	//oLNote = append(oLNote, oNote)
	res = UpdateNote(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
