package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"

	"bin-collections-api/internal/pkg/redis-pool"
	"bin-collections-api/internal/address-search"
	"bin-collections-api/internal/collections-for-address"
)

func main() {
	fmt.Println("yes hello")
	fmt.Println(os.Getenv("PORT"))

	pool := redispool.NewPool("redis")
	conn := pool.Get()

	pong, err := conn.Do("PING")

	fmt.Println(redis.String(pong, err))

	godotenv.Load(".env.yml")
	port := os.Getenv("PORT")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/address/search", addresssearch.FindAddressByPostCode).Methods("POST")
	router.HandleFunc("/collections", getcollectiondates.GetCollectionsForID).Methods("POST")
	log.Fatal(http.ListenAndServe(":" + port, router))
}
