package database

import (
	"assessment/config"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variables for the MongoDB client and database instances.
var (
	client   *mongo.Client
	database *mongo.Database
)

// Connect establishes a connection to the MongoDB server.
func Connect() error {
	// Load configuration settings.
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg.Url)

	// Set up client options with the MongoDB URI.
	clientOptions := options.Client().ApplyURI(cfg.Url)

	// Attempt to connect to the MongoDB server.
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Ping the MongoDB server to ensure it's responsive.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error pinging MongoDB: %v", err)
	}

	// Select the database to use.
	database = client.Database("organization_db")

	fmt.Println("Connected to MongoDB!")

	return nil
}

// GetClient retrieves the global MongoDB client instance.
func GetClient() *mongo.Client {
	return client
}

// GetDatabase retrieves the global MongoDB database instance.
func GetDatabase() *mongo.Database {
	return database
}
