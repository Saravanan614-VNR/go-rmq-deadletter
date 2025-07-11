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
	failOnError(err, "Connect failed")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Channel failed")
	defer ch.Close()

	// Declare DLX and DLQ
	err = ch.ExchangeDeclare("dead_letter_exchange", "direct", true, false, false, false, nil)
	failOnError(err, "DLX declare failed")

	_, err = ch.QueueDeclare("dead_letter_queue", true, false, false, false, nil)
	failOnError(err, "DLQ declare failed")

	err = ch.QueueBind("dead_letter_queue", "dlq", "dead_letter_exchange", false, nil)
	failOnError(err, "DLQ bind failed")

	// Declare main queue with DLX (NO consumer for it!)
	args := amqp091.Table{
		"x-dead-letter-exchange":    "dead_letter_exchange",
		"x-dead-letter-routing-key": "dlq",
	}
	_, err = ch.QueueDeclare("main_queue", true, false, false, false, args)
	failOnError(err, "Main queue declare failed")

	// Only consume from DLQ
	log.Println("Listening to DLQ...")
	dlqMsgs, err := ch.Consume("dead_letter_queue", "", true, false, false, false, nil)
	failOnError(err, "DLQ consume failed")

	for d := range dlqMsgs {
		log.Printf("[DLQ] Received expired message: %s", d.Body)
	}
}
