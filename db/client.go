package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var database *mongo.Database

// Connect attempts to connect to the database.
func Connect() error {
	// Create client.
	cl, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DISCORD_HACK_WEEK_2019_MONGODB")))
	if err != nil {
		return err
	}

	// Connect.
	err = client.Connect(context.Background())
	if err != nil {
		return err
	}

	// Store client and database.
	client = cl
	database = client.Database("dismod")

	return nil
}

// Disconnect disconnects from the database.
func Disconnect() {
	_ = client.Disconnect(context.Background())
}
