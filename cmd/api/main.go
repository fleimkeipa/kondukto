package main

import (
	"os"

	"github.com/docker/docker/client"
	"github.com/fleimkeipa/kondukto/cmd/handlers"
	"github.com/fleimkeipa/kondukto/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// fmt.Println("docker_host", os.Getenv("DOCKER_HOST"))
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

	e.POST("/newscan", receiver.NewScan)
	e.GET("/scan/:scan_id", receiver.GetScan)
	e.GET("/scans", receiver.GetScanAll)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
