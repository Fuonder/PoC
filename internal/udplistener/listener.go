package udplistener

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
)

// UDP Listener
type Listener struct {
	udpAddr   string
	udpConn   *net.UDPConn
	stopCh    chan os.Signal // channel for graceful stop
	httpToUDP chan []byte    // Chanel HTTP to UDP
	udpToHTTP chan []byte    // Channel UDP to HTTP
}

// Creates new UDP Listener
func NewListener(udpAddr string, httpToUDP, udpToHTTP chan []byte) *Listener {
	return &Listener{
		udpAddr:   udpAddr,
		stopCh:    make(chan os.Signal),
		httpToUDP: httpToUDP,
		udpToHTTP: udpToHTTP,
	}
}

// Start + main handler for udp
func (l *Listener) Run() error {
	udpAddr, err := net.ResolveUDPAddr("udp", l.udpAddr)
	if err != nil {
		return fmt.Errorf("error resolving UDP address: %v", err)
	}
	l.udpConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("error listening to UDP: %v", err)
	}
	defer l.udpConn.Close()

	// Get data from HTTP and Send to CLIENT via UDP
	go func() {
		for {
			select {
			case <-l.stopCh: // we receive signal to stop
				return
			case data := <-l.httpToUDP: // Got command from "Terminal" (from HTTP)

				fmt.Printf("UDP received from HTTP: %s\n", string(data))
				// Here Process what should be done with data
				//      from HTTP and how it to send to udp CLIENT
				// When write packet for active client to UDP
				_, err = l.udpConn.WriteToUDP([]byte("This is packet for client"), &net.UDPAddr{IP: net.IP{}, Port: 10}) //Ip will be field later from Storage. This is temporary.
				if err != nil {
					fmt.Printf("Error sending UDP response: %v\n", err)
				}
			default:
			}
		}
	}()

	// this is Receiver form CLIENTS
	for {
		buf := make([]byte, 1024)
		_, addr, err := l.udpConn.ReadFromUDP(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("UDP CON CLOSED, GR EXIT")
				return nil
			}
			fmt.Printf("Error reading from UDP: %v\n", err)
			return err
		}
		fmt.Printf("Received UDP message: %s from %v\n", string(buf), addr)

		// Here process response
		// Send response to HTTP
		l.udpToHTTP <- []byte("Response from client to HTTP server")

	}
}

// Stop Gracefully
func (l *Listener) Stop() error {
	err := l.udpConn.Close()
	l.stopCh <- syscall.SIGTERM
	close(l.stopCh)
	return err
}
