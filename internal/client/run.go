package client

import (
	"context"
	"net"
	"sync"
	"time"

	"l4echo/internal/log"
	"l4echo/internal/network"

	"golang.org/x/sync/semaphore"
)

func Run(configs []Config) {
	wg := &sync.WaitGroup{}
	for _, config := range configs {
		wg.Add(1)
		c := config
		go func() {
			defer wg.Done()
			run(c)
		}()
	}
	wg.Wait()
}

func run(config Config) {
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	s := semaphore.NewWeighted(int64(config.Concurrency))
	stats := &Stats{Name: config.Network}

	for i := 0; i < config.Completions; i++ {
		if err := s.Acquire(ctx, 1); err != nil {
			log.Errorf("failed to acquire: %v", err)
			continue
		}
		conn, err := net.Dial(config.Network, config.Address)
		if err != nil {
			log.Errorf("failed to dial: %v", err)
			continue
		}
		if config.Network == "tcp" {
			conn = &network.ConnWithTimeout{
				Conn:         conn,
				ReadTimeout:  config.RecvTimeout,
				WriteTimeout: config.SendTimeout,
			}
		}
		client := &Client{
			conn:  conn,
			Stats: stats,
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer s.Release(1)
			defer func() {
				if err := conn.Close(); err != nil {
					log.Errorf("failed to close: %v", err)
				}
			}()
			log.Info("starting client")
			client.Start(&config.Workload)
		}()
	}
	wg.Wait()
	log.Infof("stats: %+v", stats)
}

type Config struct {
	Concurrency int
	Completions int
	Network     string
	Address     string
	RecvTimeout time.Duration
	SendTimeout time.Duration
	Workload
}
