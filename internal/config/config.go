package config

import (
	"fmt"
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

type Config struct {
	DatabaseDSN string `yaml:"databaseDSN" required:"true"`
	Nats        struct {
		Host    string `yaml:"host" required:"true"`
		Port    string `yaml:"port" required:"true"`
		Cluster string `yaml:"cluster" required:"true"`
		Client  string `yaml:"client" required:"true"`
		Topic   string `yaml:"topic" required:"true"`
	} `required:"true" yaml:"nats"`
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		Files: []string{"./config.yml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yml": aconfigyaml.New(),
		},
	})

	if err := loader.Load(); err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		return nil, err
	}

	return &cfg, nil
}
