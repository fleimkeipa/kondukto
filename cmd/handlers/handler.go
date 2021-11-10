package handlers

import (
	"fmt"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/fleimkeipa/kondukto/pkg"
	"github.com/fleimkeipa/kondukto/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type Receiver struct {
	Cli *client.Client
}

func (r *Receiver) Handler(c echo.Context) error {
	context := struct {
		Url string `json:"url"`
	}{}
	fmt.Println("url", context.Url)

	err := c.Bind(&context)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	repo, err := pkg.ScanFunc(context.Url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := utils.RunContainer(r.Cli, repo); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"scan_id": uuid.String()})
}
