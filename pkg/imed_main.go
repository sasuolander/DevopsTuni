package pkg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mina.fi/devopstuni/pkg/externalConnection/rabbitmq"
	"net/http"
	"os"
	"strconv"
	"time"
)

func IMEDMain() {
	ch := rabbitmq.SetQos(rabbitmq.OpenConnectionAndChannel(Properties()["rabbitmq"]))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var queueName = Properties()["mainqueue2"]
	var queueListen = Properties()["fanoutqueue2"]
	log.Print(queueListen)
	log.Print(queueName)
	var condition, _ = strconv.ParseBool(Properties()["queueCreated"])
	if !condition {
		rabbitmq.CreateQueue(queueName, ch)
	}
	for true {
		item1 := rabbitmq.ConsumeEvent(ch, queueListen)
		time.Sleep(1 * time.Second)
		var message = "Got " + item1
		log.Print(queueListen)
		log.Print(message)
		log.Print("Sending to " + queueName)
		rabbitmq.PublicEvent(ch, ctx, queueName, message)
		log.Print("Wait")
	}
}

func HttpServerIMEDServer() {

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
