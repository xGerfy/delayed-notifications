package main

import (
	"context"
	"delayed/internal/config"
	"delayed/internal/rabbit/consumer"
	"delayed/internal/rabbit/publisher"
	"delayed/internal/repo"
	"delayed/internal/router"
	"delayed/internal/router/handler"
	"delayed/internal/service"
	"fmt"
	"sync"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
	}

	repo := repo.New(cfg.RedisConfig.Addr, cfg.RedisConfig.Password, cfg.RedisConfig.Db)

	pub, err := publisher.New(cfg)
	if err != nil {
		fmt.Println(err)
	}
	defer pub.Close()

	svc := service.New(repo, pub)

	handlers := handler.New(svc)
	router := router.New(handlers)

	cons, err := consumer.New(cfg, svc)
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Go(func() {
		if err := cons.Start(ctx); err != nil {
			fmt.Println(err)
		}
	})

	if err := router.Run(cfg.ServerConfig.Port); err != nil {
		fmt.Println(err)
		cancel()
	}

	go func() {
		wg.Wait()
	}()
}
