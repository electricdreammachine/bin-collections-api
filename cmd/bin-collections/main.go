package main

import (
	scrapercontroller "bin-collections-api/internal/controllers"
	config "bin-collections-api/internal/config"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(config.ProjectRootPath + "/config/.env.yml")
	port := os.Getenv("PORT")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/address/search", scrapercontroller.GetAddresses).Methods("POST")
	router.HandleFunc("/collections", scrapercontroller.GetCollectionDates).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
