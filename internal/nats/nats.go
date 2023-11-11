package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"wb-tech-level-0/internal/config"
	"wb-tech-level-0/internal/model"
	storage "wb-tech-level-0/internal/storage/postgres"
)

type Nats struct {
	config         *config.Config
	StanConnection stan.Conn
	postgres       *storage.Postgres
	ctx            context.Context
}

func NewNats(cfg *config.Config, postgres *storage.Postgres, ctx context.Context) (*Nats, error) {
	sc, err := stan.Connect(cfg.Nats.Cluster, cfg.Nats.Client)
	if err != nil {
		return nil, err
	}

	return &Nats{cfg, sc, postgres, ctx}, nil
}

func (n *Nats) Close() error {
	err := n.StanConnection.Close()
	if err != nil {
		return err
	}

	return nil
}

func (n *Nats) Subscribe() error {
	_, err := n.StanConnection.Subscribe(n.config.Nats.Topic, n.handleMessage)
	if err != nil {
		log.Print(err)
	}

	return nil
}

func (n *Nats) handleMessage(msg *stan.Msg) {
	var result model.Order

	err := json.Unmarshal(msg.Data, &result)
	if err != nil {
		log.Println(err)
	}

	err = n.postgres.AddOrder(n.ctx, result)
	if err != nil {
		log.Printf("Ошибка вставки данных в базу данных: %v", err)
		return
	}

}
