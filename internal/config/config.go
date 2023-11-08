package config

import (
	"fmt"
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

type Config struct {
	DatabaseDSN string `required:"true"`
	Nats        struct {
		Host    string `required:"true"`
		Port    string `required:"true"`
		Cluster string `required:"true"`
		Client  string `required:"true"`
		Topic   string `required:"true"`
	} `required:"true" toml:"nats" yaml:"nats"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{}
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		Files: []string{"config.yml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yml": aconfigyaml.New(),
		},
	})

	if err := loader.Load(); err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		return nil, err
	}

	return cfg, nil
}
