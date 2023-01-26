package pkg

import (
	"context"
	"fmt"
	"log"
	"mina.fi/devopstuni/pkg/externalConnection/rabbitmq"
	"mina.fi/devopstuni/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

func imedMain() {
	var queueName = utils.Properties.MAINQUEUE2
	var queueListen = utils.Properties.FANOUTQUEUE2
	var condition, _ = strconv.ParseBool(utils.Properties.QUEUECREATED)

	ch := rabbitmq.SetQos(rabbitmq.OpenConnectionAndChannel(utils.Properties.RABBITMQ))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !condition {
		rabbitmq.CreateQueue(queueName, ch)
	}

	for true {
		item1 := rabbitmq.ConsumeEvent(ch, queueListen)
		time.Sleep(1 * time.Second)
		var message = "Got " + item1
		log.Print("Sending to " + queueName + " message: " + message)
		rabbitmq.PublicEvent(ch, ctx, queueName, message)
		log.Print("Wait")
	}
}
func ServerStarterImed() {
	startTime := fmt.Sprintf("%s", time.Now().Format(time.RFC3339))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pigApi(w, r, startTime)
	})

	go imedMain()
	utils.StartServer(utils.Properties.PORT)
}
