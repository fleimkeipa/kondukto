package main

import (
	"log"

	"github.com/docker/docker/client"
	"github.com/fleimkeipa/kondukto/cmd/handlers"
	"github.com/fleimkeipa/kondukto/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(utils.Logger())

	e.POST("/newscan", receiver.NewScan)
	e.GET("/scan/:scan_id", receiver.GetScan)
	log.Fatal(e.Start(":8080"))
}
