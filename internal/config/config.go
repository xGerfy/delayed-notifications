package config

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type AppConfig struct {
	RedisConfig  RedisConfig
	RabbitConfig RabbitConfig
	ServerConfig ServerConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	Db       int
}

type RabbitConfig struct {
	Url      string
	Exchange string
	Queue    string
}

type ServerConfig struct {
	Port string
}

func New() (*AppConfig, error) {
	envFilePath := "./.env"

	cfg := config.New()

	if err := cfg.LoadEnvFiles(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to load env files: %w", err)
	}

	cfg.EnableEnv("")

	return &AppConfig{
		RedisConfig: RedisConfig{
			Addr:     cfg.GetString("redis.addr"),
			Password: cfg.GetString("redis.password"),
			Db:       cfg.GetInt("redis.db"),
		},
		RabbitConfig: RabbitConfig{
			Url:      cfg.GetString("rabbit.url"),
			Exchange: cfg.GetString("rabbit.exchange"),
			Queue:    cfg.GetString("rabbit.queue"),
		},
		ServerConfig: ServerConfig{
			Port: cfg.GetString("server.port"),
		},
	}, nil
}
