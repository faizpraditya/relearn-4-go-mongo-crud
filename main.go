package main

import (
	"context"
	"fmt"
	"go-mongo-crud/controllers"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(getMongoClient())

	r.GET("/user", uc.GetAllUsers)
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)

	// create golang server
	http.ListenAndServe("localhost:9000", r)
}

func getMongoClient() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetConnectTimeout(10 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		return nil
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Printf("Failed to ping MongoDB: %v\n", err)
		return nil
	}

	fmt.Println("Connected to MongoDB!")

	return client
}
