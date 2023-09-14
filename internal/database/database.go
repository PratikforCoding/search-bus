package db

import (
	"context"
	"os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongo() (*mongo.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	mongoUri := os.Getenv("mongo_link")
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

