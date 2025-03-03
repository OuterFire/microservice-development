package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func NewHttpServer(mux *http.ServeMux, port int, writeTimeout int) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      mux,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}
}

type ErrorResponse struct {
	Status      int    `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func WriteResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	dataJson, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(dataJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func WriteStatusCode(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func Marshal(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
