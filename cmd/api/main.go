package main

import (
	"log"

	"github.com/docker/docker/client"
	"github.com/fleimkeipa/kondukto/cmd/handlers"
	"github.com/fleimkeipa/kondukto/utils"
	"github.com/labstack/echo"
)

func main() {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	mongo, err := utils.Connect()
	if err != nil {
		panic(err)
	}

	receiver := handlers.Receiver{
		Docker: docker,
		Mongo:  mongo,
	}

	e := echo.New()
	e.Use(utils.Logger())
	e.POST("/newscan", receiver.NewScan)

	log.Fatal(e.Start(":8080"))
}
