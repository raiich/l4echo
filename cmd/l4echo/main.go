package main

import (
	"fmt"
	"time"

	"l4echo/internal/client"
	"l4echo/internal/server"
)

func main() {
	switch env.Mode {
	case "client", "":
		count := int(env.RPS * env.Duration.Seconds())
		interval := time.Duration(float64(time.Second) / env.RPS)
		cs := []client.Config{
			{
				Concurrency: env.Concurrency,
				Completions: env.Completions,
				Network:     "udp",
				Address:     env.ClientEnv.UDPAddr,
				RecvTimeout: 0,
				SendTimeout: 0,
				Workload: client.Workload{
					Count:    count,
					Interval: interval,
				},
			},
			{
				Concurrency: env.Concurrency,
				Completions: env.Completions,
				Network:     "tcp",
				Address:     env.ClientEnv.TCPAddr,
				RecvTimeout: 10 * time.Second,
				SendTimeout: 1 * time.Second,
				Workload: client.Workload{
					Count:    count,
					Interval: interval,
				},
			},
		}
		client.Run(cs)
	case "server":
		s := server.Server{
			TCPAddr: env.ServerEnv.TCPAddr,
			UDPAddr: env.ServerEnv.UDPAddr,
		}
		s.Serve()
	default:
		panic(fmt.Errorf("invalid mode: %s", env.Mode))
	}
}
