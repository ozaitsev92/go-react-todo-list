package httpserver

import (
	"net"
	"time"
)

// Option customises a Server.
type Option func(*Server)

// Port sets the TCP port to listen on (all interfaces).
func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

// Addr sets the full listen address, e.g. "127.0.0.1:8080".
func Addr(addr string) Option {
	return func(s *Server) {
		s.server.Addr = addr
	}
}

// ReadHeaderTimeout bounds the time allowed to read request headers.
func ReadHeaderTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadHeaderTimeout = timeout
	}
}

// ReadTimeout bounds the time allowed to read the entire request.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// WriteTimeout bounds the time allowed to write the response.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

// IdleTimeout bounds how long keep-alive connections stay idle.
func IdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.IdleTimeout = timeout
	}
}

// ShutdownTimeout bounds how long Shutdown waits for in-flight requests.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
