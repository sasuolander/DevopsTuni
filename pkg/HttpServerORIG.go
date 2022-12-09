package pkg

import (
	"context"
	"log"
	"strconv"
	"time"
)

func HttpServerORIG(queueName string, queueCreated bool, connectionString string) {

	ch := openConnectionAndChannel(connectionString)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var exchanges = queueName + "-ex"

	if !queueCreated {
		createExchange(exchanges, ch)
		queue1 := createQueue(queueName+"-1", ch)
		queue2 := createQueue(queueName+"-2", ch)
		exchangeBindingToQueue(ch, exchanges, queue1.Name)
		exchangeBindingToQueue(ch, exchanges, queue2.Name)
		log.Print(queue1.Name)
		log.Print(queue2.Name)
	}

	for i := range []int{1, 2, 3} {
		body := "MSG_" + strconv.Itoa(i)
		publicEventExchange(ch, ctx, exchanges, body)
		log.Printf(" [x] Sent %s", body)
		log.Print("sleep")
		time.Sleep(3 * time.Second)
	}

}
