package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"github.com/joho/godotenv"

	"bin-collections-api/internal/address-search"
	"bin-collections-api/internal/collections-for-address"
)

func main() {
	godotenv.Load(".env.yml")
	port := os.Getenv("PORT")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/address/search", addresssearch.FindAddressByPostCode).Methods("POST")
	router.HandleFunc("/collections", getcollectiondates.GetCollectionsForID).Methods("POST")
	log.Fatal(http.ListenAndServe(":" + port, router))
}
