package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"wb-tech-level-0/internal/config"
	"wb-tech-level-0/internal/http-server/handlers"
	"wb-tech-level-0/internal/nats"
	"wb-tech-level-0/internal/storage/cache"
	storage "wb-tech-level-0/internal/storage/postgres"
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	postgres := storage.NewPostgres(db)

	cache := cache.NewCache(postgres, ctx)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", handlers.NewIndex())
	router.Get("/orders/{orderID}", handlers.NewGetOrder(cache))

	go func() {
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			log.Fatal(err)
		}
	}()

	natsStream, err := nats.NewNats(cfg, postgres, ctx, cache)
	if err != nil {
		panic(err)
	}
	defer func(natsStream *nats.Nats) {
		err := natsStream.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(natsStream)

	//генерация случайных заказов
	/*	gnrt := generator.NewGenerator()

		go func() {
			for {
				order := gnrt.GenerateOrder()
				err := natsStream.Publish(order)
				if err != nil {
					log.Println("не удалось отправить заказ")
					time.Sleep(10 * time.Second)
					continue
				}
				log.Println("отправил сгенерированный заказ")

				time.Sleep(15 * time.Second)
			}
		}()*/

	go natsStream.Subscribe()

	for {
		select {
		case <-ctx.Done():
			log.Print("завершение работы")
			return
		}
	}
}
