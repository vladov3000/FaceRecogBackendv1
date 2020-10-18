package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Person is an exported struct that contains generic
// info for a person stored in our database
type Person struct {
	Name         string `json: "name, omitempty" bson: "name, omitempty"`
	Title        string `json: "title, omitempty" bson: "title, omitempty"`
	CustomFields bson.M `json: "custom, omitempty" bson.M: "custom, omitempty"`
}

func setupMongo() (context.Context, *mongo.Collection) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoKey))
	if err != nil {
		log.Printf("Failed to Get MongoDB Client: %s", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to Connect to MongoDB Context: %s", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("Failed to Ping MongoDB Client: %s", err)
	}

	peopleDatabase := client.Database("people")
	nameCollection := peopleDatabase.Collection("names")
	return ctx, nameCollection
}

func insertPerson(ctx context.Context, nameOfCollection *mongo.Collection, person Person) string {
	insertRes, err := nameOfCollection.InsertOne(ctx, person)
	if err != nil {
		log.Fatal(err)
	}

	result := insertRes.InsertedID.(primitive.ObjectID)

	log.Printf("Inserted new person with ObjectID %s into database", result.Hex())
	return result.Hex()
}

func queryPerson(ctx context.Context, nameOfCollection *mongo.Collection, itemID string) map[string]interface{} {
	if itemID == "" {
		return nil;
	}
	if len(itemID) > 24 {
		itemID = itemID[:25]
	}
	objectID, err := primitive.ObjectIDFromHex(itemID[:len(itemID)-1])
	if err != nil {
		log.Fatal(err)
	}
	filterCursor, err := nameOfCollection.Find(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Fatal(err)
	}
	var episodesFiltered []bson.M
	if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
		log.Fatal(err)
	}
	log.Printf("Found and returning info for data with ID %s", itemID)
	return episodesFiltered[0]
}