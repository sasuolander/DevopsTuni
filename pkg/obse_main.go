package pkg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mina.fi/devopstuni/pkg/externalConnection/database"
	"mina.fi/devopstuni/pkg/externalConnection/rabbitmq"
	"net/http"
	"os"
	"strconv"
	"time"
)

func OBSEMain() {
	var counter = 0

	ch1 := rabbitmq.SetQos(rabbitmq.OpenConnectionAndChannel(Properties()["rabbitmq"]))
	ch2 := rabbitmq.SetQos(rabbitmq.OpenConnectionAndChannel(Properties()["rabbitmq"]))
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var client = database.Connection(Properties()["mongoDbURL"])
	var collectionName = "observLogCollection"
	var queueName1 = Properties()["fanoutqueue1"]
	var queueName2 = Properties()["mainqueue2"]
	for true {
		var item1 = rabbitmq.ConsumeEvent(ch1, queueName1)
		counter += 1
		var count1 = counter
		message1 := fmt.Sprintf("%s %s %s from %s", time.Now().Format(time.RFC3339), strconv.Itoa(count1), item1, queueName1)
		log.Print("insert item2")
		log.Print(queueName1)
		log.Print(message1)
		database.InsertItem(client, Properties()["dbName"], message1, collectionName)
		var item2 = rabbitmq.ConsumeEvent(ch2, queueName2)
		counter += 1
		var count2 = counter
		message2 := fmt.Sprintf("%s %s %s from %s", time.Now().Format(time.RFC3339), strconv.Itoa(count2), item2, queueName2)
		log.Print("insert item2")
		log.Print(queueName2)
		log.Print(message2)
		database.InsertItem(client, Properties()["dbName"], message2, collectionName)

	}
}

func HttpServerOBSEServer() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

	})
	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
