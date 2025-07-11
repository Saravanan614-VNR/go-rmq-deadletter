package main

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare DLX and DLQ (same as before)
	err = ch.ExchangeDeclare("dead_letter_exchange", "direct", true, false, false, false, nil)
	failOnError(err, "Failed to declare DLX")

	_, err = ch.QueueDeclare("dead_letter_queue", true, false, false, false, nil)
	failOnError(err, "Failed to declare DLQ")

	err = ch.QueueBind("dead_letter_queue", "dlq", "dead_letter_exchange", false, nil)
	failOnError(err, "Failed to bind DLQ")

	// Declare main queue with DLX (no TTL here!)
	args := amqp091.Table{
		"x-dead-letter-exchange":    "dead_letter_exchange",
		"x-dead-letter-routing-key": "dlq",
	}
	_, err = ch.QueueDeclare("main_queue", true, false, false, false, args)
	failOnError(err, "Failed to declare main queue")

	// Publish with message-level TTL (5 seconds)
	body := "Message with 5-second TTL"
	err = ch.Publish(
		"",           // default exchange
		"main_queue", // routing key
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Expiration:  "5000", // in milliseconds (5 seconds)
		})
	failOnError(err, "Failed to publish")

	log.Printf("Sent: %s", body)
}
