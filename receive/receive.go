package main

import (
	"log"

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

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "не удалось зарегистрировать consumer")
	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("получено сообщение: %s", d.Body)
		}
	}()
	log.Printf("[*] Ожидание сообщений. Для выхода нажмите CTRL + C ")
	<-forever

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
