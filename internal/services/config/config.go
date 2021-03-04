package config

import (
	"os"
)

// ByKey uses godot package to load/read the .env file and return value for key
func ByKey(key string) string {
	return os.Getenv(key)
}
