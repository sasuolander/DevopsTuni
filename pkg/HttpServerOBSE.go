package pkg

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

func HttpServerOBSE(queueName1 string, queueName2 string, connectionString string, databaseString string, database string) {
	var counter = 0

	ch1 := setQos(openConnectionAndChannel(connectionString))
	ch2 := setQos(openConnectionAndChannel(connectionString))
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var client = databaseConnection(databaseString)
	var collectionName = "observLogCollection"

	for true {
		var item1 = consumeEvent(ch1, queueName1)
		counter += 1
		var count1 = counter
		message1 := fmt.Sprintf("%s %s %s from %s", time.Now().Format(time.RFC3339), strconv.Itoa(count1), item1, queueName1)
		log.Print(queueName1)
		log.Print(message1)
		insertItem(client, database, message1, collectionName)
		var item2 = consumeEvent(ch2, queueName2)
		counter += 1
		var count2 = counter
		message2 := fmt.Sprintf("%s %s %s from %s", time.Now().Format(time.RFC3339), strconv.Itoa(count2), item2, queueName2)
		insertItem(client, database, message2, collectionName)
		log.Print(queueName2)
		log.Print(message2)
	}
}
