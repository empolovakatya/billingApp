package workers

import (
	billing "billingApp"
	"billingApp/pkg/repository"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func Work() {
	//conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"worker_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			logrus.Printf(" [*] %s", d.Body)
			data := &billing.Worker{}
			err := json.Unmarshal(d.Body, data)
			if err != nil {
				response := billing.Errors{ErrMessage: fmt.Sprintf("Error decoding JSON %s", err)}
				dataResponse := billing.ErrorResponse{Data: response}
				body, _ := json.Marshal(dataResponse)
				logrus.Print(string(body))
				//d.Ack(false)
			} else {
				if data.Method == "decrease" && data.BalanceId > 0 && data.Amount > 0 {
					msg, err := repository.Decrease(*data)
					if err != nil {
						response := billing.Errors{ErrMessage: fmt.Sprintf("Failed on decrease amount %s", err)}
						dataResponse := billing.ErrorResponse{Data: response}
						body, _ := json.Marshal(dataResponse)
						logrus.Print(string(body))
						return
					}
					logrus.Print(string(msg))
					d.Ack(false)
				} else if data.Method == "increase" && data.BalanceId > 0 && data.Amount > 0 {
					msg, err := repository.Increase(*data)
					if err != nil {
						response := billing.Errors{ErrMessage: fmt.Sprintf("Failed on decrease amount %s", err)}
						dataResponse := billing.ErrorResponse{Data: response}
						body, _ := json.Marshal(dataResponse)
						logrus.Print(string(body))
						return
					}
					logrus.Print(string(msg))
					d.Ack(false)
				} else if data.Method == "send" && data.BalanceId > 0 && data.Amount > 0 && data.Receiver > 0 {
					msg, err := repository.SendToOther(*data)
					if err != nil {
						response := billing.Errors{ErrMessage: fmt.Sprintf("Failed on send amount %s", err)}
						dataResponse := billing.ErrorResponse{Data: response}
						body, _ := json.Marshal(dataResponse)
						logrus.Print(string(body))
						return
					}
					logrus.Print(string(msg))
					d.Ack(false)
				} else if data.Method == "freeze" && data.BalanceId > 0 && data.FreezedAmount > 0 {
					msg, err := repository.FreezeAmount(*data)
					if err != nil {
						response := billing.Errors{ErrMessage: fmt.Sprintf("Failed on freeze amount %s", err)}
						dataResponse := billing.ErrorResponse{Data: response}
						body, _ := json.Marshal(dataResponse)
						logrus.Print(string(body))
						return
					}
					logrus.Print(string(msg))
					d.Ack(false)
				} else if data.Method == "approve" && data.FreezeId > 0 {
					msg, err := repository.Approve(*data)
					if err != nil {
						response := billing.Errors{ErrMessage: fmt.Sprintf("Failed on approve %s", err)}
						dataResponse := billing.ErrorResponse{Data: response}
						body, _ := json.Marshal(dataResponse)
						logrus.Print(string(body))
						return
					}
					logrus.Print(string(msg))
					d.Ack(false)
				} else {
					response := billing.Errors{ErrMessage: fmt.Sprintf("Invalid Request")}
					dataResponse := billing.ErrorResponse{Data: response}
					body, _ := json.Marshal(dataResponse)
					logrus.Print(string(body))
					d.Ack(false)
				}
			}
			logrus.Print(" [*] Waiting for messages. To exit press CTRL+C")
		}
	}()

	logrus.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
