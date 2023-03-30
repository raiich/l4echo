package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	env = &Env{
		Mode: "",
		ClientEnv: ClientEnv{
			Concurrency: 32,
			Completions: 96,
			Duration:    30 * time.Second,
			RPS:         0.2,
			TCPAddr:     "localhost:7000",
			UDPAddr:     "localhost:7001",
		},
		ServerEnv: ServerEnv{
			TCPAddr: ":7000",
			UDPAddr: ":7001",
		},
	}
)

func init() {
	if err := envconfig.Process("", env); err != nil {
		panic(err)
	}
}

type Env struct {
	Mode string
	ClientEnv
	ServerEnv
}

type ClientEnv struct {
	Concurrency int
	Completions int
	Duration    time.Duration
	RPS         float64
	TCPAddr     string
	UDPAddr     string
}

type ServerEnv struct {
	TCPAddr string
	UDPAddr string
}
