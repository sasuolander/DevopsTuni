package pkg

import (
	"context"
	"log"
	"time"
)

func HttpServerIMED(queueName string, queueListen string, queueCreated bool, connectionString string) {
	ch := openConnectionAndChannel(connectionString)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !queueCreated {
		createQueue(queueName, ch)
	}
	for true {
		item1 := consumeEvent(ch, queueListen)
		time.Sleep(1 * time.Second)
		var message = "Got " + item1
		log.Print(queueListen)
		log.Print(message)
		log.Print("Sending to " + queueName)
		publicEvent(ch, ctx, queueName, message)
		log.Print("Wait")
	}
}
