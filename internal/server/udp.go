package server

import (
	"fmt"
	"net"

	"l4echo/internal/log"
	"l4echo/internal/payload"
)

func ServeUDP(address string) error {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error("error in close:", err)
		}
	}()

	log.Info("listening udp:", address)
	buf := make([]byte, 1500)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}
		p, err := payload.Validate(buf[:n])
		if err != nil {
			return err
		}
		m, err := conn.WriteTo(p, addr)
		if err != nil {
			return err
		}
		if n != m {
			return fmt.Errorf("invalid write size: %d", m)
		}
	}
}
