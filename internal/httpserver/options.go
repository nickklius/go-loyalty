package httpserver

import (
	"time"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		//if port[0:1] == ":" {
		s.server.Addr = ":8080"
		//} else {
		//	s.server.Addr = net.JoinHostPort("", port)
		//}
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
