package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	problemId    int `json:"problemId"`
	submissionId int `json:"submissionId"`
	codeId       int `json:"codeId"`
	timeLimit    int `json:"timeLimit"`
	memoryLimit  int `json:"memoryLimit"`
}

func initMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Printf("failed to dial rabbitmq: %v", err)
		return nil, err
	}
	return conn, nil
}

func openHelloChannel(conn *amqp.Connection, consumerID int) (<-chan amqp.Delivery, *amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("consumer %d: failed to open channel: %v", consumerID, err)
		return nil, nil, err
	}

	_, err = ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Printf("consumer %d: failed to declare queue: %v", consumerID, err)
		ch.Close()
		return nil, nil, err
	}

	msgs, err := ch.Consume(
		"hello", // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Printf("consumer %d: failed to register consumer: %v", consumerID, err)
		ch.Close()
		return nil, nil, err
	}

	log.Printf("consumer %d: started for topic hello", consumerID)
	return msgs, ch, nil
}
