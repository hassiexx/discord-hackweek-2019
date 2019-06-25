package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne inserts a single record into the database.
func InsertOne(collection Collection, record interface{}) error {
	_, err := database.Collection(string(collection)).InsertOne(context.Background(), record)
	return err
}

// QueryOne queries for a single record.
// The found record is marshalled into v.
func QueryOne(collection Collection, query interface{}, v interface{}) error {
	err := database.Collection(string(collection)).FindOne(context.Background(), query).Decode(v)
	return err
}

// UpsertOne upserts a record.
// Upsert in MongoDB is creating a new record if it does not exist -
// determined by the query - by merging the query fields and the update fields;
// otherwise performing an update.
func UpsertOne(collection Collection, query interface{}, update interface{}) error {
	upsert := true
	_, err := database.Collection(string(collection)).UpdateOne(context.Background(),
		query, update, &options.UpdateOptions{Upsert: &upsert})
	return err
}
