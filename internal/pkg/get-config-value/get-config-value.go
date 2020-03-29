package getconfigvalue

import (
	"os"
	"fmt"
)

// ByKey uses godot package to load/read the .env file and return value for key
func ByKey(key string) string {
	fmt.Println(os.Getenv(key))
	return os.Getenv(key)
}