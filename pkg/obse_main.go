package pkg

import (
	"context"
	"fmt"
	"log"
	"mina.fi/devopstuni/pkg/externalConnection/database"
	"mina.fi/devopstuni/pkg/externalConnection/rabbitmq"
	"mina.fi/devopstuni/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

func obseMain() {
	var counter = 0

	var dbname = utils.Properties.DBNAME
	var rabbitmqP = utils.Properties.RABBITMQ
	var client = database.Connection(utils.Properties.MONGODBURL)
	var collectionName = utils.Properties.OBSERVLOGCOLLECTION
	var queueName1 = utils.Properties.FANOUTQUEUE1
	var queueName2 = utils.Properties.MAINQUEUE2

	ch1 := rabbitmq.SetQos(rabbitmq.OpenConnectionAndChannel(rabbitmqP))
	ch2 := rabbitmq.SetQos(rabbitmq.OpenConnectionAndChannel(rabbitmqP))
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for true {
		var item1 = rabbitmq.ConsumeEvent(ch1, queueName1)
		counter += 1
		message1 := fmt.Sprintf("%s %s %s from %s", time.Now().Format(time.RFC3339), strconv.Itoa(counter), item1, queueName1)
		log.Print("Got from queue " + queueName1 + " message: " + message1)
		database.InsertItem(client, dbname, message1, collectionName)

		var item2 = rabbitmq.ConsumeEvent(ch2, queueName2)
		counter += 1
		message2 := fmt.Sprintf("%s %s %s from %s", time.Now().Format(time.RFC3339), strconv.Itoa(counter), item2, queueName2)
		log.Print("Got from queue " + queueName2 + " message: " + message2)
		database.InsertItem(client, dbname, message2, collectionName)

	}
}

func ServerStarterObse() {
	startTime := fmt.Sprintf("%s", time.Now().Format(time.RFC3339))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pigApi(w, r, startTime)
	})

	go obseMain()
	utils.StartServer(utils.Properties.PORT)
}
