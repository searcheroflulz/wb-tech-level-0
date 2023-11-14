package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestGettingAll(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/wb-tech-level-0?sslmode=disable")
	if err != nil {
		log.Printf("[ERROR] failed to connect to db: %v", err)
		return
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	postgres := NewPostgres(db)

	orders, err := postgres.Orders(ctx)
	if err != nil {
		return
	}

	for _, order := range orders {
		jsonString, _ := json.Marshal(order)
		fmt.Println(string(jsonString))
	}

}
