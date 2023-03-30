package server

import (
	"l4echo/internal/log"
)

type Server struct {
	TCPAddr string
	UDPAddr string
}

func (s *Server) Serve() {
	go func() {
		if err := ServeTCP(s.TCPAddr); err != nil {
			log.Error("failed to serve tcp:", err)
		}
	}()
	if err := ServeUDP(s.UDPAddr); err != nil {
		log.Error("error in udp: ", err)
	}
}
