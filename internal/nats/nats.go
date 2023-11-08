package nats

import (
	"github.com/nats-io/stan.go"
	"wb-tech-level-0/internal/config"
)

type Nats struct {
	config         *config.Config
	stanConnection *stan.Conn
}

func NewNats(cfg *config.Config) (*Nats, error) {
	sc, err := stan.Connect(cfg.Nats.Cluster, cfg.Nats.Client)
	if err != nil {
		return nil, err
	}

	return &Nats{cfg, &sc}, nil
}
