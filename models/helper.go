package models

import "go.mongodb.org/mongo-driver/mongo"

func getCollection(client *mongo.Client, collection string) *mongo.Collection {
	return client.Database(Database).Collection(collection)
}
