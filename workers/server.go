package workers

import (
	billing "billingApp"
	"billingApp/pkg/repository"

	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

//Server receive requests and send responses, routing requests to functions
func Server() {
	//conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/") // to use on remote host
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") //to use on localhost
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	FailOnError(err, "Failed to set QoS")

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

	//check method of request and routing to match function, publish response
	go func() {
		for d := range msgs {
			logrus.Printf(" [*] %s", d.Body)
			data := &billing.Worker{}
			err := json.Unmarshal(d.Body, data)
			data.Amount = FloatToInt(data.Amount)
			data.FreezedAmount = FloatToInt(data.FreezedAmount)
			if err != nil {
				response := billing.Errors{ErrMessage: fmt.Sprintf("Error decoding JSON %s", err)}
				dataResponse := billing.ErrorResponse{Data: response}
				body, _ := json.Marshal(dataResponse)
				err = ch.Publish(
					"",
					d.ReplyTo,
					false,
					false,
					amqp.Publishing{
						ContentType:   "application/json",
						CorrelationId: d.CorrelationId,
						Body:          body,
					})
				d.Ack(false)

			} else {
				if data.Method == "decrease" && data.BalanceId > 0 && data.Amount > 0 {
					msg, err := repository.Decrease(*data)
					if err != nil {
						body := repository.ErrorResponse("Failed on decrease %s", err)
						err = ch.Publish(
							"",
							d.ReplyTo,
							false,
							false,
							amqp.Publishing{
								ContentType:   "application/json",
								CorrelationId: d.CorrelationId,
								Body:          body,
							})
						d.Ack(false)
						return
					}
					err = ch.Publish(
						"",
						d.ReplyTo,
						false,
						false,
						amqp.Publishing{
							ContentType:   "application/json",
							CorrelationId: d.CorrelationId,
							Body:          msg,
						})
					d.Ack(false)
				} else if data.Method == "increase" && data.BalanceId > 0 && data.Amount > 0 {
					msg, err := repository.Increase(*data)
					if err != nil {
						body := repository.ErrorResponse("Failed on increase %s", err)
						err = ch.Publish(
							"",
							d.ReplyTo,
							false,
							false,
							amqp.Publishing{
								ContentType:   "application/json",
								CorrelationId: d.CorrelationId,
								Body:          body,
							})
						d.Ack(false)
						return
					}
					err = ch.Publish(
						"",
						d.ReplyTo,
						false,
						false,
						amqp.Publishing{
							ContentType:   "application/json",
							CorrelationId: d.CorrelationId,
							Body:          msg,
						})
					d.Ack(false)

				} else if data.Method == "send" && data.BalanceId > 0 && data.Amount > 0 && data.Receiver > 0 {
					msg, err := repository.SendToOther(*data)
					if err != nil {
						body := repository.ErrorResponse("Failed on send to other %s", err)
						err = ch.Publish(
							"",
							d.ReplyTo,
							false,
							false,
							amqp.Publishing{
								ContentType:   "application/json",
								CorrelationId: d.CorrelationId,
								Body:          body,
							})
						d.Ack(false)
						return
					}
					err = ch.Publish(
						"",
						d.ReplyTo,
						false,
						false,
						amqp.Publishing{
							ContentType:   "application/json",
							CorrelationId: d.CorrelationId,
							Body:          msg,
						})
					d.Ack(false)
				} else if data.Method == "freeze" && data.BalanceId > 0 && data.FreezedAmount > 0 {
					msg, err := repository.FreezeAmount(*data)
					if err != nil {
						body := repository.ErrorResponse("Failed on freeze amount %s", err)
						err = ch.Publish(
							"",
							d.ReplyTo,
							false,
							false,
							amqp.Publishing{
								ContentType:   "application/json",
								CorrelationId: d.CorrelationId,
								Body:          body,
							})
						d.Ack(false)
						return
					}
					err = ch.Publish(
						"",
						d.ReplyTo,
						false,
						false,
						amqp.Publishing{
							ContentType:   "application/json",
							CorrelationId: d.CorrelationId,
							Body:          msg,
						})
					d.Ack(false)

				} else if data.Method == "approve" && data.FreezeId > 0 {
					msg, err := repository.Approve(*data)
					if err != nil {
						body := repository.ErrorResponse("Failed on approve %s", err)
						err = ch.Publish(
							"",
							d.ReplyTo,
							false,
							false,
							amqp.Publishing{
								ContentType:   "application/json",
								CorrelationId: d.CorrelationId,
								Body:          body,
							})
						d.Ack(false)
						return
					}
					err = ch.Publish(
						"",
						d.ReplyTo,
						false,
						false,
						amqp.Publishing{
							ContentType:   "application/json",
							CorrelationId: d.CorrelationId,
							Body:          msg,
						})
					d.Ack(false)
				} else {
					body := repository.ErrorResponse("Invalid Request", err)
					err = ch.Publish(
						"",
						d.ReplyTo,
						false,
						false,
						amqp.Publishing{
							ContentType:   "application/json",
							CorrelationId: d.CorrelationId,
							Body:          body,
						})
					d.Ack(false)
				}
			}

		}

	}()
	logrus.Printf(" [*] Awaiting RPC requests")
	<-forever
}
