package main

import (
	"log"

	"github.com/mohamadafzal06/depository/handler"
	"github.com/mohamadafzal06/depository/repository/postgres"
	"github.com/mohamadafzal06/depository/service"
)

func main() {
	repo, err := postgres.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	if err := repo.Init(); err != nil {
		log.Fatal(err)
	}

	service := service.NewDepository(repo)

	handler := handler.New(":8999", service)

	handler.Run()
}
