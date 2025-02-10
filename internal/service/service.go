package service

import (
	"PoC/internal/Storage"
	"PoC/internal/httpserver"
	"PoC/internal/repository"
	"PoC/internal/udplistener"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	httpServer    *httpserver.Server
	udpServer     repository.MyService
	httpToUDP     chan []byte // Channel HTTP to UDP
	udpToHTTP     chan []byte // Channel UDP to HTTP
	stopCh        chan os.Signal
	clientStorage map[string]Storage.Client
}

func NewService(httpAddr, udpAddr string) *Service {
	httpToUDP := make(chan []byte)
	udpToHTTP := make(chan []byte)

	return &Service{
		httpServer: httpserver.NewServer(httpAddr, httpToUDP, udpToHTTP),
		udpServer:  udplistener.NewListener(udpAddr, httpToUDP, udpToHTTP),
		httpToUDP:  httpToUDP,
		udpToHTTP:  udpToHTTP,
		stopCh:     make(chan os.Signal, 1),
	}
}

func (s *Service) Run() error {
	// HTTP server Start
	errors := make(chan error)

	go func() {
		fmt.Println("starting HTTP")
		if err := s.httpServer.Run(); err != nil {
			errors <- fmt.Errorf("method <run()>: http server: %v", err)
			return
		}
	}()
	// UDP listener start
	go func() {
		fmt.Println("starting UDP")
		if err := s.udpServer.Run(); err != nil {
			errors <- fmt.Errorf("method <run()>: udp server: %v", err)
			return
			//log.Fatalf("error running udp server: %v\n", err)
		}
	}()
	signal.Notify(s.stopCh, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-s.stopCh:
			fmt.Println("graceful exit started")
			s.Stop()
			return nil
		default:
			select {
			case err := <-errors:

				fmt.Println("ERROR!! EXITING...")
				return fmt.Errorf("service: %v", err)
			default:
				continue
			}
		}

	}
}

func (s *Service) Stop() {
	if err := s.httpServer.Stop(); err != nil {
		fmt.Printf("Error stopping HTTP server: %v\n", err)
	}
	if err := s.udpServer.Stop(); err != nil {
		fmt.Printf("Error stopping UDP server: %v\n", err)
	}
	close(s.stopCh)
}
