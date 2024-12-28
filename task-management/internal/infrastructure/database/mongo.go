package database

import (
	"context"
	"log"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoClient(config *config.Config, ctx context.Context) *mongo.Client {
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatalf("❌ Error connecting to MongoDB: %v\n", err)

		return nil
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Error pinging MongoDB: %v\n", err)

		return nil
	}

	log.Println("✅ Connected to MongoDB")

	return mongoClient
}
