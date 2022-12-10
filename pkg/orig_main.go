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

func ORIGMain() {

	ch := rabbitmq.OpenConnectionAndChannel(Properties()["rabbitmq"])
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var exchanges = Properties()["exchanges"]
	var condition, _ = strconv.ParseBool(Properties()["queueCreated"])
	if !condition {
		rabbitmq.CreateExchange(exchanges, ch)
		queue1 := rabbitmq.CreateQueue(Properties()["fanoutqueue1"], ch)
		queue2 := rabbitmq.CreateQueue(Properties()["fanoutqueue2"], ch)
		rabbitmq.ExchangeBindingToQueue(ch, exchanges, queue1.Name)
		rabbitmq.ExchangeBindingToQueue(ch, exchanges, queue2.Name)
		log.Print(queue1.Name)
		log.Print(queue2.Name)
	}

	for i := range []int{1, 2, 3} {
		body := "MSG_" + strconv.Itoa(i)
		rabbitmq.PublicEventExchange(ch, ctx, exchanges, body)
		log.Printf(" [x] Sent %s", body)
		log.Print("sleep")
		time.Sleep(3 * time.Second)
	}

}

func HttpServerORIGServer() {

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
