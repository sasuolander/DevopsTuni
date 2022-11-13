package main

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func publicEvent(ch *amqp.Channel, ctx context.Context, queue amqp.Queue, body string) {
	var err error
	err = ch.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

func consumeEvent(ch *amqp.Channel, queueName string) string {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer or receive message")

	var message string

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		message = string(d.Body)
		return message
	}
	return "empty"
}

func envOrString(env string, parameter string) string {
	if len(env) != 0 {
		return env
	} else {
		return parameter
	}
}

// guest:guest@localhost:5672
func openConnectionAndChannel(connectionStringPrm string) *amqp.Channel {
	var connectionString = fmt.Sprintf("amqp://%s/", envOrString(os.Getenv("connection"), connectionStringPrm))
	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

func createQueue(name string, ch *amqp.Channel) amqp.Queue {
	queue, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return queue
}

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

func writer(filename string, message string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(f, message)
	if err != nil {
		fmt.Println(err)
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("file appended successfully")
}

func HttpServerOBSE(queueName1 string, queueName2 string) {
	var counter = 0

	ch1 := openConnectionAndChannel("guest:guest@localhost:5672")
	ch2 := openConnectionAndChannel("guest:guest@localhost:5672")
	var filename = "/tmp/log.txt"
	_, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for true {
		var item1 = consumeEvent(ch1, queueName1)
		counter += 1
		var count1 = counter
		message1 := fmt.Sprintf("%s %s %s to %s", time.Now().Format(time.RFC3339), strconv.Itoa(count1), item1, "compse140.o")
		log.Print(message1)
		writer(filename, message1)
		var item2 = consumeEvent(ch2, queueName2)
		counter += 1
		var count2 = counter
		message2 := fmt.Sprintf("%s %s %s to %s", time.Now().Format(time.RFC3339), strconv.Itoa(count2), item2, "compse140.i")
		writer(filename, message1)
		log.Print(message2)
	}
}

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

func getRoot(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := os.ReadFile("/tmp/log.txt")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write(fileBytes)
	if err != nil {
		fmt.Printf("error in server: %s\n", err)
		return
	}
	return
}

func HttpServ() {

	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("start")
	fmt.Println(os.Args[1:])

	switch os.Args[1:][0] {
	case "HttpServ":
		HttpServ()
	case "HttpServerORIG":
		HttpServerORIG("compse140.o")
	case "HttpServerOBSE":
		HttpServerOBSE("compse140.o", "compse140.i")
	case "HttpServerIMED":
		HttpServerIMED("compse140.i", "compse140.o")
	default:
		fmt.Println("unknown mode")
		errors.New("unknown mode")
	}

}
