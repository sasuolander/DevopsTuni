package pkg

import (
	"context"
	"log"
	"strconv"
	"time"
)

func HttpServerORIG(queueName string) {

	ch := openConnectionAndChannel("guest:guest@localhost:5672")
	queue := createQueue(queueName, ch)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := range []int{1, 2, 3} {
		body := "MSG_" + strconv.Itoa(i)
		publicEvent(ch, ctx, queue, body)
		log.Printf(" [x] Sent %s", body)
		log.Print("sleep")
		time.Sleep(3 * time.Second)

	}
}
