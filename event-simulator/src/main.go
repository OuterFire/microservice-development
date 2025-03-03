package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := getEnvInt("EVENT_PRODUCER_PORT")
	host := getEnvString("EVENT_PRODUCER_HOST")
	num := getEnvInt("EVENTS_PER_SECOND")
	perSecond := 1000 / num

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	for {
		time.Sleep(time.Duration(perSecond) * time.Millisecond)

		jsonBody := []byte(`{"description":"hello world"}`)
		bodyReader := bytes.NewReader(jsonBody)

		requestURL := fmt.Sprintf("http://%v:%v/event", host, port)
		res, err := http.Post(requestURL, "application/json", bodyReader)
		if err != nil {
			logger.Error("error making http request: %s\n", err)
			continue
		}

		logger.Info(fmt.Sprintf("Status Code: %v", res.StatusCode))
	}
}

func getEnvString(env string) string {
	envValue := strings.TrimSpace(os.Getenv(env))
	return envValue
}

func getEnvInt(env string) int {
	envValue := strings.TrimSpace(os.Getenv(env))
	intValue, err := strconv.Atoi(envValue)
	if err != nil && envValue != "" {
		return 0
	}
	return intValue
}
