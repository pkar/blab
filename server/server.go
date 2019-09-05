package server

import (
	"fmt"
	"log"
	"net"
)

// Server listens for new connections and interacts with rooms.
type Server struct {
	done     chan struct{}
	host     string
	listener net.Listener
	logDir   string
	port     int
	rooms    *Rooms
	shutdown chan struct{}
}

// New will initialize a server instance and its rooms.
func New(host string, port int, logDir string) (*Server, error) {
	s := &Server{
		shutdown: make(chan struct{}),
		done:     make(chan struct{}),
		host:     host,
		port:     port,
		rooms:    NewRooms(logDir),
	}
	network := "tcp"
	ln, err := net.Listen(network, fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	s.listener = ln
	return s, nil
}

// Start will start the server listener
func (s *Server) Start() {
	log.Println("INFO: listening", s.host, s.port)
	conns := make(chan net.Conn)
	go func(cns chan net.Conn) {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if err, ok := err.(*net.OpError); ok && !err.Timeout() {
					return
				}
				log.Println("ERRO: listen.Listener.Accept", err)
				continue
			}
			cns <- conn
		}
	}(conns)

	for {
		select {
		case conn := <-conns:
			go func() {
				_, err := NewClient(conn, s.rooms)
				if err != nil {
					log.Println("ERRO:", err)
					return
				}
			}()
		case <-s.shutdown:
			log.Println("closing rooms")
			s.rooms.Close()
			log.Println("closing server")
			s.done <- struct{}{}
			return
		}
	}
}

// Close will stop the server listener
func (s *Server) Close() {
	if err := s.listener.Close(); err != nil {
		log.Println("ERRO: Stop.Listener.Close", err)
	}
	s.shutdown <- struct{}{}
	<-s.done
}
