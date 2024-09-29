package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "не удалось подключиться к RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "не удалось открыть канал")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"привет",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "не удалось объявить queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnError(err, "не удалось опубликовать message")
	log.Printf("[x] Sent %s\n", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
