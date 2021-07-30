package main

import (
	billing "billingApp"
	"billingApp/workers"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	workers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	workers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"worker_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	workers.FailOnError(err, "Failed to declare a queue")

	data := billing.WorkerSender{Method: "approve", FreezeId: 1, IsApproved: true}
	body, err := json.Marshal(data)
	workers.FailOnError(err, "Error encoding JSON")
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{ContentType: "application/json",
			Body: []byte(body),
		})
	workers.FailOnError(err, "Failed to publish a message")
	logrus.Printf(" [x] Sent %s", body)

}
