package config

import (
	"os"
)

func ByKey(key string) string {
	return os.Getenv(key)
}
