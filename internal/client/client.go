package client

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"l4echo/internal/log"
	"l4echo/internal/payload"
)

type Client struct {
	conn net.Conn
	*Stats
}

type Workload struct {
	Count    int
	Interval time.Duration
}

func (c *Client) Start(w *Workload) {
	wait := make(chan struct{}, 1)
	go func() {
		if err := c.receiveLoop(w); err != nil {
			log.Error("failed to receive:", err)
		}
		close(wait)
	}()
	if err := c.sendLoop(w); err != nil {
		log.Error("failed to send:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil && err != context.DeadlineExceeded {
			log.Error("error in Loop:", err)
		}
		return
	case <-wait:
		return
	}
}

func (c *Client) sendLoop(w *Workload) error {
	p := payload.New()
	for i := 0; i < w.Count; i++ {
		p.SetSeq(uint64(i))
		p.SetTime(time.Now())
		if err := payload.Send(c.conn, p); err != nil {
			return err
		}
		c.Stats.OnSending()
		time.Sleep(w.Interval)
	}
	return nil
}

func (c *Client) receiveLoop(w *Workload) error {
	next := uint64(0)
	buf := make([]byte, payload.Size)
	for i := 0; i < w.Count; i++ {
		p, err := payload.Receive(c.conn, buf)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		seq := p.Seq()
		if seq == next {
			c.Stats.OnReceived(p.Time())
			next = seq + 1
		} else if seq > next {
			c.Stats.Dropped(seq - next)
			c.Stats.OnReceived(p.Time())
			next = seq + 1
		} else {
			c.Stats.Delayed()
		}
	}
	return nil
}
