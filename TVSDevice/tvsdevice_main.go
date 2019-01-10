package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	st "github.com/pimpina/tvsglobalb/TVSStructs" // referpath
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Device Restful")
}

func getDeviceData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["sn"])
	var odv st.Device
	odv = GetDeviceBySerialNumber(params["sn"])
	/*var odv st.ResponseResult
	odv = GetDeviceFromAPI(params["sn"])*/

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(odv)
}

func moveDepot(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//fmt.Println(params["sn"])
	//fmt.Println(params["depottoid"])
	//fmt.Println(params["reason"])
	var sn string
	var depotid, reason int64

	sn = params["sn"]
	depotid, _ = strconv.ParseInt(params["depottoid"], 10, 64)
	reason, _ = strconv.ParseInt(params["reason"], 10, 64)
	var oRes st.ResponseResult
	oRes = MoveDevice(sn, depotid, reason, params["by"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}

func pairingDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//fmt.Println(params["sn"])
	//fmt.Println(params["depottoid"])
	//fmt.Println(params["reason"])
	var devicefrom, deviceto, reason int64

	devicefrom, _ = strconv.ParseInt(params["devicefrom"], 10, 64)
	deviceto, _ = strconv.ParseInt(params["deviceto"], 10, 64)
	reason, _ = strconv.ParseInt(params["reason"], 10, 64)
	var oRes st.ResponseResult
	oRes = PairOneDeviceToAnother(devicefrom, deviceto, reason, params["by"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}

func sendCmd(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//fmt.Println(params["sn"])
	//fmt.Println(params["reason"])
	var reason int64

	reason, _ = strconv.ParseInt(params["reason"], 10, 64)
	var oRes st.ResponseResult
	oRes = SendCommandToDevice(params["sn"], reason, params["by"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}

func main() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvsdevice", index)
	mainRouter.HandleFunc("/tvsdevice/getdevicebyserialnumber/{sn}", getDeviceData)
	mainRouter.HandleFunc("/tvsdevice/movedevice/{sn}/{depottoid}/{reason}/{by}", moveDepot)
	mainRouter.HandleFunc("/tvsdevice/paironedevicetoanother/{devicefrom}/{deviceto}/{reason}/{by}", pairingDevice)
	mainRouter.HandleFunc("/tvsdevice/sendcmdtodevice/{sn}/{reason}/{by}", sendCmd)

	log.Fatal(http.ListenAndServe(":8081", mainRouter))
}
