// Localities API
//
// This is a sample localities API.
//
//	Schemes: http
//  Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Iulian Coretchi <iulian.coretchi@gmail.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"localities-api/handlers"
	"log"
	"os"
)

var authHandler *handlers.AuthHandler
var localitiesHandler *handlers.LocalitiesHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("localities")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Println(status)

	localitiesHandler = handlers.NewLocalitiesHandler(ctx, collection, redisClient)
	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func main() {
	router := gin.Default()
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.POST("/localities", localitiesHandler.NewLocalityHandler)
	router.GET("/localities", localitiesHandler.ListLocalitiesHandler)
	router.PUT("/localities/:code", localitiesHandler.UpdateLocalityHandler)
	router.DELETE("/localities/:code", localitiesHandler.DeleteLocalityHandler)
	router.GET("/localities/:code", localitiesHandler.GetOneLocalityHandler)
	router.Run()
}
