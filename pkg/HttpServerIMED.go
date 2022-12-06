package pkg

import (
	"context"
	"log"
	"time"
)

func HttpServerIMED(queueName string, queueListen string) {
	ch := openConnectionAndChannel("guest:guest@localhost:5672")

	queue := createQueue(queueName, ch)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for true {
		item1 := consumeEvent(ch, queueListen)
		time.Sleep(1 * time.Second)
		var message = "Got " + item1
		log.Print(message)
		publicEvent(ch, ctx, queue, message)
		log.Print("Wait")
	}
}
