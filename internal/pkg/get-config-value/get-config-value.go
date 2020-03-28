package getconfigvalue

import (
	"log"
	"github.com/joho/godotenv"
)

// ByKey uses godot package to load/read the .env file and return value for key
func ByKey(key string) string {
	configMap, err := godotenv.Read(".env.yml")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return configMap[key]
}