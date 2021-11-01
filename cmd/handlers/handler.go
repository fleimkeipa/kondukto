package handlers

import (
	"net/http"

	"github.com/fleimkeipa/kondukto/pkg"
	"github.com/labstack/echo"
)

type Repo struct {
	Url string `json:"url"`
}

func Handler(c echo.Context) error {
	repo := Repo{}

	err := c.Bind(&repo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := pkg.ScanFunc(repo.Url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, id)
}
