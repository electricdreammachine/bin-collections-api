package getconfigvalue

import (
	"github.com/joho/godotenv"
)

// ByKey uses godot package to load/read the .env file and return value for key
func ByKey(key string) string {
	configMap,_ := godotenv.Read(".env.yml")

	return configMap[key]
}