package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"github.com/fleimkeipa/kondukto/pkg"
	"github.com/fleimkeipa/kondukto/utils"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Receiver struct {
	Docker *client.Client
	Mongo  *mongo.Client
}

func (r *Receiver) NewScan(c echo.Context) error {
	start := time.Now()
	context := struct {
		Url string `json:"url"`
	}{}
	err := c.Bind(&context)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	repo, err := pkg.ScanFunc(context.Url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println("scanFunc", time.Since(start))

	if err := utils.RunContainer(r.Docker, repo); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println("runContainer", time.Since(start))

	result, err := utils.InsertDB(r.Mongo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println("insertDB", time.Since(start))

	return c.JSON(http.StatusOK, map[string]interface{}{"scan_id": result["_id"].(primitive.ObjectID)})
}

func (r *Receiver) GetScan(c echo.Context) error {
	id, _ := primitive.ObjectIDFromHex(c.Param("scan_id"))

	coll := r.Mongo.Database("kondukto").Collection("results")

	var result map[string]interface{}
	err := coll.FindOne(context.Background(), map[string]interface{}{
		"_id": id,
	}).Decode(&result)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
