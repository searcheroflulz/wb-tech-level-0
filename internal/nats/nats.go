package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"wb-tech-level-0/internal/config"
	"wb-tech-level-0/internal/model"
	"wb-tech-level-0/internal/storage/cache"
	storage "wb-tech-level-0/internal/storage/postgres"
)

type Nats struct {
	config         *config.Config
	stanConnection stan.Conn
	postgres       *storage.Postgres
	ctx            context.Context
	cache          *cache.Cache
}

func NewNats(cfg *config.Config, postgres *storage.Postgres, ctx context.Context, cache *cache.Cache) (*Nats, error) {
	sc, err := stan.Connect(cfg.Nats.Cluster, cfg.Nats.Client,
		stan.Pings(10, 5),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		return nil, err
	}

	return &Nats{
		config:         cfg,
		stanConnection: sc,
		postgres:       postgres,
		ctx:            ctx,
		cache:          cache}, nil
}

func (n *Nats) Close() error {
	err := n.stanConnection.Close()
	if err != nil {
		return err
	}

	return nil
}

func (n *Nats) Subscribe() {
	_, err := n.stanConnection.Subscribe(n.config.Nats.Topic, n.handleMessage,
		stan.DurableName("my-durable"))

	if err != nil {
		log.Println(err)
	}
}

func (n *Nats) handleMessage(msg *stan.Msg) {
	var result model.Order

	err := json.Unmarshal(msg.Data, &result)
	if err != nil {
		log.Println(err)
		return
	}
	if result.OrderUID == "" {
		log.Println("некорректное значение из канала")
		return
	}
	err = n.postgres.AddOrder(n.ctx, result)
	if err != nil {
		log.Printf("Ошибка вставки данных в базу данных: %v", err)
		return
	}
	go n.cache.AddOrder(result)
	log.Println("принял заказ и отправил в базу данных")
}

func (n *Nats) Publish(order *model.Order) error {
	marshal, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = n.stanConnection.Publish(n.config.Nats.Topic, marshal)
	if err != nil {
		return err
	}

	return nil
}
