package cache

import (
	"sync"
	"wb-tech-level-0/internal/model"
	storage "wb-tech-level-0/internal/storage/postgres"
)

type Cache struct {
	sync.RWMutex
	items map[string]*model.Order
	db    *storage.Postgres
}

func NewCache(db *storage.Postgres) *Cache {
	return &Cache{items: nil, db: db}
}
