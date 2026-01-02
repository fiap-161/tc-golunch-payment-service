package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	client *mongo.Client
}

func NewMongoDatabase() *MongoDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	log.Println("Connected to MongoDB successfully")

	return &MongoDatabase{
		client: client,
	}
}

func (m *MongoDatabase) GetClient() *mongo.Client {
	return m.client
}

func (m *MongoDatabase) GetDatabase() *mongo.Database {
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "golunch_payments"
	}
	return m.client.Database(dbName)
}


