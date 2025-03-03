package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvString(env string) string {
	envValue := strings.TrimSpace(os.Getenv(env))
	return envValue
}

func GetEnvInt(env string) int {
	envValue := strings.TrimSpace(os.Getenv(env))
	intValue, err := strconv.Atoi(envValue)
	if err != nil && envValue != "" {
		return 0
	}
	return intValue
}
