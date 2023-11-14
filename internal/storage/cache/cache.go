package cache

import (
	"context"
	"sync"
	"wb-tech-level-0/internal/model"
	storage "wb-tech-level-0/internal/storage/postgres"
)

type Cache struct {
	sync.RWMutex
	items map[string]model.Order
	db    *storage.Postgres
}

func NewCache(db *storage.Postgres, ctx context.Context) *Cache {
	var cache Cache
	cache.items = make(map[string]model.Order)
	cache.db = db
	orders, err := db.Orders(ctx)
	if err != nil {
		return nil
	}
	for _, order := range orders {
		cache.items[order.OrderUID] = order
	}
	return &cache
}

func (c *Cache) AddOrder(order model.Order) {
	c.Lock()
	defer c.Unlock()
	c.items[order.OrderUID] = order
}

func (c *Cache) GetOrder(id string) (*model.Order, bool) {
	c.RLock()
	defer c.RUnlock()
	item, found := c.items[id]

	// ключ не найден
	if !found {
		return nil, false
	}

	return &item, true
}
