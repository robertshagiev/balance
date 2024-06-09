package server

import (
	"log"
	"net/http"
)

type Server struct {
	handler ser
	host    string
	port    string
}

type ser interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Introduction(w http.ResponseWriter, r *http.Request)
	Debit(w http.ResponseWriter, r *http.Request)
	Transfer(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
}

func NewServer(h ser, host, port string) *Server {
	return &Server{
		handler: h,
		host:    host,
		port:    port,
	}
}

func (s *Server) Start() error {
	addr := s.host + ":" + s.port
	log.Printf("Server is starting at %s\n", addr)
	return http.ListenAndServe(addr, s.handler)
}
