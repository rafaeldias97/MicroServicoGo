package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	// Official 'mongo-go-driver' packages

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuracao Global
var ctx context.Context
var db *mongo.Database

func connect(url string, dbName string) {
	clientOptions := options.Client().ApplyURI(url)
	fmt.Println("clientOptions TYPE:", reflect.TypeOf(clientOptions))

	// Connect to the MongoDB and return Client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}
	// Declare Context type object for managing multiple API requests
	ctx, _ = context.WithTimeout(context.Background(), 15*time.Second)
	db = client.Database(dbName)
}
func conectRabbit() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()
	}
	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		defer ch.Close()
	}

	return conn, ch
}
func QueueDeclare(queueName string, ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err == nil {
		failOnError(err, "Failed to declare a queue")
	}

	return q
}

func Consume(queueName string, ch *amqp.Channel) {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err == nil {
		failOnError(err, "Failed to register a consumer")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			createUser(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func main() {
	connect("mongodb://root:ContaDigital2020!@localhost:27017", "Cliente")
	_, ch := conectRabbit()
	QueueDeclare("CadastrarCliente", ch)
	Consume("CadastrarCliente", ch)
}

func createUser(msg []byte) {
	col := db.Collection("ClienteCollection")

	result, insertErr := col.InsertOne(ctx, bson.M{
		"nome":       "Rafael Dias",
		"idade":      22,
		"created_at": time.Now(),
	})

	if insertErr != nil {
		fmt.Println("InsertOne ERROR:", insertErr)
		os.Exit(1) // safely exit script on error
	}

	fmt.Println("InsertOne() API result:", result)
	// get the inserted ID string
	fmt.Println("InsertOne() newID:", result.InsertedID)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
