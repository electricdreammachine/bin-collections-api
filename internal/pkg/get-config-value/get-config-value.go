package getconfigvalue

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

// ByKey uses godot package to load/read the .env file and return value for key
func ByKey(key string) string {
	err := godotenv.Load(".env.yml")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

return os.Getenv(key)
}