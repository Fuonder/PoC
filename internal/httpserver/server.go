package httpserver

import (
	"PoC/internal/handlers"
	"errors"
	"log"
	"net/http"
)

// HTTP server
type Server struct {
	httpAddr   string
	httpServer *http.Server
	httpToUDP  chan []byte // channel HTPP to UDP
	udpToHTTP  chan []byte // channel UDP to HTTP

}

// Creates new server
func NewServer(httpAddr string, httpToUDP, udpToHTTP chan []byte) *Server {
	return &Server{
		httpAddr:  httpAddr,
		httpToUDP: httpToUDP,
		udpToHTTP: udpToHTTP,
	}
}

// HTTP server start
func (s *Server) Run() error {
	http.HandleFunc("/send", func(rw http.ResponseWriter, r *http.Request) {
		handlers.RootHandler(rw, r, s.httpToUDP, s.udpToHTTP)
	})
	s.httpServer = &http.Server{Addr: s.httpAddr}

	if err := s.httpServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Graceful http server exit")
			return nil
		}
		return err
	}
	return nil
}

// Graceful server stop
func (s *Server) Stop() error {
	return s.httpServer.Close()
}
