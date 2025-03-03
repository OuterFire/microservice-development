package utils

import "net/http"

type MuxHandler struct {
	Mux      *http.ServeMux
	Handlers map[string]http.HandlerFunc
}

func NewMuxHandler() *MuxHandler {
	mux := http.NewServeMux()
	return &MuxHandler{
		Mux:      mux,
		Handlers: make(map[string]http.HandlerFunc),
	}
}

func (s *MuxHandler) AddHandlers(endpoint string, handler http.HandlerFunc) {
	s.Handlers[endpoint] = handler
}
