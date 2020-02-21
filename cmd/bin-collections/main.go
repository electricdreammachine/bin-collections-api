package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"

	"bin-collections-api/internal/address-search"
	"bin-collections-api/internal/collections-for-address"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/address/search", addresssearch.FindAddressByPostCode).Methods("POST")
	router.HandleFunc("/collections", getcollectiondates.GetCollectionsForID).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
