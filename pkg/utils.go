package pkg

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"strings"
)

const uri = "mongodb://user:pass@sample.host:27017/?maxPoolSize=20&w=majority"

func databaseConnection(url string) *mongo.Client {
	// Create a new client and connect to the server
	log.Print(url)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		failOnError(err, "Failed to connect")
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	return client
}

func insertItem(client *mongo.Client, database string, message string, collectionName string) {
	coll := client.Database(database).Collection(collectionName)
	doc := bson.D{{"message", message}}
	result, err := coll.InsertOne(context.TODO(), doc)
	failOnError(err, "Failed to insert")
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}

type logMessage struct {
	Message string `json:"message"`
}

func getItems(client *mongo.Client, database string, collectionName string) []logMessage {
	coll := client.Database(database).Collection(collectionName)
	opts := options.Find()
	cursor, err := coll.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		panic(err)
	}
	var results []logMessage
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem logMessage
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)

	}

	return results
}

func getRoot(w http.ResponseWriter, r *http.Request) {

	var client = databaseConnection("mongodb://localhost:27017/test")
	var collectionName = "observLogCollection"
	var result = getItems(client, "test", collectionName)
	var stringArray []string

	for _, message := range result {
		stringArray = append(stringArray, message.Message+" \n")
	}

	fmt.Print(result)
	stringByte := strings.Join(stringArray, "\x20\x00")
	w.Header().Add("Content-Type", "text/plain;charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(stringByte))
	if err != nil {
		fmt.Printf("error in server: %s\n", err)
		return
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			failOnError(err, "Failed to disconnect")
			panic(err)
		}
	}()
	return
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

func setQos(ch *amqp.Channel) *amqp.Channel {
	err := ch.Qos(1, 0, true)
	if err != nil {
		failOnError(err, "Failed to set QoS")
		return nil
	}
	return ch
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
