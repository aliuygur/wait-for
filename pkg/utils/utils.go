package utils

import (
	"os"
	"strconv"
)

func GetOsEnvOrDefault(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return value
}

func GetOsEnvNumberOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	number, _ := strconv.Atoi(value)
	return number
}
