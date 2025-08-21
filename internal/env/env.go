package env

import (
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
	key, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return key
}

func GetInt(key string, fallback int) int {
	key, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	intKey, err := strconv.Atoi(key)
	if err != nil {
		return fallback
	}
	return intKey
}
