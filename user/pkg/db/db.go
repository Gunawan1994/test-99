package db

import (
	"context"
	"fmt"
	"os"
	"strings"
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
			"mongodb://%s:%s@%s:%s/%s?authSource=admin",
			os.Getenv("MONGO_USER"),     // misal: root
			os.Getenv("MONGO_PASSWORD"), // misal: 12345
			os.Getenv("MONGO_ADDR"),     // misal: mongo (nama service Docker)
			os.Getenv("MONGO_PORT"),     // misal: 27017
			os.Getenv("MONGO_DATABASE"), // misal: mydb
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

		db := client.Database(os.Getenv("MONGO_DATABASE"))
		ctxCreate, cancelCreate := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelCreate()

		err = db.CreateCollection(ctxCreate, "users")
		if err != nil && !IsNamespaceExists(err) {
			panic(fmt.Sprintf("failed to create 'users' collection: %v", err))
		} else if err == nil {
			fmt.Println("📁 Created collection: users")
		} else {
			fmt.Println("Collection 'users' already exists")
		}

		clientInstance = client
		fmt.Println("Connected to MongoDB")
	})

	return clientInstance
}

func IsNamespaceExists(err error) bool {
	return err != nil && strings.Contains(err.Error(), "NamespaceExists")
}
