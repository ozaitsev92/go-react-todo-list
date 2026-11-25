// Package httpserver wraps net/http.Server with safe timeouts.
package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 5 * time.Second
	defaultWriteTimeout      = 5 * time.Second
	defaultIdleTimeout       = 60 * time.Second
	defaultAddr              = ":80"
	defaultShutdownTimeout   = 3 * time.Second
)

// Server is an HTTP server that reports serve errors on a channel.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// New configures a server without starting it.
func New(handler http.Handler, opts ...Option) *Server {
	s := &Server{
		server: &http.Server{
			Handler:           handler,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
			ReadTimeout:       defaultReadTimeout,
			WriteTimeout:      defaultWriteTimeout,
			IdleTimeout:       defaultIdleTimeout,
			Addr:              defaultAddr,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start serves in the background; serve errors are sent on Notify.
func (s *Server) Start() {
	go func() {
		defer close(s.notify)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err
		}
	}()
}

// Notify receives a serve error, or is closed after Shutdown.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Addr returns the configured listen address.
func (s *Server) Addr() string {
	return s.server.Addr
}

// Shutdown waits up to the shutdown timeout for in-flight requests, then stops.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("httpserver: shutdown: %w", err)
	}

	return nil
}
