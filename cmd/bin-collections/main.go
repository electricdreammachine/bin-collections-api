package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	addresssearch "bin-collections-api/internal/address-search"
	getcollectiondates "bin-collections-api/internal/collections-for-address"
)

func main() {
	fmt.Println("yes hello")
	fmt.Println(os.Getenv("PORT"))

	// pool := redispool.NewPool("redis")
	// conn := pool.Get()

	// pong, err := conn.Do("PING")

	// fmt.Println(redis.String(pong, err))

	godotenv.Load(".env.yml")
	port := os.Getenv("PORT")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/address/search", addresssearch.FindAddressByPostCode).Methods("POST")
	router.HandleFunc("/collections", getcollectiondates.GetCollectionsForID).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
