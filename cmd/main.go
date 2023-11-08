package main

import (
	"fmt"
	"log"
	"wb-tech-level-0/internal/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)
}
