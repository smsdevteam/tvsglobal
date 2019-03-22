package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	cm "github.com/smsdevteam/tvsglobal/common"
	st "github.com/smsdevteam/tvsglobal/tvsstructs"
)

const applicationname string = "tvssubmitorder"
const tagappname string = "icc-tvssubmitorder"
const taglogtype string = "info"

func index(w http.ResponseWriter, r *http.Request) {
	var req st.TVSSubmitOrdReqData
	fmt.Println(w, "Welcome to TVS Device Restful")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Error func main .. %s\n", err)
		}
	}()
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvssubmitorder/", index)
	mainRouter.HandleFunc("/tvssubmitorder/submitorder", submitorder).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", mainRouter))

}

func submitorder(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Error func submitorder .. %s\n", err)
		}
	}()
	var applog cm.Applog
	defer applog.PrintJSONLog()
	applog = cm.NewApploginfo("", applicationname, "submitorder",
		tagappname, taglogtype)
	//fmt.Println("start call submitorder")
	//fmt.Println("************************************************************************")
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
	applog = cm.NewApploginfo(oRes.Trackingno, applicationname, "submitorder",
		tagappname, taglogtype)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}
