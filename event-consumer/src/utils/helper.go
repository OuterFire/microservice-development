package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvString(env string) string {
	value := strings.TrimSpace(os.Getenv(env))
	return value
}

func GetEnvInt(env string) int {
	value := strings.TrimSpace(os.Getenv(env))
	intValue, err := strconv.Atoi(value)
	if err != nil && value != "" {
		return 0
	}
	return intValue
}
