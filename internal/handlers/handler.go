package handlers

import (
	"net/http"
	"time"
)

// HTTP Handler
func RootHandler(rw http.ResponseWriter, r *http.Request, httpToUDP chan []byte, udpToHTTP chan []byte) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var body []byte
	if _, err := r.Body.Read(body); err != nil {
		http.Error(rw, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	//
	// Send to UDP
	httpToUDP <- body
	// Wait from UDP
	select {
	case resp := <-udpToHTTP:
		_, _ = rw.Write(resp)
	case <-time.After(5 * time.Second):
		http.Error(rw, "Timeout waiting for UDP response", http.StatusGatewayTimeout)
	}
}
