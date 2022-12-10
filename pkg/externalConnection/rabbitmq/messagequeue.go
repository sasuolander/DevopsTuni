package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"mina.fi/devopstuni/pkg"
)

func PublicEvent(ch *amqp.Channel, ctx context.Context, name string, body string) {
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
	pkg.FailOnError(err, "Failed to publish a message")
}

func PublicEventExchange(ch *amqp.Channel, ctx context.Context, exchangeName string, body string) {
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
	pkg.FailOnError(err, "Failed to publish a message")
}

func ExchangeBindingToQueue(ch *amqp.Channel, exchange string, queue string) {
	err := ch.QueueBind(
		queue,    // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil,
	)
	pkg.FailOnError(err, "Failed to bind a queue")
}

func SetQos(ch *amqp.Channel) *amqp.Channel {
	err := ch.Qos(1, 0, true)
	if err != nil {
		pkg.FailOnError(err, "Failed to set QoS")
		return nil
	}
	return ch
}

func ConsumeEvent(ch *amqp.Channel, queueName string) string {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	pkg.FailOnError(err, "Failed to register a consumer or receive message")

	var message string

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		message = string(d.Body)
		return message
	}
	return "empty"
}

func OpenConnectionAndChannel(connectionStringPrm string) *amqp.Channel {
	var connectionString = fmt.Sprintf("amqp://%s", connectionStringPrm)
	log.Print(connectionString)
	conn, err := amqp.Dial(connectionString)
	pkg.FailOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	pkg.FailOnError(err, "Failed to open a channel")
	return ch
}

func CreateQueue(name string, ch *amqp.Channel) amqp.Queue {
	queue, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	pkg.FailOnError(err, "Failed to declare a queue")
	return queue
}

func CreateExchange(name string, ch *amqp.Channel) string {
	err := ch.ExchangeDeclare(
		name,     // name
		"fanout", // fanout
		false,    // delete when unused
		false,
		false, // exclusive
		true,  // no-wait
		nil,   // arguments
	)
	pkg.FailOnError(err, "Failed to declare a exchange")

	return name
}
