package network

import (
	"net"
	"time"
)

type ConnWithTimeout struct {
	net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *ConnWithTimeout) Read(p []byte) (n int, err error) {
	if err := c.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout)); err != nil {
		return -1, err
	}
	return c.Conn.Read(p)
}

func (c *ConnWithTimeout) Write(p []byte) (n int, err error) {
	if err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout)); err != nil {
		return -1, err
	}
	return c.Conn.Write(p)
}
