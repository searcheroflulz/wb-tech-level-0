package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"log"
	"wb-tech-level-0/internal/config"
	"wb-tech-level-0/internal/nats"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sqlx.Connect("postgres", cfg.DatabaseDSN)
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

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db.DB, "internal/storage/migrations"); err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	_, err = nats.NewNats(cfg)
	if err != nil {
		panic(err)
	}

}
