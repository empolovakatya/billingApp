package main

import (
	billing "billingApp"
	"billingApp/pkg/repository"
	"billingApp/workers"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"math/rand"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	workers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	workers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	workers.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	workers.FailOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	data := billing.WorkerSender{Method: "increase", BalanceId: 95, Amount: 45}
	body, err := json.Marshal(data)
	if err != nil {
		log := repository.ErrorResponse("Error encoding JSON %s", err)
		logrus.Print(string(log))
	}

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          body,
		})
	workers.FailOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			logrus.Print(string(d.Body))
			break
		}
	}
	return
}
