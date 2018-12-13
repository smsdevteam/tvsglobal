package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	tvs_s "github.com/smsdevteam/tvsglobal/tvsstructs"
)

type TVS_Note_response struct {
	noteObj    []tvs_s.Note
	ResultCode string
}

func main() {
	fmt.Println("start service")
	router := mux.NewRouter()
	router.HandleFunc("/getNote/{noteid:[0-9]+}", getNote).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
	fmt.Println("start service.....")
}

func getNote(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	fmt.Print(vars)
	noteID, err := strconv.Atoi(vars["noteid"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	log.Println("get note")
	res := getNoteDB(noteID)
	fmt.Print("res is ", res)
	respondWithJSON(w, http.StatusOK, res)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func getNoteDB(noteID int) TVS_Note_response {
	arrayNoteObj := []tvs_s.Note{}
	noteObj2 := tvs_s.Note{}
	noteObj2.CustomerID = noteID
	fmt.Println("customer is ", noteObj2.CustomerID)
	res := TVS_Note_response{}
	arrayNoteObj = append(arrayNoteObj, noteObj2)
	res.noteObj = arrayNoteObj
	res.ResultCode = "1"

	fmt.Println(arrayNoteObj)

	return res
}
