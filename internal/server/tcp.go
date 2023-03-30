package server

import (
	"errors"
	"io"
	"net"
	"time"

	"l4echo/internal/log"
	"l4echo/internal/network"
	"l4echo/internal/payload"
)

func ServeTCP(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Info("listening tcp:", address)
	defer func() {
		if err := lis.Close(); err != nil {
			log.Error("error in close:", err)
		}
	}()

	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		go start(conn)
	}
}

func start(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("error in close:", conn.RemoteAddr(), err)
		} else {
			log.Error("closed conn:", conn.RemoteAddr())
		}
	}()

	log.Info("accepted TCP conn:", conn.RemoteAddr())
	conn = &network.ConnWithTimeout{
		Conn:         conn,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	buf := make([]byte, payload.Size)
	for {
		p, err := payload.Receive(conn, buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Error("failed to receive (tcp):", err)
				return
			}
			log.Info("tcp receive loop finished")
			return
		}
		if err := payload.Send(conn, p); err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Error("failed to send (tcp): ", err)
				return
			}
			log.Info("tcp send loop finished")
			return
		}
	}
}
