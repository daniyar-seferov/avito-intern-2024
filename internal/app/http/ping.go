package http

import (
	"net/http"
)

// PingHandler struct.
type PingHandler struct{}

// NewPingHandler returns ping handler.
func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
