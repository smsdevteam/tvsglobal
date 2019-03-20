package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// util
	// referpath
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TVS Device Restful")
}

func mainold() {
	fmt.Println("Service Start...")
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/tvsqueueorder", index)
	log.Fatal(http.ListenAndServe(":8081", mainRouter))

}
