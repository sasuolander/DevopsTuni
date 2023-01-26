package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"mina.fi/devopstuni/pkg/utils"
)

func PublicEvent(ch *amqp.Channel, ctx context.Context, name string, body string) {
	var err error
	err = ch.PublishWithContext(ctx,
		"",
		name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	utils.FailOnError(err, "Failed to publish a message")
}

func PublicEventExchange(ch *amqp.Channel, ctx context.Context, exchangeName string, body string) {
	var err error
	err = ch.PublishWithContext(ctx,
		exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	utils.FailOnError(err, "Failed to publish a message")
}

func ExchangeBindingToQueue(ch *amqp.Channel, exchange string, queue string) {
	err := ch.QueueBind(
		queue,
		"",
		exchange,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to bind a queue")
}

func SetQos(ch *amqp.Channel) *amqp.Channel {
	err := ch.Qos(1, 0, true)
	utils.FailOnError(err, "Failed to set QoS")
	return ch
}

func ConsumeEvent(ch *amqp.Channel, queueName string) string {
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to register a consumer or receive message")

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
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	return ch
}

func CreateQueue(name string, ch *amqp.Channel) amqp.Queue {
	queue, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare a queue")
	return queue
}

func CreateExchange(name string, ch *amqp.Channel) string {
	err := ch.ExchangeDeclare(
		name,
		"fanout",
		false,
		false,
		false,
		true,
		nil,
	)
	utils.FailOnError(err, "Failed to declare a exchange")

	return name
}
