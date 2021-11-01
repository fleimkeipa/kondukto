package main

import (
	"log"

	"github.com/fleimkeipa/kondukto/cmd/handlers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.POST("/newscan", handlers.Handler)

	log.Fatal(e.Start(":8080"))
}
