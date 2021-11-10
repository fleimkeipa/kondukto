package handlers

import (
	"fmt"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/fleimkeipa/kondukto/utils"
	"github.com/labstack/echo"
)

type Receiver struct {
	Cli *client.Client
}

func (r *Receiver) Handler(c echo.Context) error {
	repo := struct {
		Url string `json:"url"`
	}{}
	fmt.Println("url", repo.Url)
	err := c.Bind(&repo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// id, err := pkg.ScanFunc(repo.Url)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, err.Error())
	// }

	if err := utils.ImageBuild(r.Cli); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := utils.RunContainer(r.Cli); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "id")
}
