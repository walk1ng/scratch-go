package main

import (
	"crypto/tls"
	"time"
)

type Server struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}

type Option func(*Server)

func Protocol(p string) Option {
	return func(s *Server) {
		s.Protocol = p
	}
}

func Timeout(t time.Duration) Option {
	return func(s *Server) {
		s.Timeout = t
	}
}

func MaxConns(conn int) Option {
	return func(s *Server) {
		s.MaxConns = conn
	}
}

func TLS(tlsConfig *tls.Config) Option {
	return func(s *Server) {
		s.TLS = tlsConfig
	}
}

func NewServer(addr string, port int, options ...Option) *Server {
	srv := &Server{
		Addr:     addr,
		Port:     port,
		Protocol: "tcp",
		Timeout:  30 * time.Second,
		MaxConns: 1000,
		TLS:      nil,
	}

	for _, opt := range options {
		opt(srv)
	}

	return srv
}

func main() {
	var _ *Server = NewServer("localhost", 8080, Protocol("udp"), MaxConns(101))
}
