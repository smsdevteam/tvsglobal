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
	fmt.Fprintf(w, "Welcome to TVS Contact Restful")
}

func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvscontact", index)
	mainRouter.HandleFunc("/tvscontact/getcontact/{contactid}", getContact)
	mainRouter.HandleFunc("/tvscontact/getcontactlist/{customerid}", getContactList)
	mainRouter.HandleFunc("/tvscontact/createcontact", createContact).Methods("POST")
	mainRouter.HandleFunc("/tvscontact/updatecontact", updateContact).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}

func getContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["contactid"])

	var contactResult st.Contact

	contactResult = GetContactByContactID(params["contactid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contactResult)
}

func getContactList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["customerid"])

	var lContactResult st.ListContact

	lContactResult = GetContactListByCustomerID(params["customerid"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lContactResult)
}

func createContact(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var req st.CreateContactRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	log.Println(req)

	var res st.CreateContactResponse
	res = CreateContact(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func updateContact(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var req st.UpdateContactRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	log.Println(req)

	var res st.UpdateContactResponse
	res = UpdateContact(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
