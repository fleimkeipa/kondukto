package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() (*mongo.Client, error) {
	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(
		context.TODO(),
		// options.Client().ApplyURI("mongodb+srv://"+os.Getenv("DBUSERNAME")+":"+os.Getenv("DBPASSWORD")+"@cluster0.4lioy.mongodb.net/eight-sup?retryWrites=true&w=majority"),
		options.Client().ApplyURI("mongodb://127.0.0.1:27017"),
	)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	fmt.Println("connection is ready")
	return client, nil
}

func InsertDB(client *mongo.Client) (map[string]interface{}, error) {
	// Open our jsonFile
	jsonFile, err := os.Open("tmp/src/result.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully Opened result.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	coll := client.Database("kondukto").Collection("results")

	insert, err := coll.InsertOne(context.Background(), result)
	if err != nil {
		return nil, err
	}
	result["_id"] = insert.InsertedID

	return result, nil
}
