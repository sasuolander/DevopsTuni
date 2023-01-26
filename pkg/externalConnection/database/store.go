package database

import (
	"bufio"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"mina.fi/devopstuni/pkg/utils"
	"os"
)

var basePath = "/temp"

func Connection(url string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	utils.FailOnError(err, "Failed to connect")

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
	utils.FailOnError(err, "Failed to insert")
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
		var elem logMessage
		err := cursor.Decode(&elem)
		utils.FailOnError(err, "Failed to decode")
		results = append(results, elem)
	}

	return results
}
func ConvertToStringArray(result []logMessage) []string {
	var stringArray []string
	for _, message := range result {
		stringArray = append(stringArray, message.Message+" \n")
	}

	return stringArray
}

func fileWrite(message string, file string) {
	data := []byte(message)
	err := os.WriteFile(basePath+file+".txt", data, 0)
	if err != nil {
		utils.FailOnError(err, "Failed to write")
	}

	fmt.Println("done")
}

func fileReader(file string) []logMessage {
	f, err := os.Open(basePath + file + ".txt")

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	var results []logMessage
	for scanner.Scan() {
		results = append(results, logMessage{Message: scanner.Text()})
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return results
}
