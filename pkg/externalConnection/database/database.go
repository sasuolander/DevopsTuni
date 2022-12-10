package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"mina.fi/devopstuni/pkg"
)

func Connection(url string) *mongo.Client {
	// Create a new client and connect to the server
	log.Print(url)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		pkg.FailOnError(err, "Failed to connect")
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	return client
}

func InsertItem(client *mongo.Client, database string, message string, collectionName string) {
	coll := client.Database(database).Collection(collectionName)
	doc := bson.D{{"message", message}}
	result, err := coll.InsertOne(context.TODO(), doc)
	pkg.FailOnError(err, "Failed to insert")
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}

type logMessage struct {
	Message string `json:"message"`
}

func GetItems(client *mongo.Client, database string, collectionName string) []logMessage {
	coll := client.Database(database).Collection(collectionName)
	opts := options.Find()
	cursor, err := coll.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		panic(err)
	}
	var results []logMessage
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem logMessage
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)

	}

	return results
}
