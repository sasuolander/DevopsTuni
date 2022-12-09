package pkg

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
)

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

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func publicEvent(ch *amqp.Channel, ctx context.Context, name string, body string) {
	var err error
	err = ch.PublishWithContext(ctx,
		"",    // exchange
		name,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

func publicEventExchange(ch *amqp.Channel, ctx context.Context, exchangeName string, body string) {
	var err error
	err = ch.PublishWithContext(ctx,
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

func exchangeBindingToQueue(ch *amqp.Channel, exchange string, queue string) {
	err := ch.QueueBind(
		queue,    // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
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
	var connectionString = fmt.Sprintf("amqp://%s", envOrString(os.Getenv("connection"), connectionStringPrm))
	log.Print(connectionString)
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

func createExchange(name string, ch *amqp.Channel) string {
	err := ch.ExchangeDeclare(
		name,     // name
		"fanout", // fanout
		false,    // delete when unused
		false,
		false, // exclusive
		true,  // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a exchange")

	return name
}
