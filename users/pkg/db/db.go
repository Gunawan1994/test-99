package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	onceMongo      sync.Once
)

func NewConn() *mongo.Client {
	onceMongo.Do(func() {
		uri := fmt.Sprintf(
			"mongodb://%s:%s/%s",
			os.Getenv("MONGO_ADDR"),
			os.Getenv("MONGO_PORT"),
			os.Getenv("MONGO_DATABASE"),
		)

		clientOptions := options.Client().ApplyURI(uri)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %v", err))
		}

		if err := client.Ping(ctx, nil); err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v", err))
		}

		clientInstance = client
		fmt.Println("✅ Connected to MongoDB")
	})

	return clientInstance
}
