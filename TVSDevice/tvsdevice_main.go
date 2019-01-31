package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	cm "github.com/smsdevteam/tvsglobal/Common"     // util
	st "github.com/smsdevteam/tvsglobal/TVSStructs" // referpath
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Device Restful")
}

func getDeviceBySerialNumber(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["sn"])
	var odv st.GetDeviceResponse
	odv = GetDeviceBySerialNumber(params["sn"], params["by"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(odv)
}

func getDeviceData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["sn"])
	var odv st.DeviceInfo
	odv = GetDataSerialNumber(params["sn"])
	/*var odv st.ResponseResult
	odv = GetDeviceFromAPI(params["sn"])*/

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(odv)
}

func moveDepot(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var deviceid, depotid, reason int64
	deviceid = cm.StrToInt64(params["deviceid"])
	depotid = cm.StrToInt64(params["depottoid"])
	reason = cm.StrToInt64(params["reason"])

	var oRes st.ResponseResult
	oRes = MoveDevice(deviceid, depotid, reason, params["by"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}

func pairingDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var devicefrom, deviceto, reason int64
	devicefrom = cm.StrToInt64(params["devicefrom"])
	deviceto = cm.StrToInt64(params["deviceto"])
	reason = cm.StrToInt64(params["reason"])

	var oRes st.ResponseResult
	oRes = PairOneDeviceToAnother(devicefrom, deviceto, reason, params["by"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
}

func sendCmd(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var deviceid, reason int64

	deviceid = cm.StrToInt64(params["deviceid"])
	reason = cm.StrToInt64(params["reason"])
	var oRes st.ResponseResult
	oRes = SendCommandToDevice(deviceid, reason, params["by"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oRes)
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

	// var tRes st.ResponseResult
	tRes := st.NewResponseResult()
	log.Println(*tRes)
	tRes.ErrorCode = 0
	tRes.ErrorDesc = "TEST"
	log.Println(*tRes)

	var nRes st.ResponseResult
	log.Println(nRes)
	nRes.ErrorCode = 100
	nRes.ErrorDesc = "Test#2"
	log.Println(nRes)

	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvsdevice", index)
	mainRouter.HandleFunc("/tvsdevice/getdevicebyserialnumber/{sn}/{by}", getDeviceBySerialNumber)
	mainRouter.HandleFunc("/tvsdevice/getdevicebyserialnumber/{sn}", getDeviceBySerialNumber)
	mainRouter.HandleFunc("/tvsdevice/getdevicedata/{sn}", getDeviceData)
	mainRouter.HandleFunc("/tvsdevice/movedevice/{deviceid}/{depottoid}/{reason}/{by}", moveDepot)
	mainRouter.HandleFunc("/tvsdevice/paironedevicetoanother/{devicefrom}/{deviceto}/{reason}/{by}", pairingDevice)
	mainRouter.HandleFunc("/tvsdevice/sendcmdtodevice/{deviceid}/{reason}/{by}", sendCmd)
	mainRouter.HandleFunc("/tvsdevice/createnewserialnumber", createNewSN).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", mainRouter))
}
