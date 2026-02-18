package main

import (
	"delayed/internal/rabbit/publisher"
	"delayed/internal/repo"
	"delayed/internal/router"
	"delayed/internal/router/handler"
	"delayed/internal/service"
)

func main() {
	repo := repo.New("localhost:6379", "", 0)
	pub, _ := publisher.New("amqp://guest:guest@localhost:5672/")
	svc := service.New(repo, pub)
	handlers := handler.New(svc)
	router := router.New(handlers)

	router.Run(":8080")
}
