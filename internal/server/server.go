package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"myhttpfromtcp/internal/request"
	"myhttpfromtcp/internal/response"
)

const CONNECTION_TYPE = "tcp"

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen(CONNECTION_TYPE, fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Coun't create listener '%s' by port :%d\n", CONNECTION_TYPE, port)
		return nil, err
	}

	s := &Server {
		listener: listener,
		handler:  handler,
	}

	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	w := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.STATUS_CODE_BAD_REQUEST)

		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}

	s.handler(w, req)
}
